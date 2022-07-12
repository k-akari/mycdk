package stacks

import (
	myeks "mycdk/stacks/my_eks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	constructs "github.com/aws/constructs-go/constructs/v10"
)

func NewMyEKSStack(scope constructs.Construct, id string, props *cdk.StackProps) (cdk.Stack, eks.Cluster) {
	stack := cdk.NewStack(scope, &id, props)

	vpc := myeks.NewNetwork(stack)
	eksCluster := myeks.NewEksCluster(stack, vpc)
	myeks.NewIamRolesForServiceAccounts(stack, eksCluster)
	myeks.NewDatabaseCluster(stack, vpc, eksCluster)
	myeks.NewImageBuilder(stack, props)

	return stack, eksCluster
}