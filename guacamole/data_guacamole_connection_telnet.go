package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConnectionTelnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionTelnetRead,
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
						"username":                    {Type: schema.TypeString, Computed: true},
						"password":                    {Type: schema.TypeString, Computed: true, Sensitive: true},
						"username_regex":              {Type: schema.TypeString, Computed: true},
						"password_regex":              {Type: schema.TypeString, Computed: true},
						"login_success_regex":         {Type: schema.TypeString, Computed: true},
						"login_failure_regex":         {Type: schema.TypeString, Computed: true},
						"color_scheme":                {Type: schema.TypeString, Computed: true},
						"font_name":                   {Type: schema.TypeString, Computed: true},
						"font_size":                   {Type: schema.TypeString, Computed: true},
						"max_scrollback_size":         {Type: schema.TypeString, Computed: true},
						"readonly":                    {Type: schema.TypeBool, Computed: true},
						"disable_copy":                {Type: schema.TypeBool, Computed: true},
						"disable_paste":               {Type: schema.TypeBool, Computed: true},
						"backspace":                   {Type: schema.TypeString, Computed: true},
						"terminal_type":               {Type: schema.TypeString, Computed: true},
						"typescript_path":             {Type: schema.TypeString, Computed: true},
						"typescript_name":             {Type: schema.TypeString, Computed: true},
						"typescript_auto_create_path": {Type: schema.TypeBool, Computed: true},
						"recording_path":              {Type: schema.TypeString, Computed: true},
						"recording_name":              {Type: schema.TypeString, Computed: true},
						"recording_exclude_output":    {Type: schema.TypeBool, Computed: true},
						"recording_exclude_mouse":     {Type: schema.TypeBool, Computed: true},
						"recording_include_keys":      {Type: schema.TypeBool, Computed: true},
						"recording_auto_create_path":  {Type: schema.TypeBool, Computed: true},
						"wol_send_packet":             {Type: schema.TypeBool, Computed: true},
						"wol_mac_address":             {Type: schema.TypeString, Computed: true},
						"wol_broadcast_address":       {Type: schema.TypeString, Computed: true},
						"wol_boot_wait_time":          {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceConnectionTelnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Get("identifier").(string)

	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read telnet connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read telnet connection parameters %s: %w", id, err))
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
			"username":                    params["username"],
			"password":                    params["password"],
			"username_regex":              params["username-regex"],
			"password_regex":              params["password-regex"],
			"login_success_regex":         params["login-success-regex"],
			"login_failure_regex":         params["login-failure-regex"],
			"color_scheme":                params["color-scheme"],
			"font_name":                   params["font-name"],
			"font_size":                   params["font-size"],
			"max_scrollback_size":         params["scrollback"],
			"readonly":                    stringToBool(params["read-only"]),
			"disable_copy":                stringToBool(params["disable-copy"]),
			"disable_paste":               stringToBool(params["disable-paste"]),
			"backspace":                   params["backspace"],
			"terminal_type":               params["terminal-type"],
			"typescript_path":             params["typescript-path"],
			"typescript_name":             params["typescript-name"],
			"typescript_auto_create_path": stringToBool(params["create-typescript-path"]),
			"recording_path":              params["recording-path"],
			"recording_name":              params["recording-name"],
			"recording_exclude_output":    stringToBool(params["recording-exclude-output"]),
			"recording_exclude_mouse":     stringToBool(params["recording-exclude-mouse"]),
			"recording_include_keys":      stringToBool(params["recording-include-keys"]),
			"recording_auto_create_path":  stringToBool(params["create-recording-path"]),
			"wol_send_packet":             stringToBool(params["wol-send-packet"]),
			"wol_mac_address":             params["wol-mac-addr"],
			"wol_broadcast_address":       params["wol-broadcast-addr"],
			"wol_boot_wait_time":          params["wol-boot-wait-time"],
		},
	})

	d.SetId(id)

	return nil
}
