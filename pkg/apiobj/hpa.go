package apiobj

type HpaCpuMetric struct {
	Target float64 `yaml:"target" json:"target"`
}
type HpaMemoryMetric struct {
	Target float64 `yaml:"target" json:"target"`
}
type HpaSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels" json:"matchLabels"`
}

type HpaMetric struct {
	CpuMetric    HpaCpuMetric    `yaml:"cpu" json:"cpu"`
	MemoryMetric HpaMemoryMetric `yaml:"memory" json:"memory"`
}

type HpaScaleTarget struct {
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
}

type HpaSpec struct {
	Selector       HpaSelector    `yaml:"selector" json:"selector"`
	ScaleTargetRef HpaScaleTarget `yaml:"scaleTargetRef" json:"scaleTargetRef"`
	MinReplicas    int            `yaml:"minReplicas" json:"minReplicas"`
	MaxReplicas    int            `yaml:"maxReplicas" json:"maxReplicas"`
}

type HpaStatus struct {
	Replicas      int `yaml:"replicas" json:"replicas"`
	CpuPercent    int `yaml:"cpuPercent" json:"cpuPercent"`
	MemoryPercent int `yaml:"memoryPercent" json:"memoryPercent"`
}

type Hpa struct {
	ApiVersion string    `yaml:"apiVersion" json:"apiVersion"`
	Kind       string    `yaml:"kind" json:"kind"`
	MetaData   MetaData  `yaml:"metadata" json:"metadata"`
	Spec       HpaSpec   `yaml:"spec" json:"spec"`
	Status     HpaStatus `yaml:"status" json:"status"`
}

func (h *Hpa) GetKind() string {
	return h.Kind
}
func (h *Hpa) GetName() string {
	return h.MetaData.Name
}
func (h *Hpa) GetNamespace() string {
	return h.MetaData.Namespace
}
func (h *Hpa) SetNamespace(namespace string) {
	h.MetaData.Namespace = namespace
}
