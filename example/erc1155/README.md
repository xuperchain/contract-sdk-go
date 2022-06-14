

## 基于xuperchain的ERC1155合约
允许使用者在同一个智能合约中无限量地重复使用同质化或者非同质化的代币.是可以一次性铸造多种多量同质化及非同质化资产.  
其主要特点是：  
* 既可以发行同质化也可以发行非同质化代币,当对同质化和非同质代币都有需求时都可以在此标准上发行,无需切换别的标准.    
* 可以批量转移代币资产,以及一次操作就可向不同对象转移多个代币资产,大大提高使用效率降低时间及成本.   

## ERC1155参考手册

接口参照: https://eips.ethereum.org/EIPS/eip-1155    
参考资料: https://u.naturaldao.io/be/chapter-4/eip1155    

 ## 使用示例
 ```shell
 ## 统一由 Action 方法完成调用,参数 action 的值是方法名, 例如 "action":"mint"
./bin/xchain-cli native invoke --method Action -a '{"action":"mint","to":"","tokenID":"10000001","amount":"10"}' --fee 110000 erc1155
 ```
 ## 参数说明
每个action都对应一个函数实现,参见函数的参数
```go
switch action {
	case "uri": //获取 URI
		e.uri(tokenID)
	
	case "mint": //铸造Token
		e.mint(to, tokenID, amount, data)

	case "mintBatch": //批量铸造
		e.mintBatch(to, tokenIDs, amounts, data)
		
	case "burn": //销毁Token
		e.burn(from, tokenID, amount)
		
	case "burnBatch": //批量销毁
		e.burnBatch(from, tokenIDs, amounts)
	
	case "safeTransferFrom": //转移Token
		e.safeTransferFrom(from, to, tokenID, amount, data)
		
	case "safeBatchTransferFrom": //批量转移Token
		e.safeBatchTransferFrom(from, to, tokenIDs, amounts, data)

	case "balanceOf": //查询余额
	    e.balanceOf(owner, tokenID)
		
	case "balanceOfBatch": //批量查询余额
		e.balanceOfBatch(owners, tokenIDs)
		
	case "setApprovalForAll": //设置授权
		e.setApprovalForAll(to, approved)

	case "isApprovedForAll": //是否有权限
	    e.isApprovedForAll(owner, to)

	default:
		return code.Errors("Invalid action " + action)
}

```


