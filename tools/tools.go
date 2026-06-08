package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// Pod tools details
var ListPod = mcp.NewTool(
	"list-pod",
	mcp.WithDescription("List the pod with status, label and instance"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the pod in the particular namespace"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Only return pods matching this label selector"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetPod = mcp.NewTool(
	"get-pod",
	mcp.WithDescription("Get the pod in particular namespace with status, label and instance"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the pod exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the pod to get details"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeletePod = mcp.NewTool(
	"delete-pod",
	mcp.WithDescription("Delete the pod in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the pod to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the pod to be deleted"),
	),
)

var UpdatePod = mcp.NewTool(
	"update-pod",
	mcp.WithDescription("Update the pod in particular namespace like label changes"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the pod to be updated"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pod to be updated"),
	),
	mcp.WithString(
		"label",
		mcp.Required(),
		mcp.Description("Label to be updated"),
	),
)

var CreatePod = mcp.NewTool(
	"create-pod",
	mcp.WithDescription("Create the pod in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the pod to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pod to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that pod"),
	),
	mcp.WithString(
		"containerNames",
		mcp.Required(),
		mcp.Description("Container Names for the pod"),
	),
	mcp.WithString(
		"containerImages",
		mcp.Required(),
		mcp.Description("Container Image for the pod"),
	),
	mcp.WithString(
		"containerPorts",
		mcp.Description("Container port details for the pod"),
	),
)

var PodLog = mcp.NewTool(
	"pod-log",
	mcp.WithDescription("Get the log for particular pod in the namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the pod is present"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pod to get log"),
	),
	mcp.WithNumber(
		"tailLine",
		mcp.Description("Number of log line to get"),
	),
	mcp.WithString(
		"containerName",
		mcp.Description("Container name (optional if the pod has only one container)"),
	),
)

