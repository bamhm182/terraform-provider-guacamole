package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleSharingProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharingProfileCreate,
		ReadContext:   resourceSharingProfileRead,
		UpdateContext: resourceSharingProfileUpdate,
		DeleteContext: resourceSharingProfileDelete,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of guacamole sharing profile",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of guacamole sharing profile",
				Required:    true,
			},
			"primary_connection_identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of the primary connection being shared",
				Required:    true,
			},
			"parameters": {
				Type:        schema.TypeList,
				Description: "Parameters of guacamole sharing profile",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_only": {
							Type:        schema.TypeBool,
							Description: "Whether the shared session is read-only",
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

func resourceSharingProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	profile := convertResourceDataToGuacSharingProfile(d)

	created, err := client.CreateSharingProfile(ctx, profile)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("identifier", created.Identifier)
	d.SetId(created.Identifier)

	return resourceSharingProfileRead(ctx, d, m)
}

func resourceSharingProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	profile, err := client.GetSharingProfile(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read sharing profile %s: %w", id, err))
	}

	params, err := client.GetSharingProfileParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read sharing profile parameters %s: %w", id, err))
	}
	profile.Parameters = params

	convertGuacSharingProfileToResourceData(d, profile)
	d.SetId(id)

	return nil
}

func resourceSharingProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("name", "primary_connection_identifier", "parameters") {
		profile := convertResourceDataToGuacSharingProfile(d)
		if err := client.UpdateSharingProfile(ctx, d.Id(), profile); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSharingProfileRead(ctx, d, m)
}

func resourceSharingProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteSharingProfile(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func convertResourceDataToGuacSharingProfile(d *schema.ResourceData) guacamole.SharingProfile {
	profile := guacamole.SharingProfile{
		Identifier:                  d.Get("identifier").(string),
		Name:                        d.Get("name").(string),
		PrimaryConnectionIdentifier: d.Get("primary_connection_identifier").(string),
		Attributes:                  guacamole.NullableStringMap{},
	}

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) > 0 {
		p := paramList[0].(map[string]interface{})
		profile.Parameters = map[string]string{
			"read-only": boolToString(p["read_only"].(bool)),
		}
	}

	return profile
}

func convertGuacSharingProfileToResourceData(d *schema.ResourceData, profile *guacamole.SharingProfile) {
	d.Set("identifier", profile.Identifier)
	d.Set("name", profile.Name)
	d.Set("primary_connection_identifier", profile.PrimaryConnectionIdentifier)

	parameters := map[string]interface{}{
		"read_only": stringToBool(profile.Parameters["read-only"]),
	}
	d.Set("parameters", []interface{}{parameters})
}
