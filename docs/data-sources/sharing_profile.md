---
page_title: "Sharing Profile Data Source - terraform-provider-guacamole"
subcategory: ""
description: |-
  The sharing_profile data source allows you to retrieve a guacamole sharing profile by identifier
---

# Data Source `guacamole_sharing_profile`

The sharing_profile data source allows you to retrieve a guacamole sharing profile by identifier

## Example Usage

```terraform
data "guacamole_sharing_profile" "share" {
  identifier = "42"
}
```

## Attributes Reference

The following attributes are exported.

### Base

- `identifier` - (string) Numeric identifier of the sharing profile
- `name` - (string) Name of the sharing profile
- `primary_connection_identifier` - (string) Identifier of the connection being shared

### Parameters

- `read_only` - (bool) Whether the shared session is read-only
