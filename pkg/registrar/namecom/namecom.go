package namecom

import (
	"github.com/namedotcom/go/namecom"
	"github.com/sirupsen/logrus"
)

type Client struct {
	NC *namecom.NameCom
}

func NewClient(ak string, sk string) *Client {
	return &Client{
		NC: namecom.New(ak, sk),
	}
}

func (c *Client) ListDomains() ([]string, error) {
	domains, err := c.NC.ListDomains(&namecom.ListDomainsRequest{})
	if err != nil {
		return nil, err
	}

	var domainList []string
	for _, d := range domains.Domains {
		domainList = append(domainList, d.DomainName)
	}
	return domainList, nil
}

func (c *Client) ListRecords(domain string) ([]*namecom.Record, error) {
	records, err := c.NC.ListRecords(&namecom.ListRecordsRequest{DomainName: domain})
	if err != nil {
		logrus.Fatalf("获取域名记录失败: %v", err)
	}
	return records.Records, nil
}

func (c *Client) UpdateRecord(domain string, record namecom.Record) error {
	_, err := c.NC.UpdateRecord(&namecom.Record{
		ID:         0,
		DomainName: "",
		Host:       "",
		Fqdn:       "",
		Type:       "",
		Answer:     "",
		TTL:        0,
		Priority:   0,
	})
	if err != nil {
		logrus.Fatalf("更新域名记录失败: %v", err)
	}

	return nil
}

func (c *Client) Run(operation string) {
	switch operation {
	case "list":
		// 获取当前域名注册商下有哪些域名
		domains, err := c.ListDomains()
		if err != nil {
			logrus.Fatalf("获取域名列表失败: %v", err)
		}

		// 获取每个域名的所有资源记录
		for _, d := range domains {
			logrus.WithFields(logrus.Fields{
				"domain": d,
			}).Infoln("Name.com 中拥有的域名")

			records, err := c.ListRecords(d)
			if err != nil {
				logrus.Fatal(err)
			}

			for _, rr := range records {
				logrus.WithFields(logrus.Fields{
					"fqdn": rr.Fqdn,
				}).Infof("%v 域名下的资源记录", d)
			}
		}
	case "update":
		// TODO: 更新域名记录
	}
}
