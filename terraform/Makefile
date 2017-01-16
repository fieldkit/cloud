all: terraform.tfplan

apply: terraform.tfplan
	terraform apply

terraform.tfplan:
	terraform plan -out=terraform.tfplan

clean:
	rm -rf terraform.tfplan

.PHONY: all apply clean
