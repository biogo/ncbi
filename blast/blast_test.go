// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	check "launchpad.net/gocheck"
	"testing"
	"time"
)

const tool = "biogo.ncbi/blast-testsuite"

// Helpers
func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
func stringPtr(s string) *string  { return &s }

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestLimiter(c *check.C) {
	var count int
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				Limit.Wait()
				count++
			}
		}()
	}
	time.Sleep(9 * time.Second)
	c.Check(count, check.Equals, 3)
}

var net = flag.String("net", "", "Runs tests involving network connections if given an email address.")

func (s *S) TestBlast(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	for i, t := range []struct {
		query     string
		putParams *PutParameters
		getParams *GetParameters
		expect    *Output
	}{
		{
			query:     "X14032.1",
			putParams: &PutParameters{Program: "blastn", Database: "nr"},
			getParams: nil,
			expect: &Output{
				Program:   "blastn",
				Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
				Database:  "nr",
				QueryId:   "gi|322|emb|X14032.1|",
				Parameters: Parameters{
					MatrixName:  nil,
					Expect:      10,
					Include:     nil,
					Match:       intPtr(2),
					Mismatch:    intPtr(-3),
					GapOpen:     5,
					GapExtend:   2,
					Filter:      stringPtr("L;m;"),
					PhiPattern:  nil,
					EntrezQuery: nil,
				},
				Iterations: []Iteration{
					{
						N:        1,
						QueryId:  stringPtr("gi|322|emb|X14032.1|"),
						QueryDef: stringPtr("Bovine mRNA for EDGF II (acidic eye-derived growth factor II)"),
						QueryLen: intPtr(668),
						Hits: []Hit{
							{
								N:         1,
								Id:        "gi|322|emb|X14032.1|",
								Def:       "Bovine mRNA for EDGF II (acidic eye-derived growth factor II)",
								Accession: "X14032",
								Len:       668,
								Hsps: []Hsp{
									{
										N:              1,
										BitScore:       1205.94,
										Score:          1336,
										EValue:         0,
										QueryFrom:      1,
										QueryTo:        668,
										HitFrom:        1,
										HitTo:          668,
										PhiPatternFrom: nil,
										PhiPatternTo:   nil,
										QueryFrame:     intPtr(1),
										HitFrame:       intPtr(1),
										HspIdentity:    intPtr(668),
										HspPositive:    intPtr(668),
										HspGaps:        intPtr(0),
										AlignLen:       intPtr(668),
										Density:        nil,
										QuerySeq:       stringPtr("GGATCCTCTTTCCCTTCTACTGGAGAGGAAAAGCCCTCAGCCTGCAAGCTGTTCAGCCTTGAAACAGCCACAACCAGCAGCTGCTGAGCCATGGCTGAAGGAGAAACCACGACCTTCACGGCCCTGACTGAGAAGTTTAACCTGCCTCTAGGCAATTACAAGAAGCCCAAGCTCCTCTACTGCAGCAACGGGGGCTACTTCCTGAGAATCCTCCCAGATGGCACAGTGGATGGGACGAAGGACAGGAGCGACCAGCACATTCAGCTGCAGCTCTGTGCGGAAAGCATAGGGGAGGTGTATATTAAGAGTACGGAGACTGGCCAGTTCTTGGCCATGGACACCGACGGGCTTTTGTACGGCTCACAGACACCCAATGAGGAATGTTTGTTCCTGGAAAGGTTGGAGGAAAACCATTACAACACCTACATATCCAAGAAGCATGCAGAGAAGCATTGGTTCGTTGGTCTCAAGAAGAACGGAAGGTCTAAACTCGGTCCTCGGACTCACTTCGGCCAGAAAGCCATCTTGTTTCTCCCCCTGCCAGTCTCCTCTGATTAAAGAAATCTGTTGTGGGTGCTGAGCCACTCCAGAGGAATCTGAAGGGGTCCTCACCTGGCTGACCCCAGATTGTACCCTTTACCATTGGCCGTGCTAACCCCTGGCCCACA"),
										SubjectSeq:     stringPtr("GGATCCTCTTTCCCTTCTACTGGAGAGGAAAAGCCCTCAGCCTGCAAGCTGTTCAGCCTTGAAACAGCCACAACCAGCAGCTGCTGAGCCATGGCTGAAGGAGAAACCACGACCTTCACGGCCCTGACTGAGAAGTTTAACCTGCCTCTAGGCAATTACAAGAAGCCCAAGCTCCTCTACTGCAGCAACGGGGGCTACTTCCTGAGAATCCTCCCAGATGGCACAGTGGATGGGACGAAGGACAGGAGCGACCAGCACATTCAGCTGCAGCTCTGTGCGGAAAGCATAGGGGAGGTGTATATTAAGAGTACGGAGACTGGCCAGTTCTTGGCCATGGACACCGACGGGCTTTTGTACGGCTCACAGACACCCAATGAGGAATGTTTGTTCCTGGAAAGGTTGGAGGAAAACCATTACAACACCTACATATCCAAGAAGCATGCAGAGAAGCATTGGTTCGTTGGTCTCAAGAAGAACGGAAGGTCTAAACTCGGTCCTCGGACTCACTTCGGCCAGAAAGCCATCTTGTTTCTCCCCCTGCCAGTCTCCTCTGATTAAAGAAATCTGTTGTGGGTGCTGAGCCACTCCAGAGGAATCTGAAGGGGTCCTCACCTGGCTGACCCCAGATTGTACCCTTTACCATTGGCCGTGCTAACCCCTGGCCCACA"),
										FormatMidline:  stringPtr("||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||"),
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		r, err := Put(t.query, t.putParams, tool, *net)
		c.Assert(err, check.Equals, nil, check.Commentf("Test %d", i))
		<-r.Ready()
		o, err := r.GetOutput(t.getParams, tool, *net)
		c.Assert(err, check.Equals, nil, check.Commentf("Test %d", i))
		c.Check(o.Program, check.Equals, t.expect.Program)
		c.Check(o.Reference, check.Equals, t.expect.Reference)
		c.Check(o.Database, check.Equals, t.expect.Database)
		c.Check(o.QueryId, check.Equals, t.expect.QueryId)
		c.Check(o.Parameters, check.DeepEquals, t.expect.Parameters)
		c.Assert(len(o.Iterations) > 0, check.Equals, true)
		c.Assert(len(o.Iterations[0].Hits) > 0, check.Equals, true)
		c.Check(o.Iterations[0].Hits[0], check.DeepEquals, t.expect.Iterations[0].Hits[0])
		spew.Dump(o)
	}
}
