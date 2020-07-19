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
	var asgName string
	svc := autoscaling.New(session.New())

	for _, record := range event.Records {
		fmt.Println("EventSource:", record.EventSource)
		fmt.Println("EventTimestamp:", record.SNS.Timestamp)
		// parameter pointer to record instance
		msg, err := parseEvent(&record)
		if err != nil {
			log.Println("checkEventParameters:", err)
			return
		}

		// Get resource.ResourceName based on msg dimension
		for _, resource := range msg.Trigger.Dimensions {
			if resource.ResourceType == "AutoScalingGroupName" {
				asgName = resource.ResourceName
				fmt.Println("ASG name:", asgName)
			}
		}

		output, err := getAsgParameters(svc, asgName)
		if err != nil {
			fmt.Println("getAsgParameters:", err)
			return
		}

		asgroup, err := parseToAsgInstance(output)
		if err != nil {
			fmt.Println("parseToInstance:", err)
			return
		}

		desireValue, err := getDesireValue(asgroup, msg.Trigger.ComparisonOperator)
		if err != nil {
			fmt.Println("getDesireValue:", err)
			return
		}
		_, err = remediationEvent(svc, asgroup.AutoScaleGroupName, desireValue)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	lambda.Start(handler)
}

func parseEvent(record *events.SNSEventRecord) (*Message, error) {
	msg := NewMessage()
	b := []byte(record.SNS.Message)
	err := msg.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	fmt.Println(record.SNS.Message)
	return msg, nil
}

func remediationEvent(svc *autoscaling.AutoScaling, resourceName string, desireValue int64) (bool, error) {
	// one of possible remediation

	isScale, err := scaleAsg(svc, resourceName, desireValue)
	if err != nil {
		return false, err
	}
	fmt.Println("Remediation:", isScale)
	return isScale, nil
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

func getDesireValue(asgroup *AutoScaleGroup, ComparisonOperator string) (int64, error) {
	var count int64
	var step int64

	switch ComparisonOperator {
	case "GreaterThanOrEqualToThreshold":
		step = 1
	case "GreaterThanThreshold":
		step = 1
	case "LessThanOrEqualToThreshold":
		step = -1
	case "LessThanThreshold":
		step = -1
	}

	fmt.Println("Previos ASG instances count:", len(asgroup.Instances))

	for _, item := range asgroup.Instances {
		fmt.Printf("instance ID: %s, health status: %s\n", *item.InstanceId, *item.HealthStatus)
		if *item.LifecycleState == "InService" {
			count++
		}
	}

	// Scale out
	if count < asgroup.MaxSize && step == 1 {
		fmt.Println("getDesireValue new value:", count+step)
		return count + step, nil
	}

	// Scale in
	if count > asgroup.MinSize && step == -1 {
		fmt.Println("getDesireValue new value:", count+step)
		return count + step, nil
	}
	fmt.Println("getDesireValue old value:", count)
	return count, errors.New("Error: Current desireValue could not be decresed on increased")
}
