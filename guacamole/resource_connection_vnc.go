package guacamole

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleConnectionVNC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionVNCCreate,
		ReadContext:   resourceConnectionVNCRead,
		UpdateContext: resourceConnectionVNCUpdate,
		DeleteContext: resourceConnectionVNCDelete,
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
						"hostname":                   {Type: schema.TypeString, Required: true},
						"port":                       {Type: schema.TypeString, Optional: true, Computed: true},
						"username":                   {Type: schema.TypeString, Required: true},
						"password":                   {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"readonly":                   {Type: schema.TypeBool, Optional: true, Computed: true},
						"swap_red_blue":              {Type: schema.TypeBool, Optional: true, Computed: true},
						"cursor":                     {Type: schema.TypeString, Optional: true, Computed: true},
						"color_depth":                {Type: schema.TypeString, Optional: true, Computed: true},
						"clipboard_encoding":         {Type: schema.TypeString, Optional: true, Computed: true},
						"disable_copy":               {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_paste":              {Type: schema.TypeBool, Optional: true, Computed: true},
						"destination_host":           {Type: schema.TypeString, Optional: true, Computed: true},
						"destination_port":           {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_path":             {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_name":             {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_exclude_output":   {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_exclude_mouse":    {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_include_keys":     {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_auto_create_path": {Type: schema.TypeBool, Optional: true, Computed: true},
						"sftp_enable":                {Type: schema.TypeBool, Optional: true, Computed: true},
						"sftp_root_directory":        {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_hostname":              {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_port":                  {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_host_key":              {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_username":              {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_password":              {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"sftp_private_key":           {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"sftp_passphrase":            {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"sftp_upload_directory":      {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_keepalive_interval":    {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_disable_file_download": {Type: schema.TypeBool, Optional: true, Computed: true},
						"sftp_disable_file_upload":   {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_audio":               {Type: schema.TypeBool, Optional: true, Computed: true},
						"audio_server_name":          {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_send_packet":            {Type: schema.TypeBool, Optional: true, Computed: true},
						"wol_mac_address":            {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_broadcast_address":      {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_boot_wait_time":         {Type: schema.TypeString, Optional: true, Computed: true},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionVNCCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if check := validateConnectionVNC(d); check.HasError() {
		return check
	}

	conn := buildVNCConnection(d)
	created, err := client.CreateConnection(ctx, conn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("identifier", created.Identifier)
	d.SetId(created.Identifier)

	return resourceConnectionVNCRead(ctx, d, m)
}

func resourceConnectionVNCRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read vnc connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read vnc connection parameters %s: %w", id, err))
	}

	d.Set("name", conn.Name)
	d.Set("identifier", conn.Identifier)
	d.Set("parent_identifier", conn.ParentIdentifier)
	d.Set("protocol", conn.Protocol)
	d.Set("active_connections", conn.ActiveConnections)
	d.Set("attributes", readConnectionAttributes(conn.Attributes))
	d.Set("parameters", []interface{}{
		map[string]interface{}{
			"hostname":                   params["hostname"],
			"port":                       params["port"],
			"username":                   params["username"],
			"password":                   params["password"],
			"readonly":                   stringToBool(params["read-only"]),
			"swap_red_blue":              stringToBool(params["swap-red-blue"]),
			"cursor":                     params["cursor"],
			"color_depth":                params["color-depth"],
			"clipboard_encoding":         params["clipboard-encoding"],
			"disable_copy":               stringToBool(params["disable-copy"]),
			"disable_paste":              stringToBool(params["disable-paste"]),
			"destination_host":           params["dest-host"],
			"destination_port":           params["dest-port"],
			"recording_path":             params["recording-path"],
			"recording_name":             params["recording-name"],
			"recording_exclude_output":   stringToBool(params["recording-exclude-output"]),
			"recording_exclude_mouse":    stringToBool(params["recording-exclude-mouse"]),
			"recording_include_keys":     stringToBool(params["recording-include-keys"]),
			"recording_auto_create_path": stringToBool(params["create-recording-path"]),
			"sftp_enable":                stringToBool(params["enable-sftp"]),
			"sftp_root_directory":        params["sftp-root-directory"],
			"sftp_hostname":              params["sftp-hostname"],
			"sftp_port":                  params["sftp-port"],
			"sftp_host_key":              params["sftp-host-key"],
			"sftp_username":              params["sftp-username"],
			"sftp_password":              params["sftp-password"],
			"sftp_private_key":           params["sftp-private-key"],
			"sftp_passphrase":            params["sftp-passphrase"],
			"sftp_upload_directory":      params["sftp-upload-directory"],
			"sftp_keepalive_interval":    params["sftp-alive-interval"],
			"sftp_disable_file_download": stringToBool(params["sftp-disable-file-download"]),
			"sftp_disable_file_upload":   stringToBool(params["sftp-disable-file-upload"]),
			"enable_audio":               stringToBool(params["enable-audio"]),
			"audio_server_name":          params["audio-server-name"],
			"wol_send_packet":            stringToBool(params["wol-send-packet"]),
			"wol_mac_address":            params["wol-mac-addr"],
			"wol_broadcast_address":      params["wol-broadcast-addr"],
			"wol_boot_wait_time":         params["wol-boot-wait-time"],
		},
	})

	return nil
}

func resourceConnectionVNCUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("name", "parent_identifier", "attributes", "parameters") {
		if check := validateConnectionVNC(d); check.HasError() {
			return check
		}
		conn := buildVNCConnection(d)
		if err := client.UpdateConnection(ctx, d.Id(), conn); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceConnectionVNCRead(ctx, d, m)
}

func resourceConnectionVNCDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteConnection(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func buildVNCConnection(d *schema.ResourceData) guacamole.Connection {
	conn := guacamole.Connection{
		Name:             d.Get("name").(string),
		Identifier:       d.Get("identifier").(string),
		ParentIdentifier: d.Get("parent_identifier").(string),
		Protocol:         "vnc",
	}

	conn.Attributes = buildConnectionAttributes(d.Get("attributes").([]interface{}))

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) > 0 {
		p := paramList[0].(map[string]interface{})
		conn.Parameters = map[string]string{
			"hostname":           p["hostname"].(string),
			"port":               p["port"].(string),
			"username":           p["username"].(string),
			"password":           p["password"].(string),
			"read-only":          boolToString(p["readonly"].(bool)),
			"swap-red-blue":      boolToString(p["swap_red_blue"].(bool)),
			"cursor":             p["cursor"].(string),
			"color-depth":        p["color_depth"].(string),
			"clipboard-encoding": p["clipboard_encoding"].(string),
			"disable-copy":       boolToString(p["disable_copy"].(bool)),
			"disable-paste":      boolToString(p["disable_paste"].(bool)),
			"dest-host":          p["destination_host"].(string),
			"dest-port":          p["destination_port"].(string),
			"recording-path":     p["recording_path"].(string),
			"recording-name":     p["recording_name"].(string),
			"recording-exclude-output": boolToString(p["recording_exclude_output"].(bool)),
			"recording-exclude-mouse":  boolToString(p["recording_exclude_mouse"].(bool)),
			"recording-include-keys":   boolToString(p["recording_include_keys"].(bool)),
			"create-recording-path":    boolToString(p["recording_auto_create_path"].(bool)),
			"enable-sftp":              boolToString(p["sftp_enable"].(bool)),
			"sftp-root-directory":      p["sftp_root_directory"].(string),
			"sftp-hostname":            p["sftp_hostname"].(string),
			"sftp-port":                p["sftp_port"].(string),
			"sftp-host-key":            p["sftp_host_key"].(string),
			"sftp-username":            p["sftp_username"].(string),
			"sftp-password":            p["sftp_password"].(string),
			"sftp-private-key":         p["sftp_private_key"].(string),
			"sftp-passphrase":          p["sftp_passphrase"].(string),
			"sftp-upload-directory":    p["sftp_upload_directory"].(string),
			"sftp-alive-interval":      p["sftp_keepalive_interval"].(string),
			"sftp-disable-file-download": boolToString(p["sftp_disable_file_download"].(bool)),
			"sftp-disable-file-upload":   boolToString(p["sftp_disable_file_upload"].(bool)),
			"enable-audio":             boolToString(p["enable_audio"].(bool)),
			"audio-server-name":        p["audio_server_name"].(string),
			"wol-send-packet":          boolToString(p["wol_send_packet"].(bool)),
			"wol-mac-addr":             p["wol_mac_address"].(string),
			"wol-broadcast-addr":       p["wol_broadcast_address"].(string),
			"wol-boot-wait-time":       p["wol_boot_wait_time"].(string),
		}
	}

	return conn
}

func validateConnectionVNC(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) == 0 {
		return diags
	}
	p := paramList[0].(map[string]interface{})

	intFields := map[string]string{
		"port":                    p["port"].(string),
		"destination_port":        p["destination_port"].(string),
		"sftp_port":               p["sftp_port"].(string),
		"sftp_keepalive_interval": p["sftp_keepalive_interval"].(string),
		"wol_boot_wait_time":      p["wol_boot_wait_time"].(string),
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
