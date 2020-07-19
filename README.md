# lambda_sns_event
ClowdWatch ALARM -> SNS -> LAMBDA -> Action.

- ClowdWatch ALARM configured based on CloudWatch Metrics and send event to SNS 
- SNS subscribed by Lambda func
- Labmda unmarshal string to instance &Message{}

After you can use it as you want and send any API request to remediate ALARM case.
Example:
> Alarm:
> ResourceType: AutoScalingGroupName  
> MetricName:  CPUUtilization
> Thershold: 10%

> SNS:
> Suscriber: Lambda  

> Lambda:  
> Action: Scale out up to currentAsgSize + 1  

