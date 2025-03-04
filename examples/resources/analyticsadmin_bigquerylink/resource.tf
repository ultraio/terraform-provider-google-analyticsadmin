# Copyright (c) HashiCorp, Inc.

resource "analyticsadmin_bigquerylink" "link" {
  project          = "projects/1009831384334"
  property_id      = "480154826"
  dataset_location = "europe-west1"
}
