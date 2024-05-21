package main

import (
	"context"
	"flag"
	"log"

	"github.com/Azure/terraform-provider-azapi/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	serveOpts := providerserver.ServeOpts{
		Debug:   debugMode,
		Address: "registry.terraform.io/Azure/azapi",
	}

	err := providerserver.Serve(context.Background(), provider.AzureProvider, serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
