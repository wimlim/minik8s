package apiconfig

//import "minik8s/pkg/apiserver/app/apiobj"

const(
	URL_AllNodes = "/api/v1/nodes"
	URL_Node = "/api/v1/nodes/:name"
	URL_NodeAllPods = "/api/v1/nodes/:name/pods"
	URL_NodeStatus = "/api/v1/nodes/:name/status"
	
	
	URL_GlobalPods = "/api/v1/pods"
	URL_AllPods = "/api/v1/namespaces/:namespace/pods"
	URL_Pod = "/api/v1/namespaces/:namespace/pods/:name"
	URL_PodStatus = "/api/v1/namespaces/:namespace/pods/:name/status"
	
	URL_AllServices = "/api/v1/namespaces/:namespace/services"
	URL_Service = "/api/v1/namespaces/:namespace/services/:name"
	URL_ServiceStatus = "/api/v1/namespaces/:namespace/services/:name/status"
 )
