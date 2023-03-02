package provider

import "github.com/hashicorp/terraform-plugin-framework/providerserver"
import "github.com/hashicorp/terraform-plugin-go/tfprotov6"

var testProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
	"tf-split-policies": providerserver.NewProtocol6WithError(New()()),
}
