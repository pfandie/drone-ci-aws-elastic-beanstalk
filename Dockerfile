FROM golang:alpine3.20@sha256:09742590377387b931261cbeb72ce56da1b0d750a27379f7385245b2b058b63a AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /drone-aws-eb/build
ADD . /drone-aws-eb/build
RUN go build -o /release/drone-aws-eb

FROM alpine:3.20@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d

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
