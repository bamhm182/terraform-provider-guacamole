package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConnectionRDP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionRDPRead,
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
						"hostname":                     {Type: schema.TypeString, Computed: true},
						"port":                         {Type: schema.TypeString, Computed: true},
						"username":                     {Type: schema.TypeString, Computed: true},
						"password":                     {Type: schema.TypeString, Computed: true, Sensitive: true},
						"domain":                       {Type: schema.TypeString, Computed: true},
						"security_mode":                {Type: schema.TypeString, Computed: true},
						"disable_authentication":       {Type: schema.TypeBool, Computed: true},
						"ignore_cert":                  {Type: schema.TypeBool, Computed: true},
						"gateway_hostname":             {Type: schema.TypeString, Computed: true},
						"gateway_port":                 {Type: schema.TypeString, Computed: true},
						"gateway_username":             {Type: schema.TypeString, Computed: true},
						"gateway_password":             {Type: schema.TypeString, Computed: true, Sensitive: true},
						"gateway_domain":               {Type: schema.TypeString, Computed: true},
						"initial_program":              {Type: schema.TypeString, Computed: true},
						"client_name":                  {Type: schema.TypeString, Computed: true},
						"keyboard_layout":              {Type: schema.TypeString, Computed: true},
						"timezone":                     {Type: schema.TypeString, Computed: true},
						"administrator_console":        {Type: schema.TypeBool, Computed: true},
						"width":                        {Type: schema.TypeString, Computed: true},
						"height":                       {Type: schema.TypeString, Computed: true},
						"dpi":                          {Type: schema.TypeString, Computed: true},
						"color_depth":                  {Type: schema.TypeString, Computed: true},
						"resize_method":                {Type: schema.TypeString, Computed: true},
						"readonly":                     {Type: schema.TypeBool, Computed: true},
						"disable_copy":                 {Type: schema.TypeBool, Computed: true},
						"disable_paste":                {Type: schema.TypeBool, Computed: true},
						"console_audio":                {Type: schema.TypeBool, Computed: true},
						"disable_audio":                {Type: schema.TypeBool, Computed: true},
						"enable_audio_input":           {Type: schema.TypeBool, Computed: true},
						"enable_printing":              {Type: schema.TypeBool, Computed: true},
						"printer_name":                 {Type: schema.TypeString, Computed: true},
						"enable_drive":                 {Type: schema.TypeBool, Computed: true},
						"drive_name":                   {Type: schema.TypeString, Computed: true},
						"disable_file_download":        {Type: schema.TypeBool, Computed: true},
						"disable_file_upload":          {Type: schema.TypeBool, Computed: true},
						"drive_path":                   {Type: schema.TypeString, Computed: true},
						"create_drive_path":            {Type: schema.TypeBool, Computed: true},
						"static_channels":              {Type: schema.TypeString, Computed: true},
						"enable_wallpaper":             {Type: schema.TypeBool, Computed: true},
						"enable_theming":               {Type: schema.TypeBool, Computed: true},
						"enable_font_smoothing":        {Type: schema.TypeBool, Computed: true},
						"enable_full_window_drag":      {Type: schema.TypeBool, Computed: true},
						"enable_desktop_composition":   {Type: schema.TypeBool, Computed: true},
						"enable_menu_animations":       {Type: schema.TypeBool, Computed: true},
						"disable_bitmap_caching":       {Type: schema.TypeBool, Computed: true},
						"disable_offscreen_caching":    {Type: schema.TypeBool, Computed: true},
						"disable_glyph_caching":        {Type: schema.TypeBool, Computed: true},
						"remote_app":                   {Type: schema.TypeString, Computed: true},
						"remote_app_working_directory": {Type: schema.TypeString, Computed: true},
						"remote_app_parameters":        {Type: schema.TypeString, Computed: true},
						"preconnection_id":             {Type: schema.TypeString, Computed: true},
						"preconnection_blob":           {Type: schema.TypeString, Computed: true},
						"load_balance_info":            {Type: schema.TypeString, Computed: true},
						"recording_path":               {Type: schema.TypeString, Computed: true},
						"recording_name":               {Type: schema.TypeString, Computed: true},
						"recording_exclude_output":     {Type: schema.TypeBool, Computed: true},
						"recording_exclude_mouse":      {Type: schema.TypeBool, Computed: true},
						"recording_include_keys":       {Type: schema.TypeBool, Computed: true},
						"recording_auto_create_path":   {Type: schema.TypeBool, Computed: true},
						"sftp_enable":                  {Type: schema.TypeBool, Computed: true},
						"sftp_root_directory":          {Type: schema.TypeString, Computed: true},
						"sftp_hostname":                {Type: schema.TypeString, Computed: true},
						"sftp_port":                    {Type: schema.TypeString, Computed: true},
						"sftp_host_key":                {Type: schema.TypeString, Computed: true},
						"sftp_username":                {Type: schema.TypeString, Computed: true},
						"sftp_password":                {Type: schema.TypeString, Computed: true, Sensitive: true},
						"sftp_private_key":             {Type: schema.TypeString, Computed: true, Sensitive: true},
						"sftp_passphrase":              {Type: schema.TypeString, Computed: true, Sensitive: true},
						"sftp_upload_directory":        {Type: schema.TypeString, Computed: true},
						"sftp_keepalive_interval":      {Type: schema.TypeString, Computed: true},
						"sftp_disable_file_download":   {Type: schema.TypeBool, Computed: true},
						"sftp_disable_file_upload":     {Type: schema.TypeBool, Computed: true},
						"wol_send_packet":              {Type: schema.TypeBool, Computed: true},
						"wol_mac_address":              {Type: schema.TypeString, Computed: true},
						"wol_broadcast_address":        {Type: schema.TypeString, Computed: true},
						"wol_boot_wait_time":           {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceConnectionRDPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Get("identifier").(string)

	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read rdp connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read rdp connection parameters %s: %w", id, err))
	}

	d.Set("name", conn.Name)
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

	d.SetId(id)

	return nil
}
