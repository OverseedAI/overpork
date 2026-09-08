package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

type capturedWriteRequest struct {
	path string
	body map[string]any
}

func newWriteRequestClient(t *testing.T) (*Client, *capturedWriteRequest) {
	t.Helper()

	captured := &capturedWriteRequest{}
	client := &Client{
		apiKey:    "public-key",
		secretKey: "secret-key",
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			captured.path = req.URL.Path
			if err := json.NewDecoder(req.Body).Decode(&captured.body); err != nil {
				t.Fatalf("request body is not JSON: %v", err)
			}
			return jsonHTTPResponse(`{"status":"SUCCESS"}`), nil
		})},
	}
	return client, captured
}

func TestGlueWriteRequests(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		request func(*Client, []string) error
	}{
		{
			name: "create",
			path: "/api/json/v3/domain/createGlue/example.com/ns1",
			request: func(client *Client, ips []string) error {
				return client.GlueCreate("example.com", "ns1", ips)
			},
		},
		{
			name: "update",
			path: "/api/json/v3/domain/updateGlue/example.com/ns1",
			request: func(client *Client, ips []string) error {
				return client.GlueUpdate("example.com", "ns1", ips)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, captured := newWriteRequestClient(t)
			ips := []string{"192.0.2.1", "2001:db8::1"}

			if err := tt.request(client, ips); err != nil {
				t.Fatalf("Glue%s() returned error: %v", tt.name, err)
			}
			if captured.path != tt.path {
				t.Fatalf("Glue%s() URL path = %q, want %q", tt.name, captured.path, tt.path)
			}

			wantBody := map[string]any{
				"apikey":       "public-key",
				"secretapikey": "secret-key",
				"ips":          []any{"192.0.2.1", "2001:db8::1"},
			}
			if !reflect.DeepEqual(captured.body, wantBody) {
				t.Fatalf("Glue%s() request body = %#v, want %#v", tt.name, captured.body, wantBody)
			}
		})
	}
}

func TestDNSSECCreateWriteRequestWithOptionalFields(t *testing.T) {
	client, captured := newWriteRequestClient(t)
	record := DNSSECRecord{
		KeyTag:     "12345",
		Algorithm:  "13",
		DigestType: "2",
		Digest:     "ABCDEF",
		PublicKey:  "public-key-data",
		Flags:      "257",
	}

	if err := client.DNSSECCreate("example.com", record); err != nil {
		t.Fatalf("DNSSECCreate() returned error: %v", err)
	}
	if got, want := captured.path, "/api/json/v3/dns/createDnssecRecord/example.com"; got != want {
		t.Fatalf("DNSSECCreate() URL path = %q, want %q", got, want)
	}

	wantBody := map[string]any{
		"apikey":        "public-key",
		"secretapikey":  "secret-key",
		"keyTag":        "12345",
		"alg":           "13",
		"digestType":    "2",
		"digest":        "ABCDEF",
		"keyDataPubKey": "public-key-data",
		"keyDataFlags":  "257",
	}
	if !reflect.DeepEqual(captured.body, wantBody) {
		t.Fatalf("DNSSECCreate() request body = %#v, want %#v", captured.body, wantBody)
	}
}

func TestDNSSECCreateWriteRequestWithoutOptionalFields(t *testing.T) {
	client, captured := newWriteRequestClient(t)
	record := DNSSECRecord{
		KeyTag:     "12345",
		Algorithm:  "13",
		DigestType: "2",
		Digest:     "ABCDEF",
	}

	if err := client.DNSSECCreate("example.com", record); err != nil {
		t.Fatalf("DNSSECCreate() returned error: %v", err)
	}
	if got, want := captured.path, "/api/json/v3/dns/createDnssecRecord/example.com"; got != want {
		t.Fatalf("DNSSECCreate() URL path = %q, want %q", got, want)
	}

	wantBody := map[string]any{
		"apikey":       "public-key",
		"secretapikey": "secret-key",
		"keyTag":       "12345",
		"alg":          "13",
		"digestType":   "2",
		"digest":       "ABCDEF",
	}
	if !reflect.DeepEqual(captured.body, wantBody) {
		t.Fatalf("DNSSECCreate() request body = %#v, want %#v", captured.body, wantBody)
	}
}
