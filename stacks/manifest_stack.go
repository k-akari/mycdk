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

	label := map[string]*string{"app": jsii.String("front"),}

	deployment := map[string]interface{}{
		"apiVersion": jsii.String("apps/v1"),
		"kind": jsii.String("Deployment"),
		"metadata": map[string]*string{
			"name": jsii.String("front-deployment"),
		},
		"spec": map[string]interface{}{
			"replicas": jsii.Number(3),
			"selector": map[string]interface{}{
				"matchLabels": label,
			},
			"template": map[string]interface{}{
				"metadata": map[string]map[string]*string{
					"labels": label,
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
	}

	service := map[string]interface{}{
		"apiVersion": jsii.String("v1"),
		"kind": jsii.String("Service"),
		"metadata": map[string]*string{
			"name": jsii.String("front-nodeport-service"),
		},
		"spec": map[string]interface{}{
			"type": jsii.String("NodePort"),
			"ports": []map[string]interface{}{
				{
					"name": jsii.String("http-port"),
					"protocol": jsii.String("TCP"),
					"port": jsii.Number(8080),
					"targetPort": jsii.Number(80),
				},
			},
			"selector": label,
		},
	}

	ingress := map[string]interface{}{
		"apiVersion": jsii.String("networking.k8s.io/v1"),
		"kind": jsii.String("Ingress"),
		"metadata": map[string]interface{}{
			"name": jsii.String("front-ingress"),
			"annotations": map[string]interface{}{
				"kubernetes.io/ingress.class": jsii.String("alb"),
    			"alb.ingress.kubernetes.io/scheme": jsii.String("internet-facing"), //外部公開
				"alb.ingress.kubernetes.io/target-type": jsii.String("instance"),
			},
		},
		"spec": map[string]interface{}{
			"rules": []map[string]interface{}{
				{
					"http": map[string]interface{}{
						"paths": []map[string]interface{}{
							{
								"path": jsii.String("/"),
								"pathType": jsii.String("Prefix"),
								"backend": map[string]interface{}{
									"service": map[string]interface{}{
										"name": jsii.String("front-nodeport-service"),
										"port": map[string]interface{}{
											"number": jsii.Number(8080),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// マニフェストの適用
	eks.NewKubernetesManifest(stack, jsii.String("Manifest"), &eks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			&deployment,
			&service,
			&ingress,
		},
	})

	return stack
}