// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	. "code.google.com/p/biogo.entrez/link"
	. "code.google.com/p/biogo.entrez/spell"
	. "code.google.com/p/biogo.entrez/summary"

	"flag"
	"io/ioutil"
	check "launchpad.net/gocheck"
	"testing"
	"time"
)

// Helpers
func intPtr(i int) *int          { return &i }
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

var net = flag.Bool("net", false, "Run tests involving network connections.")

func (s *S) TestDoInfo(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	i, err := DoInfo("", "biogo.entrez-testsuite", "")
	c.Check(err, check.Equals, nil)
	c.Check(i, check.DeepEquals, &Info{
		DbList: []string{
			"pubmed", "protein", "nuccore", "nucleotide",
			"nucgss", "nucest", "structure", "genome",
			"assembly", "gcassembly", "genomeprj", "bioproject",
			"biosample", "biosystems", "blastdbinfo", "books",
			"cdd", "clone", "gap", "gapplus",
			"dbvar", "epigenomics", "gene", "gds",
			"geoprofiles", "homologene", "journals", "medgen",
			"mesh", "ncbisearch", "nlmcatalog", "omia",
			"omim", "pmc", "popset", "probe",
			"proteinclusters", "pcassay", "pccompound", "pcsubstance",
			"pubmedhealth", "seqannot", "snp", "sra",
			"taxonomy", "toolkit", "toolkitall", "toolkitbook",
			"unigene", "unists", "gencoll"},
		DbInfo: nil, Err: nil})
}

func (s *S) TestDoSearch(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	sr, err := DoSearch("nuccore", "hox", nil, nil, "biogo.entrez-testsuite", "")
	c.Check(err, check.Equals, nil)
	c.Check(sr, check.Not(check.Equals), nil)
}

func (s *S) TestDoPost(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	p, err := DoPost("protein", "biogo.entrez-testsuite", "", nil, 15718680, 157427902, 119703751)
	c.Check(err, check.Equals, nil)
	c.Assert(p.History, check.NotNil)
	c.Check(p.History.QueryKey, check.DeepEquals, 1)
	c.Check(p.History.WebEnv, check.Matches, "NCID_[0-9]+_.*")
}

