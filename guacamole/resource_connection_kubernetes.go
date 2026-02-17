package guacamole

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleConnectionKubernetes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionKubernetesCreate,
		ReadContext:   resourceConnectionKubernetesRead,
		UpdateContext: resourceConnectionKubernetesUpdate,
		DeleteContext: resourceConnectionKubernetesDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the guacamole connection",
				Required:    true,
			},
			"identifier": {
				Type:        schema.TypeString,
				Description: "Numeric identifier of the guacamole connection",
				Computed:    true,
			},
			"parent_identifier": {
				Type:        schema.TypeString,
				Description: "Parent identifier of the guacamole connection",
				Optional:    true,
				Default:     "ROOT",
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Protocol type of the guacamole connection",
				Computed:    true,
			},
			"active_connections": {
				Type:        schema.TypeInt,
				Description: "Active connection count for the guacamole connection",
				Computed:    true,
			},
			"attributes": {
				Type:        schema.TypeList,
				Description: "Guacamole connection attributes",
				Optional:    true,
				MaxItems:    1,
				Elem:        connectionAttributesSchemaElem(),
			},
			"parameters": {
				Type:        schema.TypeList,
				Description: "Guacamole connection parameters",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname":                    {Type: schema.TypeString, Required: true},
						"port":                        {Type: schema.TypeString, Optional: true, Computed: true},
						"use_ssl":                     {Type: schema.TypeBool, Optional: true, Computed: true},
						"ignore_cert":                 {Type: schema.TypeBool, Optional: true, Computed: true},
						"ca_cert":                     {Type: schema.TypeString, Optional: true, Computed: true},
						"namespace":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"pod":                         {Type: schema.TypeString, Optional: true, Computed: true},
						"container":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"client_cert":                 {Type: schema.TypeString, Optional: true, Computed: true},
						"client_key":                  {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"color_scheme":                {Type: schema.TypeString, Optional: true, Computed: true},
						"font_name":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"font_size":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"max_scrollback_size":         {Type: schema.TypeString, Optional: true, Computed: true},
						"readonly":                    {Type: schema.TypeBool, Optional: true, Computed: true},
						"backspace":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"typescript_path":             {Type: schema.TypeString, Optional: true, Computed: true},
						"typescript_name":             {Type: schema.TypeString, Optional: true, Computed: true},
						"typescript_auto_create_path": {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_path":              {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_name":              {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_exclude_output":    {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_exclude_mouse":     {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_include_keys":      {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_auto_create_path":  {Type: schema.TypeBool, Optional: true, Computed: true},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionKubernetesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if check := validateConnectionKubernetes(d); check.HasError() {
		return check
	}

	conn := buildKubernetesConnection(d)
	created, err := client.CreateConnection(ctx, conn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("identifier", created.Identifier)
	d.SetId(created.Identifier)

	return resourceConnectionKubernetesRead(ctx, d, m)
}

func resourceConnectionKubernetesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read kubernetes connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read kubernetes connection parameters %s: %w", id, err))
	}

	d.Set("name", conn.Name)
	d.Set("identifier", conn.Identifier)
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

	return nil
}

func resourceConnectionKubernetesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("name", "parent_identifier", "attributes", "parameters") {
		if check := validateConnectionKubernetes(d); check.HasError() {
			return check
		}
		conn := buildKubernetesConnection(d)
		if err := client.UpdateConnection(ctx, d.Id(), conn); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceConnectionKubernetesRead(ctx, d, m)
}

func resourceConnectionKubernetesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteConnection(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func buildKubernetesConnection(d *schema.ResourceData) guacamole.Connection {
	conn := guacamole.Connection{
		Name:             d.Get("name").(string),
		Identifier:       d.Get("identifier").(string),
		ParentIdentifier: d.Get("parent_identifier").(string),
		Protocol:         "kubernetes",
	}

	conn.Attributes = buildConnectionAttributes(d.Get("attributes").([]interface{}))

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) > 0 {
		p := paramList[0].(map[string]interface{})
		conn.Parameters = map[string]string{
			"hostname":    p["hostname"].(string),
			"port":        p["port"].(string),
			"ssl":         boolToString(p["use_ssl"].(bool)),
			"ignore-cert": boolToString(p["ignore_cert"].(bool)),
			"ca-cert":     p["ca_cert"].(string),
			"namespace":   p["namespace"].(string),
			"pod":         p["pod"].(string),
			"container":   p["container"].(string),
			"client-cert": p["client_cert"].(string),
			"client-key":  p["client_key"].(string),
			"color-scheme": p["color_scheme"].(string),
			"font-name":   p["font_name"].(string),
			"font-size":   p["font_size"].(string),
			"scrollback":  p["max_scrollback_size"].(string),
			"read-only":   boolToString(p["readonly"].(bool)),
			"backspace":   p["backspace"].(string),
			"typescript-path":        p["typescript_path"].(string),
			"typescript-name":        p["typescript_name"].(string),
			"create-typescript-path": boolToString(p["typescript_auto_create_path"].(bool)),
			"recording-path":         p["recording_path"].(string),
			"recording-name":         p["recording_name"].(string),
			"recording-exclude-output": boolToString(p["recording_exclude_output"].(bool)),
			"recording-exclude-mouse":  boolToString(p["recording_exclude_mouse"].(bool)),
			"recording-include-keys":   boolToString(p["recording_include_keys"].(bool)),
			"create-recording-path":    boolToString(p["recording_auto_create_path"].(bool)),
		}
	}

	return conn
}

func validateConnectionKubernetes(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) == 0 {
		return diags
	}
	p := paramList[0].(map[string]interface{})

	intFields := map[string]string{
		"port":                p["port"].(string),
		"max_scrollback_size": p["max_scrollback_size"].(string),
	}
	for name, val := range intFields {
		if val != "" {
			if _, err := strconv.Atoi(val); err != nil {
				diags = append(diags, diag.Errorf("expected integer for parameter %s, got: %s", name, val)...)
			}
		}
	}

	return diags
}
