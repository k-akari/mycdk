package my_eks

import (
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewGitHubActions(stack constructs.Construct) {
	// Create an Open ID Connect Provider
	provider := iam.NewOpenIdConnectProvider(stack, jsii.String("Provider"), &iam.OpenIdConnectProviderProps{
		Url: jsii.String("https://token.actions.githubusercontent.com"),
		ClientIds: &[]*string {jsii.String("sts.amazonaws.com")},
		Thumbprints: &[]*string {
			jsii.String("6938fd4d98bab03faadb97b34396831e3780aea1"),
		},
	})

	// Create a Federated Principal
	principalGitHub := iam.NewFederatedPrincipal(provider.OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringLike": map[string]string{"token.actions.githubusercontent.com:sub": "repo:k-akari/*"},
		"StringEquals": map[string]string{"token.actions.githubusercontent.com:aud": "sts.amazonaws.com"},
	}, jsii.String("sts:AssumeRoleWithWebIdentity"))

	// Create a Policy for Federated Principal
    githubPolicy := iam.NewManagedPolicy(stack, jsii.String("policy-github"), &iam.ManagedPolicyProps{
      	ManagedPolicyName: jsii.String("policy-github"),
      	Document: iam.NewPolicyDocument(&iam.PolicyDocumentProps{
    		Statements: &[]iam.PolicyStatement{
          		iam.NewPolicyStatement(&iam.PolicyStatementProps{
    			    Effect: iam.Effect_ALLOW,
    			    Resources: &[]*string{jsii.String("*")},
    			    Actions: &[]*string{
						jsii.String("sts:GetCallerIdentity"),
						jsii.String("ecr:GetAuthorizationToken"),
						jsii.String("ecr:CompleteLayerUpload"),
    					jsii.String("ecr:UploadLayerPart"),
    					jsii.String("ecr:InitiateLayerUpload"),
    					jsii.String("ecr:BatchCheckLayerAvailability"),
    					jsii.String("ecr:PutImage"),
					},
				}),
			},
      	}),
    })

	// Create an Iam Role for Federated Principal
    iam.NewRole(stack, jsii.String("RoleGithub"), &iam.RoleProps{
      	AssumedBy: principalGitHub,
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-github"),
      	Description: jsii.String("Role assumed by githubPrincipal for deploying from CI using aws cdk"),
		ManagedPolicies: &[]iam.IManagedPolicy{githubPolicy,},
    })
}