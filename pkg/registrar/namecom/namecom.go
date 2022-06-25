package namecom

import (
	"github.com/DesistDaydream/dns-mgmt/pkg/config"
	"github.com/namedotcom/go/namecom"
	"github.com/sirupsen/logrus"
)

func Run(auth *config.AuthConfig, operation string) {
	nc := namecom.New(auth.AuthList["name.com"].AK, auth.AuthList["name.com"].SK)

	switch operation {
	case "list":
		// 获取当前域名注册商下有哪些域名
		domains, err := nc.ListDomains(&namecom.ListDomainsRequest{})
		if err != nil {
			logrus.Fatalf("获取域名列表失败: %v", err)
		}
		for _, d := range domains.Domains {
			logrus.WithFields(logrus.Fields{
				"domain": d.DomainName,
			}).Infoln("拥有域名")

			resp, err := nc.ListRecords(&namecom.ListRecordsRequest{DomainName: d.DomainName})
			if err != nil {
				logrus.Fatal(err)
			}

			for _, rr := range resp.Records {
				logrus.WithFields(logrus.Fields{
					"fqdn": rr.Fqdn,
				}).Infof("%v 域名下的资源记录", d.DomainName)
			}
		}
	case "update":
		nc.UpdateRecord(&namecom.Record{
			ID:         0,
			DomainName: "",
			Host:       "",
			Fqdn:       "",
			Type:       "",
			Answer:     "",
			TTL:        0,
			Priority:   0,
		})
	}
}
