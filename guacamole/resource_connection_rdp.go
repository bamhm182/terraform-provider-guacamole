package guacamole

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleConnectionRDP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionRDPCreate,
		ReadContext:   resourceConnectionRDPRead,
		UpdateContext: resourceConnectionRDPUpdate,
		DeleteContext: resourceConnectionRDPDelete,
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
						"hostname":                     {Type: schema.TypeString, Required: true},
						"port":                         {Type: schema.TypeString, Optional: true, Computed: true},
						"username":                     {Type: schema.TypeString, Required: true},
						"password":                     {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"domain":                       {Type: schema.TypeString, Optional: true, Computed: true},
						"security_mode":                {Type: schema.TypeString, Optional: true, Computed: true},
						"disable_authentication":       {Type: schema.TypeBool, Optional: true, Computed: true},
						"ignore_cert":                  {Type: schema.TypeBool, Optional: true, Computed: true},
						"gateway_hostname":             {Type: schema.TypeString, Optional: true, Computed: true},
						"gateway_port":                 {Type: schema.TypeString, Optional: true, Computed: true},
						"gateway_username":             {Type: schema.TypeString, Optional: true, Computed: true},
						"gateway_password":             {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"gateway_domain":               {Type: schema.TypeString, Optional: true, Computed: true},
						"initial_program":              {Type: schema.TypeString, Optional: true, Computed: true},
						"client_name":                  {Type: schema.TypeString, Optional: true, Computed: true},
						"keyboard_layout":              {Type: schema.TypeString, Optional: true, Computed: true},
						"timezone":                     {Type: schema.TypeString, Optional: true, Computed: true},
						"administrator_console":        {Type: schema.TypeBool, Optional: true, Computed: true},
						"width":                        {Type: schema.TypeString, Optional: true, Computed: true},
						"height":                       {Type: schema.TypeString, Optional: true, Computed: true},
						"dpi":                          {Type: schema.TypeString, Optional: true, Computed: true},
						"color_depth":                  {Type: schema.TypeString, Optional: true, Computed: true},
						"resize_method":                {Type: schema.TypeString, Optional: true, Computed: true},
						"readonly":                     {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_copy":                 {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_paste":                {Type: schema.TypeBool, Optional: true, Computed: true},
						"console_audio":                {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_audio":                {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_audio_input":           {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_printing":              {Type: schema.TypeBool, Optional: true, Computed: true},
						"printer_name":                 {Type: schema.TypeString, Optional: true, Computed: true},
						"enable_drive":                 {Type: schema.TypeBool, Optional: true, Computed: true},
						"drive_name":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"disable_file_download":        {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_file_upload":          {Type: schema.TypeBool, Optional: true, Computed: true},
						"drive_path":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"create_drive_path":            {Type: schema.TypeBool, Optional: true, Computed: true},
						"static_channels":              {Type: schema.TypeString, Optional: true, Computed: true},
						"enable_wallpaper":             {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_theming":               {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_font_smoothing":        {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_full_window_drag":      {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_desktop_composition":   {Type: schema.TypeBool, Optional: true, Computed: true},
						"enable_menu_animations":       {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_bitmap_caching":       {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_offscreen_caching":    {Type: schema.TypeBool, Optional: true, Computed: true},
						"disable_glyph_caching":        {Type: schema.TypeBool, Optional: true, Computed: true},
						"remote_app":                   {Type: schema.TypeString, Optional: true, Computed: true},
						"remote_app_working_directory": {Type: schema.TypeString, Optional: true, Computed: true},
						"remote_app_parameters":        {Type: schema.TypeString, Optional: true, Computed: true},
						"preconnection_id":             {Type: schema.TypeString, Optional: true, Computed: true},
						"preconnection_blob":           {Type: schema.TypeString, Optional: true, Computed: true},
						"load_balance_info":            {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_path":               {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_name":               {Type: schema.TypeString, Optional: true, Computed: true},
						"recording_exclude_output":     {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_exclude_mouse":      {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_include_keys":       {Type: schema.TypeBool, Optional: true, Computed: true},
						"recording_auto_create_path":   {Type: schema.TypeBool, Optional: true, Computed: true},
						"sftp_enable":                  {Type: schema.TypeBool, Optional: true, Computed: true},
						"sftp_root_directory":          {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_hostname":                {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_port":                    {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_host_key":                {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_username":                {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_password":                {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"sftp_private_key":             {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"sftp_passphrase":              {Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true},
						"sftp_upload_directory":        {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_keepalive_interval":      {Type: schema.TypeString, Optional: true, Computed: true},
						"sftp_disable_file_download":   {Type: schema.TypeBool, Optional: true, Computed: true},
						"sftp_disable_file_upload":     {Type: schema.TypeBool, Optional: true, Computed: true},
						"wol_send_packet":              {Type: schema.TypeBool, Optional: true, Computed: true},
						"wol_mac_address":              {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_broadcast_address":        {Type: schema.TypeString, Optional: true, Computed: true},
						"wol_boot_wait_time":           {Type: schema.TypeString, Optional: true, Computed: true},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionRDPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if check := validateConnectionRDP(d); check.HasError() {
		return check
	}

	conn := buildRDPConnection(d)
	created, err := client.CreateConnection(ctx, conn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("identifier", created.Identifier)
	d.SetId(created.Identifier)

	return resourceConnectionRDPRead(ctx, d, m)
}

func resourceConnectionRDPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read rdp connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read rdp connection parameters %s: %w", id, err))
	}

	d.Set("name", conn.Name)
	d.Set("identifier", conn.Identifier)
	d.Set("parent_identifier", conn.ParentIdentifier)
	d.Set("protocol", conn.Protocol)
	d.Set("active_connections", conn.ActiveConnections)
	d.Set("attributes", readConnectionAttributes(conn.Attributes))
	d.Set("parameters", []interface{}{
		map[string]interface{}{
			"hostname":                     params["hostname"],
			"port":                         params["port"],
			"username":                     params["username"],
			"password":                     params["password"],
			"domain":                       params["domain"],
			"security_mode":                params["security"],
			"disable_authentication":       stringToBool(params["disable-auth"]),
			"ignore_cert":                  stringToBool(params["ignore-cert"]),
			"gateway_hostname":             params["gateway-hostname"],
			"gateway_port":                 params["gateway-port"],
			"gateway_username":             params["gateway-username"],
			"gateway_password":             params["gateway-password"],
			"gateway_domain":               params["gateway-domain"],
			"initial_program":              params["initial-program"],
			"client_name":                  params["client-name"],
			"keyboard_layout":              params["server-layout"],
			"timezone":                     params["timezone"],
			"administrator_console":        stringToBool(params["console"]),
			"width":                        params["width"],
			"height":                       params["height"],
			"dpi":                          params["resolution"],
			"color_depth":                  params["color-depth"],
			"resize_method":                params["resize-method"],
			"readonly":                     stringToBool(params["read-only"]),
			"disable_copy":                 stringToBool(params["disable-copy"]),
			"disable_paste":                stringToBool(params["disable-paste"]),
			"console_audio":                stringToBool(params["console-audio"]),
			"disable_audio":                stringToBool(params["disable-audio"]),
			"enable_audio_input":           stringToBool(params["enable-audio-input"]),
			"enable_printing":              stringToBool(params["enable-printing"]),
			"printer_name":                 params["printer-name"],
			"enable_drive":                 stringToBool(params["enable-drive"]),
			"drive_name":                   params["drive-name"],
			"disable_file_download":        stringToBool(params["disable-file-download"]),
			"disable_file_upload":          stringToBool(params["disable-file-upload"]),
			"drive_path":                   params["drive-path"],
			"create_drive_path":            stringToBool(params["create-drive-path"]),
			"static_channels":              params["static-channels"],
			"enable_wallpaper":             stringToBool(params["enable-wallpaper"]),
			"enable_theming":               stringToBool(params["enable-theming"]),
			"enable_font_smoothing":        stringToBool(params["enable-font-smoothing"]),
			"enable_full_window_drag":      stringToBool(params["enable-full-window-drag"]),
			"enable_desktop_composition":   stringToBool(params["enable-desktop-composition"]),
			"enable_menu_animations":       stringToBool(params["enable-menu-animations"]),
			"disable_bitmap_caching":       stringToBool(params["disable-bitmap-caching"]),
			"disable_offscreen_caching":    stringToBool(params["disable-offscreen-caching"]),
			"disable_glyph_caching":        stringToBool(params["disable-glyph-caching"]),
			"remote_app":                   params["remote-app"],
			"remote_app_working_directory": params["remote-app-dir"],
			"remote_app_parameters":        params["remote-app-args"],
			"preconnection_id":             params["preconnection-id"],
			"preconnection_blob":           params["preconnection-blob"],
			"load_balance_info":            params["load-balance-info"],
			"recording_path":               params["recording-path"],
			"recording_name":               params["recording-name"],
			"recording_exclude_output":     stringToBool(params["recording-exclude-output"]),
			"recording_exclude_mouse":      stringToBool(params["recording-exclude-mouse"]),
			"recording_include_keys":       stringToBool(params["recording-include-keys"]),
			"recording_auto_create_path":   stringToBool(params["create-recording-path"]),
			"sftp_enable":                  stringToBool(params["enable-sftp"]),
			"sftp_root_directory":          params["sftp-root-directory"],
			"sftp_hostname":                params["sftp-hostname"],
			"sftp_port":                    params["sftp-port"],
			"sftp_host_key":                params["sftp-host-key"],
			"sftp_username":                params["sftp-username"],
			"sftp_password":                params["sftp-password"],
			"sftp_private_key":             params["sftp-private-key"],
			"sftp_passphrase":              params["sftp-passphrase"],
			"sftp_upload_directory":        params["sftp-upload-directory"],
			"sftp_keepalive_interval":      params["sftp-alive-interval"],
			"sftp_disable_file_download":   stringToBool(params["sftp-disable-file-download"]),
			"sftp_disable_file_upload":     stringToBool(params["sftp-disable-file-upload"]),
			"wol_send_packet":              stringToBool(params["wol-send-packet"]),
			"wol_mac_address":              params["wol-mac-addr"],
			"wol_broadcast_address":        params["wol-broadcast-addr"],
			"wol_boot_wait_time":           params["wol-boot-wait-time"],
		},
	})

	return nil
}

func resourceConnectionRDPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("name", "parent_identifier", "attributes", "parameters") {
		if check := validateConnectionRDP(d); check.HasError() {
			return check
		}
		conn := buildRDPConnection(d)
		if err := client.UpdateConnection(ctx, d.Id(), conn); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceConnectionRDPRead(ctx, d, m)
}

func resourceConnectionRDPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteConnection(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func buildRDPConnection(d *schema.ResourceData) guacamole.Connection {
	conn := guacamole.Connection{
		Name:             d.Get("name").(string),
		Identifier:       d.Get("identifier").(string),
		ParentIdentifier: d.Get("parent_identifier").(string),
		Protocol:         "rdp",
	}

	conn.Attributes = buildConnectionAttributes(d.Get("attributes").([]interface{}))

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) > 0 {
		p := paramList[0].(map[string]interface{})
		conn.Parameters = map[string]string{
			"hostname":         p["hostname"].(string),
			"port":             p["port"].(string),
			"username":         p["username"].(string),
			"password":         p["password"].(string),
			"domain":           p["domain"].(string),
			"security":         p["security_mode"].(string),
			"disable-auth":     boolToString(p["disable_authentication"].(bool)),
			"ignore-cert":      boolToString(p["ignore_cert"].(bool)),
			"gateway-hostname": p["gateway_hostname"].(string),
			"gateway-port":     p["gateway_port"].(string),
			"gateway-username": p["gateway_username"].(string),
			"gateway-password": p["gateway_password"].(string),
			"gateway-domain":   p["gateway_domain"].(string),
			"initial-program":  p["initial_program"].(string),
			"client-name":      p["client_name"].(string),
			"server-layout":    p["keyboard_layout"].(string),
			"timezone":         p["timezone"].(string),
			"console":          boolToString(p["administrator_console"].(bool)),
			"width":            p["width"].(string),
			"height":           p["height"].(string),
			"resolution":       p["dpi"].(string),
			"color-depth":      p["color_depth"].(string),
			"resize-method":    p["resize_method"].(string),
			"read-only":        boolToString(p["readonly"].(bool)),
			"disable-copy":     boolToString(p["disable_copy"].(bool)),
			"disable-paste":    boolToString(p["disable_paste"].(bool)),
			"console-audio":    boolToString(p["console_audio"].(bool)),
			"disable-audio":    boolToString(p["disable_audio"].(bool)),
			"enable-audio-input": boolToString(p["enable_audio_input"].(bool)),
			"enable-printing":  boolToString(p["enable_printing"].(bool)),
			"printer-name":     p["printer_name"].(string),
			"enable-drive":     boolToString(p["enable_drive"].(bool)),
			"drive-name":       p["drive_name"].(string),
			"disable-file-download": boolToString(p["disable_file_download"].(bool)),
			"disable-file-upload":   boolToString(p["disable_file_upload"].(bool)),
			"drive-path":       p["drive_path"].(string),
			"create-drive-path": boolToString(p["create_drive_path"].(bool)),
			"static-channels":  p["static_channels"].(string),
			"enable-wallpaper": boolToString(p["enable_wallpaper"].(bool)),
			"enable-theming":   boolToString(p["enable_theming"].(bool)),
			"enable-font-smoothing":      boolToString(p["enable_font_smoothing"].(bool)),
			"enable-full-window-drag":    boolToString(p["enable_full_window_drag"].(bool)),
			"enable-desktop-composition": boolToString(p["enable_desktop_composition"].(bool)),
			"enable-menu-animations":     boolToString(p["enable_menu_animations"].(bool)),
			"disable-bitmap-caching":     boolToString(p["disable_bitmap_caching"].(bool)),
			"disable-offscreen-caching":  boolToString(p["disable_offscreen_caching"].(bool)),
			"disable-glyph-caching":      boolToString(p["disable_glyph_caching"].(bool)),
			"remote-app":       p["remote_app"].(string),
			"remote-app-dir":   p["remote_app_working_directory"].(string),
			"remote-app-args":  p["remote_app_parameters"].(string),
			"preconnection-id":   p["preconnection_id"].(string),
			"preconnection-blob": p["preconnection_blob"].(string),
			"load-balance-info": p["load_balance_info"].(string),
			"recording-path":   p["recording_path"].(string),
			"recording-name":   p["recording_name"].(string),
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
			"wol-send-packet":  boolToString(p["wol_send_packet"].(bool)),
			"wol-mac-addr":     p["wol_mac_address"].(string),
			"wol-broadcast-addr": p["wol_broadcast_address"].(string),
			"wol-boot-wait-time": p["wol_boot_wait_time"].(string),
		}
	}

	return conn
}

func validateConnectionRDP(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	paramList := d.Get("parameters").([]interface{})
	if len(paramList) == 0 {
		return diags
	}
	p := paramList[0].(map[string]interface{})

	intFields := map[string]string{
		"port":                    p["port"].(string),
		"gateway_port":            p["gateway_port"].(string),
		"width":                   p["width"].(string),
		"height":                  p["height"].(string),
		"dpi":                     p["dpi"].(string),
		"preconnection_id":        p["preconnection_id"].(string),
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

	tz := p["timezone"].(string)
	if tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			diags = append(diags, diag.Errorf("invalid timezone: %s", tz)...)
		}
	}

	return diags
}
