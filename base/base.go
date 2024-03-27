package base

import (
	"encoding/json"
	"fmt"
)

const (
	ContainersReady string = "ContainersReady"
	PodInitialized  string = "Initialized"
	PodReady        string = "Ready"
	PodScheduled    string = "PodScheduled"
)

const (
	ConditionTrue    string = "True"
	ConditionFalse   string = "False"
	ConditionUnknown string = "Unknown"
)

const (
	ConditionUnavailable  string = "Unavailable"
	ConditionInitializing string = "Initializing"
	ConditionScheduling   string = "Scheduling"
	ConditionRunning      string = "Running"
)

type ConditionDetail struct {
	Type            string `json:"type"`
	Reason          string `json:"reason"`
	Message         string `json:"message"`
	Status          string `json:"status"`
	ConditionStatus string `json:"condition_status"`
}

type JobPodsStatus struct {
	Name       string             `json:"name"`
	Conditions []*ConditionDetail `json:"conditions"`
}

func ConvertConditionDetailStructToString(conditionDetails []*ConditionDetail) (string, error) {
	str, err := json.Marshal(conditionDetails)
	if err != nil {
		return "", fmt.Errorf("marshal condition details error: %s", err.Error())
	}
	return string(str), nil
}

func ConvertConditionDetailStringToStruct(str string) ([]*ConditionDetail, error) {
	out := make([]*ConditionDetail, 0)
	err := json.Unmarshal([]byte(str), &out)
	if err != nil {
		return nil, fmt.Errorf("unnamrshal str to condition details error: %s", err.Error())
	}
	return out, nil
}
