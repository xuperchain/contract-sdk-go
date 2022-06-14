// @Title xuperchain 的 erc1155合约
// @Description 接口参照:https://eips.ethereum.org/EIPS/eip-1155  资料参考:https://u.naturaldao.io/be/chapter-4/eip1155
// @Author 盛见网络-springrain

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuperchain/contract-sdk-go/code"
	"github.com/xuperchain/contract-sdk-go/driver"
)

// ErrNotFound is returned when key is not found
var errNotFound = "Key not found"

type erc1155 struct {
	ctx code.Context
}

//初始化方法,接口要求,必须实现
func (e *erc1155) Initialize(ctx code.Context) code.Response {
	//https://token-cdn-domain/{id}.json 字符串,非必须
	metaDataURI := ctx.Args()["metaDataURI"]
	if metaDataURI != nil && len(metaDataURI) > 0 {
		ctx.PutObject([]byte("/metaDataURI/"), metaDataURI)
	}
	return code.OK(nil)
}

func main() {
	driver.Serve(new(erc1155))
}

// 操作都通过这个方法进行处理
func (e *erc1155) Action(ctx code.Context) code.Response {
	e.ctx = ctx
	args := ctx.Args()
	action := string(args["action"])
	if action == "" {
		return code.Errors("Missing key: action")
	}

	from := string(args["from"])
	to := string(args["to"])
	//是否有授权
	approved := false
	approvedStr := string(args["approved"])
	if strings.EqualFold(approvedStr, "true") {
		approved = true
	}
	//token 拥有者
	owner := string(args["owner"])
	owners := make([]string, 0)
	ownersBytes := args["owners"]
	if ownersBytes != nil {
		json.Unmarshal(ownersBytes, &owners)
	}

	//tokenID token 类型
	var tokenID int64
	tokenIDStr := string(args["tokenID"])
	if tokenIDStr != "" {
		tokenID, _ = strconv.ParseInt(tokenIDStr, 10, 64)
	}
	tokenIDs := make([]int64, 0)
	tokenIDsBytes := args["tokenIDs"]
	if tokenIDsBytes != nil {
		json.Unmarshal(tokenIDsBytes, &tokenIDs)
	}

	//数量
	var amount int64
	amountStr := string(args["amount"])
	if amountStr != "" {
		amount, _ = strconv.ParseInt(amountStr, 10, 64)
	}
	amounts := make([]int64, 0)
	amountsBytes := args["amounts"]
	if amountsBytes != nil {
		json.Unmarshal(amountsBytes, &amounts)
	}
	//原始数据,用于处理 {to} 是{IERC1155Receiver-onERC1155Received} 合约,用于传递原始数据
	//暂时没有处理转给合约这种情况
	data := args["data"]

	switch action {
	case "uri": //获取 URI
		s, err := e.uri(tokenID)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		uri(ctx, s, tokenID)
		return code.OK([]byte(s))

	case "mint": //铸造Token
		if to == "" {
			to = ctx.Caller()
		}
		err := e.mint(to, tokenID, amount, data)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		transferSingle(ctx, ctx.Caller(), "", to, tokenID, amount)
		return code.OK(nil)

	case "mintBatch": //批量铸造
		if to == "" {
			to = ctx.Caller()
		}
		err := e.mintBatch(to, tokenIDs, amounts, data)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		transferBatch(ctx, ctx.Caller(), "", to, tokenIDs, amounts)
		return code.OK(nil)
	case "burn": //销毁Token
		if from == "" {
			from = ctx.Caller()
		}
		err := e.burn(from, tokenID, amount)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		transferSingle(ctx, ctx.Caller(), from, "", tokenID, amount)
		return code.OK(nil)
	case "burnBatch": //批量销毁
		if from == "" {
			from = ctx.Caller()
		}
		err := e.burnBatch(from, tokenIDs, amounts)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		transferBatch(ctx, ctx.Caller(), from, "", tokenIDs, amounts)
		return code.OK(nil)
	case "safeTransferFrom": //转移Token
		if from == "" {
			from = ctx.Caller()
		}
		err := e.safeTransferFrom(from, to, tokenID, amount, data)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		transferSingle(ctx, ctx.Caller(), from, to, tokenID, amount)
		return code.OK(nil)
	case "safeBatchTransferFrom": //批量转移Token
		if from == "" {
			from = ctx.Caller()
		}
		err := e.safeBatchTransferFrom(from, to, tokenIDs, amounts, data)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		transferBatch(ctx, ctx.Caller(), from, to, tokenIDs, amounts)
		return code.OK(nil)

	case "balanceOf": //查询余额
		if owner == "" {
			owner = ctx.Caller()
		}
		tokenAmount, err := e.balanceOf(owner, tokenID)
		if err != nil {
			return code.Error(err)
		}
		ta, err := json.Marshal(tokenAmount)
		if err != nil {
			return code.Error(err)
		}
		return code.OK(ta)
	case "balanceOfBatch": //批量查询余额
		tokenAmounts, err := e.balanceOfBatch(owners, tokenIDs)
		if err != nil {
			return code.Error(err)
		}
		ta, err := json.Marshal(tokenAmounts)
		if err != nil {
			return code.Error(err)
		}
		return code.OK(ta)
	case "setApprovalForAll": //设置授权
		err := e.setApprovalForAll(to, approved)
		if err != nil {
			return code.Error(err)
		}
		//调用通知事件
		approvalForAll(ctx, ctx.Caller(), to, approved)
		return code.OK(nil)

	case "isApprovedForAll": //是否有权限
		if owner == "" {
			owner = ctx.Caller()
		}
		approved, err := e.isApprovedForAll(owner, to)
		if err != nil {
			return code.Error(err)
		}
		ta, err := json.Marshal(approved)
		if err != nil {
			return code.Error(err)
		}
		return code.OK(ta)

	default:
		return code.Errors("Invalid action " + action)
	}
}

