resource "local_file" "foo" {
  content         = <<EOC
TRANSFORM_AGE_PUBLIC_KEY=${var.age_public_key}
TRANSFORM_AGE_PRIVATE_KEY=${var.age_private_key}
TRANSFORM_VAULT_ADDRESS=${var.vault_addr}
TRANSFORM_VAULT_APP_ROLE_ID=${var.vault_approle_id}
TRANSFORM_VAULT_APP_ROLE_SECRET_ID=${var.vault_approle_secret_id}
TRANSFORM_VAULT_TRANSIT_MOUNT=${var.vault_transit_mount_path}
TRANSFORM_VAULT_TRANSIT_NAME=${var.vault_transit_backend_key}
EOC
  filename        = "${var.output_dir}/.testenv"
  file_permission = "0600"
}
