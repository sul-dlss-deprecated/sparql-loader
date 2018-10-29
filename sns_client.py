import boto3


class SnsClient():
    def __init__(self, sns_endpoint, topic_arn, aws_region):
        self.topic_arn = topic_arn
        self.sns_conn = boto3.client('sns',
                                     region_name=aws_region,
                                     endpoint_url=sns_endpoint)

    def publish(self, message):
        response = self.sns_conn.publish(
            TopicArn=self.topic_arn,
            Message=message
        )
        return response
