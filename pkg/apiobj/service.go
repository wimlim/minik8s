package apiobj

const (
	ServiceKind   = "Service"
	TypeClusterIP = "ClusterIP"
	TypeNodePort  = "NodePort"
)

type ServicePort struct {
	Name       string `yaml:"name" json:"name"`
	Protocol   string `yaml:"protocol" json:"protocol"`
	Port       int    `yaml:"port" json:"port"`
	TargetPort int    `yaml:"targetPort" json:"targetPort"`
	NodePort   int    `yaml:"nodePort" json:"nodePort"`
}
type ServiceSpec struct {
	Selector  map[string]string `yaml:"selector" json:"selector"`
	Ports     []ServicePort     `yaml:"ports" json:"ports"`
	Type      string            `yaml:"type" json:"type"`
	ClusterIP string            `yaml:"clusterIP" json:"clusterIP"`
}
type ServiceStatus struct {
	LoadBalancer map[string]string `yaml:"loadBalancer" json:"loadBalancer"`
}
type Service struct {
	ApiVersion string        `yaml:"apiVersion" json:"apiVersion"`
	Kind       string        `yaml:"kind" json:"kind"`
	MetaData   MetaData      `yaml:"metadata" json:"metadata"`
	Spec       ServiceSpec   `yaml:"spec" json:"spec"`
	Status     ServiceStatus `yaml:"status" json:"status"`
}

func (s *Service) GetKind() string {
	return s.Kind
}
func (s *Service) GetName() string {
	return s.MetaData.Name
}
func (s *Service) GetNamespace() string {
	return s.MetaData.Namespace
}
func (s *Service) SetNamespace(namespace string) {
	s.MetaData.Namespace = namespace
}
