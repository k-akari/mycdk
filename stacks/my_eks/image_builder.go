package my_eks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewImageBuilder(stack constructs.Construct, props *cdk.StackProps) (repoMigration ecr.Repository) {
	// ビルドしたイメージを格納するECRリポジトリの作成
	repoApp := ecr.NewRepository(stack, jsii.String("EKSAppImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{{MaxImageCount: jsii.Number(1),},},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-app"),
	})
	repoMigration = ecr.NewRepository(stack, jsii.String("EKSMigrationImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{{MaxImageCount: jsii.Number(1),},},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-migration"),
	})
	repoWeb := ecr.NewRepository(stack, jsii.String("EKSWebImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{{MaxImageCount: jsii.Number(1),},},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-web"),
	})

	// DockerイメージをビルドしてECRリポジトリへプッシュするIamRoleを作成
	// [ref] https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-push.html#image-push-iam
	pushImagePolicy := iam.NewManagedPolicy(stack, jsii.String("PushImagePolicyForImageBuilder"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("push-image-policy-for-image-builder"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{repoApp.RepositoryArn(), repoMigration.RepositoryArn(), repoWeb.RepositoryArn()},
    				Actions: &[]*string{
    					jsii.String("ecr:CompleteLayerUpload"),
    					jsii.String("ecr:UploadLayerPart"),
    					jsii.String("ecr:InitiateLayerUpload"),
    					jsii.String("ecr:BatchCheckLayerAvailability"),
    					jsii.String("ecr:PutImage"),
					},
				}),
				iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{jsii.String("*")},
    				Actions: &[]*string{
    					jsii.String("ecr:GetAuthorizationToken"),
					},
				}),
			},
    	}),
	})
	role := iam.NewRole(stack, jsii.String("ImageBuilderRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), &iam.ServicePrincipalOpts{}),
		Description: jsii.String("Iam Role for CodeBuild Project to push image to ECR repository"),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-codebuild-for-image-builder"),
		ManagedPolicies: &[]iam.IManagedPolicy{pushImagePolicy,},
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
			"AWS_ACCOUNT": {Value: props.Env.Account},
			"AWS_REGION": {Value: props.Env.Region},
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

	return
}