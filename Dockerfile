FROM golang:alpine3.20@sha256:9a31ef0803e6afdf564edc8ba4b4e17caed22a0b1ecd2c55e3c8fdd8d8f68f98 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /drone-aws-eb/build
ADD . /drone-aws-eb/build
RUN go build -o /release/drone-aws-eb

FROM alpine:3.21@sha256:e323a465c03a31ad04374fc7239144d0fd4e2b92da6e3e0655580476d3a84621

LABEL org.label-schema.name="drone-aws-eb" \
    org.label-schema.description="A Drone plugin to deploy an Application to AWS Elastic Beanstalk" \
    org.label-schema.vcs-url="https://github.com/pfandie/drone-ci-aws-elastic-beanstalk" \
    org.label-schema.vendor="pfandie" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.docker.cmd="docker run --rm -e PLUGIN_AWS_ACCESS_KEY=<redacted-token> -e PLUGIN_AWS_SECRET_KEY=<redacted-secret> \
-e PLUGIN_BUCKET=<bucket> -e PLUGIN_BUCKET_KEY=<bucket-key> -e PLUGIN_APPLICATION=<beanstalk-app-name> -e PLUGIN_ENVIRONMENT=<beanstalk-env-name> \
-e PLUGIN_VERSION_LABEL=<version-label> -v $(pwd):$(pwd) -w $(pwd) pfandie/drone-aws-eb"

COPY --from=builder \
    /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs \
    /release/drone-aws-eb /bin/

ENTRYPOINT ["/bin/drone-aws-eb"]
