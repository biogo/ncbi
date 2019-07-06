// Copyright ©2017 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fetch is a command line entrez database query program.
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
	tool = "entrez.example"
)

var (
	clQuery = flag.String("query", "", "query specifies the search query for record retrieval (required).")
	db      = flag.String("db", "protein", "db specifies the database to search")
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

	h := entrez.History{}
	s, err := entrez.DoSearch(*db, *clQuery, nil, &h, tool, *email)
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("will retrieve %d records.\n", s.Count)

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
			r, err = entrez.Fetch(*db, p, tool, *email, &h)
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
