package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
)

var (
	AppVersion = "0.0.0-dev.0"
)

func newSession() (*session.Session, error) {

	sessCfg := aws.NewConfig().
		WithCredentials(credentials.NewEnvCredentials()).
		WithRegion(os.Getenv("AWS_REGION")).
		WithHTTPClient(http.DefaultClient).
		WithMaxRetries(aws.UseServiceDefaultRetries).
		WithLogger(aws.NewDefaultLogger()).
		WithLogLevel(aws.LogOff).
		WithEndpointResolver(endpoints.DefaultResolver())

	sess, err := session.NewSession(sessCfg)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func GetSecret(sess *session.Session, name string) (string, error) {

	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	value, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	if value.SecretString != nil {
		return *value.SecretString, nil
	}

	bytes := make([]byte, base64.StdEncoding.DecodedLen(len(value.SecretBinary)))
	len, err := base64.StdEncoding.Decode(bytes, value.SecretBinary)
	if err != nil {
		return "", err
	}

	return string(bytes[:len]), nil
}

func main() {
	app := cli.NewApp()
	app.Name = "aws-secret"
	app.Usage = "Fetch AWS Secrets Manager secrets"
	app.Version = AppVersion
	app.UsageText = fmt.Sprintf("%s <name>", app.Name)

	app.Action = func(c *cli.Context) error {
		name := c.Args().First()
		if name == "" {
			cli.ShowAppHelpAndExit(c, 1)
			return nil
		}

		sess, err := newSession()
		if err != nil {
			return err
		}

		value, err := GetSecret(sess, name)
		if err != nil {
			return err
		}

		fmt.Println(value)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
