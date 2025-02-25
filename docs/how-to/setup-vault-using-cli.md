# HOW To: Setup Vault using CLI

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![how-to-guides](../assets/breadcrum-how-to-guides.drawio.svg)](./index.md)

## About

This HOW TO Guide will give you the required steps to setup Vault as en-/decryption service for your `terraform-sops-backend` using the VAULT CLI.

## Required information

* All required information are references throughout this HOW TO in the form `%REF_NAME%`
* There are required information collected up front and other collected in the process of this HOW TO.

### Up front collected information

* We need the access URL to the Vault
* We need a mount point and secret name for the [transit secret engine](https://developer.hashicorp.com/vault/docs/secrets/transit) which will be used to en- and decrypt the terraform state.
    * If you already have an transit secret engine configured just skip the steps and use that transit secret engine by its mount point.
* We need an a [AppRole](https://developer.hashicorp.com/vault/docs/auth/approle) to let `terraform-sops-backend` authenticate against the Vault service and a corresponding [policy](https://developer.hashicorp.com/vault/docs/concepts/policies) to authorize the AppRole for en- and decryption.

| reference throughout this HOW TO | Description                                         |
| -------------------------------- | --------------------------------------------------- |
| `%VAULT_ADDR%`                   | The access URL to the Vault                         |
| `%MOUNT_POINT%`                  | The path to mount the transit secret engine to      |
| `%SECRET_NAME%`                  | The name of the secret in the transit secret engine |
| `%POLICY_NAME%`                  | The name of the authorization policy                |
| `%APP_ROLE_NAME%`                | The name of the AppRole                             |

### Information collected in the process of this HOW TO

| reference throughout this HOW TO | Description                    |
| -------------------------------- | ------------------------------ |
| `%APP_ROLE_ID%`                  | The collected AppRole ID       |
| `%APP_ROLE_SECRET_ID%`           | The collected AppRole SecretID |

## Procedure

### Setup the Vault CLI client

1. Configure your Vault CLI client to connect with your Vault service. E.g.: by setting the environment variable `VAULT_ADDR`.
2. Use your Vault CLI client to perform the required login. (This depends on your setup and is not part of this HOW TO guide)
3. Setup environment variables required for the scriptlets of this HOW TO
     ```sh
     MOUNT_POINT=%MOUNT_POINT%
     SECRET_NAME=%SECRET_NAME%
     POLICY_NAME=%POLICY_NAME%
     APP_ROLE_NAME=%APP_ROLE_NAME%
     ```

### Setup the transit secret engine

1. Create the transient mount
    ```sh
    vault secrets enable -path "${MOUNT_POINT}" transit
    ```
2. Create the transient backend key
    ```sh
    vault write -f "${MOUNT_POINT}/keys/${SECRET_NAME}"
    ```

### Setup the AppRole

1. Create a policy file for our approle `approle-policy.hcl`
    ```hcl
    path "%MOUNT_POINT%/encrypt/%SECRET_NAME%" {
      capabilities = ["update"]
    }
    path "%MOUNT_POINT%/decrypt/%SECRET_NAME%" {
      capabilities = ["update"]
    }
    ```
    * Replace the placeholders `%MOUNT_POINT%` and `%SECRET_NAME%`
2. Upload the created policy
    ```sh
    vault policy write "${POLICY_NAME}" approle-policy.hcl
    ```
3. Enable approle login
    ```sh
    vault auth enable approle
    ```
4.  Create an approle
    ```sh
    vault write "auth/approle/role/${APP_ROLE_NAME}" \
        token_type=batch \
        secret_id_ttl=0 \
        token_ttl=20m \
        token_max_ttl=30m \
        secret_id_num_uses=0 \
        policies="${POLICY_NAME}"
    ```
5. Read the approle ID
    ```sh
    vault read "auth/approle/role/${APP_ROLE_NAME}/role-id"
    ```
    * The approle ID is from now on references as `%APP_ROLE_ID%`
6. Create an approle secret
    ```sh
    vault write -f "auth/approle/role/${APP_ROLE_NAME}/secret-id"
    ```
    * The approle secret ID is from now on references as `%APP_ROLE_SECRET_ID%`

### Setup `terraform-sops-backend`

This is just the vault part of the start configuration.  
See [terraform-sops-backend cli reference](../reference/cli-parameter.md) for the full list of parameters

```sh
terraform-sops-backend start \
--vault-addr %VAULT_ADDR% \
--vault-app-role-id %APP_ROLE_ID% \
--vault-app-role-secret-id %APP_ROLE_SECRET_ID% \
--vault-transit-mount %MOUNT_POINT% \
--vault-transit-name %SECRET_NAME%
```
