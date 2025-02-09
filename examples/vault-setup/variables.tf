variable "transit_mount_path" {
  description = "mount path for the example transit engine"
  type        = string
  default     = "sops"
}
variable "transit_mount_description" {
  description = "description"
  type        = string
  default     = "Example transit engine"
}

variable "transit_backend_name" {
  description = "TODO"
  type        = string
  default     = "terraform"
}

variable "transit_role_name" {
  description = "TODO"
  type        = string
  default     = "sops"
}

variable "vault_approle_backend_path" {
  type        = string
  description = "(optional) the mount path of the approle auth method"
  default     = "approle"
}

variable "role_token_ttl" {
  description = "The TTL for role tokens - in seconds"
  type        = number
  default     = 3600
}

variable "role_max_token_ttl" {
  description = "The max TTL for role tokens - in seconds"
  type        = number
  default     = 68400
}
