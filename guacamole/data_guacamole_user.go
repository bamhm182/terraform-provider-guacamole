package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "Username of guacamole user",
				Required:    true,
			},
			"last_active": {
				Type:        schema.TypeString,
				Description: "Epoch time string of last user activity",
				Computed:    true,
			},
			"attributes": {
				Type:        schema.TypeList,
				Description: "Attributes of guacamole user",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organizational_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_window_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_window_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"valid_from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid_until": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"group_membership": {
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

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	username := d.Get("username").(string)

	user, err := client.GetUser(ctx, username)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read user %s: %w", username, err))
	}

	d.Set("last_active", fmt.Sprintf("%d", user.LastActive))
	d.Set("attributes", []interface{}{
		map[string]interface{}{
			"organizational_role": user.Attributes["guac-organizational-role"],
			"full_name":           user.Attributes["guac-full-name"],
			"email":               user.Attributes["guac-email-address"],
			"expired":             stringToBool(user.Attributes["expired"]),
			"timezone":            user.Attributes["timezone"],
			"access_window_start": user.Attributes["access-window-start"],
			"access_window_end":   user.Attributes["access-window-end"],
			"disabled":            stringToBool(user.Attributes["disabled"]),
			"valid_from":          user.Attributes["valid-from"],
			"valid_until":         user.Attributes["valid-until"],
		},
	})

	groups, err := client.GetUserGroups(ctx, username)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get user groups %s: %w", username, err))
	}
	d.Set("group_membership", groups)

	permissions, err := client.GetUserPermissions(ctx, username)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get user permissions %s: %w", username, err))
	}

	d.Set("system_permissions", permissions.SystemPermissions)

	var connections []string
	for id := range permissions.ConnectionPermissions {
		connections = append(connections, id)
	}
	d.Set("connections", connections)

	var connectionGroups []string
	for id := range permissions.ConnectionGroupPermissions {
		connectionGroups = append(connectionGroups, id)
	}
	d.Set("connection_groups", connectionGroups)

	d.SetId(username)

	return nil
}
