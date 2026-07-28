package cmd

import "testing"

func TestValidateSSLPart(t *testing.T) {
	for _, part := range []string{"", "cert", "key", "intermediate", "public"} {
		t.Run("valid_"+part, func(t *testing.T) {
			if err := validateSSLPart(part); err != nil {
				t.Fatalf("validateSSLPart(%q) returned error: %v", part, err)
			}
		})
	}

	if err := validateSSLPart("publci"); err == nil {
		t.Fatal("validateSSLPart() accepted an invalid part")
	}
}
