---
page_title: "Sharing Profile Resource - terraform-provider-guacamole"
subcategory: ""
description: |-
  The sharing_profile resource allows you to configure a guacamole sharing profile
---

# Resource `guacamole_sharing_profile`

The sharing_profile resource allows you to configure a guacamole sharing profile. A sharing profile defines secondary connection parameters used when sharing an active session, typically to allow read-only access.

## Example Usage

```terraform
resource "guacamole_connection_rdp" "example" {
  name              = "Example RDP"
  parent_identifier = "ROOT"
  parameters {
    hostname = "example.com"
    port     = 3389
  }
}

resource "guacamole_sharing_profile" "readonly" {
  name                          = "Read-Only Share"
  primary_connection_identifier = guacamole_connection_rdp.example.identifier
  parameters {
    read_only = true
  }
}
```

## Argument Reference

### Base

- `name` - (string, Required) Name of the sharing profile
- `primary_connection_identifier` - (string, Required) Identifier of the connection being shared

### Parameters

- `read_only` - (bool) Whether the shared session is read-only

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

#### Base
- `identifier` - (string) Numeric identifier of the sharing profile

## Import

Sharing profile can be imported using the `resource id`, e.g.

```shell
terraform import guacamole_sharing_profile.readonly 42
```
