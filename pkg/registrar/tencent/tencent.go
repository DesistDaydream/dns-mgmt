package tencent

import (
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type Client struct {
	Client *dnspod.Client
}

func NewClient(ak string, sk string) *Client {
	credential := common.NewCredential(
		ak,
		sk,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)

	return &Client{
		Client: client,
	}
}

func (c *Client) ListDomains() ([]string, error) {
	reqDomainList := dnspod.NewDescribeDomainListRequest()
	resp, err := c.Client.DescribeDomainList(reqDomainList)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logrus.Errorf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		logrus.Fatalln(err)
	}

	var domainList []string
	for _, d := range resp.Response.DomainList {
		domainList = append(domainList, *d.Name)
	}
	return domainList, nil
}

func (c *Client) ListRecords(domain string) ([]*dnspod.RecordListItem, error) {
	reqRecordList := dnspod.NewDescribeRecordListRequest()

	reqRecordList.Domain = common.StringPtr(domain)

	resp, err := c.Client.DescribeRecordList(reqRecordList)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logrus.Errorf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		logrus.Fatalln(err)
	}

	return resp.Response.RecordList, nil
}

func Run(c *Client, operation string) {

	switch operation {
	case "list":
		domains, err := c.ListDomains()
		if err != nil {
			logrus.Fatalf("获取域名列表失败: %v", err)
		}

		for _, d := range domains {
			logrus.WithFields(logrus.Fields{
				"domain": d,
			}).Infoln("DNSPod 中拥有的域名")

			records, err := c.ListRecords(d)
			if err != nil {
				logrus.Fatalln(err)
			}

			for _, rr := range records {
				logrus.WithFields(logrus.Fields{
					"fqdn": *rr.Name,
				}).Infof("%v 域名下的资源记录", d)
			}

		}
	case "update":
		// TODO: 更新域名记录
	}
}
