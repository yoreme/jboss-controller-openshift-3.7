/*
Copyright 2017 The Kubernetes Authors.

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

package controller

import (
	"fmt"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

func filterAndGroupClaimsByOrdinal(allClaims []*corev1.PersistentVolumeClaim, sts *appsv1.StatefulSet) map[int][]*corev1.PersistentVolumeClaim {
	claims := map[int][]*corev1.PersistentVolumeClaim{}
	for _, pvc := range allClaims {
		if pvc.DeletionTimestamp != nil {
			glog.Infof("PVC '%s' is being deleted. Ignoring it.", pvc.Name)
			continue
		}

		name, ordinal, err := extractNameAndOrdinal(pvc.Name)
		if err != nil {
			continue
		}

		for _, t := range sts.Spec.VolumeClaimTemplates {
			if name == fmt.Sprintf("%s-%s", t.Name, sts.Name) {
				if claims[ordinal] == nil {
					claims[ordinal] = []*corev1.PersistentVolumeClaim{}
				}
				claims[ordinal] = append(claims[ordinal], pvc)
			}
		}
	}
	return claims
}

func getPVCName(sts *appsv1.StatefulSet, volumeClaimName string, ordinal int) string {
	return fmt.Sprintf("%s-%s-%d", volumeClaimName, sts.Name, ordinal)
}

func extractNameAndOrdinal(pvcName string) (string, int, error) {
	idx := strings.LastIndexAny(pvcName, "-")
	if idx == -1 {
		return "", 0, fmt.Errorf("PVC does not belong to a StatefulSet")
	}

	name := pvcName[:idx]
	ordinal, err := strconv.Atoi(pvcName[idx+1:])
	if err != nil {
		return "", 0, fmt.Errorf("PVC does not belong to a StatefulSet")
	}
	return name, ordinal, nil
}
