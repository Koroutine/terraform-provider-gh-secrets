clean:
	rm -rf .terraform*

build:
	go build -o terraform-provider-gh-secrets

install:
	mv terraform-provider-gh-secrets ~/.terraform.d/plugins/koroutine/tf/gh-secrets/0.1/darwin_arm64
	terraform init