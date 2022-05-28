package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewOidcProviderStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create an Open ID Connect Provider
	provider := awsiam.NewOpenIdConnectProvider(stack, jsii.String("Provider"), &awsiam.OpenIdConnectProviderProps{
		Url: jsii.String("https://token.actions.githubusercontent.com"),
		ClientIds: &[]*string {jsii.String("sts.amazonaws.com")},
		Thumbprints: &[]*string {
			jsii.String("a031c46782e6e6c662c2c87c76da9aa62ccabd8e"),
			jsii.String("6938fd4d98bab03faadb97b34396831e3780aea1"),
		},
	})

	// Create a Federated Principal
	federatedPrincipal := awsiam.NewFederatedPrincipal(provider.OpenIdConnectProviderArn(), &map[string]interface{}{
		"StringLike": map[string]string{"token.actions.githubusercontent.com:sub": "repo:k-akari/*"},
		"StringEquals": map[string]string{"token.actions.githubusercontent.com:aud": "sts.amazonaws.com"},
	}, jsii.String("sts:AssumeRoleWithWebIdentity"))

	// Create a Policy for Federated Principal
    githubPolicy := awsiam.NewManagedPolicy(stack, jsii.String("policy-github"), &awsiam.ManagedPolicyProps{
      	ManagedPolicyName: jsii.String("policy-github"),
      	Document: awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
    		Statements: &[]awsiam.PolicyStatement{
          		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
    			    Effect: awsiam.Effect_ALLOW,
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
    awsiam.NewRole(stack, jsii.String("RoleGithub"), &awsiam.RoleProps{
      	AssumedBy: federatedPrincipal,
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-github"),
      	Description: jsii.String("Role assumed by githubPrincipal for deploying from CI using aws cdk"),
		ManagedPolicies: &[]awsiam.IManagedPolicy{githubPolicy,},
    })

	return stack
}