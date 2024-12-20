package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	log "github.com/sirupsen/logrus"
	"securitycc/internal/models"
	"securitycc/pkg/constant"
	"securitycc/pkg/logger"
	"securitycc/pkg/utils"
	"strconv"
)

type Securitycc struct{}

func (s *Securitycc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	log.Infof("--------------- Securitycc Init ---------------\n")

	_, param := stub.GetFunctionAndParameters()
	var level = logger.InfoLevel
	if len(param) != 0 {
		i, err := strconv.Atoi(param[0])
		if err != nil {
			log.Errorf("logger level should be int")
			return shim.Error(err.Error())
		}

		if i == 0 {
			level = logger.DebugLevel
		} else if i == 1 {
			level = logger.InfoLevel
		} else if i == 2 {
			level = logger.WarnLevel
		} else if i == 3 {
			level = logger.ErrorLevel
		} else {
			level = logger.InfoLevel
		}
	}
	var format = logger.FormatText
	logConfig := &logger.Config{
		Level:  level,
		Format: format,
	}
	log.Infof("初始化参数: %+v\n", *logConfig)
	logger.InitLogger(logConfig)

	cb, _ := json.Marshal(logConfig)
	// 日志信息上链
	err := stub.PutState(logger.LevelKey, cb)
	if err != nil {
		return shim.Error(err.Error())
	}

	log.Infof("--------------- example account Init ---------------\n")
	if len(param) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5, must be init example account A & B like [\"logLevel\",\"accA\",\"balanceA\",\"accB\",\"balanceB\"]")
	}
	log.Infof("--------------- all params ----------------\n", param)
	var A, B string    // account
	var Aval, Bval int // balance
	// Initialize the chaincode
	A = param[1]
	Aval, err = strconv.Atoi(param[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = param[3]
	Bval, err = strconv.Atoi(param[4])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	log.Infof("example account A = %s ,balance = %d, example account B = %s , balance = %d\n", A, Aval, B, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *Securitycc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	log.Infof("function:[%s], param:[%s]\n", function, args)

	switch function {
	case "Init":
		return s.Init(stub)
	case constant.CreateData:
		return CreateData(stub, args[0])
	case constant.ReceiveData:
		return ReceiveData(stub, args[0])
	case constant.LockData:
		return LockData(stub, args[0])
	case constant.Transfer:
		return Transfer(stub)
	case constant.Delete:
		return Delete(stub)
	case constant.Query:
		return Query(stub)
	}
	return shim.Error("Unsupported operation")
}

func CreateData(stub shim.ChaincodeStubInterface, param string) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	log.Infof("CreateData  function:[%s], param:[%s]\n", function, args)
	req := new(models.Security)
	err := json.Unmarshal([]byte(param), req)
	if err != nil {
		log.Errorf("Fail to unmarshal to Request, err: %s", err)
		return shim.Error(err.Error())
	}

	eventPayload := []string{
		req.ID,
		req.Content,
	}

	// 序列化拼装后的string数组
	b, err := json.Marshal(eventPayload)
	if err != nil {
		log.Errorf("Fail to marshal fund to payload, err: %s", err)
		return shim.Error(err.Error())
	}
	log.Infof("set resp: %s, param: %+v", constant.CreateData, b)
	return shim.Success(b)
}

func ReceiveData(stub shim.ChaincodeStubInterface, param string) pb.Response {
	log.Infof("receive data param: %s", param)

	params := utils.SliceBytesToSlice([]byte(param))

	req := &models.Security{
		ID:         params[0],
		Content:    params[1],
		CreateTime: utils.GetNowTxTimestamp(stub),
	}

	var key = fmt.Sprintf("%s-%s", constant.SecurityPrefix, req.ID)

	b, _ := json.Marshal(req)

	err := stub.PutState(key, b)
	if err != nil {
		log.Errorf("Fail to write to chain, key: %s, data:%s", key, b)
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func LockData(stub shim.ChaincodeStubInterface, param string) pb.Response {
	log.Infof("lock data param: %s", param)

	params := utils.SliceBytesToSlice([]byte(param))

	req := &models.Security{
		ID:         params[0],
		Content:    params[1],
		CreateTime: utils.GetNowTxTimestamp(stub),
	}

	var key = fmt.Sprintf("%s-%s", constant.SecurityPrefix, req.ID)

	b, _ := json.Marshal(req)

	err := stub.PutState(key, b)
	if err != nil {
		log.Errorf("Fail to write to chain, key: %s, data:%s", key, b)
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

// Transfer 转账示例-转账方法
func Transfer(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	log.Infof("transfer param: %s\n", args)

	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4, like:[\"func\",\"accA\",\"accB\",\"amount\"]")
	}

	A = args[1]
	B = args[2]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Delete 转账示例-删除账户方法
func Delete(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	log.Infof("delete param: %s\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, like [\"func\", \"acc\"]")
	}

	A := args[1]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// Query 转账示例-账户查询方法
func Query(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	log.Infof("query param: %s\n", args)

	var A string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, like [\"func\", \"acc\"] ")
	}

	A = args[1]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}
