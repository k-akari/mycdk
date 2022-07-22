package my_eks

import (
	"os"

	acm "github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewHostZone(stack constructs.Construct) {
	// パブリックホストゾーンを作成
	hostedZone := route53.NewPublicHostedZone(stack, jsii.String("PublicHostedZone"), &route53.PublicHostedZoneProps{
		ZoneName: jsii.String(os.Getenv("DOMAIN")),
		Comment: jsii.String("free sample domain"),
		CaaAmazon: jsii.Bool(true),
	})

	acm.NewCertificate(stack, jsii.String("Certificate"), &acm.CertificateProps{
		DomainName: jsii.String("mycdk-app-example.tk"),
		Validation: acm.CertificateValidation_FromDns(hostedZone),
	})
}