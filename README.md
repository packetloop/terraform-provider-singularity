[![CircleCI](https://circleci.com/gh/packetloop/terraform-provider-singularity.svg?style=svg)](https://circleci.com/gh/packetloop/terraform-provider-singularity)
[![GitHub release](https://img.shields.io/github/release/packetloop/terraform-provider-singularity.svg)](https://github.com/packetloop/terraform-provider-singularity/releases/)
[![All Contributors](https://img.shields.io/github/contributors/packetloop/terraform-provider-singularity.svg?longCache=true&style=flat-square&colorB=orange&label=all%20contributors)](#contributors)
[![Github All Releases](https://img.shields.io/github/downloads/packetloop/terraform-provider-singularity/total.svg)]()


# terraform-provider-singularity

A terraform provider to manage Mesos hubspot/Singularity objects.

## Usage:

Download this provider, pick a version you'd like from releases from
[Binary Releases](https://github.com/packetloop/terraform-provider-singularity/releases)

```bash
curl -L \
  https://github.com/packetloop/terraform-provider-singularity/releases/download/v0.1.0/terraform-provider-singularity_v0.1.0_Darwin_x86_64 \
  -o ~/.terraform.d/plugins/terraform-provider-singularity_v0.1.0 && \
  chmod +x ~/.terraform.d/plugins/terraform-provider-singularity_v0.1.0
```

```bash
provider "singularity" {
  host    = "localhost"
  version = "~> 0.1.0"
}

resource "singularity_request" "lenfree-demand" {
  request_id   = "sample-request"
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
```

More examples can be found in examples/main.tf.

## Import Resources:

Syntax

```
 terraform import singularity_request.lenfree-run <resource ID>
 terraform import singularity_docker_deploy.test-deploy-2 <resource ID>
```
