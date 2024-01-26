package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 鉴权失败返回错误
func auth(c *gin.Context, body *[]byte, path *string) error {
	// 鉴权
	verifyToken := cfg.Auth[*path]
	if verifyToken != "" {
		bearerOrHash := c.GetHeader("Authorization")
		if bearerOrHash == "" && body != nil { //hash method
			bearerOrHash = c.GetHeader("X-Hub-Signature-256")
			if bearerOrHash == "" {
				bearerOrHash = c.GetHeader("X-Signature-256")
			}
			if bearerOrHash == "" {
				err := fmt.Errorf("auth failed, no token from client in path '%s'", *path)
				c.AbortWithError(401, err)
				return err
			}

			if !checkHash(&bearerOrHash, body, &verifyToken) {
				err := fmt.Errorf("auth failed, hash not match in path '%s'", *path)
				c.AbortWithError(401, err)
				return err
			}
		} else if bearerOrHash != "Bearer "+verifyToken {
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

func checkHash(sha256Str *string, body *[]byte, verifyToken *string) bool {
	mac := hmac.New(sha256.New, []byte(*verifyToken))
	mac.Write(*body)
	expectedMAC := mac.Sum(nil)
	expectedHash := hex.EncodeToString(expectedMAC)
	return *sha256Str == "sha256="+expectedHash
}
