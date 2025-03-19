package recallaigo_test

import (
	"net/http"
	"os"
	"testing"

	recallaigo "github.com/harrison-peng/recallai-go"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// newTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// newMockedClient returns *http.Client which responds with content from given file
func newMockedClient(t *testing.T, requestMockFile string, statusCode int) *http.Client {
	return newTestClient(func(*http.Request) *http.Response {
		b, err := os.Open(requestMockFile)
		if err != nil {
			t.Fatal(err)
		}

		resp := &http.Response{
			StatusCode: statusCode,
			Body:       b,
			Header:     make(http.Header),
		}
		return resp
	})
}

func TestNewClient(t *testing.T) {
	token := recallaigo.Token("test-token")

	client := recallaigo.NewClient(token)

	if client.Token != token {
		t.Errorf("expected token %s, got %s", token, client.Token)
	}

	if client.Region != recallaigo.UsEast {
		t.Errorf("expected default region %s, got %s", recallaigo.UsEast, client.Region)
	}
}

func TestNewClientWithCustomRegion(t *testing.T) {
	token := recallaigo.Token("test-token")
	customRegion := recallaigo.UsWest

	client := recallaigo.NewClient(token, recallaigo.WithRegion(customRegion))

	if client.Token != token {
		t.Errorf("expected token %s, got %s", token, client.Token)
	}

	if client.Region != customRegion {
		t.Errorf("expected region %s, got %s", customRegion, client.Region)
	}
}
