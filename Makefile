.PHONY: terraform-fix-literraform-fix-lint:
	for file in $$(find deployments/ -type f -name '*.tf'); do terraform fmt $$file; done
