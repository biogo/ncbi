// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	check "launchpad.net/gocheck"
	"strings"
)

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
												Url:         Url{Url: "http://www.labome.org//expert/switzerland/university/klingenberg/roland-klingenberg-1568163.html"},
												LinkName:    stringPtr("Labome Researcher Resource"),
												SubjectType: []string{"author profiles"},
												Category:    []string{"Other Literature Sources"},
												Attribute:   []string{"subscription/membership/fee required"},
												Provider: Provider{
													Name:     "ExactAntigen/Labome",
													NameAbbr: "EADB",
													Id:       5753,
													Url:      Url{Url: "http://www.labome.com", Lang: "EN"},
												},
											},
											{
												Url:         Url{Url: "http://eurheartj.oxfordjournals.org/cgi/pmidlookup?view=long&pmid=19880848"},
												IconUrl:     &Url{Url: "http://www.ncbi.nlm.nih.gov/corehtml/query/egifs/http:--highwire.stanford.edu-icons-externalservices-pubmed-custom-oxfordjournals_final_free.gif", Lang: "EN"},
												SubjectType: []string{"publishers/providers"},
												Category:    []string{"Full Text Sources"},
												Attribute:   []string{"free resource", "full-text online", "publisher of information in url"},
												Provider: Provider{
													Name:     "HighWire",
													NameAbbr: "HighWire",
													Id:       3051,
													Url:      Url{Url: "http://highwire.stanford.edu", Lang: "EN"},
												},
											},
											{
												Url:         Url{Url: "http://www.nlm.nih.gov/medlineplus/atherosclerosis.html"},
												IconUrl:     &Url{Url: "http://www.ncbi.nlm.nih.gov/corehtml/query/egifs/http:--www.nlm.nih.gov-medlineplus-images-linkout_sm.gif", Lang: "EN"},
												LinkName:    stringPtr("Atherosclerosis"),
												SubjectType: []string{"consumer health"},
												Category:    []string{"Medical"},
												Attribute:   []string{"free resource"},
												Provider: Provider{
													Name:     "MedlinePlus Health Information",
													NameAbbr: "MEDPLUS",
													Id:       3162,
													Url:      Url{Url: "http://medlineplus.gov/", Lang: "EN"},
													IconUrl:  &Url{Url: "http://www.nlm.nih.gov/medlineplus/images/linkout_sm.gif", Lang: "EN"},
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
