// Copyright ©2014 The bíogo.ncbi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphic provides simple BLAST report graphic rendering.
package graphic

import (
	"code.google.com/p/biogo.ncbi/blast"
	"code.google.com/p/plotinum/vg"

	"fmt"
	"image/color"
	"math"
)

const maxInt = int(^uint(0) >> 1)

var (
	black  = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	purple = color.RGBA{R: 0xc4, G: 0x00, B: 0xff, A: 0xff}
	blue   = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	cyan   = color.RGBA{R: 0x00, G: 0xff, B: 0xff, A: 0xff}
	green  = color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
	yellow = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	orange = color.RGBA{R: 0xff, G: 0xc4, B: 0x00, A: 0xff}
	red    = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	grey   = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
	white  = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	pallete = colors{
		color.RGBA{},
		color.RGBA{},
		grey,
		red,
		orange,
		yellow,
		green,
		cyan,
		blue,
		purple,
		black,
	}
)

type colors []color.Color

func (c colors) color(f float64) color.Color {
	if f >= 1 {
		return black
	}
	return pallete[int(f*10)]
}

const (
	maxNameLen = 20

	header vg.Length = 40
	footer vg.Length = 20

	gutter      vg.Length = 5
	leftMargin  vg.Length = 150
	bodyWidth   vg.Length = 600
	rightMargin vg.Length = 50

	lineWidth vg.Length = 3
	hitGap    vg.Length = 16
	hspGap    vg.Length = 20

	fontName = "Helvetica"

	tinySize   = vg.Length(7)
	smallSize  = vg.Length(12)
	mediumSize = vg.Length(13)
	largeSize  = vg.Length(16)
)

func mustMakeFont(font vg.Font, err error) vg.Font {
	if err != nil {
		panic(err)
	}
	return font
}

var (
	tinyFont   = mustMakeFont(vg.MakeFont(fontName, tinySize))
	smallFont  = mustMakeFont(vg.MakeFont(fontName, smallSize))
	mediumFont = mustMakeFont(vg.MakeFont(fontName+"-Bold", mediumSize))
	largeFont  = mustMakeFont(vg.MakeFont(fontName, mediumSize))
)

func fracId(id, ln *int) float64 {
	if id == nil || ln == nil {
		return math.NaN()
	}
	return float64(*id) / float64(*ln)
}

// A Summary can display a graphical summary of a blast output result.
type Summary struct {
	// Legend, Aligns and Depths specify whether the legend, alignments
	// and depth plots are drawn. All are set to true by NewSummary.
	Legend bool
	Aligns bool
	Depths bool

	w, h  vg.Length
	scale float64

	query  string
	labels []string
	hits   map[string][]hspSum

	n        int
	min, max int
}

type hspSum struct {
	queryFrom int
	queryTo   int
	hitFrom   int
	hitTo     int
	identity  float64
}

// NewSummary returns a Summary of the provided blast output.
func NewSummary(o blast.Output) Summary {
	if o.Program == "" {
		return Summary{}
	}

	var (
		n int

		labels []string
		hits   = make(map[string][]hspSum)

		min = maxInt
		max = 0
	)
	for _, i := range o.Iterations {
		for _, h := range i.Hits {
			if _, ok := hits[h.Id]; !ok {
				labels = append(labels, h.Id)
			}
			for _, hsp := range h.Hsps {
				if hsp.QueryFrom > hsp.QueryTo {
					hsp.QueryFrom, hsp.QueryTo = hsp.QueryTo, hsp.QueryFrom
					hsp.HitFrom, hsp.HitTo = hsp.HitTo, hsp.HitFrom
				}
				if hsp.QueryTo > max {
					max = hsp.QueryTo
				}
				if hsp.QueryFrom < min {
					min = hsp.QueryFrom
				}
				hits[h.Id] = append(hits[h.Id], hspSum{
					queryFrom: hsp.QueryFrom,
					queryTo:   hsp.QueryTo,
					hitFrom:   hsp.HitFrom,
					hitTo:     hsp.HitTo,
					identity:  fracId(hsp.HspIdentity, hsp.AlignLen),
				})
				n++
			}
		}
	}

	if min >= max {
		min = 1
		max = o.QueryLen
	}

	return Summary{
		Legend: true,
		Aligns: true,
		Depths: true,

		scale: float64(bodyWidth) / float64(max-min),

		query:  o.QueryId,
		labels: labels,
		hits:   hits,
		n:      n,
		min:    min,
		max:    max,
	}
}

