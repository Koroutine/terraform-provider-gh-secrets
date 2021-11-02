terraform {
  required_providers {
    gh-secrets = {
      source  = "koroutine/tf/gh-secrets"
    }
  }
}

provider "gh-secrets" {
    token = "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"
}

resource "gh-secrets_value" "test" {
    repo = "koroutine/esimgo-cms"
    name = "test"
    value = "test"
}
