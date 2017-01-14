# Infrastructure

The AWS infrastructure for the project is managed using [Terraform](https://www.terraform.io/).

## AWS credentials

Terraform will use the AWS credentials for the local "floodwatch" profile. To set it up, use the [AWS Command Line Interface](https://aws.amazon.com/cli/).

```bash
aws configure --profile floodwatch
```

## Execute a plan

To update the infrastructure, update the relevant Terraform files and generate an execution plan.

```bash
make
```

If there are no errors, you can execute the plan.

```bash
make apply
```

Before generating a new execution plan – or after encountering an error – you should delete any existing plan.

```bash
make clean
```
