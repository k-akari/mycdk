package my_eks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewRepositories(stack constructs.Construct) (repoMigration ecr.Repository) {
	// アプリケーションイメージのECRリポジトリ
	ecr.NewRepository(stack, jsii.String("EKSAppImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{{MaxImageCount: jsii.Number(1),},},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-app"),
	})

	// DBマイグレーション用イメージのECRリポジトリ
	repoMigration = ecr.NewRepository(stack, jsii.String("EKSMigrationImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{{MaxImageCount: jsii.Number(1),},},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-migration"),
	})

	// NginxイメージのECRリポジトリ
	ecr.NewRepository(stack, jsii.String("EKSWebImageRepository"), &ecr.RepositoryProps{
		ImageScanOnPush: jsii.Bool(true),
		LifecycleRules: &[]*ecr.LifecycleRule{{MaxImageCount: jsii.Number(1),},},
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		RepositoryName: jsii.String("eks-web"),
	})

	return
}