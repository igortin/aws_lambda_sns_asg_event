package main

import (
	"encoding/json"
)

// Message custom structure for unmarshal record.SNS.Message
type Message struct {
	AlarmName       string      `json:"AlarmName"`
	NewStateValue   string      `jon:"NewStateValue"`
	OldStateValue   string      `json:"OldStateValue"`
	StateChangeTime string      `json:"StateChangeTime"`
	Trigger         *NewTrigger `json:"Trigger"`
}

/*
NewTrigger structure for event */
type NewTrigger struct {
	MetricName         string          `json:"MetricName"`
	Namespace          string          `json:"Namespace"`
	Statistic          string          `json:"Statistic"`
	ComparisonOperator string          `json:"ComparisonOperator"`
	Threshold          float64         `json:"Threshold"`
	Dimensions         []*NewDimension `json:"Dimensions"`
	EvaluationPeriods  float64         `json:"EvaluationPeriods"`
	Period             float64         `json:"Period"`
}

// NewDimension is object structure
type NewDimension struct {
	ResourceName string `json:"value"`
	ResourceType string `json:"name"`
}

// NewMessage return new message instance
func NewMessage() *Message {
	return &Message{}
}

// Unmarshal method fullfil instance
func (msg *Message) Unmarshal(b []byte) error {
	err := json.Unmarshal(b, msg)
	if err != nil {
		return err
	}
	return nil
}
