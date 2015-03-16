// Copyright ©2014 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/biogo/ncbi/blast"
	"github.com/biogo/ncbi/blast/graphic"

	"code.google.com/p/plotinum/vg"
	"code.google.com/p/plotinum/vg/vgsvg"

	"gopkg.in/check.v1"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func intPtr(i int) *int          { return &i }
func stringPtr(s string) *string { return &s }

var testOutputs = []struct {
	in blast.Output
	cf func(w, h vg.Length) vg.Canvas

	configure bool
	legend    bool
	aligns    bool
	depths    bool

	expect string
}{
	{
		in: blast.Output{
			Program:   "blastn",
			Version:   "BLASTN 2.2.27+",
			Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
			Database:  "nr",
			QueryId:   "33421",
			QueryDef:  "No definition line",
			QueryLen:  32,
			QuerSeq:   nil,
			Parameters: blast.Parameters{
				MatrixName:  nil,
				Expect:      1000,
				Include:     nil,
				Match:       intPtr(1),
				Mismatch:    intPtr(-3),
				GapOpen:     5,
				GapExtend:   2,
				Filter:      stringPtr("F"),
				PhiPattern:  nil,
				EntrezQuery: nil,
			},
			Iterations: []blast.Iteration{
				{
					N:        1,
					QueryId:  stringPtr("33421"),
					QueryDef: stringPtr("No definition line"),
					QueryLen: intPtr(32),
					Hits: []blast.Hit{
						{
							N:         1,
							Id:        "gi|388525227|gb|CP003531.1|",
							Def:       "Thermogladius cellulolyticus 1633, complete genome",
							Accession: "CP003531",
							Len:       1356318,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       38.1576,
									Score:          19,
									EValue:         1.72292,
									QueryFrom:      7,
									QueryTo:        29,
									HitFrom:        1187458,
									HitTo:          1187436,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(22),
									HspPositive:    intPtr(22),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(23),
									Density:        nil,
									QuerySeq:       []byte("TGTCGAACTATACGACGAGCACT"),
									SubjectSeq:     []byte("TGTCGAGCTATACGACGAGCACT"),
									FormatMidline:  []byte("|||||| ||||||||||||||||"),
								},
							},
						},
						{
							N:         2,
							Id:        "gi|354799811|gb|JN964312.1|",
							Def:       "Mus musculus targeted non-conditional, lacZ-tagged mutant allele Morn4:tm1e(KOMP)Wtsi; transgenic",
							Accession: "JN964312",
							Len:       40247,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       36.1753,
									Score:          18,
									EValue:         6.80792,
									QueryFrom:      1,
									QueryTo:        18,
									HitFrom:        26580,
									HitTo:          26597,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(18),
									HspPositive:    intPtr(18),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(18),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTATA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||||"),
								},
							},
						},
						{
							N:         24,
							Id:        "gi|356871506|emb|FO082053.1|",
							Def:       "Pichia sorbitophila strain CBS 7064 chromosome G complete sequence",
							Accession: "FO082053",
							Len:       1423303,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      13,
									QueryTo:        28,
									HitFrom:        61730,
									HitTo:          61745,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACTATACGACGAGCAC"),
									SubjectSeq:     []byte("ACTATACGACGAGCAC"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         25,
							Id:        "gi|353230524|emb|HE601625.1|",
							Def:       "Schistosoma mansoni strain Puerto Rico chromosome 2, complete genome",
							Accession: "HE601625",
							Len:       34464480,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      1,
									QueryTo:        16,
									HitFrom:        28173004,
									HitTo:          28173019,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
								{
									N:              2,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      3,
									QueryTo:        18,
									HitFrom:        28996063,
									HitTo:          28996048,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("AGAATGTCGAACTATA"),
									SubjectSeq:     []byte("AGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         26,
							Id:        "gi|341821300|emb|HE576794.1|",
							Def:       "Megasphaera elsdenii strain DSM 20460 draft genome",
							Accession: "HE576794",
							Len:       2474718,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      14,
									QueryTo:        29,
									HitFrom:        1741778,
									HitTo:          1741793,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("CTATACGACGAGCACT"),
									SubjectSeq:     []byte("CTATACGACGAGCACT"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
					},
					Statistics: &blast.Statistics{
						DbNum:    17331862,
						DbLen:    1419046940,
						HspLen:   0,
						EffSpace: 0,
						Kappa:    0.710603,
						Lambda:   1.37406,
						Entropy:  1.30725,
					},
					Message: nil,
				},
			},
			MegaStatistics: nil,
		},
		cf: func(w, h vg.Length) vg.Canvas { return vgsvg.New(w, h) },
		expect: "" +
			"<?xml version=\"1.0\"?>\n" +
			"<!-- Generated by SVGo and Plotinum VG -->\n" +
			"<svg width=\"11.111in\" height=\"3.5556in\"\n" +
			"\txmlns=\"http://www.w3.org/2000/svg\" \n" +
			"\txmlns:xlink=\"http://www.w3.org/1999/xlink\">\n" +
			"<g transform=\"scale(1, -1) translate(0, -320)\">\n" +
			"<path d=\"M0,0L1000,0L1000,320L0,320Z\" style=\"fill:#FFFFFF\" />\n" +
			"<text x=\"6.25\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:bold;font-style:normal;font-size:13pt\">33421</text>\n" +
			"<path d=\"M187.5,270L937.5,270\" style=\"fill:none;stroke:#000000;stroke-width:1.25\" />\n" +
			"<text x=\"187.5\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">1</text>\n" +
			"<text x=\"920.82\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">29</text>\n" +
			"<text x=\"837.5\" y=\"-302.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">% Identity</text>\n" +
			"<path d=\"M612.5,313.75L625,313.75L625,301.25L612.5,301.25Z\" style=\"fill:#808080\" />\n" +
			"<text x=\"612.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">20</text>\n" +
			"<path d=\"M637.5,313.75L650,313.75L650,301.25L637.5,301.25Z\" style=\"fill:#FF0000\" />\n" +
			"<text x=\"637.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">30</text>\n" +
			"<path d=\"M662.5,313.75L675,313.75L675,301.25L662.5,301.25Z\" style=\"fill:#FFC400\" />\n" +
			"<text x=\"662.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">40</text>\n" +
			"<path d=\"M687.5,313.75L700,313.75L700,301.25L687.5,301.25Z\" style=\"fill:#FFFF00\" />\n" +
			"<text x=\"687.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">50</text>\n" +
			"<path d=\"M712.5,313.75L725,313.75L725,301.25L712.5,301.25Z\" style=\"fill:#00FF00\" />\n" +
			"<text x=\"712.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">60</text>\n" +
			"<path d=\"M737.5,313.75L750,313.75L750,301.25L737.5,301.25Z\" style=\"fill:#00FFFF\" />\n" +
			"<text x=\"737.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">70</text>\n" +
			"<path d=\"M762.5,313.75L775,313.75L775,301.25L762.5,301.25Z\" style=\"fill:#0000FF\" />\n" +
			"<text x=\"762.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">80</text>\n" +
			"<path d=\"M787.5,313.75L800,313.75L800,301.25L787.5,301.25Z\" style=\"fill:#C400FF\" />\n" +
			"<text x=\"787.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">90</text>\n" +
			"<path d=\"M812.5,313.75L825,313.75L825,301.25L812.5,301.25Z\" />\n" +
			"<text x=\"812.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">100</text>\n" +
			"<text x=\"12.5\" y=\"-227.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|388525227|gb|CP00</text>\n" +
			"<path d=\"M348.21,235L937.5,235\" style=\"fill:none;stroke:#C400FF;stroke-width:3.75\" />\n" +
			"<text x=\"348.21\" y=\"-242.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">7</text>\n" +
			"<text x=\"927.77\" y=\"-242.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">29</text>\n" +
			"<text x=\"348.21\" y=\"-225\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1187458</text>\n" +
			"<text x=\"898.1\" y=\"-225\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1187436 -</text>\n" +
			"<text x=\"12.5\" y=\"-182.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|354799811|gb|JN96</text>\n" +
			"<path d=\"M187.5,190L642.86,190\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"187.5\" y=\"-197.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1</text>\n" +
			"<text x=\"633.13\" y=\"-197.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">18</text>\n" +
			"<text x=\"187.5\" y=\"-180\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">26580</text>\n" +
			"<text x=\"610.99\" y=\"-180\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">26597 +</text>\n" +
			"<text x=\"12.5\" y=\"-137.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|356871506|emb|FO0</text>\n" +
			"<path d=\"M508.93,145L910.71,145\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"508.93\" y=\"-152.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">13</text>\n" +
			"<text x=\"900.98\" y=\"-152.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28</text>\n" +
			"<text x=\"508.93\" y=\"-135\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">61730</text>\n" +
			"<text x=\"878.85\" y=\"-135\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">61745 +</text>\n" +
			"<text x=\"12.5\" y=\"-92.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|353230524|emb|HE6</text>\n" +
			"<path d=\"M187.5,100L589.29,100\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"187.5\" y=\"-107.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1</text>\n" +
			"<text x=\"579.56\" y=\"-107.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">16</text>\n" +
			"<text x=\"187.5\" y=\"-90\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28173004</text>\n" +
			"<text x=\"542.82\" y=\"-90\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28173019 +</text>\n" +
			"<path d=\"M241.07,80L642.86,80\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"241.07\" y=\"-87.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">3</text>\n" +
			"<text x=\"633.13\" y=\"-87.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">18</text>\n" +
			"<text x=\"241.07\" y=\"-70\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28996063</text>\n" +
			"<text x=\"598.59\" y=\"-70\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28996048 -</text>\n" +
			"<text x=\"12.5\" y=\"-27.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|341821300|emb|HE5</text>\n" +
			"<path d=\"M535.71,35L937.5,35\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"535.71\" y=\"-42.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">14</text>\n" +
			"<text x=\"927.77\" y=\"-42.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">29</text>\n" +
			"<text x=\"535.71\" y=\"-25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1741778</text>\n" +
			"<text x=\"895.9\" y=\"-25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1741793 +</text>\n" +
			"<text x=\"941.25\" y=\"-261.25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1/line</text>\n" +
			"<path d=\"M187.5,266.25L187.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M188.75,266.25L188.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M190,266.25L190,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M191.25,266.25L191.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M192.5,266.25L192.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M193.75,266.25L193.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M195,266.25L195,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M196.25,266.25L196.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M197.5,266.25L197.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M198.75,266.25L198.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M200,266.25L200,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M201.25,266.25L201.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M202.5,266.25L202.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M203.75,266.25L203.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M205,266.25L205,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M206.25,266.25L206.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M207.5,266.25L207.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M208.75,266.25L208.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M210,266.25L210,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M211.25,266.25L211.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M212.5,266.25L212.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M213.75,266.25L213.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M215,266.25L215,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M216.25,266.25L216.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M217.5,266.25L217.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M218.75,266.25L218.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M220,266.25L220,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M221.25,266.25L221.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M222.5,266.25L222.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M223.75,266.25L223.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M225,266.25L225,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M226.25,266.25L226.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M227.5,266.25L227.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M228.75,266.25L228.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M230,266.25L230,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M231.25,266.25L231.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M232.5,266.25L232.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M233.75,266.25L233.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M235,266.25L235,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M236.25,266.25L236.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M237.5,266.25L237.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M238.75,266.25L238.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M240,266.25L240,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M241.25,266.25L241.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M242.5,266.25L242.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M243.75,266.25L243.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M245,266.25L245,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M246.25,266.25L246.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M247.5,266.25L247.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M248.75,266.25L248.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M250,266.25L250,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M251.25,266.25L251.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M252.5,266.25L252.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M253.75,266.25L253.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M255,266.25L255,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M256.25,266.25L256.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M257.5,266.25L257.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M258.75,266.25L258.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M260,266.25L260,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M261.25,266.25L261.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M262.5,266.25L262.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M263.75,266.25L263.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M265,266.25L265,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M266.25,266.25L266.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M267.5,266.25L267.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M268.75,266.25L268.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M270,266.25L270,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M271.25,266.25L271.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M272.5,266.25L272.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M273.75,266.25L273.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M275,266.25L275,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M276.25,266.25L276.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M277.5,266.25L277.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M278.75,266.25L278.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M280,266.25L280,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M281.25,266.25L281.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M282.5,266.25L282.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M283.75,266.25L283.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M285,266.25L285,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M286.25,266.25L286.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M287.5,266.25L287.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M288.75,266.25L288.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M290,266.25L290,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M291.25,266.25L291.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M292.5,266.25L292.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M293.75,266.25L293.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M295,266.25L295,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M296.25,266.25L296.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M297.5,266.25L297.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M298.75,266.25L298.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M300,266.25L300,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M301.25,266.25L301.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M302.5,266.25L302.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M303.75,266.25L303.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M305,266.25L305,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M306.25,266.25L306.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M307.5,266.25L307.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M308.75,266.25L308.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M310,266.25L310,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M311.25,266.25L311.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M312.5,266.25L312.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M313.75,266.25L313.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M315,266.25L315,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M316.25,266.25L316.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M317.5,266.25L317.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M318.75,266.25L318.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M320,266.25L320,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M321.25,266.25L321.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M322.5,266.25L322.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M323.75,266.25L323.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M325,266.25L325,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M326.25,266.25L326.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M327.5,266.25L327.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M328.75,266.25L328.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M330,266.25L330,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M331.25,266.25L331.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M332.5,266.25L332.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M333.75,266.25L333.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M335,266.25L335,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M336.25,266.25L336.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M337.5,266.25L337.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M338.75,266.25L338.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M340,266.25L340,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M341.25,266.25L341.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M342.5,266.25L342.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M343.75,266.25L343.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M345,266.25L345,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M346.25,266.25L346.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M347.5,266.25L347.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M348.75,266.25L348.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M350,266.25L350,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M351.25,266.25L351.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M352.5,266.25L352.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M353.75,266.25L353.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M355,266.25L355,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M356.25,266.25L356.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M357.5,266.25L357.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M358.75,266.25L358.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M360,266.25L360,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M361.25,266.25L361.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M362.5,266.25L362.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M363.75,266.25L363.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M365,266.25L365,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M366.25,266.25L366.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M367.5,266.25L367.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M368.75,266.25L368.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M370,266.25L370,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M371.25,266.25L371.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M372.5,266.25L372.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M373.75,266.25L373.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M375,266.25L375,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M376.25,266.25L376.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M377.5,266.25L377.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M378.75,266.25L378.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M380,266.25L380,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M381.25,266.25L381.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M382.5,266.25L382.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M383.75,266.25L383.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M385,266.25L385,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M386.25,266.25L386.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M387.5,266.25L387.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M388.75,266.25L388.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M390,266.25L390,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M391.25,266.25L391.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M392.5,266.25L392.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M393.75,266.25L393.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M395,266.25L395,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M396.25,266.25L396.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M397.5,266.25L397.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M398.75,266.25L398.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M400,266.25L400,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M401.25,266.25L401.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M402.5,266.25L402.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M403.75,266.25L403.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M405,266.25L405,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M406.25,266.25L406.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M407.5,266.25L407.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M408.75,266.25L408.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M410,266.25L410,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M411.25,266.25L411.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M412.5,266.25L412.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M413.75,266.25L413.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M415,266.25L415,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M416.25,266.25L416.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M417.5,266.25L417.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M418.75,266.25L418.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M420,266.25L420,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M421.25,266.25L421.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M422.5,266.25L422.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M423.75,266.25L423.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M425,266.25L425,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M426.25,266.25L426.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M427.5,266.25L427.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M428.75,266.25L428.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M430,266.25L430,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M431.25,266.25L431.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M432.5,266.25L432.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M433.75,266.25L433.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M435,266.25L435,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M436.25,266.25L436.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M437.5,266.25L437.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M438.75,266.25L438.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M440,266.25L440,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M441.25,266.25L441.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M442.5,266.25L442.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M443.75,266.25L443.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M445,266.25L445,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M446.25,266.25L446.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M447.5,266.25L447.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M448.75,266.25L448.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M450,266.25L450,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M451.25,266.25L451.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M452.5,266.25L452.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M453.75,266.25L453.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M455,266.25L455,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M456.25,266.25L456.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M457.5,266.25L457.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M458.75,266.25L458.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M460,266.25L460,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M461.25,266.25L461.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M462.5,266.25L462.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M463.75,266.25L463.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M465,266.25L465,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M466.25,266.25L466.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M467.5,266.25L467.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M468.75,266.25L468.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M470,266.25L470,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M471.25,266.25L471.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M472.5,266.25L472.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M473.75,266.25L473.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M475,266.25L475,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M476.25,266.25L476.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M477.5,266.25L477.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M478.75,266.25L478.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M480,266.25L480,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M481.25,266.25L481.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M482.5,266.25L482.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M483.75,266.25L483.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M485,266.25L485,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M486.25,266.25L486.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M487.5,266.25L487.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M488.75,266.25L488.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M490,266.25L490,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M491.25,266.25L491.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M492.5,266.25L492.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M493.75,266.25L493.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M495,266.25L495,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M496.25,266.25L496.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M497.5,266.25L497.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M498.75,266.25L498.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M500,266.25L500,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M501.25,266.25L501.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M502.5,266.25L502.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M503.75,266.25L503.75,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M505,266.25L505,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M506.25,266.25L506.25,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M507.5,266.25L507.5,261.25\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M508.75,266.25L508.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M510,266.25L510,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M511.25,266.25L511.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M512.5,266.25L512.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M513.75,266.25L513.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M515,266.25L515,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M516.25,266.25L516.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M517.5,266.25L517.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M518.75,266.25L518.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M520,266.25L520,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M521.25,266.25L521.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M522.5,266.25L522.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M523.75,266.25L523.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M525,266.25L525,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M526.25,266.25L526.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M527.5,266.25L527.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M528.75,266.25L528.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M530,266.25L530,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M531.25,266.25L531.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M532.5,266.25L532.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M533.75,266.25L533.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M535,266.25L535,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M536.25,266.25L536.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M537.5,266.25L537.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M538.75,266.25L538.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M540,266.25L540,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M541.25,266.25L541.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M542.5,266.25L542.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M543.75,266.25L543.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M545,266.25L545,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M546.25,266.25L546.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M547.5,266.25L547.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M548.75,266.25L548.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M550,266.25L550,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M551.25,266.25L551.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M552.5,266.25L552.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M553.75,266.25L553.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M555,266.25L555,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M556.25,266.25L556.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M557.5,266.25L557.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M558.75,266.25L558.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M560,266.25L560,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M561.25,266.25L561.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M562.5,266.25L562.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M563.75,266.25L563.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M565,266.25L565,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M566.25,266.25L566.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M567.5,266.25L567.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M568.75,266.25L568.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M570,266.25L570,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M571.25,266.25L571.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M572.5,266.25L572.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M573.75,266.25L573.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M575,266.25L575,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M576.25,266.25L576.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M577.5,266.25L577.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M578.75,266.25L578.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M580,266.25L580,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M581.25,266.25L581.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M582.5,266.25L582.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M583.75,266.25L583.75,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M585,266.25L585,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M586.25,266.25L586.25,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M587.5,266.25L587.5,258.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M588.75,266.25L588.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M590,266.25L590,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M591.25,266.25L591.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M592.5,266.25L592.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M593.75,266.25L593.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M595,266.25L595,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M596.25,266.25L596.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M597.5,266.25L597.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M598.75,266.25L598.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M600,266.25L600,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M601.25,266.25L601.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M602.5,266.25L602.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M603.75,266.25L603.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M605,266.25L605,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M606.25,266.25L606.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M607.5,266.25L607.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M608.75,266.25L608.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M610,266.25L610,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M611.25,266.25L611.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M612.5,266.25L612.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M613.75,266.25L613.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M615,266.25L615,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M616.25,266.25L616.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M617.5,266.25L617.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M618.75,266.25L618.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M620,266.25L620,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M621.25,266.25L621.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M622.5,266.25L622.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M623.75,266.25L623.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M625,266.25L625,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M626.25,266.25L626.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M627.5,266.25L627.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M628.75,266.25L628.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M630,266.25L630,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M631.25,266.25L631.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M632.5,266.25L632.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M633.75,266.25L633.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M635,266.25L635,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M636.25,266.25L636.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M637.5,266.25L637.5,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M638.75,266.25L638.75,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M640,266.25L640,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M641.25,266.25L641.25,260\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M642.5,266.25L642.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M643.75,266.25L643.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M645,266.25L645,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M646.25,266.25L646.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M647.5,266.25L647.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M648.75,266.25L648.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M650,266.25L650,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M651.25,266.25L651.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M652.5,266.25L652.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M653.75,266.25L653.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M655,266.25L655,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M656.25,266.25L656.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M657.5,266.25L657.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M658.75,266.25L658.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M660,266.25L660,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M661.25,266.25L661.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M662.5,266.25L662.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M663.75,266.25L663.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M665,266.25L665,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M666.25,266.25L666.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M667.5,266.25L667.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M668.75,266.25L668.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M670,266.25L670,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M671.25,266.25L671.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M672.5,266.25L672.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M673.75,266.25L673.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M675,266.25L675,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M676.25,266.25L676.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M677.5,266.25L677.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M678.75,266.25L678.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M680,266.25L680,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M681.25,266.25L681.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M682.5,266.25L682.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M683.75,266.25L683.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M685,266.25L685,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M686.25,266.25L686.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M687.5,266.25L687.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M688.75,266.25L688.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M690,266.25L690,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M691.25,266.25L691.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M692.5,266.25L692.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M693.75,266.25L693.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M695,266.25L695,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M696.25,266.25L696.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M697.5,266.25L697.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M698.75,266.25L698.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M700,266.25L700,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M701.25,266.25L701.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M702.5,266.25L702.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M703.75,266.25L703.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M705,266.25L705,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M706.25,266.25L706.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M707.5,266.25L707.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M708.75,266.25L708.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M710,266.25L710,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M711.25,266.25L711.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M712.5,266.25L712.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M713.75,266.25L713.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M715,266.25L715,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M716.25,266.25L716.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M717.5,266.25L717.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M718.75,266.25L718.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M720,266.25L720,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M721.25,266.25L721.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M722.5,266.25L722.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M723.75,266.25L723.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M725,266.25L725,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M726.25,266.25L726.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M727.5,266.25L727.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M728.75,266.25L728.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M730,266.25L730,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M731.25,266.25L731.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M732.5,266.25L732.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M733.75,266.25L733.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M735,266.25L735,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M736.25,266.25L736.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M737.5,266.25L737.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M738.75,266.25L738.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M740,266.25L740,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M741.25,266.25L741.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M742.5,266.25L742.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M743.75,266.25L743.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M745,266.25L745,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M746.25,266.25L746.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M747.5,266.25L747.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M748.75,266.25L748.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M750,266.25L750,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M751.25,266.25L751.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M752.5,266.25L752.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M753.75,266.25L753.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M755,266.25L755,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M756.25,266.25L756.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M757.5,266.25L757.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M758.75,266.25L758.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M760,266.25L760,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M761.25,266.25L761.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M762.5,266.25L762.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M763.75,266.25L763.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M765,266.25L765,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M766.25,266.25L766.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M767.5,266.25L767.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M768.75,266.25L768.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M770,266.25L770,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M771.25,266.25L771.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M772.5,266.25L772.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M773.75,266.25L773.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M775,266.25L775,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M776.25,266.25L776.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M777.5,266.25L777.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M778.75,266.25L778.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M780,266.25L780,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M781.25,266.25L781.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M782.5,266.25L782.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M783.75,266.25L783.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M785,266.25L785,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M786.25,266.25L786.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M787.5,266.25L787.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M788.75,266.25L788.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M790,266.25L790,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M791.25,266.25L791.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M792.5,266.25L792.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M793.75,266.25L793.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M795,266.25L795,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M796.25,266.25L796.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M797.5,266.25L797.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M798.75,266.25L798.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M800,266.25L800,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M801.25,266.25L801.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M802.5,266.25L802.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M803.75,266.25L803.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M805,266.25L805,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M806.25,266.25L806.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M807.5,266.25L807.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M808.75,266.25L808.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M810,266.25L810,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M811.25,266.25L811.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M812.5,266.25L812.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M813.75,266.25L813.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M815,266.25L815,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M816.25,266.25L816.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M817.5,266.25L817.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M818.75,266.25L818.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M820,266.25L820,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M821.25,266.25L821.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M822.5,266.25L822.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M823.75,266.25L823.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M825,266.25L825,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M826.25,266.25L826.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M827.5,266.25L827.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M828.75,266.25L828.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M830,266.25L830,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M831.25,266.25L831.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M832.5,266.25L832.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M833.75,266.25L833.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M835,266.25L835,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M836.25,266.25L836.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M837.5,266.25L837.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M838.75,266.25L838.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M840,266.25L840,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M841.25,266.25L841.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M842.5,266.25L842.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M843.75,266.25L843.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M845,266.25L845,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M846.25,266.25L846.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M847.5,266.25L847.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M848.75,266.25L848.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M850,266.25L850,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M851.25,266.25L851.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M852.5,266.25L852.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M853.75,266.25L853.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M855,266.25L855,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M856.25,266.25L856.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M857.5,266.25L857.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M858.75,266.25L858.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M860,266.25L860,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M861.25,266.25L861.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M862.5,266.25L862.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M863.75,266.25L863.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M865,266.25L865,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M866.25,266.25L866.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M867.5,266.25L867.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M868.75,266.25L868.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M870,266.25L870,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M871.25,266.25L871.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M872.5,266.25L872.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M873.75,266.25L873.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M875,266.25L875,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M876.25,266.25L876.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M877.5,266.25L877.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M878.75,266.25L878.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M880,266.25L880,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M881.25,266.25L881.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M882.5,266.25L882.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M883.75,266.25L883.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M885,266.25L885,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M886.25,266.25L886.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M887.5,266.25L887.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M888.75,266.25L888.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M890,266.25L890,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M891.25,266.25L891.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M892.5,266.25L892.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M893.75,266.25L893.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M895,266.25L895,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M896.25,266.25L896.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M897.5,266.25L897.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M898.75,266.25L898.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M900,266.25L900,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M901.25,266.25L901.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M902.5,266.25L902.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M903.75,266.25L903.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M905,266.25L905,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M906.25,266.25L906.25,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M907.5,266.25L907.5,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M908.75,266.25L908.75,262.5\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M910,266.25L910,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M911.25,266.25L911.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M912.5,266.25L912.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M913.75,266.25L913.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M915,266.25L915,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M916.25,266.25L916.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M917.5,266.25L917.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M918.75,266.25L918.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M920,266.25L920,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M921.25,266.25L921.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M922.5,266.25L922.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M923.75,266.25L923.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M925,266.25L925,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M926.25,266.25L926.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M927.5,266.25L927.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M928.75,266.25L928.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M930,266.25L930,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M931.25,266.25L931.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M932.5,266.25L932.5,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M933.75,266.25L933.75,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M935,266.25L935,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"<path d=\"M936.25,266.25L936.25,263.75\" style=\"fill:none;stroke:#000000;stroke-width:0.3125\" />\n" +
			"</g>\n" +
			"</svg>\n",
	},
	{
		in: blast.Output{
			Program:   "blastn",
			Version:   "BLASTN 2.2.27+",
			Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
			Database:  "nr",
			QueryId:   "33421",
			QueryDef:  "No definition line",
			QueryLen:  32,
			QuerSeq:   nil,
			Parameters: blast.Parameters{
				MatrixName:  nil,
				Expect:      1000,
				Include:     nil,
				Match:       intPtr(1),
				Mismatch:    intPtr(-3),
				GapOpen:     5,
				GapExtend:   2,
				Filter:      stringPtr("F"),
				PhiPattern:  nil,
				EntrezQuery: nil,
			},
			Iterations: []blast.Iteration{
				{
					N:        1,
					QueryId:  stringPtr("33421"),
					QueryDef: stringPtr("No definition line"),
					QueryLen: intPtr(32),
					Hits: []blast.Hit{
						{
							N:         1,
							Id:        "gi|388525227|gb|CP003531.1|",
							Def:       "Thermogladius cellulolyticus 1633, complete genome",
							Accession: "CP003531",
							Len:       1356318,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       38.1576,
									Score:          19,
									EValue:         1.72292,
									QueryFrom:      7,
									QueryTo:        29,
									HitFrom:        1187458,
									HitTo:          1187436,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(22),
									HspPositive:    intPtr(22),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(23),
									Density:        nil,
									QuerySeq:       []byte("TGTCGAACTATACGACGAGCACT"),
									SubjectSeq:     []byte("TGTCGAGCTATACGACGAGCACT"),
									FormatMidline:  []byte("|||||| ||||||||||||||||"),
								},
							},
						},
						{
							N:         2,
							Id:        "gi|354799811|gb|JN964312.1|",
							Def:       "Mus musculus targeted non-conditional, lacZ-tagged mutant allele Morn4:tm1e(KOMP)Wtsi; transgenic",
							Accession: "JN964312",
							Len:       40247,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       36.1753,
									Score:          18,
									EValue:         6.80792,
									QueryFrom:      1,
									QueryTo:        18,
									HitFrom:        26580,
									HitTo:          26597,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(18),
									HspPositive:    intPtr(18),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(18),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTATA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||||"),
								},
							},
						},
						{
							N:         24,
							Id:        "gi|356871506|emb|FO082053.1|",
							Def:       "Pichia sorbitophila strain CBS 7064 chromosome G complete sequence",
							Accession: "FO082053",
							Len:       1423303,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      13,
									QueryTo:        28,
									HitFrom:        61730,
									HitTo:          61745,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACTATACGACGAGCAC"),
									SubjectSeq:     []byte("ACTATACGACGAGCAC"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         25,
							Id:        "gi|353230524|emb|HE601625.1|",
							Def:       "Schistosoma mansoni strain Puerto Rico chromosome 2, complete genome",
							Accession: "HE601625",
							Len:       34464480,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      1,
									QueryTo:        16,
									HitFrom:        28173004,
									HitTo:          28173019,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
								{
									N:              2,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      3,
									QueryTo:        18,
									HitFrom:        28996063,
									HitTo:          28996048,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("AGAATGTCGAACTATA"),
									SubjectSeq:     []byte("AGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         26,
							Id:        "gi|341821300|emb|HE576794.1|",
							Def:       "Megasphaera elsdenii strain DSM 20460 draft genome",
							Accession: "HE576794",
							Len:       2474718,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      14,
									QueryTo:        29,
									HitFrom:        1741778,
									HitTo:          1741793,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("CTATACGACGAGCACT"),
									SubjectSeq:     []byte("CTATACGACGAGCACT"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
					},
					Statistics: &blast.Statistics{
						DbNum:    17331862,
						DbLen:    1419046940,
						HspLen:   0,
						EffSpace: 0,
						Kappa:    0.710603,
						Lambda:   1.37406,
						Entropy:  1.30725,
					},
					Message: nil,
				},
			},
			MegaStatistics: nil,
		},
		cf:        func(w, h vg.Length) vg.Canvas { return vgsvg.New(w, h) },
		configure: true,
		legend:    true,
		aligns:    true,
		depths:    false,
		expect: "" +
			"<?xml version=\"1.0\"?>\n" +
			"<!-- Generated by SVGo and Plotinum VG -->\n" +
			"<svg width=\"11.111in\" height=\"3.5556in\"\n" +
			"\txmlns=\"http://www.w3.org/2000/svg\" \n" +
			"\txmlns:xlink=\"http://www.w3.org/1999/xlink\">\n" +
			"<g transform=\"scale(1, -1) translate(0, -320)\">\n" +
			"<path d=\"M0,0L1000,0L1000,320L0,320Z\" style=\"fill:#FFFFFF\" />\n" +
			"<text x=\"6.25\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:bold;font-style:normal;font-size:13pt\">33421</text>\n" +
			"<path d=\"M187.5,270L937.5,270\" style=\"fill:none;stroke:#000000;stroke-width:1.25\" />\n" +
			"<text x=\"187.5\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">1</text>\n" +
			"<text x=\"920.82\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">29</text>\n" +
			"<text x=\"837.5\" y=\"-302.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">% Identity</text>\n" +
			"<path d=\"M612.5,313.75L625,313.75L625,301.25L612.5,301.25Z\" style=\"fill:#808080\" />\n" +
			"<text x=\"612.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">20</text>\n" +
			"<path d=\"M637.5,313.75L650,313.75L650,301.25L637.5,301.25Z\" style=\"fill:#FF0000\" />\n" +
			"<text x=\"637.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">30</text>\n" +
			"<path d=\"M662.5,313.75L675,313.75L675,301.25L662.5,301.25Z\" style=\"fill:#FFC400\" />\n" +
			"<text x=\"662.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">40</text>\n" +
			"<path d=\"M687.5,313.75L700,313.75L700,301.25L687.5,301.25Z\" style=\"fill:#FFFF00\" />\n" +
			"<text x=\"687.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">50</text>\n" +
			"<path d=\"M712.5,313.75L725,313.75L725,301.25L712.5,301.25Z\" style=\"fill:#00FF00\" />\n" +
			"<text x=\"712.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">60</text>\n" +
			"<path d=\"M737.5,313.75L750,313.75L750,301.25L737.5,301.25Z\" style=\"fill:#00FFFF\" />\n" +
			"<text x=\"737.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">70</text>\n" +
			"<path d=\"M762.5,313.75L775,313.75L775,301.25L762.5,301.25Z\" style=\"fill:#0000FF\" />\n" +
			"<text x=\"762.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">80</text>\n" +
			"<path d=\"M787.5,313.75L800,313.75L800,301.25L787.5,301.25Z\" style=\"fill:#C400FF\" />\n" +
			"<text x=\"787.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">90</text>\n" +
			"<path d=\"M812.5,313.75L825,313.75L825,301.25L812.5,301.25Z\" />\n" +
			"<text x=\"812.5\" y=\"-292.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">100</text>\n" +
			"<text x=\"12.5\" y=\"-227.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|388525227|gb|CP00</text>\n" +
			"<path d=\"M348.21,235L937.5,235\" style=\"fill:none;stroke:#C400FF;stroke-width:3.75\" />\n" +
			"<text x=\"348.21\" y=\"-242.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">7</text>\n" +
			"<text x=\"927.77\" y=\"-242.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">29</text>\n" +
			"<text x=\"348.21\" y=\"-225\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1187458</text>\n" +
			"<text x=\"898.1\" y=\"-225\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1187436 -</text>\n" +
			"<text x=\"12.5\" y=\"-182.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|354799811|gb|JN96</text>\n" +
			"<path d=\"M187.5,190L642.86,190\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"187.5\" y=\"-197.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1</text>\n" +
			"<text x=\"633.13\" y=\"-197.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">18</text>\n" +
			"<text x=\"187.5\" y=\"-180\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">26580</text>\n" +
			"<text x=\"610.99\" y=\"-180\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">26597 +</text>\n" +
			"<text x=\"12.5\" y=\"-137.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|356871506|emb|FO0</text>\n" +
			"<path d=\"M508.93,145L910.71,145\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"508.93\" y=\"-152.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">13</text>\n" +
			"<text x=\"900.98\" y=\"-152.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28</text>\n" +
			"<text x=\"508.93\" y=\"-135\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">61730</text>\n" +
			"<text x=\"878.85\" y=\"-135\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">61745 +</text>\n" +
			"<text x=\"12.5\" y=\"-92.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|353230524|emb|HE6</text>\n" +
			"<path d=\"M187.5,100L589.29,100\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"187.5\" y=\"-107.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1</text>\n" +
			"<text x=\"579.56\" y=\"-107.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">16</text>\n" +
			"<text x=\"187.5\" y=\"-90\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28173004</text>\n" +
			"<text x=\"542.82\" y=\"-90\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28173019 +</text>\n" +
			"<path d=\"M241.07,80L642.86,80\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"241.07\" y=\"-87.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">3</text>\n" +
			"<text x=\"633.13\" y=\"-87.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">18</text>\n" +
			"<text x=\"241.07\" y=\"-70\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28996063</text>\n" +
			"<text x=\"598.59\" y=\"-70\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28996048 -</text>\n" +
			"<text x=\"12.5\" y=\"-27.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|341821300|emb|HE5</text>\n" +
			"<path d=\"M535.71,35L937.5,35\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"535.71\" y=\"-42.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">14</text>\n" +
			"<text x=\"927.77\" y=\"-42.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">29</text>\n" +
			"<text x=\"535.71\" y=\"-25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1741778</text>\n" +
			"<text x=\"895.9\" y=\"-25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1741793 +</text>\n" +
			"</g>\n" +
			"</svg>\n",
	},
	{
		in: blast.Output{
			Program:   "blastn",
			Version:   "BLASTN 2.2.27+",
			Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
			Database:  "nr",
			QueryId:   "33421",
			QueryDef:  "No definition line",
			QueryLen:  32,
			QuerSeq:   nil,
			Parameters: blast.Parameters{
				MatrixName:  nil,
				Expect:      1000,
				Include:     nil,
				Match:       intPtr(1),
				Mismatch:    intPtr(-3),
				GapOpen:     5,
				GapExtend:   2,
				Filter:      stringPtr("F"),
				PhiPattern:  nil,
				EntrezQuery: nil,
			},
			Iterations: []blast.Iteration{
				{
					N:        1,
					QueryId:  stringPtr("33421"),
					QueryDef: stringPtr("No definition line"),
					QueryLen: intPtr(32),
					Hits: []blast.Hit{
						{
							N:         1,
							Id:        "gi|388525227|gb|CP003531.1|",
							Def:       "Thermogladius cellulolyticus 1633, complete genome",
							Accession: "CP003531",
							Len:       1356318,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       38.1576,
									Score:          19,
									EValue:         1.72292,
									QueryFrom:      7,
									QueryTo:        29,
									HitFrom:        1187458,
									HitTo:          1187436,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(22),
									HspPositive:    intPtr(22),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(23),
									Density:        nil,
									QuerySeq:       []byte("TGTCGAACTATACGACGAGCACT"),
									SubjectSeq:     []byte("TGTCGAGCTATACGACGAGCACT"),
									FormatMidline:  []byte("|||||| ||||||||||||||||"),
								},
							},
						},
						{
							N:         2,
							Id:        "gi|354799811|gb|JN964312.1|",
							Def:       "Mus musculus targeted non-conditional, lacZ-tagged mutant allele Morn4:tm1e(KOMP)Wtsi; transgenic",
							Accession: "JN964312",
							Len:       40247,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       36.1753,
									Score:          18,
									EValue:         6.80792,
									QueryFrom:      1,
									QueryTo:        18,
									HitFrom:        26580,
									HitTo:          26597,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(18),
									HspPositive:    intPtr(18),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(18),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTATA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||||"),
								},
							},
						},
						{
							N:         24,
							Id:        "gi|356871506|emb|FO082053.1|",
							Def:       "Pichia sorbitophila strain CBS 7064 chromosome G complete sequence",
							Accession: "FO082053",
							Len:       1423303,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      13,
									QueryTo:        28,
									HitFrom:        61730,
									HitTo:          61745,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACTATACGACGAGCAC"),
									SubjectSeq:     []byte("ACTATACGACGAGCAC"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         25,
							Id:        "gi|353230524|emb|HE601625.1|",
							Def:       "Schistosoma mansoni strain Puerto Rico chromosome 2, complete genome",
							Accession: "HE601625",
							Len:       34464480,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      1,
									QueryTo:        16,
									HitFrom:        28173004,
									HitTo:          28173019,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
								{
									N:              2,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      3,
									QueryTo:        18,
									HitFrom:        28996063,
									HitTo:          28996048,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("AGAATGTCGAACTATA"),
									SubjectSeq:     []byte("AGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         26,
							Id:        "gi|341821300|emb|HE576794.1|",
							Def:       "Megasphaera elsdenii strain DSM 20460 draft genome",
							Accession: "HE576794",
							Len:       2474718,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      14,
									QueryTo:        29,
									HitFrom:        1741778,
									HitTo:          1741793,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("CTATACGACGAGCACT"),
									SubjectSeq:     []byte("CTATACGACGAGCACT"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
					},
					Statistics: &blast.Statistics{
						DbNum:    17331862,
						DbLen:    1419046940,
						HspLen:   0,
						EffSpace: 0,
						Kappa:    0.710603,
						Lambda:   1.37406,
						Entropy:  1.30725,
					},
					Message: nil,
				},
			},
			MegaStatistics: nil,
		},
		cf:        func(w, h vg.Length) vg.Canvas { return vgsvg.New(w, h) },
		configure: true,
		legend:    false,
		aligns:    true,
		depths:    false,
		expect: "" +
			"<?xml version=\"1.0\"?>\n" +
			"<!-- Generated by SVGo and Plotinum VG -->\n" +
			"<svg width=\"11.111in\" height=\"3.5556in\"\n" +
			"\txmlns=\"http://www.w3.org/2000/svg\" \n" +
			"\txmlns:xlink=\"http://www.w3.org/1999/xlink\">\n" +
			"<g transform=\"scale(1, -1) translate(0, -320)\">\n" +
			"<path d=\"M0,0L1000,0L1000,320L0,320Z\" style=\"fill:#FFFFFF\" />\n" +
			"<text x=\"6.25\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:bold;font-style:normal;font-size:13pt\">33421</text>\n" +
			"<path d=\"M187.5,270L937.5,270\" style=\"fill:none;stroke:#000000;stroke-width:1.25\" />\n" +
			"<text x=\"187.5\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">1</text>\n" +
			"<text x=\"920.82\" y=\"-285\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">29</text>\n" +
			"<text x=\"12.5\" y=\"-227.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|388525227|gb|CP00</text>\n" +
			"<path d=\"M348.21,235L937.5,235\" style=\"fill:none;stroke:#C400FF;stroke-width:3.75\" />\n" +
			"<text x=\"348.21\" y=\"-242.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">7</text>\n" +
			"<text x=\"927.77\" y=\"-242.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">29</text>\n" +
			"<text x=\"348.21\" y=\"-225\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1187458</text>\n" +
			"<text x=\"898.1\" y=\"-225\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1187436 -</text>\n" +
			"<text x=\"12.5\" y=\"-182.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|354799811|gb|JN96</text>\n" +
			"<path d=\"M187.5,190L642.86,190\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"187.5\" y=\"-197.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1</text>\n" +
			"<text x=\"633.13\" y=\"-197.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">18</text>\n" +
			"<text x=\"187.5\" y=\"-180\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">26580</text>\n" +
			"<text x=\"610.99\" y=\"-180\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">26597 +</text>\n" +
			"<text x=\"12.5\" y=\"-137.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|356871506|emb|FO0</text>\n" +
			"<path d=\"M508.93,145L910.71,145\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"508.93\" y=\"-152.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">13</text>\n" +
			"<text x=\"900.98\" y=\"-152.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28</text>\n" +
			"<text x=\"508.93\" y=\"-135\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">61730</text>\n" +
			"<text x=\"878.85\" y=\"-135\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">61745 +</text>\n" +
			"<text x=\"12.5\" y=\"-92.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|353230524|emb|HE6</text>\n" +
			"<path d=\"M187.5,100L589.29,100\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"187.5\" y=\"-107.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1</text>\n" +
			"<text x=\"579.56\" y=\"-107.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">16</text>\n" +
			"<text x=\"187.5\" y=\"-90\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28173004</text>\n" +
			"<text x=\"542.82\" y=\"-90\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28173019 +</text>\n" +
			"<path d=\"M241.07,80L642.86,80\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"241.07\" y=\"-87.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">3</text>\n" +
			"<text x=\"633.13\" y=\"-87.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">18</text>\n" +
			"<text x=\"241.07\" y=\"-70\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28996063</text>\n" +
			"<text x=\"598.59\" y=\"-70\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">28996048 -</text>\n" +
			"<text x=\"12.5\" y=\"-27.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">gi|341821300|emb|HE5</text>\n" +
			"<path d=\"M535.71,35L937.5,35\" style=\"fill:none;stroke:#000000;stroke-width:3.75\" />\n" +
			"<text x=\"535.71\" y=\"-42.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">14</text>\n" +
			"<text x=\"927.77\" y=\"-42.5\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">29</text>\n" +
			"<text x=\"535.71\" y=\"-25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1741778</text>\n" +
			"<text x=\"895.9\" y=\"-25\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:7pt\">1741793 +</text>\n" +
			"</g>\n" +
			"</svg>\n",
	},
	{
		in: blast.Output{
			Program:   "blastn",
			Version:   "BLASTN 2.2.27+",
			Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
			Database:  "nr",
			QueryId:   "33421",
			QueryDef:  "No definition line",
			QueryLen:  32,
			QuerSeq:   nil,
			Parameters: blast.Parameters{
				MatrixName:  nil,
				Expect:      1000,
				Include:     nil,
				Match:       intPtr(1),
				Mismatch:    intPtr(-3),
				GapOpen:     5,
				GapExtend:   2,
				Filter:      stringPtr("F"),
				PhiPattern:  nil,
				EntrezQuery: nil,
			},
			Iterations: []blast.Iteration{
				{
					N:        1,
					QueryId:  stringPtr("33421"),
					QueryDef: stringPtr("No definition line"),
					QueryLen: intPtr(32),
					Hits: []blast.Hit{
						{
							N:         1,
							Id:        "gi|388525227|gb|CP003531.1|",
							Def:       "Thermogladius cellulolyticus 1633, complete genome",
							Accession: "CP003531",
							Len:       1356318,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       38.1576,
									Score:          19,
									EValue:         1.72292,
									QueryFrom:      7,
									QueryTo:        29,
									HitFrom:        1187458,
									HitTo:          1187436,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(22),
									HspPositive:    intPtr(22),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(23),
									Density:        nil,
									QuerySeq:       []byte("TGTCGAACTATACGACGAGCACT"),
									SubjectSeq:     []byte("TGTCGAGCTATACGACGAGCACT"),
									FormatMidline:  []byte("|||||| ||||||||||||||||"),
								},
							},
						},
						{
							N:         2,
							Id:        "gi|354799811|gb|JN964312.1|",
							Def:       "Mus musculus targeted non-conditional, lacZ-tagged mutant allele Morn4:tm1e(KOMP)Wtsi; transgenic",
							Accession: "JN964312",
							Len:       40247,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       36.1753,
									Score:          18,
									EValue:         6.80792,
									QueryFrom:      1,
									QueryTo:        18,
									HitFrom:        26580,
									HitTo:          26597,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(18),
									HspPositive:    intPtr(18),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(18),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTATA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||||"),
								},
							},
						},
						{
							N:         24,
							Id:        "gi|356871506|emb|FO082053.1|",
							Def:       "Pichia sorbitophila strain CBS 7064 chromosome G complete sequence",
							Accession: "FO082053",
							Len:       1423303,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      13,
									QueryTo:        28,
									HitFrom:        61730,
									HitTo:          61745,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACTATACGACGAGCAC"),
									SubjectSeq:     []byte("ACTATACGACGAGCAC"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         25,
							Id:        "gi|353230524|emb|HE601625.1|",
							Def:       "Schistosoma mansoni strain Puerto Rico chromosome 2, complete genome",
							Accession: "HE601625",
							Len:       34464480,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      1,
									QueryTo:        16,
									HitFrom:        28173004,
									HitTo:          28173019,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("ACAGAATGTCGAACTA"),
									SubjectSeq:     []byte("ACAGAATGTCGAACTA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
								{
									N:              2,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      3,
									QueryTo:        18,
									HitFrom:        28996063,
									HitTo:          28996048,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(-1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("AGAATGTCGAACTATA"),
									SubjectSeq:     []byte("AGAATGTCGAACTATA"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
						{
							N:         26,
							Id:        "gi|341821300|emb|HE576794.1|",
							Def:       "Megasphaera elsdenii strain DSM 20460 draft genome",
							Accession: "HE576794",
							Len:       2474718,
							Hsps: []blast.Hsp{
								{
									N:              1,
									BitScore:       32.2105,
									Score:          16,
									EValue:         106.294,
									QueryFrom:      14,
									QueryTo:        29,
									HitFrom:        1741778,
									HitTo:          1741793,
									PhiPatternFrom: nil,
									PhiPatternTo:   nil,
									QueryFrame:     intPtr(1),
									HitFrame:       intPtr(1),
									HspIdentity:    intPtr(16),
									HspPositive:    intPtr(16),
									HspGaps:        intPtr(0),
									AlignLen:       intPtr(16),
									Density:        nil,
									QuerySeq:       []byte("CTATACGACGAGCACT"),
									SubjectSeq:     []byte("CTATACGACGAGCACT"),
									FormatMidline:  []byte("||||||||||||||||"),
								},
							},
						},
					},
					Statistics: &blast.Statistics{
						DbNum:    17331862,
						DbLen:    1419046940,
						HspLen:   0,
						EffSpace: 0,
						Kappa:    0.710603,
						Lambda:   1.37406,
						Entropy:  1.30725,
					},
					Message: nil,
				},
			},
			MegaStatistics: nil,
		},
		cf:        func(w, h vg.Length) vg.Canvas { return vgsvg.New(w, h) },
		configure: true,
		legend:    true,
		aligns:    false,
		depths:    false,
		expect: "" +
			"<?xml version=\"1.0\"?>\n" +
			"<!-- Generated by SVGo and Plotinum VG -->\n" +
			"<svg width=\"11.111in\" height=\"0.83333in\"\n" +
			"\txmlns=\"http://www.w3.org/2000/svg\" \n" +
			"\txmlns:xlink=\"http://www.w3.org/1999/xlink\">\n" +
			"<g transform=\"scale(1, -1) translate(0, -75)\">\n" +
			"<path d=\"M0,0L1000,0L1000,75L0,75Z\" style=\"fill:#FFFFFF\" />\n" +
			"<text x=\"6.25\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:bold;font-style:normal;font-size:13pt\">33421</text>\n" +
			"<path d=\"M187.5,25L937.5,25\" style=\"fill:none;stroke:#000000;stroke-width:1.25\" />\n" +
			"<text x=\"187.5\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">1</text>\n" +
			"<text x=\"920.82\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">29</text>\n" +
			"</g>\n" +
			"</svg>\n",
	},
	{
		in: blast.Output{
			Program:   "blastn",
			Version:   "BLASTN 2.2.27+",
			Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
			Database:  "nr",
			QueryId:   "33421",
			QueryDef:  "No definition line",
			QueryLen:  32,
			QuerSeq:   nil,
			Parameters: blast.Parameters{
				MatrixName:  nil,
				Expect:      1000,
				Include:     nil,
				Match:       intPtr(1),
				Mismatch:    intPtr(-3),
				GapOpen:     5,
				GapExtend:   2,
				Filter:      stringPtr("F"),
				PhiPattern:  nil,
				EntrezQuery: nil,
			},
			Iterations:     nil,
			MegaStatistics: nil,
		},
		cf: func(w, h vg.Length) vg.Canvas { return vgsvg.New(w, h) },
		expect: "" +
			"<?xml version=\"1.0\"?>\n" +
			"<!-- Generated by SVGo and Plotinum VG -->\n" +
			"<svg width=\"11.111in\" height=\"0.83333in\"\n" +
			"\txmlns=\"http://www.w3.org/2000/svg\" \n" +
			"\txmlns:xlink=\"http://www.w3.org/1999/xlink\">\n" +
			"<g transform=\"scale(1, -1) translate(0, -75)\">\n" +
			"<path d=\"M0,0L1000,0L1000,75L0,75Z\" style=\"fill:#FFFFFF\" />\n" +
			"<text x=\"6.25\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:bold;font-style:normal;font-size:13pt\">33421</text>\n" +
			"<path d=\"M187.5,25L937.5,25\" style=\"fill:none;stroke:#000000;stroke-width:1.25\" />\n" +
			"<text x=\"187.5\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">1</text>\n" +
			"<text x=\"920.82\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">32</text>\n" +
			"</g>\n" +
			"</svg>\n",
	},
	{
		in: blast.Output{},
		cf: func(w, h vg.Length) vg.Canvas { return vgsvg.New(w, h) },
		expect: "" +
			"<?xml version=\"1.0\"?>\n" +
			"<!-- Generated by SVGo and Plotinum VG -->\n" +
			"<svg width=\"11.111in\" height=\"0.83333in\"\n" +
			"\txmlns=\"http://www.w3.org/2000/svg\" \n" +
			"\txmlns:xlink=\"http://www.w3.org/1999/xlink\">\n" +
			"<g transform=\"scale(1, -1) translate(0, -75)\">\n" +
			"<path d=\"M0,0L1000,0L1000,75L0,75Z\" style=\"fill:#FFFFFF\" />\n" +
			"<text x=\"6.25\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:bold;font-style:normal;font-size:13pt\"></text>\n" +
			"<path d=\"M187.5,25L937.5,25\" style=\"fill:none;stroke:#000000;stroke-width:1.25\" />\n" +
			"<text x=\"187.5\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">0</text>\n" +
			"<text x=\"929.16\" y=\"-40\" transform=\"scale(1, -1)\"\n" +
			"\tstyle=\"font-family:Helvetica;font-weight:normal;font-style:normal;font-size:12pt\">0</text>\n" +
			"</g>\n" +
			"</svg>\n",
	},
}

func (s *S) TestSummary(c *check.C) {
	for _, t := range testOutputs {
		sum := graphic.NewSummary(t.in)
		if t.configure {
			sum.Legend = t.legend
			sum.Aligns = t.aligns
			sum.Depths = t.depths
		}
		var b bytes.Buffer
		sum.Render(t.cf).(io.WriterTo).WriteTo(&b)
		c.Check(b.String(), check.Equals, t.expect)
	}
}
