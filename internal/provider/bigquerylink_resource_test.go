package provider

import (
	"errors"
	"fmt"
	"regexp"
	"terraform-provider-google-analyticsadmin/internal/test_helper"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccBigQueryLinkResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Expect Error on invalid resource
			{
				Config: providerConfig + `
				resource "analyticsadmin_bigquerylink" "link" {
				}
				`,
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
			{
				Config: providerConfig + `
				resource "analyticsadmin_bigquerylink" "link" {
					project = "1234"
					property_id = "1234"
					dataset_location = "europe-west1"
				}
				`,
				ExpectError: regexp.MustCompile("Must match Format: 'projects/{project number}'"),
			},
			{
				Config: providerConfig + `
				resource "analyticsadmin_bigquerylink" "link" {
					project = "projects/1234"
					property_id = "1234"	
					dataset_location = "europe-west1"
					export_streams = ["mystream"]
				}
				`,
				ExpectError: regexp.MustCompile(`'properties/{property_id}/dataStreams/{stream_id}'`),
			},
			// Create a BigQuery link
			{
				Config: testAccBigQueryLinkResourceBuilder_create(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectKnownValue(
							"analyticsadmin_bigquerylink.link",
							tfjsonpath.New("property_id"),
							knownvalue.StringExact("480154826"),
						),
						plancheck.ExpectKnownValue(
							"analyticsadmin_bigquerylink.link",
							tfjsonpath.New("dataset_location"),
							knownvalue.StringExact("europe-west1"),
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"analyticsadmin_bigquerylink.link",
						tfjsonpath.New("name"),
						knownvalue.StringRegexp(regexp.MustCompile(`properties/\w+/bigQueryLinks/\w+`)),
					),
					test_helper.ExpectTimeFormat(
						"analyticsadmin_bigquerylink.link",
						tfjsonpath.New("create_time"),
						time.RFC3339,
					),
				},
			},
			// Update a BigQuery link
			{
				Config: testAccBigQueryLinkResourceBuilder_update(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"analyticsadmin_bigquerylink.link",
						tfjsonpath.New("daily_export"),
						knownvalue.Bool(true),
					),
					test_helper.ExpectTimeFormat(
						"analyticsadmin_bigquerylink.link",
						tfjsonpath.New("create_time"),
						time.RFC3339,
					),
				},
			},
			// ImportState testing
			{
				ResourceName: "analyticsadmin_bigquerylink.link",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					var name string
					for _, m := range s.Modules {
						if len(m.Resources) > 0 {
							if v, ok := m.Resources["analyticsadmin_bigquerylink.link"]; ok {
								name = v.Primary.Attributes["name"]
							}
						}
					}

					if len(name) == 0 {
						return "", errors.New("Failed to find resource in state")
					}

					return name, nil
				},
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Delete
		},
	})
}

func testAccBigQueryLinkResourceBuilder_create() string {
	link := func(name, project, property_id string) string {
		return fmt.Sprintf(`
resource "analyticsadmin_bigquerylink" "%s" {
	project = "projects/%s"							
	property_id = "%s"	
	dataset_location = "europe-west1"
}
		`, name, project, property_id)
	}

	return providerConfig + link("link", "1009831384334", "480154826")
}

func testAccBigQueryLinkResourceBuilder_update() string {
	link := func(name, project, property_id string) string {
		return fmt.Sprintf(`
resource "analyticsadmin_bigquerylink" "%s" {
	project = "projects/%s"							
	property_id = "%s"	
	dataset_location = "europe-west1"
	daily_export = true
}
		`, name, project, property_id)
	}

	return providerConfig + link("link", "1009831384334", "480154826")
}
