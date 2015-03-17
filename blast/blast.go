// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Documentation from http://www.ncbi.nlm.nih.gov/books/NBK21097/

// Package blast provides support for interaction with the NCBI BLAST service.
//
// Please see http://blast.ncbi.nlm.nih.gov/Blast.cgi?CMD=Web&PAGE_TYPE=BlastDocs&DOC_TYPE=DeveloperInfo
// for the BLAST service usage policy.
//
// Required parameters are specified by name in the function call. Optional parameters are
// passed via parameter struct values. See the 'QBlast URL API User's Guide' at
// http://www.ncbi.nlm.nih.gov/BLAST/Doc/urlapi.html for explanation of the use of these
// programs.
//
// The following two parameters should be included in all BLAST requests.
//
//  tool   Name of application making the BLAST call. Its value must be a string with no
//         internal spaces.
//
//  email  E-mail address of the BLAST user. Its value must be a string with no internal
//         spaces, and should be a valid e-mail address.
package blast

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/biogo/ncbi/ncbi"
)

var (
	ErrNoRidProvided = errors.New("blast: no RID provided")
	ErrMissingRid    = errors.New("blast: missing RID/RTOE field")
	ErrMissingStatus = errors.New("blast: missing Status field")
)

// Limit is a package level limit on requests that can be sent to the BLAST server. This
// limit is mandated by the BLAST service usage policy. Limit is exported to allow reuse
// of http.Requests provided by RequestWebReadCloser without overrunning the BLAST request limit.
// Changing the the value of Limit to allow more frequent requests may result in IP blocking
// by the BLAST servers.
var Limit = ncbi.NewLimiter(3 * time.Second)

const cmdParam = "CMD" // parameter CMD

// PutParameters is used to pass optional parameters to the Put command. The relevant documentation
// for each of these parameters is at http://www.ncbi.nlm.nih.gov/BLAST/Doc/node9.html.
type PutParameters struct {
	AutoFormat                 string   `param:"AUTO_FORMAT"`
	CompositionBasedStatistics bool     `param:"COMPOSITION_BASED_STATISTICS"`
	Database                   string   `param:"DATABASE"`
	DbGeneticCode              []int    `param:"DB_GENETIC_CODE"`
	EndPoints                  bool     `param:"ENDPOINTS"`
	EntrezQuery                string   `param:"ENTREZ_QUERY"`
	Expect                     *float64 `param:"EXPECT"`
	Filter                     string   `param:"FILTER"`
	GapCosts                   [2]int   `param:"GAPCOSTS"`
	GeneticCode                []int    `param:"GENETIC_CODE"`
	HitListSize                int      `param:"HITLIST_SIZE"`
	IThresh                    float64  `param:"I_THRESH"`
	Layout                     string   `param:"LAYOUT"`
	LCaseMask                  bool     `param:"LCASE_MASK"`
	Megablast                  bool     `param:"MEGABLAST"`
	MatrixName                 string   `param:"MATRIX_NAME"`
	NuclPenalty                int      `param:"NUCL_PENALTY"`
	NuclReward                 int      `param:"NUCL_REWARD"`
	OtherAdvanced              string   `param:"OTHER_ADVANCED"`
	PercIdent                  int      `param:"PERC_IDENT"`
	PhiPattern                 string   `param:"PHI_PATTERN"`
	Program                    string   `param:"PROGRAM"`
	Pssm                       string   `param:"PSSM"`
	QueryFile                  string   `param:"QUERY_FILE"`
	QueryBelieveDefline        bool     `param:"QUERY_BELIEVE_DEFLINE"`
	QueryFrom                  int      `param:"QUERY_FROM"`
	QueryTo                    int      `param:"QUERY_TO"`
	ResultsFile                bool     `param:"RESULTS_FILE"`
	SearchspEff                int      `param:"SEARCHSP_EFF"`
	Service                    string   `param:"SERVICE"`
	Threshold                  int      `param:"THRESHOLD"`
	UngappedAlignment          bool     `param:"UNGAPPED_ALIGNMENT"`
	WordSize                   int      `param:"WORD_SIZE"`
}

