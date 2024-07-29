package utils

import (
	"fmt"
	"sync"

	"github.com/huandu/go-clone"
)

var globalVars = make(map[string]interface{})
var rlock sync.RWMutex







func GetValue(keys ...string) (interface{}, error) {
	
	rlock.RLock()
	defer rlock.RUnlock()

	
	tempGlobalVars := globalVars
	for i, key := range keys {
		
		if tempGlobalVars[key] != nil {
			
			if i == len(keys)-1 {
				return clone.Clone(tempGlobalVars[key]), nil
			}
			
			if subMap, ok := tempGlobalVars[key].(map[string]interface{}); ok {
				tempGlobalVars = subMap
			} else {
				
				return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		} else {
			
			return nil, fmt.Errorf("参数%d '%s' 不存在", i+1, key)
		}
	}
	
	return nil, fmt.Errorf("参数不足,缺少键值参数")
}






func SetValue(value interface{}, keys ...string) error {
	
	rlock.Lock()
	defer rlock.Unlock()

	
	tempVars := globalVars
	
	for i, key := range keys {
		
		if i == len(keys)-1 {
			tempVars[key] = value
			return nil
		} else {
			
			if sub_map, ok := tempVars[key].(map[string]interface{}); ok {
				
				tempVars = sub_map
			} else {
				
				return fmt.Errorf("参数%d '%s' 不存在", i+1, key)
			}
		}
	}
	
	return fmt.Errorf("参数不足,缺少键值参数")
}



func DelValue(keys ...string) error {
    
    rlock.Lock()
    defer rlock.Unlock()
    
    
    tempVars := globalVars
    
    
    for i, key := range keys {
        
        if i == len(keys)-1 {
            delete(tempVars, key)
            return nil
        }
        
        
        if subMap, ok := tempVars[key].(map[string]interface{}); ok {
            
            tempVars = subMap
        } else {
            
            return fmt.Errorf("参数%d '%s' 不存在", i+1, key)
        }
    }
    
    
    return fmt.Errorf("参数不足,缺少键值参数")
}
func Show(){
	fmt.Println(globalVars)
}