package system_setting

import "github.com/QuantumNous/new-api/setting/config"

type EnterprisePolicySettings struct {
	Enabled                     bool                `json:"enabled"`
	EnforceGroupModelAllowlist  bool                `json:"enforce_group_model_allowlist"`
	GroupModelAllowlist         map[string][]string `json:"group_model_allowlist"`
	AllowUnconfiguredGroups     bool                `json:"allow_unconfigured_groups"`
	ForceTokenModelLimits       bool                `json:"force_token_model_limits"`
	DefaultTokenModels          []string            `json:"default_token_models"`
	DisableUnlimitedTokens      bool                `json:"disable_unlimited_tokens"`
	DefaultTokenQuota           int                 `json:"default_token_quota"`
	MaxTokenQuota               int                 `json:"max_token_quota"`
	RequireTokenExpiry          bool                `json:"require_token_expiry"`
	DefaultTokenValidDays       int                 `json:"default_token_valid_days"`
	DisableTokenCrossGroupRetry bool                `json:"disable_token_cross_group_retry"`
}

var enterprisePolicySettings = EnterprisePolicySettings{
	Enabled:                     false,
	EnforceGroupModelAllowlist:  false,
	GroupModelAllowlist:         map[string][]string{},
	AllowUnconfiguredGroups:     true,
	ForceTokenModelLimits:       false,
	DefaultTokenModels:          []string{},
	DisableUnlimitedTokens:      false,
	DefaultTokenQuota:           0,
	MaxTokenQuota:               0,
	RequireTokenExpiry:          false,
	DefaultTokenValidDays:       0,
	DisableTokenCrossGroupRetry: false,
}

func init() {
	config.GlobalConfig.Register("enterprise_policy", &enterprisePolicySettings)
}

func GetEnterprisePolicySettings() *EnterprisePolicySettings {
	if enterprisePolicySettings.GroupModelAllowlist == nil {
		enterprisePolicySettings.GroupModelAllowlist = map[string][]string{}
	}
	if enterprisePolicySettings.DefaultTokenModels == nil {
		enterprisePolicySettings.DefaultTokenModels = []string{}
	}
	return &enterprisePolicySettings
}