// GetParameters is used to pass optional parameters to the Get command. The relevant documentation
// for each of these parameters is at http://www.ncbi.nlm.nih.gov/BLAST/Doc/node9.html.
type GetParameters struct {
	FormatType string `param:"FORMAT_TYPE"` // Ignored by GetOutput: "HTML", "Text", "ASN.1" or "XML".

	Alignments           int     `param:"ALIGNMENTS"`
	AlignmentView        string  `param:"ALIGNMENT_VIEW"`
	Descriptions         int     `param:"DESCRIPTIONS"`
	EntrezLinksNewWindow bool    `param:"ENTREZ_LINKS_NEW_WINDOW"`
	ExpectLow            float64 `param:"EXPECT_LOW"`
	ExpectHigh           float64 `param:"EXPECT_HIGH"`
	FormatEntrezQuery    string  `param:"FORMAT_ENTREZ_QUERY"`
	FormatObject         string  `param:"FORMAT_OBJECT"`
	NcbiGi               bool    `param:"NCBI_GI"`
	ResultsFile          bool    `param:"RESULTS_FILE"`
	Service              string  `param:"SERVICE"`
	ShowOverview         *bool   `param:"SHOW_OVERVIEW"`
}

// WebParameters is used to pass optional parameters to the Web command. The relevant documentation
// for each of these parameters is at http://www.ncbi.nlm.nih.gov/BLAST/Doc/node9.html.
// Note there is inadequate documentation for what parameters the Web command accepts, so all are included.
type WebParameters struct {
	Alignments                 int      `param:"ALIGNMENTS"`
	AlignmentView              string   `param:"ALIGNMENT_VIEW"`
	AutoFormat                 string   `param:"AUTO_FORMAT"`
	Cmd                        string   `param:"CMD"`
	CompositionBasedStatistics bool     `param:"COMPOSITION_BASED_STATISTICS"`
	Database                   string   `param:"DATABASE"`
	DbGeneticCode              []int    `param:"DB_GENETIC_CODE"`
	Descriptions               int      `param:"DESCRIPTIONS"`
	DocType                    string   `param:"DOC_TYPE"`
	EndPoints                  bool     `param:"ENDPOINTS"`
	EntrezLinksNewWindow       bool     `param:"ENTREZ_LINKS_NEW_WINDOW"`
	EntrezQuery                string   `param:"ENTREZ_QUERY"`
	Expect                     *float64 `param:"EXPECT"`
	ExpectLow                  float64  `param:"EXPECT_LOW"`
	ExpectHigh                 float64  `param:"EXPECT_HIGH"`
	Filter                     string   `param:"FILTER"`
	FormatEntrezQuery          string   `param:"FORMAT_ENTREZ_QUERY"`
	FormatObject               string   `param:"FORMAT_OBJECT"`
	FormatType                 string   `param:"FORMAT_TYPE"`
	GapCosts                   [2]int   `param:"GAPCOSTS"`
	GeneticCode                []int    `param:"GENETIC_CODE"`
	HitListSize                int      `param:"HITLIST_SIZE"`
	IThresh                    float64  `param:"I_THRESH"`
	Layout                     string   `param:"LAYOUT"`
	LCaseMask                  bool     `param:"LCASE_MASK"`
	Megablast                  bool     `param:"MEGABLAST"`
	MatrixName                 string   `param:"MATRIX_NAME"`
	NcbiGi                     bool     `param:"NCBI_GI"`
	NuclPenalty                int      `param:"NUCL_PENALTY"`
	NuclReward                 int      `param:"NUCL_REWARD"`
	OtherAdvanced              string   `param:"OTHER_ADVANCED"`
	PageType                   string   `param:"PAGE_TYPE"`
	PercIdent                  int      `param:"PERC_IDENT"`
	PhiPattern                 string   `param:"PHI_PATTERN"`
	Program                    string   `param:"PROGRAM"`
	Pssm                       string   `param:"PSSM"`
	Query                      string   `param:"QUERY"`
	QueryFile                  string   `param:"QUERY_FILE"`
	QueryBelieveDefline        bool     `param:"QUERY_BELIEVE_DEFLINE"`
	QueryFrom                  int      `param:"QUERY_FROM"`
	QueryTo                    int      `param:"QUERY_TO"`
	Rid                        string   `param:"RID"`
	ResultsFile                bool     `param:"RESULTS_FILE"`
	SearchspEff                int      `param:"SEARCHSP_EFF"`
	Service                    string   `param:"SERVICE"`
	ShowOverview               *bool    `param:"SHOW_OVERVIEW"`
	Threshold                  int      `param:"THRESHOLD"`
	UngappedAlignment          bool     `param:"UNGAPPED_ALIGNMENT"`
	WordSize                   int      `param:"WORD_SIZE"`
}

// BlastUri is the base URL for the NCBI BLAST URL API.
const BlastUri = ncbi.Util("http://www.ncbi.nlm.nih.gov/blast/Blast.cgi")

