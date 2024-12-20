package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"securitycc/pkg/constant"
)

func WriteToChain(stub shim.ChaincodeStubInterface, key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Printf("Write to chain key: %s, value: %s\n", key, string(b))
	return stub.PutState(key, b)
}

func DeleteByKey(stub shim.ChaincodeStubInterface, key string) error {
	err := stub.DelState(key)
	if err != nil {

		return fmt.Errorf("delete state err: %s", err)
	}
	return nil
}

func GetStateByKey(stub shim.ChaincodeStubInterface, key string, data interface{}) error {
	bys, err := stub.GetState(key)
	if err != nil {
		return fmt.Errorf("get state err: %s", err)
	}
	if bys == nil {
		return errors.New(constant.StateEmpty)
	}
	err = json.Unmarshal(bys, data)
	if err != nil {
		return fmt.Errorf("unmarshal err: %s", err)
	}
	return nil
}
