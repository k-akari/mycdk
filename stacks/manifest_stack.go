package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewManifestStack(scope constructs.Construct, id string, cluster eks.Cluster, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// マニフェストの適用
	eks.NewKubernetesManifest(stack, jsii.String("EKSAutoScaler"), &eks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			{
				"apiVersion": jsii.String("apps/v1"),
				"kind": jsii.String("Deployment"),
				"metadata": map[string]*string{
					"name": jsii.String("hello-kubernetes"),
				},
				"spec": map[string]interface{}{
					"replicas": jsii.Number(6),
					"selector": map[string]map[string]*string{
						"matchLabels": {
							"app": jsii.String("hello"),
						},
					},
					"template": map[string]interface{}{
						"metadata": map[string]map[string]*string{
							"labels": {
								"app": jsii.String("hello"),
							},
						},
						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{
								{
									"name": jsii.String("hello-kubernetes"),
									"image": jsii.String("paulbouwer/hello-kubernetes:1.5"),
									"ports": []map[string]*float64{
										{"containerPort": jsii.Number(8080),},
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