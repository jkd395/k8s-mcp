package pvc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	"k8s-mcp/kubernetes/output"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type pvcData struct {
	Name         string   `json:"name,omitempty"`
	Namespace    string   `json:"namespace,omitempty"`
	Status       string   `json:"status,omitempty"`
	Capacity     string   `json:"capacity,omitempty"`
	AccessMode   []string `json:"accessMode,omitempty"`
	StorageClass string   `json:"storageClass,omitempty"`
	Volume       string   `json:"volume,omitempty"`
}

func ListPVC(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var pvcList []pvcData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing pvc in %s: %v", namespace.Name, err)), nil
			}
			if outFmt != "" {
				result, err := output.Format(outFmt, pvcs.Items)
				if err != nil {
					return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
				}
				return mcp.NewToolResultText(result), nil
			}
			for _, pvc := range pvcs.Items {
				pvcList = append(pvcList, pvcData{
					Name:      pvc.Name,
					Namespace: pvc.Namespace,
					Status:    string(pvc.Status.Phase),
				})
			}
		}
		mcpOutput, err := json.MarshalIndent(pvcList, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		pvcs, err := clientset.CoreV1().PersistentVolumeClaims(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing pvc in %s: %v", ns, err)), nil
		}
		if outFmt != "" {
			result, err := output.Format(outFmt, pvcs.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, pvc := range pvcs.Items {
			capacity := ""
			if pvc.Spec.Resources.Requests != nil {
				if qty, ok := pvc.Spec.Resources.Requests[v1.ResourceStorage]; ok {
					capacity = qty.String()
				}
			}
			pvcList = append(pvcList, pvcData{
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Capacity:  capacity,
				Status:    string(pvc.Status.Phase),
			})
		}
		mcpOutput, err := json.MarshalIndent(pvcList, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetPVC(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting pvc in %s/%s: %v", ns, name, err)), nil
	}
	outFmt := request.GetString("output", "")
	if outFmt != "" {
		result, err := output.Format(outFmt, pvc)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}
	var accMode []string
	for _, mode := range pvc.Spec.AccessModes {
		accMode = append(accMode, string(mode))
	}
	capacity := ""
	if pvc.Spec.Resources.Requests != nil {
		if q, ok := pvc.Spec.Resources.Requests[v1.ResourceStorage]; ok {
			capacity = q.String()
		}
	}
	storageClass := ""
	if pvc.Spec.StorageClassName != nil {
		storageClass = *pvc.Spec.StorageClassName
	}
	res := pvcData{
		Name:         pvc.Name,
		Namespace:    pvc.Namespace,
		Capacity:     capacity,
		AccessMode:   accMode,
		StorageClass: storageClass,
		Volume:       pvc.Spec.VolumeName,
		Status:       string(pvc.Status.Phase),
	}
	mcpOutput, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeletePVC(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.CoreV1().PersistentVolumeClaims(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting pvc in %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully pvc %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func UpdatePVC(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	size, err := request.RequireString("size")
	if err != nil {
		output := fmt.Sprintf("Provide size for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting pvc in %s/%s: %v", ns, name, err)), nil
	}

	qty, err := resource.ParseQuantity(size)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Invalid pvc size: %s", size)), nil
	}

	pvc.Spec.Resources.Requests[v1.ResourceStorage] = qty

	updatePVC, err := clientset.CoreV1().PersistentVolumeClaims(ns).Update(ctx, pvc, metav1.UpdateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in updating pvc in %s/%s with size %s: %v", ns, name, size, err)), nil
	}
	output := fmt.Sprintf("Successfully pvc %s/%s updated with size %s", updatePVC.Namespace, updatePVC.Name, size)
	return mcp.NewToolResultText(string(output)), nil
}

func CreatePVC(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	size, err := request.RequireString("size")
	if err != nil {
		output := fmt.Sprintf("Provide size for pvc")
		return mcp.NewToolResultText(string(output)), nil
	}
	accessModes := request.GetString("accessMode", "ReadWriteOnce")
	storageClass := request.GetString("storageClass", "")
	var accMode []v1.PersistentVolumeAccessMode
	am := strings.Split(accessModes, ",")
	for _, mode := range am {
		accMode = append(accMode, v1.PersistentVolumeAccessMode(mode))
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}

	var storageClassName *string
	if storageClass != "" {
		storageClassName = &storageClass
	}

	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: accMode,
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(size),
				},
			},
			StorageClassName: storageClassName,
		},
	}
	createPVC, err := clientset.CoreV1().PersistentVolumeClaims(ns).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating pvc %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully pvc %s/%s is created", createPVC.Namespace, createPVC.Name)
	return mcp.NewToolResultText(string(output)), nil
}
