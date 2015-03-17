// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

/*
<!-- ============================================
     ::DATATOOL:: Generated from "blastxml.asn"
     ::DATATOOL:: by application DATATOOL version 1.5.0
     ::DATATOOL:: on 06/06/2006 23:03:48
     ============================================ -->

<!-- NCBI_BlastOutput.dtd
  This file is built from a series of basic modules.
  The actual ELEMENT and ENTITY declarations are in the modules.
  This file is used to put them together.
-->

<!ENTITY % NCBI_Entity_module PUBLIC "-//NCBI//NCBI Entity Module//EN" "NCBI_Entity.mod.dtd">
%NCBI_Entity_module;

<!ENTITY % NCBI_BlastOutput_module PUBLIC "-//NCBI//NCBI BlastOutput Module//EN" "NCBI_BlastOutput.mod.dtd">
%NCBI_BlastOutput_module;
*/

/*
<!-- ============================================
     ::DATATOOL:: Generated from "blastxml.asn"
     ::DATATOOL:: by application DATATOOL version 2.0.0
     ::DATATOOL:: on 08/02/2010 23:05:14
     ============================================ -->

<!-- ============================================ -->
<!-- This section is mapped from module "NCBI-BlastOutput"
================================================= -->

<!--$Id: blastxml.asn 100080 2007-03-12 16:05:35Z kazimird $ -->

<!ELEMENT BlastOutput (
        BlastOutput_program,
        BlastOutput_version,
        BlastOutput_reference,
        BlastOutput_db,
        BlastOutput_query-ID,
        BlastOutput_query-def,
        BlastOutput_query-len,
        BlastOutput_query-seq?,
        BlastOutput_param,
        BlastOutput_iterations,
        BlastOutput_mbstat?)>

<!-- BLAST program: blastp, tblastx etc. -->
<!ELEMENT BlastOutput_program (#PCDATA)>

<!-- Program version  -->
<!ELEMENT BlastOutput_version (#PCDATA)>

<!-- Steven, David, Tom and others -->
<!ELEMENT BlastOutput_reference (#PCDATA)>

<!-- BLAST Database name -->
<!ELEMENT BlastOutput_db (#PCDATA)>

<!-- SeqId of query -->
<!ELEMENT BlastOutput_query-ID (#PCDATA)>

<!-- Definition line of query -->
<!ELEMENT BlastOutput_query-def (#PCDATA)>

<!-- length of query sequence -->
<!ELEMENT BlastOutput_query-len (%INTEGER;)>

<!-- query sequence itself -->
<!ELEMENT BlastOutput_query-seq (#PCDATA)>

<!-- search parameters -->
<!ELEMENT BlastOutput_param (Parameters)>

<!ELEMENT BlastOutput_iterations (Iteration*)>

<!-- Mega BLAST search statistics -->
<!ELEMENT BlastOutput_mbstat (Statistics)>
*/

// An Output holds the deserialised results of an Blast Get request.
type Output struct {
	Program        string      `xml:"BlastOutput_program"`              // BlastOutput_program
	Version        string      `xml:"BlastOutput_version"`              // BlastOutput_version
	Reference      string      `xml:"BlastOutput_reference"`            // BlastOutput_reference
	Database       string      `xml:"BlastOutput_db"`                   // BlastOutput_db
	QueryId        string      `xml:"BlastOutput_query-ID"`             // BlastOutput_query-ID
	QueryDef       string      `xml:"BlastOutput_query-def"`            // BlastOutput_query-def
	QueryLen       int         `xml:"BlastOutput_query-len"`            // BlastOutput_query-len
	QuerSeq        *string     `xml:"BlastOutput_query-seq"`            // BlastOutput_query-seq?
	Parameters     Parameters  `xml:"BlastOutput_param>Parameters"`     // BlastOutput_param
	Iterations     []Iteration `xml:"BlastOutput_iterations>Iteration"` // BlastOutput_iterations
	MegaStatistics *Statistics `xml:"BlastOutput_mbstat>Statistics"`    // BlastOutput_mbstat?
}

