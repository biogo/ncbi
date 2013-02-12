// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	check "launchpad.net/gocheck"
	"strings"
)

// Helpers
func intPtr(i int) *int          { return &i }
func stringPtr(s string) *string { return &s }

func (s *S) TestParseInfo(c *check.C) {
	for i, t := range []struct {
		retval string
		info   Info
	}{
		{
			`<?xml version="1.0"?>
	<!DOCTYPE eInfoResult PUBLIC "-//NLM//DTD eInfoResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eInfo_020511.dtd">
	<eInfoResult>
	<DbList>

		<DbName>pubmed</DbName>
		<DbName>protein</DbName>
		<DbName>nuccore</DbName>
		<DbName>nucleotide</DbName>
		<DbName>nucgss</DbName>
		<DbName>nucest</DbName>
		<DbName>structure</DbName>
		<DbName>genome</DbName>
		<DbName>assembly</DbName>
	</DbList>

	</eInfoResult>
	`,
			Info{
				DbList: []string{
					"pubmed",
					"protein",
					"nuccore",
					"nucleotide",
					"nucgss",
					"nucest",
					"structure",
					"genome",
					"assembly",
				},
				DbInfo: nil,
				Err:    nil,
			},
		},
		{
			`<?xml version="1.0"?>
<!DOCTYPE eInfoResult PUBLIC "-//NLM//DTD eInfoResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eInfo_020511.dtd">
<eInfoResult>
	<DbInfo>
		<DbName>toolkit</DbName>
		<MenuName>ToolKit</MenuName>
		<Description>ToolKit database</Description>
		<Count>265403</Count>
		<LastUpdate>2013/02/07 14:34</LastUpdate>
		<FieldList>
			<Field>
				<Name>ALL</Name>
				<FullName>All Fields</FullName>
				<Description>All terms from all searchable fields</Description>
				<TermCount>830579</TermCount>
				<IsDate>N</IsDate>
				<IsNumerical>N</IsNumerical>
				<SingleToken>N</SingleToken>
				<Hierarchy>N</Hierarchy>
				<IsHidden>N</IsHidden>
			</Field>
			<Field>
				<Name>UID</Name>
				<FullName>UID</FullName>
				<Description>Unique number assigned to publication</Description>
				<TermCount>0</TermCount>
				<IsDate>N</IsDate>
				<IsNumerical>Y</IsNumerical>
				<SingleToken>Y</SingleToken>
				<Hierarchy>N</Hierarchy>
				<IsHidden>Y</IsHidden>
			</Field>
		</FieldList>
		<LinkList>
		</LinkList>
	</DbInfo>

</eInfoResult>
`,
			Info{
				DbList: nil,
				DbInfo: []DbInfo{
					{
						DbName: "toolkit",

						MenuName: "ToolKit",

						Description: "ToolKit database",

						Count:      265403,
						LastUpdate: "2013/02/07 14:34",
						FieldList: []Field{
							{
								Name:        "ALL",
								FullName:    "All Fields",
								Description: "All terms from all searchable fields",
								TermCount:   830579,
								IsDate:      false,
								IsNumerical: false,
								SingleToken: false,
								Hierarchy:   false,
								IsHidden:    false,
							},
							{
								Name:        "UID",
								FullName:    "UID",
								Description: "Unique number assigned to publication",
								TermCount:   0,
								IsDate:      false,
								IsNumerical: true,
								SingleToken: true,
								Hierarchy:   false,
								IsHidden:    true,
							},
						},
						LinkList: nil,
					},
				},
				Err: nil,
			},
		},
		{
			`<?xml version="1.0"?>
	<!DOCTYPE eInfoResult PUBLIC "-//NLM//DTD eInfoResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eInfo_020511.dtd">
	<eInfoResult>
		<ERROR>Can not retrieve DbInfo for db=blah</ERROR>
	</eInfoResult>
	`,
			Info{Err: Error("Can not retrieve DbInfo for db=blah")},
		},
	} {
		var in Info
		err := in.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(in, check.DeepEquals, t.info, check.Commentf("Test: %d", i))
	}
}

