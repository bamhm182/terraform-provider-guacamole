package guacamole

import (
	"context"
	"fmt"
	"testing"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testProviderSharingProfile = map[string]interface{}{
	"name": "testProviderSharingProfile",
	"parameters": map[string]interface{}{
		"read_only": true,
	},
}

func TestAccGuacamoleSharingProfileBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGuacamoleSharingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGuacamoleSharingProfileConfigBasic(testProviderSharingProfile),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGuacamoleSharingProfileExists("guacamole_sharing_profile.new"),
					resource.TestCheckResourceAttr("guacamole_sharing_profile.new", "name", testProviderSharingProfile["name"].(string)),
					resource.TestCheckResourceAttr("guacamole_sharing_profile.new", "parameters.0.read_only", boolToString(testProviderSharingProfile["parameters"].(map[string]interface{})["read_only"].(bool))),
				),
			},
		},
	})
}

func testAccCheckGuacamoleSharingProfileConfigBasic(definition map[string]interface{}) string {
	params := definition["parameters"].(map[string]interface{})
	return fmt.Sprintf(`
	resource "guacamole_connection_ssh" "base" {
		name              = "testSharingProfileBaseConnection"
		parent_identifier = "ROOT"
		parameters {
			hostname = "test.example.com"
		}
	}

	resource "guacamole_sharing_profile" "new" {
		name                         = %q
		primary_connection_identifier = guacamole_connection_ssh.base.identifier
		parameters {
			read_only = %v
		}
	}
	`, definition["name"].(string), params["read_only"].(bool))
}

func testAccCheckGuacamoleSharingProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No sharing profile ID set")
		}
		return nil
	}
}

func testAccCheckGuacamoleSharingProfileDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*guacamole.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "guacamole_sharing_profile" {
			continue
		}
		_, err := c.GetSharingProfile(context.Background(), rs.Primary.ID)
		if err != nil {
			if guacamole.IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("sharing profile %s still exists in guacamole database", rs.Primary.ID)
	}
	return nil
}
