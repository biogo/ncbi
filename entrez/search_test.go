// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"strings"

	. "github.com/biogo/ncbi/entrez/search"

	"gopkg.in/check.v1"
)

func (s *S) TestParseSearch(c *check.C) {
	for i, t := range []struct {
		retval string
		search Search
		ast    Node
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
				History:  nil,
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
					&Term{
						Term:    "\"Science\"[Journal]",
						Field:   "Journal",
						Count:   162433,
						Explode: true,
					},
					&Term{
						Term:    "\"Science (80- )\"[Journal]",
						Field:   "Journal",
						Count:   10,
						Explode: true,
					},
					&Op{Operation: "OR"},
					&Term{
						Term:    "\"J Zhejiang Univ Sci\"[Journal]",
						Field:   "Journal",
						Count:   364,
						Explode: true,
					},
					&Op{Operation: "OR"},
					&Op{Operation: "GROUP"},
					&Term{
						Term:    "\"breast neoplasms\"[MeSH Terms]",
						Field:   "MeSH Terms",
						Count:   199283,
						Explode: true,
					},
					&Term{Term: "\"breast\"[All Fields]",
						Field:   "All Fields",
						Count:   322674,
						Explode: true,
					},
					&Term{
						Term:    "\"neoplasms\"[All Fields]",
						Field:   "All Fields",
						Count:   1897643,
						Explode: true,
					},
					&Op{Operation: "AND"},
					&Op{Operation: "GROUP"},
					&Op{Operation: "OR"},
					&Term{
						Term:    "\"breast neoplasms\"[All Fields]",
						Field:   "All Fields",
						Count:   199169,
						Explode: true,
					},
					&Op{Operation: "OR"},
					&Term{
						Term:    "\"breast\"[All Fields]",
						Field:   "All Fields",
						Count:   322674,
						Explode: true,
					},
					&Term{
						Term:    "\"cancer\"[All Fields]",
						Field:   "All Fields",
						Count:   1166779,
						Explode: true,
					},
					&Op{Operation: "AND"},
					&Op{Operation: "GROUP"},
					&Op{Operation: "OR"},
					&Term{
						Term:    "\"breast cancer\"[All Fields]",
						Field:   "All Fields",
						Count:   156855,
						Explode: true,
					},
					&Op{Operation: "OR"},
					&Op{Operation: "GROUP"},
					&Op{Operation: "AND"},
					&Term{
						Term:    "2008[pdat]",
						Field:   "pdat",
						Count:   828593,
						Explode: true,
					},
					&Op{Operation: "AND"},
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
				Err: nil,
			},
			&Op{
				Operation: "AND",
				Operands: []Node{
					&Op{
						Operation: "AND",
						Operands: []Node{
							&Op{
								Operation: "GROUP",
								Operands: []Node{
									&Op{
										Operation: "OR",
										Operands: []Node{
											&Op{
												Operation: "OR",
												Operands: []Node{
													&Term{
														Term:    "\"Science\"[Journal]",
														Field:   "Journal",
														Count:   162433,
														Explode: true,
													},
													&Term{
														Term:    "\"Science (80- )\"[Journal]",
														Field:   "Journal",
														Count:   10,
														Explode: true,
													},
												},
											},
											&Term{
												Term:    "\"J Zhejiang Univ Sci\"[Journal]",
												Field:   "Journal",
												Count:   364,
												Explode: true,
											},
										},
									},
								},
							},
							&Op{
								Operation: "GROUP",
								Operands: []Node{
									&Op{
										Operation: "OR",
										Operands: []Node{
											&Op{
												Operation: "OR",
												Operands: []Node{
													&Op{
														Operation: "OR",
														Operands: []Node{
															&Op{
																Operation: "OR",
																Operands: []Node{
																	&Term{
																		Term:    "\"breast neoplasms\"[MeSH Terms]",
																		Field:   "MeSH Terms",
																		Count:   199283,
																		Explode: true,
																	},
																	&Op{
																		Operation: "GROUP",
																		Operands: []Node{
																			&Op{
																				Operation: "AND",
																				Operands: []Node{
																					&Term{
																						Term:    "\"breast\"[All Fields]",
																						Field:   "All Fields",
																						Count:   322674,
																						Explode: true,
																					},
																					&Term{
																						Term:    "\"neoplasms\"[All Fields]",
																						Field:   "All Fields",
																						Count:   1897643,
																						Explode: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
															&Term{
																Term:    "\"breast neoplasms\"[All Fields]",
																Field:   "All Fields",
																Count:   199169,
																Explode: true,
															},
														},
													},
													&Op{
														Operation: "GROUP",
														Operands: []Node{
															&Op{
																Operation: "AND",
																Operands: []Node{
																	&Term{
																		Term:    "\"breast\"[All Fields]",
																		Field:   "All Fields",
																		Count:   322674,
																		Explode: true,
																	},
																	&Term{
																		Term:    "\"cancer\"[All Fields]",
																		Field:   "All Fields",
																		Count:   1166779,
																		Explode: true,
																	},
																},
															},
														},
													},
												},
											},
											&Term{
												Term:    "\"breast cancer\"[All Fields]",
												Field:   "All Fields",
												Count:   156855,
												Explode: true,
											},
										},
									},
								},
							},
						},
					},
					&Term{
						Term:    "2008[pdat]",
						Field:   "pdat",
						Count:   828593,
						Explode: true,
					},
				},
			},
		},
		{`<?xml version="1.0" ?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult>
	<ERROR>Empty term and query_key - nothing todo</ERROR>
</eSearchResult>
`,
			Search{Err: stringPtr("Empty term and query_key - nothing todo")},
			nil,
		},
		{`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE eSearchResult PUBLIC "-//NLM//DTD eSearchResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSearch_020511.dtd">
<eSearchResult>
	<ERROR>Invalid db name specified: pub</ERROR>
</eSearchResult>
`,
			Search{Err: stringPtr("Invalid db name specified: pub")},
			nil,
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
				History:          nil,
				IdList:           nil,
				Translations:     nil,
				TranslationStack: nil,
				QueryTranslation: nil,
				Err:              nil,
				NotFound: &NotFound{
					Field: []string{"jungle", "pat"},
				},
			},
			nil,
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
				History:          nil,
				IdList:           nil,
				Translations:     nil,
				TranslationStack: nil,
				QueryTranslation: stringPtr("(nonjournal[journal] AND nonyear[date])"),
				Err:              nil,
				NotFound: &NotFound{
					Phrase: []string{"nonjournal[journal]", "nonyear[date]"},
				},
				Warnings: &Warnings{
					Message: []string{"No items found."},
				},
			},
			nil,
		},
	} {
		var s Search
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&s)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(s, check.DeepEquals, t.search, check.Commentf("Test: %d", i))
		if s.TranslationStack != nil {
			n, err := s.TranslationStack.AST()
			c.Check(err, check.Equals, nil)
			c.Check(n, check.DeepEquals, t.ast)
		}
	}
}
