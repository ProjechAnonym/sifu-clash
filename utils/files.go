package utils

import (
	"os"
	"path/filepath"
)

func FileWrite(content []byte, dst string) error {
	// 检查目标文件目录是否存在,若不存在则创建
	if _, err := os.Stat(filepath.Dir(dst)); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(dst),0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 打开(若不存在则创建)文件,准备进行写操作
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer func() {
		// 确保文件在函数返回前关闭,避免资源泄露
		if err := file.Close(); err != nil {
			LoggerCaller("文件无法关闭", err, 1)
		}
	}()
	if err != nil {
		return err
	}

	// 将内容写入文件
	_, err = file.WriteString(string(content))
	if err != nil {
		return err
	}
	// 操作成功,返回nil
	return nil
}

func DirCreate(src string) error{
    // 检查临时目录是否存在,不存在则创建临时目录
    if _,err := os.Stat(src); err != nil {
        if os.IsNotExist(err){
            if err := os.MkdirAll(src,0755); err != nil{
                return err
            }
        }else{
            return err
        }
    }
    return nil
}

func FileDelete(dst string) error{
	// 检查目标文件是否存在,若存在则删除
	_, err := os.Stat(dst)
	if err != nil {
        if os.IsNotExist(err) {
            // 文件不存在,不需要删除,直接返回
            return nil
        } else {
            // 其他错误,例如权限问题,返回错误
            return err
        }
    }

    // 尝试删除文件
    if err := os.RemoveAll(dst); err != nil {
        return err
    }
	return nil
}