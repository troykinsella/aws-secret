package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	AppVersion = "0.0.0-dev.0"
)

func GetSecret(name string) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

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

		value, err := GetSecret(name)
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
