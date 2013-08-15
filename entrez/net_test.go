// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"code.google.com/p/biogo.ncbi/entrez/link"
	. "code.google.com/p/biogo.ncbi/entrez/spell"
	. "code.google.com/p/biogo.ncbi/entrez/summary"

	"flag"
	"io/ioutil"
	check "launchpad.net/gocheck"
	"testing"
	"time"
)

const tool = "biogo.ncbi/entrez-testsuite"

// Helpers
func intPtr(i int) *int          { return &i }
func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }

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
	time.Sleep(3 * time.Second)
	c.Check(count < 10, check.Equals, true)
}

var net = flag.String("net", "", "Runs tests involving network connections if given an email address.")

func (s *S) TestDoInfo(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	i, err := DoInfo("", tool, *net)
	c.Check(err, check.Equals, nil)
	var (
		// These are the databases we expect to exist always.
		// Others seem to come and go with frightening regularity.
		// This is a reasonable sample, and not intended to be exhaustive.
		dbListCore = map[string]struct{}{
			"books":      struct{}{},
			"genome":     struct{}{},
			"homologene": struct{}{},
			"nuccore":    struct{}{},
			"nucest":     struct{}{},
			"nucleotide": struct{}{},
			"protein":    struct{}{},
			"pubmed":     struct{}{},
			"structure":  struct{}{},
			"taxonomy":   struct{}{},
			"unigene":    struct{}{},
			"unists":     struct{}{},
		}
		dbListRetrieved = make(map[string]struct{})
	)
	for _, db := range i.DbList {
		if _, ok := dbListCore[db]; ok {
			dbListRetrieved[db] = struct{}{}
		}
	}

	// We don't trust this value.
	i.DbList = nil

	c.Check(dbListRetrieved, check.DeepEquals, dbListCore)
	c.Check(i, check.DeepEquals, &Info{DbInfo: nil, Err: ""})
}

func (s *S) TestDoSearch(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	sr, err := DoSearch("nuccore", "hox", nil, nil, tool, *net)
	c.Check(err, check.Equals, nil)
	c.Check(sr, check.Not(check.Equals), nil)
}

func (s *S) TestDoPost(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	p, err := DoPost("protein", tool, *net, nil, 15718680, 157427902, 119703751)
	c.Check(err, check.Equals, nil)
	c.Assert(p.History, check.NotNil)
	c.Check(p.QueryKey, check.DeepEquals, 1)
	c.Check(p.WebEnv, check.Matches, "NCID_[0-9]+_.*")
}

