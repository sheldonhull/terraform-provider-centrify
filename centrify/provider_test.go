package centrify

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"centrifyvault": testAccProvider,
	}
}
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
func testAccPreCheck(t *testing.T) {

	if v := os.Getenv("VAULT_URL"); v == "" {
		t.Fatal("PAS URL must be set for acceptance tests")
	}

	if v := os.Getenv("VAULT_SCOPE"); v == "" {
		t.Fatal("SCOPE must be set for acceptance tests")
	}

	if v := os.Getenv("VAULT_USEDMC"); v == "" {

		if v := os.Getenv("VAULT_TOKEN"); v == "" {

			if v := os.Getenv("VAULT_APPID"); v == "" {
				t.Fatal("APPID must be set for acceptance tests")
			}
			if v := os.Getenv("VAULT_USERNAME"); v == "" {
				t.Fatal("USERNAME must be set for acceptance tests")
			}
			if v := os.Getenv("VAULT_PASSWORD"); v == "" {
				t.Fatal("PASSWORD must be set for acceptance tests")
			}
		}
	}
}
