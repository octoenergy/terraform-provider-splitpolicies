package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/octoenergy/tf-split-policies/internal/provider"
)

func main() {

	opts := providerserver.ServeOpts {
	    Address: "github.com/octoenergy/tf-split-policies",
	    }

	err := providerserver.Serve(context.Background(), provider.New(), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
