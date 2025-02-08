locals {
  basedir         = get_env("BASEDIR", abspath("${get_terragrunt_dir()}/../.."))
  terragrunt_root = get_env("EXAMPLES_DIR", abspath("${get_terragrunt_dir()}/.."))
}

terraform {
  source = ".//."
}

dependencies {
  paths = [
    "${local.terragrunt_root}/vault-setup"
  ]
}

dependency "vault_setup" {
  config_path = "${local.terragrunt_root}/vault-setup"
}

generate "configuration_input" {
  path              = "terragrunt-input.auto.tfvars.json"
  if_exists         = "overwrite"
  disable_signature = true
  contents = jsonencode(
    {
      output_dir                = local.basedir
      vault_approle_id          = dependency.vault_setup.outputs.approle_id
      vault_approle_secret_id   = dependency.vault_setup.outputs.approle_secret_id
      vault_transit_mount_path  = dependency.vault_setup.outputs.transit_mount_path
      vault_transit_backend_key = dependency.vault_setup.outputs.transit_backend_key
    }
  )
}
