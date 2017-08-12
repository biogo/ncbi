![bíogo](https://raw.githubusercontent.com/biogo/biogo/master/biogo.png)

# NCBI

[![Build Status](https://travis-ci.org/biogo/ncbi.svg?branch=master)](https://travis-ci.org/biogo/ncbi) [![GoDoc](https://godoc.org/github.com/biogo/ncbi?status.svg)](https://godoc.org/github.com/biogo/ncbi)

## Installation

The NCBI package requires a functioning [Go compiler installation](https://golang.org/doc/install).

        $ go get github.com/biogo/ncbi/...

## Overview

ncbi provides API interfaces to NCBI services.

* [Entrez Utility Programs](https://www.ncbi.nlm.nih.gov/books/NBK25501/)

* [BLAST](https://ncbi.github.io/blast-cloud/)

## Citing ##

If you use bíogo/ncbi, please cite Kortschak and Adelson "bíogo/ncbi: interfaces to NCBI services for the Go language", doi:[10.21105/joss.00234](http://dx.doi.org/10.21105/joss.00234), and Kortschak and Adelson "bíogo: a simple high-performance bioinformatics toolkit for the Go language", doi:[10.1101/005033](http://biorxiv.org/content/early/2014/05/12/005033).

## Example usage

### Entrez

This is a simple illustration of using the Entrez Utility Programs to retrieve a large set of sequences to a file.
The complete code is [here](paper/examples/fetch/fetch.go).

```
package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"github.com/biogo/ncbi"
	"github.com/biogo/ncbi/entrez"
)

const (
	db   = "protein"
	tool = "entrez.example"
)

var (
	clQuery = flag.String("query", "", "query specifies the search query for record retrieval (required).")
	rettype = flag.String("rettype", "fasta", "rettype specifies the format of the returned data.")
	retmax  = flag.Int("retmax", 500, "retmax specifies the number of records to be retrieved per request.")
	out     = flag.String("out", "", "out specifies destination of the returned data (default to stdout).")
	email   = flag.String("email", "", "email specifies the email address to be sent to the server (required).")
	retries = flag.Int("retry", 5, "retry specifies the number of attempts to retrieve the data.")
	help    = flag.Bool("help", false, "help prints this message.")
)

func main() {
	ncbi.SetTimeout(0)

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *email == "" || *clQuery == "" {
		flag.Usage()
		os.Exit(1)
	}

	h := entrez.History{}
	s, err := entrez.DoSearch(db, *clQuery, nil, &h, tool, *email)
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("will retrieve %d records.\n", s.Count)

	var of *os.File
	if *out == "" {
		of = os.Stdout
	} else {
		of, err = os.Create(*out)
		if err != nil {
			log.Printf("error: %v\n", err)
			os.Exit(1)
		}
		defer of.Close()
	}

	var (
		buf   = &bytes.Buffer{}
		p     = &entrez.Parameters{RetMax: *retmax, RetType: *rettype, RetMode: "text"}
		bn, n int64
	)
	for p.RetStart = 0; p.RetStart < s.Count; p.RetStart += p.RetMax {
		log.Printf("attempting to retrieve %d records starting from %d with %d retries.\n", p.RetMax, p.RetStart, *retries)
		var t int
		for t = 0; t < *retries; t++ {
			buf.Reset()
			var (
				r   io.ReadCloser
				_bn int64
			)
			r, err = entrez.Fetch(db, p, tool, *email, &h)
			if err != nil {
				if r != nil {
					r.Close()
				}
				log.Printf("failed to retrieve on attempt %d... error: %v ... retrying.\n", t, err)
				continue
			}
			_bn, err = io.Copy(buf, r)
			bn += _bn
			r.Close()
			if err == nil {
				break
			}
			log.Printf("failed to buffer on attempt %d... error: %v ... retrying.\n", t, err)
		}
		if err != nil {
			os.Exit(1)
		}

		log.Printf("retrieved records with %d retries... writing out.\n", t)
		_n, err := io.Copy(of, buf)
		n += _n
		if err != nil {
			log.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
	if bn != n {
		log.Printf("writethrough mismatch: %d != %d\n", bn, n)
	}
}
```

### BLAST

The following example provides a simple function used to perform a BLAST search from within a larger program.
A complete example is available [here](paper/examples/blast/blast.go).

```
// tool is required by the BLAST server.
const tool  = "blast.example"

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
```

## Getting help

Help or similar requests are preferred on the biogo-user Google Group.

https://groups.google.com/forum/#!forum/biogo-user

## Contributing

If you find any bugs, feel free to file an issue on the github issue tracker.
Pull requests are welcome, though if they involve changes to API or addition of features, please first open a discussion at the biogo-dev Google Group.

https://groups.google.com/forum/#!forum/biogo-dev

## Library Structure and Coding Style

The coding style should be aligned with normal Go idioms as represented in the
Go core libraries.

## Copyright and License

Copyright ©2011-2013 The bíogo Authors except where otherwise noted. All rights
reserved. Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.

The bíogo logo is derived from Bitstream Charter, Copyright ©1989-1992
Bitstream Inc., Cambridge, MA.

BITSTREAM CHARTER is a registered trademark of Bitstream Inc.
