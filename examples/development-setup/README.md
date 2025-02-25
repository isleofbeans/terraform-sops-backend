<!-- BEGIN_TF_DOCS -->
## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.10)

- <a name="requirement_local"></a> [local](#requirement\_local) (~> 2.5)

## Providers

The following providers are used by this module:

- <a name="provider_local"></a> [local](#provider\_local) (~> 2.5)

## Resources

The following resources are used by this module:

- [local_file.foo](https://registry.terraform.io/providers/hashicorp/local/latest/docs/resources/file) (resource)

## Required Inputs

The following input variables are required:

### <a name="input_output_dir"></a> [output\_dir](#input\_output\_dir)

Description: Directory to create configuration files in.

Type: `string`

### <a name="input_vault_approle_id"></a> [vault\_approle\_id](#input\_vault\_approle\_id)

Description: AppRole ID to use when using Vault to en- or decrypt.

Type: `string`

### <a name="input_vault_approle_secret_id"></a> [vault\_approle\_secret\_id](#input\_vault\_approle\_secret\_id)

Description: AppRole secret ID to use when using Vault to en- or decrypt.

Type: `string`

### <a name="input_vault_transit_backend_key"></a> [vault\_transit\_backend\_key](#input\_vault\_transit\_backend\_key)

Description: Name of the secret inside the transit engine.

Type: `string`

### <a name="input_vault_transit_mount_path"></a> [vault\_transit\_mount\_path](#input\_vault\_transit\_mount\_path)

Description: Path the transit engine is mounted to.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_age_private_key"></a> [age\_private\_key](#input\_age\_private\_key)

Description: (optional) private AGE key to use.  
This is used by go test only and can stay as it is.

Make sure to create a proper key pair for production use.

Type: `string`

Default: `"AGE-SECRET-KEY-1Z22A6EL3ECQC96ZDMPD5KRPUX32SCAMU2DJGV3Q48PXN3ZW535VQFQDEF9"`

### <a name="input_age_public_key"></a> [age\_public\_key](#input\_age\_public\_key)

Description: (optional) public AGE key to use.  
This is used by go test only and can stay as it is.

Make sure to create a proper key pair for production use.

Type: `string`

Default: `"age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08"`

### <a name="input_vault_addr"></a> [vault\_addr](#input\_vault\_addr)

Description: (optional) Address the Vault service is listening on.

Type: `string`

Default: `"http://127.0.0.1:8200"`
<!-- END_TF_DOCS -->
