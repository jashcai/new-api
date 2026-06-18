package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func TestApplyPrivateDeploymentModeOverridesPublicSiteOptions(t *testing.T) {
	oldPrivateMode := common.PrivateDeploymentMode
	oldOptionMap := common.OptionMap
	oldRegisterEnabled := common.RegisterEnabled
	oldPasswordRegisterEnabled := common.PasswordRegisterEnabled
	oldGitHubOAuthEnabled := common.GitHubOAuthEnabled
	oldLinuxDOOAuthEnabled := common.LinuxDOOAuthEnabled
	oldWeChatAuthEnabled := common.WeChatAuthEnabled
	oldTelegramOAuthEnabled := common.TelegramOAuthEnabled
	oldDiscordEnabled := system_setting.GetDiscordSettings().Enabled
	oldQuotaForNewUser := common.QuotaForNewUser
	oldQuotaForInviter := common.QuotaForInviter
	oldQuotaForInvitee := common.QuotaForInvitee
	oldTopUpLink := common.TopUpLink
	oldDemoSiteEnabled := operation_setting.DemoSiteEnabled
	oldComplianceConfirmed := operation_setting.GetPaymentSetting().ComplianceConfirmed
	oldComplianceTermsVersion := operation_setting.GetPaymentSetting().ComplianceTermsVersion
	oldPayMethods := operation_setting.PayMethods
	oldAmountDiscount := operation_setting.GetPaymentSetting().AmountDiscount
	oldCheckinEnabled := operation_setting.GetCheckinSetting().Enabled
	t.Cleanup(func() {
		common.PrivateDeploymentMode = oldPrivateMode
		common.OptionMap = oldOptionMap
		common.RegisterEnabled = oldRegisterEnabled
		common.PasswordRegisterEnabled = oldPasswordRegisterEnabled
		common.GitHubOAuthEnabled = oldGitHubOAuthEnabled
		common.LinuxDOOAuthEnabled = oldLinuxDOOAuthEnabled
		common.WeChatAuthEnabled = oldWeChatAuthEnabled
		common.TelegramOAuthEnabled = oldTelegramOAuthEnabled
		system_setting.GetDiscordSettings().Enabled = oldDiscordEnabled
		common.QuotaForNewUser = oldQuotaForNewUser
		common.QuotaForInviter = oldQuotaForInviter
		common.QuotaForInvitee = oldQuotaForInvitee
		common.TopUpLink = oldTopUpLink
		operation_setting.DemoSiteEnabled = oldDemoSiteEnabled
		operation_setting.GetPaymentSetting().ComplianceConfirmed = oldComplianceConfirmed
		operation_setting.GetPaymentSetting().ComplianceTermsVersion = oldComplianceTermsVersion
		operation_setting.PayMethods = oldPayMethods
		operation_setting.GetPaymentSetting().AmountDiscount = oldAmountDiscount
		operation_setting.GetCheckinSetting().Enabled = oldCheckinEnabled
	})

	common.PrivateDeploymentMode = true
	common.RegisterEnabled = true
	common.PasswordRegisterEnabled = true
	common.GitHubOAuthEnabled = true
	common.LinuxDOOAuthEnabled = true
	common.WeChatAuthEnabled = true
	common.TelegramOAuthEnabled = true
	system_setting.GetDiscordSettings().Enabled = true
	common.QuotaForNewUser = 100
	common.QuotaForInviter = 200
	common.QuotaForInvitee = 300
	common.TopUpLink = "https://pay.example.com"
	operation_setting.DemoSiteEnabled = true
	operation_setting.GetPaymentSetting().ComplianceConfirmed = true
	operation_setting.GetPaymentSetting().ComplianceTermsVersion = operation_setting.CurrentComplianceTermsVersion
	operation_setting.PayMethods = []map[string]string{{"type": "alipay"}}
	operation_setting.GetPaymentSetting().AmountDiscount = map[int]float64{100: 0.9}
	operation_setting.GetCheckinSetting().Enabled = true
	common.OptionMap = map[string]string{
		"RegisterEnabled":                         "true",
		"PasswordRegisterEnabled":                 "true",
		"GitHubOAuthEnabled":                      "true",
		"LinuxDOOAuthEnabled":                     "true",
		"WeChatAuthEnabled":                       "true",
		"TelegramOAuthEnabled":                    "true",
		"discord.enabled":                         "true",
		"QuotaForNewUser":                         "100",
		"QuotaForInviter":                         "200",
		"QuotaForInvitee":                         "300",
		"TopUpLink":                               "https://pay.example.com",
		"PayMethods":                              `[{"type":"alipay"}]`,
		"DemoSiteEnabled":                         "true",
		"payment_setting.compliance_confirmed":    "true",
		"payment_setting.amount_discount":         `{"100":0.9}`,
		"checkin_setting.enabled":                 "true",
		"payment_setting.compliance_confirmed_at": "123",
	}

	applyPrivateDeploymentModeOverrides()

	if common.RegisterEnabled || common.PasswordRegisterEnabled {
		t.Fatal("expected public registration to be disabled")
	}
	if common.GitHubOAuthEnabled || common.LinuxDOOAuthEnabled || common.WeChatAuthEnabled || common.TelegramOAuthEnabled || system_setting.GetDiscordSettings().Enabled {
		t.Fatal("expected public OAuth providers to be disabled")
	}
	if common.QuotaForNewUser != 0 || common.QuotaForInviter != 0 || common.QuotaForInvitee != 0 {
		t.Fatalf("expected invitation and new-user quotas to be zero, got %d/%d/%d", common.QuotaForNewUser, common.QuotaForInviter, common.QuotaForInvitee)
	}
	if common.TopUpLink != "" || operation_setting.DemoSiteEnabled {
		t.Fatal("expected public top-up link and demo mode to be disabled")
	}
	if operation_setting.GetPaymentSetting().ComplianceConfirmed {
		t.Fatal("expected payment compliance to be forced unconfirmed")
	}
	if len(operation_setting.PayMethods) != 0 {
		t.Fatalf("expected no pay methods, got %#v", operation_setting.PayMethods)
	}
	if len(operation_setting.GetPaymentSetting().AmountDiscount) != 0 {
		t.Fatalf("expected no payment discounts, got %#v", operation_setting.GetPaymentSetting().AmountDiscount)
	}
	if operation_setting.GetCheckinSetting().Enabled {
		t.Fatal("expected check-in to be disabled")
	}
	if common.OptionMap["RegisterEnabled"] != "false" ||
		common.OptionMap["PayMethods"] != "[]" ||
		common.OptionMap["discord.enabled"] != "false" ||
		common.OptionMap["payment_setting.compliance_confirmed"] != "false" ||
		common.OptionMap["checkin_setting.enabled"] != "false" {
		t.Fatalf("expected OptionMap overrides, got %#v", common.OptionMap)
	}
}
