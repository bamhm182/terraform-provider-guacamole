package guacamole

import (
	"context"
	"fmt"

	"github.com/bamhm182/go-guacamole/guacamole"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guacamoleUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "Identifier of guacamole user group",
				Required:    true,
				ForceNew:    true,
			},
			"attributes": {
				Type:        schema.TypeList,
				Description: "Attributes of guacamole user group",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled": {
							Type:        schema.TypeBool,
							Description: "Whether group is disabled",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"group_membership": {
				Type:        schema.TypeSet,
				Description: "Groups this user group is a member of (parent groups)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"system_permissions": {
				Type:        schema.TypeSet,
				Description: "System permissions assigned to user group",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"connections": {
				Type:        schema.TypeSet,
				Description: "Connection identifiers this group has permission to read",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"connection_groups": {
				Type:        schema.TypeSet,
				Description: "Connection group identifiers this group has permission to read",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	group := convertResourceDataToGuacUserGroup(d)

	created, err := client.CreateUserGroup(ctx, group)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.Identifier)

	// Parent group membership
	groupMembership := setToStringSlice(d.Get("group_membership").(*schema.Set))
	if len(groupMembership) > 0 {
		var ops []guacamole.PatchOperation
		for _, g := range groupMembership {
			ops = append(ops, guacamole.AddGroupMembership(g))
		}
		if err := client.UpdateUserGroupParentGroups(ctx, created.Identifier, ops); err != nil {
			_ = client.DeleteUserGroup(ctx, created.Identifier)
			return diag.FromErr(fmt.Errorf("set group membership: %w", err))
		}
	}

	// Permissions
	var permOps []guacamole.PatchOperation

	sysPerms := setToStringSlice(d.Get("system_permissions").(*schema.Set))
	if len(sysPerms) > 0 {
		if check := stringInSlice(validSystemPermissions(), sysPerms); check.HasError() {
			_ = client.DeleteUserGroup(ctx, created.Identifier)
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
		if err := client.UpdateUserGroupPermissions(ctx, created.Identifier, permOps); err != nil {
			_ = client.DeleteUserGroup(ctx, created.Identifier)
			return diag.FromErr(fmt.Errorf("set permissions: %w", err))
		}
	}

	return resourceUserGroupRead(ctx, d, m)
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	id := d.Id()
	group, err := client.GetUserGroup(ctx, id)
	if err != nil {
		if guacamole.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("read user group %s: %w", id, err))
	}

	d.Set("identifier", group.Identifier)

	attributes := map[string]interface{}{
		"disabled": stringToBool(group.Attributes["disabled"]),
	}
	d.Set("attributes", []interface{}{attributes})

	parentGroups, err := client.GetUserGroupParentGroups(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get parent groups: %w", err))
	}
	d.Set("group_membership", parentGroups)

	permissions, err := client.GetUserGroupPermissions(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get user group permissions: %w", err))
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

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)

	if d.HasChanges("attributes") {
		group := convertResourceDataToGuacUserGroup(d)
		if err := client.UpdateUserGroup(ctx, d.Id(), group); err != nil {
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
			if err := client.UpdateUserGroupParentGroups(ctx, d.Id(), ops); err != nil {
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
			if err := client.UpdateUserGroupPermissions(ctx, d.Id(), ops); err != nil {
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
			if err := client.UpdateUserGroupPermissions(ctx, d.Id(), ops); err != nil {
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
			if err := client.UpdateUserGroupPermissions(ctx, d.Id(), ops); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceUserGroupRead(ctx, d, m)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*guacamole.Client)
	if err := client.DeleteUserGroup(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func convertResourceDataToGuacUserGroup(d *schema.ResourceData) guacamole.UserGroup {
	group := guacamole.UserGroup{
		Identifier: d.Get("identifier").(string),
	}

	attrList := d.Get("attributes").([]interface{})
	if len(attrList) > 0 {
		attrs := attrList[0].(map[string]interface{})
		group.Attributes = guacamole.NullableStringMap{
			"disabled": boolToString(attrs["disabled"].(bool)),
		}
	}

	return group
}
