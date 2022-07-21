package my_eks

import (
	"os"

	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewHostZone(stack constructs.Construct) {
	// パブリックホストゾーンを作成
	route53.NewPublicHostedZone(stack, jsii.String("PublicHostedZone"), &route53.PublicHostedZoneProps{
		ZoneName: jsii.String(os.Getenv("DOMAIN")),
		Comment: jsii.String("free sample domain"),
		CaaAmazon: jsii.Bool(true),
	})
}