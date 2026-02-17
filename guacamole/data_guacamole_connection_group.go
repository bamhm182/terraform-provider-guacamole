package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConnectionGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionGroupRead,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of guacamole connection group",
				Required:    true,
			},
			"parent_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_connections": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections_per_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_session_affinity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"member_connection_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"member_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceConnectionGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	identifier := d.Get("identifier").(string)

	// Use tree to get children along with the group itself
	tree, err := client.GetConnectionGroupTree(ctx, identifier)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read connection group %s: %w", identifier, err))
	}

	d.Set("name", tree.Name)
	d.Set("parent_identifier", tree.ParentIdentifier)
	d.Set("type", tree.Type)
	d.Set("active_connections", tree.ActiveConnections)
	d.Set("attributes", []interface{}{
		map[string]interface{}{
			"max_connections":         tree.Attributes["max-connections"],
			"max_connections_per_user": tree.Attributes["max-connections-per-user"],
			"enable_session_affinity": stringToBool(tree.Attributes["enable-session-affinity"]),
		},
	})

	var memberGroups []interface{}
	for _, g := range tree.ChildConnectionGroups {
		memberGroups = append(memberGroups, map[string]interface{}{
			"identifier":         g.Identifier,
			"parent_identifier":  g.ParentIdentifier,
			"name":               g.Name,
			"type":               g.Type,
			"active_connections": g.ActiveConnections,
		})
	}
	d.Set("member_connection_groups", memberGroups)

	var memberConnections []interface{}
	for _, c := range tree.ChildConnections {
		memberConnections = append(memberConnections, map[string]interface{}{
			"identifier":         c.Identifier,
			"parent_identifier":  c.ParentIdentifier,
			"name":               c.Name,
			"protocol":           c.Protocol,
			"active_connections": c.ActiveConnections,
		})
	}
	d.Set("member_connections", memberConnections)

	d.SetId(identifier)

	return nil
}
