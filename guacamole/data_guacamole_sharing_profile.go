package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSharingProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharingProfileRead,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of guacamole sharing profile",
				Required:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_connection_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSharingProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	identifier := d.Get("identifier").(string)

	profile, err := client.GetSharingProfile(ctx, identifier)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read sharing profile %s: %w", identifier, err))
	}

	params, err := client.GetSharingProfileParameters(ctx, identifier)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read sharing profile parameters %s: %w", identifier, err))
	}
	profile.Parameters = params

	d.Set("name", profile.Name)
	d.Set("primary_connection_identifier", profile.PrimaryConnectionIdentifier)
	d.Set("parameters", []interface{}{
		map[string]interface{}{
			"read_only": stringToBool(profile.Parameters["read-only"]),
		},
	})

	d.SetId(identifier)

	return nil
}
