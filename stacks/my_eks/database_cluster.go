package my_eks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewDatabaseCluster(stack constructs.Construct, eksCluster eks.Cluster) {
	// DBクラスターの作成
	dbCluster := rds.NewDatabaseCluster(stack, jsii.String("DatabaseCluster"), &rds.DatabaseClusterProps{
		Engine: rds.DatabaseClusterEngine_AuroraPostgres(&rds.AuroraPostgresClusterEngineProps{
			Version: rds.AuroraPostgresEngineVersion_VER_13_6(),
		}),
		ClusterIdentifier: jsii.String("cluster-identifier"),
		InstanceIdentifierBase: jsii.String("db-instance-identifier"),
		Credentials: rds.Credentials_FromGeneratedSecret(jsii.String("postgres"), &rds.CredentialsBaseOptions{
			SecretName: jsii.String("database-secrets"),
		}),
		InstanceProps: &rds.InstanceProps{
			Vpc: eksCluster.Vpc(),
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
	dbCluster.Connections().AllowFrom(eksCluster, ec2.Port_Tcp(jsii.Number(5432)), jsii.String("Allow access to Database clster from EKS cluster"))
}