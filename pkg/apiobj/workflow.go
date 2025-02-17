package apiobj

const (
	FunctionType = "Function"
	ChoiceType   = "Choice"
	EndNode      = "end"
)

type FuncNode struct {
	FuncNameSpace string   `yaml:"funcNameSpace" json:"funcNameSpace"`
	FuncName      string   `yaml:"funcName" json:"funcName"`
	FuncParam     []string `yaml:"funcParam" json:"funcParam"`
	NextNodeName  string   `yaml:"nextNodeName" json:"nextNodeName"`
}

type ChoiceNode struct {
	TrueNodeName    string                 `yaml:"trueNodeName" json:"trueNodeName"`
	TrueEntryParam  map[string]interface{} `yaml:"trueEntryParam" json:"trueEntryParam"`
	FalseNodeName   string                 `yaml:"falseNodeName" json:"falseNodeName"`
	FalseEntryParam map[string]interface{} `yaml:"falseEntryParam" json:"falseEntryParam"`
	ChoiceParam     string                 `yaml:"choiceParam" json:"choiceParam"`
	Expression      string                 `yaml:"expression" json:"expression"`
}

type WorkflowNode struct {
	Name       string     `yaml:"name" json:"name"`
	Type       string     `yaml:"type" json:"type"`
	FuncNode   FuncNode   `yaml:"funcNode" json:"funcNode"`
	ChoiceNode ChoiceNode `yaml:"choiceNode" json:"choiceNode"`
}

type WorkflowSpec struct {
	EntryName     string                 `yaml:"entryName" json:"entryName"`
	EntryParam    map[string]interface{} `yaml:"entryParam" json:"entryParam"`
	WorkflowNodes []WorkflowNode         `yaml:"workflowNodes" json:"workflowNodes"`
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