// fillParams adds elements to v based on the "param" tag of p if the value is not the
// zero value for that type.
func fillParams(cmd string, p interface{}, v url.Values) {
	defer func() {
		v[cmdParam] = []string{cmd}
	}()
	pv := reflect.ValueOf(p)
	if pv.IsNil() {
		return
	}
	pv = pv.Elem()
	n := pv.NumField()
	t := pv.Type()
	for i := 0; i < n; i++ {
		tf := t.Field(i)
		if tf.PkgPath != "" {
			continue
		}
		tag := tf.Tag.Get("param")
		if tag != "" {
			in := pv.Field(i).Interface()
			switch cv := in.(type) {
			case int:
				if cv != 0 {
					v[tag] = []string{fmt.Sprint(cv)}
				}
			case float64:
				if cv != 0 {
					v[tag] = []string{fmt.Sprint(cv)}
				}
			case *float64:
				if cv != nil {
					v[tag] = []string{fmt.Sprint(*cv)}
				}
			case string:
				if cv != "" {
					v[tag] = []string{cv}
				}
			case bool:
				if cv {
					v[tag] = []string{"yes"}
				}
			case [2]int:
				if cv != [2]int{} {
					v[tag] = []string{fmt.Sprintf("%d %d", cv[0], cv[1])}
				}
			case []int:
				if cv != nil {
					s := make([]string, len(cv))
					for i, c := range cv {
						s[i] = fmt.Sprint(c)
					}
					v[tag] = []string{strings.Join(s, ",")}
				}
			case *bool:
				if cv != nil {
					if *cv {
						v[tag] = []string{"yes"}
					} else {
						v[tag] = []string{"no"}
					}
				}
			default:
				panic("cannot reach")
			}
		}
	}
}

// RequestWebReadCloser returns an io.ReadCloser that reads from the stream returned by a Web request
// of the the given page. It is the responsibility of the caller to close the returned stream.
func RequestWebReadCloser(page string, p *WebParameters, tool, email string) (io.ReadCloser, error) {
	v := url.Values{}
	fillParams("Web", p, v)
	if page != "" {
		v["PAGE"] = []string{page}
	}
	resp, err := BlastUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Put submits a request for a BLAST job to the NCBI BLAST server and returns the associated
// Rid containing the RID for the request.
func Put(query string, p *PutParameters, tool, email string) (*Rid, error) {
	v := url.Values{}
	if query != "" {
		v["QUERY"] = []string{query}
	}
	fillParams("Put", p, v)
	rid := Rid{}
	resp, err := BlastUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	err = rid.unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return &rid, nil
}

// SearchInfo returns status information for the search request corresponding to r.
func (r *Rid) SearchInfo(tool, email string) (*SearchInfo, error) {
	v := url.Values{}
	if r.rid != "" {
		v["RID"] = []string{r.rid}
	} else {
		return nil, ErrNoRidProvided
	}
	v[cmdParam] = []string{"Get"}
	v["FORMAT_OBJECT"] = []string{"SearchInfo"}
	<-r.Ready()
	resp, err := BlastUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	s := SearchInfo{Rid: r}
	err = s.unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetOutput returns an Output filled with data obtained from an Get request for the request
// corresponding to r.
func (r *Rid) GetOutput(p *GetParameters, tool, email string) (*Output, error) {
	v := url.Values{}
	if r.rid != "" {
		v["RID"] = []string{r.rid}
	} else {
		return nil, ErrNoRidProvided
	}
	fillParams("Get", p, v)
	v["FORMAT_TYPE"] = []string{"XML"}
	o := Output{}
	r.limit.Wait()
	err := BlastUri.GetXML(v, tool, email, Limit, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetReadCloser returns an io.ReadCloser that reads from the stream returned by a Get request
// corresponding to r. It is the responsibility of the caller to close the returned stream.
func (r *Rid) GetReadCloser(p *GetParameters, tool, email string) (io.ReadCloser, error) {
	v := url.Values{}
	if r.rid != "" {
		v["RID"] = []string{r.rid}
	} else {
		return nil, ErrNoRidProvided
	}
	fillParams("Get", p, v)
	r.limit.Wait()
	resp, err := BlastUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete deletes the the request and results corresponding to r from the NCBI BLAST server.
func (r *Rid) Delete(tool, email string) error {
	v := url.Values{}
	if r.rid != "" {
		v["RID"] = []string{r.rid}
	} else {
		return ErrNoRidProvided
	}
	v[cmdParam] = []string{"Delete"}
	resp, err := BlastUri.Get(v, tool, email, Limit)
	if err != nil {
		return err
	}
	return resp.Close()
}

// RequestInfo returns an Info with up-to-date information about NCBI BLAST services.
func RequestInfo(target string, tool, email string) (*Info, error) {
	v := url.Values{}
	if target != "" {
		v["TARGET"] = []string{target}
	}
	v[cmdParam] = []string{"Info"}
	var i Info
	resp, err := BlastUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	err = i.unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
