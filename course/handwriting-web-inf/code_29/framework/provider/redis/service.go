/**
 * @Author: vincent
 * @Description:
 * @File:  service
 * @Version: 1.0.0
 * @Date: 2021/11/21 10:25
 */

package redis

import (
	"go-examples/course/handwriting-web-inf/code_29/framework"
	"go-examples/course/handwriting-web-inf/code_29/framework/contract"
	"sync"

	"github.com/go-redis/redis/v8"
)

// HadeRedis  代表hade框架的redis实现
type HadeRedis struct {
	container framework.Container      // 服务容器
	clients   map[string]*redis.Client // key为uniqKey，value为redis.Client（连接池）

	lock *sync.RWMutex // 读写锁
}

// NewHadeRedis 代表实例化Client
func NewHadeRedis(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}

	return &HadeRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

// GetClient 获取Client实例
func (app *HadeRedis) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetBaseConfig(app.container)

	// option对opt进行修改
	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn，就生成dsn
	key := config.UniqKey()

	// 判断是否已经实例化了redis.Client
	app.lock.RLock()
	if db, ok := app.clients[key]; ok {
		// 已经实例化了就返回
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	// 没有实例化，那么就要进行实例化操作
	app.lock.Lock()
	defer app.lock.Unlock()

	// 实例化
	client := redis.NewClient(config.Options)

	// 挂载到map中，结束配置
	app.clients[key] = client

	return client, nil
}
