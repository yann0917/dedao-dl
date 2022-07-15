package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/yann0917/dedao-dl/services"
)

const (
	// EnvConfigDir 配置路径环境变量
	EnvConfigDir = "DEDAO_GO_CONFIG_DIR"
	// ConfigName 配置文件名
	ConfigName = "config.json"
)

var (
	configFilePath = filepath.Join(GetConfigDir(), ConfigName)

	// Instance 配置信息 全局调用
	Instance = New(configFilePath)
)

// DedaoUsers user
type DedaoUsers []*Dedao

// ConfigsData Configs data
type ConfigsData struct {
	AcitveUID      string
	DownloadPath   string
	Users          DedaoUsers
	activeUser     *Dedao
	configFilePath string
	configFile     *os.File
	fileMu         sync.Mutex
	service        *services.Service
}

type configJSONExport struct {
	AcitveUID string
	Users     DedaoUsers
}

// Init 初始化配置
func (c *ConfigsData) Init() error {
	if c.configFilePath == "" {
		return errors.New("配置文件未找到")
	}

	// 初始化默认配置
	c.initDefaultConfig()
	// 从配置文件中加载配置
	err := c.loadConfigFromFile()
	if err != nil {
		return err
	}

	// 初始化登陆用户信息
	err = c.initActiveUser()
	if err != nil {
		return nil
	}

	fmt.Printf("xxxxxx %#v\n", c.activeUser)
	if c.activeUser != nil {
		c.service = c.activeUser.New()
	}

	return nil
}

func (c *ConfigsData) initActiveUser() error {
	// 如果已经初始化过，则跳过
	if c.AcitveUID != "" && c.activeUser != nil && c.activeUser.UIDHazy == c.AcitveUID {
		return nil
	}

	if c.AcitveUID == "" && c.activeUser != nil {
		c.AcitveUID = c.activeUser.UIDHazy
		return nil
	}

	if c.AcitveUID != "" {
		for _, user := range c.Users {
			if user.UIDHazy == c.AcitveUID {
				c.activeUser = user
				return nil
			}
		}
	}

	if c.AcitveUID == "" && len(c.Users) == 0 {
		c.activeUser = new(Dedao)
	}

	if len(c.Users) > 0 {
		return errors.New("存在登录的用户，可以进行切换登录用户")
	}

	return errors.New("未登陆")
}

// Save 保存配置
func (c *ConfigsData) Save() error {
	err := c.lazyOpenConfigFile()
	if err != nil {
		fmt.Println(err)
		return err
	}

	c.fileMu.Lock()
	defer c.fileMu.Unlock()

	// 保存配置的数据
	var conf configJSONExport
	conf.AcitveUID = c.AcitveUID
	conf.Users = c.Users

	data, err := jsoniter.MarshalIndent(conf, "", " ")

	if err != nil {
		panic(err)
	}

	// 减掉多余的部分
	err = c.configFile.Truncate(int64(len(data)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = c.configFile.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = c.configFile.Write(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (c *ConfigsData) initDefaultConfig() {
	// todo 默认配置
	c.Save()
}

func (c *ConfigsData) loadConfigFromFile() error {
	err := c.lazyOpenConfigFile()
	if err != nil {
		return err
	}

	info, err := c.configFile.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		return c.Save()
	}

	c.fileMu.Lock()
	defer c.fileMu.Unlock()

	_, err = c.configFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil
	}

	// 从配置文件中加载配置
	decoder := jsoniter.NewDecoder(c.configFile)
	var conf configJSONExport
	decoder.Decode(&conf)

	c.AcitveUID = conf.AcitveUID
	c.Users = conf.Users
	return nil
}

func (c *ConfigsData) lazyOpenConfigFile() (err error) {
	if c.configFile != nil {
		return nil
	}
	c.fileMu.Lock()
	os.MkdirAll(filepath.Dir(c.configFilePath), 0700)
	c.configFile, err = os.OpenFile(c.configFilePath, os.O_CREATE|os.O_RDWR, 0600)
	c.fileMu.Unlock()

	if err != nil {
		if os.IsPermission(err) {
			return
		}
		if os.IsExist(err) {
			return
		}
		return
	}
	return
}

// New config
func New(configFilePath string) *ConfigsData {
	c := &ConfigsData{
		configFilePath: configFilePath,
	}

	return c
}

// GetConfigDir config file dir
func GetConfigDir() string {
	configDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return configDir
}

// ActiveUserService user
func (c *ConfigsData) ActiveUserService() *services.Service {
	if c.service == nil {
		c.service = c.activeUser.New()
	}
	return c.service
}

// SetUser set user
func (c *ConfigsData) SetUser(u *Dedao) (*Dedao, error) {
	ser := services.NewService(&u.CookieOptions)
	user, err := ser.User()
	if err != nil {
		return nil, err
	}

	c.DeleteUser(&User{UIDHazy: user.UIDHazy})

	dedao := &Dedao{
		User: User{
			UIDHazy: user.UIDHazy,
			Name:    user.Nickname,
			Avatar:  user.Avatar,
		},
		CookieOptions: u.CookieOptions,
	}
	c.Users = append(c.Users, dedao)
	c.setActiveUser(dedao)
	return dedao, nil
}

// DeleteUser delete
func (c *ConfigsData) DeleteUser(u *User) {
	for k, user := range c.Users {
		if user.UIDHazy == u.UIDHazy {
			c.Users = append(c.Users[:k], c.Users[k+1:]...)
			break
		}
	}
}

// ActiveUser active user
func (c *ConfigsData) ActiveUser() *Dedao {
	return c.activeUser
}

func (c *ConfigsData) setActiveUser(u *Dedao) {
	c.AcitveUID = u.UIDHazy
	c.activeUser = u
}

// LoginUserCount 登录用户数量
func (c *ConfigsData) LoginUserCount() int {
	return len(c.Users)
}

// SwitchUser switch user
func (c *ConfigsData) SwitchUser(u *User) error {
	for _, user := range c.Users {
		if user.UIDHazy == u.UIDHazy {
			c.setActiveUser(user)
			c.Save()
			// remove cache file
			dir, _ := ioutil.ReadDir("./.cache")
			for _, d := range dir {
				os.RemoveAll(path.Join([]string{".cache", d.Name()}...))
			}
			return nil
		}
	}

	return errors.New("用户不存在")
}
