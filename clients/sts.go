package clients

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

func NewSTSClient(sdkConfig *sdk.Config, key, secret string) (*STSClient, error) {
	client, err := sts.NewClientWithOptions("us-east-1", sdkConfig, credentials.NewAccessKeyCredential(key, secret))
	if err != nil {
		return nil, err
	}
	return &STSClient{client: client}, nil
}

type STSClient struct {
	client *sts.Client
}

func (c *STSClient) AssumeRole(userName, roleARN string) (*sts.AssumeRoleResponse, error) {
	assumeRoleReq := sts.CreateAssumeRoleRequest()
	assumeRoleReq.RoleArn = roleARN
	assumeRoleReq.RoleSessionName = userName
	return c.client.AssumeRole(assumeRoleReq)
}
