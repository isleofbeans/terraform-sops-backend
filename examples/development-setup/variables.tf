variable "output_dir" {
  description = "Directory to create configuration files in."
  type        = string
}

variable "age_public_key" {
  description = <<EOD
(optional) public AGE key to use.  
This is used by go test only and can stay as it is.

Make sure to create a proper key pair for production use.
EOD
  type        = string
  default     = "age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08"
}

variable "age_private_key" {
  description = <<EOD
(optional) private AGE key to use.  
This is used by go test only and can stay as it is.

Make sure to create a proper key pair for production use.
EOD
  type        = string
  default     = "AGE-SECRET-KEY-1Z22A6EL3ECQC96ZDMPD5KRPUX32SCAMU2DJGV3Q48PXN3ZW535VQFQDEF9"
  sensitive   = true
}

variable "vault_addr" {
  description = "(optional) Address the Vault service is listening on."
  type        = string
  default     = "http://127.0.0.1:8200"
}

variable "vault_approle_id" {
  description = "AppRole ID to use when using Vault to en- or decrypt."
  type        = string
  sensitive   = true
}

variable "vault_approle_secret_id" {
  description = "AppRole secret ID to use when using Vault to en- or decrypt."
  type        = string
  sensitive   = true
}

variable "vault_transit_mount_path" {
  description = "Path the transit engine is mounted to."
  type        = string
}

variable "vault_transit_backend_key" {
  description = "Name of the secret inside the transit engine."
  type        = string
}
