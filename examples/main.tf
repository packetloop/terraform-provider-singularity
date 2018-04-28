provider "singularity" {
  host = "localhost/singularity"
}

resource "singularity_request" "my-server" {
  request_id             = "lenfree-test-tf"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
  instances              = 1
}

resource "singularity_request" "lenfree-run" {
  request_id   = "lenfree-test-runonce"
  request_type = "RUN_ONCE"
  instances    = 5
}

resource "singularity_request" "lenfree-scheduled" {
  request_id             = "lenfree-test-scheduled"
  request_type           = "SCHEDULED"
  num_retries_on_failure = 3
  schedule               = "0 7 * * *"
  schedule_type          = "CRON"
  instances              = 2
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

resource "singularity_request" "lenfree-demand" {
  request_id   = "lenfree-ondemand-2"
  request_type = "ON_DEMAND"
}

resource "singularity_docker_deploy" "test-deploy" {
  deploy_id        = "mydeploy"
  force_pull_image = false
  network          = "bridge"
  image            = "golang:latest"
  cpu              = 2
  memory           = 128
  command          = "bash"
  args             = ["-xc", "date"]
  request_id       = "${singularity_request.lenfree-demand.id}"
}
