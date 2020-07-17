# lambda_sns_event
ClowdWatch ALARM -> SNS -> LAMBDA -> Action.
- ClowdWatch ALARM configured based on CloudWatch Metrics and send event to SNS 
- SNS subscribed by Lambda func
- Labmda unmarshal string to instance &Message{}

After you can use it as you want and send any API request to remediate ALARM situation.
# aws_lambda_sns_asg_event
