provider "singularity" {
  host = "localhost"
  port = 443
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
  instances              = 1
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
  instances    = 2
}

resource "singularity_docker_deploy" "test-deploy" {
  deploy_id  = "mydeploy4ab"
  command    = "bash"
  args       = ["-xc", "sleep 10000"]
  request_id = "${singularity_request.lenfree-service.id}"

  resources {
    cpus      = 2
    memory_mb = 128
  }

  envs {
    MYENV = "Test"
    OWNER = "lenfree"
  }

  docker_info {
    force_pull_image = false
    network          = "BRIDGE"
    image            = "golang:latest"

    port_mapping {
      host_port           = 0
      container_port      = 10001
      container_port_type = "LITERAL"
      host_port_type      = "FROM_OFFER"
      protocol            = "tcp"
    }
  }
}

resource "singularity_docker_deploy" "test-deploy-2" {
  deploy_id  = "mydeploy2"
  command    = "bash"
  args       = ["-xc", "env"]
  request_id = "${singularity_request.lenfree-demand.id}"

  resources {
    cpus      = 2
    memory_mb = 128
  }

  envs {
    MYENV = "Test"
    OWNER = "lenfree"
  }

  docker_info {
    force_pull_image = false
    network          = "BRIDGE"
    image            = "golang:latest"

    port_mapping {
      host_port           = 0
      container_port      = 10001
      container_port_type = "LITERAL"
      host_port_type      = "FROM_OFFER"
      protocol            = "tcp"
    }
  }

  volume {
    host_path      = "/outside/path"
    container_path = "/inside/path"
    mode           = "RO"
  }

  uri {
    path       = "file:///app/config.json"
    cache      = false
    executable = false
    extract    = true
  }
}

resource "singularity_request" "imoussa-demand" {
  request_id      = "imoussa-ondemand-2"
  request_type    = "ON_DEMAND"
  instances       = 2
  slave_placement = "SEPARATE_BY_DEPLOY"
}
