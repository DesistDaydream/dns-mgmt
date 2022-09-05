package namecom

import (
	"strconv"

	"github.com/DesistDaydream/dns-mgmt/pkg/fileparse"
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

func (c *Client) DeleteRecord(domain string, records []*namecom.Record) {
	for _, rr := range records {
		_, err := c.NC.DeleteRecord(&namecom.DeleteRecordRequest{
			DomainName: domain,
			ID:         rr.ID,
		})

		if err != nil {
			logrus.Errorf("删除资源记录失败：%v", err)
		} else {
			logrus.WithFields(logrus.Fields{
				"域名": rr.DomainName,
				"记录": rr.Host,
				"ID": rr.ID,
			}).Infof("删除资源记录成功")
		}
	}
}

func (c *Client) CreateRecord(domain string, data *fileparse.ExcelData) {
	for _, rr := range data.Rows {
		var (
			ttl      int64
			priority int64
		)
		ttl, err := strconv.ParseInt(rr.TTL, 10, 32)
		if err != nil {
			logrus.Errorln(err)
		}

		if rr.Priority != "" {
			priority, err = strconv.ParseInt(rr.Priority, 10, 32)
			if err != nil {
				logrus.Errorln(err)
			}
		}

		r, err := c.NC.CreateRecord(&namecom.Record{
			DomainName: domain,
			Type:       rr.Type,
			Host:       rr.Host,
			Answer:     rr.Answer,
			TTL:        uint32(ttl),
			Priority:   uint32(priority),
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"域名": domain,
				"记录": rr.Host,
				"原因": err,
			}).Errorf("创建资源记录失败")
		} else {
			logrus.WithFields(logrus.Fields{
				"ID": r.ID,
				"域名": domain,
				"记录": r.Host,
			}).Infof("创建资源记录成功")
		}
	}
}

func Run(c *Client, operation string, file string, sheet string) {
	// 获取当前域名注册商下有哪些域名
	domains, err := c.ListDomains()
	if err != nil {
		logrus.Fatalf("获取域名列表失败: %v", err)
	}

	switch operation {
	case "list":
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
					"id":   rr.ID,
					"fqdn": rr.Fqdn,
				}).Infof("%v 域名下的资源记录", d)
			}
		}
	case "update":
		// 从 Excel 文件中获取需要添加的资源记录
		data, err := fileparse.NewExcelData(file, sheet)
		if err != nil {
			logrus.Fatalf("从文件中获取域名的资源记录出错：%v", err)
		}

		for _, domain := range domains {
			records, err := c.ListRecords(domain)
			if err != nil {
				logrus.Fatal(err)
			}

			// 删除所有资源记录
			c.DeleteRecord(domain, records)
			c.CreateRecord(domain, data)
		}
	}
}
