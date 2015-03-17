// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"strings"

	"gopkg.in/check.v1"
)

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
				History: &History{
					QueryKey: 1,
					WebEnv:   "NCID_1_298287560_130.14.18.48_5555_1360293704_91337037",
				},
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
				History: &History{
					QueryKey: 1,
					WebEnv:   "NCID_1_299062774_130.14.18.97_5555_1360293760_1713152879",
				},
				Err: stringPtr("IDs contain invalid characters which was treated as delimiters."),
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
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&p)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(p, check.DeepEquals, t.post, check.Commentf("Test: %d", i))
	}
}
