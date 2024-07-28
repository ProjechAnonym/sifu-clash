package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sifu-clash/models"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DiskDb *gorm.DB
var MemoryDb *gorm.DB
func GetDatabase() error{
	project_dir, err := GetValue("project-dir"); 
	if err != nil {
		return err
	}
	if DiskDb, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/sifu-box.db", project_dir)), &gorm.Config{}); err != nil{
		return err
	}
	if MemoryDb,err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}); err != nil{
		return err
	}
	DiskDb.AutoMigrate(&models.Host{})
	MemoryDb.AutoMigrate(&models.Provider{},&models.Ruleset{})
	return nil
}

func LoadConfig(dst string,class string) error {
	projectDir, err := GetValue("project-dir"); 
	if err != nil {
		return err
	}
	viper.SetConfigFile(filepath.Join(projectDir.(string),dst))
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	switch class {
	case "mode":
		var server models.Server
		if err = viper.Unmarshal(&server); err != nil {
			return err
		}
		SetValue(server,class)
	case "proxy":
		var proxy models.Proxy
		if err = viper.Unmarshal(&proxy); err != nil {
			return err
		}
		if len(proxy.Providers) != 0  {
			if err = MemoryDb.Create(&proxy.Providers).Error; err != nil {
				return err
			}	
		}
		if len(proxy.Rulesets) != 0 {
			if err = MemoryDb.Create(&proxy.Rulesets).Error; err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("类型'%s'不正确",class)
	}
	return nil
}
func LoadTemplate() error{
	projectDir, err := GetValue("project-dir"); 
	if err != nil {
		return err
	}
	templateDir, err := os.Open(filepath.Join(projectDir.(string),"template"))
	if err != nil {
		return err
	}
	defer templateDir.Close()
	files, err := templateDir.ReadDir(-1) // -1 表示读取所有条目
	if err != nil {
		return err
	}
	templateMap := make(map[string]models.Template)
	for _, file := range files{
		var template models.Template
		fileName := strings.Split(file.Name(), ".")[0]
		viper.SetConfigFile(filepath.Join(projectDir.(string),"template",file.Name()))
		if err = viper.ReadInConfig();err != nil {
			return err
		}
		if err = viper.Unmarshal(&template); err != nil {
			return err
		}
		templateMap[fileName] = template
	}
	if err := SetValue(templateMap,"templates"); err != nil {
		return err
	}
	return nil
}
func GetProjectDir() string {
	// base_dir := filepath.Dir(os.Args[0])
	// base_dir := "E:/Myproject/sifu-clash"
	base_dir := "/root/sifu-clash"
	return base_dir
}
