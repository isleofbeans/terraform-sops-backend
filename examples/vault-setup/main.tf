
resource "vault_mount" "transit" {
  path        = var.transit_mount_path
  type        = "transit"
  description = var.transit_mount_description
  options = {
    convergent_encryption = false
  }
}

resource "vault_transit_secret_backend_key" "sops" {
  backend = vault_mount.transit.path
  name    = var.transit_backend_name
  type    = "aes256-gcm96"
}

resource "vault_policy" "transit" {
  name = "transit-policy"

  policy = <<EOT
path "${var.transit_mount_path}/encrypt/${var.transit_backend_name}" {
  capabilities = ["update"]
}
path "${var.transit_mount_path}/decrypt/${var.transit_backend_name}" {
  capabilities = ["update"]
}
EOT
}

resource "vault_auth_backend" "approle" {
  type = "approle"
}

resource "vault_approle_auth_backend_role" "transit" {
  backend        = vault_auth_backend.approle.path
  role_name      = var.transit_role_name
  token_policies = [vault_policy.transit.name]
  token_ttl      = var.role_token_ttl
  token_max_ttl  = var.role_max_token_ttl
}

resource "vault_approle_auth_backend_role_secret_id" "transit" {
  backend   = var.vault_approle_backend_path
  role_name = vault_approle_auth_backend_role.transit.role_name
}
