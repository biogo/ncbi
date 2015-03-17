// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"strings"

	"github.com/biogo/ncbi/entrez/global"

	"gopkg.in/check.v1"
)

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
				Results: []global.Result{
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
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&g)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(g, check.DeepEquals, t.global, check.Commentf("Test: %d", i))
	}
}