/*
<!ELEMENT Iteration (
        Iteration_iter-num,
        Iteration_query-ID?,
        Iteration_query-def?,
        Iteration_query-len?,
        Iteration_hits?,
        Iteration_stat?,
        Iteration_message?)>

<!-- iteration number -->
<!ELEMENT Iteration_iter-num (%INTEGER;)>

<!-- SeqId of query -->
<!ELEMENT Iteration_query-ID (#PCDATA)>

<!-- Definition line of query -->
<!ELEMENT Iteration_query-def (#PCDATA)>

<!-- length of query sequence -->
<!ELEMENT Iteration_query-len (%INTEGER;)>

<!-- Hits one for every db sequence -->
<!ELEMENT Iteration_hits (Hit*)>

<!-- search statistics             -->
<!ELEMENT Iteration_stat (Statistics)>

<!-- Some (error?) information -->
<!ELEMENT Iteration_message (#PCDATA)>
*/

// An Iteration holds the iteration data for a Blast result.
type Iteration struct {
	N          int         `xml:"Iteration_iter-num"`        // Iteration_iter-num
	QueryId    *string     `xml:"Iteration_query-ID"`        // Iteration_query-ID?
	QueryDef   *string     `xml:"Iteration_query-def"`       // Iteration_query-def?
	QueryLen   *int        `xml:"Iteration_query-len"`       // Iteration_query-len?
	Hits       []Hit       `xml:"Iteration_hits>Hit"`        // Iteration_hits?
	Statistics *Statistics `xml:"Iteration_stat>Statistics"` // Iteration_stat?
	Message    *string     `xml:"Iteration_message"`         // Iteration_message?
}

/*
<!ELEMENT Parameters (
        Parameters_matrix?,
        Parameters_expect,
        Parameters_include?,
        Parameters_sc-match?,
        Parameters_sc-mismatch?,
        Parameters_gap-open,
        Parameters_gap-extend,
        Parameters_filter?,
        Parameters_pattern?,
        Parameters_entrez-query?)>

<!-- Matrix used (-M) -->
<!ELEMENT Parameters_matrix (#PCDATA)>

<!-- Expectation threshold (-e) -->
<!ELEMENT Parameters_expect (%REAL;)>

<!-- Inclusion threshold (-h) -->
<!ELEMENT Parameters_include (%REAL;)>

<!-- match score for NT (-r) -->
<!ELEMENT Parameters_sc-match (%INTEGER;)>

<!-- mismatch score for NT (-q) -->
<!ELEMENT Parameters_sc-mismatch (%INTEGER;)>

<!-- Gap opening cost (-G) -->
<!ELEMENT Parameters_gap-open (%INTEGER;)>

<!-- Gap extension cost (-E) -->
<!ELEMENT Parameters_gap-extend (%INTEGER;)>

<!-- Filtering options (-F) -->
<!ELEMENT Parameters_filter (#PCDATA)>

<!-- PHI-BLAST pattern -->
<!ELEMENT Parameters_pattern (#PCDATA)>

<!-- Limit of request to Entrez query -->
<!ELEMENT Parameters_entrez-query (#PCDATA)>
*/

// A Parameters holds the parameter information for a Blast result.
type Parameters struct {
	MatrixName  *string  `xml:"Parameters_matrix"`       // Parameters_matrix?
	Expect      float64  `xml:"Parameters_expect"`       // Parameters_expect
	Include     *float64 `xml:"Parameters_include"`      // Parameters_include?
	Match       *int     `xml:"Parameters_sc-match"`     // Parameters_sc-match?
	Mismatch    *int     `xml:"Parameters_sc-mismatch"`  // Parameters_sc-mismatch?
	GapOpen     int      `xml:"Parameters_gap-open"`     // Parameters_gap-open
	GapExtend   int      `xml:"Parameters_gap-extend"`   // Parameters_gap-extend
	Filter      *string  `xml:"Parameters_filter"`       // Parameters_filter?
	PhiPattern  *string  `xml:"Parameters_pattern"`      // Parameters_pattern?
	EntrezQuery *string  `xml:"Parameters_entrez-query"` // Parameters_entrez-query?
}

