package stacks_test

import (
	"mycdk/stacks"
	"os"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewDNSStack(t *testing.T) {
	app := cdk.NewApp(nil)

	// テスト対象のスタックテンプレートを用意
	testStack := stacks.NewDNSStack(app, "TestStack", &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},})
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::Route53::HostedZone"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::Route53::RecordSet"), jsii.Number(2));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::Route53::HostedZone"), map[string]interface{}{
		"HostedZoneConfig": map[string]interface{}{
			"Comment": assertions.Match_AnyValue(),
		},
		"Name": os.Getenv("DOMAIN") + ".",
	})
	template.HasResourceProperties(jsii.String("AWS::Route53::RecordSet"), map[string]interface{}{
		"Name": os.Getenv("DOMAIN") + ".",
		"Type": "A",
		"AliasTarget": map[string]interface{}{
			"DNSName": assertions.Match_AnyValue(),
			"HostedZoneId": assertions.Match_AnyValue(),
		},
		"HostedZoneId": map[string]interface{}{
			"Ref": assertions.Match_AnyValue(),
		},
	})
}
