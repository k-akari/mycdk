package my_eks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewIamRolesForServiceAccounts(stack constructs.Construct, cluster eks.Cluster) {
	// Secrets ManagerからSecretリソース作成するPodに付与するIAMロール
	principalESO := iam.NewFederatedPrincipal(cluster.OpenIdConnectProvider().OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringEquals": cdk.NewCfnJson(stack, jsii.String("ConditionForAccountToAccessSecrets"), &cdk.CfnJsonProps{
			Value: map[string]string{
				*cluster.ClusterOpenIdConnectIssuer()+":sub": "system:serviceaccount:main:account-to-access-secrets",
			},
		}),
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
    AssumedBy: principalESO,
    RoleName: jsii.String("create-secret-from-secrets-manager-role"),
		ManagedPolicies: &[]iam.IManagedPolicy{secretAccessPolicy,},
  	})

	// Route53のレコードを作成・更新・削除するExternal DNSのPodに付与するIAMロール
	principalExternalDNS := iam.NewFederatedPrincipal(cluster.OpenIdConnectProvider().OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringEquals": cdk.NewCfnJson(stack, jsii.String("ConditionForAccountForExternalDNS"), &cdk.CfnJsonProps{
			Value: map[string]string{
				*cluster.ClusterOpenIdConnectIssuer()+":sub": "system:serviceaccount:main:account-for-external-dns",
			},
		}),
	}, jsii.String("sts:AssumeRoleWithWebIdentity"))

	changeRecordSetsPolicy := iam.NewManagedPolicy(stack, jsii.String("ChangeRecordSetsPolicy"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("change-record-sets-policy"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{jsii.String("arn:aws:route53:::hostedzone/*")},
    				Actions: &[]*string{
    					jsii.String("route53:ChangeResourceRecordSets"),
					},
				}),
			},
    	}),
	})
	listRecordSetsPolicy := iam.NewManagedPolicy(stack, jsii.String("ListRecordSetsPolicy"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("list-record-sets-policy"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{jsii.String("*")},
    				Actions: &[]*string{
    					jsii.String("route53:ListHostedZones"),
    					jsii.String("route53:ListResourceRecordSets"),
					},
				}),
			},
    	}),
	})

	iam.NewRole(stack, jsii.String("EditRoute53RecordRole"), &iam.RoleProps{
    AssumedBy: principalExternalDNS,
    RoleName: jsii.String("role-for-external-dns"),
		ManagedPolicies: &[]iam.IManagedPolicy{changeRecordSetsPolicy, listRecordSetsPolicy},
  	})
}