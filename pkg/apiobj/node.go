package apiobj

type node struct {
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
	IP         string   `yaml:"ip" json:"ip"`
}

func (n *node) GetKind() string {
	return n.Kind
}
func (n *node) GetName() string {
	return n.MetaData.Name
}
