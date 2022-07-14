package my_eks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewDBMigrator(stack constructs.Construct, repo ecr.Repository, dbCluster rds.DatabaseCluster, props *cdk.StackProps) {
	// DockerイメージをプルするIamRoleを作成
	pullImagePolicy := iam.NewManagedPolicy(stack, jsii.String("PullImagePolicyForDBMigrator"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("pull-image-policy-for-db-migrator"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{repo.RepositoryArn()},
    				Actions: &[]*string{
    					jsii.String("ecr:BatchGetImage"),
    					jsii.String("ecr:GetDownloadUrlForLayer"),
					},
				}),
			},
    	}),
	})
	role := iam.NewRole(stack, jsii.String("DBMigratorRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), &iam.ServicePrincipalOpts{}),
		Description: jsii.String("Iam Role for CodeBuild Project to pull image from ECR repository"),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-for-db-migrator"),
		ManagedPolicies: &[]iam.IManagedPolicy{pullImagePolicy,},
    })

	// CodeBuild Project用のセキュリティグループの作成
	sgProject := ec2.NewSecurityGroup(stack, jsii.String("SecurityGroupForDBMigrator"), &ec2.SecurityGroupProps{
		Vpc: dbCluster.Vpc(),
		AllowAllOutbound: jsii.Bool(true),
		Description: jsii.String("Security Group for DB Migrator Project"),
		SecurityGroupName: jsii.String("SecurityGroupForDBMigrator"),
	})

	// ECRからプルしたイメージを元にDBマイグレートを行うプロジェクトを作成
	codebuild.NewProject(stack, jsii.String("DBMigratorProject"), &codebuild.ProjectProps{
		//AllowAllOutbound: jsii.Bool(true),
		BuildSpec: codebuild.BuildSpec_FromObject(&map[string]interface{}{
			"version": jsii.String("0.2"),
			"phases": map[string]interface{}{
				"build": map[string]interface{}{
					"commands": []*string{
						jsii.String("whoami"),
						jsii.String("pwd"),
						jsii.String("printenv"),
						jsii.String("/app/cmd"),
					},
				},
			},
		}),
		ConcurrentBuildLimit: jsii.Number(1),
		Environment: &codebuild.BuildEnvironment{
			BuildImage: codebuild.LinuxBuildImage_FromEcrRepository(repo, jsii.String("latest")),
			ComputeType: codebuild.ComputeType_SMALL,
			Privileged: jsii.Bool(true),
		},
		EnvironmentVariables: &map[string]*codebuild.BuildEnvironmentVariable{
			"AWS_ACCOUNT": {Value: props.Env.Account},
			"AWS_REGION": {Value: props.Env.Region},
			// [ref] https://docs.aws.amazon.com/ja_jp/codebuild/latest/userguide/build-spec-ref.html#build-spec.env.secrets-manager
			"DB_HOST": {Value: *dbCluster.Secret().SecretFullArn() + ":host", Type: codebuild.BuildEnvironmentVariableType_SECRETS_MANAGER},
			"DB_PORT": {Value: *dbCluster.Secret().SecretFullArn() + ":port", Type: codebuild.BuildEnvironmentVariableType_SECRETS_MANAGER},
			"DB_PASSWORD": {Value: *dbCluster.Secret().SecretFullArn() + ":password", Type: codebuild.BuildEnvironmentVariableType_SECRETS_MANAGER},
			"DB_USER": {Value: *dbCluster.Secret().SecretFullArn() + ":username", Type: codebuild.BuildEnvironmentVariableType_SECRETS_MANAGER},
			"DB_NAME": {Value: *dbCluster.Secret().SecretFullArn() + ":dbname", Type: codebuild.BuildEnvironmentVariableType_SECRETS_MANAGER},
		},
		ProjectName: jsii.String("DBMigratorProject"),
		QueuedTimeout: cdk.Duration_Hours(jsii.Number(1)),
		Role: role,
		Vpc: dbCluster.Vpc(),
		SubnetSelection: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
		},
		SecurityGroups: &[]ec2.ISecurityGroup{sgProject},
		Timeout: cdk.Duration_Minutes(jsii.Number(20)),
	})

	// DBMigrateプロジェクトからDBクラスターへのアクセスを許可する
	for _, sg := range *dbCluster.SecurityGroups() {
		sg.AddIngressRule(
			ec2.Peer_SecurityGroupId(sgProject.SecurityGroupId(), jsii.String("")),
			ec2.Port_Tcp(jsii.Number(5432)),
			jsii.String("Allow access to DB Cluster from CodeBuild Project"),
			jsii.Bool(true),
		)
	}
}