// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"strings"

	. "github.com/biogo/ncbi/entrez/summary"

	"gopkg.in/check.v1"
)

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
				Documents: []Document{
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
			},
		},
	} {
		var s Summary
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&s)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(s, check.DeepEquals, t.summary, check.Commentf("Test: %d", i))
	}
}
