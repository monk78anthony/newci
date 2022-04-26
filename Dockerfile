FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/demo
ENTRYPOINT ["/bin/demo"]
