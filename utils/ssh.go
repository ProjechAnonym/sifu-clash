package utils

import (
	"fmt"
	"net"
	"net/url"
	"sifu-clash/models"

	"golang.org/x/crypto/ssh"
)
func tempKeyCallback(host models.Host, key ssh.PublicKey) error {
    // 检查服务器是否有记录的指纹信息
	if host.Fingerprint == ""{
        // 计算并获取公钥的SHA256指纹
		fingerPrint := ssh.FingerprintSHA256(key)
        // 更新数据库中服务器的指纹信息
		if err := DiskDb.Model(&host).Where("url = ?",host.Url).Update("fingerprint",fingerPrint).Error; err != nil{
			return err
		}
	}else{
        // 如果服务器有记录的指纹信息,则与提供的公钥指纹进行比较
		if host.Fingerprint != ssh.FingerprintSHA256(key){
			return fmt.Errorf("fingerprint mismatch")
		}
	}
    // 公钥验证成功,返回nil
	return nil
}

func InitClient(host models.Host) (*ssh.ClientConfig,string,error) {
    // 解析服务器的URL以获取主机名
    hostUrl,err := url.Parse(host.Url)
    if err != nil{
        return nil,"",err
    }
    // 构造SSH服务器地址,默认端口为22
    addr := hostUrl.Hostname() + ":22"
    // 配置SSH客户端配置
    config := &ssh.ClientConfig{
        User: host.Username,
        Auth: []ssh.AuthMethod{ssh.Password(host.Password)},
        // 定制化HostKeyCallback,用于验证服务器的公钥
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            // 调用自定义的公钥验证函数
            if err := tempKeyCallback(host,key); err != nil {
                return err
            }
            return nil
        },
    }
    // 返回配置好的SSH客户端配置、服务器地址和nil错误
    return config,addr,nil
}