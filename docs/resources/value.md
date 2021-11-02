---
page_title: "value Resource - terraform-provider-gh-secrets"
subcategory: ""
description: |-
  GitHub Actions secret
---

# Resource `gh-secrets_value`

This allows you to deploy secrets to GitHub Actions

## Example Usage

```terraform
resource "gh-secrets_value" "secret" {
    repo = "me/example"
    name = "password"
    value = "myAwesomePassword"
}
```

## Argument Reference

- `repo` - (Required) Target GitHub repo (with owner).
- `name` - (Required) Name of secret.
- `value` - (Required) Value of secret.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

### Coffee

- `image` - The coffee's image URL path.
- `name` - The coffee name.
- `price` - The coffee price.
- `teaser` - The coffee teaser.