// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"strings"

	"github.com/biogo/ncbi/entrez/info"

	"gopkg.in/check.v1"
)

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
				Err:    "",
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
				DbInfo: &info.DbInfo{
					DbName:      "toolkit",
					MenuName:    "ToolKit",
					Description: "ToolKit database",
					Count:       265403,
					LastUpdate:  "2013/02/07 14:34",
					FieldList: []info.Field{
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
				Err: "",
			},
		},
		{
			`<?xml version="1.0"?>
			<!DOCTYPE eInfoResult PUBLIC "-//NLM//DTD eInfoResult, 11 May 2002//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eInfo_020511.dtd">
			<eInfoResult>
				<ERROR>Can not retrieve DbInfo for db=blah</ERROR>
			</eInfoResult>
			`,
			Info{Err: "Can not retrieve DbInfo for db=blah"},
		},
	} {
		var in Info
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&in)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(in, check.DeepEquals, t.info, check.Commentf("Test: %d", i))
	}
}
