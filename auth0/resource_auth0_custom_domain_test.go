package auth0

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/auth0/terraform-provider-auth0/auth0/internal/random"
)

func init() {
	resource.AddTestSweepers("auth0_custom_domain", &resource.Sweeper{
		Name: "auth0_custom_domain",
		F: func(_ string) (err error) {
			api, err := Auth0()
			if err != nil {
				return
			}
			domains, err := api.CustomDomain.List()
			if err != nil {
				return
			}
			for _, domain := range domains {
				log.Printf("[DEBUG] ➝ %s", domain.GetDomain())
				if strings.Contains(domain.GetDomain(), "auth.uat.terraform-provider-auth0.com") {
					if e := api.CustomDomain.Delete(domain.GetID()); e != nil {
						multierror.Append(err, e)
					}
					log.Printf("[DEBUG] ✗ %s", domain.GetDomain())
				}
			}
			return
		},
	})
}

func TestAccCustomDomain(t *testing.T) {
	rand := random.String(6)

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"auth0": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: random.Template(testAccCustomDomain, rand),
				Check: resource.ComposeTestCheckFunc(
					random.TestCheckResourceAttr("auth0_custom_domain.my_custom_domain", "domain", "{{.random}}.auth.uat.terraform-provider-auth0.com", rand),
					resource.TestCheckResourceAttr("auth0_custom_domain.my_custom_domain", "type", "auth0_managed_certs"),
					resource.TestCheckResourceAttr("auth0_custom_domain.my_custom_domain", "status", "pending_verification"),
				),
			},
		},
	})
}

const testAccCustomDomain = `
resource "auth0_custom_domain" "my_custom_domain" {
  domain = "{{.random}}.auth.uat.terraform-provider-auth0.com"
  type = "auth0_managed_certs"
}
`
