package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

const (
	messagesTemplate = `// @flow
{{range .}}
type {{.Name}}Type = {
{{- range .Properties}}
{{- if .Message }}
  {{.Name}}?: {{.Type}}Type,
{{- else }}
  {{.Name}}?: {{.Type}},
{{- end}}
{{- end}}
}

export class {{.Name}} {
{{- range .Properties}}
  {{.Name}}: {{ if .Message }}?{{end}}{{.Type}};
{{- end}}
  constructor(obj: {{.Name}}Type){
{{- range .Properties}}
    this.{{.Name}} = obj.{{.Name}} ?
      {{if .Message }}new {{.Type}}(obj.{{.Name}}) : null{{else}}this.{{.Name}} = obj.{{.Name}} : {{.Default}}{{end}}
{{- end}}
  }
}
{{end}}`
)

type FlowProperty struct {
	Name, Type, Default string
	Message             bool
}

type FlowClass struct {
	Name       string
	Properties []FlowProperty
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	request := &plugin_go.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, request); err != nil {
		panic(err)
	}

	messageTypes := map[string]*descriptor.DescriptorProto{}
	for _, file := range request.ProtoFile {
		for _, messageType := range file.MessageType {
			messageTypes[messageType.GetName()] = messageType
		}
	}

	flowClasses := make([]FlowClass, 0, len(messageTypes))
	for messageTypeName, messageType := range messageTypes {

		flowClass := FlowClass{
			Name:       messageTypeName,
			Properties: make([]FlowProperty, len(messageType.GetField())),
		}

		for i, messageTypeField := range messageType.GetField() {
			flowClass.Properties[i].Name = messageTypeField.GetName()
			switch messageTypeField.GetType().String() {
			case "TYPE_DOUBLE", "TYPE_FLOAT", "TYPE_INT64", "TYPE_UINT64", "TYPE_INT32", "TYPE_FIXED64", "TYPE_FIXED32", "TYPE_UINT32", "TYPE_SFIXED32", "TYPE_SFIXED64", "TYPE_SINT32", "TYPE_SINT64":
				flowClass.Properties[i].Type = "number"
				flowClass.Properties[i].Default = "0"
			case "TYPE_BYTES", "TYPE_ENUM", "TYPE_STRING":
				flowClass.Properties[i].Type = "string"
				flowClass.Properties[i].Default = "''"
			case "TYPE_MESSAGE":
				flowClass.Properties[i].Type = strings.TrimPrefix(messageTypeField.GetTypeName(), ".")
				flowClass.Properties[i].Message = true
			}
		}

		flowClasses = append(flowClasses, flowClass)
	}

	javascript := new(bytes.Buffer)
	template.Must(template.New("messages").Parse(messagesTemplate)).Execute(javascript, flowClasses)

	response := &plugin_go.CodeGeneratorResponse{
		File: []*plugin_go.CodeGeneratorResponse_File{
			&plugin_go.CodeGeneratorResponse_File{
				Name:    proto.String("flow.js"),
				Content: proto.String(javascript.String()),
			},
		},
	}

	data, err = proto.Marshal(response)
	if err != nil {
		panic(err)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}
