package initialize

import (
	"context"
	"fmt"
	"log"

	"go-react-admin/global"

	"github.com/go-redis/redis/v8"
)

// InitRedis 初始化Redis连接
func InitRedis() {
	// 创建Redis客户端
	config := global.GlobalConfig.Redis
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})

	// 测试连接
	ctx := context.Background()
	_, err := global.RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}

	fmt.Println("Redis连接成功")
}
