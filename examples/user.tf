/*
resource "centrifyvault_user" "testuser1" {
    username = "testuser1@example.com"
    email = "testuser1@centrify.com"
    displayname = "Test User 1"
    description = "Test user 1"
    password = "xxxxxxxxx"
    //password_never_expire = true
    //force_password_change_next = true
    //oauth_client = true
    send_email_invite = true
    office_number = "+65 98273622"
    home_number = "+65 98273622"
    mobile_number = "+65 98273622"
    //redirect_mfa_user_id = data.centrifyvault_user.admin.id
    //manager_username = "admin@example.com"
    roles = [
        data.centrifyvault_role.lab_infra_admin.id,
        data.centrifyvault_role.sso_role.id
    ]
}

resource "centrifyvault_user" "testuser2" {
    username = "testuser2@example.com"
    email = "testuser2@example.com"
    displayname = "Test User 2"
    description = "Test user 2"
    password = "xxxxxxxxx"
    //password_never_expire = true
    //force_password_change_next = true
    //oauth_client = true
    send_email_invite = true
    office_number = "+65 98273622"
    home_number = "+65 98273622"
    mobile_number = "+65 98273622"
    //redirect_mfa_user_id = centrifyvault_user.testuser1.id
    //manager_username = "admin@example.com"
}
*/