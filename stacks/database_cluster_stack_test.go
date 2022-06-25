package stacks_test

import (
	"testing"

	"mycdk/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewDatabaseClusterStack(t *testing.T) {
	app := cdk.NewApp(nil)

	// クロススタック参照のデータを用意
	refStack := cdk.NewStack(app, jsii.String("ReferenceStack"), nil)
	vpc := ec2.NewVpc(refStack, jsii.String("VPC"), &ec2.VpcProps{
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
	})
	cluster := eks.NewCluster(refStack, jsii.String("EKSCluster"), &eks.ClusterProps{
		Version: eks.KubernetesVersion_V1_21(),
		Vpc: vpc,
	})

	// テスト対象のスタックとテンプレートを用意
	testStack := stacks.NewDatabaseClusterStack(app, "TestStack", vpc, cluster, nil)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::RDS::DBCluster"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::RDS::DBInstance"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroup"), jsii.Number(1));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::RDS::DBCluster"), map[string]interface{}{
		"Engine": "aurora-postgresql",
		"CopyTagsToSnapshot": true,
		"DatabaseName": "EksDatabaseName",
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
	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"SecurityGroupIngress": []map[string]interface{}{
			{
				"Description": assertions.Match_AnyValue(),
				"FromPort": 80,
				"ToPort": 80,
				"IpProtocol": "tcp",
				"SourceSecurityGroupId": map[string]interface{}{"Fn::ImportValue": assertions.Match_AnyValue()},
			},
		},
		"VpcId": map[string]interface{}{
			"Fn::ImportValue": assertions.Match_AnyValue(),
		},
	})
}
