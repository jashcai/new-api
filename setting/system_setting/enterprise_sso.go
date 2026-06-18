package system_setting

import "github.com/QuantumNous/new-api/setting/config"

type EnterpriseSSOSettings struct {
	Enabled            bool              `json:"enabled"`
	GroupClaim         string            `json:"group_claim"`
	GroupMappings      map[string]string `json:"group_mappings"`
	DefaultGroup       string            `json:"default_group"`
	DirectGroupMatch   bool              `json:"direct_group_match"`
	SyncOnLogin        bool              `json:"sync_on_login"`
	SyncProfileOnLogin bool              `json:"sync_profile_on_login"`
}

var enterpriseSSOSettings = EnterpriseSSOSettings{
	Enabled:            false,
	GroupClaim:         "groups",
	GroupMappings:      map[string]string{},
	DefaultGroup:       "",
	DirectGroupMatch:   false,
	SyncOnLogin:        true,
	SyncProfileOnLogin: false,
}

func init() {
	config.GlobalConfig.Register("enterprise_sso", &enterpriseSSOSettings)
}

func GetEnterpriseSSOSettings() *EnterpriseSSOSettings {
	if enterpriseSSOSettings.GroupClaim == "" {
		enterpriseSSOSettings.GroupClaim = "groups"
	}
	if enterpriseSSOSettings.GroupMappings == nil {
		enterpriseSSOSettings.GroupMappings = map[string]string{}
	}
	return &enterpriseSSOSettings
}