/**
  TransferSingle 或 TransferBatch 事件: 必须在token转移时发出,包括ZERO参与的铸造或销毁(参见标准的"安全转移规则"部分)
  铸造token时, {from} 必须设置为 "" 零值地址
  销毁token时, {to} 必须设置为 "" 零值地址

  @param operator 必须是被批准进行转账的账户/合约的地址,需要from地址提前授权给{operator},参见ApprovalForAll方法
  @param from 转出token的账户地址
  @param to 接收token的账户地址
  @param tokenID 转移的token类型
  @param amount  转移的token数量
*/
// event TransferSingle(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
func transferSingle(ctx code.Context, operator string, from string, to string, tokenID int64, amount int64) {
	eventMap := make(map[string]interface{})
	eventMap["operator"] = operator
	eventMap["from"] = from
	eventMap["to"] = to
	eventMap["tokenID"] = tokenID
	eventMap["amount"] = amount
	//事件名和ERC-1155保持一致
	ctx.EmitJSONEvent("TransferSingle", eventMap)
}

/**
  TransferSingle 或 TransferBatch 事件: 必须在token转移时发出,包括ZERO参与的铸造或销毁(参见标准的"安全转移规则"部分)
  铸造token时, {from} 必须设置为 "" 零值地址
  销毁token时, {to} 必须设置为 "" 零值地址

  @param operator 必须是被批准进行转账的账户/合约的地址,需要from地址提前授权给{operator},参见ApprovalForAll方法
  @param from 转出token的账户地址
  @param to 接收token的账户地址
  @param tokenIDs 转移的token类型数组
  @param amounts  转移的token数量数组,和{tokenIDs}顺序对应
*/
// event TransferBatch(address indexed _operator, address indexed _from, address indexed _to, uint256[] _ids, uint256[] _values);
func transferBatch(ctx code.Context, operator string, from string, to string, tokenIDs []int64, amounts []int64) {
	eventMap := make(map[string]interface{})
	eventMap["operator"] = operator
	eventMap["from"] = from
	eventMap["to"] = to
	eventMap["tokenIDs"] = tokenIDs
	eventMap["amounts"] = amounts
	//事件名和ERC-1155保持一致
	ctx.EmitJSONEvent("TransferBatch", eventMap)
}

