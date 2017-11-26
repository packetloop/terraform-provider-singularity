provider "singularity" {
  host = "localhost.com"
}

resource "singularity_request" "my-server" {
  request_id             = "lenfree-test-tf"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
}
