package main

import (
	"io/ioutil"
	"os"
	"securitycc/chaincode"
	"securitycc/pkg/logger"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	log "github.com/sirupsen/logrus"
)

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	// 日志等级默认info，格式为text。
	cfg := &logger.Config{
		Level:  logger.InfoLevel,
		Format: logger.FormatText,
	}
	logger.InitLogger(cfg)
	// See chaincode.env.example
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	// 设置日志显示行号等信息
	log.SetReportCaller(true)
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	log.SetFormatter(&log.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})

	server := &shim.ChaincodeServer{
		CCID:     config.CCID,
		Address:  config.Address,
		CC:       &chaincode.Securitycc{},
		TLSProps: getTLSProperties(),
	}

	if err := server.Start(); err != nil {
		log.Panicf("error starting sampleChaincode-go chaincode: %s", err)
	}
}

func getTLSProperties() shim.TLSProperties {
	// Check if chaincode is TLS enabled
	tlsEnabledStr := getEnvOrDefault("CHAINCODE_TLS_ENABLED", "false")
	key := getEnvOrDefault("CHAINCODE_TLS_KEY", "")
	cert := getEnvOrDefault("CHAINCODE_TLS_CERT", "")
	encKey := getEnvOrDefault("CHAINCODE_TLS_ENCKEY", "")
	encCert := getEnvOrDefault("CHAINCODE_TLS_ENCCERT", "")
	clientCACert := getEnvOrDefault("CHAINCODE_CLIENT_CA_CERT", "")

	// convert tlsDisabledStr to boolean
	tlsDisabled := !getBoolOrDefault(tlsEnabledStr, false)
	var keyBytes, certBytes, encKeyBytes, encCertBytes, clientCACertBytes []byte
	var err error

	if !tlsDisabled {
		keyBytes, err = ioutil.ReadFile(key)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
		certBytes, err = ioutil.ReadFile(cert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
		_, err = os.Stat(encKey)
		if err == nil {
			encKeyBytes, err = ioutil.ReadFile(encKey)
			if err != nil {
				log.Panicf("error while reading the crypto file: %s", err)
			}
		}
		_, err = os.Stat(encKey)
		if err == nil {
			encCertBytes, err = ioutil.ReadFile(encCert)
			if err != nil {
				log.Panicf("error while reading the crypto file: %s", err)
			}
		}
	}
	// Did not request for the peer cert verification
	if clientCACert != "" {
		clientCACertBytes, err = ioutil.ReadFile(clientCACert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}

	return shim.TLSProperties{
		Disabled:      tlsDisabled,
		Key:           keyBytes,
		Cert:          certBytes,
		EncKey:        encKeyBytes,
		EncCert:       encCertBytes,
		ClientCACerts: clientCACertBytes,
	}
}

func getEnvOrDefault(env, defaultVal string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		value = defaultVal
	}
	return value
}

// Note that the method returns default value if the string
// cannot be parsed!
func getBoolOrDefault(value string, defaultVal bool) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}
