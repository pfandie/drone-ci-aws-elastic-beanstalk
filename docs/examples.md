# Drone Example

- [parameter references](references.md)

## Full Example

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
      region: eu-central-1
      application: awesome-application
      environment: awesome-environment
      version_label: v1.0.0
      description: awesome-version-description
      tags:
        - Project=Project Name
        - Team=MyTeam
        - Department=Web-Team
      update: true
      env_description: false
      wait_for_update: true
      create: true
      preprocess: false
```
