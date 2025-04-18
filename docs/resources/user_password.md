---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kestra_user_password Resource - terraform-provider-kestra"
subcategory: ""
description: |-
  Manages a Kestra User Basic Auth Password.
  -> This resource is only available on the Enterprise Edition https://kestra.io/enterprise
---

# kestra_user_password (Resource)

Manages a Kestra User Basic Auth Password.

-> This resource is only available on the [Enterprise Edition](https://kestra.io/enterprise)

## Example Usage

```terraform
resource "kestra_user_password" "example" {
  user_id  = "4by6NvSLcPXFhCj8nwbZOM"
  password = "my-random-password"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `password` (String, Sensitive) The user password.
- `user_id` (String) The user id.

### Read-Only

- `id` (String) The ID of this resource.
