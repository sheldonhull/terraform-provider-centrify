# Terraform Provider for Centrify Vault

The Terraform Provider for Centrify Vault is a Terraform plugin that allows other Terraform providers to retrieve vaulted password or secret from Centrify Vault. It also enables full configuration management of Centrify Vault.


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x or higher
- [Go](https://golang.org/doc/install) 1.13 or higher (to build the provider plugin)


## Building The Provider

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-centrify`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/centrify/terraform-provider terraform-provider-centrify
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-centrify
$ make build
```

To install the provider in your home directory

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-centrify
$ make install
```

## Using The Provider

The provider needs to be configured with the proper credentials before it can be used.

### Example Usage

```terraform
# Configure CentrifyVault Provider
provider "centrifyvault" {
    url = "https://tenantid.my.centrify.net"
    appid = "CentrifyCLI"
    scope = "terraform"
    token = "xxxxxxxxxxxxxxxxxx"
}
```

### Authentication and Argument Reference

The Provider supports OAuth2 and DMC authentication methods.

* **url** - (Required) This is the cloud tenant or on-prem PAS URL. It must be provided, but it can also be sourced from the VAULT_URL environment variable.
* **appid** - (Required) This is the OAuth application ID configured in Centrify Vault. It must be provided, but it can also be sourced from the VAULT_APPID environment variable.
* **scope** - (Required) This is either the OAuth or DMC scope. It must be provided, but it can also be sourced from the VAULT_SCOPE environment variable.
* **token** - (Optional) This is the Oauth token. It can also be sourced from the VAULT_TOKEN environment variable.
* **username** - (Optional) Authorized user to retrieve Oauth token. It can also be sourced from the VAULT_USERNAME environment variable.
* **username** - (Optional) Authorized user's password for retrieving Oauth token. It can also be sourced from the VAULT_PASSWORD environment variable.
* **use_dmc** - (Optional) Whether to use DMC authentication. It can also be sourced from the VAULT_USEDMC environment variable. The default is false.
* **skip_cert_verify** - (Optional) Whether to skip certificate validation. It is used for testing against on-prem PAS deployment which uses self-signed certificate. It can also be sourced from the VAULT_SKIPCERTVERIFY environment variable. The default is false.

## Checkout Credentials

Following example shows how to retrieve password for a vaulted Linux account.

```terraform
// data source for system "centos"
data "centrifyvault_vaultsystem" "centos" {
    name = "centos"
    fqdn = "centos.demo.lab"
    computer_class = "Unix"
}

// data source for account "local_account" in "centos"
// Specify checkout = true to checkout the password
data "centrifyvault_vaultaccount" "centos_local_account" {
    name = "local_account"
    host_id = data.centrifyvault_vaultsystem.centos.id
    checkout = true
}

// Output retrieved password for account "local_account"
output "centos_local_account" {
    value = data.centrifyvault_vaultaccount.centos_local_account.password
}

```

Following example shows how to retrieve a secret.

```terraform
// data source for secret named "My secret key"
// Specify checkout = true to checkout the secret
data "centrifyvault_vaultsecret" "my_secret_key" {
    secret_name = "My secret key"
    checkout = true
}

// Output retrieved secret content
output "pas_admin_credential" {
    value = data.centrifyvault_vaultsecret.my_secret_key.secret_text
}
```
