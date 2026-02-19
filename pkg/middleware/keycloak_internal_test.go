package middleware

import "testing"

func TestAudienceOrAzpAllowed_AllowsAudString(t *testing.T) {
	raw := map[string]any{"aud": "erp-frontend"}
	if !audienceOrAzpAllowed(raw, []string{"erp-frontend"}) {
		t.Fatalf("expected aud to be allowed")
	}
}

func TestAudienceOrAzpAllowed_AllowsAudArray(t *testing.T) {
	raw := map[string]any{"aud": []any{"account", "erp-frontend"}}
	if !audienceOrAzpAllowed(raw, []string{"erp-frontend"}) {
		t.Fatalf("expected aud array to be allowed")
	}
}

func TestAudienceOrAzpAllowed_AllowsAzp(t *testing.T) {
	raw := map[string]any{"azp": "erp-frontend", "aud": "account"}
	if !audienceOrAzpAllowed(raw, []string{"erp-frontend"}) {
		t.Fatalf("expected azp to be allowed")
	}
}

func TestAudienceOrAzpAllowed_RejectsWhenNoMatch(t *testing.T) {
	raw := map[string]any{"azp": "other", "aud": "account"}
	if audienceOrAzpAllowed(raw, []string{"erp-frontend"}) {
		t.Fatalf("expected to be rejected")
	}
}
