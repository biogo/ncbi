// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

import (
	"encoding/xml"
	"strings"

	"gopkg.in/check.v1"
)

func (s *S) TestParseOutput(c *check.C) {
	for i, t := range []struct {
		retval string
		output Output
	}{
		{
			`<?xml version="1.0"?>
<!DOCTYPE BlastOutput PUBLIC "-//NCBI//NCBI BlastOutput/EN" "http://www.ncbi.nlm.nih.gov/dtd/NCBI_BlastOutput.dtd">
<BlastOutput>
  <BlastOutput_program>blastn</BlastOutput_program>
  <BlastOutput_version>BLASTN 2.2.27+</BlastOutput_version>
  <BlastOutput_reference>Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&amp;auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), &quot;Gapped BLAST and PSI-BLAST: a new generation of protein database search programs&quot;, Nucleic Acids Res. 25:3389-3402.</BlastOutput_reference>
  <BlastOutput_db>nr</BlastOutput_db>
  <BlastOutput_query-ID>33421</BlastOutput_query-ID>
  <BlastOutput_query-def>No definition line</BlastOutput_query-def>
  <BlastOutput_query-len>32</BlastOutput_query-len>
  <BlastOutput_param>
    <Parameters>
      <Parameters_expect>1000</Parameters_expect>
      <Parameters_sc-match>1</Parameters_sc-match>
      <Parameters_sc-mismatch>-3</Parameters_sc-mismatch>
      <Parameters_gap-open>5</Parameters_gap-open>
      <Parameters_gap-extend>2</Parameters_gap-extend>
      <Parameters_filter>F</Parameters_filter>
    </Parameters>
  </BlastOutput_param>
<BlastOutput_iterations>
<Iteration>
  <Iteration_iter-num>1</Iteration_iter-num>
  <Iteration_query-ID>33421</Iteration_query-ID>
  <Iteration_query-def>No definition line</Iteration_query-def>
  <Iteration_query-len>32</Iteration_query-len>
<Iteration_hits>
<Hit>
  <Hit_num>1</Hit_num>
  <Hit_id>gi|388525227|gb|CP003531.1|</Hit_id>
  <Hit_def>Thermogladius cellulolyticus 1633, complete genome</Hit_def>
  <Hit_accession>CP003531</Hit_accession>
  <Hit_len>1356318</Hit_len>
  <Hit_hsps>
    <Hsp>
      <Hsp_num>1</Hsp_num>
      <Hsp_bit-score>38.1576</Hsp_bit-score>
      <Hsp_score>19</Hsp_score>
      <Hsp_evalue>1.72292</Hsp_evalue>
      <Hsp_query-from>7</Hsp_query-from>
      <Hsp_query-to>29</Hsp_query-to>
      <Hsp_hit-from>1187458</Hsp_hit-from>
      <Hsp_hit-to>1187436</Hsp_hit-to>
      <Hsp_query-frame>1</Hsp_query-frame>
      <Hsp_hit-frame>-1</Hsp_hit-frame>
      <Hsp_identity>22</Hsp_identity>
      <Hsp_positive>22</Hsp_positive>
      <Hsp_gaps>0</Hsp_gaps>
      <Hsp_align-len>23</Hsp_align-len>
      <Hsp_qseq>TGTCGAACTATACGACGAGCACT</Hsp_qseq>
      <Hsp_hseq>TGTCGAGCTATACGACGAGCACT</Hsp_hseq>
      <Hsp_midline>|||||| ||||||||||||||||</Hsp_midline>
    </Hsp>
  </Hit_hsps>
</Hit>
<Hit>
  <Hit_num>2</Hit_num>
  <Hit_id>gi|354799811|gb|JN964312.1|</Hit_id>
  <Hit_def>Mus musculus targeted non-conditional, lacZ-tagged mutant allele Morn4:tm1e(KOMP)Wtsi; transgenic</Hit_def>
  <Hit_accession>JN964312</Hit_accession>
  <Hit_len>40247</Hit_len>
  <Hit_hsps>
    <Hsp>
      <Hsp_num>1</Hsp_num>
      <Hsp_bit-score>36.1753</Hsp_bit-score>
      <Hsp_score>18</Hsp_score>
      <Hsp_evalue>6.80792</Hsp_evalue>
      <Hsp_query-from>1</Hsp_query-from>
      <Hsp_query-to>18</Hsp_query-to>
      <Hsp_hit-from>26580</Hsp_hit-from>
      <Hsp_hit-to>26597</Hsp_hit-to>
      <Hsp_query-frame>1</Hsp_query-frame>
      <Hsp_hit-frame>1</Hsp_hit-frame>
      <Hsp_identity>18</Hsp_identity>
      <Hsp_positive>18</Hsp_positive>
      <Hsp_gaps>0</Hsp_gaps>
      <Hsp_align-len>18</Hsp_align-len>
      <Hsp_qseq>ACAGAATGTCGAACTATA</Hsp_qseq>
      <Hsp_hseq>ACAGAATGTCGAACTATA</Hsp_hseq>
      <Hsp_midline>||||||||||||||||||</Hsp_midline>
    </Hsp>
  </Hit_hsps>
</Hit>
<Hit>
  <Hit_num>24</Hit_num>
  <Hit_id>gi|356871506|emb|FO082053.1|</Hit_id>
  <Hit_def>Pichia sorbitophila strain CBS 7064 chromosome G complete sequence</Hit_def>
  <Hit_accession>FO082053</Hit_accession>
  <Hit_len>1423303</Hit_len>
  <Hit_hsps>
    <Hsp>
      <Hsp_num>1</Hsp_num>
      <Hsp_bit-score>32.2105</Hsp_bit-score>
      <Hsp_score>16</Hsp_score>
      <Hsp_evalue>106.294</Hsp_evalue>
      <Hsp_query-from>13</Hsp_query-from>
      <Hsp_query-to>28</Hsp_query-to>
      <Hsp_hit-from>61730</Hsp_hit-from>
      <Hsp_hit-to>61745</Hsp_hit-to>
      <Hsp_query-frame>1</Hsp_query-frame>
      <Hsp_hit-frame>1</Hsp_hit-frame>
      <Hsp_identity>16</Hsp_identity>
      <Hsp_positive>16</Hsp_positive>
      <Hsp_gaps>0</Hsp_gaps>
      <Hsp_align-len>16</Hsp_align-len>
      <Hsp_qseq>ACTATACGACGAGCAC</Hsp_qseq>
      <Hsp_hseq>ACTATACGACGAGCAC</Hsp_hseq>
      <Hsp_midline>||||||||||||||||</Hsp_midline>
    </Hsp>
  </Hit_hsps>
</Hit>
<Hit>
  <Hit_num>25</Hit_num>
  <Hit_id>gi|353230524|emb|HE601625.1|</Hit_id>
  <Hit_def>Schistosoma mansoni strain Puerto Rico chromosome 2, complete genome</Hit_def>
  <Hit_accession>HE601625</Hit_accession>
  <Hit_len>34464480</Hit_len>
  <Hit_hsps>
    <Hsp>
      <Hsp_num>1</Hsp_num>
      <Hsp_bit-score>32.2105</Hsp_bit-score>
      <Hsp_score>16</Hsp_score>
      <Hsp_evalue>106.294</Hsp_evalue>
      <Hsp_query-from>1</Hsp_query-from>
      <Hsp_query-to>16</Hsp_query-to>
      <Hsp_hit-from>28173004</Hsp_hit-from>
      <Hsp_hit-to>28173019</Hsp_hit-to>
      <Hsp_query-frame>1</Hsp_query-frame>
      <Hsp_hit-frame>1</Hsp_hit-frame>
      <Hsp_identity>16</Hsp_identity>
      <Hsp_positive>16</Hsp_positive>
      <Hsp_gaps>0</Hsp_gaps>
      <Hsp_align-len>16</Hsp_align-len>
      <Hsp_qseq>ACAGAATGTCGAACTA</Hsp_qseq>
      <Hsp_hseq>ACAGAATGTCGAACTA</Hsp_hseq>
      <Hsp_midline>||||||||||||||||</Hsp_midline>
    </Hsp>
    <Hsp>
      <Hsp_num>2</Hsp_num>
      <Hsp_bit-score>32.2105</Hsp_bit-score>
      <Hsp_score>16</Hsp_score>
      <Hsp_evalue>106.294</Hsp_evalue>
      <Hsp_query-from>3</Hsp_query-from>
      <Hsp_query-to>18</Hsp_query-to>
      <Hsp_hit-from>28996063</Hsp_hit-from>
      <Hsp_hit-to>28996048</Hsp_hit-to>
      <Hsp_query-frame>1</Hsp_query-frame>
      <Hsp_hit-frame>-1</Hsp_hit-frame>
      <Hsp_identity>16</Hsp_identity>
      <Hsp_positive>16</Hsp_positive>
      <Hsp_gaps>0</Hsp_gaps>
      <Hsp_align-len>16</Hsp_align-len>
      <Hsp_qseq>AGAATGTCGAACTATA</Hsp_qseq>
      <Hsp_hseq>AGAATGTCGAACTATA</Hsp_hseq>
      <Hsp_midline>||||||||||||||||</Hsp_midline>
    </Hsp>
  </Hit_hsps>
</Hit>
<Hit>
  <Hit_num>26</Hit_num>
  <Hit_id>gi|341821300|emb|HE576794.1|</Hit_id>
  <Hit_def>Megasphaera elsdenii strain DSM 20460 draft genome</Hit_def>
  <Hit_accession>HE576794</Hit_accession>
  <Hit_len>2474718</Hit_len>
  <Hit_hsps>
    <Hsp>
      <Hsp_num>1</Hsp_num>
      <Hsp_bit-score>32.2105</Hsp_bit-score>
      <Hsp_score>16</Hsp_score>
      <Hsp_evalue>106.294</Hsp_evalue>
      <Hsp_query-from>14</Hsp_query-from>
      <Hsp_query-to>29</Hsp_query-to>
      <Hsp_hit-from>1741778</Hsp_hit-from>
      <Hsp_hit-to>1741793</Hsp_hit-to>
      <Hsp_query-frame>1</Hsp_query-frame>
      <Hsp_hit-frame>1</Hsp_hit-frame>
      <Hsp_identity>16</Hsp_identity>
      <Hsp_positive>16</Hsp_positive>
      <Hsp_gaps>0</Hsp_gaps>
      <Hsp_align-len>16</Hsp_align-len>
      <Hsp_qseq>CTATACGACGAGCACT</Hsp_qseq>
      <Hsp_hseq>CTATACGACGAGCACT</Hsp_hseq>
      <Hsp_midline>||||||||||||||||</Hsp_midline>
    </Hsp>
  </Hit_hsps>
</Hit>
</Iteration_hits>
  <Iteration_stat>
    <Statistics>
      <Statistics_db-num>17331862</Statistics_db-num>
      <Statistics_db-len>1419046940</Statistics_db-len>
      <Statistics_hsp-len>0</Statistics_hsp-len>
      <Statistics_eff-space>0</Statistics_eff-space>
      <Statistics_kappa>0.710603</Statistics_kappa>
      <Statistics_lambda>1.37406</Statistics_lambda>
      <Statistics_entropy>1.30725</Statistics_entropy>
    </Statistics>
  </Iteration_stat>
</Iteration>
</BlastOutput_iterations>
</BlastOutput>

`,
			Output{
				Program:   "blastn",
				Version:   "BLASTN 2.2.27+",
				Reference: "Stephen F. Altschul, Thomas L. Madden, Alejandro A. Sch&auml;ffer, Jinghui Zhang, Zheng Zhang, Webb Miller, and David J. Lipman (1997), \"Gapped BLAST and PSI-BLAST: a new generation of protein database search programs\", Nucleic Acids Res. 25:3389-3402.",
				Database:  "nr",
				QueryId:   "33421",
				QueryDef:  "No definition line",
				QueryLen:  32,
				QuerSeq:   nil,
				Parameters: Parameters{
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
				Iterations: []Iteration{
					{
						N:        1,
						QueryId:  stringPtr("33421"),
						QueryDef: stringPtr("No definition line"),
						QueryLen: intPtr(32),
						Hits: []Hit{
							{
								N:         1,
								Id:        "gi|388525227|gb|CP003531.1|",
								Def:       "Thermogladius cellulolyticus 1633, complete genome",
								Accession: "CP003531",
								Len:       1356318,
								Hsps: []Hsp{
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
								Hsps: []Hsp{
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
								Hsps: []Hsp{
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
								Hsps: []Hsp{
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
								Hsps: []Hsp{
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
						Statistics: &Statistics{
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
		},
	} {
		var o Output
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&o)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(o, check.DeepEquals, t.output, check.Commentf("Test: %d", i))
	}
}
