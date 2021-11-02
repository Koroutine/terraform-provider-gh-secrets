---
page_title: "Provider: gh-secrets"
subcategory: ""
description: |-
  Super simple provider for deploying secrets to GitHub Actions
---

# GH-Secrets Provider

Super simple provider for deploying secrets to GitHub actions, using latest Terraform sdk

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "gh-secrets" {
  token = "" # or GITHUB_TOKEN
}
```

## Schema

### Required

- **token** (String) Personal Access Token from GitHub https://github.com/settings/tokens