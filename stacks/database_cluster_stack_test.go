package stacks_test

import (
	"testing"

	"mycdk/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewDatabaseClusterStack(t *testing.T) {
	app := cdk.NewApp(nil)
	networkStack := cdk.NewStack(app, jsii.String("NetworkStack"), nil)
	vpc := ec2.NewVpc(networkStack, jsii.String("VPC"), &ec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
		MaxAzs: jsii.Number(3),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-eks-cluster-public"),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-eks-cluster-private-with-nat"),
				SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
			},
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-eks-cluster-private-isolated"),
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
			},
		},
		VpcName: jsii.String("vpc-for-eks-cluster"),
	})

	databaseClusterStack := stacks.NewDatabaseClusterStack(app, "MyStack", vpc, nil)

	template := assertions.Template_FromStack(databaseClusterStack)

	template.ResourceCountIs(jsii.String("AWS::RDS::DBCluster"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::RDS::DBInstance"), jsii.Number(1));

	template.HasResourceProperties(jsii.String("AWS::RDS::DBCluster"), map[string]interface{}{
		"Engine": "aurora-postgresql",
		"CopyTagsToSnapshot": true,
		"DatabaseName": "EksDatabaseClusterName",
		"DeletionProtection": false,
		"EngineVersion": "13.6",
		"Port": 5432,
	})
	template.HasResourceProperties(jsii.String("AWS::RDS::DBInstance"), map[string]interface{}{
		"DBInstanceClass": "db.t3.medium",
		"AllowMajorVersionUpgrade": false,
		"AutoMinorVersionUpgrade": true,
		"DeleteAutomatedBackups": false,
		"EnablePerformanceInsights": false,
		"Engine": "aurora-postgresql",
		"EngineVersion": "13.6",
		"PubliclyAccessible": false,
	})
}
