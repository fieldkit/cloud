package logging

import (
	"context"
	"reflect"
	"strings"
)

func HandlerContext(ctx context.Context, queue string, handlerType reflect.Type, messageType reflect.Type) context.Context {
	name := handlerType.String()
	if handlerType.Kind() == reflect.Ptr {
		name = handlerType.Elem().String()
	}
	return context.WithValue(context.WithValue(ctx, queueKey, queue), handlerKey, name)
}

func CreateFacilityForType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	p := t.PkgPath()
	return strings.Replace(p, packagePrefix, "", 1)
}
