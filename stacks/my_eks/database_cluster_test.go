package my_eks_test

import (
	myeks "mycdk/stacks/my_eks"
	"os"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewDatabaseCluster(t *testing.T) {
	app := cdk.NewApp(nil)

	// テスト対象のスタックテンプレートを用意
	testStack := cdk.NewStack(app, jsii.String("TestStack"), &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},})
	vpc := ec2.NewVpc(testStack, jsii.String("VPC"), &ec2.VpcProps{
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
	cluster := eks.NewCluster(testStack, jsii.String("EKSCluster"), &eks.ClusterProps{
		Version: eks.KubernetesVersion_V1_21(),
		Vpc: vpc,
	})
	myeks.NewDatabaseCluster(testStack, cluster)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::RDS::DBCluster"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::RDS::DBInstance"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::SecretsManager::Secret"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::SecretsManager::SecretTargetAttachment"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroup"), jsii.Number(2)); // テスト用に作成したsgAlbも含まれるため

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::RDS::DBCluster"), map[string]interface{}{
		"Engine": "aurora-postgresql",
		"CopyTagsToSnapshot": true,
		"DatabaseName": "EksDatabaseName",
		"DBClusterIdentifier": "cluster-identifier",
		"DeletionProtection": false,
		"EngineVersion": "13.6",
		"Port": 5432,
	})
	template.HasResourceProperties(jsii.String("AWS::RDS::DBInstance"), map[string]interface{}{
		"DBInstanceClass": "db.t3.medium",
		"AllowMajorVersionUpgrade": false,
		"AutoMinorVersionUpgrade": true,
		"DBClusterIdentifier": map[string]interface{}{
			"Ref": assertions.Match_StringLikeRegexp(jsii.String("DatabaseCluster")),
		},
		"DBInstanceIdentifier": assertions.Match_StringLikeRegexp(jsii.String("db-instance-identifier")),
		"DeleteAutomatedBackups": false,
		"EnablePerformanceInsights": false,
		"Engine": "aurora-postgresql",
		"EngineVersion": "13.6",
		"PubliclyAccessible": false,
	})
	template.HasResourceProperties(jsii.String("AWS::SecretsManager::Secret"), map[string]interface{}{
		"Name": "database-secrets",
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"SecurityGroupIngress": []map[string]interface{}{
			{
				"Description": assertions.Match_AnyValue(),
				"FromPort": 5432,
				"ToPort": 5432,
				"IpProtocol": "tcp",
				"SourceSecurityGroupId": assertions.Match_AnyValue(),
			},
		},
		"VpcId": map[string]interface{}{
			"Ref": assertions.Match_AnyValue(),
		},
	})
}
