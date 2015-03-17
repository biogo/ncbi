// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// All documentation in this package should agree with documents in the
// 'Entrez Programming Utilities Help': http://www.ncbi.nlm.nih.gov/books/NBK25501/
// In case of disagreement that document is correct unless otherwise specified.

// Please people, let's start using JSON for these kinds of things.

// Package entrez provides support for interaction with the NCBI Entrez Utility Programs (E-utilities).
//
// Please see http://www.ncbi.nlm.nih.gov/books/n/helpeutils/chapter2/ for the E-utility
// usage policy.
//
// Required parameters are specified by name in the function call. Optional parameters are
// passed via Parameter and History values. See the 'Entrez Programming Utilities Help' at
// http://www.ncbi.nlm.nih.gov/books/NBK25501/ for detailed explanation of the use of these
// programs.
//
// The following two parameters should be included in all E-utility requests.
//
//  tool   Name of application making the E-utility call. Its value must be a string with no
//         internal spaces.
//
//  email  E-mail address of the E-utility user. Its value must be a string with no internal
//         spaces, and should be a valid e-mail address.
package entrez

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/biogo/ncbi/ncbi"
)

// The E-utilities default to "pubmed". Some functions mark which db was used because E-utilities
// don't, so this is needed.
const defaultDb = "pubmed"

// Limit is a package level limit on requests that can be sent to the Entrez server. This
// limit is mandated by chapter 2 of the E-utilities manual. Limit is exported to allow reuse
// of http.Requests provided by NewRequest without overrunning the Entrez request limit.
// Changing the the value of Limit to allow more frequent requests may result in IP blocking
// by the Entrez servers.
var Limit = ncbi.NewLimiter(time.Second / 3)

var (
	ErrNoIdProvided = errors.New("entrez: no id provided")
	ErrNoQuery      = errors.New("entrez: no query")
)

// Parameters is used to pass optional parameters to E-utility programs. The relevant documentation
// for each of these parameters is at http://www.ncbi.nlm.nih.gov/books/n/helpeutils/chapter4/.
type Parameters struct {
	RetMode    string `param:"retmode"`
	RetType    string `param:"rettype"`
	RetStart   int    `param:"retstart"`
	RetMax     int    `param:"retmax"`
	Strand     int    `param:"strand"`
	SeqStart   int    `param:"seqstart"`
	SeqStop    int    `param:"seqstop"`
	Complexity int    `param:"complexity"`
	LinkName   string `param:"linkname"`
	Holding    string `param:"holding"`
	DateType   string `param:"datetype"`
	RelDate    string `param:"reldate"`
	MinDate    string `param:"mindate"`
	MaxDate    string `param:"maxdate"`
}

// History stores an Entrez Web Environment and query key. The zero values of QueryKey and WebEnv
// indicate unset values.
type History struct {
	QueryKey int    `xml:"QueryKey"`
	WebEnv   string `xml:"WebEnv"`
}

const (
	// Base is the base URL for the NCBI Entrez Programming Utilities (E-utilities) API.
	Base = "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/"

	//  * Provides a list of the names of all valid Entrez databases.
	//  * Provides statistics for a single database, including lists of indexing fields and available
	//    link names.
	InfoUri = ncbi.Util(Base + "einfo.fcgi")

	//  * Provides a list of UIDs matching a text query.
	//  * Posts the results of a search on the History server.
	//  * Downloads all UIDs from a dataset stored on the History server.
	//  * Combines or limits UID datasets stored on the History server.
	//  * Sorts sets of UIDs.
	SearchUri = ncbi.Util(Base + "esearch.fcgi")

	//  * Uploads a list of UIDs to the Entrez History server.
	//  * Appends a list of UIDs to an existing set of UID lists attached to a Web Environment.
	PostUri = ncbi.Util(Base + "epost.fcgi")

	//  * Returns document summaries (DocSums) for a list of input UIDs.
	//  * Returns DocSums for a set of UIDs stored on the Entrez History server.
	SummaryUri = ncbi.Util(Base + "esummary.fcgi")

	//  * Returns formatted data records for a list of input UIDs.
	//  * Returns formatted data records for a set of UIDs stored on the Entrez History server.
	FetchUri = ncbi.Util(Base + "efetch.fcgi")

	//  * Returns UIDs linked to an input set of UIDs in either the same or a different Entrez database.
	//  * Returns UIDs linked to other UIDs in the same Entrez database that match an Entrez query.
	//  * Checks for the existence of Entrez links for a set of UIDs within the same database.
	//  * Lists the available links for a UID.
	//  * Lists LinkOut URLs and attributes for a set of UIDs.
	//  * Lists hyperlinks to primary LinkOut providers for a set of UIDs.
	//  * Creates hyperlinks to the primary LinkOut provider for a single UID.
	LinkUri = ncbi.Util(Base + "elink.fcgi")

	//  * Provides the number of records retrieved in all Entrez databases by a single text query.
	GlobalUri = ncbi.Util(Base + "egquery.fcgi")

	//  * Provides spelling suggestions for terms within a single text query in a given database.
	SpellUri = ncbi.Util(Base + "espell.fcgi")

	//  * Retrieves PubMed IDs (PMIDs) that correspond to a set of input citation queries.
	CitMatchUri = ncbi.Util(Base + "ecitmatch.cgi")
)

