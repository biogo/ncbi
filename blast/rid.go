// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/biogo/ncbi"

	"golang.org/x/net/html"
)

// Blast server usage policy requires that users not poll for any single RID more often than once
// a minute. The blast package Get requests honour this policy.
const RidPollLimit = 60 * time.Second

// Rid implements RID recovery and waiting functions associated with Blast Put and Get requests.
type Rid struct {
	rid   string
	rtoe  time.Time
	delay <-chan time.Time
	limit *ncbi.Limiter
}

// NewRid returns a Rid with the given request ID string. The returned Rid has a
// zero RTOE but is subject to the usage policy poll limiter. It is intended to
// be used to retrieve results from queries submitted without a call to Put.
func NewRid(rid string) *Rid {
	delay := make(chan time.Time)
	close(delay)
	return &Rid{
		rid:   rid,
		delay: delay,
		limit: ncbi.NewLimiter(RidPollLimit),
	}
}

var (
	messageIDBytes = []byte("Message ID")
	errorBytes     = []byte("Error:")
)

type ErrBadRequest string

func (e ErrBadRequest) Error() string { return string(e) }

func (rid *Rid) unmarshal(r io.Reader) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return z.Err()
		}
		if tt == html.TextToken {
			text := z.Text()
			if bytes.HasPrefix(text, messageIDBytes) && bytes.Contains(text, errorBytes) {
				rid.setElapsedDelay()
				return ErrBadRequest(text)
			}
		}
		if tt != html.CommentToken {
			continue
		}
		d := z.Token().Data
		if !strings.Contains(d, "QBlastInfoBegin") {
			continue
		}
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
			rid.setElapsedDelay()
			return ErrMissingRid
		}
		rid.limit = ncbi.NewLimiter(RidPollLimit)
		return nil
	}
}

// String returns the string representation of the Rid.
func (r *Rid) String() string {
	if r == nil {
		return "<nil>"
	}
	return r.rid
}

// TimeOfExecution returns the expected time until the request can be satisfied.
func (r *Rid) TimeOfExecution() time.Duration {
	now := time.Now()
	if now.Before(r.rtoe) {
		return r.rtoe.Sub(now)
	}
	return 0
}

// setElapsedDelay sets r's delay to have already elapsed.
func (r *Rid) setElapsedDelay() {
	delay := make(chan time.Time)
	close(delay)
	r.delay = delay
}

// Ready returns a time.Time chan that will send when the estimated time for the
// Put request to be satisfied has elapsed. If the request has failed the channel
// is returned closed.
func (r *Rid) Ready() <-chan time.Time {
	defer func() {
		delay := make(chan time.Time)
		close(delay)
		r.delay = delay
	}()
	return r.delay
}

// SearchInfo holds search status information.
type SearchInfo struct {
	*Rid
	Status   string
	HaveHits bool
}

func (s *SearchInfo) String() string {
	return fmt.Sprintf("RID:%s Status:%s Hits:%v", s.Rid, s.Status, s.HaveHits)
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

		// The following block pulls out the following element and makes the delay equal to X seconds.
		//
		//  <p class="WAITING">This page will be automatically updated in <b>X</b> seconds</p>
		//
		// If the integer parse fails it means NCBI have changed the html returned by SearchInfo so we
		// silently fall back to the policy wait time.
		if tt == html.StartTagToken {
			attr := z.Token().Attr
			if len(attr) > 0 && attr[0].Val == "WAITING" {
				for i := 0; i < 3; i++ {
					z.Next()
				}
				data := z.Token().Data
				rt, err := strconv.ParseInt(data, 10, 64)
				if err != nil {
					continue
				}
				secs := time.Duration(rt) * time.Second
				s.Rid.delay = time.After(secs)
				s.Rid.rtoe = time.Now().Add(secs)
				continue
			}
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
