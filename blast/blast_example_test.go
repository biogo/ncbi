// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast_test

import (
	"errors"
	"fmt"

	"github.com/biogo/ncbi/blast"
)

const (
	tool  = "blast-example"
	email = "foo@bar.org"
)

// Example submits a query to the BLAST server, waits for the server's estimated time of
// execution and retrieves the search status. If the search is ready the results are then
// retrieved and returned. If errors are returned during data retrieval from the server,
// retrieval is retried with up to retry attempts; all server requests honour the request
// frequency policy specified in the BLAST usage guidelines.
func Example(query string, retry int, pp *blast.PutParameters, gp *blast.GetParameters) (*blast.Output, error) {
	// Put the query request to the BLAST server.
	r, err := blast.Put(query, pp, tool, email)
	if err != nil {
		return nil, err
	}

	var o *blast.Output
	for k := 0; k < retry; k++ {
		// Wait for RTOE to elapse and get search status.
		var s *blast.SearchInfo
		s, err = r.SearchInfo(tool, email)
		if err != nil {
			return nil, err
		}

		// Output search status.
		fmt.Println(s)

		switch s.Status {
		case "WAITING":
			continue
		case "FAILED":
			return nil, fmt.Errorf("search: %s failed", r)
		case "UNKNOWN":
			return nil, fmt.Errorf("search: %s expired", r)
		case "READY":
			if !s.HaveHits {
				return nil, fmt.Errorf("search: %s no hits", r)
			}
		default:
			return nil, errors.New("unknown error")
		}

		// We have hits, so get the BLAST output.
		o, err = r.GetOutput(gp, tool, email)
		if err == nil {
			return o, err
		}
	}

	return nil, fmt.Errorf("%s exceeded retries", r)
}
