package tccp

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned/fake"
	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_Resource_Cloudformation_GetCloudFormationTags(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description  string
		installation string
		obj          v1alpha1.AWSConfig
		expectedTags []*awscloudformation.Tag
	}{
		{
			description:  "basic match",
			installation: "test-install",
			obj: v1alpha1.AWSConfig{
				Spec: v1alpha1.AWSConfigSpec{
					Cluster: v1alpha1.Cluster{
						ID: "5xchu",
						Customer: v1alpha1.ClusterCustomer{
							ID: "giantswarm",
						},
					},
				},
			},
			expectedTags: []*awscloudformation.Tag{
				{
					Key:   aws.String("kubernetes.io/cluster/5xchu"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("giantswarm.io/cluster"),
					Value: aws.String("5xchu"),
				},
				{
					Key:   aws.String("giantswarm.io/organization"),
					Value: aws.String("giantswarm"),
				},
				{
					Key:   aws.String("giantswarm.io/installation"),
					Value: aws.String("test-install"),
				},
			},
		},
	}

	c := Config{}

	c.EncrypterBackend = "kms"
	c.G8sClient = fake.NewSimpleClientset()
	c.Logger = microloggertest.New()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c.InstallationName = tc.installation
			r, err := New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			tags := r.getCloudFormationTags(tc.obj)

			if len(tags) != len(tc.expectedTags) {
				t.Fatalf("Expected %d tags, found %d", len(tc.expectedTags), len(tags))
			}

			for _, tag := range tc.expectedTags {
				if !containsTag(tag, tags) {
					t.Fatalf("Expected cloud formation contains tag %v in the slice %v", tag, tags)
				}
			}
		})
	}
}

func containsTag(tag *awscloudformation.Tag, tags []*awscloudformation.Tag) bool {
	for _, inTag := range tags {
		if reflect.DeepEqual(tag, inTag) {
			return true
		}
	}

	return false
}
