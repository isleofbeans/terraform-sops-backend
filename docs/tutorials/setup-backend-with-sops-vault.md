# Tutorial: Setup terraform project with Vault encrypted state on gitlab.com

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![tutorials](../assets/breadcrum-tutorials.drawio.svg)](./index.md)

## About

This tutorial will guide you through the steps to setup the Terraform SOPS backend to save your [Vault](https://developer.hashicorp.com/vault) encrypted state in the terraform state backend provided by <https://gitlab.com>.

After finishing this tutorial you'll have started Terraform SOPS backend on your machine configured to
* use <https://gitlab.com> as terraform state backend.
* use a Vault service to encrypt and decrypt all terraform states send to the backend.

An AGE public key is provided in the configuration to have a fallback if the Vault service is not available.  
Handling such scenarios is not topic of this tutorial.

## Actions

### Collect required data from gitlab.com

1. Go to <https://gitlab.com/-/user_settings/personal_access_tokens>
2. Create an API scoped access token ![create-api-scoped-token](../assets/create-api-scoped-token.png)
3. Save the value of the token for later use.
    * The token value is from now on referenced as `%gitlab-api-token%`
4. Open the gitlab.com project you want to use for your terraform state.
5. Get the project ID from `Settings` > `General`.
    * The ID is from now on referenced as `%project-id%`.

### Collect required data from your host.

1. Identify the network IP address of your host.
   * e.g.: using `ip addr`
   * The IP address is from now on references as `%host-ip-address%`

### Setup your environment

```sh
export TF_HTTP_USERNAME=oauth
export TF_HTTP_PASSWORD=%gitlab-api-token%
export VAULT_TOKEN=dev-token
export VAULT_ADDR=http://%host-ip-address%:8200
CONTAINER_COMMAND=podman # change to docker if you prefer to use docker
PROJECT_ID=%project-id%
```

### Setup a local Vault instance in develop mode (do not use for production)

1. Start vault service
    ```sh
    ${CONTAINER_COMMAND} run --rm -d --name vault -e VAULT_DEV_ROOT_TOKEN_ID=dev-token -p 8200:8200 docker.io/hashicorp/vault:latest
    ```
2. Open a shell in the vault service container
    ```sh
    ${CONTAINER_COMMAND} exec -it vault sh
    ```
3. Login to the container local service
    ```sh
    export VAULT_ADDR=http://localhost:8200
    export VAULT_TOKEN=dev-token
    ```
4. Create the transient mount
    ```sh
    vault secrets enable -path sops transit
    ```
5. Create the transient backend key
    ```sh
    vault write -f sops/keys/terraform
    ```
6. Create the policy for our approle `approle-policy.hcl`
    ```hcl
    path "sops/encrypt/terraform" {
      capabilities = ["update"]
    }
    path "sops/decrypt/terraform" {
      capabilities = ["update"]
    }
    ```
7. Upload the created policy
    ```sh
    vault policy write sops approle-policy.hcl
    ```
8. Enable approle login
    ```sh
    vault auth enable approle
    ```
9.  Create an approle
    ```sh
    vault write auth/approle/role/sops \
        token_type=batch \
        secret_id_ttl=0 \
        token_ttl=20m \
        token_max_ttl=30m \
        secret_id_num_uses=0 \
        policies=sops
    ```
10. Read the approle ID
    ```sh
    vault read auth/approle/role/sops/role-id
    ```
    * The approle ID is from now on references as `%approle-id%`
11. Create an approle secret
    ```sh
    vault write -f auth/approle/role/sops/secret-id
    ```
    * The approle secret ID is from now on references as `%approle-secret-id%`
12. Leave the container using `CTRL+d`

### Perform the tutorial

1. Prepare your terraform-sops-backend configuration as `terraform-sops-backend.env`
    ```sh
    # GitLab uses POST instead of LOCK to acquire a state lock
    BACKEND_LOCK_METHOD=POST
    # GitLab uses DELETE instead of UNLOCK to release a state lock
    BACKEND_UNLOCK_METHOD=DELETE
    BACKEND_URL=https://gitlab.com
    TRANSFORM_AGE_PUBLIC_KEY=age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08
    TRANSFORM_VAULT_ADDRESS=http://%host-ip-address%:8200
    TRANSFORM_VAULT_APP_ROLE_ID=%approle-id%
    TRANSFORM_VAULT_APP_ROLE_SECRET_ID=%approle-secret-id%
    ```
2. Start the container in background
    ```sh
    ${CONTAINER_COMMAND} run --rm -d --name terraform-sops-backend --env-file terraform-sops-backend.env  -p 8080:8080 ghcr.io/isleofbeans/terraform-sops-backend:latest
    ```
5. Create a directory for your terraform module
    ```sh
    mkdir vault-encrypted-terraform-state
    cd vault-encrypted-terraform-state
    ```
6. Create a `main.tf` with all the terraform content.
    ```terraform
    resource "random_password" "this" {
      length = 21
    }

    output "password" {
      sensitive = true
      value = random_password.this.result
    }

    terraform {
      # REPLACE %project-id% with your actual project ID from GitLab
      # in vi use :%s/.project-id./%project-id%/g
      backend "http" {
        address        = "http://localhost:8080/api/v4/projects/%project-id%/terraform/state/sops-vault"
        lock_address   = "http://localhost:8080/api/v4/projects/%project-id%/terraform/state/sops-vault/lock"
        unlock_address = "http://localhost:8080/api/v4/projects/%project-id%/terraform/state/sops-vault/lock"
        retry_wait_min = 5
      }
    }

    terraform {
      required_version = "~> 1.10"
      required_providers {
        random = {
          source = "hashicorp/random"
          version = "~> 3.6"
        }
      }
    }
    ```
7. Initialize terraform
    ```sh
    terraform init
    ```
8. Run terraform apply
    ```sh
    terraform apply -auto-approve
    ```
9.  Curl the state from gitlab.com to see it is actually encrypted
    ```sh
    curl -u "${TF_HTTP_USERNAME}:${TF_HTTP_PASSWORD}" "https://gitlab.com/api/v4/projects/${PROJECT_ID}/terraform/state/sops-vault"
    ```
10. Curl the state from localhost:8080 to see terraform-sops-backend is decrypting the state
    ```sh
    curl -u "${TF_HTTP_USERNAME}:${TF_HTTP_PASSWORD}" "http://localhost:8080/api/v4/projects/${PROJECT_ID}/terraform/state/sops-vault"
    ```
11. Show output from terraform
    ```sh
    terraform output -raw password
    ```

### Cleanup

To clean up your workplace stop the running container providing the terraform-sops-backend service.

```sh
${CONTAINER_COMMAND} stop terraform-sops-backend
${CONTAINER_COMMAND} stop vault
```
