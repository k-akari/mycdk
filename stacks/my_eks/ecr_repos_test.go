package my_eks_test

import (
	myeks "mycdk/stacks/my_eks"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestImageBuilderStack(t *testing.T) {
	app := cdk.NewApp(nil)
	
	// テスト対象のスタックテンプレートを用意
	testStack := cdk.NewStack(app, jsii.String("TestStack"), nil)
	myeks.NewRepositories(testStack)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::ECR::Repository"), jsii.Number(3));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "eks-app",
	})
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "eks-migration",
	})
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "eks-web",
	})
}
