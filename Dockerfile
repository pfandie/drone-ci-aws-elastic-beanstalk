FROM golang:1.19-alpine3.15@sha256:eabc3aca6f6c4386369b5b067c9c210aeccd39e76907fa2f8f774fd59d83425a AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /drone-aws-eb/build
ADD . /drone-aws-eb/build
RUN go build -o /release/drone-aws-eb

FROM alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

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
