package main

import (
	"context"
	"os"

	"drone-aws-elastic-beanstalk/plugin"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var version = "2.0.0"

func main() {
	app := cli.Command{}
	app.Name = "Drone-CI AWS Elastic Beanstalk plugin"
	app.Usage = "Drone-CI AWS Elastic Beanstalk plugin"
	app.Description = "Deploy applications to AWS Elastic Beanstalk"
	app.Action = run
	app.Version = version

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "aws-access-key",
			Usage:   "aws access key for deployment",
			Sources: cli.EnvVars("PLUGIN_AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"),
		},
		&cli.StringFlag{
			Name:    "aws-secret-key",
			Usage:   "aws secret key for deployment",
			Sources: cli.EnvVars("PLUGIN_AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"),
		},
		&cli.StringFlag{
			Name:    "s3-bucket",
			Usage:   "bucket which contains the deployment file",
			Sources: cli.EnvVars("PLUGIN_BUCKET"),
		},
		&cli.StringFlag{
			Name:    "s3-bucket-key",
			Usage:   "deployment file in bucket including path",
			Sources: cli.EnvVars("PLUGIN_BUCKET_KEY"),
		},
		&cli.StringFlag{
			Name:    "region",
			Usage:   "aws region of the beanstalk application",
			Value:   "eu-central-1",
			Sources: cli.EnvVars("PLUGIN_REGION"),
		},
		&cli.StringFlag{
			Name:    "eb-app-name",
			Usage:   "beanstalk application name",
			Sources: cli.EnvVars("PLUGIN_APPLICATION"),
		},
		&cli.StringFlag{
			Name:    "eb-env-name",
			Usage:   "name of environment to update in aws eb app",
			Sources: cli.EnvVars("PLUGIN_ENVIRONMENT"),
		},
		&cli.StringFlag{
			Name:    "eb-version-label",
			Usage:   "version name of aws eb application",
			Sources: cli.EnvVars("PLUGIN_VERSION_LABEL"),
		},
		&cli.StringFlag{
			Name:    "eb-description",
			Usage:   "description appended to the app-version",
			Sources: cli.EnvVars("PLUGIN_DESCRIPTION"),
		},
		&cli.StringSliceFlag{
			Name:    "eb-tags",
			Usage:   "tags appended to the app-version",
			Sources: cli.EnvVars("PLUGIN_TAGS"),
		},
		&cli.BoolFlag{
			Name:    "eb-env-update",
			Usage:   "update environment with given app-version - defaults to true",
			Value:   true,
			Sources: cli.EnvVars("PLUGIN_UPDATE"),
		},
		&cli.BoolFlag{
			Name:    "eb-env-description",
			Usage:   "add description to environment - defaults to true",
			Value:   true,
			Sources: cli.EnvVars("PLUGIN_ENV_DESCRIPTION"),
		},
		&cli.BoolFlag{
			Name:    "eb-wait-for-update",
			Usage:   "wait for updated environment completed, requires update to be true - defaults to false",
			Value:   false,
			Sources: cli.EnvVars("PLUGIN_WAIT_FOR_UPDATE"),
		},
		&cli.BoolFlag{
			Name:    "eb-auto-create",
			Usage:   "create new app-environment, if not already exists - defaults to false",
			Value:   false,
			Sources: cli.EnvVars("PLUGIN_CREATE"),
		},
		&cli.BoolFlag{
			Name:    "eb-process",
			Usage:   "validate and preprocess manifest - defaults to false",
			Value:   false,
			Sources: cli.EnvVars("PLUGIN_PREPROCESS"),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to initialize Plugin")
	}
}

func run(ctx context.Context, c *cli.Command) error {
	p := plugin.Plugin{
		AccessKey:      c.String("aws-access-key"),
		SecretKey:      c.String("aws-secret-key"),
		Bucket:         c.String("s3-bucket"),
		BucketKey:      c.String("s3-bucket-key"),
		Region:         c.String("region"),
		AppName:        c.String("eb-app-name"),
		EnvName:        c.String("eb-env-name"),
		Version:        c.String("eb-version-label"),
		Description:    c.String("eb-description"),
		EnvDescription: c.Bool("eb-env-description"),
		Tags:           c.StringSlice("eb-tags"),
		Update:         c.Bool("eb-env-update"),
		WaitForUpdate:  c.Bool("eb-wait-for-update"),
		Create:         c.Bool("eb-auto-create"),
		PreProcess:     c.Bool("eb-process"),
	}

	return p.Exec()
}
