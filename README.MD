# Drone AWS Elastic Beanstalk Plugin

A Drone-CI plugin to deploy and update zipped Applications to
[AWS Elastic Beanstalk](https://aws.amazon.com/elasticbeanstalk/)

## Example Usage

Following examples shows the config in a .drone.yml file

### Simplest example

```yaml
kind: pipeline
name: default

steps:
  - name: deploy to aws-beanstalk
    image: pfandie/drone-aws-eb
    environment:
      AWS_ACCESS_KEY_ID:
        from_secret: aws_access_key_id
      AWS_SECRET_ACCESS_KEY:
        from_secret: aws_secret_access_key
    settings:
      bucket: application-storage
      bucket_key: bundle.zip
      application: awesome-application
      environment: awesome-environment
      version_label: v1.0.0
```

- [more examples](docs/examples.md)
- [parameter references](docs/references.md)

## Build


```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o release/drone-aws-eb
```

## Docker image:

```bash
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --label org.label-schema.version=$(git tag | sort -V | tail -1 | tr -d 'v') \
  -t pfandie/drone-aws-eb .
```

## Usage

```bash
 docker run --rm \
  -e PLUGIN_AWS_ACCESS_KEY="<redacted-token>" \
  -e PLUGIN_AWS_SECRET_KEY="<redacted-secret>" \
  -e PLUGIN_BUCKET="<bucket>" \
  -e PLUGIN_BUCKET_KEY="<bucket-key>" \
  -e PLUGIN_APPLICATION="<beanstalk-app-name>" \
  -e PLUGIN_ENVIRONMENT="<beanstalk-env-name>" \
  -e PLUGIN_VERSION_LABEL="<version-label>" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  pfandie/drone-aws-eb
```
