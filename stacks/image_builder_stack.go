package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewImageBuilderStack(scope constructs.Construct, id string, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// DockerイメージをビルドしてECRへプッシュするIamRoleを作成
	role := iam.NewRole(stack, jsii.String("EKSCodeBuildRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), &iam.ServicePrincipalOpts{}),
		Description: jsii.String("Iam Role for CodeBuild Project to push image to ECR repository"),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-codebuild-for-image-builder"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryPowerUser")),
		},
    })

	// 指定のリポジトリのmainブランチへの更新をトリガーとして、DockerイメージをビルドしECRリポジトリへイメージをプッシュするプロジェクトを作成
	codebuild.NewProject(stack, jsii.String("EKSImageBuildProject"), &codebuild.ProjectProps{
		BuildSpec: codebuild.BuildSpec_FromSourceFilename(jsii.String("buildspec.yml")),
		ConcurrentBuildLimit: jsii.Number(1),
		Environment: &codebuild.BuildEnvironment{
			ComputeType: codebuild.ComputeType_SMALL,
			Privileged: jsii.Bool(true),
		},
		EnvironmentVariables: &map[string]*codebuild.BuildEnvironmentVariable{
			"AWS_ACCOUNT": {Value: sprops.Env.Account},
			"AWS_REGION": {Value: sprops.Env.Region},
		},
		ProjectName: jsii.String("EKSImageBuildProject"),
		QueuedTimeout: cdk.Duration_Hours(jsii.Number(1)),
		Role: role,
		Source: codebuild.Source_GitHub(&codebuild.GitHubSourceProps{
			Owner: jsii.String("k-akari"),
			Repo: jsii.String("ent-example"),
			BranchOrRef: jsii.String("main"),
			CloneDepth: jsii.Number(1),
			Webhook: jsii.Bool(true),
			WebhookFilters: &[]codebuild.FilterGroup{
				codebuild.FilterGroup_InEventOf(codebuild.EventAction_PUSH).AndBranchIs(jsii.String("main")),
			},
		}),
		Timeout: cdk.Duration_Minutes(jsii.Number(20)),
	})

	// ビルドしたアプリケーションイメージを格納するリポジトリの作成
	ecr.NewRepository(stack, jsii.String("EKSAppImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{
			{
				MaxImageCount: jsii.Number(1),
			},
		},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-app"),
	})

	return stack
}