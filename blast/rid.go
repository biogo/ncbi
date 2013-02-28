// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

import (
	"code.google.com/p/biogo.ncbi"
	"code.google.com/p/biogo.ncbi/html"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Blast server usage policy requires that users not poll for any single RID more often than once
// a minute. blast package Get requests honour this policy.
const RidPollLimit = 60 * time.Second

// Rid implements RID recovery and waiting functions associated with Blast Put and Get requests.
type Rid struct {
	rid   string
	rtoe  time.Time
	delay <-chan time.Time
	limit *ncbi.Limiter
}

func (rid *Rid) unmarshal(r io.Reader) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return z.Err()
		}
		if tt == html.CommentToken {
			d := z.Token().Data
			if strings.Contains(d, "QBlastInfoBegin") {
				for _, l := range strings.Split(d, "\n") {
					l = strings.TrimSpace(l)
					kv := strings.Split(l, " = ")
					if len(kv) != 2 {
						continue
					}
					switch kv[0] {
					case "RID":
						rid.rid = kv[1]
					case "RTOE":
						rt, err := strconv.ParseInt(kv[1], 10, 64)
						if err != nil {
							return err
						}
						secs := time.Duration(rt) * time.Second
						rid.delay = time.After(secs)
						rid.rtoe = time.Now().Add(secs)
					}
				}
				if rid.rid == "" || rid.delay == nil {
					delay := make(chan time.Time)
					close(delay)
					rid.delay = delay
					return ErrMissingRid
				}
				rid.limit = ncbi.NewLimiter(RidPollLimit)
				return nil
			}
		}
	}

	panic("cannot reach")
}

// String returns the string representation of the Rid.
func (r *Rid) String() string { return r.rid }

// TimeOfExecution returns the expected time until the request can be satisfied.
func (r *Rid) TimeOfExecution() time.Duration {
	now := time.Now()
	if now.Before(r.rtoe) {
		return r.rtoe.Sub(now)
	}
	return 0
}

// Ready returns a time.Time chan that will send when the estimated time for the
// Put request to be satisfied has elapsed. If the request has failed the channel
// is returned closed.
func (r *Rid) Ready() <-chan time.Time {
	return r.delay
}

// SearchInfo holds search status information.
type SearchInfo struct {
	*Rid
	Status   string
	HaveHits bool
}

func (s *SearchInfo) String() string {
	return fmt.Sprintf("%s Status:%s Hits:%v", s.Rid, s.Status, s.HaveHits)
}

func (s *SearchInfo) unmarshal(r io.Reader) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			err := z.Err()
			if err == io.EOF {
				break
			}
			return err
		}
		if tt == html.CommentToken {
			d := z.Token().Data
			if strings.Contains(d, "QBlastInfoBegin") {
				for _, l := range strings.Split(d, "\n") {
					l = strings.TrimSpace(l)
					kv := strings.Split(l, "=")
					if len(kv) != 2 {
						continue
					}
					switch kv[0] {
					case "Status":
						s.Status = kv[1]
					case "ThereAreHits":
						s.HaveHits = true
					}
				}
			}
		}
	}

	if s.Status == "" {
		return ErrMissingStatus
	}
	return nil
}
