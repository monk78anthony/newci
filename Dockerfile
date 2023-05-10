FROM 767653220718.dkr.ecr.us-east-1.amazonaws.com/golang:1.20.4-alpine3.17 AS build

RUN go env -w GOPROXY=direct
RUN apk update && apk add --no-cache git
RUN mkdir -p /go/src
WORKDIR /go/src
COPY . /go/src
RUN go get "github.com/go-sql-driver/mysql"
COPY . .
RUN CGO_ENABLED=0 go build -o /src/demo

ENTRYPOINT ["/src/demo"]
