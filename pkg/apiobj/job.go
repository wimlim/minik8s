package apiobj

type JobSpec struct {
	Partition     string `yaml:"partition" json:"partition"`
	Nodes         int    `yaml:"nodes" json:"nodes"`
	NtasksPerNode int    `yaml:"ntasksPerNode" json:"ntasksPerNode"`
	CpusPerTask   int    `yaml:"cpusPerTask" json:"cpusPerTask"`
	Gres          string `yaml:"gres" json:"gres"`
}

type Job struct {
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
	Spec       JobSpec  `yaml:"spec" json:"spec"`
	File       string   `yaml:"file" json:"file"`
	Script     string   `yaml:"script" json:"script"`
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
