FROM golang:1.14-buster

RUN apt-get update && apt-get install -y curl git-lfs git && git lfs install

RUN git clone https://github.com/rticommunity/rticonnextdds-connector-go.git /rticonnextdds-connector-go

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0

RUN cd /rticonnextdds-connector-go && go mod download

CMD golangci-lint run -v && bash run_test.sh

WORKDIR /rticonnextdds-connector-go
