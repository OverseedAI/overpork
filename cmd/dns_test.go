package cmd

import "testing"

func TestNormalizeDNSSubdomain(t *testing.T) {
	tests := map[string]string{
		"@":      "",
		"":       "",
		"www":    "www",
		"_dmarc": "_dmarc",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := normalizeDNSSubdomain(input); got != want {
				t.Fatalf("normalizeDNSSubdomain(%q) = %q, want %q", input, got, want)
			}
		})
	}
}
