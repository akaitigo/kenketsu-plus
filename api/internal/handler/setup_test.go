package handler_test

import (
	"net/http"
	"os"
	"testing"
)

// testAdminKey is the admin API key configured for the handler test binary.
const testAdminKey = "test-admin-key"

// TestMain configures the admin API key for the whole test binary. RequireAdminKey
// is fail-closed (#20): without ADMIN_API_KEY set, admin-protected endpoints return
// 503. Setting it once here (rather than per test) keeps the endpoint tests parallel;
// each request still authenticates via the X-Admin-Key header (see setAdminKey).
func TestMain(m *testing.M) {
	if err := os.Setenv("ADMIN_API_KEY", testAdminKey); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// setAdminKey adds the admin API key header required by admin-protected endpoints.
func setAdminKey(req *http.Request) {
	req.Header.Set("X-Admin-Key", testAdminKey)
}
