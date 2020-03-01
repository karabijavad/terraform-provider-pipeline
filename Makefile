PIPELINE_VERSION = 0.37.0
OPENAPI_GENERATOR_VERSION = v4.2.2

terraform-provider-pipeline: *.go
	go build -o terraform-provider-pipeline

tf: terraform-provider-pipeline
	terraform init
	terraform apply

.PHONY: generate-pipeline-client
generate-pipeline-client:
	curl https://raw.githubusercontent.com/banzaicloud/pipeline/${PIPELINE_VERSION}/apis/pipeline/pipeline.yaml > pipeline-openapi.yaml
	rm -rf .gen/pipeline
	docker run --rm \
		-v ${CURDIR}:/local \
		openapitools/openapi-generator-cli:${OPENAPI_GENERATOR_VERSION} \
		generate \
		--additional-properties packageName=pipeline \
		--additional-properties withGoCodegenComment=true \
		-i /local/pipeline-openapi.yaml \
		-g go \
		-o /local/.gen/pipeline
	sudo chown -R $(shell id -u).$(shell id -g) local
	sudo chown -R $(shell id -u).$(shell id -g) .gen
	echo "package pipeline\n\nconst PipelineVersion = \"${PIPELINE_VERSION}\"" > .gen/pipeline/version.go
	sed 's#jsonCheck = .*#jsonCheck = regexp.MustCompile(`(?i:(?:application|text)/(?:(?:vnd\\.[^;]+\\+)|(?:problem\\+))?json)`)#' .gen/pipeline/client.go > .gen/pipeline/client.go.new
	mv .gen/pipeline/client.go.new .gen/pipeline/client.go

all: terraform-provider-pipeline tf
