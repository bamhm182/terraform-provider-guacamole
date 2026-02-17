package guacamole

import (
	"context"
	"fmt"
	"testing"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testProviderUserGroup = map[string]interface{}{
	"identifier":         "testProviderUserGroup",
	"system_permissions": validSystemPermissions(),
}

func TestAccGuacamoleUserGroupBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGuacamoleUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGuacamoleUserGroupConfigBasic(toHclString(testProviderUserGroup, true)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGuacamoleUserGroupExists("guacamole_user_group.new"),
					resource.TestCheckResourceAttr("guacamole_user_group.new", "identifier", testProviderUserGroup["identifier"].(string)),
					testAccCheckTestSliceVals("guacamole_user_group.new", "system_permissions", validSystemPermissions()),
				),
			},
		},
	})
}

func testAccCheckGuacamoleUserGroupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*guacamole.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "guacamole_user_group" {
			continue
		}

		identifier := rs.Primary.ID

		_, err := c.GetUserGroup(context.Background(), identifier)
		if err != nil {
			if guacamole.IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("user group %s still exists in guacamole database", identifier)
	}

	return nil
}

func testAccCheckGuacamoleUserGroupConfigBasic(definition string) string {
	return fmt.Sprintf(`
	resource "guacamole_user_group" "new" %s
	`, definition)
}

func TestAccGuacamoleUserGroupMembership(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGuacamoleUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGuacamoleUserGroupMembershipConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGuacamoleUserGroupExists("guacamole_user_group.parent"),
					resource.TestCheckResourceAttr("guacamole_user_group.parent", "identifier", "testParentGroup"),
					testAccCheckTestSliceVals("guacamole_user_group.parent", "member_groups", []string{"testMemberGroup"}),
					testAccCheckTestSliceVals("guacamole_user_group.parent", "member_users", []string{"testMemberUser"}),
				),
			},
		},
	})
}

func testAccCheckGuacamoleUserGroupMembershipConfig() string {
	return `
	resource "guacamole_user_group" "member_group" {
		identifier = "testMemberGroup"
	}

	resource "guacamole_user" "member_user" {
		username = "testMemberUser"
	}

	resource "guacamole_user_group" "parent" {
		identifier   = "testParentGroup"
		member_groups = [guacamole_user_group.member_group.identifier]
		member_users  = [guacamole_user.member_user.username]
	}
	`
}

func testAccCheckGuacamoleUserGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No group identifier set")
		}

		return nil
	}
}