func (s *S) TestParseSearch(c *check.C) {
	for i, t := range []struct {
		retval string
		search Search
	}{
		{
			`<?xml version="1.0" ?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult>
	<Count>6</Count>
	<RetMax>6</RetMax>
	<RetStart>0</RetStart>
	<IdList>
		<Id>19008416</Id>
		<Id>18927361</Id>
		<Id>18787170</Id>
		<Id>18487186</Id>
		<Id>18239126</Id>
		<Id>18239125</Id>
	</IdList>
	<TranslationSet>
		<Translation>
			<From>science[journal]</From>
			<To>"Science"[Journal] OR "Science (80- )"[Journal] OR "J Zhejiang Univ Sci"[Journal]</To>
		</Translation>
		<Translation>
			<From>breast cancer</From>
			<To>"breast neoplasms"[MeSH Terms] OR ("breast"[All Fields] AND "neoplasms"[All Fields]) OR "breast neoplasms"[All Fields] OR ("breast"[All Fields] AND "cancer"[All Fields]) OR "breast cancer"[All Fields]</To>
		</Translation>
	</TranslationSet>
	<TranslationStack>
		<TermSet>
			<Term>"Science"[Journal]</Term>
			<Field>Journal</Field>
			<Count>162433</Count>
			<Explode>Y</Explode>
		</TermSet>
		<TermSet>
			<Term>"Science (80- )"[Journal]</Term>
			<Field>Journal</Field>
			<Count>10</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>OR</OP>
		<TermSet>
		<Term>"J Zhejiang Univ Sci"[Journal]</Term>
			<Field>Journal</Field>
			<Count>364</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>OR</OP>
		<OP>GROUP</OP>
		<TermSet>
			<Term>"breast neoplasms"[MeSH Terms]</Term>
			<Field>MeSH Terms</Field>
			<Count>199283</Count>
			<Explode>Y</Explode>
		</TermSet>
		<TermSet>
			<Term>"breast"[All Fields]</Term>
			<Field>All Fields</Field>
			<Count>322674</Count>
			<Explode>Y</Explode>
		</TermSet>
		<TermSet>
			<Term>"neoplasms"[All Fields]</Term>
			<Field>All Fields</Field>
			<Count>1897643</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>AND</OP>
		<OP>GROUP</OP>
		<OP>OR</OP>
		<TermSet>
			<Term>"breast neoplasms"[All Fields]</Term>
			<Field>All Fields</Field>
			<Count>199169</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>OR</OP>
		<TermSet>
			<Term>"breast"[All Fields]</Term>
			<Field>All Fields</Field>
			<Count>322674</Count>
			<Explode>Y</Explode>
		</TermSet>
		<TermSet>
			<Term>"cancer"[All Fields]</Term>
			<Field>All Fields</Field>
			<Count>1166779</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>AND</OP>
		<OP>GROUP</OP>
		<OP>OR</OP>
		<TermSet>
			<Term>"breast cancer"[All Fields]</Term>
			<Field>All Fields</Field>
			<Count>156855</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>OR</OP>
		<OP>GROUP</OP>
		<OP>AND</OP>
		<TermSet>
			<Term>2008[pdat]</Term>
			<Field>pdat</Field>
			<Count>828593</Count>
			<Explode>Y</Explode>
		</TermSet>
		<OP>AND</OP>
	</TranslationStack>
	<QueryTranslation>("Science"[Journal] OR "Science (80- )"[Journal] OR "J Zhejiang Univ Sci"[Journal]) AND ("breast neoplasms"[MeSH Terms] OR ("breast"[All Fields] AND "neoplasms"[All Fields]) OR "breast neoplasms"[All Fields] OR ("breast"[All Fields] AND "cancer"[All Fields]) OR "breast cancer"[All Fields]) AND 2008[pdat]</QueryTranslation>
</eSearchResult>
`,
			Search{
				Count:    6,
				RetMax:   6,
				RetStart: 0,
				QueryKey: nil,
				WebEnv:   nil,
				IdList:   []int{19008416, 18927361, 18787170, 18487186, 18239126, 18239125},
				Translations: []Translation{
					{
						From: "science[journal]",
						To:   `"Science"[Journal] OR "Science (80- )"[Journal] OR "J Zhejiang Univ Sci"[Journal]`,
					},
					{
						From: "breast cancer",
						To: `` +
							`"breast neoplasms"[MeSH Terms]` +
							` OR ` +
							`(` +
							`` + `"breast"[All Fields]` +
							`` + ` AND ` +
							`` + `"neoplasms"[All Fields]` +
							`)` +
							`` + ` OR ` +
							`` + `"breast neoplasms"[All Fields]` +
							`` + ` OR ` +
							`(` +
							`` + `"breast"[All Fields] AND "cancer"[All Fields]` +
							`)` +
							` OR ` +
							`"breast cancer"[All Fields]`,
					},
				},
				TranslationStack: []Node{
					Term{
						Term:    "\"Science\"[Journal]",
						Field:   "Journal",
						Count:   162433,
						Explode: true,
					},
					Term{
						Term:    "\"Science (80- )\"[Journal]",
						Field:   "Journal",
						Count:   10,
						Explode: true,
					},
					Op("OR"),
					Term{
						Term:    "\"J Zhejiang Univ Sci\"[Journal]",
						Field:   "Journal",
						Count:   364,
						Explode: true,
					},
					Op("OR"),
					Op("GROUP"),
					Term{
						Term:    "\"breast neoplasms\"[MeSH Terms]",
						Field:   "MeSH Terms",
						Count:   199283,
						Explode: true,
					},
					Term{Term: "\"breast\"[All Fields]",
						Field:   "All Fields",
						Count:   322674,
						Explode: true,
					},
					Term{
						Term:    "\"neoplasms\"[All Fields]",
						Field:   "All Fields",
						Count:   1897643,
						Explode: true,
					},
					Op("AND"),
					Op("GROUP"),
					Op("OR"),
					Term{
						Term:    "\"breast neoplasms\"[All Fields]",
						Field:   "All Fields",
						Count:   199169,
						Explode: true,
					},
					Op("OR"),
					Term{
						Term:    "\"breast\"[All Fields]",
						Field:   "All Fields",
						Count:   322674,
						Explode: true,
					},
					Term{
						Term:    "\"cancer\"[All Fields]",
						Field:   "All Fields",
						Count:   1166779,
						Explode: true,
					},
					Op("AND"),
					Op("GROUP"),
					Op("OR"),
					Term{
						Term:    "\"breast cancer\"[All Fields]",
						Field:   "All Fields",
						Count:   156855,
						Explode: true,
					},
					Op("OR"),
					Op("GROUP"),
					Op("AND"),
					Term{
						Term:    "2008[pdat]",
						Field:   "pdat",
						Count:   828593,
						Explode: true,
					},
					Op("AND"),
				},
				QueryTranslation: stringPtr(`` +
					`(` +
					`` + `"Science"[Journal]` +
					`` + ` OR ` +
					`` + `"Science (80- )"[Journal]` +
					`` + ` OR ` +
					`` + `"J Zhejiang Univ Sci"[Journal]` +
					`)` +
					` AND ` +
					`(` +
					`` + `"breast neoplasms"[MeSH Terms] OR ` +
					`` + `(` +
					`` + `` + `"breast"[All Fields]` +
					`` + `` + ` AND ` +
					`` + `` + `"neoplasms"[All Fields]` +
					`` + `)` +
					`` + ` OR ` +
					`` + `"breast neoplasms"[All Fields]` +
					`` + ` OR ` +
					`` + `(` +
					`` + `` + `"breast"[All Fields]` +
					`` + `` + ` AND ` +
					`` + `` + `"cancer"[All Fields]` +
					`` + `)` +
					`` + ` OR ` +
					`` + `"breast cancer"[All Fields]` +
					`)` +
					` AND ` +
					`2008[pdat]`),
				Err:      nil,
				Errors:   nil,
				Warnings: nil,
			},
		},
		{`<?xml version="1.0" ?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult>
	<ERROR>Empty term and query_key - nothing todo</ERROR>
</eSearchResult>
`,
			Search{Err: Error("Empty term and query_key - nothing todo")},
		},
		{`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult>
	<ERROR>Invalid db name specified: pub</ERROR>
</eSearchResult>
`,
			Search{Err: Error("Invalid db name specified: pub")},
		},
		{
			`<?xml version="1.0" ?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult><Count>20</Count><RetMax>20</RetMax><RetStart>0</RetStart>
	<ErrorList>
		<FieldNotFound>jungle</FieldNotFound>
		<FieldNotFound>pat</FieldNotFound>
	</ErrorList>
</eSearchResult>
`,
			Search{
				Count:            20,
				RetMax:           20,
				RetStart:         0,
				QueryKey:         nil,
				WebEnv:           nil,
				IdList:           nil,
				Translations:     nil,
				TranslationStack: nil,
				QueryTranslation: nil,
				Err:              nil,
				Errors: []error{
					NotFoundError{Type: "field not found", Value: "jungle"},
					NotFoundError{Type: "field not found", Value: "pat"},
				},
				Warnings: nil,
			},
		},
		{
			`<?xml version="1.0" ?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult>
	<Count>0</Count>
	<RetMax>0</RetMax>
	<RetStart>0</RetStart>
	<IdList>
	</IdList>
	<TranslationSet/>
	<QueryTranslation>(nonjournal[journal] AND nonyear[date])</QueryTranslation>
	<ErrorList>
		<PhraseNotFound>nonjournal[journal]</PhraseNotFound>
		<PhraseNotFound>nonyear[date]</PhraseNotFound>
	</ErrorList>
	<WarningList>
		<OutputMessage>No items found.</OutputMessage>
	</WarningList>
</eSearchResult>
`,
			Search{
				Count:            0,
				RetMax:           0,
				RetStart:         0,
				QueryKey:         nil,
				WebEnv:           nil,
				IdList:           nil,
				Translations:     nil,
				TranslationStack: nil,
				QueryTranslation: stringPtr("(nonjournal[journal] AND nonyear[date])"),
				Err:              nil,
				Errors: []error{
					NotFoundError{Type: "phrase not found", Value: "nonjournal[journal]"},
					NotFoundError{Type: "phrase not found", Value: "nonyear[date]"},
				},
				Warnings: []error{
					Warning{Type: "output message", Value: "No items found."},
				},
			},
		},
	} {
		var s Search
		err := s.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(s, check.DeepEquals, t.search, check.Commentf("Test: %d", i))
	}
}

