output "approle_id" {
  description = "TODO"
  value       = vault_approle_auth_backend_role.transit.role_id
  sensitive   = true
}

output "approle_secret_id" {
  description = "TODO"
  value       = vault_approle_auth_backend_role_secret_id.transit.secret_id
  sensitive   = true
}

output "transit_mount_path" {
  description = "TODO"
  value       = vault_mount.transit.path
}

output "transit_backend_key" {
  description = "TODO"
  value       = vault_transit_secret_backend_key.sops.name
}
