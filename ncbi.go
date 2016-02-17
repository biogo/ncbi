// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ncbi provides support for interaction with the NCBI services, Entrez and Blast.
//
// Please check the relevant usage policy when using these services. Note that the Blast
// and Entrez server requests are subject to frequency limits.
//
// Required parameters are specified by name in the function call.
//
// The following two parameters should be included in all requests.
//
//  tool   Name of application making the call. Its value must be a string with no internal
//         spaces.
//
//  email  E-mail address of the user. Its value must be a string with no internal spaces,
//         and should be a valid e-mail address.
package ncbi

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var client = http.Client{Timeout: 10 * time.Second}

// SetTimeout sets the HTTP client timeout duration. The default timeout is 10 seconds.
func SetTimeout(d time.Duration) {
	client.Timeout = d
}

// Limiter implements a thread-safe event frequency limit.
type Limiter struct {
	m     sync.Mutex
	delay time.Duration
	next  time.Time
}

// NewLimiter returns a Limiter that will wait for the specified duration between Wait calls.
func NewLimiter(d time.Duration) *Limiter {
	return &Limiter{delay: d}
}

// Wait blocks until the Limiter's specified duration has passed since the last Wait call.
func (d *Limiter) Wait() {
	d.m.Lock()
	defer d.m.Unlock()
	now := time.Now()
	if d.next.After(now) {
		time.Sleep(d.next.Sub(now))
		now = time.Now()
	}
	d.next = now.Add(d.delay)
}

// Util implements low level request generator for interaction with the NCBI services. It is the
// clients responsibility to provide appropriate program parameters and deserialise the returned
// data using the appropriate unmarshaling method.
type Util string

// NewRequest returns an http.Request for the utility, ut using the given method. Parameters to be
// sent to the utility program should be places in db, v, tool and email. NewRequest is subject to
// a limit that prevents requests being sent more frequently than allowed by l. The limit is easy to
// circumvent, though circumvention may result in IP blocking by the NCBI servers, so please do not
// do this.
func (ut Util) NewRequest(method, db string, v url.Values, tool, email string, l *Limiter) (*http.Request, error) {
	if db != "" {
		v["db"] = []string{db}
	}
	u, err := ut.Prepare(v, tool, email)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	l.Wait()
	return req, nil
}

// Prepare constructs a URL with the base provided by ut and the parameters provided by v, tool and email.
func (ut Util) Prepare(v url.Values, tool, email string) (*url.URL, error) {
	u, err := url.Parse(string(ut))
	if err != nil {
		return nil, err
	}
	if tool != "" {
		v["tool"] = []string{tool}
	}
	if email != "" {
		v["email"] = []string{email}
	}
	u.RawQuery = v.Encode()
	return u, nil
}

// GetMethodLimit is the maximum length of a constructed URL that will be retrieved by
// the high level API functions using the GET method.
var GetMethodLimit = 2048

// GetResponse performs a GET or POST method call to the URI in ut, passing the parameters in v,
// tool and email. The decision on which method to use is based on the length of the
// constructed URL the value of GetMethodLimit. An http.Response is returned for a successful
// request. It is the caller's responsibility to close the response body.
func (ut Util) GetResponse(v url.Values, tool, email string, l *Limiter) (*http.Response, error) {
	u, err := ut.Prepare(v, tool, email)
	if err != nil {
		return nil, err
	}
	l.Wait()
	if len(ut)+len(u.RawQuery) < GetMethodLimit {
		return client.Get(u.String())
	}
	buf := strings.NewReader(u.RawQuery)
	u.RawQuery = ""
	return client.Post(u.String(), "", buf)
}

// GetXML performs a GET or POST method call to the URI in ut, passing the parameters in v,
// tool and email. The returned stream is unmarshaled into d. The decision on which
// method to use is based on the length of the constructed URL the value of GetMethodLimit.
func (ut Util) GetXML(v url.Values, tool, email string, l *Limiter, d interface{}) error {
	resp, err := ut.GetResponse(v, tool, email, l)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return xml.NewDecoder(resp.Body).Decode(d)
}

// Get performs a GET or POST method call to the URI in ut, passing the parameters in v,
// tool and email. The decision on which method to use is based on the length of the
// constructed URL the value of GetMethodLimit. An io.ReadCloser is returned for a successful
// request. It is the caller's responsibility to close this.
func (ut Util) Get(v url.Values, tool, email string, l *Limiter) (io.ReadCloser, error) {
	resp, err := ut.GetResponse(v, tool, email, l)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
