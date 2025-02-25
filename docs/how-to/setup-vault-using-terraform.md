# HOW To: Setup Vault using terraform

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![how-to-guides](../assets/breadcrum-how-to-guides.drawio.svg)](./index.md)

## About

This HOW TO Guide will give you the required steps to setup Vault as en-/decryption service for your `terraform-sops-backend` using terraform.

## Example terraform implementation

See the [vault-setup example](../../examples/vault-setup/).  
The example is creating:

- [vault_approle_auth_backend_role.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/approle_auth_backend_role) (resource)
- [vault_approle_auth_backend_role_secret_id.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/approle_auth_backend_role_secret_id) (resource)
- [vault_auth_backend.approle](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/auth_backend) (resource)
- [vault_mount.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/mount) (resource)
- [vault_policy.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/policy) (resource)
- [vault_transit_secret_backend_key.sops](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/transit_secret_backend_key) (resource)
