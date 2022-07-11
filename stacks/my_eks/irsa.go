package my_eks

import (
	"fmt"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewIamRolesForServiceAccounts(stack constructs.Construct, cluster eks.Cluster) {
	// Secrets ManagerからSecretリソース作成するPodに付与するIAMロール
	stringEquals := cdk.NewCfnJson(stack, jsii.String("ConditionJson"), &cdk.CfnJsonProps{
		Value: map[string]string{
			*cluster.ClusterOpenIdConnectIssuer()+":sub": fmt.Sprintf("system:serviceaccount:%s:account-to-access-secrets", stack.Node().TryGetContext(jsii.String("namespace"))),
		},
	})

	federatedPrincipal := iam.NewFederatedPrincipal(cluster.OpenIdConnectProvider().OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringEquals": stringEquals,
	}, jsii.String("sts:AssumeRoleWithWebIdentity"))

	secretAccessPolicy := iam.NewManagedPolicy(stack, jsii.String("SecretsManagerAccessPolicy"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("secrets-access-policy"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{jsii.String("*")},
    				Actions: &[]*string{
    					jsii.String("secretsmanager:GetResourcePolicy"),
						jsii.String("secretsmanager:GetSecretValue"),
						jsii.String("secretsmanager:DescribeSecret"),
						jsii.String("secretsmanager:ListSecretVersionIds"),
					},
				}),
			},
    	}),
	})
	
	iam.NewRole(stack, jsii.String("CreateSecretFromSecretsManagerRole"), &iam.RoleProps{
    AssumedBy: federatedPrincipal,
    RoleName: jsii.String("create-secret-from-secrets-manager-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{secretAccessPolicy,},
  	})
}