package main

import (
	"context"
	"flag"
	"log"

	"github.com/Azure/terraform-provider-azapi/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

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
