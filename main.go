package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.SNSEvent) {
	for _, record := range event.Records {
		snsRecord := record.SNS

		fmt.Printf("[%s %s]\n", record.EventSource, snsRecord.Timestamp)

		msg := &Message{}
		err := json.Unmarshal([]byte(record.SNS.Message), msg)
		if err != nil {
			log.Println("Check structure Message:", err)
		}
		fmt.Println("Success:", reflect.TypeOf(msg))
		fmt.Println("AlarmName:", msg.AlarmName)
		fmt.Println("Trigger MetricName:", msg.Trigger.MetricName)

		for _, item := range msg.Trigger.Dimensions {
			fmt.Println("Success:", reflect.TypeOf(item))
			fmt.Println("Success:", item.ResourceType)
			fmt.Println("Success:", item.ResourceName)
		}
	}
}

func main() {
	lambda.Start(handler)
}
