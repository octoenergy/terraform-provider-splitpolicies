package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestDataSourceEmpty(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `data "splitpolicies" "test" {}`,
				ExpectError: regexp.MustCompile("The argument \"policies\" is required, but no definition was found."),
			},
		},
	})
}

func TestSmallDataSourceOneChunk(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testSmallDataSourceOneChunkConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.%", "1"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.#", "3"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.0", "one"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.1", "two"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.2", "three"),
				),
			},
		},
	})
}

const testSmallDataSourceOneChunkConfig = `
data "splitpolicies" "test" {
  policies = ["one", "two", "three"]
}
`

func TestSmallDataSourceManyChunks(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testSmallDataSourceManyChunksConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.%", "2"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.#", "2"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.1.#", "1"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.0", "one"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.0.1", "two"),
					resource.TestCheckResourceAttr("data.splitpolicies.test", "chunks.1.0", "three"),
				),
			},
		},
	})
}

const testSmallDataSourceManyChunksConfig = `
data "splitpolicies" "test" {
  policies = ["one", "two", "three"]
  maximum_chunk_size = 6
}
`