func (s *S) TestFetch(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	for i, t := range []struct { //db, tool, email string, id ...int
		db      string
		rettype string
		retmode string
		tool    string
		email   string
		ids     []int
		expect  string
	}{
		{
			"protein", "fasta", "text", "biogo.entrez-testsuite", "", []int{15718680, 157427902, 119703751},
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
				">gi|119703751|ref|NP_034713.2| tyrosine-protein kinase ITK/TSK [Mus musculus]\n" +
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
		rc, err := Fetch(t.db, &Parameters{RetMode: t.rettype, RetType: t.rettype}, t.tool, t.email, nil, t.ids...)
		c.Assert(err, check.Equals, nil, check.Commentf("Test %d", i))
		b, err := ioutil.ReadAll(rc)
		rc.Close()
		c.Assert(err, check.Equals, nil)
		c.Check(string(b), check.Equals, t.expect, check.Commentf("Test %d", i))
	}
}

func (s *S) TestDoSummary(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	sum, err := DoSummary("protein", nil, "biogo.entrez-testsuite", "", nil, 15718680, 157427902, 119703751)
	c.Check(err, check.Equals, nil)
	c.Check(sum, check.DeepEquals, &Summary{
		Database: "protein",
		Docs: []Doc{
			Doc{
				Id: 15718680,
				Items: []Item{
					{Value: "NP_005537", Name: "Caption", Type: "String"},
					{Value: "tyrosine-protein kinase ITK/TSK [Homo sapiens]", Name: "Title", Type: "String"},
					{Value: "gi|15718680|ref|NP_005537.3|[15718680]", Name: "Extra", Type: "String"},
					{Value: "15718680", Name: "Gi", Type: "Integer"},
					{Value: "1999/06/09", Name: "CreateDate", Type: "String"},
					{Value: "2013/01/05", Name: "UpdateDate", Type: "String"},
					{Value: "512", Name: "Flags", Type: "Integer"},
					{Value: "9606", Name: "TaxId", Type: "Integer"},
					{Value: "620", Name: "Length", Type: "Integer"},
					{Value: "live", Name: "Status", Type: "String"},
					{Value: "", Name: "ReplacedBy", Type: "String"},
					{Value: "  ", Name: "Comment", Type: "String"},
				},
			},
			Doc{
				Id: 157427902,
				Items: []Item{
					{Value: "NP_001098858", Name: "Caption", Type: "String"},
					{Value: "tyrosine-protein kinase ITK/TSK [Bos taurus]", Name: "Title", Type: "String"},
					{Value: "gi|157427902|ref|NP_001098858.1|[157427902]", Name: "Extra", Type: "String"},
					{Value: "157427902", Name: "Gi", Type: "Integer"},
					{Value: "2007/09/24", Name: "CreateDate", Type: "String"},
					{Value: "2012/08/27", Name: "UpdateDate", Type: "String"},
					{Value: "512", Name: "Flags", Type: "Integer"},
					{Value: "9913", Name: "TaxId", Type: "Integer"},
					{Value: "620", Name: "Length", Type: "Integer"},
					{Value: "live", Name: "Status", Type: "String"},
					{Value: "", Name: "ReplacedBy", Type: "String"},
					{Value: "  ", Name: "Comment", Type: "String"},
				},
			},
			Doc{
				Id: 119703751,
				Items: []Item{
					{Value: "NP_034713", Name: "Caption", Type: "String"},
					{Value: "tyrosine-protein kinase ITK/TSK [Mus musculus]", Name: "Title", Type: "String"},
					{Value: "gi|119703751|ref|NP_034713.2|[119703751]", Name: "Extra", Type: "String"},
					{Value: "119703751", Name: "Gi", Type: "Integer"},
					{Value: "2000/01/25", Name: "CreateDate", Type: "String"},
					{Value: "2013/01/20", Name: "UpdateDate", Type: "String"},
					{Value: "512", Name: "Flags", Type: "Integer"},
					{Value: "10090", Name: "TaxId", Type: "Integer"},
					{Value: "619", Name: "Length", Type: "Integer"},
					{Value: "live", Name: "Status", Type: "String"},
					{Value: "", Name: "ReplacedBy", Type: "String"},
					{Value: "  ", Name: "Comment", Type: "String"},
				},
			},
		},
		Err: nil,
	})
}

func (s *S) TestDoLink(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	for i, t := range []struct {
		db     string
		fromDb string
		cmd    string
		query  string
		params *Parameters
		tool   string
		email  string
		hist   *History
		ids    [][]int
		expect *Link
	}{
		{
			"protein", "gene", "", "", nil, "biogo.entrez-testsuite", "", nil, [][]int{{15718680, 157427902}},
			&Link{
				LinkSets: []LinkSet{
					LinkSet{
						DbFrom: "protein",
						IdList: []int{15718680, 157427902},
						LinkSetDbs: []LinkSetDb{
							LinkSetDb{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []LinkId{
									LinkId{Id: 522311, HasLinkOut: false, HasNeighbor: false, Score: 0},
									LinkId{Id: 3702, HasLinkOut: false, HasNeighbor: false, Score: 0},
								},
							},
						},
						IdUrls:   nil,
						IdChecks: nil,
						Err:      nil,
					},
				},
				Err: nil,
			},
		},
		{
			"protein", "gene", "", "", nil, "biogo.entrez-testsuite", "", nil, [][]int{{15718680}, {157427902}},
			&Link{
				LinkSets: []LinkSet{
					{
						DbFrom: "protein",
						IdList: []int{15718680},
						LinkSetDbs: []LinkSetDb{
							{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []LinkId{
									{Id: 3702, HasLinkOut: false, HasNeighbor: false, Score: 0},
								},
							},
						},
						IdUrls:   nil,
						IdChecks: nil,
						Err:      nil,
					},
					{
						DbFrom: "protein",
						IdList: []int{157427902},
						LinkSetDbs: []LinkSetDb{
							{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []LinkId{
									{Id: 522311, HasLinkOut: false, HasNeighbor: false, Score: 0},
								},
							},
						},
						IdUrls:   nil,
						IdChecks: nil,
						Err:      nil,
					},
				},
				Err: nil,
			},
		},
	} {
		l, err := DoLink(t.db, t.fromDb, t.cmd, t.query, t.params, t.tool, t.email, t.hist, t.ids...)
		c.Check(err, check.Equals, nil, check.Commentf("Test %d", i))
		c.Check(l, check.DeepEquals, t.expect, check.Commentf("Test %d", i))
	}
}

func (s *S) TestDoGlobal(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	g, err := DoGlobal("toolkit", "biogo.entrez-testsuite", "")
	c.Check(err, check.Equals, nil)
	c.Check(g, check.Not(check.Equals), nil)
}

func (s *S) TestDoSpell(c *check.C) {
	if !*net {
		c.Skip("Network tests not requested.")
	}
	sp, err := DoSpell("", "asthmaa OR alergies", "biogo.entrez-testsuite", "")
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
		Err: nil,
	})
}