func (s *S) TestFetch(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	for i, t := range []struct { //db, tool, email string, id ...int
		db      string
		rettype string
		retmode string
		ids     []int
		expect  string
	}{
		{
			"protein", "fasta", "text", []int{15718680, 157427902, 119703751},
			"" +
				">gi|15718680|ref|NP_005537.3| tyrosine-protein kinase ITK/TSK [Homo sapiens]\n" +
				"MNNFILLEEQLIKKSQQKRRTSPSNFKVRFFVLTKASLAYFEDRHGKKRTLKGSIELSRIKCVEIVKSDI\n" +
				"SIPCHYKYPFQVVHDNYLLYVFAPDRESRQRWVLALKEETRNNNSLVPKYHPNFWMDGKWRCCSQLEKLA\n" +
				"TGCAQYDPTKNASKKPLPPTPEDNRRPLWEPEETVVIALYDYQTNDPQELALRRNEEYCLLDSSEIHWWR\n" +
				"VQDRNGHEGYVPSSYLVEKSPNNLETYEWYNKSISRDKAEKLLLDTGKEGAFMVRDSRTAGTYTVSVFTK\n" +
				"AVVSENNPCIKHYHIKETNDNPKRYYVAEKYVFDSIPLLINYHQHNGGGLVTRLRYPVCFGRQKAPVTAG\n" +
				"LRYGKWVIDPSELTFVQEIGSGQFGLVHLGYWLNKDKVAIKTIREGAMSEEDFIEEAEVMMKLSHPKLVQ\n" +
				"LYGVCLEQAPICLVFEFMEHGCLSDYLRTQRGLFAAETLLGMCLDVCEGMAYLEEACVIHRDLAARNCLV\n" +
				"GENQVIKVSDFGMTRFVLDDQYTSSTGTKFPVKWASPEVFSFSRYSSKSDVWSFGVLMWEVFSEGKIPYE\n" +
				"NRSNSEVVEDISTGFRLYKPRLASTHVYQIMNHCWKERPEDRPAFSRLLRQLAEIAESGL\n" +
				"\n" +
				">gi|157427902|ref|NP_001098858.1| tyrosine-protein kinase ITK/TSK [Bos taurus]\n" +
				"MNNFILLEEQLIKKSQQKRRTSPSNFKVRFFVLTKTSLAYFEDRHGKKRTLKGSIELSRIKCVEIVKSDI\n" +
				"IIPCHYKYPFQVVHDNYLLYVFAPDRESRQRWVLALKEETRNNNSLVPKYHPNFWLDGRWRCCAQMEKLA\n" +
				"VGCAQYDPTKNASKKPLPPTPEDNRRSLRELEETVVIALYDYQTNDPQELMLQRNEEYYLLDSSEIHWWR\n" +
				"VQDRNGHEGYVPSSYLVEKSPNNLETYEWYNKNISRDKAEKLLLDTGKEGAFMVRDSRTPGTYTVSVFTK\n" +
				"AIVSENNPCIKHYHIKETNDNPKRYYVAEKYVFDSIPLLINYHQHNGGGLVTRLRYPVCSWRQKAPVTAG\n" +
				"LRYGKWVIDPSELTFVQEIGSGQFGLVHLGYWLNKDKVAIKTIQEGAMSEEDFIEEAEVMMKLSHPKLVQ\n" +
				"LYGVCLEQAPICLVFEFMEHGCLSDYLRSQRGLFAAETLLGMCLDVCEGMAYLEEACVIHRDLAARNCLV\n" +
				"GENQVIKVSDFGMTRFVLDDQYTSSTGTKFPVKWASPEVFSFSRYSSKSDVWSFGVLMWEVFSEGKIPYE\n" +
				"NRSNSEVVEDITTGFRLYKPRLASQHIYQIMNHCWKEKPEDRPPFSRLLSQLAEIAELGL\n" +
				"\n" +
				">gi|119703751|ref|NP_034713.2| tyrosine-protein kinase ITK/TSK isoform 2 [Mus musculus]\n" +
				"MNNFILLEEQLIKKSQQKRRTSPSNFKVRFFVLTKASLAYFEDRHGKKRTLKGSIELSRIKCVEIVKSDI\n" +
				"SIPCHYKYPFQVVHDNYLLYVFAPDCESRQRWVLTLKEETRNNNSLVSKYHPNFWMDGRWRCCSQLEKPA\n" +
				"VGCAPYDPSKNASKKPLPPTPEDNRRSFQEPEETLVIALYDYQTNDPQELALRCDEEYYLLDSSEIHWWR\n" +
				"VQDKNGHEGYAPSSYLVEKSPNNLETYEWYNKSISRDKAEKLLLDTGKEGAFMVRDSRTPGTYTVSVFTK\n" +
				"AIISENPCIKHYHIKETNDSPKRYYVAEKYVFDSIPLLIQYHQYNGGGLVTRLRYPVCSWRQKAPVTAGL\n" +
				"RYGKWVIQPSELTFVQEIGSGQFGLVHLGYWLNKDKVAIKTIQEGAMSEEDFIEEAEVMMKLSHPKLVQL\n" +
				"YGVCLEQAPICLVFEFMEHGCLSDYLRSQRGLFAAETLLGMCLDVCEGMAYLEKACVIHRDLAARNCLVG\n" +
				"ENQVIKVSDFGMTRFVLDDQYTSSTGTKFPVKWASPEVFSFSRYSSKSDVWSFGVLMWEVFSEGKIPYEN\n" +
				"RSNSEVVEDISTGFRLYKPRLASCHVYQIMNHCWKEKPEDRPPFSQLLSQLAEIAEAGL\n" +
				"\n",
		},
	} {
		rc, err := Fetch(t.db, &Parameters{RetMode: t.retmode, RetType: t.rettype}, tool, *net, nil, t.ids...)
		c.Assert(err, check.Equals, nil, check.Commentf("Test %d", i))
		b, err := ioutil.ReadAll(rc)
		rc.Close()
		c.Assert(err, check.Equals, nil)
		c.Check(string(b), check.Equals, t.expect, check.Commentf("Test %d", i))
	}
}

