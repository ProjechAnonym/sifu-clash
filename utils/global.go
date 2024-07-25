package utils

import (
	"fmt"
	"sync"

	"github.com/huandu/go-clone"
)

var globalVars = make(map[string]interface{})
var rlock sync.RWMutex

// GetValue 根据提供的键序列从全局变量中获取最终对应的值
// 如果所有键都存在,则返回最终的值和nil错误
// 如果任何键不存在,则返回nil和相应的错误
// keys: 键序列,用于深度查找全局变量中的值
// 返回值: 最终查找到的值,如果未找到则为nil
// 返回错误: 如果任何键不存在或查找过程中遇到非映射类型,则返回错误信息
func GetValue(keys ...string) (interface{}, error) {
	// 以读锁方式保护全局变量访问,确保并发安全
	rlock.RLock()
	defer rlock.RUnlock()

	// 创建tempGlobalVars的副本以避免直接操作全局变量,减少并发风险
	tempGlobalVars := globalVars
	for i, key := range keys {
		// 检查当前键是否存在,如果不存在则返回错误
		if tempGlobalVars[key] != nil {
			// 如果是最后一个键,返回对应的值的副本和nil错误
			if i == len(keys)-1 {
				return clone.Clone(tempGlobalVars[key]), nil
			}
			// 如果当前值是映射,更新tempGlobalVars以继续下一层查找
			if subMap, ok := tempGlobalVars[key].(map[string]interface{}); ok {
				tempGlobalVars = subMap
			} else {
				// 如果当前值不是映射,返回错误
				return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		} else {
			// 如果当前键不存在,返回错误
			return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
		}
	}
	// 如果没有提供任何键或查找过程没有找到任何值,返回错误
	return nil, fmt.Errorf("参数不足,缺少键值参数")
}


// SetValue 用于递归地设置嵌套字典中指定键的值
// 它首先锁定全局变量以确保并发安全,然后根据提供的键路径更新值
// 参数value是要设置的新值,keys是键的序列,用于在嵌套字典中定位最终要设置的值
// 函数返回一个错误,如果在更新过程中遇到任何问题,例如找不到中间键对应的子字典
func SetValue(value interface{}, keys ...string) error {
	// 锁定全局变量以确保并发安全
	rlock.Lock()
	defer rlock.Unlock()

	// 创建一个局部变量tempVars,用于在更新过程中操作,以避免直接修改全局变量
	tempVars := globalVars
	// 遍历键序列
	for i, key := range keys {
		// 如果当前键是最后一个键,则设置值并返回nil,表示更新成功
		if i == len(keys)-1 {
			tempVars[key] = value
			return nil
		} else {
			// 尝试将当前键对应的值转换为map[string]interface{}类型
			if sub_map, ok := tempVars[key].(map[string]interface{}); ok {
				// 如果转换成功,将tempVars更新为这个子字典,为下一次迭代做准备
				tempVars = sub_map
			} else {
				// 如果转换失败,表示当前键不是期望的字典类型,返回错误
				return fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		}
	}
	// 如果遍历完成后还没有返回,表示键序列不正确,没有指定要设置的值,返回错误
	return fmt.Errorf("参数不足,缺少键值参数")
}
// DelValue 递归删除全局变量中指定的键值
// keys: 需要删除的键值路径,多个键值用数组表示,用于深入到特定的嵌套级别
// 返回值: 如果成功删除,则返回nil；如果删除过程中遇到任何问题,则返回错误
func DelValue(keys ...string) error {
    // 加锁以确保并发安全地访问全局变量
    rlock.Lock()
    defer rlock.Unlock()
    
    // 创建tempVars作为globalVars的本地副本,以避免直接修改全局变量
    tempVars := globalVars
    
    // 遍历keys数组,根据键值路径深入到tempVars的适当嵌套级别
    for i, key := range keys {
        // 如果当前键是最后一个键,则删除它并返回nil,表示删除成功
        if i == len(keys)-1 {
            delete(tempVars, key)
            return nil
        }
        
        // 尝试将当前键对应的值转换为map[string]interface{}类型
        if subMap, ok := tempVars[key].(map[string]interface{}); ok {
            // 如果转换成功,则将tempVars更新为这个子map,继续遍历
            tempVars = subMap
        } else {
            // 如果转换失败,则返回错误,表示指定的键不是map类型
            return fmt.Errorf("参数%d '%s' 不存在", i+1, key)
        }
    }
    
    // 如果遍历完成后没有返回错误,则表示没有提供足够的键值路径来删除值,返回相应错误
    return fmt.Errorf("参数不足,缺少键值参数")
}
func Show(){
	fmt.Println(globalVars)
}