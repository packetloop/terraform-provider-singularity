provider "singularity" {
  host = "localhost"
}

resource "singularity_request" "my-server" {
  address = "1.2.3.4"
}
