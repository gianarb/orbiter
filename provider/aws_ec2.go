package provider

import (
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitalocean/godo"
	"github.com/gianarb/orbiter/autoscaler"
)

type AwsEc2Provider struct {
	client *ec2.EC2
	config map[string]string
}

func NewAwsEc2Provider(c map[string]string) (autoscaler.Provider, error) {
	var p AwsEc2Provider
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", ""),
	})
	if err != nil {
		return p, err
	}
	svc := ec2.New(sess)
	p = AwsEc2Provider{
		client: svc,
		config: c,
	}
	return p, nil
}

func (p AwsEc2Provider) Scale(serviceId string, target int, direction bool) error {
	if direction == true {
		// Specify the details of the instance that you want to create.
		_, err := p.client.RunInstances(&ec2.RunInstancesInput{
			// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
			ImageId:      aws.String("ami-e7527ed7"),
			InstanceType: aws.String("t2.micro"),
			MinCount:     aws.Int64(2),
			MaxCount:     aws.Int64(2),
			//UserData:       "",
			SecurityGroups: []*string{},
			//SubnetId:       "",
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"provider": "aws_ec2",
				"error":    err,
			}).Warn("We are not able to create new instances")
			return err
		}
	} else {
		params := &ec2.StopInstancesInput{
			InstanceIds: []*string{ // Required
				aws.String("String"), // Required
				// More values...
			},
			DryRun: aws.Bool(true),
			Force:  aws.Bool(true),
		}
		_, err := p.client.StopInstances(params)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"provider": "aws_ec2",
				"error":    err,
			}).Warn("We are not able to delete the instances")
			return err
		}
	}
	return nil
}
func (p AwsEc2Provider) isGoodToBeDeleted(droplet godo.Droplet, serviceId string) bool {
	return true
}
