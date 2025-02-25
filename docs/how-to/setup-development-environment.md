# HOW To: Setup setup development environment

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![how-to-guides](../assets/breadcrum-how-to-guides.drawio.svg)](./index.md)

## About

This HOW TO Guide will give you the required steps to setup the development environment expected by `terraform-sops-backend`.

`terraform-sops-backend` expects to have golang set up and a local running Vault to use for automated tests.  
We will guide you through the steps to achieve this.

## Procedure

### Setup required binaries

The `terraform-sops-backend` comes with a [ASDF configuration](../../.tool-versions) that can be used to install all required binaries.  
you are free to use whatever way to install this binaries but this HOW TO will show the procedure using [ASDF](https://asdf-vm.com/guide/getting-started.html).

1. Add all required ASDF plugins
     ```sh
     asdf plugin add golang
     asdf plugin add vault
     asdf plugin add pre-commit
     asdf plugin add terraform-docs
     asdf plugin add terraform
     asdf plugin add terragrunt
     asdf plugin add tflint
     asdf plugin add age
     asdf plugin add sops
     ```
2. Install all required binaries
     ```sh
     asdf install
     ```

Additionaly there are two go binaries required by `pre-commit` to lint the go code of this repository

```sh
go install golang.org/x/lint/golint
go install golang.org/x/tools/cmd/goimports
```

### Setup the required services to perform testing

1. activate the environment tooling
    ```sh
    source examples/bin/activate
    ```
2. setup the development environment
    ```sh
    setup_development
    ```

### Verify the development environment

```sh
go test ./...
```

### Tear down the development environment

```sh
deactivate
```
