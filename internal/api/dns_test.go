package api

import "testing"

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
