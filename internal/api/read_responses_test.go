package api

import (
	"net/http"
	"reflect"
	"testing"
)

func TestDomainGetDecodesNestedDomain(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got, want := req.URL.String(), BaseURL+"/domain/get/example.com"; got != want {
				t.Fatalf("DomainGet() URL = %q, want %q", got, want)
			}

			return jsonHTTPResponse(`{
				"status":"SUCCESS",
				"domain":{
					"domain":"example.com",
					"status":"ACTIVE",
					"tld":"com",
					"createDate":"2024-01-02 03:04:05",
					"expireDate":"2027-01-02 03:04:05",
					"securityLock":"1",
					"whoisPrivacy":"1",
					"autoRenew":"0",
					"notLocal":0
				}
			}`), nil
		})},
	}

	got, err := client.DomainGet("example.com")
	if err != nil {
		t.Fatalf("DomainGet() returned error: %v", err)
	}
	want := &Domain{
		Domain:       "example.com",
		Status:       "ACTIVE",
		TLD:          "com",
		CreateDate:   "2024-01-02 03:04:05",
		ExpireDate:   "2027-01-02 03:04:05",
		SecurityLock: "1",
		WhoisPrivacy: "1",
		AutoRenew:    "0",
		NotLocal:     0,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DomainGet() = %#v, want %#v", got, want)
	}
}

func TestDomainGetDecodesDocumentedNumericFlagsAsStrings(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return jsonHTTPResponse(`{
				"status":"SUCCESS",
				"domain":{
					"domain":"example.com",
					"status":"ACTIVE",
					"tld":"com",
					"createDate":"2024-01-02 03:04:05",
					"expireDate":"2027-01-02 03:04:05",
					"securityLock":1,
					"whoisPrivacy":0,
					"autoRenew":1,
					"notLocal":0
				}
			}`), nil
		})},
	}

	got, err := client.DomainGet("example.com")
	if err != nil {
		t.Fatalf("DomainGet() returned error: %v", err)
	}
	if got.SecurityLock != "1" {
		t.Errorf("DomainGet().SecurityLock = %q, want %q", got.SecurityLock, "1")
	}
	if got.WhoisPrivacy != "0" {
		t.Errorf("DomainGet().WhoisPrivacy = %q, want %q", got.WhoisPrivacy, "0")
	}
	if got.AutoRenew != "1" {
		t.Errorf("DomainGet().AutoRenew = %q, want %q", got.AutoRenew, "1")
	}
}

func TestGlueListDecodesHostTuples(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got, want := req.URL.String(), BaseURL+"/domain/getGlue/example.com"; got != want {
				t.Fatalf("GlueList() URL = %q, want %q", got, want)
			}

			return jsonHTTPResponse(`{
				"status":"SUCCESS",
				"hosts":[
					["ns1.example.com",{"v4":["192.0.2.1"],"v6":["2001:db8::1"]}],
					["nameservers.example.com",{"v4":["192.0.2.2","192.0.2.3"],"v6":[]}]
				]
			}`), nil
		})},
	}

	got, err := client.GlueList("example.com")
	if err != nil {
		t.Fatalf("GlueList() returned error: %v", err)
	}
	want := []GlueRecord{
		{Subdomain: "ns1", IPs: []string{"192.0.2.1", "2001:db8::1"}},
		{Subdomain: "nameservers", IPs: []string{"192.0.2.2", "192.0.2.3"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GlueList() = %#v, want %#v", got, want)
	}
}

func TestDNSSECListDecodesKeyedRecordsDeterministically(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got, want := req.URL.String(), BaseURL+"/dns/getDnssecRecords/example.com"; got != want {
				t.Fatalf("DNSSECList() URL = %q, want %q", got, want)
			}

			return jsonHTTPResponse(`{
				"status":"SUCCESS",
				"records":{
					"54321":{"alg":"8","digestType":"2","digest":"BBBB","pubKey":"key-b"},
					"12345":{"alg":"13","digestType":"2","digest":"AAAA","pubKey":"key-a"}
				}
			}`), nil
		})},
	}

	got, err := client.DNSSECList("example.com")
	if err != nil {
		t.Fatalf("DNSSECList() returned error: %v", err)
	}
	want := []DNSSECRecord{
		{KeyTag: "12345", Algorithm: "13", DigestType: "2", Digest: "AAAA", PublicKey: "key-a"},
		{KeyTag: "54321", Algorithm: "8", DigestType: "2", Digest: "BBBB", PublicKey: "key-b"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DNSSECList() = %#v, want %#v", got, want)
	}
}
