package api

import "testing"

func TestFilterDNSRecordsByType(t *testing.T) {
	records := []DNSRecord{
		{ID: "1", Type: "A", Name: "example.com"},
		{ID: "2", Type: "AAAA", Name: "example.com"},
		{ID: "3", Type: "A", Name: "www.example.com"},
	}

	got := filterDNSRecordsByType(records, "a")
	if len(got) != 2 || got[0].ID != "1" || got[1].ID != "3" {
		t.Fatalf("filterDNSRecordsByType() = %#v, want A records 1 and 3", got)
	}

	got = filterDNSRecordsByType(records, "MX")
	if got == nil || len(got) != 0 {
		t.Fatalf("filterDNSRecordsByType() = %#v, want a non-nil empty slice", got)
	}
}

func TestDNSByNameTypeEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		domain     string
		recordType string
		subdomain  string
		want       string
	}{
		{
			name:       "edit with subdomain",
			action:     "editByNameType",
			domain:     "example.com",
			recordType: "A",
			subdomain:  "www",
			want:       "/dns/editByNameType/example.com/A/www",
		},
		{
			name:       "edit root domain omits subdomain segment",
			action:     "editByNameType",
			domain:     "example.com",
			recordType: "A",
			subdomain:  "",
			want:       "/dns/editByNameType/example.com/A",
		},
		{
			name:       "delete with subdomain",
			action:     "deleteByNameType",
			domain:     "example.com",
			recordType: "MX",
			subdomain:  "mail",
			want:       "/dns/deleteByNameType/example.com/MX/mail",
		},
		{
			name:       "delete root domain omits subdomain segment",
			action:     "deleteByNameType",
			domain:     "example.com",
			recordType: "MX",
			subdomain:  "",
			want:       "/dns/deleteByNameType/example.com/MX",
		},
		{
			name:       "retrieve with subdomain",
			action:     "retrieveByNameType",
			domain:     "example.com",
			recordType: "TXT",
			subdomain:  "_dmarc",
			want:       "/dns/retrieveByNameType/example.com/TXT/_dmarc",
		},
		{
			name:       "retrieve root domain omits subdomain segment",
			action:     "retrieveByNameType",
			domain:     "example.com",
			recordType: "TXT",
			subdomain:  "",
			want:       "/dns/retrieveByNameType/example.com/TXT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dnsByNameTypeEndpoint(tt.action, tt.domain, tt.recordType, tt.subdomain)
			if got != tt.want {
				t.Errorf("dnsByNameTypeEndpoint() = %q, want %q", got, tt.want)
			}
		})
	}
}
