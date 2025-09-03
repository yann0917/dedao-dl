package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
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
type CourseIDMap map[int]map[string]interface{}

// ConfigsData Configs data
type ConfigsData struct {
	ActiveUID      string
	DownloadPath   string
	Users          DedaoUsers
	activeUser     *Dedao
	configFilePath string
	configFile     *os.File
	fileMu         sync.Mutex
	service        *services.Service
	badgerDB       *utils.BadgerDB
}

type configJSONExport struct {
	ActiveUID    string
	Users        DedaoUsers
	DownloadPath string
}

// Init 初始化配置
func (c *ConfigsData) Init() error {
	if c.configFilePath == "" {
		return errors.New("配置文件未找到")
	}

	// 初始化 BadgerDB
	badgerDBPath := utils.GetDefaultBadgerDBPath()
	db, err := utils.GetBadgerDB(badgerDBPath)
	if err != nil {
		return fmt.Errorf("初始化 BadgerDB 失败: %w", err)
	}
	c.badgerDB = db

	// 从配置文件中加载配置
	err = c.loadConfigFromFile()
	if err != nil {
		return err
	}

	// 初始化登陆用户信息
	err = c.initActiveUser()
	if err != nil {
		return nil
	}

	if c.activeUser != nil {
		c.service = c.activeUser.New()
	}

	return nil
}

func (c *ConfigsData) initActiveUser() error {
	// 如果已经初始化过，则跳过
	if c.ActiveUID != "" && c.activeUser != nil && c.activeUser.UIDHazy == c.ActiveUID {
		return nil
	}

	if c.ActiveUID == "" && c.activeUser != nil {
		c.ActiveUID = c.activeUser.UIDHazy
		return nil
	}

	if c.ActiveUID != "" {
		for _, user := range c.Users {
			if user.UIDHazy == c.ActiveUID {
				c.activeUser = user
				return nil
			}
		}
	}

	if c.ActiveUID == "" && len(c.Users) == 0 {
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
	conf := configJSONExport{
		ActiveUID:    c.ActiveUID,
		Users:        c.Users,
		DownloadPath: c.DownloadPath,
	}

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
	_ = decoder.Decode(&conf)

	c.ActiveUID = conf.ActiveUID
	c.Users = conf.Users
	c.DownloadPath = conf.DownloadPath
	return nil
}

func (c *ConfigsData) lazyOpenConfigFile() (err error) {
	if c.configFile != nil {
		return nil
	}
	c.fileMu.Lock()
	err = os.MkdirAll(filepath.Dir(c.configFilePath), 0700)
	c.configFile, err = os.OpenFile(c.configFilePath, os.O_CREATE|os.O_RDWR, 0600)
	c.fileMu.Unlock()

	if err != nil {
		if os.IsPermission(err) {
			return
		}
		if os.IsExist(err) {
			return
		}
		if fi, _ := os.Stat(c.configFilePath); fi.IsDir() {
			err = errors.New(c.configFilePath + "配置文件不能是目录")
			log.Fatalln(err)
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
		if c.activeUser == nil {
			c.activeUser = new(Dedao)
		}
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
	c.ActiveUID = u.UIDHazy
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
			err := c.Save()
			return err
		}
	}

	return errors.New("用户不存在")
}

// SetCourseCache 保存课程数据到 Badger 数据库
func (c *ConfigsData) SetCourseCache(category string, id int, course services.Course) error {
	// 遍历 map，将每个键值对存储到 BadgerDB 中
	key := utils.FormatKey(category, id)

	ttl := 2 * time.Hour
	if err := c.badgerDB.SetWithTTL(key, course, ttl); err != nil {
		return fmt.Errorf("保存 %s 数据到 BadgerDB 失败: %w", key, err)
	}
	return nil
}

// GetCourseCache 获取指定类别和 ID 的课程数据
func (c *ConfigsData) GetCourseCache(category string, id int) *services.Course {
	key := utils.FormatKey(category, id)
	var course services.Course
	err := c.badgerDB.Get(key, &course)
	if err != nil {
		// 如果获取失败，返回空对象
		return nil
	}
	return &course
}

// GetAllCourseCache 获取指定类别的所有课程数据
func (c *ConfigsData) GetAllCourseCache(category string) (map[int]*services.Course, error) {
	// 从 Badger 获取所有指定前缀的数据
	rawResults, err := c.badgerDB.GetAllByPrefix(category, nil)
	if err != nil {
		return nil, err
	}

	// 初始化结果 map
	results := make(map[int]*services.Course)

	// 遍历原始结果
	for idStr, rawValue := range rawResults {
		// 将 ID 字符串转为整数
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		// 将原始数据转换为 JSON，再解析为 CourseV2 对象
		rawJson, err := jsoniter.Marshal(rawValue)
		if err != nil {
			continue
		}

		var course services.Course
		if err := jsoniter.Unmarshal(rawJson, &course); err != nil {
			continue
		}

		results[id] = &course
	}

	return results, nil
}

// DeleteCourseCache 删除指定类别和 ID 的课程缓存
func (c *ConfigsData) DeleteCourseCache(category string, id int) error {
	key := utils.FormatKey(category, id)
	return c.badgerDB.Delete(key)
}

// DeleteAllCourseCache 删除指定类别的所有课程缓存
func (c *ConfigsData) DeleteAllCourseCache(category string) error {
	return c.badgerDB.DeleteWithPrefix(category)
}

// ExistsCourseCache 检查指定类别和 ID 的课程缓存是否存在
func (c *ConfigsData) ExistsCourseCache(category string, id int) (bool, error) {
	key := utils.FormatKey(category, id)
	return c.badgerDB.Exists(key)
}
