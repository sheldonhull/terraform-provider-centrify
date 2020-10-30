/*
data "centrifyvault_policy" "Invited_Users" {
    name = "Invited Users"
}

data "centrifyvault_policy" "User_Login_Policy" {
    name = "LAB User Login Policy"
}

data "centrifyvault_policy" "Machine_Login_Policy" {
    name = "LAB Machine Login Policy"
}

data "centrifyvault_policy" "System_Administrator_Login_Policy" {
    name = "LAB System Administrator Login Policy"
}

data "centrifyvault_policy" "Default_Policy" {
    name = "Default Policy"
}

data "centrifyvault_policy" "Deny_Login_Policy" {
    name = "LAB Deny Login Policy"
}

resource "centrifyvault_policyorder" "policy_order" {
    policy_order = [
        data.centrifyvault_policy.Invited_Users.id,
        data.centrifyvault_policy.User_Login_Policy.id,
        centrifyvault_policy.test_policy1.id,
        data.centrifyvault_policy.Machine_Login_Policy.id,
        data.centrifyvault_policy.System_Administrator_Login_Policy.id,
        data.centrifyvault_policy.Default_Policy.id,
        data.centrifyvault_policy.Deny_Login_Policy.id,
    ]
}



resource "centrifyvault_policy" "test_policy1" {
    name = "Test Policy 1"
    description = "Test Policy 1"
    link_type = "Role"
    policy_assignment = [
        // Role
        //data.centrifyvault_role.lab_infra_admin.id,
        //data.centrifyvault_role.sso_role.id,
        //data.centrifyvault_manualset.lab_systems.id,
        // System Set
        //format("Server|%s", data.centrifyvault_manualset.lab_systems.id),
        // Database Set
        //format("VaultDatabase|%s", data.centrifyvault_manualset.lab_databases.id),
        //"VaultDatabase|@SQL Server",
        // Domain Set
        //"VaultDomain|@All Domains",
        // Account Set
        //"VaultAccount|@Database Accounts",
        //format("VaultAccount|%s", data.centrifyvault_manualset.lab_all_accounts.id),
        // Secret Set
        //"DataVault|@Text Generic Secrets",
        // SSHKey Set
        //"SshKeys|@Managed SshKeys",
    ]
    
    settings {
        
        centrify_services {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            // Session Parameters
            session_lifespan = 23
            allow_session_persist = true
            default_session_persist = true
            persist_session_lifespan = 30
            // Other Settings
            allow_iwa = true
            iwa_set_cookie = true
            iwa_satisfies_all = true
            use_certauth = true
            certauth_skip_challenge = true
            certauth_set_cookie = true
            certauth_satisfies_all = true
            allow_no_mfa_mech = true
            auth_rule_federated = false
            federated_satisfies_all = true
            block_auth_from_same_device = true
            continue_failed_sessions = true
            stop_auth_on_prev_failed = true
            remember_last_factor = true
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "DayOfWeek"
                    condition = "OpIsDayOfWeek"
                    value = "L,1,3,4"
                }
                rule {
                    filter = "Browser"
                    condition = "OpNotEqual"
                    value = "Firefox"
                }
                rule {
                    filter = "CountryCode"
                    condition = "OpNotEqual"
                    value = "GA"
                }
            }
        }

        centrify_client {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }
        
        centrify_css_server {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            pass_through_mode = 2
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }

        centrify_css_workstation {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }

        centrify_css_elevation {
            authentication_enabled = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }

        self_service {
            // Password Reset
            account_selfservice_enabled = true
            password_reset_enabled = true
            pwreset_allow_for_aduser = true
            pwreset_with_cookie_only = true
            login_after_reset = true
            pwreset_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            max_reset_attempts = 5
            // Account Unlock
            account_unlock_enabled = true
            unlock_allow_for_aduser = true
            unlock_with_cookie_only = true
            show_locked_message = true
            unlock_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            // Active Directory Self Service Settings
            use_ad_admin = true
            ad_admin_user = "ad_admin"
            admin_user_password {
                type = "SafeString"
                value = "xxxxxxxxxxx"
            }
            // Additional Policy Parameters
            max_reset_allowed = 6
            max_time_allowed = 50
        }

        password_settings {
            // Password Requirements
            min_length = 12
            max_length = 24
            require_digit = true
            require_mix_case = true
            require_symbol = true
            // Display Requirements
            show_password_complexity = true
            complexity_hint = "Whatever ......."
            // Additional Requirements
            //no_of_repeated_char_allowed = 2
            check_weak_password = true
            allow_include_username = true
            allow_include_displayname = true
            require_unicode = true
            // Password Age
            min_age_in_days = 10
            max_age_in_days = 90
            password_history = 10
            expire_soft_notification = 35
            expire_hard_notification = 72
            expire_notification_mobile = true
            // Capture Settings
            bad_attempt_threshold = 5
            capture_window = 20
            lockout_duration = 30
        }

        oath_otp {
            allow_otp = true
        }

        radius {
            allow_radius = true
            require_challenges = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            send_vendor_attributes = true
            allow_external_radius = true
        }

        user_account {
            allow_user_change_password = true
            password_change_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            show_fido2 = true
            fido2_prompt = "FIDO2 Key"
            fido2_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            show_otp = true
            otp_prompt = "Google Authenticator"
            otp_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            configure_security_questions = true
            prevent_dup_answers = false
            user_defined_questions = 3
            admin_defined_questions = 2
            min_char_in_answer = 2
            question_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            allow_phone_pin_change = true
            min_phone_pin_length = 6
            phone_pin_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            allow_mfa_redirect_change = true
            user_profile_auth_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            default_language = "en"
        }

        mobile_device {
            allow_enrollment = true
            permit_non_compliant_device = true
            enable_invite_enrollment = true
            allow_notify_multi_devices = true
            enable_debug = true
            location_tracking = true
            force_fingerprint = true
            allow_fallback_pin = true
            //require_passcode = true
            auto_lock_timeout = 15
            lock_app_on_exit = true
        }

        system_set {
            // Account Policy
            checkout_lifetime = 60
            // System Policy
            allow_remote_access = true
            allow_rdp_clipboard = true
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            // Security Settings
            allow_multiple_checkouts = true
            enable_password_rotation = true
            password_rotate_interval = 80
            enable_password_rotation_after_checkin = true
            minimum_password_age = 30
            minimum_sshkey_age = 30
            enable_sshkey_rotation = true
            sshkey_rotate_interval = 90
            sshkey_algorithm = "RSA_2048"
            // Maintenance Settings
            enable_password_history_cleanup = true
            password_historycleanup_duration = 120
            enable_sshkey_history_cleanup = true
            sshkey_historycleanup_duration = 120
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }

        database_set {
            // Account Security
            checkout_lifetime = 60
            // Security Settings
            allow_multiple_checkouts = true
            enable_password_rotation = true
            password_rotate_interval = 90
            enable_password_rotation_after_checkin = true
            minimum_password_age = 70
            // Maintenance Settings
            enable_password_history_cleanup = true
            password_historycleanup_duration = 120
        }

        domain_set {
            // Domain Security
            checkout_lifetime = 60
            // Security Settings
            allow_multiple_checkouts = true
            enable_password_rotation = true
            password_rotate_interval = 91
            enable_password_rotation_after_checkin = true
            minimum_password_age = 70
            // Maintenance Settings
            enable_password_history_cleanup = true
            password_historycleanup_duration = 120
        }

        account_set {
            checkout_lifetime = 60
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }

        secret_set {
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }

        sshkey_set {
            default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }
        
    }
    
}
*/