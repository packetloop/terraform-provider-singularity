go-singularity
--------------

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![Build status](https://circleci.com/gh/lenfree/go-mesos-singularity.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/lenfree/go-mesos-singularity)

[godocs]: https://godoc.org/github.com/lenfree/go-mesos-singularity

A Go binding for Mesos hubspot/Singularity binding. Since I couldn't
manage to find one, hence, write a new one. One of the intention of
having this package is so I could write a Terraform provider to
interface with this.

# Status [ WORK IN PROGRESS ]

## Usage:

Import package
```bash
go get github.com/lenfree/go-mesos-singularity
```

For package dependency management, we use dep:
```bash
go get -u github.com/golang/dep/cmd/dep
```

If new package is required, pls run below command
after go get. For more information about dep, please
visit this URL https://github.com/golang/dep.
```bash
dep ensure
```

Run test:
```bash
make test
```

To maintain codebase quality with static checks and analysis:
```bash
make run
```

Examples:
```go
package main

import (
	"fmt"

	singularity "github.com/lenfree/go-singularity"
)

func main() {
	c := singularity.Config{
		Host: "singularity.net/singularity",
	}
	client := singularity.New(c)
	r, _ := client.GetRequests()
	for _, i := range r {
		body, _ := client.GetRequestByID(i.Request.ID)
		fmt.Println(body)
	}
}
```


## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request