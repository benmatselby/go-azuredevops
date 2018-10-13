package azuredevops_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/benmatselby/go-azuredevops/azuredevops"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/testing"
)

// Pulled from https://github.com/google/go-github/blob/master/github/github_test.go
func setup() (client *azuredevops.Client, mux *http.ServeMux, serverURL string, teardown func()) {
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

	// The client being tested and is configured to use test server.
	client = azuredevops.NewClient("AZURE_DEVOPS_Account", "AZURE_DEVOPS_Project", "AZURE_DEVOPS_TOKEN")

	url, _ := url.Parse(server.URL + baseURLPath)
	client.BaseURL = url.String()
	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, expectedBody []byte) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if !bytes.Equal(body, expectedBody) {
		t.Errorf("Request body expected %s got %s", expectedBody, body)
	}
}

func testURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL; got.String() != want {
		t.Errorf("request URL is %s, want %s", got, want)
	}
}

func Test_NewClient(t *testing.T) {
	c := azuredevops.NewClient("AZURE_DEVOPS_ACCOUNT", "AZURE_DEVOPS_Project", "AZURE_DEVOPS_TOKEN")

	if c.Account != "AZURE_DEVOPS_ACCOUNT" {
		t.Errorf("Client.Account = %s; expected %s", c.Account, "AZURE_DEVOPS_ACCOUNT")
	}

	if c.Project != "AZURE_DEVOPS_Project" {
		t.Errorf("Client.Project = %s; expected %s", c.Project, "AZURE_DEVOPS_Project")
	}

	if c.AuthToken != "AZURE_DEVOPS_TOKEN" {
		t.Errorf("Client.Token = %s; expected %s", c.AuthToken, "AZURE_DEVOPS_TOKEN")
	}
}
