package vsts_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/benmatselby/go-vsts/vsts"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/testing"
)

// Pulled from https://github.com/google/go-github/blob/master/github/github_test.go
func setup() (client *vsts.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the VSTS client being tested and is
	// configured to use test server.
	client = vsts.NewClient("VSTS_Account", "VSTS_Project", "VSTS_Token")

	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url.String()
	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL; got.String() != want {
		t.Errorf("request URL is %s, want %s", got, want)
	}
}

func TestVsts_NewClient(t *testing.T) {
	c := vsts.NewClient("VSTS_Account", "VSTS_Project", "VSTS_Token")

	if c.Account != "VSTS_Account" {
		t.Errorf("Client.Account = %s; expected %s", c.Account, "VSTS_Account")
	}

	if c.Project != "VSTS_Project" {
		t.Errorf("Client.Project = %s; expected %s", c.Project, "VSTS_Project")
	}

	if c.AuthToken != "VSTS_Token" {
		t.Errorf("Client.Token = %s; expected %s", c.AuthToken, "VSTS_Token")
	}
}
