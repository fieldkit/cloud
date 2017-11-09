default: none

none:
	@echo "Try:"
	@echo "make image"
	@echo "make generate"

image:
	docker build -t conservify/fk-cloud-server .

generate:
	# https://github.com/goadesign/goa/issues/923
	mv vendor vendor-temp
	go generate || true
	mv vendor-temp vendor


