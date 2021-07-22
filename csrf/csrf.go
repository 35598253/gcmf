package csrf

import (
	"net/http"
	"nmrich/rcmf/core/rc"
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/grand"
)

// CsrfMiddle is the configuration struct for CSRF feature.
type CsrfMiddle struct {
	TokenRequestKey string
	CookieName      string
}

// CSRF creates and returns a CSRF middleware with incoming configuration.
func (c *CsrfMiddle) CSRF(r *ghttp.Request) {
	if c.CookieName == "" {
		c.CookieName = "rcmf_csrf"
	}
	TokenLength := 30
	if c.TokenRequestKey == "" {
		c.TokenRequestKey = "X-CSRF-Token"
	}
	// Read the token in the request cookie
	tokenInCookie := r.Cookie.Get(c.CookieName)
	if tokenInCookie == "" {
		// Generate a random token
		tokenInCookie = grand.S(TokenLength)
	}

	// Read the token attached to the request
	// Read priority: Router < Query < Body < Form < Custom < Header
	tokenInRequestData := r.Header.Get(c.TokenRequestKey)
	if tokenInRequestData == "" {
		tokenInRequestData = r.GetString(c.TokenRequestKey)
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		// No verification required
	default:
		// Authentication token
		if !strings.EqualFold(tokenInCookie, tokenInRequestData) {
			rc.JsonExit(r, 9, "CSRF ERROR")
		}
	}

	// Set cookie in response
	r.Cookie.SetCookie(c.CookieName, tokenInCookie, "", "/", 0)
	r.Middleware.Next()
}
