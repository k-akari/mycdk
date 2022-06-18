package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewGitHubActionsStack(scope constructs.Construct, id string, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// Create an Open ID Connect Provider
	provider := iam.NewOpenIdConnectProvider(stack, jsii.String("Provider"), &iam.OpenIdConnectProviderProps{
		Url: jsii.String("https://token.actions.githubusercontent.com"),
		ClientIds: &[]*string {jsii.String("sts.amazonaws.com")},
		Thumbprints: &[]*string {
			jsii.String("a031c46782e6e6c662c2c87c76da9aa62ccabd8e"),
			jsii.String("6938fd4d98bab03faadb97b34396831e3780aea1"),
		},
	})

	// Create a Federated Principal
	federatedPrincipal := iam.NewFederatedPrincipal(provider.OpenIdConnectProviderArn(), &map[string]interface{}{
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
    			    	jsii.String("ec2:AuthorizeSecurityGroupEgress"),
    			    	jsii.String("ec2:AuthorizeSecurityGroupIngress"),
    			    	jsii.String("ec2:RevokeSecurityGroupEgress"),
    			    	jsii.String("ec2:RevokeSecurityGroupIngress"),
    			    	jsii.String("ssm:SendCommand"),
						jsii.String("sts:GetCallerIdentity"),
					},
				}),
			},
      	}),
    })

	// Create an Iam Role for Federated Principal
    iam.NewRole(stack, jsii.String("RoleGithub"), &iam.RoleProps{
      	AssumedBy: federatedPrincipal,
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-github"),
      	Description: jsii.String("Role assumed by githubPrincipal for deploying from CI using aws cdk"),
		ManagedPolicies: &[]iam.IManagedPolicy{githubPolicy,},
    })

	return stack
}