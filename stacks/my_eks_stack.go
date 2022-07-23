package stacks

import (
	myeks "mycdk/stacks/my_eks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	constructs "github.com/aws/constructs-go/constructs/v10"
)

func NewMyEKSStack(scope constructs.Construct, id string, props *cdk.StackProps) (stack cdk.Stack) {
	stack = cdk.NewStack(scope, &id, props)

	vpc, vpcEndpoint := myeks.NewNetwork(stack)
	eksCluster := myeks.NewEksCluster(stack, vpc, vpcEndpoint)
	myeks.NewIamRolesForServiceAccounts(stack, eksCluster)
	dbCluster := myeks.NewDatabaseCluster(stack, eksCluster)
	repoMigrate := myeks.NewImageBuilder(stack, props)
	myeks.NewDBMigrator(stack, repoMigrate, dbCluster, vpcEndpoint, props)
	myeks.NewHostZone(stack)

	return
}