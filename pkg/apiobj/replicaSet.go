package apiobj

type ReplicaSetSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}
type PodTemplateSpec struct {
	MetaData MetaData `json:"metadata"`
	Spec     PodSpec  `json:"spec"`
}

type ReplicaSetSpec struct {
	Replicas int                `json:"replicas"`
	Selector ReplicaSetSelector `json:"selector"`
	Template PodTemplateSpec    `json:"template"`
}

type ReplicaSetStatus struct {
	SpecReplicas  int `json:"specReplicas"`
	ReadyReplicas int `json:"readyReplicas"`
}

type ReplicaSet struct {
	ApiVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	MetaData   MetaData         `json:"metadata"`
	Spec       ReplicaSetSpec   `json:"spec"`
	Status     ReplicaSetStatus `json:"status"`
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
