package apiobj

type Node struct {
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
	IP         string   `yaml:"ip" json:"ip"`
}

func (n *Node) GetKind() string {
	return n.Kind
}
func (n *Node) GetName() string {
	return n.MetaData.Name
}
