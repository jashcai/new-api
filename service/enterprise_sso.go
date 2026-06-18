package service

import (
	"strings"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/oauth"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/tidwall/gjson"
)

const oauthClaimsJSONKey = "claims_json"

func ApplyEnterpriseSSOToNewUser(user *model.User, oauthUser *oauth.OAuthUser) bool {
	group, ok := ResolveEnterpriseSSOGroup(oauthUser)
	if !ok || group == "" {
		return false
	}
	user.Group = group
	return true
}

func SyncEnterpriseSSOUser(user *model.User, oauthUser *oauth.OAuthUser) (bool, error) {
	setting := system_setting.GetEnterpriseSSOSettings()
	if !setting.Enabled || !setting.SyncOnLogin || user == nil || user.Id == 0 {
		return false, nil
	}

	changed := false
	if group, ok := ResolveEnterpriseSSOGroup(oauthUser); ok && group != "" && group != user.Group {
		user.Group = group
		changed = true
	}
	if setting.SyncProfileOnLogin && oauthUser != nil {
		if oauthUser.DisplayName != "" && oauthUser.DisplayName != user.DisplayName {
			user.DisplayName = oauthUser.DisplayName
			changed = true
		}
		if oauthUser.Email != "" && oauthUser.Email != user.Email {
			user.Email = oauthUser.Email
			changed = true
		}
	}
	if !changed {
		return false, nil
	}
	return true, user.Update(false)
}

func ResolveEnterpriseSSOGroup(oauthUser *oauth.OAuthUser) (string, bool) {
	setting := system_setting.GetEnterpriseSSOSettings()
	if !setting.Enabled {
		return "", false
	}

	for _, candidate := range enterpriseSSOClaimValues(oauthUser, setting.GroupClaim) {
		if mapped := strings.TrimSpace(setting.GroupMappings[candidate]); mapped != "" {
			return mapped, true
		}
		if setting.DirectGroupMatch {
			return candidate, true
		}
	}

	if fallback := strings.TrimSpace(setting.DefaultGroup); fallback != "" {
		return fallback, true
	}
	return "", false
}

func enterpriseSSOClaimValues(oauthUser *oauth.OAuthUser, claimPath string) []string {
	if oauthUser == nil || oauthUser.Extra == nil {
		return nil
	}
	claimsJSON, _ := oauthUser.Extra[oauthClaimsJSONKey].(string)
	if claimsJSON == "" {
		return nil
	}
	if claimPath == "" {
		claimPath = "groups"
	}
	result := gjson.Get(claimsJSON, claimPath)
	if !result.Exists() {
		return nil
	}
	values := make([]string, 0)
	if result.IsArray() {
		for _, item := range result.Array() {
			values = appendClaimValue(values, item.String())
		}
		return values
	}
	return appendClaimValue(values, result.String())
}

func appendClaimValue(values []string, raw string) []string {
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}
