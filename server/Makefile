default: none

none:
	@echo "Try 'make generate'"

generate:
  # https://github.com/goadesign/goa/issues/923
	mv vendor vendor-temp
	go generate || true
	mv vendor-temp vendor
