package utils

import (
	"io"
	"os"
	"path/filepath"
)

func FileWrite(content []byte, dst string) error {
	
	if _, err := os.Stat(filepath.Dir(dst)); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(dst),0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer func() {
		
		if err := file.Close(); err != nil {
			LoggerCaller("文件无法关闭", err, 1)
		}
	}()
	if err != nil {
		return err
	}

	
	_, err = file.WriteString(string(content))
	if err != nil {
		return err
	}
	
	return nil
}

func DirCreate(src string) error{
    
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
	
	_, err := os.Stat(dst)
	if err != nil {
        if os.IsNotExist(err) {
            
            return nil
        } else {
            
            return err
        }
    }

    
    if err := os.RemoveAll(dst); err != nil {
        return err
    }
	return nil
}

func FileCopy(src, dst string) error{
    
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    
    defer srcFile.Close()
    
    
    targetFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    
    defer targetFile.Close()
    
    
    if _,err = io.Copy(targetFile,srcFile);err != nil{
        return err
    }
    
    
    return nil
}

func FileRead(src string) ([]byte,error){
    
    srcFile, err := os.Open(src)
    if err != nil {
        return nil,err
    }
    
    defer srcFile.Close()
    content,err := io.ReadAll(srcFile)
    if err != nil {
        return nil,err
    }
    return content,nil
}