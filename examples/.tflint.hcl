plugin "terraform" {
  enabled = true
  preset  = "all"
}

# # This Plugin is deactivated until we no longer need to download it from Github
# plugin "opa" {
#   enabled = true
#   version = "0.4.0"
#   source  = "github.com/terraform-linters/tflint-ruleset-opa"
# }

rule "terraform_naming_convention" {
  enabled = false
}

rule "terraform_module_pinned_source" {
  enabled = false
}
