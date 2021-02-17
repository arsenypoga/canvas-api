package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	got := NewClient("domain", "authToken")
	header := http.Header{}
	header.Add("authorization", "Bearer "+"authToken")

	expected := CanvasClient{
		Domain:  "domain",
		client:  http.DefaultClient,
		headers: &header,
	}

	assert.Equal(t, &expected, got)
}

func TestCanvasClient_ClientURL(t *testing.T) {
	got := NewClient("domain", "authToken").ClientURL()
	assert.Equal(t, "https://domain.instructure.com", got)
}
