package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to
	// them).
	app := newTestApplication(t)
	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}
