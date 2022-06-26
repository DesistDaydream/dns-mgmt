package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// 认证信息配置
type AuthConfig struct {
	AuthList []AuthList `json:"authList" yaml:"authList"`
}
type AuthList struct {
	Reg string `json:"reg" yaml:"reg"`
	AK  string `json:"ak" yaml:"ak"`
	SK  string `json:"sk" yaml:"sk"`
}

func NewAuthInfo(file string) (auth *AuthConfig) {
	// 读取认证信息
	fileByte, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileByte, &auth)
	if err != nil {
		panic(err)
	}
	return auth
}
