package conf

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// AuthMethod is a wrapper around authorization methods for git authentication.
type AuthMethod struct {
	*http.BasicAuth
	*http.TokenAuth
}

// GetAuth returns the non nil authorization method.
func (m *AuthMethod) GetAuth() http.AuthMethod {
	if m == nil {
		return nil
	}

	if m.BasicAuth != nil {
		return m.BasicAuth
	}

	return m.TokenAuth
}

func (m *AuthMethod) String() string {
	if m == nil {
		return ""
	}

	return m.GetAuth().String()
}
