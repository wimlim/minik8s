package apiobj

type FunctionSpec struct {
	Content []byte `yaml:"content" json:"content"`
}

type Function struct {
	ApiVersion string       `yaml:"apiVersion" json:"apiVersion"`
	Kind       string       `yaml:"kind" json:"kind"`
	MetaData   MetaData     `yaml:"metadata" json:"metadata"`
	Spec       FunctionSpec `yaml:"spec" json:"spec"`
}

func (f *Function) GetKind() string {
	return f.Kind
}
func (f *Function) GetName() string {
	return f.MetaData.Name
}
func (f *Function) GetNamespace() string {
	return f.MetaData.Namespace
}
func (f *Function) SetNamespace(namespace string) {
	f.MetaData.Namespace = namespace
}
