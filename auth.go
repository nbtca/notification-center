package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 鉴权失败返回错误
func auth(c *gin.Context, path *string) error {
	// 鉴权
	if cfg.AuthBearer[*path] != "" {
		token := c.GetHeader("Authorization")
		if token != "Bearer "+cfg.AuthBearer[*path] {
			err := fmt.Errorf("auth failed, token not match in path '%s'", *path)
			c.AbortWithError(401, err)
			return err
		}
	} else { //no token for this path 配置中没有设置这个路径的密钥
		err := fmt.Errorf("auth failed, no token. please add token in config.json serverside for path '%s'. example \"auth\":{\"%s\":\"xxx\"}", *path, *path)
		c.AbortWithError(401, err)
		return err
	}
	return nil
}
