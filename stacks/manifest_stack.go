package stacks

import (
	myeks "mycdk/stacks/my_eks"

	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	jsii "github.com/aws/jsii-runtime-go"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	constructs "github.com/aws/constructs-go/constructs/v10"
)

func NewManifestStack(scope constructs.Construct, id string, cluster eks.Cluster, repos myeks.Repositories, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// マニフェストの適用
	eks.NewKubernetesManifest(stack, jsii.String("EKSManifest"), &eks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			{
				"apiVersion": jsii.String("apps/v1"),
				"kind": jsii.String("Deployment"),
				"metadata": map[string]*string{
					"name": jsii.String("front-deployment"),
				},
				"spec": map[string]interface{}{
					"replicas": jsii.Number(3),
					"selector": map[string]map[string]*string{
						"matchLabels": {
							"app": jsii.String("front"),
						},
					},
					"template": map[string]interface{}{
						"metadata": map[string]map[string]*string{
							"labels": {
								"app": jsii.String("front"),
							},
						},
						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{
								{
									"name": jsii.String("app"),
									"image": repos.App.RepositoryUri(),
									"ports": []map[string]*float64{
										{"containerPort": jsii.Number(8080),},
									},
								},
								{
									"name": jsii.String("web"),
									"image": repos.Web.RepositoryUri(),
									"ports": []map[string]*float64{
										{"containerPort": jsii.Number(80),},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	return stack
}