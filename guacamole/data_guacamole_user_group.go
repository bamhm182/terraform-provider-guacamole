package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of guacamole user group",
				Required:    true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"parent_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"system_permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connections": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connection_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Get("identifier").(string)

	group, err := client.GetUserGroup(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read user group %s: %w", id, err))
	}

	d.Set("identifier", group.Identifier)
	d.Set("attributes", []interface{}{
		map[string]interface{}{
			"disabled": stringToBool(group.Attributes["disabled"]),
		},
	})

	parentGroups, err := client.GetUserGroupParentGroups(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get parent groups %s: %w", id, err))
	}
	d.Set("parent_groups", parentGroups)

	memberGroups, err := client.GetUserGroupMemberGroups(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get member groups %s: %w", id, err))
	}
	d.Set("member_groups", memberGroups)

	memberUsers, err := client.GetUserGroupMemberUsers(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get member users %s: %w", id, err))
	}
	d.Set("member_users", memberUsers)

	permissions, err := client.GetUserGroupPermissions(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get user group permissions %s: %w", id, err))
	}

	d.Set("system_permissions", permissions.SystemPermissions)

	var connections []string
	for connID := range permissions.ConnectionPermissions {
		connections = append(connections, connID)
	}
	d.Set("connections", connections)

	var connectionGroups []string
	for cgID := range permissions.ConnectionGroupPermissions {
		connectionGroups = append(connectionGroups, cgID)
	}
	d.Set("connection_groups", connectionGroups)

	d.SetId(id)

	return nil
}
