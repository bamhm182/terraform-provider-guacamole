package guacamole

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleConnectionGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionGroupCreate,
		ReadContext:   resourceConnectionGroupRead,
		UpdateContext: resourceConnectionGroupUpdate,
		DeleteContext: resourceConnectionGroupDelete,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of guacamole connection group",
				Computed:    true,
			},
			"parent_identifier": {
				Type:        schema.TypeString,
				Description: "Parent identifier of guacamole connection group",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of guacamole connection group",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of guacamole connection group (ORGANIZATIONAL or BALANCING)",
				Optional:    true,
				Default:     "ORGANIZATIONAL",
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"active_connections": {
				Type:        schema.TypeInt,
				Description: "Active connections of guacamole connection group",
				Computed:    true,
			},
			"attributes": {
				Type:        schema.TypeList,
				Description: "Attributes of guacamole connection group",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_connections": {
							Type:        schema.TypeString,
							Description: "Maximum number of total simultaneous connections allowed",
							Optional:    true,
							Computed:    true,
						},
						"max_connections_per_user": {
							Type:        schema.TypeString,
							Description: "Maximum number of simultaneous connections allowed per user",
							Optional:    true,
							Computed:    true,
						},
						"enable_session_affinity": {
							Type:        schema.TypeBool,
							Description: "Enable session affinity",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if check := validateConnectionGroup(d); check.HasError() {
		return check
	}

	group := convertResourceDataToGuacConnectionGroup(d)

	created, err := client.CreateConnectionGroup(ctx, group)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("identifier", created.Identifier)
	d.SetId(created.Identifier)

	return resourceConnectionGroupRead(ctx, d, m)
}

func resourceConnectionGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	group, err := client.GetConnectionGroup(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read connection group %s: %w", id, err))
	}

	convertGuacConnectionGroupToResourceData(d, group)
	d.SetId(id)

	return nil
}

func resourceConnectionGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("name", "parent_identifier", "type", "attributes") {
		if check := validateConnectionGroup(d); check.HasError() {
			return check
		}
		group := convertResourceDataToGuacConnectionGroup(d)
		if err := client.UpdateConnectionGroup(ctx, d.Id(), group); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceConnectionGroupRead(ctx, d, m)
}

func resourceConnectionGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteConnectionGroup(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func convertResourceDataToGuacConnectionGroup(d *schema.ResourceData) guacamole.ConnectionGroup {
	group := guacamole.ConnectionGroup{
		Identifier:       d.Get("identifier").(string),
		ParentIdentifier: d.Get("parent_identifier").(string),
		Name:             d.Get("name").(string),
		Type:             strings.ToUpper(d.Get("type").(string)),
	}

	attrList := d.Get("attributes").([]interface{})
	if len(attrList) > 0 {
		attrs := attrList[0].(map[string]interface{})
		group.Attributes = guacamole.NullableStringMap{
			"max-connections":          attrs["max_connections"].(string),
			"max-connections-per-user": attrs["max_connections_per_user"].(string),
			"enable-session-affinity":  boolToString(attrs["enable_session_affinity"].(bool)),
		}
	}

	return group
}

func convertGuacConnectionGroupToResourceData(d *schema.ResourceData, group *guacamole.ConnectionGroup) {
	d.Set("identifier", group.Identifier)
	d.Set("parent_identifier", group.ParentIdentifier)
	d.Set("name", group.Name)
	d.Set("type", group.Type)
	d.Set("active_connections", group.ActiveConnections)

	attributes := map[string]interface{}{
		"max_connections":          group.Attributes["max-connections"],
		"max_connections_per_user": group.Attributes["max-connections-per-user"],
		"enable_session_affinity":  stringToBool(group.Attributes["enable-session-affinity"]),
	}
	d.Set("attributes", []interface{}{attributes})
}

func validateConnectionGroup(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	groupType := strings.ToUpper(d.Get("type").(string))
	if check := stringInSlice([]string{"ORGANIZATIONAL", "BALANCING"}, []string{groupType}); check.HasError() {
		diags = append(diags, check...)
	}

	attrList := d.Get("attributes").([]interface{})
	if len(attrList) > 0 {
		attrs := attrList[0].(map[string]interface{})
		for _, key := range []string{"max_connections", "max_connections_per_user"} {
			val := attrs[key].(string)
			if val != "" {
				if _, err := strconv.Atoi(val); err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid entry",
						Detail:   fmt.Sprintf("Expected string integer for %s, got: %s", key, val),
					})
				}
			}
		}
	}

	return diags
}