func (s *S) TestParsePost(c *check.C) {
	for i, t := range []struct {
		retval string
		post   Post
	}{
		{
			`<?xml version="1.0"?>
<!DOCTYPE ePostResult PUBLIC "-//NLM//DTD ePostResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/ePost_020511.dtd">
<ePostResult>
	<QueryKey>1</QueryKey>
	<WebEnv>NCID_1_298287560_130.14.18.48_5555_1360293704_91337037</WebEnv>
</ePostResult>
`,
			Post{
				QueryKey: intPtr(1),
				WebEnv:   stringPtr("NCID_1_298287560_130.14.18.48_5555_1360293704_91337037"),
			},
		},
		{
			`<?xml version="1.0"?>
<!DOCTYPE ePostResult PUBLIC "-//NLM//DTD ePostResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/ePost_020511.dtd">
<ePostResult>
	<QueryKey>1</QueryKey>
	<WebEnv>NCID_1_299062774_130.14.18.97_5555_1360293760_1713152879</WebEnv>
	<ERROR>IDs contain invalid characters which was treated as delimiters.</ERROR>
</ePostResult>
`,
			Post{
				QueryKey: intPtr(1),
				WebEnv:   stringPtr("NCID_1_299062774_130.14.18.97_5555_1360293760_1713152879"),
				Err:      Error("IDs contain invalid characters which was treated as delimiters."),
			},
		},
		{
			`<?xml version="1.0"?>
<!DOCTYPE ePostResult PUBLIC "-//NLM//DTD ePostResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/ePost_020511.dtd">
<ePostResult>
	<InvalidIdList>
		<Id>19008416</Id>
		<Id>18927361</Id>
		<Id>18787170</Id>
		<Id>18487186</Id>
		<Id>18239126</Id>
		<Id>18239125</Id>
	</InvalidIdList>
</ePostResult>
`,
			Post{
				InvalidIds: []int{19008416, 18927361, 18787170, 18487186, 18239126, 18239125},
			},
		},
	} {
		var p Post
		err := p.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(p, check.DeepEquals, t.post, check.Commentf("Test: %d", i))
	}
}

