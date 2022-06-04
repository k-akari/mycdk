package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewImageBuilderStack(scope constructs.Construct, id string, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// Create 1 VPC with 1 public subnet
	vpc := ec2.NewVpc(stack, jsii.String("VPCImageBuilder"), &ec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
		MaxAzs: jsii.Number(1),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-image-builder"),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
		},
		VpcName: jsii.String("vpc-for-image-builder"),
	})

	// Create a CodeBuild project
	codebuild.NewProject(stack, jsii.String("ProjectEntExample"), &codebuild.ProjectProps{
		AllowAllOutbound: jsii.Bool(true),
		BuildSpec: codebuild.BuildSpec_FromObject(&map[string]interface{}{
			"version": jsii.String("0.2"),
		}),
		Environment: &codebuild.BuildEnvironment{
			ComputeType: codebuild.ComputeType_SMALL,
			Privileged: jsii.Bool(true),
		},
		ProjectName: jsii.String("ImageBuilerEntExample"),
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
		SubnetSelection: &ec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
		Vpc: vpc,
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