package manifest

import (
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	jsii "github.com/aws/jsii-runtime-go"

	constructs "github.com/aws/constructs-go/constructs/v10"
)

func NewCronJobManifest(stack constructs.Construct, cluster eks.Cluster) {
	label := map[string]*string{"app": jsii.String("cronjob"),}

	sampleJob1 := map[string]interface{}{
		"apiVersion": jsii.String("batch/v1"),
		"kind": jsii.String("CronJob"),
		"metadata": map[string]*string{
			"name": jsii.String("sample-job-1"),
		},
		"spec": map[string]interface{}{
			"schedule": jsii.String("*/15 * * * *"),
			"concurrencyPolicy": jsii.String("Allow"),
			"successfulJobsHistoryLimit": jsii.Number(4),
			"failedJobsHistoryLimit": jsii.Number(2),
			"startingDeadlineSeconds": jsii.Number(9000),
			"jobTemplate": map[string]interface{}{
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"metadata": map[string]map[string]*string{
							"labels": label,
						},
						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{
								{
									"name": jsii.String("job-container-1"),
									"image": jsii.String("busybox:1.28"),
									"imagePullPolicy": jsii.String("IfNotPresent"),
									"command": []*string{
										jsii.String("/bin/sh"),
										jsii.String("-c"),
										jsii.String("date; echo Hello from the Kubernetes cluster"),
									},
								},
							},
							"restartPolicy": jsii.String("OnFailure"),
						},
					},
				},
			},
		},
	}

	// マニフェストの適用
	eks.NewKubernetesManifest(stack, jsii.String("CronJobManifest"), &eks.KubernetesManifestProps{
		Cluster: cluster,
		Manifest: &[]*map[string]interface{}{
			&sampleJob1,
		},
	})
}