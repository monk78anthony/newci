FROM golang:1.14-alpine AS build

RUN mkdir /src
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /src/demo

FROM scratch
COPY --from=build /src/demo /src/demo

ENTRYPOINT ["/src/demo"]
