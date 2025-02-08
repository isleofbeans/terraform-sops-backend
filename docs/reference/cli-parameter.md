# CLI parameter

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![reference](../assets/breadcrum-reference.drawio.svg)](./index.md)

## `terraform-sops-backend`

terraform-sops-backend serves as an intermediate terraform HTTP backend.
It is encrypting and decrypting the state file before it is passing it on
to the backend terraform HTTP backend

```
Usage:
  terraform-sops-backend [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Starting the service

Flags:
      --config string   config file (default "/etc/terraform-sops-backend/conf.yaml")
  -h, --help            help for terraform-sops-backend

Use "terraform-sops-backend [command] --help" for more information about a command.
```

## `terraform-sops-backend start`

Starts the web service for the terraform SOPS backend.

This uses the default interface to accept incoming traffic

```
Usage:
  terraform-sops-backend start [flags]

Flags:
      --age-private-key string            TRANSFORM_AGE_PRIVATE_KEY (optional) private AGE key to decrypt terraform state
      --age-public-key string             TRANSFORM_AGE_PUBLIC_KEY (required) public AGE key to encrypt terraform state
      --backend-lock-method string        BACKEND_LOCK_METHOD (optional) lock method to use with the backend terraform state server (default "LOCK")
      --backend-unlock-method string      BACKEND_UNLOCK_METHOD (optional) unlock method to use with the backend terraform state server (default "UNLOCK")
      --backend-url string                BACKEND_URL (required) base url to connect with the backend terraform state server
  -h, --help                              help for start
      --log-json                          LOG_JSON (optional) if logging has to use json format
      --log-level string                  LOG_LEVEL (optional) active log level one of [TRACE, DEBUG, INFO, WARN, ERROR, OFF] (default "INFO")
      --port string                       SERVER_PORT (optional) port the service is listening to (default "8080")
      --vault-addr string                 TRANSFORM_VAULT_ADDRESS (optional) vault address to de- and encrypt terraform state
      --vault-app-role-id string          TRANSFORM_VAULT_APP_ROLE_ID (optional) (required if --vault-addr != "") AppRole ID to authenticate with vault
      --vault-app-role-secret-id string   TRANSFORM_VAULT_APP_ROLE_SECRET_ID (optional) (required if --vault-addr != "") AppRole secret ID to authenticate with vault
      --vault-transit-mount string        TRANSFORM_VAULT_TRANSIT_MOUNT (optional) mount point of the transit engine to use (default "sops")
      --vault-transit-name string         TRANSFORM_VAULT_TRANSIT_NAME (optional) name of the transit engine secret to use (default "terraform")

Global Flags:
      --config string   config file (default "/etc/terraform-sops-backend/conf.yaml")
```
