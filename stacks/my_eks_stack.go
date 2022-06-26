package stacks

import (
	myeks "mycdk/stacks/my_eks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	constructs "github.com/aws/constructs-go/constructs/v10"
)

func NewMyEKSStack(scope constructs.Construct, id string, props *cdk.StackProps) (cdk.Stack, eks.Cluster, myeks.Repositories) {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	vpc := myeks.NewNetwork(stack)
	sgAlb := myeks.NewLoadBalancer(stack, vpc)
	eksCluster := myeks.NewEksCluster(stack, vpc, sgAlb)
	myeks.NewDatabaseCluster(stack, vpc, eksCluster)
	repos := myeks.NewImageBuilder(stack, &sprops)

	return stack, eksCluster, repos
}