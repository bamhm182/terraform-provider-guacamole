package guacamole

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Provider returns the Guacamole Terraform provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GUACAMOLE_URL", nil),
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"password"},
				DefaultFunc:  schema.EnvDefaultFunc("GUACAMOLE_USERNAME", nil),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"username"},
				AtLeastOneOf: []string{"password", "token"},
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GUACAMOLE_PASSWORD", nil),
			},
			"token": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"password", "token"},
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GUACAMOLE_TOKEN", nil),
			},
			"data_source": {
				Type:             schema.TypeString,
				Optional:         true,
				RequiredWith:     []string{"token"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"postgresql", "mysql"}, true)),
				DefaultFunc:      schema.EnvDefaultFunc("GUACAMOLE_DATA_SOURCE", nil),
			},
			"disable_tls_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("GUACAMOLE_DISABLE_TLS", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"guacamole_user":                  guacamoleUser(),
			"guacamole_user_group":            guacamoleUserGroup(),
			"guacamole_connection_ssh":        guacamoleConnectionSSH(),
			"guacamole_connection_telnet":     guacamoleConnectionTelnet(),
			"guacamole_connection_rdp":        guacamoleConnectionRDP(),
			"guacamole_connection_vnc":        guacamoleConnectionVNC(),
			"guacamole_connection_kubernetes": guacamoleConnectionKubernetes(),
			"guacamole_connection_group":      guacamoleConnectionGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"guacamole_user":                  dataSourceUser(),
			"guacamole_user_group":            dataSourceUserGroup(),
			"guacamole_connection_ssh":        dataSourceConnectionSSH(),
			"guacamole_connection_telnet":     dataSourceConnectionTelnet(),
			"guacamole_connection_rdp":        dataSourceConnectionRDP(),
			"guacamole_connection_vnc":        dataSourceConnectionVNC(),
			"guacamole_connection_kubernetes": dataSourceConnectionKubernetes(),
			"guacamole_connection_group":      dataSourceConnectionGroup(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	rawURL := strings.TrimRight(d.Get("url").(string), "/")
	disableTLS := d.Get("disable_tls_verification").(bool)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	if disableTLS {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
	}

	token := d.Get("token").(string)
	if token != "" {
		dataSource := d.Get("data_source").(string)
		client := guacamole.NewClientWithToken(rawURL, token, dataSource, httpClient)
		return client, diags
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	client := guacamole.NewClientWithHTTPClient(rawURL, httpClient)
	if err := client.Authenticate(ctx, username, password); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to authenticate with Guacamole",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return client, diags
}
