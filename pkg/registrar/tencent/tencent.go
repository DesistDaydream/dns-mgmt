package tencent

import (
	"github.com/DesistDaydream/dns-mgmt/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func Run(auth *config.AuthConfig, operation string) {
	credential := common.NewCredential(
		auth.AuthList["tencent"].AK,
		auth.AuthList["tencent"].SK,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)

	switch operation {
	case "list":
		reqDomainList := dnspod.NewDescribeDomainListRequest()

		resp, err := client.DescribeDomainList(reqDomainList)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			logrus.Errorf("An API error has returned: %s", err)
			return
		}
		if err != nil {
			logrus.Fatalln(err)
		}

		for _, d := range resp.Response.DomainList {
			logrus.WithFields(logrus.Fields{
				"domain": *d.Name,
			}).Infoln("拥有域名")

			reqRecordList := dnspod.NewDescribeRecordListRequest()

			reqRecordList.Domain = common.StringPtr(*d.Name)

			resp, err := client.DescribeRecordList(reqRecordList)
			if _, ok := err.(*errors.TencentCloudSDKError); ok {
				logrus.Errorf("An API error has returned: %s", err)
				return
			}
			if err != nil {
				logrus.Fatalln(err)
			}

			for _, rr := range resp.Response.RecordList {
				logrus.WithFields(logrus.Fields{
					"fqdn": *rr.Name,
				}).Infof("%v 域名下的资源记录", *d.Name)
			}
		}
	case "update":
	}
}
