package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataStreamsDatasource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Expect Error on invalid resource
			{
				Config: providerConfig + `
				data "analyticsadmin_datastreams" "streams" {
				}
				`,
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
			// Create a BigQuery link
			{
				Config: providerConfig + `
				data "analyticsadmin_datastreams" "streams" {
					property_id = 480154826
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.analyticsadmin_datastreams.streams",
						tfjsonpath.New("data_streams"),
						knownvalue.Null(),
					),
				},
			},
			// Delete
		},
	})
}
