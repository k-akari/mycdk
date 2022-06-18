package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewDatabaseClusterStack(scope constructs.Construct, id string, vpc ec2.Vpc, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	rds.NewDatabaseCluster(stack, jsii.String("EksDatabaseCluster"), &rds.DatabaseClusterProps{
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
		DefaultDatabaseName: jsii.String("EksDatabaseClusterName"),
		DeletionProtection: jsii.Bool(false),
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
	})

	return stack
}