variable "transit_mount_path" {
  description = "(optional) Path the transit engine is mounted to."
  type        = string
  default     = "sops"
}
variable "transit_mount_description" {
  description = "(optional) Description to the transit engine mount"
  type        = string
  default     = "Example transit engine"
}

variable "transit_backend_name" {
  description = "(optional) Name of the secret backend key in the transit engine mount."
  type        = string
  default     = "terraform"
}

variable "transit_role_name" {
  description = "(optional) Name of the AppRole used to en- and decrypt terraform state"
  type        = string
  default     = "sops"
}

variable "vault_approle_backend_path" {
  description = "(optional) Path there the AppRole auth method is mounted"
  type        = string
  default     = "approle"
}

variable "role_token_ttl" {
  description = "(optional) The TTL for role tokens - in seconds"
  type        = number
  default     = 3600
}

variable "role_max_token_ttl" {
  description = "(optional) The max TTL for role tokens - in seconds"
  type        = number
  default     = 68400
}