func (s *S) TestParseSummary(c *check.C) {
	for i, t := range []struct {
		retval  string
		summary Summary
	}{
		{
			`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE eSummaryResult PUBLIC "-//NLM//DTD eSummaryResult, 29 October 2004//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSummary_041029.dtd">
<eSummaryResult>
<DocSum>
	<Id>6678417</Id>
	<Item Name="Caption" Type="String">NP_033443</Item>
	<Item Name="Title" Type="String">thyroid peroxidase precursor [Mus musculus]</Item>
	<Item Name="Extra" Type="String">gi|6678417|ref|NP_033443.1|[6678417]</Item>
	<Item Name="Gi" Type="Integer">6678417</Item>
	<Item Name="CreateDate" Type="String">2000/01/04</Item>
	<Item Name="UpdateDate" Type="String">2012/12/12</Item>
	<Item Name="Flags" Type="Integer">512</Item>
	<Item Name="TaxId" Type="Integer">10090</Item>
	<Item Name="Length" Type="Integer">914</Item>
	<Item Name="Status" Type="String">live</Item>
	<Item Name="ReplacedBy" Type="String"></Item>
	<Item Name="Comment" Type="String"><![CDATA[  ]]></Item>
</DocSum>

<DocSum>
	<Id>9507199</Id>
	<Item Name="Caption" Type="String">NP_062226</Item>
	<Item Name="Title" Type="String">thyroid peroxidase precursor [Rattus norvegicus]</Item>
	<Item Name="Extra" Type="String">gi|9507199|ref|NP_062226.1|[9507199]</Item>
	<Item Name="Gi" Type="Integer">9507199</Item>
	<Item Name="CreateDate" Type="String">2000/07/22</Item>
	<Item Name="UpdateDate" Type="String">2012/03/24</Item>
	<Item Name="Flags" Type="Integer">512</Item>
	<Item Name="TaxId" Type="Integer">10116</Item>
	<Item Name="Length" Type="Integer">914</Item>
	<Item Name="Status" Type="String">replaced</Item>
	<Item Name="ReplacedBy" Type="String">NP_062226.2</Item>
	<Item Name="Comment" Type="String"><![CDATA[ This record was replaced or removed. ]]></Item>
</DocSum>

<DocSum>
	<Id>28558982</Id>
	<Item Name="Caption" Type="String">NP_000538</Item>
	<Item Name="Title" Type="String">thyroid peroxidase isoform a precursor [Homo sapiens]</Item>
	<Item Name="Extra" Type="String">gi|28558982|ref|NP_000538.3|[28558982]</Item>
	<Item Name="Gi" Type="Integer">28558982</Item>
	<Item Name="CreateDate" Type="String">1999/03/19</Item>
	<Item Name="UpdateDate" Type="String">2013/01/07</Item>
	<Item Name="Flags" Type="Integer">512</Item>
	<Item Name="TaxId" Type="Integer">9606</Item>
	<Item Name="Length" Type="Integer">933</Item>
	<Item Name="Status" Type="String">live</Item>
	<Item Name="ReplacedBy" Type="String"></Item>
	<Item Name="Comment" Type="String"><![CDATA[  ]]></Item>
</DocSum>

<DocSum>
	<Id>28558984</Id>
	<Item Name="Caption" Type="String">NP_783650</Item>
	<Item Name="Title" Type="String">thyroid peroxidase isoform b precursor [Homo sapiens]</Item>
	<Item Name="Extra" Type="String">gi|28558984|ref|NP_783650.1|[28558984]</Item>
	<Item Name="Gi" Type="Integer">28558984</Item>
	<Item Name="CreateDate" Type="String">2003/02/25</Item>
	<Item Name="UpdateDate" Type="String">2013/01/07</Item>
	<Item Name="Flags" Type="Integer">512</Item>
	<Item Name="TaxId" Type="Integer">9606</Item>
	<Item Name="Length" Type="Integer">876</Item>
	<Item Name="Status" Type="String">live</Item>
	<Item Name="ReplacedBy" Type="String"></Item>
	<Item Name="Comment" Type="String"><![CDATA[  ]]></Item>
</DocSum>

<DocSum>
	<Id>28558988</Id>
	<Item Name="Caption" Type="String">NP_783652</Item>
	<Item Name="Title" Type="String">thyroid peroxidase isoform d precursor [Homo sapiens]</Item>
	<Item Name="Extra" Type="String">gi|28558988|ref|NP_783652.1|[28558988]</Item>
	<Item Name="Gi" Type="Integer">28558988</Item>
	<Item Name="CreateDate" Type="String">2003/02/25</Item>
	<Item Name="UpdateDate" Type="String">2013/01/07</Item>
	<Item Name="Flags" Type="Integer">512</Item>
	<Item Name="TaxId" Type="Integer">9606</Item>
	<Item Name="Length" Type="Integer">889</Item>
	<Item Name="Status" Type="String">live</Item>
	<Item Name="ReplacedBy" Type="String"></Item>
	<Item Name="Comment" Type="String"><![CDATA[  ]]></Item>
</DocSum>

<DocSum>
	<Id>28558990</Id>
	<Item Name="Caption" Type="String">NP_783653</Item>
	<Item Name="Title" Type="String">thyroid peroxidase isoform e precursor [Homo sapiens]</Item>
	<Item Name="Extra" Type="String">gi|28558990|ref|NP_783653.1|[28558990]</Item>
	<Item Name="Gi" Type="Integer">28558990</Item>
	<Item Name="CreateDate" Type="String">2003/02/25</Item>
	<Item Name="UpdateDate" Type="String">2013/01/07</Item>
	<Item Name="Flags" Type="Integer">512</Item>
	<Item Name="TaxId" Type="Integer">9606</Item>
	<Item Name="Length" Type="Integer">760</Item>
	<Item Name="Status" Type="String">live</Item>
	<Item Name="ReplacedBy" Type="String"></Item>
	<Item Name="Comment" Type="String"><![CDATA[  ]]></Item>
</DocSum>

</eSummaryResult>
`,
			Summary{
				Docs: []Doc{
					{
						Id: 6678417,
						Items: []Item{
							{Name: "Caption", Type: "String", Value: "NP_033443"},
							{Name: "Title", Type: "String", Value: "thyroid peroxidase precursor [Mus musculus]"},
							{Name: "Extra", Type: "String", Value: "gi|6678417|ref|NP_033443.1|[6678417]"},
							{Name: "Gi", Type: "Integer", Value: "6678417"},
							{Name: "CreateDate", Type: "String", Value: "2000/01/04"},
							{Name: "UpdateDate", Type: "String", Value: "2012/12/12"},
							{Name: "Flags", Type: "Integer", Value: "512"},
							{Name: "TaxId", Type: "Integer", Value: "10090"},
							{Name: "Length", Type: "Integer", Value: "914"},
							{Name: "Status", Type: "String", Value: "live"},
							{Name: "ReplacedBy", Type: "String", Value: ""},
							{Name: "Comment", Type: "String", Value: "  "},
						},
					},
					{
						Id: 9507199,
						Items: []Item{
							{Name: "Caption", Type: "String", Value: "NP_062226"},
							{Name: "Title", Type: "String", Value: "thyroid peroxidase precursor [Rattus norvegicus]"},
							{Name: "Extra", Type: "String", Value: "gi|9507199|ref|NP_062226.1|[9507199]"},
							{Name: "Gi", Type: "Integer", Value: "9507199"},
							{Name: "CreateDate", Type: "String", Value: "2000/07/22"},
							{Name: "UpdateDate", Type: "String", Value: "2012/03/24"},
							{Name: "Flags", Type: "Integer", Value: "512"},
							{Name: "TaxId", Type: "Integer", Value: "10116"},
							{Name: "Length", Type: "Integer", Value: "914"},
							{Name: "Status", Type: "String", Value: "replaced"},
							{Name: "ReplacedBy", Type: "String", Value: "NP_062226.2"},
							{Name: "Comment", Type: "String", Value: " This record was replaced or removed. "},
						},
					},
					{
						Id: 28558982,
						Items: []Item{
							{Name: "Caption", Type: "String", Value: "NP_000538"},
							{Name: "Title", Type: "String", Value: "thyroid peroxidase isoform a precursor [Homo sapiens]"},
							{Name: "Extra", Type: "String", Value: "gi|28558982|ref|NP_000538.3|[28558982]"},
							{Name: "Gi", Type: "Integer", Value: "28558982"},
							{Name: "CreateDate", Type: "String", Value: "1999/03/19"},
							{Name: "UpdateDate", Type: "String", Value: "2013/01/07"},
							{Name: "Flags", Type: "Integer", Value: "512"},
							{Name: "TaxId", Type: "Integer", Value: "9606"},
							{Name: "Length", Type: "Integer", Value: "933"},
							{Name: "Status", Type: "String", Value: "live"},
							{Name: "ReplacedBy", Type: "String", Value: ""},
							{Name: "Comment", Type: "String", Value: "  "},
						},
					},
					{
						Id: 28558984,
						Items: []Item{
							{Name: "Caption", Type: "String", Value: "NP_783650"},
							{Name: "Title", Type: "String", Value: "thyroid peroxidase isoform b precursor [Homo sapiens]"},
							{Name: "Extra", Type: "String", Value: "gi|28558984|ref|NP_783650.1|[28558984]"},
							{Name: "Gi", Type: "Integer", Value: "28558984"},
							{Name: "CreateDate", Type: "String", Value: "2003/02/25"},
							{Name: "UpdateDate", Type: "String", Value: "2013/01/07"},
							{Name: "Flags", Type: "Integer", Value: "512"},
							{Name: "TaxId", Type: "Integer", Value: "9606"},
							{Name: "Length", Type: "Integer", Value: "876"},
							{Name: "Status", Type: "String", Value: "live"},
							{Name: "ReplacedBy", Type: "String", Value: ""},
							{Name: "Comment", Type: "String", Value: "  "},
						},
					},
					{
						Id: 28558988,
						Items: []Item{
							{Name: "Caption", Type: "String", Value: "NP_783652"},
							{Name: "Title", Type: "String", Value: "thyroid peroxidase isoform d precursor [Homo sapiens]"},
							{Name: "Extra", Type: "String", Value: "gi|28558988|ref|NP_783652.1|[28558988]"},
							{Name: "Gi", Type: "Integer", Value: "28558988"},
							{Name: "CreateDate", Type: "String", Value: "2003/02/25"},
							{Name: "UpdateDate", Type: "String", Value: "2013/01/07"},
							{Name: "Flags", Type: "Integer", Value: "512"},
							{Name: "TaxId", Type: "Integer", Value: "9606"},
							{Name: "Length", Type: "Integer", Value: "889"},
							{Name: "Status", Type: "String", Value: "live"},
							{Name: "ReplacedBy", Type: "String", Value: ""},
							{Name: "Comment", Type: "String", Value: "  "},
						},
					},
					{
						Id: 28558990,
						Items: []Item{
							{Name: "Caption", Type: "String", Value: "NP_783653"},
							{Name: "Title", Type: "String", Value: "thyroid peroxidase isoform e precursor [Homo sapiens]"},
							{Name: "Extra", Type: "String", Value: "gi|28558990|ref|NP_783653.1|[28558990]"},
							{Name: "Gi", Type: "Integer", Value: "28558990"},
							{Name: "CreateDate", Type: "String", Value: "2003/02/25"},
							{Name: "UpdateDate", Type: "String", Value: "2013/01/07"},
							{Name: "Flags", Type: "Integer", Value: "512"},
							{Name: "TaxId", Type: "Integer", Value: "9606"},
							{Name: "Length", Type: "Integer", Value: "760"},
							{Name: "Status", Type: "String", Value: "live"},
							{Name: "ReplacedBy", Type: "String", Value: ""},
							{Name: "Comment", Type: "String", Value: "  "},
						},
					},
				},
				Err: nil,
			},
		},
	} {
		var s Summary
		err := s.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(s, check.DeepEquals, t.summary, check.Commentf("Test: %d", i))
	}
}

