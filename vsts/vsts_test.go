package vsts

import (
	"testing"
)

func TestClient_New(t *testing.T) {
	c := NewClient(
		"my-account",
		"my-project",
		"my-token",
	)

	if c.Account != "my-account" {
		t.Errorf("Client.Account = %s; expected %s", c.Account, "my-account")
	}

	if c.Project != "my-project" {
		t.Errorf("Client.Project = %s; expected %s", c.Project, "my-project")
	}

	if c.AuthToken != "my-token" {
		t.Errorf("Client.Token = %s; expected %s", c.AuthToken, "my-token")
	}
}
