package main

import (
	"os"

	"mycdk/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	jsii "github.com/aws/jsii-runtime-go"
)

func main() {
	app := cdk.NewApp(nil)
	props := &cdk.StackProps{Env: env(),}

	// EKSインフラの構築
	stacks.NewMyEKSStack(app, "MyEKSStack", props)
	//

	// GitHubからOIDC認証でAWSへアクセスするための権限設定
	stacks.NewGitHubActionsStack(app, "GitHubActionsStack", props)
	//

	// SSM接続でログインできて、インターネットアクセスが可能なEC2インスタンスを作成
	stacks.NewMyEc2Stack(app, "MyEc2Stack", props)
	//
	
	app.Synth(nil)
}

func env() *cdk.Environment {
	return &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
