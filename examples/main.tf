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

resource "singularity_request" "my-server-run" {
  request_id   = "lenfree-test-runonce"
  request_type = "RUN_ONCE"
}

resource "singularity_request" "my-server-demand" {
  request_id   = "lenfree-test-ondemand"
  request_type = "ON_DEMAND"
}

resource "singularity_request" "my-server-scheduled" {
  request_id   = "lenfree-test-scheduled"
  request_type = "SCHEDULED"
}

resource "singularity_request" "my-server-service" {
  request_id   = "lenfree-test-service"
  request_type = "service"
}

resource "singularity_request" "my-server-worker" {
  request_id   = "lenfree-test-worker"
  request_type = "worker"
}
