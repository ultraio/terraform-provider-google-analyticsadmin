package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "analyticsadmin" {
	scopes = ["https://www.googleapis.com/auth/analytics.edit"]
}`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"analyticsadmin": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestAccProvider_config(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Config execution require to create one resource or datasource
			{
				Config: `
provider "analyticsadmin" {}
resource "analyticsadmin_bigquerylink" "link" { 
    project = "projects/fake"
	property_id = "fake"
	dataset_location = "fake" 
}`,
				ExpectError: regexp.MustCompile("The argument \"scopes\" is required, but no definition was found"),
			},
		},
	})
}
