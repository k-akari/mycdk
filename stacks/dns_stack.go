package stacks

import (
	"fmt"
	"os"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	targets "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewDNSStack(scope constructs.Construct, id string, props *cdk.StackProps) cdk.Stack {
	stack := cdk.NewStack(scope, &id, props)

	hostedZone := route53.NewPublicHostedZone(stack, jsii.String("PublicHostedZone"), &route53.PublicHostedZoneProps{
		ZoneName: jsii.String(os.Getenv("DOMAIN")),
		Comment: jsii.String("free sample domain"),
		CaaAmazon: jsii.Bool(true),
	})

	alb := elbv2.ApplicationLoadBalancer_FromLookup(stack, jsii.String("ApplicationLoadBalancer"), &elbv2.ApplicationLoadBalancerLookupOptions{
		LoadBalancerTags: &map[string]*string{
			"ingress.k8s.aws/stack": jsii.String(fmt.Sprintf("default/%s", stack.Node().TryGetContext(jsii.String("ingressName")))),
		},
	})

	route53.NewARecord(stack, jsii.String("ARecord"), &route53.ARecordProps{
		Zone: hostedZone,
		Ttl: cdk.Duration_Seconds(jsii.Number(300)),
		Target: route53.RecordTarget_FromAlias(targets.NewLoadBalancerTarget(alb)),
	})

	return stack
}