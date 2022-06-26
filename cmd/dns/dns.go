package main

import (
	"github.com/DesistDaydream/dns-mgmt/pkg/config"
	"github.com/DesistDaydream/dns-mgmt/pkg/logging"
	"github.com/DesistDaydream/dns-mgmt/pkg/registrar/namecom"
	"github.com/DesistDaydream/dns-mgmt/pkg/registrar/tencent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func Run(ak string, sk string, operation string, domReg string) {
	switch domReg {
	case "name.com":
		c := namecom.NewClient(ak, sk)
		c.Run(operation)
	case "tencent":
		c := tencent.NewClient(ak, sk)
		c.Run(operation)
	default:
		logrus.Fatalf("不支持的域名注册商: %v", domReg)
	}
}

func main() {
	operation := pflag.StringP("operation", "o", "", "操作类型: [update, list]")
	domReg := pflag.StringP("dom-reg", "r", "all", "域名注册商, 可选值: all, name.com, tencent")
	authFile := pflag.StringP("auth-file", "f", "dd.yaml", "配置文件路径")
	// 添加命令行标志
	logFlags := logging.LoggingFlags{}
	logFlags.AddFlags()
	pflag.Parse()

	// 初始化日志
	if err := logging.LogInit(logFlags.LogLevel, logFlags.LogOutput, logFlags.LogFormat); err != nil {
		logrus.Fatal("set log level error")
	}

	// 读取认证信息
	auth := config.NewAuthInfo(*authFile)

	if *domReg == "all" {
		for _, r := range auth.AuthList {
			Run(r.AK, r.SK, *operation, r.Reg)
		}
	} else {
		for _, r := range auth.AuthList {
			if r.Reg == *domReg {
				Run(r.AK, r.SK, *operation, *domReg)
			} else {
				continue
			}
		}
	}
}
