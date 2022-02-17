package plugin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	AccessKey      string
	SecretKey      string
	Bucket         string
	BucketKey      string
	Region         string
	AppName        string
	EnvName        string
	EnvDescription bool
	Version        string
	Description    string
	Tags           []string
	Update         bool
	WaitForUpdate  bool
	Create         bool
	PreProcess     bool
}

var (
	desiredStatus = "Ready"
	actualStatus  = ""
)

func (p Plugin) Exec() error {
	p.preValidateBasicValues()

	log.WithFields(log.Fields{
		"environment":        p.EnvName,
		"app-name":           p.AppName,
		"version-label":      p.Version,
		"source-bundle":      fmt.Sprintf("%s/%s", p.Bucket, p.BucketKey),
		"region":             p.Region,
		"update-environment": p.Update,
		"wait-for-update":    p.WaitForUpdate,
		"auto-create":        p.Create,
		"process":            p.PreProcess,
	}).Info("Trying to create/update with:")

	// init aws conf
	cfg := p.getAwsConf()
	bs := elasticbeanstalk.NewFromConfig(cfg)

	s3Location := &types.S3Location{
		S3Bucket: aws.String(p.Bucket),
		S3Key:    aws.String(p.BucketKey),
	}

	appVersionInput := &elasticbeanstalk.CreateApplicationVersionInput{
		VersionLabel:          aws.String(p.Version),
		ApplicationName:       aws.String(p.AppName),
		AutoCreateApplication: aws.Bool(p.Create),
		Process:               aws.Bool(p.PreProcess),
		SourceBundle:          s3Location,
	}

	if p.Description != "" {
		appVersionInput.Description = aws.String(p.Description)
	}

	if len(p.Tags) > 0 {
		var tags []types.Tag
		for _, t := range p.Tags {
			s := strings.Split(t, "=")
			tag := types.Tag{
				Key:   aws.String(s[0]),
				Value: aws.String(s[1]),
			}
			tags = append(tags, tag)
		}

		appVersionInput.Tags = tags
	}

	// creates new beanstalk application version
	_, err := bs.CreateApplicationVersion(context.TODO(), appVersionInput)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error occurred while creating application: ")
		return err
	}

	if p.Update == true {
		envInput := &elasticbeanstalk.DescribeEnvironmentsInput{
			ApplicationName:  aws.String(p.AppName),
			EnvironmentNames: []string{*aws.String(p.EnvName)},
		}

		describedEnv, err := bs.DescribeEnvironments(
			context.TODO(),
			envInput,
		)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("An error occurred while env describe:")
		}

		// check if environment is in Ready status
		envStatus := fmt.Sprintf("%v", describedEnv.Environments[0].Status)
		if envStatus != desiredStatus {
			log.WithFields(log.Fields{
				"status": envStatus,
			}).Fatal("Cannot apply Update... Environment is not in 'Ready' state")
		}

		if p.EnvName == "" {
			log.Fatal("Environment-Name must not be empty!")
		}

		updateEnvInput := &elasticbeanstalk.UpdateEnvironmentInput{
			VersionLabel:    aws.String(p.Version),
			ApplicationName: aws.String(p.AppName),
			EnvironmentName: aws.String(p.EnvName),
		}

		if p.EnvDescription {
			updateEnvInput.Description = aws.String(p.Description)
		}

		// updates given environment with application version
		_, err = bs.UpdateEnvironment(context.TODO(), updateEnvInput)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("An error occurred while updating environment: ")
		}

		// wait until environment reaches 'Ready' state
		if p.WaitForUpdate {
			for true {
				waitForEnv, _ := bs.DescribeEnvironments(
					context.TODO(),
					envInput,
				)
				actualStatus = fmt.Sprintf("%v", waitForEnv.Environments[0].Status)

				if desiredStatus == actualStatus {
					log.Info("Environment update finished.")
					break
				}

				log.WithFields(log.Fields{
					"status": actualStatus,
				}).Info("Environment Status:")

				// wait 10 seconds for next status-check
				time.Sleep(10 * time.Second)
			}
		}
	}

	log.WithFields(log.Fields{
		"application": p.AppName,
	}).Info("Deployment finished")

	return nil
}

// preValidateBasicValues
// checks basic values are set
func (p Plugin) preValidateBasicValues() {
	if p.AppName == "" {
		log.Fatal("Application Name must not be empty!")
	}

	if p.Update == true && p.EnvName == "" {
		log.Fatal("Environment Name must not be empty when update is set to true!")
	}

	if p.Bucket == "" || p.BucketKey == "" {
		log.Fatal("Fatal, Bucket or Bucket-Key must not be empty!")
	}

	if p.Version == "" {
		log.Fatal("Version Label must not be empty!")
	}
}

// getAwsConf
// returns an AWS Config
func (p Plugin) getAwsConf() aws.Config {
	var err error
	cfg := aws.Config{}

	if p.AccessKey != "" || p.SecretKey != "" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider(
			p.AccessKey, p.SecretKey, "")
	} else {
		log.Info("One or more Keys/Secrets are missing, fallback to Environment/Instance Profile")
	}

	cfg, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(cfg.Credentials),
		config.WithRegion(p.Region))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("A problem occurred while loading config")
	}

	return cfg
}