func (s *S) TestParseLink(c *check.C) {
	for i, t := range []struct {
		retval string
		link   Link
	}{
		{
			`<?xml version="1.0"?>
<!DOCTYPE eLinkResult PUBLIC "-//NLM//DTD eLinkResult, 23 November 2010//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eLink_101123.dtd">
<eLinkResult>

	<LinkSet>
		<DbFrom>protein</DbFrom>
		<IdList>
			<Id>15718680</Id>
			<Id>157427902</Id>
		</IdList>
		<LinkSetDb>
			<DbTo>gene</DbTo>
			<LinkName>protein_gene</LinkName>
			<Link>
				<Id>522311</Id>
			</Link>
			<Link>
				<Id>3702</Id>
			</Link>
		</LinkSetDb>
	</LinkSet>

</eLinkResult>
`,
			Link{
				LinkSets: []LinkSet{
					{
						DbFrom: "protein",
						IdList: []int{
							15718680,
							157427902,
						},
						LinkSetDbs: []LinkSetDb{
							{
								DbTo:     "gene",
								LinkName: "protein_gene",
								Link: []LinkId{
									{Id: 522311, HasLinkOut: false, HasNeighbor: false, Score: 0},
									{Id: 3702, HasLinkOut: false, HasNeighbor: false, Score: 0},
								},
							},
						},
						IdUrls: nil,
						Err:    nil,
					},
				},
				Err: nil,
			},
		},
		{
			`<?xml version="1.0"?>
<!DOCTYPE eLinkResult PUBLIC "-//NLM//DTD eLinkResult, 23 November 2010//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eLink_101123.dtd">
<eLinkResult>

	<LinkSet>
		<DbFrom>pubmed</DbFrom>
		<IdList>
			<Id>20210808</Id>
			<Id>20210909</Id>
		</IdList>
		<LinkSetDb>
			<DbTo>pubmed</DbTo>
			<LinkName>pubmed_pubmed</LinkName>
			<Link>
				<Id>15876306</Id>
				<Score>75133399</Score>
			</Link>
			<Link>
				<Id>20816181</Id>
				<Score>25095241</Score>
			</Link>
			<Link>
				<Id>21053465</Id>
				<Score>24834712</Score>
			</Link>
			<Link>
				<Id>22032786</Id>
				<Score>24243731</Score>
			</Link>
			<Link>
				<Id>22374193</Id>
				<Score>23718577</Score>
			</Link>
			<Link>
				<Id>19387030</Id>
				<Score>23425951</Score>
			</Link>
			<Link>
				<Id>21978852</Id>
				<Score>22647663</Score>
			</Link>
			<Link>
				<Id>22857403</Id>
				<Score>19564745</Score>
			</Link>
		</LinkSetDb>
		<LinkSetDb>
			<DbTo>pubmed</DbTo>
			<LinkName>pubmed_pubmed_reviews_five</LinkName>
			<Link>
				<Id>12376064</Id>
				<Score>56460889</Score>
			</Link>
			<Link>
				<Id>15125698</Id>
				<Score>50774274</Score>
			</Link>
			<Link>
				<Id>10931782</Id>
				<Score>50227044</Score>
			</Link>
			<Link>
				<Id>10096822</Id>
				<Score>48788287</Score>
			</Link>
			<Link>
				<Id>12582308</Id>
				<Score>48635669</Score>
			</Link>
		</LinkSetDb>
	</LinkSet>
</eLinkResult>
`,
			Link{
				LinkSets: []LinkSet{
					{
						DbFrom: "pubmed",
						IdList: []int{
							20210808,
							20210909,
						},
						LinkSetDbs: []LinkSetDb{
							{
								DbTo:     "pubmed",
								LinkName: "pubmed_pubmed",
								Link: []LinkId{
									{Id: 15876306, HasLinkOut: false, HasNeighbor: false, Score: 75133399},
									{Id: 20816181, HasLinkOut: false, HasNeighbor: false, Score: 25095241},
									{Id: 21053465, HasLinkOut: false, HasNeighbor: false, Score: 24834712},
									{Id: 22032786, HasLinkOut: false, HasNeighbor: false, Score: 24243731},
									{Id: 22374193, HasLinkOut: false, HasNeighbor: false, Score: 23718577},
									{Id: 19387030, HasLinkOut: false, HasNeighbor: false, Score: 23425951},
									{Id: 21978852, HasLinkOut: false, HasNeighbor: false, Score: 22647663},
									{Id: 22857403, HasLinkOut: false, HasNeighbor: false, Score: 19564745},
								},
							},
							{
								DbTo:     "pubmed",
								LinkName: "pubmed_pubmed_reviews_five",
								Link: []LinkId{
									{Id: 12376064, HasLinkOut: false, HasNeighbor: false, Score: 56460889},
									{Id: 15125698, HasLinkOut: false, HasNeighbor: false, Score: 50774274},
									{Id: 10931782, HasLinkOut: false, HasNeighbor: false, Score: 50227044},
									{Id: 10096822, HasLinkOut: false, HasNeighbor: false, Score: 48788287},
									{Id: 12582308, HasLinkOut: false, HasNeighbor: false, Score: 48635669},
								},
							},
						},
						IdUrls: nil,
						Err:    nil,
					},
				},
				Err: nil},
		},
		{
			`<?xml version="1.0"?>
<!DOCTYPE eLinkResult PUBLIC "-//NLM//DTD eLinkResult, 23 November 2010//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eLink_101123.dtd">
<eLinkResult>

<LinkSet>
	<DbFrom>pubmed</DbFrom>
	<IdUrlList>
		<IdUrlSet>
			<Id>19880848</Id>
			<ObjUrl>

				<Url>http://www.labome.org//expert/switzerland/university/klingenberg/roland-klingenberg-1568163.html</Url>
				<LinkName>Labome Researcher Resource</LinkName>
				<SubjectType>author profiles</SubjectType>
				<Category>Other Literature Sources</Category>
				<Attribute>subscription/membership/fee required</Attribute>
				<Provider>
					<Name>ExactAntigen/Labome</Name>
					<NameAbbr>EADB</NameAbbr>
					<Id>5753</Id>
					<Url LNG="EN">http://www.labome.com</Url>
				</Provider>
			</ObjUrl>
			<ObjUrl>

				<Url>http://eurheartj.oxfordjournals.org/cgi/pmidlookup?view=long&amp;pmid=19880848</Url>
				<IconUrl LNG="EN">http://www.ncbi.nlm.nih.gov/corehtml/query/egifs/http:--highwire.stanford.edu-icons-externalservices-pubmed-custom-oxfordjournals_final_free.gif</IconUrl>
				<SubjectType>publishers/providers</SubjectType>
				<Category>Full Text Sources</Category>
				<Attribute>free resource</Attribute>
				<Attribute>full-text online</Attribute>
				<Attribute>publisher of information in url</Attribute>
				<Provider>
					<Name>HighWire</Name>
					<NameAbbr>HighWire</NameAbbr>
					<Id>3051</Id>
					<Url LNG="EN">http://highwire.stanford.edu</Url>
				</Provider>
			</ObjUrl>
			<ObjUrl>

				<Url>http://www.nlm.nih.gov/medlineplus/atherosclerosis.html</Url>
				<IconUrl LNG="EN">http://www.ncbi.nlm.nih.gov/corehtml/query/egifs/http:--www.nlm.nih.gov-medlineplus-images-linkout_sm.gif</IconUrl>
				<LinkName>Atherosclerosis</LinkName>
				<SubjectType>consumer health</SubjectType>
				<Category>Medical</Category>
				<Attribute>free resource</Attribute>
				<Provider>
					<Name>MedlinePlus Health Information</Name>
					<NameAbbr>MEDPLUS</NameAbbr>
					<Id>3162</Id>
					<Url LNG="EN">http://medlineplus.gov/</Url>
					<IconUrl LNG="EN">http://www.nlm.nih.gov/medlineplus/images/linkout_sm.gif</IconUrl>
				</Provider>
			</ObjUrl>
		</IdUrlSet>
	</IdUrlList>
</LinkSet>

</eLinkResult>
`,
			Link{
				LinkSets: []LinkSet{
					{
						DbFrom: "pubmed",
						IdUrls: []IdUrlList{
							{
								IdUrlSets: []IdUrlSet{
									{
										Id: 19880848,
										ObjUrls: []ObjUrl{
											{
												Url:         "http://www.labome.org//expert/switzerland/university/klingenberg/roland-klingenberg-1568163.html",
												LinkName:    stringPtr("Labome Researcher Resource"),
												SubjectType: []string{"author profiles"},
												Category:    []string{"Other Literature Sources"},
												Attribute:   []string{"subscription/membership/fee required"},
												Provider: Provider{
													Name:     "ExactAntigen/Labome",
													NameAbbr: "EADB",
													Id:       5753,
													Url:      "http://www.labome.com",
												},
											},
											{
												Url:         "http://eurheartj.oxfordjournals.org/cgi/pmidlookup?view=long&pmid=19880848",
												IconUrl:     stringPtr("http://www.ncbi.nlm.nih.gov/corehtml/query/egifs/http:--highwire.stanford.edu-icons-externalservices-pubmed-custom-oxfordjournals_final_free.gif"),
												SubjectType: []string{"publishers/providers"},
												Category:    []string{"Full Text Sources"},
												Attribute:   []string{"free resource", "full-text online", "publisher of information in url"},
												Provider: Provider{
													Name:     "HighWire",
													NameAbbr: "HighWire",
													Id:       3051,
													Url:      "http://highwire.stanford.edu",
												},
											},
											{
												Url:         "http://www.nlm.nih.gov/medlineplus/atherosclerosis.html",
												IconUrl:     stringPtr("http://www.ncbi.nlm.nih.gov/corehtml/query/egifs/http:--www.nlm.nih.gov-medlineplus-images-linkout_sm.gif"),
												LinkName:    stringPtr("Atherosclerosis"),
												SubjectType: []string{"consumer health"},
												Category:    []string{"Medical"},
												Attribute:   []string{"free resource"},
												Provider: Provider{
													Name:     "MedlinePlus Health Information",
													NameAbbr: "MEDPLUS",
													Id:       3162,
													Url:      "http://medlineplus.gov/",
													IconUrl:  stringPtr("http://www.nlm.nih.gov/medlineplus/images/linkout_sm.gif"),
												},
											},
										},
									},
								},
								Err: nil,
							},
						},
						Err: nil,
					},
				},
				Err: nil,
			},
		},
		{
			`<?xml version="1.0"?>
		<!DOCTYPE eLinkResult PUBLIC "-//NLM//DTD eLinkResult, 23 November 2010//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eLink_101123.dtd">
		<eLinkResult>
			<LinkSet>
				<DbFrom>nuccore</DbFrom>
				<IdCheckList>
					<Id HasNeighbor="Y">21614549</Id>
					<Id HasNeighbor="N">219152114</Id>
				</IdCheckList>
			</LinkSet>
		</eLinkResult>
		`,
			Link{
				LinkSets: []LinkSet{
					{
						DbFrom:     "nuccore",
						IdList:     nil,
						LinkSetDbs: nil,
						IdUrls:     nil,
						IdChecks: &IdChecks{
							Ids: []Id{
								{Id: 21614549, HasLinkOut: false, HasNeighbor: true},
								{Id: 219152114, HasLinkOut: false, HasNeighbor: false},
							},
							Err: nil,
						},
						Err: nil,
					},
				},
				Err: nil,
			},
		},
		// 		{
		// 			`<?xml version="1.0"?>
		// <!DOCTYPE eLinkResult PUBLIC "-//NLM//DTD eLinkResult, 23 November 2010//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eLink_101123.dtd">
		// <eLinkResult>
		// 	<LinkSet>
		// 		<DbFrom>protein</DbFrom>
		// 		<IdCheckList>
		// 			<IdLinkSet>
		// 				<Id>15718680</Id>
		// 				<LinkInfo>
		// 					<DbTo>pubmed</DbTo>
		// 					<LinkName>protein_pubmed</LinkName>
		// 					<MenuTag>PubMed Links</MenuTag>
		// 					<HtmlTag>PubMed</HtmlTag>
		// 					<Priority>128</Priority>
		// 				</LinkInfo>
		// 				<LinkInfo>
		// 					<DbTo>pubmed</DbTo>
		// 					<LinkName>protein_pubmed_refseq</LinkName>
		// 					<MenuTag>PubMed (RefSeq) Links</MenuTag>
		// 					<HtmlTag>PubMed (RefSeq)</HtmlTag>
		// 					<Priority>128</Priority>
		// 				</LinkInfo>
		// 				<LinkInfo>
		// 					<DbTo>pubmed</DbTo>
		// 					<LinkName>protein_pubmed_weighted</LinkName>
		// 					<MenuTag>PubMed (Weighted) Links</MenuTag>
		// 					<HtmlTag>PubMed (Weighted)</HtmlTag>
		// 					<Priority>128</Priority>
		// 				</LinkInfo>
		// 				<LinkInfo>
		// 					<DbTo>LinkOut</DbTo>
		// 					<LinkName>ExternalLink</LinkName>
		// 					<MenuTag>LinkOut</MenuTag>
		// 					<HtmlTag>LinkOut</HtmlTag>
		// 					<Priority>255</Priority>
		// 				</LinkInfo>
		// 			</IdLinkSet>
		// 			<IdLinkSet>
		// 				<Id>157427902</Id>
		// 				<LinkInfo>
		// 					<DbTo>pubmed</DbTo>
		// 					<LinkName>protein_pubmed</LinkName>
		// 					<MenuTag>PubMed Links</MenuTag>
		// 					<HtmlTag>PubMed</HtmlTag>
		// 					<Priority>128</Priority>
		// 				</LinkInfo>
		// 				<LinkInfo>
		// 					<DbTo>pubmed</DbTo>
		// 					<LinkName>protein_pubmed_refseq</LinkName>
		// 					<MenuTag>PubMed (RefSeq) Links</MenuTag>
		// 					<HtmlTag>PubMed (RefSeq)</HtmlTag>
		// 					<Priority>128</Priority>
		// 				</LinkInfo>
		// 				<LinkInfo>
		// 					<DbTo>pubmed</DbTo>
		// 					<LinkName>protein_pubmed_weighted</LinkName>
		// 					<MenuTag>PubMed (Weighted) Links</MenuTag>
		// 					<HtmlTag>PubMed (Weighted)</HtmlTag>
		// 					<Priority>128</Priority>
		// 				</LinkInfo>
		// 				<LinkInfo>
		// 					<DbTo>LinkOut</DbTo>
		// 					<LinkName>ExternalLink</LinkName>
		// 					<MenuTag>LinkOut</MenuTag>
		// 					<HtmlTag>LinkOut</HtmlTag>
		// 					<Priority>255</Priority>
		// 				</LinkInfo>
		// 			</IdLinkSet>
		// 		</IdCheckList>
		// 	</LinkSet>
		// </eLinkResult>
		// `,
		// 			Link{},
		// 		},
	} {
		var l Link
		err := l.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(l, check.DeepEquals, t.link, check.Commentf("Test: %d", i))
	}
}

