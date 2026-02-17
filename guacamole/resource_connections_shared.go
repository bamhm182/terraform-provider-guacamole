package guacamole

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// connectionAttributesSchemaElem returns the common schema resource for
// connection-level attributes (guacd proxy settings, load balancing, etc.).
func connectionAttributesSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"guacd_hostname": {
				Type:        schema.TypeString,
				Description: "Guacd proxy hostname",
				Optional:    true,
				Computed:    true,
			},
			"guacd_port": {
				Type:        schema.TypeString,
				Description: "Guacd proxy port",
				Optional:    true,
				Computed:    true,
			},
			"guacd_encryption": {
				Type:        schema.TypeString,
				Description: "Guacd proxy encryption type",
				Optional:    true,
				Computed:    true,
			},
			"failover_only": {
				Type:        schema.TypeBool,
				Description: "Use load balancing for failover only",
				Optional:    true,
				Computed:    true,
			},
			"weight": {
				Type:        schema.TypeString,
				Description: "Load balancing connection weight",
				Optional:    true,
				Computed:    true,
			},
			"max_connections": {
				Type:        schema.TypeString,
				Description: "Maximum concurrent total connections",
				Optional:    true,
				Computed:    true,
			},
			"max_connections_per_user": {
				Type:        schema.TypeString,
				Description: "Maximum concurrent connections per user",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// readConnectionAttributes maps the Attributes NullableStringMap back into the
// Terraform "attributes" list block.
func readConnectionAttributes(attrs map[string]string) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"guacd_hostname":           attrs["guacd-hostname"],
			"guacd_port":               attrs["guacd-port"],
			"guacd_encryption":         attrs["guacd-encryption"],
			"failover_only":            stringToBool(attrs["failover-only"]),
			"weight":                   attrs["weight"],
			"max_connections":          attrs["max-connections"],
			"max_connections_per_user": attrs["max-connections-per-user"],
		},
	}
}

// buildConnectionAttributes converts the Terraform "attributes" list block
// into Guacamole's NullableStringMap key names.
func buildConnectionAttributes(attrList []interface{}) map[string]string {
	if len(attrList) == 0 {
		return nil
	}
	attrs := attrList[0].(map[string]interface{})
	return map[string]string{
		"guacd-hostname":           attrs["guacd_hostname"].(string),
		"guacd-port":               attrs["guacd_port"].(string),
		"guacd-encryption":         attrs["guacd_encryption"].(string),
		"failover-only":            boolToString(attrs["failover_only"].(bool)),
		"weight":                   attrs["weight"].(string),
		"max-connections":          attrs["max_connections"].(string),
		"max-connections-per-user": attrs["max_connections_per_user"].(string),
	}
}
