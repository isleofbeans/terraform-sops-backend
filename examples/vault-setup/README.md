<!-- BEGIN_TF_DOCS -->
## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.10)

- <a name="requirement_vault"></a> [vault](#requirement\_vault) (~> 4.6)

## Providers

The following providers are used by this module:

- <a name="provider_vault"></a> [vault](#provider\_vault) (~> 4.6)

## Resources

The following resources are used by this module:

- [vault_approle_auth_backend_role.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/approle_auth_backend_role) (resource)
- [vault_approle_auth_backend_role_secret_id.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/approle_auth_backend_role_secret_id) (resource)
- [vault_auth_backend.approle](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/auth_backend) (resource)
- [vault_mount.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/mount) (resource)
- [vault_policy.transit](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/policy) (resource)
- [vault_transit_secret_backend_key.sops](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/transit_secret_backend_key) (resource)

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_role_max_token_ttl"></a> [role\_max\_token\_ttl](#input\_role\_max\_token\_ttl)

Description: (optional) The max TTL for role tokens - in seconds

Type: `number`

Default: `68400`

### <a name="input_role_token_ttl"></a> [role\_token\_ttl](#input\_role\_token\_ttl)

Description: (optional) The TTL for role tokens - in seconds

Type: `number`

Default: `3600`

### <a name="input_transit_backend_name"></a> [transit\_backend\_name](#input\_transit\_backend\_name)

Description: (optional) Name of the secret backend key in the transit engine mount.

Type: `string`

Default: `"terraform"`

### <a name="input_transit_mount_description"></a> [transit\_mount\_description](#input\_transit\_mount\_description)

Description: (optional) Description to the transit engine mount

Type: `string`

Default: `"Example transit engine"`

### <a name="input_transit_mount_path"></a> [transit\_mount\_path](#input\_transit\_mount\_path)

Description: (optional) Path the transit engine is mounted to.

Type: `string`

Default: `"sops"`

### <a name="input_transit_role_name"></a> [transit\_role\_name](#input\_transit\_role\_name)

Description: (optional) Name of the AppRole used to en- and decrypt terraform state

Type: `string`

Default: `"sops"`

### <a name="input_vault_approle_backend_path"></a> [vault\_approle\_backend\_path](#input\_vault\_approle\_backend\_path)

Description: (optional) Path there the AppRole auth method is mounted

Type: `string`

Default: `"approle"`

## Outputs

The following outputs are exported:

### <a name="output_approle_id"></a> [approle\_id](#output\_approle\_id)

Description: The created AppRole ID

### <a name="output_approle_secret_id"></a> [approle\_secret\_id](#output\_approle\_secret\_id)

Description: The created AppRole secret ID

### <a name="output_transit_backend_key"></a> [transit\_backend\_key](#output\_transit\_backend\_key)

Description: Name of the secret backend key in the transit engine mount.

### <a name="output_transit_mount_path"></a> [transit\_mount\_path](#output\_transit\_mount\_path)

Description: Path the transit engine is mounted to.
<!-- END_TF_DOCS -->
