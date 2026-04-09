package main

import (
	"context"
	"log"

	"github.com/Azure/terraform-provider-azapi/internal/provider"
	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
)

func main() {
	ctx := context.Background()
	gen, err := tffwdocs.NewGenerator(ctx, &provider.Provider{})
	if err != nil {
		log.Fatal(err)
	}

	if err := gen.Lint(nil); err != nil {
		log.Fatal(err)
	}

	if err := gen.WriteAll(ctx, "./docs", nil); err != nil {
		log.Fatal(err)
	}

	// Special handling of data plane resource
	if err := genDataPlaneResource(ctx, gen); err != nil {
		log.Fatal(err)
	}
}
