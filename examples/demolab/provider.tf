provider "centrifyvault" {
    url = "https://pas.centrify.lab"
    appid = "CentrifyCLI"
    scope = "all"
    username = "admin@centrify.lab"
    password = var.PASSWORD
    //logpath = "centrifyvault.log"
    skip_cert_verify = true
}
