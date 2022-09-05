package fileparse

import (
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// Excel 中每行的数据
type ExcelRowData struct {
	Type     string `json:"type" xlsx:"记录类型"`
	Host     string `json:"host" xlsx:"主机记录"`
	Answer   string `json:"answer" xlsx:"记录值"`
	TTL      string `json:"ttl" xlsx:"TTL值"`
	Priority string `json:"priority" xlsx:"优先级"`
}

// Excel 中的数据
type ExcelData struct {
	Rows []ExcelRowData `json:"data"`
}

func NewExcelData(file string, domainName string) (*ExcelData, error) {
	var ed ExcelData
	f, err := excelize.OpenFile(file)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logrus.Errorln(err)
			return
		}
	}()

	// 逐行读取Excel文件
	rows, err := f.GetRows(domainName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  file,
			"sheet": domainName,
		}).Errorf("读取中sheet页异常: %v", err)
		return nil, err
	}

	for k, row := range rows {
		// 跳过第一行
		if k == 0 {
			continue
		}
		logrus.WithFields(logrus.Fields{
			"k":   k,
			"row": row,
		}).Debugf("检查每一条需要处理的解析记录")

		// 将每一行中的的每列数据赋值到结构体重
		var erd ExcelRowData
		erd.Type = row[0]
		erd.Host = row[1]
		erd.Answer = row[2]
		erd.TTL = row[3]

		// 尝试第5列的值，若无法获取则设置为空
		if len(row) > 4 {
			erd.Priority = row[4]
		} else {
			erd.Priority = ""
		}

		ed.Rows = append(ed.Rows, erd)
	}

	return &ExcelData{
		Rows: ed.Rows,
	}, nil
}
