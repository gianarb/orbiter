package provider

import (
	"fmt"
	"log"

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
		runResult, err := p.client.RunInstances(&ec2.RunInstancesInput{
			// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
			ImageId:        aws.String("ami-e7527ed7"),
			InstanceType:   aws.String("t2.micro"),
			MinCount:       aws.Int64(1),
			MaxCount:       aws.Int64(1),
			UserData:       "",
			SecurityGroups: []*string{},
			SubnetId:       "",
		})

		if err != nil {
			log.Println("Could not create instance", err)
			return
		}

		log.Println("Created instance", *runResult.Instances[0].InstanceId)
	} else {
		params := &ec2.StopInstancesInput{
			InstanceIds: []*string{ // Required
				aws.String("String"), // Required
				// More values...
			},
			DryRun: aws.Bool(true),
			Force:  aws.Bool(true),
		}
		resp, err := p.client.StopInstances(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
		}

	}

}
func (p AwsEc2Provider) isGoodToBeDeleted(droplet godo.Droplet, serviceId string) bool {
}
