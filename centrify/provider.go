package centrify

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var logPath string

// Provider returns a schema.Provider for Centrify Vault.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_URL", ""),
				Description: "Centrify Vault URL",
			},
			"appid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_APPID", ""),
				Description: "Application ID",
			},
			"scope": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_SCOPE", ""),
				Description: "OAuth2 scope",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_USERNAME", ""),
				Description: "Username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_PASSWORD", ""),
				Description: "Password",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_TOKEN", ""),
				Description: "OAuth or DMC token",
			},
			"use_dmc": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_USEDMC", false),
				Description: "Whether to use DMC",
			},
			"logpath": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_LOGPATH", ""),
				Description: "Path of log file",
			},
			"skip_cert_verify": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_SKIPCERTVERIFY", false),
				Description: "Whether to skip certification verification",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"centrifyvault_user":                  dataSourceUser(),
			"centrifyvault_role":                  dataSourceRole(),
			"centrifyvault_policy":                dataSourcePolicy(),
			"centrifyvault_manualset":             dataSourceManualSet(),
			"centrifyvault_passwordprofile":       dataSourcePasswordProfile(),
			"centrifyvault_authenticationprofile": dataSourceAuthenticationProfile(),
			"centrifyvault_connector":             dataSourceConnector(),
			"centrifyvault_vaultdomain":           dataSourceVaultDomain(),
			"centrifyvault_vaultsystem":           dataSourceVaultSystem(),
			"centrifyvault_vaultaccount":          dataSourceVaultAccount(),
			"centrifyvault_vaultsecret":           dataSourceVaultSecret(),
			"centrifyvault_vaultsecretfolder":     dataSourceVaultSecretFolder(),
			"centrifyvault_sshkey":                dataSourceSSHKey(),
			"centrifyvault_directoryservice":      dataSourceDirectoryService(),
			"centrifyvault_directoryobject":       dataSourceDirectoryObject(),
			"centrifyvault_multiplexedaccount":    dataSourceMultiplexedAccount(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"centrifyvault_user":                  resourceUser(),
			"centrifyvault_role":                  resourceRole(),
			"centrifyvault_policyorder":           resourcePolicyLinks(),
			"centrifyvault_policy":                resourcePolicy(),
			"centrifyvault_manualset":             resourceManualSet(),
			"centrifyvault_passwordprofile":       resourcePasswordProfile(),
			"centrifyvault_authenticationprofile": resourceAuthenticationProfile(),
			"centrifyvault_vaultdomain":           resourceVaultDomain(),
			"centrifyvault_vaultsystem":           resourceVaultSystem(),
			"centrifyvault_vaultdatabase":         resourceVaultDatabase(),
			"centrifyvault_vaultaccount":          resourceVaultAccount(),
			"centrifyvault_vaultsecret":           resourceVaultSecret(),
			"centrifyvault_vaultsecretfolder":     resourceVaultSecretFolder(),
			"centrifyvault_sshkey":                resourceSSHKey(),
			"centrifyvault_desktopapp":            resourceDesktopApp(),
			"centrifyvault_multiplexedaccount":    resourceMultiplexedAccount(),
			"centrifyvault_service":               resourceService(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	LogD.Printf("Running providerConfigure...")

	config := Config{
		URL:            d.Get("url").(string),
		AppID:          d.Get("appid").(string),
		Scope:          d.Get("scope").(string),
		Username:       d.Get("username").(string),
		Password:       d.Get("password").(string),
		Token:          d.Get("token").(string),
		UseDMC:         d.Get("use_dmc").(bool),
		LogPath:        d.Get("logpath").(string),
		SkipCertVerify: d.Get("skip_cert_verify").(bool),
	}
	logPath = config.LogPath

	if err := config.Valid(); err != nil {
		return nil, err
	}

	restClient, err := config.getClient()

	if err != nil {
		return nil, fmt.Errorf("Unable to get oauth rest client: %v", err)
	}
	LogD.Printf("Connected to Centrify Vault %s", config.URL)

	return restClient, nil
}
