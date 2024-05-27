package apiobj

type JobSpec struct {
	Replicas int                `yaml:"replicas" json:"replicas"`
	Selector ReplicaSetSelector `yaml:"selector" json:"selector"`
	Template PodTemplateSpec    `yaml:"template" json:"template"`
}

type Job struct {
	ApiVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind" json:"kind"`
	MetaData   MetaData         `yaml:"metadata" json:"metadata"`
	Spec       ReplicaSetSpec   `yaml:"spec" json:"spec"`
	Status     ReplicaSetStatus `yaml:"status" json:"status"`
}

func (j *Job) GetKind() string {
	return j.Kind
}
func (j *Job) GetName() string {
	return j.MetaData.Name
}
func (j *Job) GetNamespace() string {
	return j.MetaData.Namespace
}
func (j *Job) SetNamespace(namespace string) {
	j.MetaData.Namespace = namespace
}
