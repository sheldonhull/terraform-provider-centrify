
resource "centrifyvault_vaultsecret" "pas_admin_credential" {
    secret_name = "Centrify PAS Admin Credential"
    description = "admin@centrify.lab"
    secret_text = var.PASSWORD
    type = "Text"
}
