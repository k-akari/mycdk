package stacks

import (
	"os"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ecrassets "github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	events "github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	targets "github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewNotifyAWSBillingStack(scope constructs.Construct, id string, props *cdk.StackProps) (stack cdk.Stack) {
	stack = cdk.NewStack(scope, &id, props)

	// ********************************************************************************
	// 1. Dockerイメージを作成してECRへプッシュ
	// ********************************************************************************
	// [NewDockerImageAsset](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets#NewDockerImageAsset)
	// [DockerImageAssetProps](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets#DockerImageAssetProps)
	dockerImageAsset := ecrassets.NewDockerImageAsset(stack, jsii.String("NotifyAWSBillingImageAsset"), &ecrassets.DockerImageAssetProps{
		Directory: jsii.String("./images/notifybilling/"), // ビルドコンテキストをリポジトリルートからのパスで指定する
  		File: jsii.String("Dockerfile"), // Dockerfile名をDirectoryプロパティで指定したディレクトリからの相対パスで指定する
		Platform: ecrassets.Platform_LINUX_ARM64(), // 今回はM1 Macbookからビルドを行うのでPlatform_LINUX_ARM64を指定する
  	})
	// ***注意点***
	// ビルドしたイメージのプッシュ先リポジトリはCDKのデフォルトリポジトリです。
	// プッシュ先リポジトリを指定する機能は2022/8/15時点ではありません。
	// 指定したリポジトリにイメージを格納したい場合はサードパーティ製のツール（AWS公式ドキュメントでも説明されている）を使う必要があり、このツールではデフォルトリポジトリへプッシュされたイメージを指定したリポジトリへコピーしてくれます。
	// ただし、このサードパーティツールはTypeScriptで開発されており、Goのライブラリでは用意されていません。
	// 詳しくは[公式ドキュメント](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets#section-readme)をご参照ください。

	// ********************************************************************************
	// 2. ECRへプッシュしたイメージを実行するLambda関数の作成
	// ********************************************************************************
	// 2-1. AWSのコストと使用量を取得するポリシーの作成
	// [NewManagedPolicy](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsiam#NewManagedPolicy)
	// [ManagedPolicyProps](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsiam#ManagedPolicyProps)
	readBillingPolicy := iam.NewManagedPolicy(stack, jsii.String("ReadAWSBillingPolicy"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("read-aws-billing-policy"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
			Statements: &[]iam.PolicyStatement{
				iam.NewPolicyStatement(&iam.PolicyStatementProps{
					Effect: iam.Effect_ALLOW,
					Resources: &[]*string{jsii.String("*")},
					Actions: &[]*string{
						jsii.String("ce:GetCostAndUsage"),
					},
				}),
			},
		}),
	})

	// 2-2. Lambda関数の実行ロールの作成
	// [NewRole](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsiam#NewRole)
	// [RoleProps](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awsiam#RoleProps)
	role := iam.NewRole(stack, jsii.String("NotifyAWSBillingRole"), &iam.RoleProps{
		AssumedBy: iam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &iam.ServicePrincipalOpts{}), // Lambda関数がこのロールを引き受けられるようにする
		RoleName: jsii.String("notify-aws-billing-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{readBillingPolicy,}, // 2-1で作成したAWSのコストと使用量を取得するポリシーをロールに紐づける
  	})

	// 2-3. LINE_NOTIFY_TOKENを環境変数から取得
	accessToken, ok := os.LookupEnv("LINE_NOTIFY_TOKEN") // 秘匿情報なので環境変数にして値を渡している
	if !ok {
		panic("LINE_NOTIFY_TOKEN NOT FOUND") // 環境変数の設定が漏れていた場合にすぐ気づけるようパニックを起こす（これがしたかったのでos.Getenvを使っていない）
	}

	// 2-4. Lambda関数の作成
	// [NewDockerImageFunction](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awslambda#NewDockerImageFunction)
	// [DockerImageFunctionProps](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awslambda#DockerImageFunctionProps)
	handler := lambda.NewDockerImageFunction(stack, jsii.String("DockerImageFunction"), &lambda.DockerImageFunctionProps{
		Code: lambda.DockerImageCode_FromEcr(dockerImageAsset.Repository(), &lambda.EcrImageCodeProps{
			TagOrDigest: dockerImageAsset.AssetHash(),
		}), // このLambda関数が実行するイメージをECRリポジトリとイメージのハッシュ値で指定する
		Architecture: lambda.Architecture_ARM_64(), // 今回はM1 Macbookからビルドしたイメージを実行するので、Lambda関数のアーキテクチャもそれに合わせる
		Environment: &map[string]*string{
			"LINE_NOTIFY_TOKEN": jsii.String(accessToken),
		},
		FunctionName: jsii.String("notify-aws-billing"),
		MemorySize: jsii.Number(1024),
		Role: role, // 2-2で作成したIAMロールを紐づける
		Timeout: cdk.Duration_Seconds(jsii.Number(10)),
	})
	
	// ********************************************************************************
	// 3. Lambda関数を定期実行させるスケジュールの作成
	// ********************************************************************************
	events.NewRule(stack, jsii.String("NotifyAWSBillingRule"), &events.RuleProps{
		Enabled: jsii.Bool(true),
		RuleName: jsii.String("notify-aws-billing"),
		Schedule: events.Schedule_Cron(&events.CronOptions{
			Minute: jsii.String("0"),
			Hour: jsii.String("1"),
			WeekDay: jsii.String("MON"),
		}), // JST で毎週月曜日の AM10:00 に定期実行
		Targets: &[]events.IRuleTarget{
			targets.NewLambdaFunction(handler, &targets.LambdaFunctionProps{
				RetryAttempts: jsii.Number(3),
			}),
		},
	})

	return
}