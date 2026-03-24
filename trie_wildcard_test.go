package nanoserve

import (
	"testing"
)

func TestWildcardSupport(t *testing.T) {
	router := NewTrieRouter()

	handlerCalled := false
	handler := func(c *Context) error {
		handlerCalled = true
		return c.Text("Hello")
	}

	router.Insert("GET", "/user/*", handler)

	// Test case 1: Exact match on segment before wildcard
	match := router.Search("GET", "/user/profile")
	if match == nil || len(match.Handler) == 0 {
		t.Fatal("Expected match for /user/profile")
	}
	match.Handler[0](nil)
	if !handlerCalled {
		t.Error("Handler not called for /user/profile")
	}

	handlerCalled = false
	// Test case 2: Deeper match
	match = router.Search("GET", "/user/settings/privacy")
	if match == nil || len(match.Handler) == 0 {
		t.Fatal("Expected match for /user/settings/privacy")
	}
	match.Handler[len(match.Handler)-1](nil)
	if !handlerCalled {
		t.Error("Handler not called for /user/settings/privacy")
	}

	// Test case 3: Middleware with wildcard
	middlewareCalled := false
	middleware := func(c *Context) error {
		middlewareCalled = true
		return nil
	}
	router.AddMiddleware("/admin/*", middleware)
	router.Insert("GET", "/admin/dashboard", func(c *Context) error {
		return nil
	})

	match = router.Search("GET", "/admin/dashboard")
	if match == nil || len(match.Handler) < 1 {
		t.Fatal("Expected match for /admin/dashboard")
	}

	// Check if middleware is in the handlers chain
	foundMiddleware := false
	for _, h := range match.Handler {
		// We can't easily compare functions in Go, but we can call them and check side effects
		h(nil)
		if middlewareCalled {
			foundMiddleware = true
			break
		}
	}
	if !foundMiddleware {
		t.Error("Middleware not called for /admin/dashboard")
	}
}
