package core

import v1 "k8s.io/api/core/v1"

func (in *Pod) Ready() bool {
	for _, condition := range in.Status.Conditions {
		if condition.Type == v1.PodReady {
			if condition.Status == v1.ConditionTrue {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
