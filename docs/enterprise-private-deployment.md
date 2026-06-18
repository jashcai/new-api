# Enterprise Private Deployment Baseline

This project can run with a conservative private-deployment baseline by setting:

```env
PRIVATE_DEPLOYMENT_MODE=true
```

When enabled, runtime option loading still reads database-backed settings, but public-site features are forced off after every option sync:

- public registration and password registration
- built-in public OAuth providers such as GitHub, Discord, LinuxDO, WeChat, and Telegram
- invitation/new-user reward quotas
- payment compliance, online top-up methods, redemption entrypoints, payment discounts, and top-up links
- demo-site mode and check-in rewards

OIDC and custom OAuth providers remain available so operators can configure enterprise SSO.

## Enterprise SSO Group Sync

Enterprise SSO group sync maps IdP claims to the existing New API user `group` field during OIDC or custom OAuth login.

Example options:

```text
enterprise_sso.enabled=true
enterprise_sso.group_claim=groups
enterprise_sso.group_mappings={"Engineering":"eng","Finance":"finance"}
enterprise_sso.default_group=default
enterprise_sso.direct_group_match=false
enterprise_sso.sync_on_login=true
enterprise_sso.sync_profile_on_login=false
```

Behavior:

- `group_claim` is read from the OAuth/OIDC userinfo response. Arrays and comma-separated strings are supported.
- `group_mappings` maps IdP group names to local user groups.
- `default_group` is used when no claim value matches a mapping.
- `direct_group_match=true` allows claim values to be used as local group names without explicit mapping.
- `sync_on_login=true` keeps existing users aligned with the IdP on each SSO login.
- `sync_profile_on_login=true` also updates display name and email from the IdP.

## Enterprise Policy

Enterprise policy uses the existing user `group` as the department/project boundary. It can restrict visible and callable models by group, and can enforce safe defaults for user-created API tokens.

Example options:

```text
enterprise_policy.enabled=true
enterprise_policy.enforce_group_model_allowlist=true
enterprise_policy.group_model_allowlist={"engineering":["gpt-4o-mini","text-embedding-3-small"],"finance":["gpt-4o-mini"],"security":["*"]}
enterprise_policy.allow_unconfigured_groups=false
enterprise_policy.force_token_model_limits=true
enterprise_policy.default_token_models=["gpt-4o-mini"]
enterprise_policy.disable_unlimited_tokens=true
enterprise_policy.default_token_quota=100000
enterprise_policy.max_token_quota=1000000
enterprise_policy.require_token_expiry=true
enterprise_policy.default_token_valid_days=90
enterprise_policy.disable_token_cross_group_retry=true
```

Behavior:

- `group_model_allowlist` filters `/v1/models` responses and is also enforced in the relay distributor, so direct API calls cannot bypass the console.
- `*` grants all models for a group. Other entries should use exact model names after the project's normal model-name matching.
- `allow_unconfigured_groups=false` makes the allowlist fail closed for groups that are not listed.
- `force_token_model_limits=true` automatically enables per-token model limits from `default_token_models`; if no defaults are configured, it falls back to the token group's allowlist.
- `disable_unlimited_tokens=true`, `max_token_quota`, and token expiry settings prevent long-lived or unbounded API keys in private deployments.
- Token create/update/delete and token-key view operations are recorded as security audit logs without storing secret key values in logs.

## CORS

Private mode disables wildcard cross-origin access. If no origins are configured, browser access is same-origin only.

For split frontend/backend deployments, configure exact origins:

```env
CORS_ALLOWED_ORIGINS=https://ai-gateway.example.com,https://console.example.com
```

`FRONTEND_BASE_URL` is also accepted as an allowed origin. Path and query parts are ignored; only scheme, host, and port are used.

To override default allowed request headers:

```env
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,New-Api-User,X-Api-Key,X-Goog-Api-Key,Anthropic-Version,Mj-Api-Secret,OpenAI-Organization,OpenAI-Beta,X-Oneapi-Request-Id
```

## Session Cookies

For HTTPS deployments, set:

```env
SESSION_COOKIE_SECURE=true
```

This marks the session cookie as Secure while preserving the existing SameSite=Strict and HttpOnly settings.