/*
<!ELEMENT Statistics (
        Statistics_db-num,
        Statistics_db-len,
        Statistics_hsp-len,
        Statistics_eff-space,
        Statistics_kappa,
        Statistics_lambda,
        Statistics_entropy)>

<!-- Number of sequences in BLAST db -->
<!ELEMENT Statistics_db-num (%INTEGER;)>

<!-- Length of BLAST db -->
<!ELEMENT Statistics_db-len (%INTEGER;)>

<!-- Effective HSP length -->
<!ELEMENT Statistics_hsp-len (%INTEGER;)>

<!-- Effective search space -->
<!ELEMENT Statistics_eff-space (%REAL;)>

<!-- Karlin-Altschul parameter K -->
<!ELEMENT Statistics_kappa (%REAL;)>

<!-- Karlin-Altschul parameter Lambda -->
<!ELEMENT Statistics_lambda (%REAL;)>

<!-- Karlin-Altschul parameter H -->
<!ELEMENT Statistics_entropy (%REAL;)>
*/

// A Statistics holds the search and database statistics for a Blast result.
type Statistics struct {
	DbNum    int     `xml:"Statistics_db-num"`    // Statistics_db-num
	DbLen    int64   `xml:"Statistics_db-len"`    // Statistics_db-len
	HspLen   int     `xml:"Statistics_hsp-len"`   // Statistics_hsp-len
	EffSpace float64 `xml:"Statistics_eff-space"` // Statistics_eff-space
	Kappa    float64 `xml:"Statistics_kappa"`     // Statistics_kappa
	Lambda   float64 `xml:"Statistics_lambda"`    // Statistics_lambda
	Entropy  float64 `xml:"Statistics_entropy"`   // Statistics_entropy
}

/*
<!ELEMENT Hit (
        Hit_num,
        Hit_id,
        Hit_def,
        Hit_accession,
        Hit_len,
        Hit_hsps?)>

<!-- hit number -->
<!ELEMENT Hit_num (%INTEGER;)>

<!-- SeqId of subject -->
<!ELEMENT Hit_id (#PCDATA)>

<!-- definition line of subject -->
<!ELEMENT Hit_def (#PCDATA)>

<!-- accession -->
<!ELEMENT Hit_accession (#PCDATA)>

<!-- length of subject -->
<!ELEMENT Hit_len (%INTEGER;)>

<!-- all HSP regions for the given subject -->
<!ELEMENT Hit_hsps (Hsp*)>
*/

// A Hit holds the details of an individual Blast hit.
type Hit struct {
	N         int    `xml:"Hit_num"`       // Hit_num
	Id        string `xml:"Hit_id"`        // Hit_id
	Def       string `xml:"Hit_def"`       // Hit_def
	Accession string `xml:"Hit_accession"` // Hit_accession
	Len       int    `xml:"Hit_len"`       // Hit_len
	Hsps      []Hsp  `xml:"Hit_hsps>Hsp"`  // Hit_hsps?
}

