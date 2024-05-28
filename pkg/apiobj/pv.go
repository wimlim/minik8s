package apiobj

const (
	NfsMntPath = "/mnt"
)

type PVCapacity struct {
	Storage string `yaml:"storage" json:"storage"`
}

type Nfs struct {
	Path   string `yaml:"path" json:"path"`
	Server string `yaml:"server" json:"server"`
}

type PVSpec struct {
	AccessModes      []string   `yaml:"accessModes" json:"accessModes"`
	Capacity         PVCapacity `yaml:"capacity" json:"capacity"`
	StorageClassName string     `yaml:"storageClassName" json:"storageClassName"`
	Nfs              Nfs        `yaml:"nfs" json:"nfs"`
}

type PV struct {
	APIVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
	Spec       PVSpec   `yaml:"spec" json:"spec"`
}

func (p *PV) GetKind() string {
	return p.Kind
}
func (p *PV) GetName() string {
	return p.MetaData.Name
}
func (p *PV) GetNamespace() string {
	return p.MetaData.Namespace
}
func (p *PV) SetNamespace(namespace string) {
	p.MetaData.Namespace = namespace
}
