package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConnectionVNC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionVNCRead,
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
						"hostname":                   {Type: schema.TypeString, Computed: true},
						"port":                       {Type: schema.TypeString, Computed: true},
						"username":                   {Type: schema.TypeString, Computed: true},
						"password":                   {Type: schema.TypeString, Computed: true, Sensitive: true},
						"readonly":                   {Type: schema.TypeBool, Computed: true},
						"swap_red_blue":              {Type: schema.TypeBool, Computed: true},
						"cursor":                     {Type: schema.TypeString, Computed: true},
						"color_depth":                {Type: schema.TypeString, Computed: true},
						"clipboard_encoding":         {Type: schema.TypeString, Computed: true},
						"disable_copy":               {Type: schema.TypeBool, Computed: true},
						"disable_paste":              {Type: schema.TypeBool, Computed: true},
						"destination_host":           {Type: schema.TypeString, Computed: true},
						"destination_port":           {Type: schema.TypeString, Computed: true},
						"recording_path":             {Type: schema.TypeString, Computed: true},
						"recording_name":             {Type: schema.TypeString, Computed: true},
						"recording_exclude_output":   {Type: schema.TypeBool, Computed: true},
						"recording_exclude_mouse":    {Type: schema.TypeBool, Computed: true},
						"recording_include_keys":     {Type: schema.TypeBool, Computed: true},
						"recording_auto_create_path": {Type: schema.TypeBool, Computed: true},
						"sftp_enable":                {Type: schema.TypeBool, Computed: true},
						"sftp_root_directory":        {Type: schema.TypeString, Computed: true},
						"sftp_hostname":              {Type: schema.TypeString, Computed: true},
						"sftp_port":                  {Type: schema.TypeString, Computed: true},
						"sftp_host_key":              {Type: schema.TypeString, Computed: true},
						"sftp_username":              {Type: schema.TypeString, Computed: true},
						"sftp_password":              {Type: schema.TypeString, Computed: true, Sensitive: true},
						"sftp_private_key":           {Type: schema.TypeString, Computed: true, Sensitive: true},
						"sftp_passphrase":            {Type: schema.TypeString, Computed: true, Sensitive: true},
						"sftp_upload_directory":      {Type: schema.TypeString, Computed: true},
						"sftp_keepalive_interval":    {Type: schema.TypeString, Computed: true},
						"sftp_disable_file_download": {Type: schema.TypeBool, Computed: true},
						"sftp_disable_file_upload":   {Type: schema.TypeBool, Computed: true},
						"enable_audio":               {Type: schema.TypeBool, Computed: true},
						"audio_server_name":          {Type: schema.TypeString, Computed: true},
						"wol_send_packet":            {Type: schema.TypeBool, Computed: true},
						"wol_mac_address":            {Type: schema.TypeString, Computed: true},
						"wol_broadcast_address":      {Type: schema.TypeString, Computed: true},
						"wol_boot_wait_time":         {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceConnectionVNCRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Get("identifier").(string)

	conn, err := client.GetConnection(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read vnc connection %s: %w", id, err))
	}

	params, err := client.GetConnectionParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read vnc connection parameters %s: %w", id, err))
	}

	d.Set("name", conn.Name)
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

	d.SetId(id)

	return nil
}