func (s *S) TestParseGlobal(c *check.C) {
	for i, t := range []struct {
		retval string
		global Global
	}{
		{
			`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE Result PUBLIC "-//NLM//DTD eSearchResult, January 2004//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/egquery.dtd">
<!-- Internal subset (illegally) redefines Term and eGQueryResult. Turns a couple dozen validation errors into two validation errors.  -->
<!-- Disabled internal subset because it doesn't really fix anything.
<!ELEMENT component (#PCDATA)>
<!ATTLIST component
    id     CDATA    #REQUIRED
    label  CDATA    #IMPLIED>

<!ELEMENT Term (#PCDATA|component)*>

<!ELEMENT eGQueryResult (ResultItem|component)+>

-->
<Result>

    <Term>health</Term>

    <eGQueryResult>
        <ResultItem><DbName>pubmed</DbName><MenuName>PubMed</MenuName><Count>2398129</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>pmc</DbName><MenuName>PubMed Central</MenuName><Count>1217453</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>mesh</DbName><MenuName>MeSH</MenuName><Count>239</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>books</DbName><MenuName>Books</MenuName><Count>76226</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>omim</DbName><MenuName>OMIM</MenuName><Count>356</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>omia</DbName><MenuName>OMIA</MenuName><Count>3</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>ncbisearch</DbName><MenuName>Site Search</MenuName><Count>24198</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>nuccore</DbName><MenuName>Nucleotide</MenuName><Count>2841517</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>nucgss</DbName><MenuName>GSS</MenuName><Count>913</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>nucest</DbName><MenuName>EST</MenuName><Count>716932</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>protein</DbName><MenuName>Protein</MenuName><Count>5641794</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>genome</DbName><MenuName>Genome</MenuName><Count>37</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>structure</DbName><MenuName>Structure</MenuName><Count>760</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>taxonomy</DbName><MenuName>Taxonomy</MenuName><Count>0</Count><Status>Term or Database is not found</Status></ResultItem>
        <ResultItem><DbName>snp</DbName><MenuName>SNP</MenuName><Count>19</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>dbvar</DbName><MenuName>dbVar</MenuName><Count>1749</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>epigenomics</DbName><MenuName>Epigenomics</MenuName><Count>399</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>gene</DbName><MenuName>Gene</MenuName><Count>15243</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>sra</DbName><MenuName>SRA</MenuName><Count>5314</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>biosystems</DbName><MenuName>BioSystems</MenuName><Count>130</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>unigene</DbName><MenuName>UniGene</MenuName><Count>0</Count><Status>Term or Database is not found</Status></ResultItem>
        <ResultItem><DbName>cdd</DbName><MenuName>CDD</MenuName><Count>9</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>clone</DbName><MenuName>Clone</MenuName><Count>0</Count><Status>Term or Database is not found</Status></ResultItem>
        <ResultItem><DbName>unists</DbName><MenuName>UniSTS</MenuName><Count>2</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>popset</DbName><MenuName>PopSet</MenuName><Count>4221</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>geoprofiles</DbName><MenuName>GEO Profiles</MenuName><Count>287541</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>gds</DbName><MenuName>GEO DataSets</MenuName><Count>12262</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>homologene</DbName><MenuName>HomoloGene</MenuName><Count>48</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>pccompound</DbName><MenuName>PubChem Compound</MenuName><Count>53</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>pcsubstance</DbName><MenuName>PubChem Substance</MenuName><Count>15170</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>pcassay</DbName><MenuName>PubChem BioAssay</MenuName><Count>3614</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>nlmcatalog</DbName><MenuName>NLM Catalog</MenuName><Count>254830</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>probe</DbName><MenuName>Probe</MenuName><Count>43</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>gap</DbName><MenuName>dbGaP</MenuName><Count>105672</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>proteinclusters</DbName><MenuName>Protein Clusters</MenuName><Count>0</Count><Status>Term or Database is not found</Status></ResultItem>

        <ResultItem><DbName>bioproject</DbName><MenuName>BioProject</MenuName><Count>2339</Count><Status>Ok</Status></ResultItem>
        <ResultItem><DbName>biosample</DbName><MenuName>BioSample</MenuName><Count>93830</Count><Status>Ok</Status></ResultItem>
    </eGQueryResult>

</Result>
`,
			Global{
				Query: "health",
				Results: []Result{
					{Database: "pubmed", MenuName: "PubMed", Count: 2398129, Status: "Ok"},
					{Database: "pmc", MenuName: "PubMed Central", Count: 1217453, Status: "Ok"},
					{Database: "mesh", MenuName: "MeSH", Count: 239, Status: "Ok"},
					{Database: "books", MenuName: "Books", Count: 76226, Status: "Ok"},
					{Database: "omim", MenuName: "OMIM", Count: 356, Status: "Ok"},
					{Database: "omia", MenuName: "OMIA", Count: 3, Status: "Ok"},
					{Database: "ncbisearch", MenuName: "Site Search", Count: 24198, Status: "Ok"},
					{Database: "nuccore", MenuName: "Nucleotide", Count: 2841517, Status: "Ok"},
					{Database: "nucgss", MenuName: "GSS", Count: 913, Status: "Ok"},
					{Database: "nucest", MenuName: "EST", Count: 716932, Status: "Ok"},
					{Database: "protein", MenuName: "Protein", Count: 5641794, Status: "Ok"},
					{Database: "genome", MenuName: "Genome", Count: 37, Status: "Ok"},
					{Database: "structure", MenuName: "Structure", Count: 760, Status: "Ok"},
					{Database: "taxonomy", MenuName: "Taxonomy", Count: 0, Status: "Term or Database is not found"},
					{Database: "snp", MenuName: "SNP", Count: 19, Status: "Ok"},
					{Database: "dbvar", MenuName: "dbVar", Count: 1749, Status: "Ok"},
					{Database: "epigenomics", MenuName: "Epigenomics", Count: 399, Status: "Ok"},
					{Database: "gene", MenuName: "Gene", Count: 15243, Status: "Ok"},
					{Database: "sra", MenuName: "SRA", Count: 5314, Status: "Ok"},
					{Database: "biosystems", MenuName: "BioSystems", Count: 130, Status: "Ok"},
					{Database: "unigene", MenuName: "UniGene", Count: 0, Status: "Term or Database is not found"},
					{Database: "cdd", MenuName: "CDD", Count: 9, Status: "Ok"},
					{Database: "clone", MenuName: "Clone", Count: 0, Status: "Term or Database is not found"},
					{Database: "unists", MenuName: "UniSTS", Count: 2, Status: "Ok"},
					{Database: "popset", MenuName: "PopSet", Count: 4221, Status: "Ok"},
					{Database: "geoprofiles", MenuName: "GEO Profiles", Count: 287541, Status: "Ok"},
					{Database: "gds", MenuName: "GEO DataSets", Count: 12262, Status: "Ok"},
					{Database: "homologene", MenuName: "HomoloGene", Count: 48, Status: "Ok"},
					{Database: "pccompound", MenuName: "PubChem Compound", Count: 53, Status: "Ok"},
					{Database: "pcsubstance", MenuName: "PubChem Substance", Count: 15170, Status: "Ok"},
					{Database: "pcassay", MenuName: "PubChem BioAssay", Count: 3614, Status: "Ok"},
					{Database: "nlmcatalog", MenuName: "NLM Catalog", Count: 254830, Status: "Ok"},
					{Database: "probe", MenuName: "Probe", Count: 43, Status: "Ok"},
					{Database: "gap", MenuName: "dbGaP", Count: 105672, Status: "Ok"},
					{Database: "proteinclusters", MenuName: "Protein Clusters", Count: 0, Status: "Term or Database is not found"},
					{Database: "bioproject", MenuName: "BioProject", Count: 2339, Status: "Ok"},
					{Database: "biosample", MenuName: "BioSample", Count: 93830, Status: "Ok"},
				},
			},
		},
	} {
		var g Global
		err := g.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(g, check.DeepEquals, t.global, check.Commentf("Test: %d", i))
	}
}

func (s *S) TestParseSpell(c *check.C) {
	for i, t := range []struct {
		retval string
		spell  Spell
	}{
		{
			`<?xml version="1.0"?>
<!DOCTYPE eSpellResult PUBLIC "-//NLM//DTD eSpellResult, 23 November 2004//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSpell.dtd">
<eSpellResult>
	<Database>pubmed</Database>
	<Query>asthmaa OR alergies</Query>
	<CorrectedQuery>asthma or allergies</CorrectedQuery>
	<SpelledQuery>
		<Original></Original>
		<Replaced>asthma</Replaced>
		<Original> OR </Original>
		<Replaced>allergies</Replaced>
	</SpelledQuery>
	<ERROR/>
</eSpellResult>
`,
			Spell{
				Database:  "pubmed",
				Query:     "asthmaa OR alergies",
				Corrected: "asthma or allergies",
				Replace: []Replacement{
					Old(""),
					New("asthma"),
					Old(" OR "),
					New("allergies"),
				},
			},
		},
	} {
		var sp Spell
		err := sp.Unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(sp, check.DeepEquals, t.spell, check.Commentf("Test: %d", i))
	}
}
