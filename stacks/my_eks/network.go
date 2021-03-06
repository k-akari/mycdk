package my_eks

import (
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewNetwork(stack constructs.Construct) (vpc ec2.Vpc, vpcEndpoint ec2.InterfaceVpcEndpoint ) {
	// 3AZにまたがるVPCの作成
	// AZ毎にパブリックサブネットとNATゲートウェイへルートを向けたプライベートサブネットと完全に独立したプライベートサブネットを1つずつ作成
	vpc = ec2.NewVpc(stack, jsii.String("VPC"), &ec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
		MaxAzs: jsii.Number(2),
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

	// VPCエンドポイントの作成
	vpcEndpoint = vpc.AddInterfaceEndpoint(jsii.String("VPCEndpoint"), &ec2.InterfaceVpcEndpointOptions{
		Service: ec2.InterfaceVpcEndpointAwsService_ECR(),
		LookupSupportedAzs: jsii.Bool(true),
	})

	return
}