/**
  {owner}授权{to}事件: 可以代为操作{owner}的所有token,默认没有授权

  @param owner token拥有者账户
  @param to 被授权的账户
  @param approved 是否允许
*/
//event ApprovalForAll(address indexed _owner, address indexed _operator, bool _approved);
func approvalForAll(ctx code.Context, owner string, to string, approved bool) {
	eventMap := make(map[string]interface{})
	eventMap["owner"] = owner
	eventMap["to"] = to
	eventMap["approved"] = approved
	//事件名和ERC-1155保持一致
	ctx.EmitJSONEvent("ApprovalForAll", eventMap)
}

/**
  更新token的URI时,调用此事件方法(非必须,还待理解ERC1155Metadata_URI)
  参考文档:https://eips.ethereum.org/EIPS/eip-1155#metadata-extensions
*/
// event URI(string _value, uint256 indexed _id);
func uri(ctx code.Context, uri string, tokenID int64) {
	eventMap := make(map[string]interface{})
	eventMap["uri"] = uri
	eventMap["tokenID"] = tokenID
	//事件名和ERC-1155保持一致
	ctx.EmitJSONEvent("URI", eventMap)
}

//ERC1155Metadata_URI 接口的 function uri(uint256 _id) external view returns (string memory);
//应该是 返回一个 https://token-cdn-domain/{id}.json 字符串,id会使用tokenID 替换成64位字符串,
//不够使用0补齐,类似https://token-cdn-domain/000000000000000000000000000000000000000000000000000000000004cce0.json
func (e *erc1155) uri(tokenID int64) (string, error) {
	uriBytes, err := e.ctx.GetObject([]byte("/metaDataURI/"))
	if err != nil {
		return "", fmt.Errorf("uri-->e.ctx.GetObject:%w", err)
	}
	return string(uriBytes), nil
}

/**
  为 {to} 账户铸造 {amount} 数量的  {tokenID}
  触发 {transferSingle} 事件
  {to} 账户不能是 "",如果是合约,必须实现 {IERC1155Receiver-onERC1155Received} , 并且返回值准确

  TODO 后期考虑是不是只有指定账户才能铸造Token

  @param to 铸造token的账户地址
  @param tokenID 铸造的token类型
  @param amount  铸造的token数量
  @param data 当参数 {to} 是智能合约时,交易发送者(合约调用方)必须将参数 {data} 的原始数据,不做任何改动地传递给{to}合约onERC1155Received接口函数的参数 "data"
*/
func (e *erc1155) mint(to string, tokenID int64, amount int64, data []byte) error {
	if to == "" || tokenID < 1 || amount < 1 {
		return fmt.Errorf("mint: error to is %v  ,tokenID is %v, amount is %v ", to, tokenID, amount)
	}
	tokenIDStr := strconv.FormatInt(tokenID, 10)

	// {tokenID} 数量的 key
	balanceTokenKey := "/token/" + tokenIDStr
	// 增加 {tokenID}的 {amount}
	err := e.addAmount(balanceTokenKey, amount)
	if err != nil {
		return fmt.Errorf("mint-->e.addAmount:%w", err)
	}
	//{to}拥有的{tokenID}数量的key
	balanceToKey := "/balance/" + tokenIDStr + "/" + to
	// 增加{to} 的 {tokenID}数量
	err = e.addAmount(balanceToKey, amount)
	if err != nil {
		return fmt.Errorf("mint-->e.addAmount:%w", err)
	}
	//transferSingle 事件通知由 Action 入口函数统一处理
	return nil
}

