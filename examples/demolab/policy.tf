
// Order of policies
resource "centrifyvault_policyorder" "policy_order" {
    policy_order = [
        centrifyvault_policy.critical_system_set.id,
        centrifyvault_policy.breakglass_account_set.id,
        centrifyvault_policy.user_login.id,
        centrifyvault_policy.sysadmin_login.id,
        data.centrifyvault_policy.default_policy.id,
        centrifyvault_policy.deny_login.id,
    ]
}

// User Login Policy
resource "centrifyvault_policy" "user_login" {
    name = "${var.PREFIX}User Login Policy"
    description = "Login policy for lab users. It applies to users who are assigned with PAS role such as requester and approver."
    link_type = "Role"
    policy_assignment = [
        centrifyvault_role.infra_admin.id,
        centrifyvault_role.infra_owner.id
    ]
    settings {
        centrify_services {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.default_other_login.id
            allow_session_persist = true
            challenge_rule {
                authentication_profile_id = centrifyvault_authenticationprofile.twofa_authprofile.id
                rule {
                    filter = "Browser"
                    condition = "OpNotEqual"
                    value = "Chrome"
                }
            }
        }
        oath_otp {
            allow_otp = true
        }
        user_account {
            allow_user_change_password = true
            password_change_auth_profile_id = centrifyvault_authenticationprofile.stepup_authprofile.id
            show_fido2 = true
            fido2_prompt = "FIDO2 Security Key"
            fido2_auth_profile_id = centrifyvault_authenticationprofile.stepup_authprofile.id
            show_otp = true
            otp_prompt = "Google Authenticator"
            otp_auth_profile_id = centrifyvault_authenticationprofile.stepup_authprofile.id
            allow_mfa_redirect_change = true
        }
    }
}

// System Administrator Login Policy
resource "centrifyvault_policy" "sysadmin_login" {
    name = "${var.PREFIX}System Administrator Login Policy"
    description = "Login policy for default PAS system administrator. Only password authentication is required and no 2FA for easy usage."
    link_type = "Role"
    policy_assignment = [
        data.centrifyvault_role.sysadmin.id
    ]
    settings {
        centrify_services {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.default_other_login.id
            allow_session_persist = true
        }
    }
}

// Deny Login Policy
resource "centrifyvault_policy" "deny_login" {
    name = "${var.PREFIX}Deny Login Policy"
    description = "Catch all policy that denies users who don't have PAS role from logging in. This policy must be placed at the bottom of policy list."
    link_type = "Role"
    policy_assignment = []
    settings {
        centrify_services {
            authentication_enabled = true
            default_profile_id = "-1"
            allow_iwa = false
            use_certauth = false
            block_auth_from_same_device = true
            continue_failed_sessions = false
        }
    }
}

// Critical System Set Policy
resource "centrifyvault_policy" "critical_system_set" {
    name = "${var.PREFIX}Critical System Set Policy"
    description = "Apply step-up authentication to critical system upon login."
    link_type = "Collection"
    policy_assignment = [
        format("Server|%s", centrifyvault_manualset.critical_systems.id),
    ]
    settings {
        system_set {
            default_profile_id = centrifyvault_authenticationprofile.stepup_authprofile.id
        }
    }
}

// Break Glass Account Set Policy
resource "centrifyvault_policy" "breakglass_account_set" {
    name = "${var.PREFIX}Break Glass Account Set Policy"
    description = "Apply step-up authentication to break glass account upon password checkout."
    link_type = "Collection"
    policy_assignment = [
        format("VaultAccount|%s", centrifyvault_manualset.breakglass_accounts.id),
    ]
    settings {
        account_set {
            default_profile_id = centrifyvault_authenticationprofile.stepup_authprofile.id
        }
    }
}
