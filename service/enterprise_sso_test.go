package service

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/oauth"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func TestResolveEnterpriseSSOGroupMappedArrayClaim(t *testing.T) {
	settings := system_setting.GetEnterpriseSSOSettings()
	old := *settings
	t.Cleanup(func() {
		*settings = old
	})

	settings.Enabled = true
	settings.GroupClaim = "groups"
	settings.GroupMappings = map[string]string{
		"Engineering": "eng",
		"Finance":     "fin",
	}
	settings.DefaultGroup = "default"
	settings.DirectGroupMatch = false

	group, ok := ResolveEnterpriseSSOGroup(&oauth.OAuthUser{
		Extra: map[string]any{
			"claims_json": `{"groups":["Finance","Engineering"]}`,
		},
	})
	if !ok || group != "fin" {
		t.Fatalf("expected first mapped group fin, got group=%q ok=%v", group, ok)
	}
}

func TestResolveEnterpriseSSOGroupDefaultFallback(t *testing.T) {
	settings := system_setting.GetEnterpriseSSOSettings()
	old := *settings
	t.Cleanup(func() {
		*settings = old
	})

	settings.Enabled = true
	settings.GroupClaim = "department"
	settings.GroupMappings = map[string]string{}
	settings.DefaultGroup = "default"
	settings.DirectGroupMatch = false

	group, ok := ResolveEnterpriseSSOGroup(&oauth.OAuthUser{
		Extra: map[string]any{
			"claims_json": `{"department":"Unknown"}`,
		},
	})
	if !ok || group != "default" {
		t.Fatalf("expected default fallback group, got group=%q ok=%v", group, ok)
	}
}

func TestApplyEnterpriseSSOToNewUser(t *testing.T) {
	settings := system_setting.GetEnterpriseSSOSettings()
	old := *settings
	t.Cleanup(func() {
		*settings = old
	})

	settings.Enabled = true
	settings.GroupClaim = "department"
	settings.GroupMappings = map[string]string{
		"Platform": "platform",
	}
	settings.DefaultGroup = ""
	settings.DirectGroupMatch = false

	user := &model.User{Group: "default"}
	changed := ApplyEnterpriseSSOToNewUser(user, &oauth.OAuthUser{
		Extra: map[string]any{
			"claims_json": `{"department":"Platform"}`,
		},
	})
	if !changed || user.Group != "platform" {
		t.Fatalf("expected user group to be platform, got group=%q changed=%v", user.Group, changed)
	}
}
