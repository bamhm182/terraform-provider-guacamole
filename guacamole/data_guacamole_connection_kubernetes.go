package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConnectionKubernetes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionKubernetesRead,
		Schema: map[string]*schema.Schema{
			"identifier":         {Type: schema.TypeString, Required: true},
			"name":               {Type: schema.TypeString, Computed: true},
			"parent_identifier":  {Type: schema.TypeString, Computed: true},
			"protocol":           {Type: schema.TypeString, Computed: true},
			"active_connections": {Type: schema.TypeInt, Computed: true},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     connectionAttributesSchemaElem(),
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname":                    {Type: schema.TypeString, Computed: true},
						"port":                        {Type: schema.TypeString, Computed: true},
						"use_ssl":                     {Type: schema.TypeBool, Computed: true},
						"ignore_cert":                 {Type: schema.TypeBool, Computed: true},
						"ca_cert":                     {Type: schema.TypeString, Computed: true},
						"namespace":                   {Type: schema.TypeString, Computed: true},
						"pod":                         {Type: schema.TypeString, Computed: true},
						"container":                   {Type: schema.TypeString, Computed: true},
						"client_cert":                 {Type: schema.TypeString, Computed: true},
						"client_key":                  {Type: schema.TypeString, Computed: true, Sensitive: true},
						"color_scheme":                {Type: schema.TypeString, Computed: true},
						"font_name":                   {Type: schema.TypeString, Computed: true},
						"font_size":                   {Type: schema.TypeString, Computed: true},
						"max_scrollback_size":         {Type: schema.TypeString, Computed: true},
						"readonly":                    {Type: schema.TypeBool, Computed: true},
						"backspace":                   {Type: schema.TypeString, Computed: true},
						"typescript_path":             {Type: schema.TypeString, Computed: true},
						"typescript_name":             {Type: schema.TypeString, Computed: true},
						"typescript_auto_create_path": {Type: schema.TypeBool, Computed: true},
						"recording_path":              {Type: schema.TypeString, Computed: true},
						"recording_name":              {Type: schema.TypeString, Computed: true},
						"recording_exclude_output":    {Type: schema.TypeBool, Computed: true},
						"recording_exclude_mouse":     {Type: schema.TypeBool, Computed: true},
						"recording_include_keys":      {Type: schema.TypeBool, Computed: true},
						"recording_auto_create_path":  {Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceConnectionKubernetesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Get("identifier").(string)

	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read kubernetes connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read kubernetes connection parameters %s: %w", id, err))
	}

	d.Set("name", conn.Name)
	d.Set("parent_identifier", conn.ParentIdentifier)
	d.Set("protocol", conn.Protocol)
	d.Set("active_connections", conn.ActiveConnections)
	d.Set("attributes", readConnectionAttributes(conn.Attributes))
	d.Set("parameters", []interface{}{
		map[string]interface{}{
			"hostname":                    params["hostname"],
			"port":                        params["port"],
			"use_ssl":                     stringToBool(params["ssl"]),
			"ignore_cert":                 stringToBool(params["ignore-cert"]),
			"ca_cert":                     params["ca-cert"],
			"namespace":                   params["namespace"],
			"pod":                         params["pod"],
			"container":                   params["container"],
			"client_cert":                 params["client-cert"],
			"client_key":                  params["client-key"],
			"color_scheme":                params["color-scheme"],
			"font_name":                   params["font-name"],
			"font_size":                   params["font-size"],
			"max_scrollback_size":         params["scrollback"],
			"readonly":                    stringToBool(params["read-only"]),
			"backspace":                   params["backspace"],
			"typescript_path":             params["typescript-path"],
			"typescript_name":             params["typescript-name"],
			"typescript_auto_create_path": stringToBool(params["create-typescript-path"]),
			"recording_path":              params["recording-path"],
			"recording_name":              params["recording-name"],
			"recording_exclude_output":    stringToBool(params["recording-exclude-output"]),
			"recording_exclude_mouse":     stringToBool(params["recording-exclude-mouse"]),
			"recording_include_keys":      stringToBool(params["recording-include-keys"]),
			"recording_auto_create_path":  stringToBool(params["create-recording-path"]),
		},
	})

	d.SetId(id)

	return nil
}
