package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func ValidateEnterpriseModelAccess(userGroup string, usingGroup string, modelName string) error {
	if IsEnterpriseModelAllowed(userGroup, usingGroup, modelName) {
		return nil
	}
	return fmt.Errorf("enterprise policy denies model %s for group %s", modelName, effectivePolicyGroup(userGroup, usingGroup))
}

func IsEnterpriseModelAllowed(userGroup string, usingGroup string, modelName string) bool {
	settings := system_setting.GetEnterprisePolicySettings()
	if !settings.Enabled || !settings.EnforceGroupModelAllowlist || strings.TrimSpace(modelName) == "" {
		return true
	}

	candidateGroups := enterprisePolicyCandidateGroups(userGroup, usingGroup)
	for _, group := range candidateGroups {
		allowed, configured := enterpriseGroupAllowsModel(group, modelName, settings)
		if allowed {
			return true
		}
		if !configured && settings.AllowUnconfiguredGroups {
			return true
		}
	}
	return false
}

func FilterEnterpriseModelsForGroups(userGroup string, groups []string, models []string) []string {
	settings := system_setting.GetEnterprisePolicySettings()
	if !settings.Enabled || !settings.EnforceGroupModelAllowlist {
		return models
	}
	if len(groups) == 0 {
		groups = enterprisePolicyCandidateGroups(userGroup, "")
	}

	filtered := make([]string, 0, len(models))
	for _, modelName := range models {
		allowed := false
		for _, group := range groups {
			if IsEnterpriseModelAllowed(userGroup, group, modelName) {
				allowed = true
				break
			}
		}
		if allowed {
			filtered = append(filtered, modelName)
		}
	}
	return filtered
}

func ApplyEnterpriseTokenPolicy(token *model.Token, userGroup string) error {
	if token == nil {
		return nil
	}
	settings := system_setting.GetEnterprisePolicySettings()
	if !settings.Enabled {
		return nil
	}

	if settings.DisableTokenCrossGroupRetry {
		token.CrossGroupRetry = false
	}

	if settings.DisableUnlimitedTokens && token.UnlimitedQuota {
		return fmt.Errorf("enterprise policy disables unlimited quota tokens")
	}

	if !token.UnlimitedQuota {
		if settings.DefaultTokenQuota > 0 && token.RemainQuota == 0 {
			token.RemainQuota = settings.DefaultTokenQuota
		}
		if settings.MaxTokenQuota > 0 && token.RemainQuota > settings.MaxTokenQuota {
			return fmt.Errorf("token quota exceeds enterprise maximum %d", settings.MaxTokenQuota)
		}
	}

	if settings.RequireTokenExpiry && token.ExpiredTime == -1 {
		if settings.DefaultTokenValidDays <= 0 {
			return fmt.Errorf("enterprise policy requires token expiration")
		}
		token.ExpiredTime = time.Now().AddDate(0, 0, settings.DefaultTokenValidDays).Unix()
	}

	effectiveGroup := effectivePolicyGroup(userGroup, token.Group)
	if settings.ForceTokenModelLimits && !token.ModelLimitsEnabled {
		modelLimits := effectiveDefaultTokenModels(effectiveGroup, settings)
		if len(modelLimits) > 0 {
			token.ModelLimitsEnabled = true
			token.ModelLimits = strings.Join(modelLimits, ",")
		}
	}

	if token.ModelLimitsEnabled && settings.EnforceGroupModelAllowlist {
		for _, modelName := range token.GetModelLimits() {
			modelName = strings.TrimSpace(modelName)
			if modelName == "" {
				continue
			}
			if !IsEnterpriseModelAllowed(userGroup, token.Group, modelName) {
				return fmt.Errorf("token model %s is not allowed for group %s", modelName, effectiveGroup)
			}
		}
	}
	return nil
}

func effectiveDefaultTokenModels(group string, settings *system_setting.EnterprisePolicySettings) []string {
	defaultModels := normalizeModelList(settings.DefaultTokenModels, true)
	if len(defaultModels) > 0 {
		return defaultModels
	}

	groupModels := settings.GroupModelAllowlist[group]
	return normalizeModelList(groupModels, true)
}

func enterprisePolicyCandidateGroups(userGroup string, usingGroup string) []string {
	usingGroup = strings.TrimSpace(usingGroup)
	userGroup = strings.TrimSpace(userGroup)
	if usingGroup == "auto" {
		autoGroups := GetUserAutoGroup(userGroup)
		if len(autoGroups) > 0 {
			return autoGroups
		}
	}
	if usingGroup != "" {
		return []string{usingGroup}
	}
	if userGroup != "" {
		return []string{userGroup}
	}
	return []string{"default"}
}

func effectivePolicyGroup(userGroup string, group string) string {
	group = strings.TrimSpace(group)
	if group != "" {
		return group
	}
	userGroup = strings.TrimSpace(userGroup)
	if userGroup != "" {
		return userGroup
	}
	return "default"
}

func enterpriseGroupAllowsModel(group string, modelName string, settings *system_setting.EnterprisePolicySettings) (bool, bool) {
	allowlist, configured := settings.GroupModelAllowlist[group]
	if !configured {
		return false, false
	}

	normalizedModel := ratio_setting.FormatMatchingModelName(strings.TrimSpace(modelName))
	for _, allowedModel := range allowlist {
		allowedModel = strings.TrimSpace(allowedModel)
		if allowedModel == "*" {
			return true, true
		}
		if ratio_setting.FormatMatchingModelName(allowedModel) == normalizedModel {
			return true, true
		}
	}
	return false, true
}

func normalizeModelList(models []string, omitWildcard bool) []string {
	normalized := make([]string, 0, len(models))
	seen := map[string]bool{}
	for _, modelName := range models {
		modelName = strings.TrimSpace(modelName)
		if modelName == "" {
			continue
		}
		if omitWildcard && modelName == "*" {
			continue
		}
		modelName = ratio_setting.FormatMatchingModelName(modelName)
		if seen[modelName] {
			continue
		}
		seen[modelName] = true
		normalized = append(normalized, modelName)
	}
	return normalized
}