/**
  为 {to} 账户 批量 铸造 {amounts} 数量 的 {tokenIDs}
  触发 {transferBatch} 事件
  {to} 账户不能是 "",如果是合约,必须实现 {IERC1155Receiver-onERC1155Received} , 并且返回值准确

  @param to 铸造token的账户地址
  @param tokenIDs 铸造的token类型
  @param amounts  铸造的token数量
  @param data 当参数 {to} 是智能合约时,交易发送者(合约调用方)必须将参数 {data} 的原始数据,不做任何改动地传递给{to}合约onERC1155Received接口函数的参数 "data"
*/
func (e *erc1155) mintBatch(to string, tokenIDs []int64, amounts []int64, data []byte) error {
	if len(tokenIDs) != len(amounts) {
		return fmt.Errorf("mintBatch: tokenIDs len is %v ,amounts  len is %v ", len(tokenIDs), len(amounts))
	}
	for i := 0; i < len(tokenIDs); i++ {
		tokenID := tokenIDs[i]
		amount := amounts[i]
		err := e.mint(to, tokenID, amount, data)
		if err != nil {
			return fmt.Errorf("mintBatch-->e.mint:%w", err)
		}
	}
	//transferBatch 事件通知由 Action 入口函数统一处理
	return nil
}

/**
  销毁 {from} 账户 {amount} 个 {tokenID}
  {from}账户拥有的{tokenID}数量必须大于等于{amount}
  触发 {transferSingle} 事件
  {from} 账户不能是 "",如果是合约,必须实现 {IERC1155Receiver-onERC1155Received} , 并且返回值准确

  @param from 销毁token的账户地址
  @param tokenID 销毁的token类型
  @param amount  销毁的token数量
*/
func (e *erc1155) burn(from string, tokenID int64, amount int64) error {
	if from == "" || tokenID < 1 || amount < 1 {
		return fmt.Errorf("burn: error from is %v  ,tokenID is %v, amount is %v ", from, tokenID, amount)
	}
	// {from} 账号的 {tokenID} 数量减少 {amount} 个
	err := e.subAmount(from, tokenID, amount)
	if err != nil {
		return fmt.Errorf("burn-->e.subAmount:%w", err)
	}
	//{tokenID} 数量的 key
	balanceTokenKey := "/token/" + strconv.FormatInt(tokenID, 10)
	//根据key 数量减少{amount} 个
	err = e.subAmountByKey(balanceTokenKey, amount)
	if err != nil {
		return fmt.Errorf("burn-->e.subAmountByKey:%w", err)
	}

	//transferSingle 事件通知由 Action 入口函数统一处理
	return nil
}

/**
  批量销毁 {from} 账户 {amounts} 个 {tokenIDs}
  触发 {transferBatch} 事件
  {from} 账户不能是 "",如果是合约,必须实现 {IERC1155Receiver-onERC1155Received} , 并且返回值准确

  @param from 销毁token的账户地址
  @param tokenIDs 销毁的token类型
  @param amounts  销毁的token数量
*/
func (e *erc1155) burnBatch(from string, tokenIDs []int64, amounts []int64) error {
	if len(tokenIDs) != len(amounts) {
		return fmt.Errorf("burnBatch: tokenIDs len is %v ,amounts  len is %v ", len(tokenIDs), len(amounts))
	}
	for i := 0; i < len(tokenIDs); i++ {
		tokenID := tokenIDs[i]
		amount := amounts[i]
		err := e.burn(from, tokenID, amount)
		if err != nil {
			return fmt.Errorf("burnBatch-->e.burn:%w", err)
		}
	}
	//transferBatch 事件通知由 Action 入口函数统一处理
	return nil
}

