package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// BadgerDB 是对 badger 数据库的封装
type BadgerDB struct {
	db     *badger.DB
	dbPath string
	mu     sync.RWMutex
}

var (
	// 全局 BadgerDB 实例
	badgerInstance *BadgerDB
	once           sync.Once
)

// 确保数据库目录存在
func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// GetBadgerDB 获取全局 BadgerDB 实例
func GetBadgerDB(dbPath string) (*BadgerDB, error) {
	once.Do(func() {
		var err error
		badgerInstance, err = NewBadgerDB(dbPath)
		if err != nil {
			log.Fatalf("初始化 BadgerDB 失败: %v", err)
		}
	})
	return badgerInstance, nil
}

// NewBadgerDB 创建一个新的 BadgerDB 实例
func NewBadgerDB(dbPath string) (*BadgerDB, error) {
	// 确保数据库目录存在
	if err := ensureDir(dbPath); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	options := badger.DefaultOptions(dbPath)
	options.Logger = nil // 禁用日志

	db, err := badger.Open(options)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 启动一个 goroutine 定期执行垃圾回收
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	return &BadgerDB{
		db:     db,
		dbPath: dbPath,
	}, nil
}

// Close 关闭数据库
func (b *BadgerDB) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.db.Close()
}

// Set 存储键值对
func (b *BadgerDB) Set(key string, value interface{}) error {
	return b.SetWithTTL(key, value, 0) // 0 表示永不过期
}

// SetWithTTL 存储键值对，并设置过期时间
func (b *BadgerDB) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 将 value 转换为 JSON 字符串
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), data)
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}
		return txn.SetEntry(entry)
	})
}

// Get 获取键对应的值
func (b *BadgerDB) Get(key string, value interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var data []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			data = append([]byte{}, val...)
			return nil
		})
	})

	if err != nil {
		return err
	}

	// 将 JSON 字符串转换为结构体
	return json.Unmarshal(data, value)
}

// Delete 删除键值对
func (b *BadgerDB) Delete(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// Exists 检查键是否存在
func (b *BadgerDB) Exists(key string) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var exists bool
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			exists = false
			return nil
		}
		if err != nil {
			return err
		}
		exists = true
		return nil
	})

	return exists, err
}

// GetKeysWithPrefix 获取指定前缀的所有键
func (b *BadgerDB) GetKeysWithPrefix(prefix string) ([]string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var keys []string
	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := item.Key()
			keys = append(keys, string(key))
		}
		return nil
	})

	return keys, err
}

// GetAllByPrefix 获取指定前缀的所有键值对
func (b *BadgerDB) GetAllByPrefix(prefix string, valueType interface{}) (map[string]interface{}, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make(map[string]interface{})
	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := string(item.Key())

			// 提取 ID 部分
			idStr := strings.TrimPrefix(key, prefix+":")

			err := item.Value(func(val []byte) error {
				// 解析 JSON 数据到 map
				newValue := make(map[string]interface{})
				if err := json.Unmarshal(val, &newValue); err != nil {
					return err
				}
				result[idStr] = newValue
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err
}

// DeleteWithPrefix 删除指定前缀的所有键值对
func (b *BadgerDB) DeleteWithPrefix(prefix string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 先获取所有匹配的键
	keysToDelete := [][]byte{}
	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			keysToDelete = append(keysToDelete, key)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 然后批量删除
	return b.db.Update(func(txn *badger.Txn) error {
		for _, key := range keysToDelete {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

// FormatKey 格式化键
func FormatKey(category string, id int) string {
	return fmt.Sprintf("%s:%d", category, id)
}

// GetDefaultBadgerDBPath 获取默认的 BadgerDB 路径
func GetDefaultBadgerDBPath() string {
	configDir, err := os.Getwd()
	if err != nil {
		configDir = os.TempDir()
	}
	return filepath.Join(configDir, ".cache", "db")
}
