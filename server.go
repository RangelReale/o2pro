// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"code.google.com/p/go-uuid/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	DefaultExpireAfter = "8h" // Duration string for time.ParseDuration()
	DefaultLogger      = log.New(os.Stdout, "[o2pro] ", log.Ltime|log.Ldate|log.Lshortfile)
	DefaultScopes      = []string{"all"}
)

// A Storage back end saves and retrieves authorizations to persistent storage,
// perhaps with caching.
type Storage interface {
	SaveAuth(auth *Auth) error
	GetAuth(token string) (*Auth, error)
	Activate() error // Called when Server is started
}

// An Authorizer decides whether to grant an authorization request based on
// client's credentials.
type Authorizer func(username, password string, scopes []string) (bool, error)

// A Server is an authorization service that can issue Oauth2-style bearer
// tokens.
type Server struct {
	Storage
	Scopes        []string      // All scopes supported by this server
	DefaultScopes []string      // Issued if no specific scope(s) requested
	Duration      time.Duration // Lifetime for an authorization
	Logger        *log.Logger
	Authorizer    Authorizer
}

// NewAuth issues a new Authorization based on an AuthRequest.
// func (s *Server) NewAuth(owner string, req AuthRequest) (Authorization, error) {
func (s *Server) NewAuth(t AuthTemplate) (Auth, error) {
	a := Auth{
		Token:      uuid.NewUUID().String(),
		Username:   t.Username,
		Scopes:     t.Scopes,
		Expiration: time.Now().Add(s.Duration),
		Note:       t.Note,
	}
	err := s.SaveAuth(&a)
	return a, err
}

func (s *Server) Error(w http.ResponseWriter, error string, code int) {

}

// Authorize may grant an authorization to a client.  The Authorizer function
// decides whether to make the grant. ErrNotAuthorized is returned if
// authorization is denied.
func (s *Server) Authorize(t AuthTemplate, password string) (Auth, error) {
	var a Auth
	ok, err := s.Authorizer(t.Username, password, t.Scopes)
	if err != nil {
		return a, err
	}
	if !ok {
		return a, ErrNotAuthorized
	}
	return s.NewAuth(t)
}
