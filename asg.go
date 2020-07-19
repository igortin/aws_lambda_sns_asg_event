package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func getAsgParameters(svc *autoscaling.AutoScaling, asgName string) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(int64(1)),
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}

	output, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func parseToAsgInstance(data *autoscaling.DescribeAutoScalingGroupsOutput) (*AutoScaleGroup, error) {
	var asg AutoScaleGroup

	if len(data.AutoScalingGroups) == 0 {
		return nil, errors.New("autoscaling.DescribeAutoScalingGroupsOutput array  length = 0, should be 1")
	}

	if len(data.AutoScalingGroups) > 1 {
		return nil, errors.New("autoscaling.DescribeAutoScalingGroupsOutput array length > 2, should be 1")
	}

	// iterate []*Group
	for _, group := range data.AutoScalingGroups {

		// init instance of structure AutoScaleGroup{}
		asg = AutoScaleGroup{
			AutoScaleGroupName:   *group.AutoScalingGroupName,
			LaunchTemplate:       *group.LaunchTemplate,
			TargetGroupARNs:      group.TargetGroupARNs,
			MinSize:              *group.MinSize,
			MaxSize:              *group.MaxSize,
			ServiceLinkedRoleARN: *group.ServiceLinkedRoleARN,
			DesiredSize:          *group.DesiredCapacity,
			Instances:            group.Instances,
		}
	}

	return &asg, nil
}
