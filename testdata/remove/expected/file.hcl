resource "example_resource" "hello" {
  name    = "hello-world"
  message = "Hello, World!"

  metadata {
    created_by = "example"
  }
}