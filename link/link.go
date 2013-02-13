// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// <!--
//                 This is the Current DTD for Entrez eLink
// $Id: eLink_101123.dtd 349314 2012-01-09 23:26:00Z fialkov $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT	ERROR			(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Info			(#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT	Id				(#PCDATA)>	<!-- \d+ -->
// <!ATTLIST	Id
// 			HasLinkOut  (Y|N)	#IMPLIED
// 			HasNeighbor (Y|N)	#IMPLIED
// 			>
//
// <!ELEMENT	Score			(#PCDATA)>	<!-- \d+ -->
// <!ELEMENT	DbFrom			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	DbTo			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	LinkName		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	WebEnv			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	MenuTag			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	HtmlTag			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	Priority		(#PCDATA)>	<!-- \S+ -->
//
// <!ELEMENT	IdList		(Id*)>
//
// <!-- cmd=neighbor -->
// <!ELEMENT	Link		(Id, Score?)>
// <!ELEMENT	QueryKey		(#PCDATA)>
//
// <!ELEMENT	LinkSetDb	(DbTo, LinkName, (Link*|Info), ERROR?)>
// <!ELEMENT	LinkSetDbHistory	(DbTo, LinkName, (QueryKey|Info), ERROR?)>
//
// <!-- cmd=llinks -->
//
// <!ELEMENT	Url			    (#PCDATA)>	<!-- \S+ -->
// <!ATTLIST	Url			LNG     (DA|DE|EN|EL|ES|FR|IT|IW|JA|NL|NO|RU|SV|ZH)     "EN">
//
// <!ELEMENT	IconUrl			(#PCDATA)>	<!-- \S+ -->
// <!ATTLIST	IconUrl		LNG     (DA|DE|EN|EL|ES|FR|IT|IW|JA|NL|NO|RU|SV|ZH)     "EN">
//
// <!ELEMENT	SubjectType		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Category		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Attribute		(#PCDATA)>	<!-- .+ -->
// <!--ELEMENT	LinkName		(#PCDATA)-->	<!--defined in neighbor section--><!-- \S+ -->
// <!ELEMENT	Name			(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	NameAbbr		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	SubProvider		(#PCDATA)>
//
// <!ELEMENT   FirstChar		(#PCDATA)>
//
// <!ELEMENT	Provider (
// 				Name,
// 				NameAbbr,
// 				Id,
// 				Url,
// 				IconUrl?
// 			)>
//
// <!ELEMENT	ObjUrl	(
// 				Url,
// 				IconUrl?,
// 				LinkName?,
//              SubjectType*,
// 				Category*,
//              Attribute*,
//              Provider,
//              SubProvider?
// 			)>
//
// <!ELEMENT	IdUrlSet	(Id,(ObjUrl+|Info))>
//
// <!ELEMENT   FirstChars  (FirstChar*)>
//
// <!ELEMENT	LinkInfo	(DbTo, LinkName, MenuTag?, HtmlTag?, Url?, Priority)>
// <!ELEMENT	IdLinkSet	(Id, LinkInfo*)>
// <!ELEMENT	IdUrlList	(IdUrlSet* | FirstChars*)>
//
// <!-- cmd=ncheck & lcheck & acheck -->
// <!ELEMENT	IdCheckList	((Id|IdLinkSet)*,ERROR?)>
//
// <!-- Common -->
// <!ELEMENT	LinkSet		(DbFrom,
// 				((IdList?, ((ERROR?, LinkSetDb)*  |  (LinkSetDbHistory*, WebEnv))) | IdUrlList | IdCheckList | ERROR), ERROR?
// 				)>
//
// <!ELEMENT	eLinkResult	(LinkSet*, ERROR?)>

// What does 'IMPLIED' mean on a boolean?

type LinkId struct {
	Id          int
	HasLinkOut  bool
	HasNeighbor bool
	Score       int
}

