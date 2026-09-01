package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonHTTPResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestPricingListUsesPublicAPIEndpoint(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got, want := req.URL.String(), BaseURL+"/pricing/get"; got != want {
				t.Fatalf("PricingList() URL = %q, want %q", got, want)
			}

			var body map[string]string
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Fatalf("PricingList() request body is not JSON: %v", err)
			}
			if len(body) != 0 {
				t.Fatalf("PricingList() request body = %#v, want empty object", body)
			}

			return jsonHTTPResponse(`{"status":"SUCCESS","pricing":{"com":{"registration":"9.73","renewal":"9.73","transfer":"9.73"}}}`), nil
		})},
	}

	pricing, err := client.PricingList()
	if err != nil {
		t.Fatalf("PricingList() returned error: %v", err)
	}
	if got := pricing["com"].Registration; got != "9.73" {
		t.Fatalf("PricingList() registration = %q, want %q", got, "9.73")
	}
}

func TestDomainCheckSendsCredentials(t *testing.T) {
	client := &Client{
		apiKey:    "public-key",
		secretKey: "secret-key",
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got, want := req.URL.String(), BaseURL+"/domain/checkDomain/example.com"; got != want {
				t.Fatalf("DomainCheck() URL = %q, want %q", got, want)
			}

			var body map[string]string
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Fatalf("DomainCheck() request body is not JSON: %v", err)
			}
			if body["apikey"] != "public-key" || body["secretapikey"] != "secret-key" {
				t.Fatalf("DomainCheck() request body omitted credentials: %#v", body)
			}

			return jsonHTTPResponse(`{"status":"SUCCESS","response":{"avail":"yes","price":"9.73"}}`), nil
		})},
	}

	available, price, err := client.DomainCheck("example.com")
	if err != nil {
		t.Fatalf("DomainCheck() returned error: %v", err)
	}
	if !available || price != 9.73 {
		t.Fatalf("DomainCheck() = (%v, %.2f), want (true, 9.73)", available, price)
	}
}
