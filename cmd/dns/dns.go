package main

import (
	"github.com/DesistDaydream/dns-mgmt/pkg/config"
	"github.com/DesistDaydream/dns-mgmt/pkg/logging"
	"github.com/DesistDaydream/dns-mgmt/pkg/registrar/namecom"
	"github.com/DesistDaydream/dns-mgmt/pkg/registrar/tencent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	operation := pflag.StringP("operation", "o", "", "操作类型: [update, list]")
	domReg := pflag.StringP("dom-reg", "r", "", "域名注册商, 可选值: name.com")
	// 添加命令行标志
	logFlags := logging.LoggingFlags{}
	logFlags.AddFlags()
	pflag.Parse()

	// 初始化日志
	if err := logging.LogInit(logFlags.LogLevel, logFlags.LogOutput, logFlags.LogFormat); err != nil {
		logrus.Fatal("set log level error")
	}

	// 读取认证信息
	auth := config.NewAuthInfo("auth.yaml")
	// 判断传入的域名是否存在在认证信息中
	if !auth.IsDomainExist(*domReg) {
		logrus.Fatalf("认证信息中不存在 %v 域名, 请检查认证信息文件或命令行参数的值", *domReg)
	}

	switch *domReg {
	case "name.com":
		namecom.Run(auth, *operation)
	case "tencent":
		tencent.Run(auth, *operation)
	default:
		logrus.Fatalf("不支持的域名注册商: %v", *domReg)
	}

}
