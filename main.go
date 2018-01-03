package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func AddAutoRecover(instanceid, accountid, snstopic, region string, oversea bool) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := cloudwatch.New(sess)
	var (
		name               = "awsec2-" + instanceid + "-High-Status-Check-Failed-System-"
		statistic          = "Maximum"
		period             = aws.Int64(60)
		evaluationperiods  = aws.Int64(2)
		threshold          = 1.0
		comparisonoperator = "GreaterThanOrEqualToThreshold"
		metricname         = "StatusCheckFailed_System"
		namespace          = "AWS/EC2"
		actionOne          string
		actionTwo          string
	)
	if oversea {
		actionOne = "arn:aws:sns:" + region + ":" + accountid + ":" + snstopic
		actionTwo = "arn:aws:automate:" + region + ":ec2:recover"
	} else {
		actionOne = "arn:aws-cn:sns:" + region + ":" + accountid + ":" + snstopic
		actionTwo = "arn:aws-cn:automate:" + region + ":ec2:recover"
	}

	input := &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(name),
		Statistic:          aws.String(statistic),
		Period:             period,
		EvaluationPeriods:  evaluationperiods,
		Threshold:          aws.Float64(threshold),
		ComparisonOperator: aws.String(comparisonoperator),
		MetricName:         aws.String(metricname),
		Namespace:          aws.String(namespace),
		AlarmActions: []*string{
			aws.String(actionOne),
			aws.String(actionTwo),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceid),
			},
		},
	}
	results, err := svc.PutMetricAlarm(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			default:
				fmt.Println(err.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	if results != nil {
		fmt.Printf("Add Auto-recover Status To instance %s\n", instanceid)
	}

}
func GetInstanceId(region string) ([]string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}
	results, err := svc.DescribeInstances(input)
	id := make([]string, 0)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			default:
				fmt.Println(err.Error())
			}

		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	} else {
		for i := 0; i < len(results.Reservations); i++ {
			val := *(results.Reservations[i].Instances[0].InstanceId)
			id = append(id, val)
		}
	}
	return id, err
}

func main() {
	accountid := flag.String("accountid", "", "Aws account id")
	snstopic := flag.String("snstopic", "", "sns topic name")
	region := flag.String("region", "", "aws region name")
	oversea := flag.Bool("oversea", false, "[Optional parameters] Whether for aws overseas area")
	instanceid := flag.String("instanceid", "", "[Optional parameters] Specify a single instance")
	flag.Parse()

	if *accountid == "" || *snstopic == "" || *region == "" {
		flag.Usage()
	} else {

		idinstance, err := GetInstanceId(*region)
		if err != nil {
			fmt.Println(err)
		} else {
			if *(instanceid) == "" {
				for i := 0; i < len(idinstance); i++ {
					AddAutoRecover(idinstance[i], *accountid, *snstopic, *region, *oversea)
				}
			} else {
				AddAutoRecover(*instanceid, *accountid, *snstopic, *region, *oversea)
			}
		}
	}
}
