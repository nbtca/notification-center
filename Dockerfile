# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY consolefixfunc consolefixfunc

RUN CGO_ENABLED=0 GOOS=linux go build -o /webhook

##
## Run the tests in the container
##

FROM build-stage AS run-test-stage
RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM ubuntu AS build-release-stage

WORKDIR /

COPY --from=build-stage /webhook /webhook

EXPOSE 18080

# call wenhook with /config/config.json
ENTRYPOINT ["/webhook","/config/config.json"]