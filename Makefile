# note: call scripts from /scripts

GO=$(shell which go)
DOCKER=DOCKER_BUILDKIT=1 $(shell which docker)
LINTER=$(shell which golangci-lint)

BUILD_IMAGE_NAME=rticonnexdds-connector-go
RUNTIME_CONTAINER=runtime-container

WORKING_DIR=/rticonnextdds-connector-go/

DOCKER_RUNTIME_CMD=\
		${DOCKER} run -i --rm \
			--name ${RUNTIME_CONTAINER} \
			-v `pwd`:${WORKING_DIR} \
			${BUILD_IMAGE_NAME} \

.DEFAULT_GOAL := test-local

.PHONY: .docker
.docker:
	${DOCKER} build \
		--build-arg working_dir=${WORKING_DIR} \
		-t ${BUILD_IMAGE_NAME} .

.PHONY: test-local
test-local:
	DYLD_LIBRARY_PATH=rticonnextdds-connector/lib/osx-x64 \
	LD_LIBRARY_PATH=rticonnextdds-connector/lib/linux-x64 \
	${GO} test -v -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: test
test: .docker
	${DOCKER_RUNTIME_CMD} test-local

.PHONY: lint-local
lint-local:
	${LINTER} run ./...

.PHONY: lint
lint:
	${DOCKER_RUNTIME_CMD} lint-local