/**
  将 {amount} 个 {tokenID} 从 {from}账户 转移到 {to}账户 (安全的)
  调用者必须得到授权才可以从参数 from 指定的账户转账token
  如果参数 {to} 设定的地址为零地址,则交易必须回滚
  对参数 {tokenID} 所指的token,如果其转出账户的余额小于参数 {amount} 所定义的金额,则交易必须回滚
  如果转账交易出现任何其它错误,交易也必须回滚
  必须触发TransferSingle事件以反映账户余额的变化(请参看TransferSingle and TransferBatch event rules章节)
  当上述条件都满足时,safeTransferFrom函数就必须检查参数 {to} 是否是一个智能合约地址(比如检查code size是否大于0),如果是,就必须在参数{to}的智能合约上调用{IERC1155Receiver-onERC1155Received}函数并执行相应的操作(请参看"onERC1155Received rules"章节)

  @param from 转出token的账户地址
  @param to 接收token的账户地址,不能是零值
  @param tokenID 转移的token类型
  @param amount  转移的token数量
  @param data 当参数 {to} 是智能合约时,交易发送者(合约调用方)必须将参数 {data} 的原始数据,不做任何改动地传递给{to}合约{IERC1155Receiver-onERC1155Received}接口函数的参数 "data"
*/
//function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _value, bytes calldata _data) external;
func (e *erc1155) safeTransferFrom(from string, to string, tokenID int64, amount int64, data []byte) error {
	if to == "" || tokenID < 1 || amount < 1 || from == to {
		return fmt.Errorf("safeTransferFrom: error from is %v , to is %v  ,tokenID is %v, amount is %v ", from, to, tokenID, amount)
	}
	// {from} 账号的 {tokenID} 数量减少 {amount} 个
	err := e.subAmount(from, tokenID, amount)
	if err != nil {
		return fmt.Errorf("safeTransferFrom-->e.subAmount:%w", err)
	}
	// {to} 拥有的token数量
	balanceToKey := "/balance/" + strconv.FormatInt(tokenID, 10) + "/" + to
	// 增加{to} 的 {tokenID}数量
	err = e.addAmount(balanceToKey, amount)
	if err != nil {
		return fmt.Errorf("safeTransferFrom-->e.addAmount:%w", err)
	}
	//transferSingle 事件通知由 Action 入口函数统一处理
	return nil
}

/**
  将 {amounts} 个 {tokenIDs} 从 {from}账户 转移到 {to}账户 (安全的)
  调用者必须得到授权才可以从参数 from 指定的账户转账token
  如果参数 {to} 设定的地址为零地址,则交易必须回滚
  如果参数{tokenIDs}的长度和参数{amounts}的长度不同,则交易必须回滚
  对参数{tokenIDs}所指的token,如果其任一转出账户的余额小于该交易所对应的参数{amounts}所定义的金额,则交易必须回滚
  如果转账交易出现任何其它错误,交易也必须回滚
  必须触发TransferSingle或TransferBatch事件以反映账户余额的变化(请参看TransferSingle and TransferBatch event rules章节)
  所有账户余额的变化和事件的触发必须按其被提交的顺序发生(即tokenIDs[0]/amounts[0]在tokenIDs[1]/amounts[1]之前发生……)
  当上述条件都满足时,safeTransferFrom函数就必须检查参数 {to} 是否是一个智能合约地址(比如检查code size是否大于0),如果是,就必须在参数{to}的智能合约上调用{IERC1155Receiver-onERC1155Received}函数并执行相应的操作(请参看"onERC1155Received rules"章节)

  @param from 转出token的账户地址
  @param to 接收token的账户地址
  @param tokenIDs 转移的token类型数组
  @param amounts  转移的token数量数组,和{tokenIDs}顺序对应
  @param data 当参数 {to} 是智能合约时,交易发送者(合约调用方)必须将参数 {data} 的原始数据,不做任何改动地传递给{to}合约{IERC1155Receiver-onERC1155Received}接口函数的参数 "data"
*/
//function safeBatchTransferFrom(address _from, address _to, uint256[] calldata _ids, uint256[] calldata _values, bytes calldata _data) external;
func (e *erc1155) safeBatchTransferFrom(from string, to string, tokenIDs []int64, amounts []int64, data []byte) error {
	if len(tokenIDs) != len(amounts) {
		return fmt.Errorf("safeBatchTransferFrom: tokenIDs len is %v ,amounts  len is %v ", len(tokenIDs), len(amounts))
	}
	for i := 0; i < len(tokenIDs); i++ {
		tokenID := tokenIDs[i]
		amount := amounts[i]
		err := e.safeTransferFrom(from, to, tokenID, amount, data)
		if err != nil {
			return fmt.Errorf("safeBatchTransferFrom-->e.safeTransferFrom:%w", err)
		}
	}
	//transferBatch 事件通知由 Action 入口函数统一处理
	return nil
}

