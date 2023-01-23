FROM 767653220718.dkr.ecr.us-east-1.amazonaws.com/golang:1.18.2-alpine AS build

RUN mkdir /src
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /src/demo

FROM scratch
COPY --from=build /src/demo /src/demo

ENTRYPOINT ["/src/demo"]