/*
<!ELEMENT Hsp (
        Hsp_num,
        Hsp_bit-score,
        Hsp_score,
        Hsp_evalue,
        Hsp_query-from,
        Hsp_query-to,
        Hsp_hit-from,
        Hsp_hit-to,
        Hsp_pattern-from?,
        Hsp_pattern-to?,
        Hsp_query-frame?,
        Hsp_hit-frame?,
        Hsp_identity?,
        Hsp_positive?,
        Hsp_gaps?,
        Hsp_align-len?,
        Hsp_density?,
        Hsp_qseq,
        Hsp_hseq,
        Hsp_midline?)>

<!-- HSP number -->
<!ELEMENT Hsp_num (%INTEGER;)>

<!-- score (in bits) of HSP -->
<!ELEMENT Hsp_bit-score (%REAL;)>

<!-- score of HSP -->
<!ELEMENT Hsp_score (%REAL;)>

<!-- e-value of HSP -->
<!ELEMENT Hsp_evalue (%REAL;)>

<!-- start of HSP in query -->
<!ELEMENT Hsp_query-from (%INTEGER;)>

<!-- end of HSP -->
<!ELEMENT Hsp_query-to (%INTEGER;)>

<!-- start of HSP in subject -->
<!ELEMENT Hsp_hit-from (%INTEGER;)>

<!-- end of HSP in subject -->
<!ELEMENT Hsp_hit-to (%INTEGER;)>

<!-- start of PHI-BLAST pattern -->
<!ELEMENT Hsp_pattern-from (%INTEGER;)>

<!-- end of PHI-BLAST pattern -->
<!ELEMENT Hsp_pattern-to (%INTEGER;)>

<!-- translation frame of query -->
<!ELEMENT Hsp_query-frame (%INTEGER;)>

<!-- translation frame of subject -->
<!ELEMENT Hsp_hit-frame (%INTEGER;)>

<!-- number of identities in HSP -->
<!ELEMENT Hsp_identity (%INTEGER;)>

<!-- number of positives in HSP -->
<!ELEMENT Hsp_positive (%INTEGER;)>

<!-- number of gaps in HSP -->
<!ELEMENT Hsp_gaps (%INTEGER;)>

<!-- length of the alignment used -->
<!ELEMENT Hsp_align-len (%INTEGER;)>

<!-- score density -->
<!ELEMENT Hsp_density (%INTEGER;)>

<!-- alignment string for the query (with gaps) -->
<!ELEMENT Hsp_qseq (#PCDATA)>

<!-- alignment string for subject (with gaps) -->
<!ELEMENT Hsp_hseq (#PCDATA)>

<!-- formating middle line -->
<!ELEMENT Hsp_midline (#PCDATA)>
*/

// An Hsp holds the details of an individual Blast high scoring pair.
type Hsp struct {
	N              int     `xml:"Hsp_num"`          // Hsp_num
	BitScore       float64 `xml:"Hsp_bit-score"`    // Hsp_bit-score
	Score          float64 `xml:"Hsp_score"`        // Hsp_score
	EValue         float64 `xml:"Hsp_evalue"`       // Hsp_evalue
	QueryFrom      int     `xml:"Hsp_query-from"`   // Hsp_query-from
	QueryTo        int     `xml:"Hsp_query-to"`     // Hsp_query-to
	HitFrom        int     `xml:"Hsp_hit-from"`     // Hsp_hit-from
	HitTo          int     `xml:"Hsp_hit-to"`       // Hsp_hit-to
	PhiPatternFrom *int    `xml:"Hsp_pattern-from"` // Hsp_pattern-from?
	PhiPatternTo   *int    `xml:"Hsp_pattern-to"`   // Hsp_pattern-to?
	QueryFrame     *int    `xml:"Hsp_query-frame"`  // Hsp_query-frame?
	HitFrame       *int    `xml:"Hsp_hit-frame"`    // Hsp_hit-frame?
	HspIdentity    *int    `xml:"Hsp_identity"`     // Hsp_identity?
	HspPositive    *int    `xml:"Hsp_positive"`     // Hsp_positive?
	HspGaps        *int    `xml:"Hsp_gaps"`         // Hsp_gaps?
	AlignLen       *int    `xml:"Hsp_align-len"`    // Hsp_align-len?
	Density        *int    `xml:"Hsp_density"`      // Hsp_density?
	QuerySeq       []byte  `xml:"Hsp_qseq"`         // Hsp_qseq
	SubjectSeq     []byte  `xml:"Hsp_hseq"`         // Hsp_hseq
	FormatMidline  []byte  `xml:"Hsp_midline"`      // Hsp_midline?
}
