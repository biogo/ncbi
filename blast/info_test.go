// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

import (
	"strings"

	"gopkg.in/check.v1"
)

func (s *S) TestParseInfo(c *check.C) {
	for i, t := range []struct {
		retval string
		info   Info
	}{
		{
			`


<p><!--
QBlastInfoBegin
	Status=INFO_DB
# Number of databases
23

# exlclusive databases

nr              1       TRUE
nr              2       FALSE
est_human       3       FALSE
est_mouse       4       FALSE
est_others      5       FALSE
htg             6       FALSE
gss             7       FALSE
sts             8       FALSE
pataa		9	TRUE
patnt		10	FALSE

# non-exlclusive databases

swissprot       11      TRUE
est             12      FALSE
pdb             13      TRUE
pdb             14      FALSE
month           15      TRUE
month.nt        16      FALSE
month.est       17      FALSE
month.gss       18      FALSE
month.htgs      19      FALSE
month.sts       20      FALSE
month.pataa     21      TRUE
month.patnt     22      FALSE

# new

chromosome      23      TRUE

# end of the file

QBlastInfoEnd
--><p>



`,
			Info(`
QBlastInfoBegin
	Status=INFO_DB
# Number of databases
23

# exlclusive databases

nr              1       TRUE
nr              2       FALSE
est_human       3       FALSE
est_mouse       4       FALSE
est_others      5       FALSE
htg             6       FALSE
gss             7       FALSE
sts             8       FALSE
pataa		9	TRUE
patnt		10	FALSE

# non-exlclusive databases

swissprot       11      TRUE
est             12      FALSE
pdb             13      TRUE
pdb             14      FALSE
month           15      TRUE
month.nt        16      FALSE
month.est       17      FALSE
month.gss       18      FALSE
month.htgs      19      FALSE
month.sts       20      FALSE
month.pataa     21      TRUE
month.patnt     22      FALSE

# new

chromosome      23      TRUE

# end of the file

QBlastInfoEnd
`),
		},
		{`


UNKNOWN
<hr><p id="blastErr"><font color="red">Results for RID  not found.</font></p><hr>
<p><!--
QBlastInfoBegin
	Status=UNKNOWN
QBlastInfoEnd
--></p>




`,
			Info(`
QBlastInfoBegin
	Status=UNKNOWN
QBlastInfoEnd
`),
		},
	} {
		var info Info
		err := info.unmarshal(strings.NewReader(t.retval))
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(info, check.DeepEquals, t.info, check.Commentf("Test: %d", i))
	}
}
