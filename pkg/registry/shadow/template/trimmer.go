/*
Copyright 2021 The Clusternet Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"fmt"
	"regexp"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"
)

func trimCommonMetadata(result *unstructured.Unstructured) {
	// metadata.uid cannot be trimmed, which will be used for checking when patching.
	// metadata.uid is set to empty when deploying to child clusters.
	unstructured.RemoveNestedField(result.Object, "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(result.Object, "metadata", "managedFields")
	unstructured.RemoveNestedField(result.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(result.Object, "metadata", "selfLink")
}

func trimCoreService(result *unstructured.Unstructured) {
	serviceType, found, err := unstructured.NestedString(result.Object, "spec", "type")
	if !found || err != nil {
		return
	}

	switch corev1.ServiceType(serviceType) {
	case corev1.ServiceTypeNodePort, corev1.ServiceTypeLoadBalancer:
		// corev1.Service will init node ports when creating NodePort or LoadBalancer
		items, found, err := unstructured.NestedSlice(result.Object, "spec", "ports")
		if !found || err != nil {
			return
		}
		for _, item := range items {
			servicePort, ok := item.(map[string]interface{})
			if !ok {
				return
			}
			unstructured.RemoveNestedField(servicePort, "nodePort")
		}

		err = unstructured.SetNestedSlice(result.Object, items, "spec", "ports")
		if err != nil {
			klog.ErrorDepth(2, fmt.Sprintf("failed to trim Service %s/%s: %v", result.GetNamespace(), result.GetName(), err))
		}
	}
}

func trimBatchJob(result *unstructured.Unstructured) {
	unstructured.RemoveNestedField(result.Object, "spec", "selector", "matchLabels", "controller-uid")
	unstructured.RemoveNestedField(result.Object, "spec", "template", "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(result.Object, "spec", "template", "metadata", "labels", "controller-uid")
}

var defaultTokenVolumeNameRe = "default-token-[0-9a-z]{5}"

func trimCoreV1Pod(result *unstructured.Unstructured) {
	isSaAutoMount, found, err := unstructured.NestedBool(result.Object, "spec", "automountServiceAccountToken")
	if err != nil {
		return
	}
	if !found || (found && isSaAutoMount) {
		// remove default token volume
		items, found, err := unstructured.NestedSlice(result.Object, "spec", "volumes")
		if !found || err != nil {
			return
		}
		foundDefaultTokenVolume := false
		defaultTokenVolumeName := ""
		newItems := make([]interface{}, 0)
		for _, item := range items {
			volume, ok := item.(map[string]interface{})
			if !ok {
				return
			}
			vName, vFound, vErr := unstructured.NestedString(volume, "name")
			if !vFound || vErr != nil {
				return
			}
			isMatched, reErr := regexp.Match(defaultTokenVolumeNameRe, []byte(vName))
			if reErr != nil {
				return
			}
			if isMatched {
				foundDefaultTokenVolume = true
				defaultTokenVolumeName = vName
				continue
			}
			newItems = append(newItems, item)
		}
		err = unstructured.SetNestedSlice(result.Object, newItems, "spec", "volumes")
		if err != nil {
			klog.ErrorDepth(2, fmt.Sprintf("failed to trim Pod volumes %s/%s: %v", result.GetNamespace(), result.GetName(), err))
			return
		}

		if foundDefaultTokenVolume {
			newContainers := make([]interface{}, 0)
			// remove default token volume
			containerItems, found, err := unstructured.NestedSlice(result.Object, "spec", "containers")
			if !found || err != nil {
				return
			}
			for _, containerItem := range containerItems {
				container, ok := containerItem.(map[string]interface{})
				if !ok {
					return
				}
				volumeMountItems, vmFound, vmErr := unstructured.NestedSlice(container, "volumeMounts")
				if vmErr != nil {
					return
				}
				if !vmFound {
					continue
				}
				newVmItems := make([]interface{}, 0)
				for _, vmItem := range volumeMountItems {
					volumeMount, ok := vmItem.(map[string]interface{})
					if !ok {
						return
					}
					vName, vFound, vErr := unstructured.NestedString(volumeMount, "name")
					if !vFound || vErr != nil {
						return
					}
					if defaultTokenVolumeName == vName {
						continue
					}
					newVmItems = append(newVmItems, vmItem)
				}
				err = unstructured.SetNestedSlice(container, newVmItems, "volumeMounts")
				if err != nil {
					klog.ErrorDepth(2, fmt.Sprintf("failed to trim Pod %s/%s volumemMounts: %v", result.GetNamespace(), result.GetName(), err))
					return
				}
				newContainers = append(newContainers, containerItem)
			}
			err = unstructured.SetNestedSlice(result.Object, newContainers, "spec", "containers")
			if err != nil {
				klog.ErrorDepth(2, fmt.Sprintf("failed to trim Pod %s/%s volumemMounts: %v", result.GetNamespace(), result.GetName(), err))
				return
			}
		}
	}
	unstructured.RemoveNestedField(result.Object, "spec", "preemptionPolicy")
	unstructured.RemoveNestedField(result.Object, "status")
}
