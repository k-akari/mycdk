package main

import (
	"os"

	"github.com/k-akari/services/test_env/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	app := awscdk.NewApp(nil)

	stacks.NewTestEnvStack(app, "TestEnvStack", &awscdk.StackProps{
		Env: env(),
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
