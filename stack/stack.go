// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stack is a helper package for the entrez package.
package stack

import (
	"fmt"
)

// stack is used for accounting XML tags during parsing.
type Stack []string

func (st Stack) Drop() Stack { return st[:len(st)-1] }
func (st Stack) Empty() bool { return len(st) == 0 }
func (st Stack) Pair(s string) (Stack, error) {
	t, st := st[len(st)-1], st[:len(st)-1]
	if s != t {
		return st, fmt.Errorf("entrez: tag name mismatch %q != %q", s, t)

	}
	return st, nil
}
func (st Stack) Peek(i int) string {
	i++
	if i > len(st) {
		return ""
	}
	return st[len(st)-i]
}
func (st Stack) Pop() (string, Stack) { return st[len(st)-1], st[:len(st)-1] }
func (st Stack) Push(s string) Stack  { return append(st, s) }