/**
  获取 {owner} 账户中的token余额,需要授权才能查询账户余额
  @param owner  token持有者账户
  @param tokenID token的ID/类型
  @return 返回 {owner} 账户中的token余额
*/
//function balanceOf(address _owner, uint256 _id) external view returns (uint256);
func (e *erc1155) balanceOf(owner string, tokenID int64) (int64, error) {
	if owner == "" || tokenID < 1 {
		return 0, fmt.Errorf("balanceOf: error owner is %v  ,tokenID is %v", owner, tokenID)
	}

	operator := e.ctx.Caller()
	isApproved, err := e.isApprovedForAll(owner, operator)
	if err != nil {
		return 0, fmt.Errorf("balanceOf-->e.isApprovedForAll:%w", err)
	}
	if !isApproved { //没有授权
		return 0, fmt.Errorf("balanceOf-->e.isApprovedForAll:operator(%v) is not owner nor approved", operator)
	}
	//{owner}拥有的{tokenID}数量
	balanceOwnerKey := "/balance/" + strconv.FormatInt(tokenID, 10) + "/" + owner
	balanceOwnerBytes, err := e.ctx.GetObject([]byte(balanceOwnerKey))
	if err != nil {
		return 0, fmt.Errorf("balanceOf-->e.ctx.GetObject:%w", err)
	}
	var balanceOwner int64
	err = json.Unmarshal(balanceOwnerBytes, &balanceOwner)
	if err != nil {
		return 0, fmt.Errorf("balanceOf-->json.Unmarshal:%w", err)
	}
	return balanceOwner, nil
}

/**
  获取多个账户/token对的余额
   @param owners token持有者账户数组
   @param tokenIDs token类型数组
   @return  owner 所请求的token类型的余额(即每个 (owner,tokenID)对 的余额）
*/
//function balanceOfBatch(address[] calldata _owners, uint256[] calldata _ids) external view returns (uint256[] memory);
func (e *erc1155) balanceOfBatch(owners []string, tokenIDs []int64) ([]int64, error) {
	if owners == nil || tokenIDs == nil || len(owners) < 1 || len(owners) != len(tokenIDs) {
		return nil, fmt.Errorf("balanceOfBatch: parameter error, please check ")
	}
	_allLen := len(owners)
	balances := make([]int64, _allLen)
	for i := 0; i < _allLen; i++ {
		owner := owners[i]
		tokenID := tokenIDs[i]
		balance, err := e.balanceOf(owner, tokenID)
		if err != nil {
			return nil, fmt.Errorf("balanceOfBatch-->e.balanceOf:%w", err)
		}
		balances[i] = balance
	}
	return balances, nil
}

/**
  授权或撤销其他操作者{to}管理其所有的token
  成功时必须调用 ApprovalForAll 事件函数
  @param to   被授权人账户
  @param approved  true 给{to}授权,false 撤销 {to}授权
*/
//function setApprovalForAll(address _operator, bool _approved) external;
func (e *erc1155) setApprovalForAll(to string, approved bool) error {
	if to == "" {
		return fmt.Errorf("setApprovalForAll: to is empty ")
	}
	owner := e.ctx.Caller()
	if owner == to {
		return fmt.Errorf("setApprovalForAll: setting approval status for self")
	}
	approvalKey := "/approval/" + owner + "/" + to
	approvedBytes, err := json.Marshal(approved)
	if err != nil {
		return fmt.Errorf("setApprovalForAll-->json.Marshal:%w", err)
	}
	err = e.ctx.PutObject([]byte(approvalKey), approvedBytes)
	//ApprovalForAll 事件函数
	return err

}

