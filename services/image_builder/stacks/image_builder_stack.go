package stacks

import (
	"os"

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

	// Create a IanRole to push image to ECR repository
	role := iam.NewRole(stack, jsii.String("RoleCodeBuildImageBuilder"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), &iam.ServicePrincipalOpts{}),
		Description: jsii.String("Iam Role for CodeBuild Project to push image to ECR repository"),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-codebuild-for-image-builder"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryPowerUser")),
		},
    })

	// Create a CodeBuild project
	codebuild.NewProject(stack, jsii.String("ProjectEntExample"), &codebuild.ProjectProps{
		BuildSpec: codebuild.BuildSpec_FromSourceFilename(jsii.String("buildspec.yml")),
		ConcurrentBuildLimit: jsii.Number(1),
		Environment: &codebuild.BuildEnvironment{
			ComputeType: codebuild.ComputeType_SMALL,
			Privileged: jsii.Bool(true),
		},
		EnvironmentVariables: &map[string]*codebuild.BuildEnvironmentVariable{
			"AWS_ACCOUNT": &codebuild.BuildEnvironmentVariable{Value: os.Getenv("CDK_DEFAULT_ACCOUNT")},
			"AWS_REGION": &codebuild.BuildEnvironmentVariable{Value: os.Getenv("CDK_DEFAULT_REGION")},
		},
		ProjectName: jsii.String("ImageBuilerEntExample"),
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

	// Create a repository to store image
	ecr.NewRepository(stack, jsii.String("RepositoryImageBuilder"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{
			{
				MaxImageCount: jsii.Number(1),
			},
		},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("ent-example"),
	})

	return stack
}