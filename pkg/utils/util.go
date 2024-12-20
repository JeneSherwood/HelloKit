package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"time"
)

// GetParam 通用参数获取
func GetParam(s string, ans interface{}) (err error) {
	if ans == nil {
		return nil
	}
	err = json.Unmarshal([]byte(s), ans)
	if err != nil {
		return err
	}
	return nil
}

// bytes2string
func SliceBytesToSlice(payload []byte) []string {
	var cnt int
	var s = make([]byte, 0)
	var arr = make([]string, 0)

	for _, v := range payload {
		if v == 34 {
			cnt++
			continue
		}
		if cnt == 1 {
			s = append(s, v)
		}

		if cnt == 2 {
			cnt = 0
			arr = append(arr, string(s))
			s = make([]byte, 0)
			continue
		}
	}
	return arr
}

// GetNowTxTimestamp 获取当前交易时间戳
func GetNowTxTimestamp(stub shim.ChaincodeStubInterface) uint64 {
	ts, err := stub.GetTxTimestamp()
	if err != nil {
		fmt.Printf("Error getting transaction timestamp: %s", err)
		return 0
	}
	return uint64(time.Unix(ts.Seconds, int64(ts.Nanos)).UnixNano() / 1e6)
}
