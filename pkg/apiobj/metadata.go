package apiobj

type MetaData struct {
	UID      string            `yaml:"uid" json:"uid"`
	Name      string            `yaml:"name" json:"name"`
	Namespace string            `yaml:"namespace" json:"namespace"`
	Labels    map[string]string `yaml:"labels" json:"labels"`
}