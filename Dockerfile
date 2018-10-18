FROM golang:1.11


COPY . /go/src/iot-device-repository
WORKDIR /go/src/iot-device-repository

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./iot-device-repository