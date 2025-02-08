terraform {
  source = ".//."
}
generate "provider_configuration" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
    provider "vault" {
      address = "http://127.0.0.1:8200"
      token   = "${get_env("VAULT_DEV_ROOT_TOKEN_ID", "dev-token")}"
    }
    EOF
}
