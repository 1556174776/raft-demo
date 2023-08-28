package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	"github.com/vision9527/raft-demo/fsm"
)

func InitRouter(raft *raft.Raft, fsm *fsm.Fsm) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("release")

	r.GET("/set", CommonSet(raft, fsm)) // 简单的set操作(key-value)
	r.GET("/get", CommonGet(raft, fsm)) // 简单的get操作

	return r
}

// 单次写入
func CommonSet(raft *raft.Raft, fsm *fsm.Fsm) func(*gin.Context) {
	return func(c *gin.Context) {
		pair := c.Request.URL.Query()
		key := pair.Get("key")
		value := pair.Get("value")
		if key == "" || value == "" {
			fmt.Fprintf(c.Writer, "error key or value")
			return
		}

		data := "set" + "," + key + "," + value
		future := raft.Apply([]byte(data), 5*time.Second) // 使用fsm.Apply()完成数据的写入(因为写操作会改变数据库状态)
		if err := future.Error(); err != nil {
			fmt.Fprintf(c.Writer, "error:"+err.Error())
			return
		}
		fmt.Fprintf(c.Writer, "ok")
	}
}

// 单次读取
func CommonGet(raft *raft.Raft, fsm *fsm.Fsm) func(*gin.Context) {
	return func(c *gin.Context) {
		pair := c.Request.URL.Query()
		key := pair.Get("key")
		if key == "" {
			fmt.Fprintf(c.Writer, "error key")
			return
		}
		value := fsm.DataBase.Get(key)
		fmt.Fprintf(c.Writer, value)
	}
}
