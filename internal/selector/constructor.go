package selector

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"math"
	"securitycc/pkg/constant"
)

// Constructor couchdb查询构造器
type Constructor struct {
	Selector map[string]interface{} `json:"selector"`
	Sort     []map[string]string    `json:"sort,omitempty"`
	Skip     int                    `json:"skip,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Bookmark string                 `json:"bookmark"`
}

// NewIdConstructor 获取正则id的正则selector
func NewIdConstructor(id string) *Constructor {
	constructor := &Constructor{
		Selector: map[string]interface{}{
			"docType": constant.TtmTable,
		},
	}
	constructor.Selector["_id"] = map[string]string{
		"$regex": fmt.Sprintf("%s.*", id),
	}
	return constructor
}

func (s *Constructor) InjectPageNoSize(pageNo, pageSize int) {
	// 如果有一个为小于等于 0，不分页
	if pageNo <= 0 || pageSize <= 0 {
		s.Limit = -1
		s.Skip = 0
		return
	}
	s.Skip = (pageNo - 1) * pageSize
	s.Limit = pageSize
}

func (s *Constructor) GetQuery(stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	selectorBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return stub.GetQueryResult(string(selectorBytes))
}

func (s *Constructor) GetCouchdbIter(stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	selectorBytes, err := json.Marshal(s)
	if err != nil {
		return nil, nil, err
	}

	return stub.GetQueryResultWithPagination(string(selectorBytes), int32(s.Limit), s.Bookmark)
}

func (s *Constructor) GetCount(stub shim.ChaincodeStubInterface) (int32, error) {
	const pageSize = constant.PageSize
	var (
		respMeta     *pb.QueryResponseMetadata
		pageNum      = int32(0)
		prePage      int
		lastPage     = 1
		lastBookmark string
	)

	for respMeta == nil || respMeta.FetchedRecordsCount == pageSize || respMeta.FetchedRecordsCount == 0 {
		s.Skip = (lastPage - 1) * int(pageSize)
		selectorByte, err := json.Marshal(s)
		if err != nil {
			return 0, err
		}
		_, respMeta, err = stub.GetQueryResultWithPagination(string(selectorByte), pageSize, "")
		if err != nil {
			return 0, err
		}
		if respMeta == nil {
			break
		}

		if respMeta.FetchedRecordsCount == 0 { //开始往回查找
			tmp, err := BinarySearch(prePage, lastPage, int(pageSize), *s, stub)
			if err != nil {
				return 0, err
			}
			return int32(tmp), nil
		} else { //继续倍增查找
			pageNum++
			if lastBookmark == respMeta.Bookmark {
				respMeta.FetchedRecordsCount = pageSize
				break
			}
			lastBookmark = respMeta.Bookmark
			prePage = lastPage
			lastPage = int(math.Pow(float64(2), float64(pageNum)))
		}
	}
	if respMeta == nil {
		return int32(prePage-1) * pageSize, nil
	}
	return int32(prePage-1)*pageSize + respMeta.FetchedRecordsCount, nil
}

// 二分往回查找
func BinarySearch(prePage, lastPage, pageSize int, selector Constructor, stub shim.ChaincodeStubInterface) (int, error) {
	var mid int
	var respMeta *pb.QueryResponseMetadata
	var err error
	var selectorByte []byte
	for prePage < lastPage {
		mid = int(math.Ceil(float64(prePage+lastPage) / 2.0))
		selector.Skip = (mid - 1) * pageSize
		selectorByte, err = json.Marshal(selector)
		if err != nil {
			return 0, err
		}
		_, respMeta, err = stub.GetQueryResultWithPagination(string(selectorByte), int32(pageSize), "")
		if err != nil {
			return 0, err
		}
		if respMeta == nil || mid == lastPage {
			break
		}
		if respMeta.FetchedRecordsCount == int32(pageSize) {
			prePage = mid
		} else {
			lastPage = mid
		}
	}
	return (mid-1)*pageSize + int(respMeta.FetchedRecordsCount), nil
}
