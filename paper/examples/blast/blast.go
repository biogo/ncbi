// Copyright ©2017 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blast is a simple command line remote BLAST database query program.
// The program will print the entire returned data structure resulting
// from the BLAST search query provided.
//
// blast requires that github.com/kortschak/utter is installed.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/biogo/ncbi/blast"
	"github.com/kortschak/utter"
)

// tool is required by the BLAST server.
const tool = "blast.example"

var (
	clQuery = flag.String("query", "", "query specifies the search query file (required).")
	prog    = flag.String("prog", "blastn", "prog specifies the blast program to run.")
	db      = flag.String("db", "nr", "db specifies the blast database to search.")
	out     = flag.String("out", "", "out specifies destination of the returned data (default to stdout).")
	email   = flag.String("email", "", "email specifies the email address to be sent to the server (required).")
	retries = flag.Int("retry", 5, "retry specifies the number of attempts to retrieve the data.")
	dump    = flag.Bool("dump", false, "dump specifies that the complete data structure should be dumped.")
	help    = flag.Bool("help", false, "help prints this message.")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *email == "" || *clQuery == "" {
		flag.Usage()
		os.Exit(1)
	}

	var of *os.File
	var err error
	if *out == "" {
		of = os.Stdout
	} else {
		of, err = os.Create(*out)
		if err != nil {
			log.Fatalf("failed to create output file: %v\n", err)
		}
		defer of.Close()
	}

	b, err := ioutil.ReadFile(*clQuery)
	if err != nil {
		log.Fatalf("failed to read query %q: %v", *clQuery, err)
	}

	pp := blast.PutParameters{
		Program:  *prog,
		Database: *db,
	}

	o, err := BLAST(string(b), *retries, &pp, nil, *email)
	if err != nil {
		log.Fatalf("failed to run remote BLAST query: %v", err)
	}

	if *dump {
		utter.Config.BytesWidth = 8
		utter.Config.ElideType = true
		utter.Fdump(of, o)
		return
	}
	for _, it := range o.Iterations {
		for _, hit := range it.Hits {
			fmt.Fprintf(of, "%s\t%s\t%d\n", hit.Id, hit.Def, hit.Len)
			for _, hsp := range hit.Hsps {
				fmt.Fprintf(of, "\tE value = %v\n", hsp.EValue)
				fmt.Fprintf(of, "\t\t%v-%v\t%s\n", hsp.QueryFrom, hsp.QueryTo, hsp.QuerySeq)
				fmt.Fprintf(of, "\t\t\t%s\n", hsp.FormatMidline)
				fmt.Fprintf(of, "\t\t%v-%v\t%s\n", hsp.HitFrom, hsp.HitTo, hsp.SubjectSeq)
			}
		}
	}
}

// BLAST submits a query to the BLAST server, waits for the server's estimated time of
// execution and retrieves the search status. If the search is ready the results are then
// retrieved and returned. If errors are returned during data retrieval from the server,
// retrieval is retried with up to retry attempts; all server requests honour the request
// frequency policy specified in the BLAST usage guidelines.
func BLAST(query string, retry int, pp *blast.PutParameters, gp *blast.GetParameters, email string) (*blast.Output, error) {
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
