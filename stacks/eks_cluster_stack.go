package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewEksClusterStack(scope constructs.Construct, id string, vpc ec2.Vpc, props *cdk.StackProps) (cdk.Stack, eks.Cluster) {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// EKSコントロールプレーンに付与するIAMロール
	masterRole := iam.NewRole(stack, jsii.String("EKSMasterRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), &iam.ServicePrincipalOpts{}),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("eks-master-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
    })

	// EKSクラスター
	cluster := eks.NewCluster(stack, jsii.String("EKSCluster"), &eks.ClusterProps{
		ClusterName: jsii.String("eks-cluster"),
		DefaultCapacity: jsii.Number(0), // デフォルトインスタンスは作らない
		EndpointAccess: eks.EndpointAccess_PUBLIC(),
		MastersRole: masterRole, // クラスターのマスターロール
		Version: eks.KubernetesVersion_V1_21(), // kubernetesのバージョン
		Vpc: vpc, // EKSクラスターをデプロイするVPC
	})

	// Nodeに付与するIAMロール
	nodeRole := iam.NewRole(stack, jsii.String("EKSNodeRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &iam.ServicePrincipalOpts{}),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("eks-node-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
		},
    })

	// EKSクラスターにNodeグループを追加
	cluster.AddNodegroupCapacity(jsii.String("EKSNodeGroupCapacity"), &eks.NodegroupOptions{
		AmiType: eks.NodegroupAmiType_AL2_X86_64,
		CapacityType: eks.CapacityType_SPOT,
		DesiredSize: jsii.Number(3),
		DiskSize: jsii.Number(10),
		InstanceTypes: &[]ec2.InstanceType{
			ec2.NewInstanceType(jsii.String("t2.micro")),
			ec2.NewInstanceType(jsii.String("t2.small")),
			ec2.NewInstanceType(jsii.String("t2.medium")),
			ec2.NewInstanceType(jsii.String("t3.micro")),
			ec2.NewInstanceType(jsii.String("t3.small")),
			ec2.NewInstanceType(jsii.String("t3.medium")),
		},
		Labels: &map[string]*string {
			"app": jsii.String("practice"),
		},
		MaxSize: jsii.Number(6),
		MinSize: jsii.Number(3),
		NodegroupName: jsii.String("eks-node-group"),
		NodeRole: nodeRole,
		Subnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
		},
		Tags: &map[string]*string {
			"Service": jsii.String("service_name"),
			"Environment": jsii.String("production"),
		},
	})

	// IAMユーザーがクラスターと対話するにはsystem:masters アクセス許可を付与する必要がある。
	// https://docs.aws.amazon.com/ja_jp/eks/latest/userguide/add-user-role.html#aws-auth-users
	user := iam.User_FromUserName(stack, jsii.String("ImportedUserByName"), jsii.String("akari"))
	cluster.AwsAuth().AddUserMapping(user, &eks.AwsAuthMapping{
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	// AWS CLI利用時にMFA認証必須化するために使用しているロールとsystem:mastersを結びつける
	mfaRole := iam.Role_FromRoleName(stack, jsii.String("ImportedRoleByName"), jsii.String("AdminRole"))
	cluster.AwsAuth().AddRoleMapping(mfaRole, &eks.AwsAuthMapping{
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	return stack, cluster
}