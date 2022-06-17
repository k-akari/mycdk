package main

import (
	"os"

	"github.com/k-akari/services/eks_cluster/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	jsii "github.com/aws/jsii-runtime-go"
)

func main() {
	app := cdk.NewApp(nil)

	stacks.NewEksClusterStack(app, "EksClusterStack", &cdk.StackProps{
		Env: env(),
	})

	app.Synth(nil)
}

func env() *cdk.Environment {
	return &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
