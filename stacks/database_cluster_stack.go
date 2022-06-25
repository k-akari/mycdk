package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewDatabaseClusterStack(scope constructs.Construct, id string, vpc ec2.Vpc, eksCluster eks.Cluster, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// DBクラスターの作成
	dbCluster := rds.NewDatabaseCluster(stack, jsii.String("DatabaseCluster"), &rds.DatabaseClusterProps{
		Engine: rds.DatabaseClusterEngine_AuroraPostgres(&rds.AuroraPostgresClusterEngineProps{
			Version: rds.AuroraPostgresEngineVersion_VER_13_6(),
		}),
		InstanceProps: &rds.InstanceProps{
			Vpc: vpc,
			VpcSubnets: &ec2.SubnetSelection{
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
			},
			AllowMajorVersionUpgrade: jsii.Bool(false),
			AutoMinorVersionUpgrade: jsii.Bool(true),
			DeleteAutomatedBackups: jsii.Bool(false),
			EnablePerformanceInsights: jsii.Bool(false),
			InstanceType: ec2.NewInstanceType(jsii.String("t3.medium")),
			PubliclyAccessible: jsii.Bool(false),
		},
		Instances: jsii.Number(1),
		Port: jsii.Number(5432),
		DefaultDatabaseName: jsii.String("EksDatabaseName"),
		DeletionProtection: jsii.Bool(false),
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
	})

	// EKSクラスターからDBクラスターへのアクセスを許可する
	for _, sg := range *dbCluster.SecurityGroups() {
		sg.AddIngressRule(
			ec2.Peer_SecurityGroupId(eksCluster.ClusterSecurityGroupId(), jsii.String("")),
			ec2.Port_Tcp(jsii.Number(80)),
			jsii.String("Allow access from EKS Node Group"),
			jsii.Bool(true),
		)
	}

	return stack
}