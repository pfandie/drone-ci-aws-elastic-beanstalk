# Param Reference

**aws_access_key**\
AWS access key

**aws_secret_key**\
AWS access key

**bucket**\
s3 bucket containing the source bundle

**bucket_key**\
name/key of the source bundle (*including path*)

**region**\
region of the application, defaults to eu-central-1

**application**\
application name, required

**environment**\
environment name, required if update is true

**version_label**\
version label to identify deployment

**description**\
describes the application version, optional

**tags**\
tags to add to the application version

**update**\
flag to enable/disable environment updates, defaults to true

**env_description**\
flag to use description also for the environment, defaults to true

**wait_for_update**\
flag to wait in build step until update is completed, defaults to false

**create**\
flag to create a new application, if not exists, defaults to false

**preprocess**\
flag to validate/process manifest, defaults to false 
