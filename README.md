terraform-provider-singularity
===============================

A terraform provider to manage Mesos hubspot/Singularity objects.

[![CircleCI](https://circleci.com/gh/packetloop/terraform-provider-singularity.svg?style=svg)](https://circleci.com/gh/packetloop/terraform-provider-singularity)

Usage:
------

```bash
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

More examples can be found in examples/main.tf.

# Import Resource 
-----------------

Syntax

```
 terraform import singularity_request.lenfree-run <resource ID>
 
```

For example:
```
 terraform import singularity_request.lenfree-run lenfree--atlas-feedback-s3--backup-artifacts                                                                                              
```

