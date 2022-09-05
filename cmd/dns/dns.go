package main

import (
	"github.com/DesistDaydream/dns-mgmt/pkg/config"
	"github.com/DesistDaydream/dns-mgmt/pkg/logging"
	"github.com/DesistDaydream/dns-mgmt/pkg/registrar/namecom"
	"github.com/DesistDaydream/dns-mgmt/pkg/registrar/tencent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type flagsConfig struct {
	Operation  string
	domReg     string
	authFile   string
	excelFile  string
	excelSheet string
}

func (flags *flagsConfig) AddFlags() {
	pflag.StringVarP(&flags.Operation, "operation", "o", "", "操作类型: [update, list]")
	pflag.StringVarP(&flags.domReg, "dom-reg", "r", "all", "域名注册商, 可选值: all, name.com, tencent")
	pflag.StringVarP(&flags.authFile, "auth-file", "f", "dd.yaml", "配置文件路径")
	pflag.StringVarP(&flags.excelFile, "excel-file", "F", "/mnt/d/Documents/WPS Cloud Files/1054253139/团队文档/东部王国/域名/102205.xyz_dns_records.xlsx", "配置文件路径")
	pflag.StringVarP(&flags.excelSheet, "excel-sheet", "S", "102205.xyz_dns_records", "配置文件路径")
}

func Run(ak string, sk string, flags *flagsConfig, domReg string) {
	switch domReg {
	case "name.com":
		c := namecom.NewClient(ak, sk)
		namecom.Run(c, flags.Operation, flags.excelFile, flags.excelSheet)
	case "tencent":
		c := tencent.NewClient(ak, sk)
		c.Run(flags.Operation)
	default:
		logrus.Fatalf("不支持的域名注册商: %v", domReg)
	}
}

func main() {
	flags := flagsConfig{}
	flags.AddFlags()
	// 添加命令行标志
	logFlags := logging.LoggingFlags{}
	logFlags.AddFlags()
	pflag.Parse()

	// 初始化日志
	if err := logging.LogInit(logFlags.LogLevel, logFlags.LogOutput, logFlags.LogFormat); err != nil {
		logrus.Fatal("set log level error")
	}

	// 读取认证信息
	auth := config.NewAuthInfo(flags.authFile)

	if flags.domReg == "all" {
		for _, r := range auth.AuthList {
			Run(r.AK, r.SK, &flags, r.Reg)
		}
	} else {
		for _, r := range auth.AuthList {
			if r.Reg == flags.domReg {
				Run(r.AK, r.SK, &flags, flags.domReg)
			} else {
				continue
			}
		}
	}
}
