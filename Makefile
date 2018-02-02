.PHONY: test
test: dep
	HOST=$(HOST) TF_ACC=$(TF_ACC) go test -v ./...

.PHONY: vendor
vendor:
	dep ensure

.PHONY: dep
dep:
ifndef HOST
$(error HOST is not set)
endif

.PHONY: build
build:
	go build -o examples/terraform-provider-singularity

.PHONY: tf-init
tf-init:
	terraform init
	mv .terraform /tmp
	mv terraform.tfstate /tmp