package service

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func withEnterprisePolicySettings(t *testing.T, fn func(settings *system_setting.EnterprisePolicySettings)) {
	t.Helper()
	settings := system_setting.GetEnterprisePolicySettings()
	old := *settings
	t.Cleanup(func() {
		*settings = old
	})
	fn(settings)
}

func TestEnterpriseModelAccessUsesGroupAllowlist(t *testing.T) {
	withEnterprisePolicySettings(t, func(settings *system_setting.EnterprisePolicySettings) {
		settings.Enabled = true
		settings.EnforceGroupModelAllowlist = true
		settings.GroupModelAllowlist = map[string][]string{
			"finance": []string{"gpt-4o-mini", "text-embedding-3-small"},
		}
		settings.AllowUnconfiguredGroups = false
	})

	if !IsEnterpriseModelAllowed("finance", "", "gpt-4o-mini") {
		t.Fatal("expected finance group to access gpt-4o-mini")
	}
	if IsEnterpriseModelAllowed("finance", "", "gpt-4o") {
		t.Fatal("expected finance group to be denied gpt-4o")
	}
}

func TestEnterpriseModelAccessAllowsUnconfiguredGroupsWhenConfigured(t *testing.T) {
	withEnterprisePolicySettings(t, func(settings *system_setting.EnterprisePolicySettings) {
		settings.Enabled = true
		settings.EnforceGroupModelAllowlist = true
		settings.GroupModelAllowlist = map[string][]string{
			"finance": []string{"gpt-4o-mini"},
		}
		settings.AllowUnconfiguredGroups = true
	})

	if !IsEnterpriseModelAllowed("engineering", "", "gpt-4o") {
		t.Fatal("expected unconfigured group to be allowed")
	}
}

func TestApplyEnterpriseTokenPolicySetsDefaults(t *testing.T) {
	withEnterprisePolicySettings(t, func(settings *system_setting.EnterprisePolicySettings) {
		settings.Enabled = true
		settings.EnforceGroupModelAllowlist = true
		settings.GroupModelAllowlist = map[string][]string{
			"engineering": []string{"gpt-4o-mini", "text-embedding-3-small"},
		}
		settings.AllowUnconfiguredGroups = false
		settings.ForceTokenModelLimits = true
		settings.DefaultTokenQuota = 1000
		settings.MaxTokenQuota = 2000
		settings.RequireTokenExpiry = true
		settings.DefaultTokenValidDays = 30
		settings.DisableTokenCrossGroupRetry = true
	})

	token := &model.Token{
		ExpiredTime:     -1,
		CrossGroupRetry: true,
	}
	if err := ApplyEnterpriseTokenPolicy(token, "engineering"); err != nil {
		t.Fatalf("expected token policy to apply, got %v", err)
	}
	if token.RemainQuota != 1000 {
		t.Fatalf("expected default token quota 1000, got %d", token.RemainQuota)
	}
	if token.ExpiredTime <= 0 {
		t.Fatalf("expected default expiration, got %d", token.ExpiredTime)
	}
	if token.CrossGroupRetry {
		t.Fatal("expected cross-group retry to be disabled")
	}
	if !token.ModelLimitsEnabled {
		t.Fatal("expected model limits to be enabled")
	}
	if !strings.Contains(token.ModelLimits, "gpt-4o-mini") {
		t.Fatalf("expected default model limits from group allowlist, got %q", token.ModelLimits)
	}
}

func TestApplyEnterpriseTokenPolicyRejectsDisallowedTokenModel(t *testing.T) {
	withEnterprisePolicySettings(t, func(settings *system_setting.EnterprisePolicySettings) {
		settings.Enabled = true
		settings.EnforceGroupModelAllowlist = true
		settings.GroupModelAllowlist = map[string][]string{
			"finance": []string{"gpt-4o-mini"},
		}
		settings.AllowUnconfiguredGroups = false
	})

	token := &model.Token{
		ModelLimitsEnabled: true,
		ModelLimits:        "gpt-4o",
	}
	if err := ApplyEnterpriseTokenPolicy(token, "finance"); err == nil {
		t.Fatal("expected disallowed token model to be rejected")
	}
}
