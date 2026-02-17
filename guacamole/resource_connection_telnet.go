package guacamole

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleConnectionTelnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionTelnetCreate,
		ReadContext:   resourceConnectionTelnetRead,
		UpdateContext: resourceConnectionTelnetUpdate,
		DeleteContext: resourceConnectionTelnetDelete,
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
						"username":                    {Type: schema.TypeString, Required: true},
						"password":                    {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"username_regex":              {Type: schema.TypeString, Optional: true, Computed: true},
						"password_regex":              {Type: schema.TypeString, Optional: true, Computed: true},
						"login_success_regex":         {Type: schema.TypeString, Optional: true, Computed: true},
						"login_failure_regex":         {Type: schema.TypeString, Optional: true, Computed: true},
						"color_scheme":                {Type: schema.TypeString, Optional: true, Computed: true},
						"font_name":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"font_size":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"max_scrollback_size":         {Type: schema.TypeString, Optional: true, Computed: true},
						"readonly":                    {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_copy":                {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_paste":               {Type: schema.TypeBool, Optional: true, Computed: true},
						"backspace":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"terminal_type":               {Type: schema.TypeString, Optional: true, Computed: true},
						"typescript_path":             {Type: schema.TypeString, Optional: true, Computed: true},
						"typescript_name":             {Type: schema.TypeString, Optional: true, Computed: true},
						"typescript_auto_create_path": {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_path":              {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_name":              {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_exclude_output":    {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_exclude_mouse":     {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_include_keys":      {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_auto_create_path":  {Type: schema.TypeBool, Optional: true, Computed: true},
						"wol_send_packet":             {Type: schema.TypeBool, Optional: true, Computed: true},
						"wol_mac_address":             {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_broadcast_address":       {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_boot_wait_time":          {Type: schema.TypeString, Optional: true, Computed: true},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionTelnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if check := validateConnectionTelnet(d); check.HasError() {
		return check
	}

	conn := buildTelnetConnection(d)
	created, err := client.CreateConnection(ctx, conn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("identifier", created.Identifier)
	d.SetId(created.Identifier)

	return resourceConnectionTelnetRead(ctx, d, m)
}

func resourceConnectionTelnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read telnet connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read telnet connection parameters %s: %w", id, err))
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

	return nil
}

func resourceConnectionTelnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("name", "parent_identifier", "attributes", "parameters") {
		if check := validateConnectionTelnet(d); check.HasError() {
			return check
		}
		conn := buildTelnetConnection(d)
		if err := client.UpdateConnection(ctx, d.Id(), conn); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceConnectionTelnetRead(ctx, d, m)
}

func resourceConnectionTelnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteConnection(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func buildTelnetConnection(d *schema.ResourceData) guacamole.Connection {
	conn := guacamole.Connection{
		Name:             d.Get("name").(string),
		Identifier:       d.Get("identifier").(string),
		ParentIdentifier: d.Get("parent_identifier").(string),
		Protocol:         "telnet",
	}

	conn.Attributes = buildConnectionAttributes(d.Get("attributes").([]interface{}))

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) > 0 {
		p := paramList[0].(map[string]interface{})
		conn.Parameters = map[string]string{
			"hostname":               p["hostname"].(string),
			"port":                   p["port"].(string),
			"username":               p["username"].(string),
			"password":               p["password"].(string),
			"username-regex":         p["username_regex"].(string),
			"password-regex":         p["password_regex"].(string),
			"login-success-regex":    p["login_success_regex"].(string),
			"login-failure-regex":    p["login_failure_regex"].(string),
			"color-scheme":           p["color_scheme"].(string),
			"font-name":              p["font_name"].(string),
			"font-size":              p["font_size"].(string),
			"scrollback":             p["max_scrollback_size"].(string),
			"read-only":              boolToString(p["readonly"].(bool)),
			"disable-copy":           boolToString(p["disable_copy"].(bool)),
			"disable-paste":          boolToString(p["disable_paste"].(bool)),
			"backspace":              p["backspace"].(string),
			"terminal-type":          p["terminal_type"].(string),
			"typescript-path":        p["typescript_path"].(string),
			"typescript-name":        p["typescript_name"].(string),
			"create-typescript-path": boolToString(p["typescript_auto_create_path"].(bool)),
			"recording-path":         p["recording_path"].(string),
			"recording-name":         p["recording_name"].(string),
			"recording-exclude-output": boolToString(p["recording_exclude_output"].(bool)),
			"recording-exclude-mouse":  boolToString(p["recording_exclude_mouse"].(bool)),
			"recording-include-keys":   boolToString(p["recording_include_keys"].(bool)),
			"create-recording-path":    boolToString(p["recording_auto_create_path"].(bool)),
			"wol-send-packet":          boolToString(p["wol_send_packet"].(bool)),
			"wol-mac-addr":             p["wol_mac_address"].(string),
			"wol-broadcast-addr":       p["wol_broadcast_address"].(string),
			"wol-boot-wait-time":       p["wol_boot_wait_time"].(string),
		}
	}

	return conn
}

func validateConnectionTelnet(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) == 0 {
		return diags
	}
	p := paramList[0].(map[string]interface{})

	intFields := map[string]string{
		"port":             p["port"].(string),
		"wol_boot_wait_time": p["wol_boot_wait_time"].(string),
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