type unmarshaler interface {
	Unmarshal(io.Reader) error
}

func get(ut ncbi.Util, v url.Values, tool, email string, d interface{}) error {
	return ut.GetXML(v, tool, email, Limit, d)
}

// fillParams adds elements to v based on the "param" tag of p if the value is not the
// zero value for that type.
func fillParams(p *Parameters, v url.Values) {
	if p == nil {
		return
	}
	pv := reflect.ValueOf(p).Elem()
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
			case string:
				if cv != "" {
					v[tag] = []string{cv}
				}
			default:
				panic("cannot reach")
			}
		}
	}
}

// DoInfo returns an Info filled with data obtained from an EInfo query of the specified
// db or all databases if db is an empty string.
func DoInfo(db, tool, email string) (*Info, error) {
	v := url.Values{}
	if db != "" {
		v["db"] = []string{db}
	}
	i := Info{}
	err := get(InfoUri, v, tool, email, &i)
	if err != nil {
		return nil, err
	}
	if i.Err != "" {
		return &i, errors.New(i.Err)
	}
	return &i, nil
}

// DoSearch returns a Search filled with data obtained from an ESearch query of the
// specified db. If h is not nil the search will use the Entrez history server and will
// be filled with the history results of the ESearch query. If h.WebEnv is not empty,
// it will be passed to ESearch as the web environment and if h.QueryKey is not zero,
// it will be passed as the query key.
func DoSearch(db, query string, p *Parameters, h *History, tool, email string) (*Search, error) {
	v := url.Values{}
	if db != "" {
		v["db"] = []string{db}
	} else {
		db = defaultDb
	}
	if query != "" {
		v["term"] = []string{query}
	}
	fillParams(p, v)
	s := Search{Database: db, History: h}
	if h != nil {
		v["usehistory"] = []string{"y"}
		if h.WebEnv != "" {
			v["webenv"] = []string{h.WebEnv}
			if h.QueryKey != 0 {
				v["query_key"] = []string{fmt.Sprint(h.QueryKey)}
			}
		}
	}
	err := get(SearchUri, v, tool, email, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// DoPost returns a Post filled with the response from an EPost action on the specified
// id list. If h is not nil, its WebEnv field is passed as the E-utilies webenv parameter,
// and if h.QueryKey is zero, h will be filled with the history result from the EPost request.
func DoPost(db, tool, email string, h *History, id ...int) (*Post, error) {
	if len(id) == 0 {
		return nil, ErrNoIdProvided
	}
	ids := make([]string, len(id))
	for i, uid := range id {
		ids[i] = fmt.Sprint(uid)
	}
	v := url.Values{"id": []string{strings.Join(ids, ",")}}
	if db != "" {
		v["db"] = []string{db}
	}
	p := Post{}
	if h != nil && h.WebEnv != "" {
		v["webenv"] = []string{h.WebEnv}
	} else if h != nil && h.QueryKey == 0 {
		p.History = h
	}
	err := get(PostUri, v, tool, email, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Fetch returns an io.ReadCloser that reads from the stream returned by an EFetch of the
// the given id list or history. It is the responsibility of the caller to close this.
func Fetch(db string, p *Parameters, tool, email string, h *History, id ...int) (io.ReadCloser, error) {
	if len(id) == 0 && h == nil {
		return nil, ErrNoIdProvided
	}
	ids := make([]string, len(id))
	for i, uid := range id {
		ids[i] = fmt.Sprint(uid)
	}
	v := url.Values{"id": ids}
	if db != "" {
		v["db"] = []string{db}
	}
	fillParams(p, v)
	if h != nil && h.WebEnv != "" && h.QueryKey != 0 {
		v["webenv"] = []string{h.WebEnv}
		v["query_key"] = []string{fmt.Sprint(h.QueryKey)}
	} else if len(id) == 0 {
		return nil, ErrNoIdProvided
	}
	resp, err := FetchUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DoSummary returns a Summary filled with the response from an ESummary query on the specified
// id list. If h is not nil and its fields are non-zero, its field values are passed to ESummary.
// DoSummary returns an error if both h is nil and id has length zero.
func DoSummary(db string, p *Parameters, tool, email string, h *History, id ...int) (*Summary, error) {
	if len(id) == 0 && h == nil {
		return nil, ErrNoIdProvided
	}
	ids := make([]string, len(id))
	for i, uid := range id {
		ids[i] = fmt.Sprint(uid)
	}
	v := url.Values{"id": []string{strings.Join(ids, ",")}}
	if db != "" {
		v["db"] = []string{db}
	} else {
		db = defaultDb
	}
	fillParams(p, v)
	if h != nil && h.WebEnv != "" && h.QueryKey != 0 {
		v["webenv"] = []string{h.WebEnv}
		v["query_key"] = []string{fmt.Sprint(h.QueryKey)}
	} else if len(id) == 0 {
		return nil, ErrNoIdProvided
	}
	s := Summary{Database: db}
	err := get(SummaryUri, v, tool, email, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// DoLink returns a Link filled with the response from an ELink action on the specified
// ids list. If h is not nil and its fields are non-zero, its field values are passed to
// ESummary. DoSummary returns an error if both h is nil and ids has length zero.
func DoLink(fromDb, toDb, cmd, query string, p *Parameters, tool, email string, h *History, ids ...[]int) (*Link, error) {
	if len(ids) == 0 && h == nil {
		return nil, ErrNoIdProvided
	}

	idls := make([]string, len(ids))
	for i, id := range ids {
		if len(id) == 0 {
			continue
		}
		ids := make([]string, len(id))
		for i, uid := range id {
			ids[i] = fmt.Sprint(uid)
		}
		idls[i] = strings.Join(ids, ",")
	}
	v := url.Values{"id": idls}

	if toDb != "" {
		v["db"] = []string{toDb}
	}
	if fromDb != "" {
		v["dbfrom"] = []string{fromDb}
	}
	if cmd != "" {
		v["cmd"] = []string{cmd}
	}
	if query != "" {
		v["term"] = []string{query}
	}
	fillParams(p, v)
	if h != nil && h.WebEnv != "" && h.QueryKey != 0 {
		v["webenv"] = []string{h.WebEnv}
		v["query_key"] = []string{fmt.Sprint(h.QueryKey)}
	} else if len(ids) == 0 {
		return nil, ErrNoIdProvided
	}
	l := Link{}
	err := get(LinkUri, v, tool, email, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// DoGlobal returns a Global filled with the response from an EGQuery query.
func DoGlobal(query, tool, email string) (*Global, error) {
	if query == "" {
		return nil, ErrNoQuery
	}
	v := url.Values{"term": []string{query}}
	g := Global{}
	err := get(GlobalUri, v, tool, email, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// DoSpell returns a Spell filled with the response from an ESpell query.
func DoSpell(db, query string, tool, email string) (*Spell, error) {
	v := url.Values{}
	if db != "" {
		v["db"] = []string{db}
	}
	if query != "" {
		v["term"] = []string{query}
	}
	sp := Spell{}
	err := get(SpellUri, v, tool, email, &sp)
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

// CitQuery represents a single element of a CitMatch citation query.
type CitQuery struct {
	JournalTitle string
	Year         string
	Volume       string
	FirstPage    string
	AuthorName   string
}

// DoCitMatch returns a map[string]int associating keys provided in the query
// to the citations requested in the query. If email is set, the response will
// also be sent to that address.
func DoCitMatch(query map[string]CitQuery, tool, email string) (map[string]int, error) {
	v := url.Values{"db": []string{"pubmed"}, "retmode": []string{"xml"}}
	if query != nil {
		var buf bytes.Buffer
		for key, cit := range query {
			fmt.Fprintf(&buf, "%s|%s|%s|%s|%s|%s|\r",
				cit.JournalTitle,
				cit.Year,
				cit.Volume,
				cit.FirstPage,
				cit.AuthorName,
				key,
			)
		}
		v["bdata"] = []string{buf.String()}
	}
	r, err := CitMatchUri.Get(v, tool, email, Limit)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res := make(map[string]int)
	buf, err := ioutil.ReadAll(r)
	for _, rec := range bytes.Split(buf, []byte{'\n'}) {
		if len(rec) == 0 {
			continue
		}
		f := bytes.Split(rec, []byte{'|'})
		res[string(f[5])], err = strconv.Atoi(string(f[6]))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
