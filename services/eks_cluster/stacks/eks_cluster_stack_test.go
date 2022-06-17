package stacks_test

import (
	"testing"

	"github.com/k-akari/services/eks_cluster/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestImageBuilderStack(t *testing.T) {
	app := cdk.NewApp(nil)
	stack := stacks.NewImageBuilderStack(app, "MyStack", nil)

	template := assertions.Template_FromStack(stack)

	template.ResourceCountIs(jsii.String("AWS::IAM::Role"), jsii.Number(1));
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "ent-example",
	})
}
