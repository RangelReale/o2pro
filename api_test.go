// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/restclient"
	"net/http/httptest"
	"net/url"
	"testing"
)

// An Authorizer implementation that always authorizes owner "jtkirk", and never
// authorizes anyone else.
func kirkAuthorizer(u *url.Userinfo, r AuthRequest) (bool, error) {
	if u.Username() == "jtkirk" {
		return true, nil
	}
	return false, nil
}

func TestAuthRequest(t *testing.T) {
	//
	// Prepare handler
	//
	s, _ := setup(t)
	s.Authorizer = kirkAuthorizer
	h := s.AuthReqHandler()
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	// rec := httptest.NewRecorder()
	//
	// REST request
	//
	scopes := []string{"enterprise", "intrepid"}
	u := url.UserPassword("jtkirk", "Beam me up, Scotty!")
	areq := AuthRequest{
		Scopes: scopes,
		Note:   "foo bar baz",
	}
	var a Authorization
	var e interface{}
	rr := restclient.RequestResponse{
		Url:      hserv.URL,
		Method:   "POST",
		Userinfo: u,
		Data:     &areq,
		Result:   &a,
		Error:    &e,
	}
	c := restclient.New()
	c.UnsafeBasicAuth = true
	status, err := c.Do(&rr)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 200, status)
	assert.Equal(t, "jtkirk", a.Owner)
	assert.NotEqual(t, nil, uuid.Parse(a.Token))
	for _, scope := range scopes {
		_, ok := a.Scopes[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}