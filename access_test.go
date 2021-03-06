// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/restclient"
	"net/http"
	"net/http/httptest"
	"testing"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func doTestRequireScope(p *Provider, t *testing.T) {
	h := p.RequireScope(fooHandler, "enterprise")
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	//
	// Valid Scope
	//
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	note := "foo bar baz"
	auth, _ := p.NewAuthz(username, note, scopes)
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+auth.Token)
	rr := restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
		Header: &header,
	}
	status, err := restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, status)
	//
	// No Token
	//
	rr = restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
	}
	status, err = restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
	//
	// Bad Header
	//
	header = make(http.Header)
	header.Add("Authorization", "foobar")
	rr = restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
		Header: &header,
	}
	status, err = restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
	//
	// Unauthorized Scope
	//
	h1 := p.RequireScope(fooHandler, "foobar") // Not among the authorized scopes
	hserv1 := httptest.NewServer(h1)
	defer hserv1.Close()
	header = make(http.Header)
	header.Add("Authorization", "Bearer "+auth.Token)
	rr = restclient.RequestResponse{
		Url:    hserv1.URL,
		Method: "GET",
		Header: &header,
	}
	status, err = restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
}

func doTestRequireAuthc(p *Provider, t *testing.T) {
	h := p.RequireAuthc(fooHandler)
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	//
	// Valid Authentication
	//
	auth, _ := p.NewAuthz("jtkirk", "", nil)
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+auth.Token)
	rr := restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
		Header: &header,
	}
	status, err := restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, status)
	//
	// Invalid Auth Token
	//
	header = make(http.Header)
	header.Add("Authorization", "Bearer foorbar") // "foobar" is not a valid token
	rr = restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
		Header: &header,
	}
	status, err = restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
	//
	// No Auth Token
	//
	rr = restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
	}
	status, err = restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
}
