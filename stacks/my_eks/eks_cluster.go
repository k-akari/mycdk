package my_eks

import (
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewEksCluster(stack constructs.Construct, vpc ec2.Vpc, vpcEndpoint ec2.InterfaceVpcEndpoint) (cluster eks.Cluster) {
	// EKSコントロールプレーンに付与するIAMロールの作成
	masterRole := iam.NewRole(stack, jsii.String("EKSMasterRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), &iam.ServicePrincipalOpts{}),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("eks-master-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
		},
    })

	// EKSクラスターの作成
	cluster = eks.NewCluster(stack, jsii.String("EKSCluster"), &eks.ClusterProps{
		AlbController: &eks.AlbControllerOptions{
			Version: eks.AlbControllerVersion_V2_4_1(),
		},
		ClusterName: jsii.String("eks-cluster"),
		DefaultCapacity: jsii.Number(0), // デフォルトインスタンスは作らない
		EndpointAccess: eks.EndpointAccess_PUBLIC(),
		MastersRole: masterRole,
		Version: eks.KubernetesVersion_Of(jsii.String("1.22")),
		Vpc: vpc,
	})

	// ノードグループがプライベートリンクを利用してECRからイメージを取得する
	vpcEndpoint.Connections().AllowFrom(cluster, ec2.Port_AllTraffic(), jsii.String("Allow access to VPC endpoint from EKS cluster"))

	// Nodeに付与するIAMロールの作成
	nodeRole := iam.NewRole(stack, jsii.String("EKSNodeRole"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &iam.ServicePrincipalOpts{}),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("eks-node-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore")),
		},
    })

	// ノードグループの起動テンプレートを作成
	launchTemplate := ec2.NewLaunchTemplate(stack, jsii.String("EKSNodesLaunchTemplate"), &ec2.LaunchTemplateProps{
		DetailedMonitoring: jsii.Bool(false),
		DisableApiTermination: jsii.Bool(false),
		EbsOptimized: jsii.Bool(false),
		HibernationConfigured: jsii.Bool(false),
		LaunchTemplateName: jsii.String("eks-nodes-launch-template"),
		NitroEnclaveEnabled: jsii.Bool(false),
	})

	// EKSクラスターにNodeグループを追加
	cluster.AddNodegroupCapacity(jsii.String("EKSNodeGroup"), &eks.NodegroupOptions{
		AmiType: eks.NodegroupAmiType_AL2_X86_64,
		CapacityType: eks.CapacityType_SPOT,
		DesiredSize: jsii.Number(2),
		InstanceTypes: &[]ec2.InstanceType{
			ec2.NewInstanceType(jsii.String("m5a.large")),
			ec2.NewInstanceType(jsii.String("m5.large")),
			ec2.NewInstanceType(jsii.String("m5ad.large")),
			ec2.NewInstanceType(jsii.String("m5d.large")),
			ec2.NewInstanceType(jsii.String("m5n.large")),
			ec2.NewInstanceType(jsii.String("m5dn.large")),
		},
		Labels: &map[string]*string {
			"app": jsii.String("practice"),
		},
		LaunchTemplateSpec: &eks.LaunchTemplateSpec{
			Id: launchTemplate.LaunchTemplateId(),
			Version: launchTemplate.LatestVersionNumber(),
		},
		MaxSize: jsii.Number(6),
		MinSize: jsii.Number(2),
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

	return
}