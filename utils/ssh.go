package utils

import (
	"fmt"
	"net"
	"net/url"
	"sifu-clash/models"

	"golang.org/x/crypto/ssh"
)
func tempKeyCallback(host models.Host, key ssh.PublicKey) error {
    
	if host.Fingerprint == ""{
        
		fingerPrint := ssh.FingerprintSHA256(key)
        
		if err := DiskDb.Model(&host).Where("url = ?",host.Url).Update("fingerprint",fingerPrint).Error; err != nil{
			return err
		}
	}else{
        
		if host.Fingerprint != ssh.FingerprintSHA256(key){
			return fmt.Errorf("fingerprint mismatch")
		}
	}
    
	return nil
}

func InitClient(host models.Host) (*ssh.ClientConfig,string,error) {
    
    hostUrl,err := url.Parse(host.Url)
    if err != nil{
        return nil,"",err
    }
    
    addr := hostUrl.Hostname() + ":22"
    
    config := &ssh.ClientConfig{
        User: host.Username,
        Auth: []ssh.AuthMethod{ssh.Password(host.Password)},
        
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            
            if err := tempKeyCallback(host,key); err != nil {
                return err
            }
            return nil
        },
    }
    
    return config,addr,nil
}