/**
  {owner} 给 {to} 的全部授权状态
  @param owner     token的拥有者账户
  @param operator  被授权人账户
  @return          true 是已授权, false 是无授权
*/
//function isApprovedForAll(address _owner, address _operator) external view returns (bool);
func (e *erc1155) isApprovedForAll(owner string, operator string) (bool, error) {
	if owner == "" || operator == "" {
		return false, fmt.Errorf("isApprovedForAll:owner or operator is empty ")
	}

	if owner == operator {
		return true, nil
	}
	approvalKey := "/approval/" + owner + "/" + operator
	isApproved, err := e.ctx.GetObject([]byte(approvalKey))
	if err != nil {
		return false, fmt.Errorf("setApprovalForAll-->e.ctx.GetObject:%w", err)
	}
	approval := false
	err = json.Unmarshal(isApproved, &approval)
	return approval, fmt.Errorf("setApprovalForAll-->json.Unmarshal:%w", err)
}

//-------------------------兼容ERC-165 开始(目前接收者必须是账户,暂时不支持是合约)---------------------------//
// 接口兼容 ERC-165 (i.e. `bytes4(keccak256('supportsInterface(bytes4)'))`).
func (e *erc1155) supportsInterface(interfaceID [4]byte) bool {
	//return n == 0x01ffc9a7 || n == 0x4e2312e0
	return true
}

/**
  用于ERC-1155合约转账接收者{to}是合约的场景.实现 ERC-165,处理单个ERC-1155 token类型的接收
  必须在调用转账函数safeTransferFrom完成,TransferSingle事件通知之后,并且事件要反映账户的余额变化,再使用接收者合约调用{IERC1155Receiver-onERC1155Received}函数
  本函数执行成功应返回 0xf23a6e61,其他认为异常,调用者的交易回滚
  如果调用被拒绝,交易函数必须回滚(This function MUST revert if it rejects the transfer)
  接收者不是智能合约:不应该被外部账户(Externally Owned Account简称EOA)调用
  不应该在除挖矿或转账之外的其它操作中被调用
  作为接收者的合约没有实现ERC1155TokenReceiver接口中相应的函数,交易必须被回滚并给出警告信息
  接收合约实现了ERC1155TokenReceiver中相应的接口函数但返回一个未知值或者抛出错误,交易必须被回滚

  @param operator 发起转账的地址,必须是被批准进行转账的账户/合约的地址,需要from地址提前授权给{operator},参见ApprovalForAll方法
  @param from 转出token的账户地址
  @param tokenID 转移的token类型
  @param amount  转移的token数量
  @param data      交易发送者传递的原始 {data} 参数
  @return          如果成功返回 `bytes4(keccak256("onERC1155Received(address,address,uint256,uint256,bytes)"))` ,其他返回值都要回滚
*/
//function onERC1155Received(address _operator, address _from, uint256 _id, uint256 _value, bytes calldata _data) external returns(bytes4);
func (e *erc1155) onERC1155Received(operator string, from string, tokenID int64, amount int64, data []byte) [4]byte {

	//如果处理成功
	//bytes4(keccak256("onERC1155Received(address,address,uint256,uint256,bytes)"))` (i.e. 0xf23a6e61)
	b := [4]byte{242, 58, 110, 97}
	return b
}

/**
  用于ERC-1155合约转账接收者{to}是合约的场景.实现 ERC-165,处理多个ERC-1155 token类型的接收
  必须在调用转账函数safeBatchTransferFrom完成,TransferBatch事件通知之后,并且事件要反映账户的余额变化,再使用接收者合约调用onERC1155BatchReceived函数
  本函数执行成功应返回 0xbc197c81,其他认为异常,调用者的交易回滚
  如果调用被拒绝,交易函数必须回滚(This function MUST revert if it rejects the transfer)
  接收者不是智能合约:不应该被外部账户(Externally Owned Account简称EOA)调用
  不应该在除挖矿或转账之外的其它操作中被调用
  作为接收者的合约没有实现ERC1155TokenReceiver接口中相应的函数,交易必须被回滚并给出警告信息
  接收合约实现了ERC1155TokenReceiver中相应的接口函数但返回一个未知值或者抛出错误,交易必须被回滚

  @param operator 发起转账的地址,必须是被批准进行转账的账户/合约的地址,需要from地址提前授权给{operator},参见ApprovalForAll方法
  @param from 转出token的账户地址
  @param tokenIDs 转移的token类型数组(长度和{amounts}保持一致)
  @param amounts  转移的token数量数组(长度和{tokenIDs}保持一致)
  @param data      交易发送者传递的原始 {data} 参数
  @return        如果成功返回 `bytes4(keccak256("onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)"))` 0xbc197c81,其他返回值都要回滚
*/
//function onERC1155BatchReceived(address _operator, address _from, uint256[] calldata _ids, uint256[] calldata _values, bytes calldata _data) external returns(bytes4);
func (e *erc1155) onERC1155BatchReceived(operator string, from string, tokenIDs []int64, amounts []int64, data []byte) [4]byte {

	//如果处理成功
	//bytes4(keccak256("onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)"))` (i.e. 0xbc197c81)
	b := [4]byte{188, 25, 124, 129}
	return b
}

