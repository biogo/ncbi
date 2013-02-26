// Copyright ©2013 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blast

import (
	"code.google.com/p/biogo.ncbi/html"
	"io"
)

type Info string

func (i *Info) unmarshal(r io.Reader) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return z.Err()
		}
		if tt == html.CommentToken {
			*i = Info(z.Token().Data)
			return nil
		}
	}

	panic("cannot reach")
}