func (li *LinkId) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			if t.Name.Local == "Id" {
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "HasLinkOut":
						switch b := attr.Value; b {
						case "Y", "N":
							li.HasLinkOut = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					case "HasNeighbor":
						switch b := attr.Value; b {
						case "Y", "N":
							li.HasNeighbor = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					default:
						return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
					}
				}
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				li.Id = id
			case "Score":
				s, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				li.Score = s
			case "Link":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type LinkSetDb struct {
	DbTo     string
	LinkName string
	Link     []LinkId
}

func (ls *LinkSetDb) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			if t.Name.Local == "Link" {
				var li LinkId
				err := li.Unmarshal(dec, st[len(st)-1:])
				if (li != LinkId{}) {
					ls.Link = append(ls.Link, li)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
				continue
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "DbTo":
				ls.DbTo = string(t)
			case "LinkName":
				ls.LinkName = string(t)
			case "LinkSetDb":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type Url struct {
	Url  string
	Lang string
}

func (u *Url) getLang(t xml.StartElement) error {
	for _, attr := range t.Attr {
		switch attr.Name.Local {
		case "LNG":
			switch attr.Value {
			case "DA", "DE", "EN", "EL", "ES", "FR", "IT", "IW", "JA", "NL", "NO", "RU", "SV", "ZH":
				u.Lang = attr.Value
			default:
				return fmt.Errorf("eutil: unknown language id: %q", attr.Value)
			}
		default:
			return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
		}
	}
	return nil
}

type Provider struct {
	Name     string
	NameAbbr string
	Id       int
	Url      Url
	IconUrl  *Url
}

func (p *Provider) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			switch t.Name.Local {
			case "Url":
				err = p.Url.getLang(t)
			case "IconUrl":
				p.IconUrl = &Url{}
				err = p.IconUrl.getLang(t)
			}
			if err != nil {
				return err
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Name":
				p.Name = string(t)
			case "NameAbbr":
				p.NameAbbr = string(t)
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				p.Id = id
			case "Url":
				p.Url.Url = string(t)
			case "IconUrl":
				p.IconUrl.Url = string(t)
			case "Provider":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type ObjUrl struct {
	Url         Url
	IconUrl     *Url
	LinkName    *string
	SubjectType []string
	Category    []string
	Attribute   []string
	Provider    Provider
	SubProvider *string
}

func (ou *ObjUrl) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			switch t.Name.Local {
			case "Provider":
				err := ou.Provider.Unmarshal(dec, st[len(st)-1:])
				if err != nil {
					return err
				}
				st = st.Drop()
			case "Url":
				err = ou.Url.getLang(t)
			case "IconUrl":
				ou.IconUrl = &Url{}
				err = ou.IconUrl.getLang(t)
			}
			if err != nil {
				return err
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Url":
				ou.Url.Url = string(t)
			case "IconUrl":
				ou.IconUrl.Url = string(t)
			case "LinkName":
				s := string(t)
				ou.LinkName = &s
			case "SubjectType":
				ou.SubjectType = append(ou.SubjectType, string(t))
			case "Category":
				ou.Category = append(ou.Category, string(t))
			case "Attribute":
				ou.Attribute = append(ou.Attribute, string(t))
			case "SubProvider":
				s := string(t)
				ou.SubProvider = &s
			case "ObjUrl":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type IdUrlSet struct {
	Id      int
	ObjUrls []ObjUrl
	Info    string
}

func (us *IdUrlSet) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			if t.Name.Local == "ObjUrl" {
				var ou ObjUrl
				err := ou.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(ou, IdUrlSet{}) {
					us.ObjUrls = append(us.ObjUrls, ou)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
				continue
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				us.Id = id
			case "Info":
				us.Info = string(t)
			case "IdUrlSet", "ObjUrl":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type IdUrlList struct {
	IdUrlSets []IdUrlSet
	Err       error
}

func (ul *IdUrlList) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			if t.Name.Local == "IdUrlSet" {
				var us IdUrlSet
				err := us.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(us, IdUrlSet{}) {
					ul.IdUrlSets = append(ul.IdUrlSets, us)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
				continue
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "ERROR":
				ul.Err = errors.New(string(t))
			case "IdUrlList":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type LinkInfo struct {
	DbTo     string
	LinkName string
	MenuTag  *string
	HtmlTag  *string
	Url      *Url
	Priority string
}

func (li *LinkInfo) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			if t.Name.Local == "Url" {
				li.Url = &Url{}
				err = li.Url.getLang(t)
				if err != nil {
					return err
				}
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "DbTo":
				li.DbTo = string(t)
			case "LinkName":
				li.LinkName = string(t)
			case "MenuTag":
				s := string(t)
				li.MenuTag = &s
			case "HtmlTag":
				s := string(t)
				li.HtmlTag = &s
			case "Url":
				li.Url.Url = string(t)
			case "Priority":
				li.Priority = string(t)
			case "LinkInfo":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type Id struct {
	Id          int
	HasLinkOut  bool
	HasNeighbor bool
}

type IdLinkSet struct {
	Id       Id
	LinkInfo []LinkInfo
}

func (is *IdLinkSet) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	var hasLinkOut, hasNeighbor bool
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			switch t.Name.Local {
			case "Id":
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "HasLinkOut":
						switch b := attr.Value; b {
						case "Y", "N":
							hasLinkOut = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					case "HasNeighbor":
						switch b := attr.Value; b {
						case "Y", "N":
							hasNeighbor = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					default:
						return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
					}
				}
			case "LinkInfo":
				var li LinkInfo
				err := li.Unmarshal(dec, st[len(st)-1:])
				if (li != LinkInfo{}) {
					is.LinkInfo = append(is.LinkInfo, li)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				is.Id = Id{Id: id, HasLinkOut: hasLinkOut, HasNeighbor: hasNeighbor}
			case "IdLinkSet", "LinkInfo":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type IdChecks struct {
	Ids        []Id
	IdLinkSets []IdLinkSet
	Err        error
}

func (ic *IdChecks) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	var hasLinkOut, hasNeighbor bool
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			switch t.Name.Local {
			case "Id":
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "HasLinkOut":
						switch b := attr.Value; b {
						case "Y", "N":
							hasLinkOut = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					case "HasNeighbor":
						switch b := attr.Value; b {
						case "Y", "N":
							hasNeighbor = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					default:
						return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
					}
				}
			case "IdLinkSet":
				var is IdLinkSet
				err := is.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(is, IdLinkSet{}) {
					ic.IdLinkSets = append(ic.IdLinkSets, is)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "ERROR":
				ic.Err = errors.New(string(t))
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				ic.Ids = append(ic.Ids, Id{Id: id, HasLinkOut: hasLinkOut, HasNeighbor: hasNeighbor})
			case "IdCheckList":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

type LinkSet struct {
	DbFrom     string
	IdList     []int
	LinkSetDbs []LinkSetDb
	IdUrls     []IdUrlList
	IdChecks   *IdChecks
	Err        error
}

func (li *LinkSet) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			switch t.Name.Local {
			case "IdUrlList":
				var ul IdUrlList
				err := ul.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(ul, IdUrlList{}) {
					li.IdUrls = append(li.IdUrls, ul)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			case "LinkSetDb":
				var ls LinkSetDb
				err := ls.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(ls, LinkSetDb{}) {
					li.LinkSetDbs = append(li.LinkSetDbs, ls)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			case "IdCheckList":
				var id IdChecks
				err := id.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(id, IdChecks{}) {
					li.IdChecks = &id
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "DbFrom":
				li.DbFrom = string(t)
			case "Id":
				if st.Peek(1) != "IdList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				li.IdList = append(li.IdList, id)
			case "LinkSet", "IdList", "IdUrlList":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}
