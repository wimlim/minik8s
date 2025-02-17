package apiobj

import (
	"reflect"
)

type ApiObject interface {
	GetKind() string
	GetName() string
	GetNamespace() string
	SetNamespace(string)
}

var KindStr2Type = map[string]reflect.Type{
	"Pod":     reflect.TypeOf(&Pod{}).Elem(),
	"Service": reflect.TypeOf(&Service{}).Elem(),
	// "Dns":        reflect.TypeOf(&Dns{}).Elem(),
	// "Node":       reflect.TypeOf(&Node{}).Elem(),
	// "Job":        reflect.TypeOf(&Job{}).Elem(),
	// "Replicaset": reflect.TypeOf(&ReplicaSet{}).Elem(),
	// "Hpa":        reflect.TypeOf(&Hpa{}).Elem(),
}
