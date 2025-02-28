package main

import (
	"context"
	"log"
	"terraform-provider-google-analyticsadmin/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/ultraio/google-analyticsadmin",
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
