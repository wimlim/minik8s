package apiobj

type ReplicaSetSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels" json:"matchLabels"`
}
type PodTemplateSpec struct {
	MetaData MetaData `yaml:"metadata" json:"metadata"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
}

type ReplicaSetSpec struct {
	Replicas int                `yaml:"replicas" json:"replicas"`
	Selector ReplicaSetSelector `yaml:"selector" json:"selector"`
	Template PodTemplateSpec    `yaml:"template" json:"template"`
}

type ReplicaSetStatus struct {
	SpecReplicas  int `yaml:"specReplicas" json:"specReplicas"`
	ReadyReplicas int `yaml:"readyReplicas" json:"readyReplicas"`
}

type ReplicaSet struct {
	ApiVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind" json:"kind"`
	MetaData   MetaData         `yaml:"metadata" json:"metadata"`
	Spec       ReplicaSetSpec   `yaml:"spec" json:"spec"`
	Status     ReplicaSetStatus `yaml:"status" json:"status"`
}

func (r *ReplicaSet) GetKind() string {
	return r.Kind
}
func (r *ReplicaSet) GetName() string {
	return r.MetaData.Name
}
func (r *ReplicaSet) GetNamespace() string {
	return r.MetaData.Namespace
}
func (r *ReplicaSet) SetNamespace(namespace string) {
	r.MetaData.Namespace = namespace
}
