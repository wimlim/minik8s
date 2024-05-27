package apiobj

type PVCResource struct {
	Requests PVCRequest `yaml:"requests" json:"requests"`
}

type PVCRequest struct {
	Storage string `yaml:"storage" json:"storage"`
}

type PVCSpec struct {
	AccessModes      []string    `yaml:"accessModes" json:"accessModes"`
	Resources        PVCResource `yaml:"resources" json:"resources"`
	StorageClassName string      `yaml:"storageClassName" json:"storageClassName"`
}

type PVC struct {
	APIVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
	Spec       PVCSpec  `yaml:"spec" json:"spec"`
}

func (p *PVC) GetKind() string {
	return p.Kind
}
func (p *PVC) GetName() string {
	return p.MetaData.Name
}
func (p *PVC) GetNamespace() string {
	return p.MetaData.Namespace
}
func (p *PVC) SetNamespace(namespace string) {
	p.MetaData.Namespace = namespace
}
