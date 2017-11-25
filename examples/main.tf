provider "singularity" {
  host = "localhost"
}

resource "singularity_request" "my-server" {
  request_id             = "lenfree-test-tf"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 13 * * *"
  schedule_type          = "CRON"
}

resource "singularity_request" "fail" {
  request_id             = "lenfree-test-fail"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 13 * * *"
  schedule_type          = "CRON"
}
