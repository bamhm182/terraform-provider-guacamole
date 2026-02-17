package guacamole

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "Username of guacamole user",
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password of guacamole user",
				Optional:    true,
				Sensitive:   true,
			},
			"last_active": {
				Type:        schema.TypeString,
				Description: "Epoch time string of last user activity",
				Computed:    true,
			},
			"attributes": {
				Type:        schema.TypeList,
				Description: "Attributes of guacamole user",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organizational_role": {
							Type:        schema.TypeString,
							Description: "Organizational role of user",
							Optional:    true,
							Computed:    true,
						},
						"full_name": {
							Type:        schema.TypeString,
							Description: "Full name of user",
							Optional:    true,
							Computed:    true,
						},
						"email": {
							Type:        schema.TypeString,
							Description: "Email of user",
							Optional:    true,
							Computed:    true,
						},
						"expired": {
							Type:        schema.TypeBool,
							Description: "Whether the user is expired",
							Optional:    true,
						},
						"timezone": {
							Type:        schema.TypeString,
							Description: "Timezone of user",
							Optional:    true,
							Computed:    true,
						},
						"access_window_start": {
							Type:        schema.TypeString,
							Description: "Access window start time for user",
							Optional:    true,
							Computed:    true,
						},
						"access_window_end": {
							Type:        schema.TypeString,
							Description: "Access window end time for user",
							Optional:    true,
							Computed:    true,
						},
						"disabled": {
							Type:        schema.TypeBool,
							Description: "Whether account is disabled",
							Optional:    true,
							Computed:    true,
						},
						"valid_from": {
							Type:        schema.TypeString,
							Description: "Start date for when user is valid",
							Optional:    true,
							Computed:    true,
						},
						"valid_until": {
							Type:        schema.TypeString,
							Description: "End date for when user is valid",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"group_membership": {
				Type:        schema.TypeSet,
				Description: "Groups this user is a member of",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"system_permissions": {
				Type:        schema.TypeSet,
				Description: "System permissions assigned to user",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"connections": {
				Type:        schema.TypeSet,
				Description: "Connection identifiers a user has permission to read",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"connection_groups": {
				Type:        schema.TypeSet,
				Description: "Connection group identifiers a user has permission to read",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	check := validateUser(d)
	if check.HasError() {
		return check
	}

	user := convertResourceDataToGuacUser(d)

	created, err := client.CreateUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.Username)

	// Group membership
	groupMembership := setToStringSlice(d.Get("group_membership").(*schema.Set))
	if len(groupMembership) > 0 {
		var ops []guacamole.PatchOperation
		for _, g := range groupMembership {
			ops = append(ops, guacamole.AddGroupMembership(g))
		}
		if err := client.UpdateUserGroups(ctx, created.Username, ops); err != nil {
			_ = client.DeleteUser(ctx, created.Username)
			return diag.FromErr(fmt.Errorf("set group membership: %w", err))
		}
	}

	// System permissions + connection/connection-group permissions
	var permOps []guacamole.PatchOperation

	sysPerms := setToStringSlice(d.Get("system_permissions").(*schema.Set))
	if len(sysPerms) > 0 {
		if check := stringInSlice(validSystemPermissions(), sysPerms); check.HasError() {
			_ = client.DeleteUser(ctx, created.Username)
			return check
		}
		for _, p := range sysPerms {
			permOps = append(permOps, guacamole.AddSystemPermission(p))
		}
	}

	for _, conn := range setToStringSlice(d.Get("connections").(*schema.Set)) {
		permOps = append(permOps, guacamole.AddConnectionPermission(conn, guacamole.PermissionRead))
	}
	for _, cg := range setToStringSlice(d.Get("connection_groups").(*schema.Set)) {
		permOps = append(permOps, guacamole.AddConnectionGroupPermission(cg, guacamole.PermissionRead))
	}

	if len(permOps) > 0 {
		if err := client.UpdateUserPermissions(ctx, created.Username, permOps); err != nil {
			_ = client.DeleteUser(ctx, created.Username)
			return diag.FromErr(fmt.Errorf("set permissions: %w", err))
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	userID := d.Id()
	user, err := client.GetUser(ctx, userID)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read user %s: %w", userID, err))
	}

	d.Set("username", user.Username)
	d.Set("last_active", strconv.FormatInt(user.LastActive, 10))

	attributes := map[string]interface{}{
		"organizational_role": user.Attributes["guac-organizational-role"],
		"full_name":           user.Attributes["guac-full-name"],
		"email":               user.Attributes["guac-email-address"],
		"expired":             stringToBool(user.Attributes["expired"]),
		"timezone":            user.Attributes["timezone"],
		"access_window_start": user.Attributes["access-window-start"],
		"access_window_end":   user.Attributes["access-window-end"],
		"disabled":            stringToBool(user.Attributes["disabled"]),
		"valid_from":          user.Attributes["valid-from"],
		"valid_until":         user.Attributes["valid-until"],
	}
	d.Set("attributes", []interface{}{attributes})

	groups, err := client.GetUserGroups(ctx, userID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get user groups: %w", err))
	}
	d.Set("group_membership", groups)

	permissions, err := client.GetUserPermissions(ctx, userID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get user permissions: %w", err))
	}

	d.Set("system_permissions", permissions.SystemPermissions)

	var connections []string
	for id := range permissions.ConnectionPermissions {
		connections = append(connections, id)
	}
	d.Set("connections", connections)

	var connectionGroups []string
	for id := range permissions.ConnectionGroupPermissions {
		connectionGroups = append(connectionGroups, id)
	}
	d.Set("connection_groups", connectionGroups)

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("username", "password", "attributes") {
		if check := validateUser(d); check.HasError() {
			return check
		}
		user := convertResourceDataToGuacUser(d)
		if err := client.UpdateUser(ctx, d.Id(), user); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("group_membership") {
		old, new := d.GetChange("group_membership")
		oldGroups := setToStringSlice(old.(*schema.Set))
		newGroups := setToStringSlice(new.(*schema.Set))

		var ops []guacamole.PatchOperation
		for _, g := range sliceDiff(oldGroups, newGroups, false) {
			ops = append(ops, guacamole.RemoveGroupMembership(g))
		}
		addGroups := sliceDiff(newGroups, oldGroups, false)
		if len(addGroups) > 0 {
			if check := checkForDuplicates(addGroups); check.HasError() {
				return check
			}
			for _, g := range addGroups {
				ops = append(ops, guacamole.AddGroupMembership(g))
			}
		}
		if len(ops) > 0 {
			if err := client.UpdateUserGroups(ctx, d.Id(), ops); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("system_permissions") {
		old, new := d.GetChange("system_permissions")
		oldPerms := setToStringSlice(old.(*schema.Set))
		newPerms := setToStringSlice(new.(*schema.Set))

		var ops []guacamole.PatchOperation
		for _, p := range sliceDiff(oldPerms, newPerms, false) {
			ops = append(ops, guacamole.RemoveSystemPermission(p))
		}
		add := sliceDiff(newPerms, oldPerms, false)
		if len(add) > 0 {
			if check := stringInSlice(validSystemPermissions(), add); check.HasError() {
				return check
			}
			for _, p := range add {
				ops = append(ops, guacamole.AddSystemPermission(p))
			}
		}
		if len(ops) > 0 {
			if err := client.UpdateUserPermissions(ctx, d.Id(), ops); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("connections") {
		old, new := d.GetChange("connections")
		oldConns := setToStringSlice(old.(*schema.Set))
		newConns := setToStringSlice(new.(*schema.Set))

		var ops []guacamole.PatchOperation
		for _, c := range sliceDiff(oldConns, newConns, false) {
			ops = append(ops, guacamole.RemoveConnectionPermission(c, guacamole.PermissionRead))
		}
		for _, c := range sliceDiff(newConns, oldConns, false) {
			ops = append(ops, guacamole.AddConnectionPermission(c, guacamole.PermissionRead))
		}
		if len(ops) > 0 {
			if err := client.UpdateUserPermissions(ctx, d.Id(), ops); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("connection_groups") {
		old, new := d.GetChange("connection_groups")
		oldCGs := setToStringSlice(old.(*schema.Set))
		newCGs := setToStringSlice(new.(*schema.Set))

		var ops []guacamole.PatchOperation
		for _, cg := range sliceDiff(oldCGs, newCGs, false) {
			ops = append(ops, guacamole.RemoveConnectionGroupPermission(cg, guacamole.PermissionRead))
		}
		for _, cg := range sliceDiff(newCGs, oldCGs, false) {
			ops = append(ops, guacamole.AddConnectionGroupPermission(cg, guacamole.PermissionRead))
		}
		if len(ops) > 0 {
			if err := client.UpdateUserPermissions(ctx, d.Id(), ops); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteUser(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// convertResourceDataToGuacUser builds a guacamole.User from Terraform state.
func convertResourceDataToGuacUser(d *schema.ResourceData) guacamole.User {
	user := guacamole.User{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	attrList := d.Get("attributes").([]interface{})
	if len(attrList) > 0 {
		attrs := attrList[0].(map[string]interface{})
		user.Attributes = guacamole.NullableStringMap{
			"guac-organizational-role": attrs["organizational_role"].(string),
			"guac-full-name":           attrs["full_name"].(string),
			"guac-email-address":       attrs["email"].(string),
			"expired":                  boolToString(attrs["expired"].(bool)),
			"timezone":                 attrs["timezone"].(string),
			"access-window-start":      attrs["access_window_start"].(string),
			"access-window-end":        attrs["access_window_end"].(string),
			"disabled":                 boolToString(attrs["disabled"].(bool)),
			"valid-from":               attrs["valid_from"].(string),
			"valid-until":              attrs["valid_until"].(string),
		}
	}

	return user
}

// validateUser validates user-specific attributes before creating or updating.
func validateUser(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	attrList := d.Get("attributes").([]interface{})
	if len(attrList) == 0 {
		return diags
	}
	attrs := attrList[0].(map[string]interface{})

	tz := attrs["timezone"].(string)
	if tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid timezone",
				Detail:   fmt.Sprintf("Unable to process timezone string: %s", tz),
			})
		}
	}

	if vf, ok := d.GetOk("attributes.0.valid_from"); ok {
		if check := validateTimestring(vf.(string), "valid_from"); check.HasError() {
			diags = append(diags, check...)
		}
	}
	if vu, ok := d.GetOk("attributes.0.valid_until"); ok {
		if check := validateTimestring(vu.(string), "valid_until"); check.HasError() {
			diags = append(diags, check...)
		}
	}

	return diags
}

// setToStringSlice converts a *schema.Set of strings to []string.
func setToStringSlice(s *schema.Set) []string {
	items := s.List()
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.(string))
	}
	return out
}

// validateGroups checks that all provided group identifiers exist in Guacamole.
func validateGroups(ctx context.Context, client *guacamole.Client, groups []string) diag.Diagnostics {
	userGroups, err := client.ListUserGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	var invalid []string
	for _, g := range groups {
		if _, ok := userGroups[g]; !ok {
			invalid = append(invalid, g)
		}
	}
	if len(invalid) > 0 {
		return diag.Errorf("invalid user group(s) for group_membership: %s", strings.Join(invalid, ", "))
	}
	return nil
}