func (s *S) TestDoSummary(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	sum, err := DoSummary("protein", nil, tool, *net, nil, 15718680, 157427902, 119703751)
	c.Check(err, check.Equals, nil)
	c.Check(sum.Database, check.Equals, "protein")
	expect := &Summary{
		Database: "protein",
		Documents: []Document{
			Document{
				Id: 15718680,
				Items: []Item{
					{Value: "NP_005537", Name: "Caption", Type: "String"},
					{Value: "tyrosine-protein kinase ITK/TSK [Homo sapiens]", Name: "Title", Type: "String"},
					{Value: "gi|15718680|ref|NP_005537.3|[15718680]", Name: "Extra", Type: "String"},
					{Value: "15718680", Name: "Gi", Type: "Integer"},
					{Value: "1999/06/09", Name: "CreateDate", Type: "String"},
					{Value: "not tested", Name: "UpdateDate", Type: "String"},
					{Value: "512", Name: "Flags", Type: "Integer"},
					{Value: "9606", Name: "TaxId", Type: "Integer"},
					{Value: "620", Name: "Length", Type: "Integer"},
					{Value: "live", Name: "Status", Type: "String"},
					{Value: "", Name: "ReplacedBy", Type: "String"},
					{Value: "  ", Name: "Comment", Type: "String"},
				},
			},
			Document{
				Id: 157427902,
				Items: []Item{
					{Value: "NP_001098858", Name: "Caption", Type: "String"},
					{Value: "tyrosine-protein kinase ITK/TSK [Bos taurus]", Name: "Title", Type: "String"},
					{Value: "gi|157427902|ref|NP_001098858.1|[157427902]", Name: "Extra", Type: "String"},
					{Value: "157427902", Name: "Gi", Type: "Integer"},
					{Value: "2007/09/24", Name: "CreateDate", Type: "String"},
					{Value: "not tested", Name: "UpdateDate", Type: "String"},
					{Value: "512", Name: "Flags", Type: "Integer"},
					{Value: "9913", Name: "TaxId", Type: "Integer"},
					{Value: "620", Name: "Length", Type: "Integer"},
					{Value: "live", Name: "Status", Type: "String"},
					{Value: "", Name: "ReplacedBy", Type: "String"},
					{Value: "  ", Name: "Comment", Type: "String"},
				},
			},
			Document{
				Id: 119703751,
				Items: []Item{
					{Value: "NP_034713", Name: "Caption", Type: "String"},
					{Value: "tyrosine-protein kinase ITK/TSK isoform 2 [Mus musculus]", Name: "Title", Type: "String"},
					{Value: "gi|119703751|ref|NP_034713.2|[119703751]", Name: "Extra", Type: "String"},
					{Value: "119703751", Name: "Gi", Type: "Integer"},
					{Value: "2000/01/25", Name: "CreateDate", Type: "String"},
					{Value: "not tested", Name: "UpdateDate", Type: "String"},
					{Value: "512", Name: "Flags", Type: "Integer"},
					{Value: "10090", Name: "TaxId", Type: "Integer"},
					{Value: "619", Name: "Length", Type: "Integer"},
					{Value: "live", Name: "Status", Type: "String"},
					{Value: "", Name: "ReplacedBy", Type: "String"},
					{Value: "  ", Name: "Comment", Type: "String"},
				},
			},
		},
	}
	for i, d := range sum.Documents {
		c.Check(d.Id, check.Equals, expect.Documents[i].Id)
		for j, it := range d.Items {
			if it.Name == "UpdateDate" {
				continue
			}
			c.Check(it, check.Equals, expect.Documents[i].Items[j])
		}
	}
}

func (s *S) TestDoLink(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	for i, t := range []struct {
		db     string
		fromDb string
		cmd    string
		query  string
		params *Parameters
		hist   *History
		ids    [][]int
		expect *Link
	}{
		{
			"protein", "gene", "", "", nil, nil, [][]int{{15718680, 157427902}},
			&Link{
				LinkSets: []link.LinkSet{
					{
						DbFrom: "protein",
						IdList: []link.Id{
							{Id: 15718680},
							{Id: 157427902},
						},
						Neighbor: []link.LinkSetDb{
							{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []link.Link{
									{Id: link.Id{Id: 522311}, Score: nil},
									{Id: link.Id{Id: 3702}, Score: nil},
								},
							},
						},
					},
				},
				Err: nil,
			},
		},
		{
			"protein", "gene", "", "", nil, nil, [][]int{{15718680}, {157427902}},
			&Link{
				LinkSets: []link.LinkSet{
					{
						DbFrom: "protein",
						IdList: []link.Id{
							{Id: 15718680},
						},
						Neighbor: []link.LinkSetDb{
							{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []link.Link{
									{Id: link.Id{Id: 3702}, Score: nil},
								},
							},
						},
					},
					{
						DbFrom: "protein",
						IdList: []link.Id{
							{Id: 157427902},
						},
						Neighbor: []link.LinkSetDb{
							{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []link.Link{
									{Id: link.Id{Id: 522311}, Score: nil},
								},
							},
						},
					},
				},
				Err: nil,
			},
		},
	} {
		l, err := DoLink(t.db, t.fromDb, t.cmd, t.query, t.params, tool, *net, t.hist, t.ids...)
		c.Check(err, check.Equals, nil, check.Commentf("Test %d", i))
		c.Check(l, check.DeepEquals, t.expect, check.Commentf("Test %d", i))
	}
}

func (s *S) TestDoGlobal(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	g, err := DoGlobal("toolkit", tool, *net)
	c.Check(err, check.Equals, nil)
	c.Check(g, check.Not(check.Equals), nil)
}

func (s *S) TestDoSpell(c *check.C) {
	if *net == "" {
		c.Skip("Network tests not requested.")
	}
	sp, err := DoSpell("", "asthmaa OR alergies", tool, *net)
	c.Check(err, check.Equals, nil)
	c.Check(sp, check.DeepEquals, &Spell{
		Database:  "pubmed",
		Query:     "asthmaa OR alergies",
		Corrected: "asthma or allergies",
		Replace: []Replacement{
			New("asthma"),
			Old(" OR "),
			New("allergies"),
		},
		Err: "",
	})
}
