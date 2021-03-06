package endpoints

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2ClientMock struct {
	ec2iface.EC2API

	isError          bool
	privateIPAddress string
}

func (e EC2ClientMock) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if e.isError {
		return nil, fmt.Errorf("error!!")
	}

	output := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						PrivateIpAddress: aws.String(e.privateIPAddress),
						State: &ec2.InstanceState{
							Code: aws.Int64(ec2RunningState),
						},
					},
				},
			},
		},
	}

	return output, nil
}
