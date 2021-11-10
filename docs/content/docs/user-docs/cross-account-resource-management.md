---
title: "Manage Resources In Multiple AWS Accounts (CARM)"
description: "Managing resources in different AWS accounts"
lead: ""
draft: false
menu:
  docs:
    parent: "getting-started"
weight: 50
toc: true
---

ACK service controllers can manage resources in different AWS accounts. To enable and start using this feature, as an administrator, you will need to:

  1. Configure the AWS accounts where the resources will be managed
  2. Map AWS accounts with the Role ARNs that need to be assumed
  3. Annotate namespaces with AWS Account IDs

For detailed information about how ACK service controllers manage resources in multiple AWS accounts, please refer to the Cross-Account Resource Management (CARM) [design document](https://github.com/aws-controllers-k8s/community/blob/main/docs/design/proposals/carm/cross-account-resource-management.md).

## Step 1: Configure your AWS accounts

AWS account administrators should create and configure IAM roles to allow ACK service controllers to assume roles in different AWS accounts.

To allow account A (000000000000) to create AWS S3 buckets in account B (111111111111), you can use the following commands:
```bash
# Using account B credentials
aws iam create-role --role-name s3FullAccess \
  --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"AWS": "arn:aws:iam::000000000000:role/roleA-production"}, "Action": "sts:AssumeRole"}]}'
aws iam attach-role-policy --role-name s3FullAccess \
  --policy-arn 'arn:aws:iam::aws:policy/service-role/AmazonS3FullAccess'
```

## Step 2: Map AWS accounts to their associated role ARNs

Create a `ConfigMap` to associate each AWS Account ID with the role ARN that needs to be assumed in order to manage resources in that particular account.

```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: ack-role-account-map
  namespace: ack-system
data:
  "111111111111": arn:aws:iam::111111111111:role/s3FullAccess
EOF
```

## Step 3: Bind accounts to namespaces

To bind AWS accounts to a specific namespace you will have to annotate the namespace with an AWS account ID. For example:
```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: production
  annotations:
    services.k8s.aws/owner-account-id: "111111111111"
EOF
```

For existing namespaces, you can run:
```bash
kubectl annotate namespace production services.k8s.aws/owner-account-id=111111111111
```

### Create resources in different AWS accounts

Next, create your custom resources (CRs) in the associated namespace.

For example, to create an S3 bucket in account B, run the following command:
```bash
cat <<EOF | kubectl apply -f -
apiVersion: s3.services.k8s.aws/v1alpha1
kind: Bucket
metadata:
  name: my-bucket
  namespace: production
spec:
  name: my-bucket
EOF
```

## OpenShift multiple AWS account pre-installation
### Summary
When ACK service controllers are installed via OperatorHub, a cluster administrator will need to perform the following pre-installation steps to provide the controller any credentials and authentication context it needs to interact with the AWS API.

Rather than setting up a `ServiceAccount` like in the EKS instructions above, you need to use IAM users and policies. You will then set the required authentication credentials inside a `ConfigMap` and a `Secret`.

The following directions will use the Elasticache controller as an example, but the instructions should apply to any ACK controller. Just make sure to appropriately name any values that include `elasticache` in them.

Some notes for the following instructions:
* The first user account, which is provided to the ACK controller and will assume the roles of any other accounts, will be identified with profile 000000000000
* The second account will be identified with profile 111111111111
* The instructions will only outline setting up a second account, but you should be able to follow the same instructions to set up user accounts and namespaces
* Substitute actual ID values and namespaces to fit your actual environment values

### Step 1: Create the user for the ACK controller and enable programmatic access

Create a user with the `aws` CLI (named `ack-elasticache-service-controller` in our example):
```bash
aws --profile 000000000000 iam create-user --user-name ack-elasticache-service-controller
```

Make a note of the ARN for this user for use in the next step. An easy way to do so is to export it to an environment variable:
```bash
USER_ARN='arn:aws:iam::000000000000:user/ack-elasticache-service-controller'
```

Enable programmatic access for the user you just created:
```bash
aws --profile 000000000000 iam create-access-key --user-name ack-elasticache-service-controller
```

You should see output with important credentials:
```json
{
    "AccessKey": {
        "UserName": "ack-elasticache-service-controller",
        "AccessKeyId": "00000000000000000000",
        "Status": "Active",
        "SecretAccessKey": "abcdefghIJKLMNOPQRSTUVWXYZabcefghijklMNO",
        "CreateDate": "2021-09-30T19:54:38+00:00"
    }
}
```

This is the user that will end up representing our ACK service controller, which means these are the credentials we’ll eventually pass to our controller. Save or note `AccessKeyId` and `SecretAccessKey` in addition to the ARN you already saved for use in later steps.

### Step 2: Configure role assumption for the controller user

Create the following user policy file, which will allow the `ack-elasticache-service-controller` to assume the role of other users:

```bash
cat > user-policy.json <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListRoles",
                "sts:AssumeRole"
            ],
            "Resource": "*"
        }
    ]
}
EOF
```

Then use the file to create a policy in account 000000000000:
```bash
aws --profile 000000000000 iam create-policy \
  --policy-name ack-can-assume \
  --policy-document file://user-policy.json
```

You should see a new policy ARN in the output of the previous command. Copy it and use it to attach the policy to the `ack-elasticache-service-controller` user:
```bash
aws --profile 000000000000 iam attach-user-policy \
  --user-name ack-elasticache-service-controller \
  --policy-arn "arn:aws:iam::000000000000:policy/ack-can-assume"
```

Now you must create a trust policy in any secondary accounts, i.e. account 111111111111, which allows `ack-elasticache-service-controller` to assume their role. Create the following file with the `$USER_ARN` you saved in step 1:
```bash
cat > trust-policy.json <<EOF
{
    "Version": "2012-10-17",
    "Statement": {
        "Effect": "Allow",
        "Principal": {
            "AWS": "$USER_ARN"
        },
        "Action": "sts:AssumeRole"
    }
}
EOF
```

Use this file to create the role in the secondary account:
```bash
aws --profile 111111111111 iam create-role \
  --role-name ack-manage-elasticache \
  --assume-role-policy-document file://trust-policy.json
```

Save the role ARN from this output for use in step 4, i.e.:
```bash
ROLE_ARN="arn:aws:iam::111111111111:role/ack-manage-elasticache"
```


Give the role permissions using a policy (select the appropriate permission ARN for your use case):
```bash
aws --profile 111111111111 iam attach-role-policy \
    --role-name ack-manage-elasticache \
    --policy-arn 'arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess'
```

You can verify that the policy applied properly with:
```bash
aws --profile 111111111111 iam list-attached-role-policies \
    --role-name ack-manage-elasticache
```

### Step 3: Create the default ACK namespace

Create the namespace for any ACK controllers you might install. The controllers as they are packaged in OperatorHub and OLM expect the namespace to be `ack-system`.
```bash
oc new-project ack-system
```

### Step 4: Configure the ACK controller to change accounts based on namespace

To bind AWS accounts to a specific namespace you will have to annotate the namespace with an AWS account ID. To do so in an existing namespace:
```bash
oc annotate namespace production services.k8s.aws/owner-account-id=111111111111
```

You must also create a `ConfigMap` named `ack-role-account-map` in the `ack-system` namespace where the controller is running, using the `$ROLE_ARN` you saved earlier in step 2:
```bash
oc create configmap \
  --namespace ack-system \
  --from-literal=111111111111=$ROLE_ARN \
  ack-role-account-map
```

### Step 5: Create required `ConfigMap` and `Secret` in OpenShift

Enter the `ack-system` namespace. Create a file, `config.txt`, with the following variables, leaving `ACK_WATCH_NAMESPACE` blank so the controller can properly watch all namespaces, and change any other values to suit your needs:

```bash
ACK_ENABLE_DEVELOPMENT_LOGGING=true
ACK_LOG_LEVEL=debug
ACK_WATCH_NAMESPACE=
AWS_REGION=us-west-2
ACK_RESOURCE_TAGS=hellofromocp
```

Now use `config.txt` to create a `ConfigMap` in your OpenShift cluster:
```bash
oc create configmap \
--namespace ack-system \
--from-env-file=config.txt ack-user-config
```

Save another file, `secrets.txt`, with the authentication values for the first user, which you should have saved from earlier when you created the user's access keys:
```bash
AWS_ACCESS_KEY_ID=00000000000000000000
AWS_SECRET_ACCESS_KEY=abcdefghIJKLMNOPQRSTUVWXYZabcefghijklMNO
```

Use `secrets.txt` to create a `Secret` in your OpenShift cluster:
```bash
oc create secret generic \
--namespace ack-system \
--from-env-file=secrets.txt ack-user-secrets
```

{{% hint type="warning" title="Warning" %}}
If you change the name of either the `ConfigMap` or the `Secret` from the values given above, i.e. `ack-user-config` and `ack-user-secrets`, then installations from OperatorHub will not function properly. The Deployment for the controller is preconfigured for these key values.
{{% /hint %}}

### Step 6: Install the controller

Now you can follow the instructions for [installing the controller using OperatorHub](../install/#install-an-ack-service-controller-with-operatorhub-in-red-hat-openshift).

### Step 7: Create resources in different AWS accounts

Now you can apply Custom Resource Definitions for your different accounts. For example, to create a `UserGroup` in Elasticache:
```bash
oc apply -f - <<EOF
apiVersion: elasticache.services.k8s.aws/v1alpha1
kind: UserGroup
metadata:
  name: example
  namespace: production
spec:
  engine: redis
  userGroupID: multidemoid
  userIDs:
  - default
EOF
```

The cluster should now work the same as [earlier on this page](cross-account-resource-management#create-resources-in-different-aws-accounts) for non-OpenShift installations.


## Next Steps
Checkout the [RBAC and IAM permissions overview](../authorization) to understand how ACK manages authorization