// Namespace tools details
var ListNS = mcp.NewTool(
	"list-ns",
	mcp.WithDescription("List the namespace in the kubernetes cluster with status"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetNS = mcp.NewTool(
	"get-ns",
	mcp.WithDescription("Get the particular namespace in the kubernetes cluster with status"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the namespace to get details for"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteNS = mcp.NewTool(
	"delete-ns",
	mcp.WithDescription("Delete the particular namespace in the kubernetes cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the namespace to be deleted"),
	),
)

var UpdateNS = mcp.NewTool(
	"update-ns",
	mcp.WithDescription("Update the namespace like label and annotation changes"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the namespace to be updated"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be updated"),
	),
	mcp.WithString(
		"annotation",
		mcp.Description("annotation to be updated"),
	),
)

var CreateNS = mcp.NewTool(
	"create-ns",
	mcp.WithDescription("Create the namespace"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the namespace to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be add in hte namespace"),
	),
)

// Deployment tools details
var ListDeployment = mcp.NewTool(
	"list-deployment",
	mcp.WithDescription("List the deployment with available instance with label"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the deployment in the particular namespace"),
	),
	mcp.WithString(
		"label",
		mcp.Description("The deployment should be listed only if this particular label is exist"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetDeployment = mcp.NewTool(
	"get-deployment",
	mcp.WithDescription("Get the deployment in particular namespace with available instance and label"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace to get the deployment"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the deployment to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteDeployment = mcp.NewTool(
	"delete-deployment",
	mcp.WithDescription("Delete the deployment in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the deployment to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the deployment to be deleted"),
	),
)

var UpdateDeployment = mcp.NewTool(
	"update-deployment",
	mcp.WithDescription("Update the deployment in particular namespace like label, annotation, replica and image changes"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the deployment to be updated"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the deployment to be updated"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be updated"),
	),
	mcp.WithString(
		"annotation",
		mcp.Description("annotation to be updated"),
	),
	mcp.WithNumber(
		"replica",
		mcp.Description("Replica to be updated"),
	),
	mcp.WithString(
		"containerName",
		mcp.Description("Container Name to update the image"),
	),
	mcp.WithString(
		"image",
		mcp.Description("Image to be updated"),
	),
)

var CreateDeployment = mcp.NewTool(
	"create-deployment",
	mcp.WithDescription("Create the deployment in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the deployment to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the deployment to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that deployment"),
	),
	mcp.WithNumber(
		"replica",
		mcp.Description("Number of replica"),
	),
	mcp.WithString(
		"containerNames",
		mcp.Required(),
		mcp.Description("Container Names for the deployment"),
	),
	mcp.WithString(
		"containerImages",
		mcp.Required(),
		mcp.Description("Container Image for the deployment"),
	),
	mcp.WithString(
		"containerPorts",
		mcp.Description("Container port details for the deployment"),
	),
)

// Service tools details
var ListService = mcp.NewTool(
	"list-service",
	mcp.WithDescription("List the service with type"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the service in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetService = mcp.NewTool(
	"get-service",
	mcp.WithDescription("Get the particular service with type and IP"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the service should be listed"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the service to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteService = mcp.NewTool(
	"delete-service",
	mcp.WithDescription("Delete the particular service in the namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the service to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the service to be deleted"),
	),
)

var UpdateService = mcp.NewTool(
	"update-service",
	mcp.WithDescription("Update the service in particular namespace like selector label and type changes"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the service to be updated"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the service to be updated"),
	),
	mcp.WithString(
		"selectorLabel",
		mcp.Description("Selector label to be updated"),
	),
	mcp.WithString(
		"svctype",
		mcp.Description("Service type to be updated"),
	),
)

var CreateService = mcp.NewTool(
	"create-service",
	mcp.WithDescription("Create the service in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the service to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the service to be created"),
	),
	mcp.WithString(
		"selectorLabel",
		mcp.Required(),
		mcp.Description("Selector Label for the service"),
	),
	mcp.WithString(
		"svcPort",
		mcp.Required(),
		mcp.Description("Service port name and port details for service"),
	),
	mcp.WithString(
		"targetPort",
		mcp.Required(),
		mcp.Description("Target port details"),
	),
	mcp.WithString(
		"svcType",
		mcp.Description("Service type need to create, if not provided it will take default service type"),
	),
)

// Statefulset tools details
var ListStatefulset = mcp.NewTool(
	"list-statefulset",
	mcp.WithDescription("List the statefulset with available instance and label"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the statefulset in the particular namespace"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Get the statefulset only if this particular label is exist"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetStatefulset = mcp.NewTool(
	"get-statefulset",
	mcp.WithDescription("Get the particular statefulset in particular namespace with available instance and labels"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the statefulset should be listed"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the statefulset to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteStatefulset = mcp.NewTool(
	"delete-statefulset",
	mcp.WithDescription("Delete the particular statefulset in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the statefulset to be deletes"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the statefulset to be deleted"),
	),
)

var UpdateStatefulset = mcp.NewTool(
	"update-statefulset",
	mcp.WithDescription("Update the statefulset in particular namespace like label, annotation, replica and image changes"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the statefulset to be updated"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the statefulset to be updated"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be updated"),
	),
	mcp.WithString(
		"annotation",
		mcp.Description("annotation to be updated"),
	),
	mcp.WithNumber(
		"replica",
		mcp.Description("Replica to be updated"),
	),
	mcp.WithString(
		"containerName",
		mcp.Description("Container Name to update the image"),
	),
	mcp.WithString(
		"image",
		mcp.Description("Image to be updated"),
	),
)

var CreateStatefulset = mcp.NewTool(
	"create-statefulset",
	mcp.WithDescription("Create the statefulset in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the statefulset to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the statefulset to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that statefulset"),
	),
	mcp.WithString(
		"containerNames",
		mcp.Description("Container Names for the statefulset"),
	),
	mcp.WithString(
		"containerImages",
		mcp.Required(),
		mcp.Description("Container Image for the statefulset"),
	),
	mcp.WithNumber(
		"containerPorts",
		mcp.Description("Container port for the statefulset"),
	),
	mcp.WithString(
		"storageValue",
		mcp.Required(),
		mcp.Description("Pvc size for the statefulset"),
	),
	mcp.WithString(
		"mountPath",
		mcp.Required(),
		mcp.Description("mount path for the statefulset container to mount the pvc"),
	),
	mcp.WithString(
		"pvcName",
		mcp.Description("Name of the pvc for statefulset"),
	),
	mcp.WithString(
		"svcType",
		mcp.Description("Servcie type for statefulset service"),
	),
	mcp.WithNumber(
		"svcPort",
		mcp.Description("Service Port for the statefulset service"),
	),
	mcp.WithNumber(
		"replica",
		mcp.Description("Number of replica for statefulset"),
	),
)

// Daemonset tools details
var ListDaemonset = mcp.NewTool(
	"list-daemonset",
	mcp.WithDescription("List the daemonset with available instance and label"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the daemonset in the particular namespace"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Get the daemonset only if this particular label is exist"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetDaemonset = mcp.NewTool(
	"get-daemonset",
	mcp.WithDescription("Get the daemonset in particular namespace with available instance and label"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace to get the daemonset"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the daemonset to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteDaemonset = mcp.NewTool(
	"delete-daemonset",
	mcp.WithDescription("Delete the daemonset in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the daemonset to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the daemonset to be deleted"),
	),
)

var UpdateDaemonset = mcp.NewTool(
	"update-daemonset",
	mcp.WithDescription("Update the daemonset in particular namespace like label, annotation and image changes"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the daemonset to be updated"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the daemonset to be updated"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be updated"),
	),
	mcp.WithString(
		"annotation",
		mcp.Description("annotation to be updated"),
	),
	mcp.WithString(
		"containerName",
		mcp.Description("Container Name to update the image"),
	),
	mcp.WithString(
		"image",
		mcp.Description("Image to be updated"),
	),
)

var CreateDaemonset = mcp.NewTool(
	"create-daemonset",
	mcp.WithDescription("Create the daemonset in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in whcih the daemonset to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the daemonset to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that daemonset"),
	),
	mcp.WithString(
		"containerNames",
		mcp.Required(),
		mcp.Description("Container Names for the daemonset"),
	),
	mcp.WithString(
		"containerImages",
		mcp.Required(),
		mcp.Description("Container Image for the daemonset"),
	),
	mcp.WithString(
		"containerPorts",
		mcp.Description("Container port details for the daemonset"),
	),
)

// Configmap tools details
var ListConfigmap = mcp.NewTool(
	"list-configmap",
	mcp.WithDescription("List the configmap"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the configmap in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetConfigmap = mcp.NewTool(
	"get-configmap",
	mcp.WithDescription("Get the configmap in particular namespace with data"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the configmap to get"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the configmap to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteConfigmap = mcp.NewTool(
	"delete-configmap",
	mcp.WithDescription("Delete the configmap in particular namespace "),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the configmap to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the configmap to be deleted"),
	),
)

var CreateConfigmap = mcp.NewTool(
	"create-configmap",
	mcp.WithDescription("Create the configmap in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the configmap to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the configmap to be created"),
	),
	mcp.WithString(
		"data",
		mcp.Required(),
		mcp.Description("Data of the configmap to be created for"),
	),
)

// Secret tools details
var ListSecret = mcp.NewTool(
	"list-secret",
	mcp.WithDescription("List the secret"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the secret in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetSecret = mcp.NewTool(
	"get-secret",
	mcp.WithDescription("Get the secret in particular namespace with data"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the secret to get"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the secret to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteSecret = mcp.NewTool(
	"delete-secret",
	mcp.WithDescription("Delete the secret in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the secret to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the secret to be deleted"),
	),
)

var CreateSecret = mcp.NewTool(
	"create-secret",
	mcp.WithDescription("Create the secret in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the secret to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the secret to be created"),
	),
	mcp.WithString(
		"data",
		mcp.Required(),
		mcp.Description("Data of the secret to be created for"),
	),
)

// Node tools details
var ListNode = mcp.NewTool(
	"list-node",
	mcp.WithDescription("List the node in the kubernetes cluster with status"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetNode = mcp.NewTool(
	"get-node",
	mcp.WithDescription("Get the particular node in the kubernetes cluster with status, kubernetes version, architecture and os"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the node to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteNode = mcp.NewTool(
	"delete-node",
	mcp.WithDescription("Delete the particular node in the kubernetes cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the node to be deleted"),
	),
)

var UpdateNode = mcp.NewTool(
	"update-node",
	mcp.WithDescription("Update the node like label changes"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the node to be updated"),
	),
	mcp.WithString(
		"label",
		mcp.Required(),
		mcp.Description("Label to be updated"),
	),
)

// ServiceAccount tools details
var ListSA = mcp.NewTool(
	"list-serviceAccount",
	mcp.WithDescription("List the serviceAccount"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the serviceAccount in the particular namespace"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label of the serviceAccount, if we need to list the service account with particualr label exist"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetSA = mcp.NewTool(
	"get-serviceAccount",
	mcp.WithDescription("Get the serviceAccount in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the serviceAccount to get"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the serviceAccount to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteSA = mcp.NewTool(
	"delete-serviceAccount",
	mcp.WithDescription("Delete the serviceAccount in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the serviceAccount to delete"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the serviceAccount to delete"),
	),
)

var CreateSA = mcp.NewTool(
	"create-serviceAccount",
	mcp.WithDescription("Create the serviceAccount in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the serviceAccount to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the serviceAccount to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label of the serviceAccount, if we need to create the service account with particualr label"),
	),
)

// PVC tools details
var ListPVC = mcp.NewTool(
	"list-pvc",
	mcp.WithDescription("List the pvc"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the pvc in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetPVC = mcp.NewTool(
	"get-pvc",
	mcp.WithDescription("Get the pvc in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the pvc to get"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pvc to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeletePVC = mcp.NewTool(
	"delete-pvc",
	mcp.WithDescription("Delete the pvc in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace in which the pvc to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pvc to delete"),
	),
)

var UpdatePVC = mcp.NewTool(
	"update-pvc",
	mcp.WithDescription("update the pvc in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the pvc to update"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pvc to update"),
	),
	mcp.WithString(
		"size",
		mcp.Required(),
		mcp.Description("size of the pvc to update"),
	),
)

var CreatePVC = mcp.NewTool(
	"create-pvc",
	mcp.WithDescription("Create the pvc in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the pvc to create"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pvc to create"),
	),
	mcp.WithString(
		"size",
		mcp.Required(),
		mcp.Description("Size of the pvc to create"),
	),
	mcp.WithString(
		"storageClass",
		mcp.Description("Name of the storageClass for pvc to create (optional, uses default if empty)"),
	),
	mcp.WithString(
		"accessMode",
		mcp.Description("AccessModes of the pvc to create"),
	),
)

// PV tools details
var ListPV = mcp.NewTool(
	"list-pv",
	mcp.WithDescription("List the entire pv"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetPV = mcp.NewTool(
	"get-pv",
	mcp.WithDescription("Get the pv in particular name"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pv to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeletePV = mcp.NewTool(
	"delete-pv",
	mcp.WithDescription("Delete the particular pv"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the pv to delete"),
	),
)

// Role tools details
var ListRole = mcp.NewTool(
	"list-role",
	mcp.WithDescription("List the role"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the role in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetRole = mcp.NewTool(
	"get-role",
	mcp.WithDescription("Get the role in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the role to get"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the role to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteRole = mcp.NewTool(
	"delete-role",
	mcp.WithDescription("Delete role in particular namespace"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the role to be deleted"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the role"),
	),
)

// RoleBinding tools details
var ListRB = mcp.NewTool(
	"list-rolebinding",
	mcp.WithDescription("List the rolebinding"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the rolebinding in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetRB = mcp.NewTool(
	"get-rolebinding",
	mcp.WithDescription("Get the rolebinding in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the rolebinding to get"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the rolebinding to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteRB = mcp.NewTool(
	"delete-rolebinding",
	mcp.WithDescription("Delete rolebinding in particular namespace"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the rolebinding to be deleted"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Namespace of the rolebinding"),
	),
)

// ClusterRole tools details
var ListCR = mcp.NewTool(
	"list-clusterrole",
	mcp.WithDescription("List all the clusterrole in the cluster"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetCR = mcp.NewTool(
	"get-clusterrole",
	mcp.WithDescription("Get the particular clusterrole"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the clusterrole to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteCR = mcp.NewTool(
	"delete-clusterrole",
	mcp.WithDescription("Delete cluster role in the cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the cluster role to be deleted"),
	),
)

// ClusterRoleBinding tools details
var ListCRB = mcp.NewTool(
	"list-clusterrolebinding",
	mcp.WithDescription("List all the clusterrolebinding in the cluster"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetCRB = mcp.NewTool(
	"get-clusterrolebinding",
	mcp.WithDescription("Get the particular clusterrolebinding"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the clusterrolebinding to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteCRB = mcp.NewTool(
	"delete-clusterrolebinding",
	mcp.WithDescription("Delete cluster role binding in the cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the cluster role binding to be deleted"),
	),
)

// StorageClass tools details
var ListSC = mcp.NewTool(
	"list-storageClass",
	mcp.WithDescription("List the storageClass in the entier cluster"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetSC = mcp.NewTool(
	"get-storageClass",
	mcp.WithDescription("Get the particular storafeClass"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the storageClass to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteSC = mcp.NewTool(
	"delete-storeageclass",
	mcp.WithDescription("Delete storage class in the cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the storageclass to be deleted"),
	),
)

// cutom tool details
var Custom = mcp.NewTool(
	"custom",
	mcp.WithDescription("Working with custom resource"),
	mcp.WithString(
		"kind",
		mcp.Required(),
		mcp.Description("Kind of the custom resource"),
	),
	mcp.WithString(
		"method",
		mcp.Required(),
		mcp.Description("Method to work on that custom resource"),
	),
	mcp.WithString(
		"name",
		mcp.Description("Name of the custom resource"),
	),
	mcp.WithString(
		"namespace",
		mcp.Description("Namespace in which custom resource exits"),
	),
	mcp.WithString(
		"jsondata",
		mcp.Description("Json data to create the custom resource"),
	),
)

// CRD tools details
var ListCRD = mcp.NewTool(
	"list-crd",
	mcp.WithDescription("List the crds in the cluster"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetCRD = mcp.NewTool(
	"get-crd",
	mcp.WithDescription("Get the particular crd in the cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the crd to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteCRD = mcp.NewTool(
	"delete-crd",
	mcp.WithDescription("Delete the particular crd in the cluster"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the crd to be deleted"),
	),
)

var CreateCRDWithJson = mcp.NewTool(
	"create-crd-with-json",
	mcp.WithDescription("Create crd with json or yaml data"),
	mcp.WithString(
		"jsondata",
		mcp.Required(),
		mcp.Description("Json or yaml data to create crd"),
	),
)

// Create Resource with Json tool details
var CreateResourceWithJSon = mcp.NewTool(
	"create-resource-with-json",
	mcp.WithDescription("Create any resource in kubernetes with json or yaml data (supports both JSON and YAML)"),
	mcp.WithString(
		"jsondata",
		mcp.Required(),
		mcp.Description("Json or yaml data to create resource"),
	),
)

// Ingress tools details
var ListIngress = mcp.NewTool(
	"list-ingress",
	mcp.WithDescription("List the ingress with hosts and labels"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the ingress in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetIngress = mcp.NewTool(
	"get-ingress",
	mcp.WithDescription("Get the particular ingress in particular namespace with hosts and labels"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the ingress exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the ingress to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteIngress = mcp.NewTool(
	"delete-ingress",
	mcp.WithDescription("Delete the particular ingress in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the ingress to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the ingress to be deleted"),
	),
)

var CreateIngress = mcp.NewTool(
	"create-ingress",
	mcp.WithDescription("Create the ingress in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the ingress to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the ingress to be created"),
	),
	mcp.WithString(
		"host",
		mcp.Description("Host for the ingress rule"),
	),
	mcp.WithString(
		"serviceName",
		mcp.Description("Name of the backend service"),
	),
	mcp.WithNumber(
		"servicePort",
		mcp.Description("Port of the backend service"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that ingress"),
	),
	mcp.WithString(
		"ingressClassName",
		mcp.Description("Ingress class name (default: nginx)"),
	),
)

// HPA tools details
var ListHPA = mcp.NewTool(
	"list-hpa",
	mcp.WithDescription("List the hpa with target and replicas"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the hpa in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetHPA = mcp.NewTool(
	"get-hpa",
	mcp.WithDescription("Get the particular hpa in particular namespace with metrics"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the hpa exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the hpa to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteHPA = mcp.NewTool(
	"delete-hpa",
	mcp.WithDescription("Delete the particular hpa in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the hpa to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the hpa to be deleted"),
	),
)

var CreateHPA = mcp.NewTool(
	"create-hpa",
	mcp.WithDescription("Create the hpa in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the hpa to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the hpa to be created"),
	),
	mcp.WithString(
		"targetKind",
		mcp.Description("Kind of the target workload (Deployment, StatefulSet, etc.)"),
	),
	mcp.WithString(
		"targetName",
		mcp.Required(),
		mcp.Description("Name of the target workload to scale"),
	),
	mcp.WithNumber(
		"minReplicas",
		mcp.Description("Minimum number of replicas"),
	),
	mcp.WithNumber(
		"maxReplicas",
		mcp.Required(),
		mcp.Description("Maximum number of replicas"),
	),
	mcp.WithNumber(
		"cpuTarget",
		mcp.Description("Target CPU utilization percentage"),
	),
)

// Job tools details
var ListJob = mcp.NewTool(
	"list-job",
	mcp.WithDescription("List the job with completions and status"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the job in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetJob = mcp.NewTool(
	"get-job",
	mcp.WithDescription("Get the particular job in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the job exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the job to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteJob = mcp.NewTool(
	"delete-job",
	mcp.WithDescription("Delete the particular job in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the job to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the job to be deleted"),
	),
)

var CreateJob = mcp.NewTool(
	"create-job",
	mcp.WithDescription("Create the job in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the job to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the job to be created"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that job"),
	),
	mcp.WithString(
		"containerName",
		mcp.Required(),
		mcp.Description("Container name for the job"),
	),
	mcp.WithString(
		"containerImage",
		mcp.Required(),
		mcp.Description("Container image for the job"),
	),
	mcp.WithString(
		"command",
		mcp.Description("Command to run in the container (space separated)"),
	),
	mcp.WithNumber(
		"backoffLimit",
		mcp.Description("Backoff limit for the job (default: 6)"),
	),
	mcp.WithNumber(
		"completions",
		mcp.Description("Number of completions (default: 1)"),
	),
	mcp.WithNumber(
		"parallelism",
		mcp.Description("Parallelism for the job (default: 1)"),
	),
)

// CronJob tools details
var ListCronJob = mcp.NewTool(
	"list-cronjob",
	mcp.WithDescription("List the cronjob with schedule and status"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the cronjob in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetCronJob = mcp.NewTool(
	"get-cronjob",
	mcp.WithDescription("Get the particular cronjob in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the cronjob exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the cronjob to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteCronJob = mcp.NewTool(
	"delete-cronjob",
	mcp.WithDescription("Delete the particular cronjob in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the cronjob to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the cronjob to be deleted"),
	),
)

var CreateCronJob = mcp.NewTool(
	"create-cronjob",
	mcp.WithDescription("Create the cronjob in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the cronjob to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the cronjob to be created"),
	),
	mcp.WithString(
		"schedule",
		mcp.Required(),
		mcp.Description("Schedule for the cronjob (e.g., '*/5 * * * *')"),
	),
	mcp.WithString(
		"label",
		mcp.Description("Label to be added in that cronjob"),
	),
	mcp.WithString(
		"containerName",
		mcp.Required(),
		mcp.Description("Container name for the cronjob"),
	),
	mcp.WithString(
		"containerImage",
		mcp.Required(),
		mcp.Description("Container image for the cronjob"),
	),
	mcp.WithString(
		"command",
		mcp.Description("Command to run in the container (space separated)"),
	),
)

// NetworkPolicy tools details
var ListNetworkPolicy = mcp.NewTool(
	"list-networkpolicy",
	mcp.WithDescription("List the networkpolicy with pod selector and policy types"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the networkpolicy in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetNetworkPolicy = mcp.NewTool(
	"get-networkpolicy",
	mcp.WithDescription("Get the particular networkpolicy in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the networkpolicy exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the networkpolicy to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var DeleteNetworkPolicy = mcp.NewTool(
	"delete-networkpolicy",
	mcp.WithDescription("Delete the particular networkpolicy in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the networkpolicy to be deleted"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the networkpolicy to be deleted"),
	),
)

var CreateNetworkPolicy = mcp.NewTool(
	"create-networkpolicy",
	mcp.WithDescription("Create the networkpolicy in particular namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the networkpolicy to be created"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the networkpolicy to be created"),
	),
	mcp.WithString(
		"podSelector",
		mcp.Description("Pod selector labels (e.g., 'app=myapp,tier=frontend')"),
	),
	mcp.WithString(
		"policyTypes",
		mcp.Description("Policy types (e.g., 'Ingress,Egress')"),
	),
)

// Event tools
var ListEvent = mcp.NewTool(
	"list-event",
	mcp.WithDescription("List events across namespaces or in a specific namespace"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the event in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetEvent = mcp.NewTool(
	"get-event",
	mcp.WithDescription("Get a particular event in a namespace"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the event exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the event to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

// ResourceQuota tools
var ListResourceQuota = mcp.NewTool(
	"list-resourcequota",
	mcp.WithDescription("List resource quotas across namespaces or in a specific namespace"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the resourcequota in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetResourceQuota = mcp.NewTool(
	"get-resourcequota",
	mcp.WithDescription("Get a particular resource quota with hard and used values"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the resourcequota exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the resourcequota to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

// LimitRange tools
var ListLimitRange = mcp.NewTool(
	"list-limitrange",
	mcp.WithDescription("List limit ranges across namespaces or in a specific namespace"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the limitrange in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetLimitRange = mcp.NewTool(
	"get-limitrange",
	mcp.WithDescription("Get a particular limit range with limits configuration"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the limitrange exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the limitrange to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

// Endpoint tools
var ListEndpoint = mcp.NewTool(
	"list-endpoint",
	mcp.WithDescription("List endpoints across namespaces or in a specific namespace"),
	mcp.WithString(
		"namespace",
		mcp.Description("List the endpoint in the particular namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetEndpoint = mcp.NewTool(
	"get-endpoint",
	mcp.WithDescription("Get a particular endpoint with addresses and ports"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace in which the endpoint exists"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the endpoint to get"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

// ComponentStatus tools
var ListComponentStatus = mcp.NewTool(
	"list-componentstatus",
	mcp.WithDescription("List health status of Kubernetes control plane components"),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

var GetComponentStatus = mcp.NewTool(
	"get-componentstatus",
	mcp.WithDescription("Get health status of a specific Kubernetes control plane component"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Name of the component (e.g., etcd-0, kube-apiserver)"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full K8s object, yaml=full K8s object"),
	),
)

// Top tools (metrics)
var TopPod = mcp.NewTool(
	"top-pod",
	mcp.WithDescription("Show CPU and memory usage of pods (requires metrics-server)"),
	mcp.WithString(
		"namespace",
		mcp.Description("Show pods in the particular namespace"),
	),
)

var TopNode = mcp.NewTool(
	"top-node",
	mcp.WithDescription("Show CPU and memory usage of nodes (requires metrics-server)"),
)

// Cluster Health tools
var GetClusterHealth = mcp.NewTool(
	"cluster-health",
	mcp.WithDescription("Show overall cluster health: node summary + control plane component pod status"),
)

var ListNodeHealth = mcp.NewTool(
	"node-health",
	mcp.WithDescription("Show detailed health of all nodes: ready status, kubelet version, allocatable resources, labels"),
)

// Diagnosis tools
var DescribePod = mcp.NewTool(
	"describe-pod",
	mcp.WithDescription("Deep inspect a pod: container states, restart counts, conditions, events, resource requests/limits"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace of the pod"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the pod"),
	),
)

var DescribeNode = mcp.NewTool(
	"describe-node",
	mcp.WithDescription("Deep inspect a node: conditions with timestamps, capacity, allocatable, pods, taints, annotations"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the node"),
	),
)

var ListNodePods = mcp.NewTool(
	"list-node-pods",
	mcp.WithDescription("List all pods running on a specific node"),
	mcp.WithString(
		"node",
		mcp.Required(),
		mcp.Description("Name of the node"),
	),
)

var DescribeService = mcp.NewTool(
	"describe-service",
	mcp.WithDescription("Deep inspect a service: endpoints, ports, selector, type, session affinity, events"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace of the service"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description(		"The name of the service"),
	),
)

var DescribeDeployment = mcp.NewTool(
	"describe-deployment",
	mcp.WithDescription("Deep inspect a deployment: conditions (Available/Progressing/ReplicaFailure), strategy, selector, rollout config, events, pods"),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("The namespace of the deployment"),
	),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("The name of the deployment"),
	),
)

// API Server Diagnosis tools
var CheckAPIServerHealth = mcp.NewTool(
	"check-apiserver-health",
	mcp.WithDescription("Probe API server health endpoints: /healthz, /livez, /readyz with verbose output. Works for both Pod-based and binary deployments."),
)

var CheckAPIServerMetrics = mcp.NewTool(
	"check-apiserver-metrics",
	mcp.WithDescription("Fetch and analyze API server /metrics: inflight requests, request rate by verb, error rate (4xx/5xx), and latency percentiles (p50/p90/p99). Works for both Pod-based and binary deployments."),
)

// Helm tools
var ListHelmReleases = mcp.NewTool(
	"helm-list-releases",
	mcp.WithDescription("List Helm releases"),
	mcp.WithString(
		"namespace",
		mcp.Description("Filter by namespace"),
	),
	mcp.WithBoolean(
		"allNamespaces",
		mcp.Description("List releases across all namespaces"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full object, yaml=full object"),
	),
)

var GetHelmRelease = mcp.NewTool(
	"helm-get-release",
	mcp.WithDescription("Get details of a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
	mcp.WithString(
		"output",
		mcp.Description("Output format: empty=summary, json=full object, yaml=full object"),
	),
)

var GetHelmReleaseValues = mcp.NewTool(
	"helm-get-values",
	mcp.WithDescription("Get the values of a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
)

var InstallHelmRelease = mcp.NewTool(
	"helm-install",
	mcp.WithDescription("Install a Helm chart"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Target namespace"),
	),
	mcp.WithString(
		"chart",
		mcp.Required(),
		mcp.Description("Chart reference (e.g., stable/nginx-ingress, or path)"),
	),
	mcp.WithString(
		"version",
		mcp.Description("Chart version"),
	),
	mcp.WithString(
		"values",
		mcp.Description("Inline YAML/JSON values to pass to the chart"),
	),
	mcp.WithString(
		"set",
		mcp.Description("Set values on the command line (key1=val1,key2=val2)"),
	),
)

var UpgradeHelmRelease = mcp.NewTool(
	"helm-upgrade",
	mcp.WithDescription("Upgrade a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
	mcp.WithString(
		"chart",
		mcp.Required(),
		mcp.Description("New chart reference"),
	),
	mcp.WithString(
		"version",
		mcp.Description("Chart version"),
	),
	mcp.WithString(
		"values",
		mcp.Description("Inline YAML/JSON values to pass to the chart"),
	),
	mcp.WithString(
		"set",
		mcp.Description("Set values on the command line (key1=val1,key2=val2)"),
	),
	mcp.WithBoolean(
		"reuseValues",
		mcp.Description("Reuse previous values (default true)"),
	),
)

var UninstallHelmRelease = mcp.NewTool(
	"helm-uninstall",
	mcp.WithDescription("Uninstall a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
)

var RollbackHelmRelease = mcp.NewTool(
	"helm-rollback",
	mcp.WithDescription("Rollback a Helm release to a previous revision"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
	mcp.WithNumber(
		"revision",
		mcp.Description("Revision to rollback to (default: previous)"),
	),
)

var GetHelmReleaseHistory = mcp.NewTool(
	"helm-history",
	mcp.WithDescription("Get the revision history of a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
	mcp.WithNumber(
		"max",
		mcp.Description("Maximum number of revisions to show"),
	),
)

var GetHelmReleaseManifest = mcp.NewTool(
	"helm-get-manifest",
	mcp.WithDescription("Get the rendered manifest of a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
)

var GetHelmReleaseNotes = mcp.NewTool(
	"helm-get-notes",
	mcp.WithDescription("Get the notes of a Helm release"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Release name"),
	),
	mcp.WithString(
		"namespace",
		mcp.Required(),
		mcp.Description("Release namespace"),
	),
)

var ListHelmRepos = mcp.NewTool(
	"helm-list-repos",
	mcp.WithDescription("List Helm repositories"),
)

var AddHelmRepo = mcp.NewTool(
	"helm-add-repo",
	mcp.WithDescription("Add a Helm repository"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Repository name"),
	),
	mcp.WithString(
		"url",
		mcp.Required(),
		mcp.Description("Repository URL"),
	),
)

var RemoveHelmRepo = mcp.NewTool(
	"helm-remove-repo",
	mcp.WithDescription("Remove a Helm repository"),
	mcp.WithString(
		"name",
		mcp.Required(),
		mcp.Description("Repository name"),
	),
)

var UpdateHelmRepos = mcp.NewTool(
	"helm-update-repos",
	mcp.WithDescription("Update Helm repositories (download index files for all or a specific repo)"),
	mcp.WithString(
		"name",
		mcp.Description("Repository name to update (omit to update all)"),
	),
)

var SearchHelmRepo = mcp.NewTool(
	"helm-search-repo",
	mcp.WithDescription("Search for charts in Helm repositories"),
	mcp.WithString(
		"keyword",
		mcp.Required(),
		mcp.Description("Search keyword"),
	),
)
