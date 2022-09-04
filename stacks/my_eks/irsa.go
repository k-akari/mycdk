package my_eks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewIamRolesForServiceAccounts(stack constructs.Construct, cluster eks.Cluster, repoMigration ecr.Repository) {
	// Secrets ManagerからSecretリソース作成するPodに付与するIAMロール
	principalESO := iam.NewFederatedPrincipal(cluster.OpenIdConnectProvider().OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringLike": cdk.NewCfnJson(stack, jsii.String("ConditionForAccountToAccessSecrets"), &cdk.CfnJsonProps{
			Value: map[string]string{
				*cluster.ClusterOpenIdConnectIssuer()+":sub": "system:serviceaccount:main:*",
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
		"StringLike": cdk.NewCfnJson(stack, jsii.String("ConditionForAccountForExternalDNS"), &cdk.CfnJsonProps{
			Value: map[string]string{
				*cluster.ClusterOpenIdConnectIssuer()+":sub": "system:serviceaccount:main:*",
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

	// ArgoCD Image UpdaterがECRリポジトリのイメージタグ情報を取得するためのIAMロール
	principalArgoCDImageUpdater := iam.NewFederatedPrincipal(cluster.OpenIdConnectProvider().OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringLike": cdk.NewCfnJson(stack, jsii.String("ConditionForAccountForArgoCDImageUpdater"), &cdk.CfnJsonProps{
			Value: map[string]string{
				*cluster.ClusterOpenIdConnectIssuer()+":sub": "system:serviceaccount:argocd:*", // "system:serviceaccount:<argocd-image-updaterを配置するnamespace>:<argocd-image-updaterに設定するServiceAccount名>"
			},
		}),
	}, jsii.String("sts:AssumeRoleWithWebIdentity"))

	updateImagePolicy := iam.NewManagedPolicy(stack, jsii.String("UpdateImagePolicy"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("argocd-update-image-policy"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{jsii.String("*")},
    				Actions: &[]*string{
    					jsii.String("ecr:GetAuthorizationToken"),
    					jsii.String("ecr:ListImages"),
    					jsii.String("ecr:BatchGetImage"),
    					jsii.String("ecr:GetDownloadUrlForLayer"),
					},
				}),
			},
    	}),
	})

	iam.NewRole(stack, jsii.String("ArgoCDImageUpdaterRole"), &iam.RoleProps{
    	AssumedBy: principalArgoCDImageUpdater,
    	RoleName: jsii.String("role-for-image-updater"),
		ManagedPolicies: &[]iam.IManagedPolicy{updateImagePolicy,},
  	})

	// DB Migrateを実行するPodがDockerイメージをプルするためのIamRole
	principalDBMigrator := iam.NewFederatedPrincipal(cluster.OpenIdConnectProvider().OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringLike": cdk.NewCfnJson(stack, jsii.String("ConditionForAccountForDBMigrator"), &cdk.CfnJsonProps{
			Value: map[string]string{
				*cluster.ClusterOpenIdConnectIssuer()+":sub": "system:serviceaccount:main:*",
			},
		}),
	}, jsii.String("sts:AssumeRoleWithWebIdentity"))

	pullImagePolicy := iam.NewManagedPolicy(stack, jsii.String("PullImagePolicyForDBMigrator"), &iam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("pull-image-policy-for-db-migrator"),
		Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    				Effect: iam.Effect_ALLOW,
    				Resources: &[]*string{repoMigration.RepositoryArn()},
    				Actions: &[]*string{
    					jsii.String("ecr:BatchGetImage"),
    					jsii.String("ecr:GetDownloadUrlForLayer"),
					},
				}),
			},
    	}),
	})

	iam.NewRole(stack, jsii.String("DBMigratorRole"), &iam.RoleProps{
      	AssumedBy: principalDBMigrator,
      	RoleName: jsii.String("role-for-db-migrator"),
		ManagedPolicies: &[]iam.IManagedPolicy{pullImagePolicy,},
    })
}