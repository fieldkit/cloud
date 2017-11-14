default:
	go build server.go

image:
	docker build -t conservify/fk-cloud-server .

generate:
	# https://github.com/goadesign/goa/issues/923
	mv vendor vendor-temp
	(cd $(GOPATH)/src/github.com/fieldkit/cloud/server && go generate) || true
	mv vendor-temp vendor


