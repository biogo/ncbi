// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast_test

import (
	"code.google.com/p/biogo.ncbi/blast"
	"errors"
	"fmt"
)

var (
	query     = "X14032.1"
	tool      = "example-blast"
	email     = "foo@bar.org"
	retries   = 3
	putParams = &blast.PutParameters{Program: "nblast", Database: "nr"}
	getParams = &blast.GetParameters{}
)

func Example_1() (*blast.Output, error) {
	// Put the query request on the BLAST server.
	r, err := blast.Put(query, putParams, tool, email)
	if err != nil {
		return nil, err
	}

	var o *blast.Output
	// Try to retrieve results up to maximum number of retires times.
	for k := 0; k < retries; k++ {
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
		o, err = r.GetOutput(getParams, tool, email)
		if err == nil {
			break
		}
	}

	return o, err
}
