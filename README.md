# IaC Repository using CDK for Golang

## 1. How to build a development environment
1. Build the image and daemonize the container
```bash
docker compose build
docker compose up -d
```

2. Click the green icon in the lower left corner of VS Code and click Reopen in Container

## 2. How to deploy an infrastructure of EKS on EC2
1. Deploy MyEKSStack
```bash
cdk deploy MyEKSStack --profile akari_mfa
```

2. Open the CodeBuild page of management console and build the nginx and app images

3. Deploy ManifestStack
```bash
cdk deploy MyEKSStack --profile akari_mfa
```

4. Update kubeconfig
```bash
aws eks update-kubeconfig --region region-code --name cluster-name --profile akari_mfa
```
