output "approle_id" {
  description = "The created AppRole ID"
  value       = vault_approle_auth_backend_role.transit.role_id
  sensitive   = true
}

output "approle_secret_id" {
  description = "The created AppRole secret ID"
  value       = vault_approle_auth_backend_role_secret_id.transit.secret_id
  sensitive   = true
}

output "transit_mount_path" {
  description = "Path the transit engine is mounted to."
  value       = vault_mount.transit.path
}

output "transit_backend_key" {
  description = "Name of the secret backend key in the transit engine mount."
  value       = vault_transit_secret_backend_key.sops.name
}
