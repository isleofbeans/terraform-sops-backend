variable "output_dir" {
  description = "TODO"
  type        = string
}

variable "age_public_key" {
  description = "TODO"
  type        = string
  default     = "age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08"
}

variable "age_private_key" {
  description = "TODO"
  type        = string
  default     = "AGE-SECRET-KEY-1Z22A6EL3ECQC96ZDMPD5KRPUX32SCAMU2DJGV3Q48PXN3ZW535VQFQDEF9"
  sensitive   = true
}

variable "vault_addr" {
  description = "TODO"
  type        = string
  default     = "http://127.0.0.1:8200"
}

variable "vault_approle_id" {
  description = "TODO"
  type        = string
  sensitive   = true
}

variable "vault_approle_secret_id" {
  description = "TODO"
  type        = string
  sensitive   = true
}

variable "vault_transit_mount_path" {
  description = "TODO"
  type        = string
}

variable "vault_transit_backend_key" {
  description = "TODO"
  type        = string
}
