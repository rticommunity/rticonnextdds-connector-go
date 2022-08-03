FROM golang:1.17

WORKDIR /app

ENV GO111MODULE=on
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/rticonnextdds-connector-go/rticonnextdds-connector/lib/linux-x64

RUN apt-get update && apt-get install -y git 

RUN git clone https://github.com/rticommunity/rticonnextdds-connector-go.git && \
    cd /app/rticonnextdds-connector-go && \
    go build /app/rticonnextdds-connector-go/examples/simple/writer/writer.go

CMD /app/rticonnextdds-connector-go/writer
