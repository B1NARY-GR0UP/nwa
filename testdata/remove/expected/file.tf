output "hello_world" {
  value = "Hello, World!"
}

# Define a local variable
locals {
  greeting = "Hello"
  subject  = "Terraform"
}

# Output using variables
output "greeting" {
  value = "${local.greeting}, ${local.subject}!"
}