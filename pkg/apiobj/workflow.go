package apiobj

type FuncNode struct {
	FuncNamespace string `yaml:"funcNamespace" json:"funcNamespace"`
	FuncName      string `yaml:"funcName" json:"funcName"`
	NextNodeName  string `yaml:"nextNodeName" json:"nextNodeName"`
}

type ChoiceNode struct {
	TrueNodeName  string `yaml:"trueNodeName" json:"trueNodeName"`
	FalseNodeName string `yaml:"falseNodeName" json:"falseNodeName"`
	Condition     string `yaml:"condition" json:"condition"`
}

type WorkflowNode struct {
	Name       string     `yaml:"name" json:"name"`
	Type       string     `yaml:"type" json:"type"`
	FuncNode   FuncNode   `yaml:"funcNode" json:"funcNode"`
	ChoiceNode ChoiceNode `yaml:"choiceNode" json:"choiceNode"`
}

type WorkflowSpec struct {
	EntryName     string         `yaml:"entryName" json:"entryName"`
	EntryParam    string         `yaml:"entryParam" json:"entryParam"`
	WorkflowNodes []WorkflowNode `yaml:"nodes" json:"nodes"`
}

type WorkflowStatus struct {
	Result string `yaml:"result" json:"result"`
}

type Workflow struct {
	ApiVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       string         `yaml:"kind" json:"kind"`
	MetaData   MetaData       `yaml:"metadata" json:"metadata"`
	Spec       WorkflowSpec   `yaml:"spec" json:"spec"`
	Status     WorkflowStatus `yaml:"status" json:"status"`
}

func (w *Workflow) GetKind() string {
	return w.Kind
}
func (w *Workflow) GetName() string {
	return w.MetaData.Name
}
func (w *Workflow) GetNamespace() string {
	return w.MetaData.Namespace
}
func (w *Workflow) SetNamespace(namespace string) {
	w.MetaData.Namespace = namespace
}
