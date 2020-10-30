
resource "centrifyvault_role" "infra_admin" {
    name = "${var.PREFIX}Infra Admins"
    description = "Requester who can request access to all lab systems."
    adminrights = [
        "Privileged Access Service User",
    ]
    member {
        id = data.centrifyvault_directoryobject.grp_infra_admins.id
        type = "Group"
    }
}


resource "centrifyvault_role" "infra_owner" {
    name = "${var.PREFIX}Infra Owners"
    description = "Approver who can approve access to access lab systems."
    adminrights = [
        "Privileged Access Service User",
    ]
    member {
        id = data.centrifyvault_directoryobject.grp_infra_owners.id
        type = "Group"
    }
}

/*
resource "centrifyvault_role" "mfa_machines_users" {
    name = "${var.PREFIX}MFA Machines & Users"
    description = "Machines and users who are enforced MFA for direct access without going through PAS."
    adminrights = [
        "Computer Login and Privilege Elevation",
    ]
}

resource "centrifyvault_role" "cloud_normal_user" {
    name = "${var.PREFIX}Cloud Normal Users"
    description = "AD accounts who can login to non-domain joined machines but without any privileges."
    adminrights = [
        "Computer Login and Privilege Elevation",
    ]
}

resource "centrifyvault_role" "cloud_local_admin" {
    name = "${var.PREFIX}Cloud Local Admins"
    description = "AD accounts that are granted local administrator access to non-domain joined machines."
    adminrights = [
        "Computer Login and Privilege Elevation",
    ]
}
*/