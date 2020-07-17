package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func handler(ctx context.Context, event events.SNSEvent) {
	svc := autoscaling.New(session.New())
	for _, record := range event.Records {
		fmt.Println("EventSource:", record.EventSource)
		fmt.Println("EventTimestamp:", record.SNS.Timestamp)
		// parameter pointer to record instance
		msg, err := parseEvent(&record)
		if err != nil {
			log.Println("checkEventParameters:", err)
		}

		isTrue, err := Remediation(svc, msg)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Remediation:", isTrue)
	}
}

func main() {
	lambda.Start(handler)
}

func parseEvent(record *events.SNSEventRecord) (*Message, error) {
	// Init instance
	msg := NewMessage()
	// Create byte array
	b := []byte(record.SNS.Message)

	err := msg.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Remediation func to align problem
func Remediation(svc *autoscaling.AutoScaling, msg *Message) (bool, error) {
	var isTrue bool
	if msg.Trigger.MetricName == "CPUUtilization" {
		for _, resource := range msg.Trigger.Dimensions {
			if resource.ResourceType == "AutoScalingGroupName" {
				isScale, err := scaleAsg(svc, resource.ResourceName, int64(msg.Trigger.Threshold+1))
				if err != nil {
					return false, err
				}
				isTrue = isScale
			} else {
				log.Println("Unknow resource.ResourceType")
			}
		}
	} else {
		fmt.Println("Unknown msg.Trigger.MetricName")
	}
	return isTrue, nil
}

func showEventReason(msg *Message) {
	fmt.Println("AlarmName:", msg.AlarmName)
	fmt.Println("MetricName:", msg.Trigger.MetricName)
	fmt.Printf("Threshold %v value was exceed:", msg.Trigger.Threshold)
}

func scaleAsg(svc *autoscaling.AutoScaling, asgName string, desireValue int64) (bool, error) {
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(desireValue),
		HonorCooldown:        aws.Bool(true),
	}
	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case autoscaling.ErrCodeScalingActivityInProgressFault:
				return false, errors.New(fmt.Sprintln(autoscaling.ErrCodeScalingActivityInProgressFault, aerr.Error()))
			case autoscaling.ErrCodeResourceContentionFault:
				return false, errors.New(fmt.Sprintln(autoscaling.ErrCodeResourceContentionFault, aerr.Error()))
			default:
				return false, errors.New(aerr.Error())
			}
		} else {
			return false, errors.New(aerr.Error())
		}
	}
	return true, nil
}