//-------------------------兼容ERC-165 结束-------------------------------//

/**
  给指定的 {key} 增加 {amount} 数量
  操作者必须有{from}授权

  @param key 需要操作的key
  @param amount  增加的token数量
*/
func (e *erc1155) addAmount(key string, amount int64) error {
	var balanceToken int64 = 0
	//获取现有的数量,如果不存在key,使用默认值
	balanceTokenBytes, getObjectErr := e.ctx.GetObject([]byte(key))
	if getObjectErr != nil {
		if strings.Contains(getObjectErr.Error(), errNotFound) {
			balanceToken = 0
		} else {
			return fmt.Errorf("addAmount-->e.ctx.GetObject:%w", getObjectErr)
		}

	} else {
		err := json.Unmarshal(balanceTokenBytes, &balanceToken)
		if err != nil {
			return fmt.Errorf("addAmount-->json.Unmarshal:%w", err)
		}
	}

	// 增加{tokenID}数量
	value := balanceToken + amount
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("addAmount-->json.Marshal:%w", err)
	}
	err = e.ctx.PutObject([]byte(key), valueBytes)
	return err
}

/**
  {from} 的 {tokenID} 数量减少 {amount} 个
  操作者必须有{from}授权

  @param from 转出token的账户地址
  @param tokenID 减少的token类型
  @param amount  减少的token数量
*/
func (e *erc1155) subAmount(from string, tokenID int64, amount int64) error {
	operator := e.ctx.Caller()
	isApproved, err := e.isApprovedForAll(from, operator)
	if err != nil {
		return fmt.Errorf("subAmount-->e.isApprovedForAll:%w", err)
	}
	if !isApproved { //没有授权
		return fmt.Errorf("subAmount-->e.isApprovedForAll: operator(%v) is not owner nor approved", operator)
	}
	//{tokenID} 数量的 key
	balanceFromKey := "/balance/" + strconv.FormatInt(tokenID, 10) + "/" + from
	//根据key 数量减少{amount} 个
	return e.subAmountByKey(balanceFromKey, amount)
}

/**
  给指定的 {key} 数量减少 {amount} 个
  操作者必须有{from}授权

  @param key 需要操作的key
  @param amount  减少的token数量
*/
func (e *erc1155) subAmountByKey(key string, amount int64) error {
	balanceTokenBytes, err := e.ctx.GetObject([]byte(key))
	if err != nil {
		return fmt.Errorf("subAmountByKey-->e.e.ctx.GetObject:%w", err)
	}
	var balanceToken int64
	err = json.Unmarshal(balanceTokenBytes, &balanceToken)
	if err != nil {
		return fmt.Errorf("subAmountByKey-->json.Unmarshal:%w", err)
	}
	// 减少{tokenID}的数量
	value := balanceToken - amount
	if value < 0 {
		return fmt.Errorf("subAmountByKey: balanceToken - amount < 0 , is %v ", balanceToken)
	}
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("subAmountByKey-->json.Marshal:%w", err)
	}
	err = e.ctx.PutObject([]byte(key), valueBytes)
	return err
}
