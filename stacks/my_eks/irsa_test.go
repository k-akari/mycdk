package my_eks_test

import (
	myeks "mycdk/stacks/my_eks"
	"os"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewIamRolesForServiceAccounts(t *testing.T) {
	app := cdk.NewApp(nil)

	// テスト対象のスタックテンプレートを用意
	testStack := cdk.NewStack(app, jsii.String("TestStack"), &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},})
	cluster := eks.NewCluster(testStack, jsii.String("EKSCluster"), &eks.ClusterProps{
		Version: eks.KubernetesVersion_Of(jsii.String("1.22")), // kubernetesのバージョン
	})
	myeks.NewIamRolesForServiceAccounts(testStack, cluster)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::IAM::Role"), jsii.Number(7)); // eks.NewClusterによって作成されるIamRoleが6つ
	template.ResourceCountIs(jsii.String("AWS::IAM::ManagedPolicy"), jsii.Number(1));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::IAM::ManagedPolicy"), map[string]interface{}{
		"PolicyDocument": map[string]interface{}{
			"Statement": assertions.Match_AnyValue(),
			"Version": assertions.Match_AnyValue(),
		},
		"ManagedPolicyName": "secrets-access-policy",
	})
	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": assertions.Match_AnyValue(),
			"Version": assertions.Match_AnyValue(),
		},
		"ManagedPolicyArns": assertions.Match_AnyValue(),
		"RoleName": "create-secret-from-secrets-manager-role",
	})
}
