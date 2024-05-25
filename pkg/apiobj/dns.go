package apiobj

type Path struct {
	SubPath     string `yaml:"subPath" json:"subPath"`
	ServiceIp   string `yaml:"serviceIp" json:"serviceIp"`
	ServiceName string `yaml:"serviceName" json:"serviceName"`
	ServicePort int    `yaml:"servicePort" json:"servicePort"`
}

type DnsSpec struct {
	Host  string `yaml:"host" json:"host"`
	Paths []Path `yaml:"paths" json:"paths"`
}

type DnsStatus struct {
	Phase string `yaml:"phase" json:"phase"`
}

type Dns struct {
	ApiVersion string    `yaml:"apiVersion " json:"apiVersion"`
	Kind       string    `yaml:"kind" json:"kind"`
	MetaData   MetaData  `yaml:"metadata" json:"metadata"`
	Spec       DnsSpec   `yaml:"spec" json:"spec"`
	Status     DnsStatus `yaml:"status" json:"status"`
}

func (d *Dns) GetKind() string {
	return d.Kind
}
func (d *Dns) GetName() string {
	return d.MetaData.Name
}
func (d *Dns) GetNamespace() string {
	return d.MetaData.Namespace
}
func (d *Dns) SetNamespace(namespace string) {
	d.MetaData.Namespace = namespace
}
