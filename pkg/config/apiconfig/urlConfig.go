package apiconfig

//import "minik8s/pkg/apiserver/app/apiobj"

const (
	URL_AllNodes    = "/api/v1/nodes"
	URL_Node        = "/api/v1/nodes/:name"
	URL_NodeAllPods = "/api/v1/nodes/:name/pods"
	URL_NodeStatus  = "/api/v1/nodes/:name/status"

	URL_GlobalPods = "/api/v1/pods"
	URL_AllPods    = "/api/v1/namespaces/:namespace/pods"
	URL_Pod        = "/api/v1/namespaces/:namespace/pods/:name"
	URL_PodStatus  = "/api/v1/namespaces/:namespace/pods/:name/status"

	URL_GlobalServices = "/api/v1/services"
	URL_AllServices    = "/api/v1/namespaces/:namespace/services"
	URL_Service        = "/api/v1/namespaces/:namespace/services/:name"
	URL_ServiceStatus  = "/api/v1/namespaces/:namespace/services/:name/status"

	URL_GlobalReplicaSets = "/api/v1/replicasets"
	URL_AllReplicaSets    = "/api/v1/namespaces/:namespace/replicasets"
	URL_ReplicaSet        = "/api/v1/namespaces/:namespace/replicasets/:name"
	URL_ReplicaSetStatus  = "/api/v1/namespaces/:namespace/replicasets/:name/status"

	URL_GlobalHpas = "/api/v1/hpas"
	URL_AllHpas    = "/api/v1/namespaces/:namespace/hpas"
	URL_Hpa        = "/api/v1/namespaces/:namespace/hpas/:name"
	URL_HpaStatus  = "/api/v1/namespaces/:namespace/hpas/:name/status"

	URL_GlobalDns = "/api/v1/dns"
	URL_AllDns    = "/api/v1/namespaces/:namespace/dns"
	URL_Dns       = "/api/v1/namespaces/:namespace/dns/:name"
	URL_DnsStatus = "/api/v1/namespaces/:namespace/dns/:name/status"

	URL_GlobalFunctions = "/api/v1/functions"
	URL_AllFunctions    = "/api/v1/namespaces/:namespace/functions"
	URL_Function        = "/api/v1/namespaces/:namespace/functions/:name"

	URL_GlobalWorkflows = "/api/v1/workflows"
	URL_AllWorkflows    = "/api/v1/namespaces/:namespace/workflows"
	URL_Workflow        = "/api/v1/namespaces/:namespace/workflows/:name"
	URL_WorkflowStatus  = "/api/v1/namespaces/:namespace/workflows/:name/status"

	URL_GlobalPVs = "/api/v1/pvs"
	URL_AllPVs    = "/api/v1/namespaces/:namespace/pvs"
	URL_PV        = "/api/v1/namespaces/:namespace/pvs/:name"

	URL_GlobalPVCs = "/api/v1/pvcs"
	URL_AllPVCs    = "/api/v1/namespaces/:namespace/pvcs"
	URL_PVC        = "/api/v1/namespaces/:namespace/pvcs/:name"
)

var Kind2URL = map[string]string{
	"Node":       URL_Node,
	"Pod":        URL_Pod,
	"Service":    URL_Service,
	"ReplicaSet": URL_ReplicaSet,
	"Hpa":        URL_Hpa,
	"Dns":        URL_Dns,
	"Function":   URL_Function,
	"Workflow":   URL_Workflow,
	"PV":         URL_PV,
	"PVC":        URL_PVC,
}
