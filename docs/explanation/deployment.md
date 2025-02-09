# How dose it look from the deployment perspective

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![explanation](../assets/breadcrum-explanation.drawio.svg)](./index.md)

The deployment is separated into two environments.  
One is controlled by the user of terraform-sops-backend and the other is controlled by the provider of the used terraform HTTP backend.

## System deployment

![architecture](../assets/architecture.drawio.svg)

## Participants

* `Platform Engineer`
    * Is performing terraform / terragrunt actions on a `Terraform / Terragrunt module`
* `Terraform / Terragrunt module`
    * Is performing read, write, lock and unlock actions against a terraform backend
    * In a terraform-sops-backend scenario this is a terraform HTTP backend represented by a `terraform-sops-backend` service
* `terraform-sops-backend`
    * Is encrypting terraform states send by the client
    * Is decrypting terraform states received from the backend terraform HTTP backend
        * Uses Vault and / or AGE depending on configuration
    * Is forwarding incoming requests to the `terraform HTTP backend` with encrypted terraform states
