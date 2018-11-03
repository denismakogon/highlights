package main

type Request struct {
	StatUUID string `json:"stat_uuid"`
}

type Threshold struct {
	ThresholdValue float64 `json:"threshold_value"`
	StatUUID       string  `json:"stat_uuid"`
}
