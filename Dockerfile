FROM golang:alpine3.20@sha256:f591145352ef7cd7d7e2b4e1d4a6fd4dd2ac72c405b689d6e8147339105a9e3a AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /drone-aws-eb/build
ADD . /drone-aws-eb/build
RUN go build -o /release/drone-aws-eb

FROM alpine:3.20@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

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
