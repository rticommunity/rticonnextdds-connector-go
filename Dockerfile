FROM golang:1.17

ARG working_dir

WORKDIR ${working_dir}
VOLUME [ ${working_dir} ]

ENV GO111MODULE=on

RUN apt-get update && apt-get install -y \
        curl \
        git \
    && rm -rf /var/lib/apt/lists/*

COPY scripts/install_golangci-lint.sh .
RUN ./install_golangci-lint.sh

COPY go.* .
RUN go mod download

ENTRYPOINT [ "make" ]
