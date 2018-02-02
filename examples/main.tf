provider "singularity" {
    host = "localhost"
}

resource "singularity_request" "my-server" {
  request_id             = "lenfree-test-tf"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
}

resource "singularity_request" "lenfree-run" {
  request_id   = "lenfree-test-runonce"
  request_type = "RUN_ONCE"
  instances    = 5
}

resource "singularity_request" "lenfree-demand" {
  request_id   = "lenfree-ondemand"
  request_type = "ON_DEMAND"
}

resource "singularity_request" "lenfree-scheduled" {
  request_id             = "lenfree-test-scheduled"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
}

resource "singularity_request" "lenfree-service" {
  request_id   = "lenfree-test-service"
  request_type = "SERVICE"
  instances    = 2
}

resource "singularity_request" "lenfree-worker" {
  request_id   = "lenfree-test-worker"
  request_type = "WORKER"
  instances    = 2
}
