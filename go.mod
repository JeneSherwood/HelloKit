module securitycc

go 1.15

require (
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20220131132609-1476cf1d3206
	github.com/hyperledger/fabric-protos-go v0.0.0-20220315113721-7dc293e117f7
	github.com/kr/pretty v0.3.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/grpc v1.46.2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace (
	github.com/Hyperledger-TWGC/ccs-gm v0.1.1 => 192.168.8.1/hyperledger/ccs-gm v1.0.0-alpha1-1-yx
	github.com/hyperledger/fabric-chaincode-go => 192.168.8.1/hyperledger/fabric-chaincode-go v0.0.0-20211124072959-8deb0be6daa5
	google.golang.org/grpc => 192.168.8.1/hyperledger/grpc v1.29.1-alpha1-1-yx
)
