package code

import (
	"math/big"

	"github.com/xuperchain/contract-sdk-go/pb"
)

// Context is the context in which the contract runs
type Context interface {
	// Args return call args of a contract call
	Args() map[string][]byte

	//  Caller return caller of an contract call
	//  it returns direct caller in in-contract-call, initiator otherwise
	Caller() string

	//  Initiator return the  initiator of a contract call
	Initiator() string

	// AuthRequire return auth address required to verify a contract call
	// it returns more than one  addresss in multi sign situation
	AuthRequire() []string

	// PutObject put a object into state, all the value is  transaction isolated,you can read the KV pair
	// in the same transaction, you can visit it through NewIterator,
	// but you can not read it in other transaction until the transaction is commited
	// error may be no nil, a typical case is networking error
	PutObject(key []byte, value []byte) error

	// GetObject return a object, see PutObject for more details
	GetObject(key []byte) ([]byte, error)

	// DeleteObject delete a object, see PutObject for more details
	DeleteObject(key []byte) error

	// NewIterator return a range iterator
	// you can use PrefixRange to generate start and limit if you wanna iterate over a specify prefix
	NewIterator(start, limit []byte) Iterator

	// QueryTx return  a transaction by txid
	QueryTx(txid string) (*pb.Transaction, error)

	// QueryBlock return a block by blockid
	QueryBlock(blockid string) (*pb.Block, error)

	// Transfer transfer an amount to a address
	Transfer(to string, amount *big.Int) error

	// TransferAmount get transfer amount of an ContractCall
	TransferAmount() (*big.Int, error)

	// Call start a in-contract call,module can be wasm|native|evm, dependding on contract you call
	Call(module, contract, method string, args map[string][]byte) (*Response, error)
	// CrossQueryStart a cross chain query, see https://xuper.baidu.com/n/xuperdoc/ CrossQuery section for more information.
	CrossQuery(uri string, args map[string][]byte) (*Response, error)

	// EmitEvent emit a event with name
	EmitEvent(name string, body []byte) error

	//  EmmitJSONEvent emit an event, you can use any JSON marshalable body,  body will be marshaled using encoding/json
	//  you can also use a object having a custom JSON marshaler
	EmitJSONEvent(name string, body interface{}) error

	// Logf emmit a log entry, see fmt.Sprintf for log format
	Logf(fmt string, args ...interface{})
}

// Iterator iterates over key/value pairs in key order
type Iterator interface {
	Key() []byte
	Value() []byte
	Next() bool
	Error() error
	// Iterator 必须在使用完毕后关闭
	Close()
}

// PrefixRange returns key range that satisfy the given prefix
func PrefixRange(prefix []byte) ([]byte, []byte) {
	var limit []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			limit = make([]byte, i+1)
			copy(limit, prefix)
			limit[i] = c + 1
			break
		}
	}
	return prefix, limit
}
