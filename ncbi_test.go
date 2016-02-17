// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ncbi

import (
	"testing"
	"time"

	"gopkg.in/check.v1"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestLimiter(c *check.C) {
	var count int
	Limit := NewLimiter(time.Second / 3)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				Limit.Wait()
				count++
			}
		}()
	}
	time.Sleep(3 * time.Second)
	c.Check(count < 10, check.Equals, true)
}
