package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	log "github.com/sirupsen/logrus"
	"securitycc/chaincode"
)

func main_internal() {
	// 设置日志显示行号等信息
	log.SetReportCaller(true)
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	log.SetFormatter(&log.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})

	err := shim.Start(new(chaincode.Securitycc))
	if err != nil {
		log.Errorf("Fail to start chaincode: %s\n", err)
	}
}