// Render returns a vg.Canvas that has had the receiver's summary information
// rendered to it. The function cf must return a vg.Canvas that is x by h in size.
func (s Summary) Render(cf func(w, h vg.Length) vg.Canvas) vg.Canvas {
	w := leftMargin + bodyWidth + rightMargin
	h := header + footer
	if s.Aligns {
		h += hitGap*vg.Length(s.n) + hspGap*vg.Length(len(s.hits))
	}

	c := cf(w, h)
	s.w = w
	s.h = h

	s.renderHeader(c)
	if s.Legend && s.Aligns && s.n > 0 {
		s.renderLegend(c)
	}
	s.renderAlignments(c)

	return c
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// invert inverts the y-axis coordinates so that the origin is at the top of the page.
func (s Summary) invert(y vg.Length) vg.Length { return s.h - y }

// renderHeader renders the header to c.
func (s Summary) renderHeader(c vg.Canvas) {
	var p vg.Path
	p.Move(0, 0)
	p.Line(s.w, 0)
	p.Line(s.w, s.h)
	p.Line(0, s.h)
	p.Close()
	c.SetColor(white)
	c.Fill(p)

	c.SetColor(black)
	c.FillString(mediumFont, gutter, s.invert(header-smallSize), s.query[:min(len(s.query), maxNameLen)])

	p = p[:0]
	p.Move(leftMargin, s.invert(header))
	p.Line(leftMargin+bodyWidth, s.invert(header))
	c.SetLineWidth(1)
	c.Stroke(p)

	c.FillString(smallFont, leftMargin, s.invert(header-smallSize), fmt.Sprint(s.min))
	queryMax := fmt.Sprint(fmt.Sprint(s.max))
	offQueryMax := smallFont.Width(queryMax)
	c.FillString(smallFont, leftMargin+bodyWidth-offQueryMax, s.invert(header-smallSize), queryMax)
}

// renderLegend renders the legend to c.
func (s Summary) renderLegend(c vg.Canvas) {
	c.SetColor(black)
	c.FillString(smallFont, leftMargin+bodyWidth-80, s.invert(smallSize+2), "% Identity")
	var p vg.Path
	for i := 20; i <= 100; i += 10 {
		x := leftMargin + bodyWidth/2 + vg.Length(i)*2

		p.Move(x, s.invert(gutter))
		p.Line(x+10, s.invert(gutter))
		p.Line(x+10, s.invert(gutter+10))
		p.Line(x, s.invert(gutter+10))
		p.Close()
		c.SetColor(pallete.color(float64(i) / 100))
		c.Fill(p)
		p = p[:0]
		c.SetColor(black)
		c.FillString(tinyFont, x, s.invert(tinySize+gutter+10), fmt.Sprint(i))
	}
}

// xCoordOf returns the scales and translated x-coordinate of a query position.
func (s Summary) xCoordOf(p int) vg.Length {
	return vg.Length(p-s.min)*vg.Length(s.scale) + leftMargin
}

// renderAlignments renders the alignments to c.
func (s Summary) renderAlignments(c vg.Canvas) {
	var (
		v     vg.Length
		depth = make(map[int]int)
		p     vg.Path
	)

	for _, id := range s.labels {
		v += hspGap
		if s.Aligns {
			c.SetColor(black)
			c.FillString(smallFont, 2*gutter, s.invert(header+v+2*tinySize), id[:min(len(id), maxNameLen)])
		}

		c.SetLineWidth(lineWidth)
		for _, hsp := range s.hits[id] {
			v += hitGap
			x1 := s.xCoordOf(hsp.queryFrom)
			x2 := s.xCoordOf(hsp.queryTo)
			y := header + v
			for i := int(x1); i < int(x2); i++ {
				depth[i]++
			}

			if s.Aligns {
				p.Move(x1, s.invert(y-smallSize/1.5))
				p.Line(x2, s.invert(y-smallSize/1.5))
				c.SetColor(pallete.color(hsp.identity))
				c.Stroke(p)
				p = p[:0]

				c.SetColor(black)
				c.FillString(tinyFont, x1, s.invert(y-2*tinySize), fmt.Sprint(hsp.queryFrom))
				queryTo := fmt.Sprint(hsp.queryTo)
				offQueryTo := tinyFont.Width(queryTo)
				c.FillString(tinyFont, x2-offQueryTo, s.invert(y-2*tinySize), queryTo)
				c.FillString(tinyFont, x1, s.invert(y), fmt.Sprint(hsp.hitFrom))
				hitTo := fmt.Sprint(hsp.hitTo, strand(hsp.hitTo-hsp.hitFrom))
				offHitTo := tinyFont.Width(hitTo)
				c.FillString(tinyFont, x2-offHitTo, s.invert(y), hitTo)
			}
		}
	}

	if s.Depths && s.n > 0 {
		c.SetLineWidth(0.25)
		var (
			maxd int
			min  = maxInt
			max  = 0
		)
		for p, d := range depth {
			if d > maxd {
				maxd = d
			}
			if p > max {
				max = p
			}
			if p < min {
				min = p
			}
		}
		dScale := maxd/10 + 1

		const offset = 3

		c.FillString(tinyFont, leftMargin+bodyWidth+offset, s.invert(header+tinySize), fmt.Sprintf("%d/line", dScale))
		for i := min; i <= max; i++ {
			d, ok := depth[i]
			if !ok {
				continue
			}
			l := vg.Length(d / dScale)
			p.Move(vg.Length(i), s.invert(header+offset))
			p.Line(vg.Length(i), s.invert(header+offset+l))
			c.Stroke(p)
			p = p[:0]
		}
	}
}

type strand int

func (s strand) String() string {
	if s < 0 {
		return "-"
	}
	return "+"
}
