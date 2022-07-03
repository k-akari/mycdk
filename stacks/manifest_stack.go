package stacks

import (
	manifest "mycdk/stacks/manifest"
	myeks "mycdk/stacks/my_eks"

	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	constructs "github.com/aws/constructs-go/constructs/v10"
)

func NewManifestStack(scope constructs.Construct, id string, cluster eks.Cluster, repos myeks.Repositories, props *cdk.StackProps) cdk.Stack {
	stack := cdk.NewStack(scope, &id, props)

	ingressName := manifest.NewMainManifest(stack, cluster, repos)
	manifest.NewCronJobManifest(stack, cluster)
	manifest.NewDNS(stack, ingressName)

	return stack
}