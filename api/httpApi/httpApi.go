package api

import (
	"net/http"
	"time"

	"raftClient/api/httpApi/app"
	errorMsg "raftClient/api/httpApi/error"
	"raftClient/fsm"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
)

func InitRouter(raft *raft.Raft, fsm *fsm.Fsm) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("release")

	r.GET("/set", CommonSet(raft, fsm)) // 简单的set操作(key-value)
	r.GET("/get", CommonGet(raft, fsm)) // 简单的get操作

	r.POST("/newTx", CreateNewTransaction(raft, fsm)) // 根据请求的内容创建一笔交易

	return r
}

// 单次写入
func CommonSet(raft *raft.Raft, fsm *fsm.Fsm) func(*gin.Context) {
	return func(c *gin.Context) {
		appG := app.Gin{c}
		pair := c.Request.URL.Query()
		key := pair.Get("key")
		value := pair.Get("value")
		if key == "" || value == "" {
			appG.Response(http.StatusOK, errorMsg.ERROR, "error key or value")
			return
		}

		data := "set" + "-" + key + "-" + value
		future := raft.Apply([]byte(data), 5*time.Second) // 使用fsm.Apply()完成数据的写入(因为写操作会改变数据库状态)
		if err := future.Error(); err != nil {
			appG.Response(http.StatusOK, errorMsg.ERROR, err.Error())
			return
		}
		appG.Response(http.StatusOK, errorMsg.SUCCESS, "set successfully")
	}
}

// 单次读取
func CommonGet(raft *raft.Raft, fsm *fsm.Fsm) func(*gin.Context) {
	return func(c *gin.Context) {
		appG := app.Gin{c}
		pair := c.Request.URL.Query()
		key := pair.Get("key")
		if key == "" {
			appG.Response(http.StatusOK, errorMsg.ERROR, "error key")
			return
		}
		value := fsm.DataBase.Get(key)
		appG.Response(http.StatusOK, errorMsg.SUCCESS, value)
	}
}

// 生成交易后立即执行,即一轮共识执行一笔交易(TODO:当前无论读写交易都视为raft写操作)
func CreateNewTransaction(raft *raft.Raft, fsm *fsm.Fsm) func(*gin.Context) {
	return func(c *gin.Context) {
		appG := app.Gin{c}
		contractName := c.PostForm("contractName") // 合约名
		functionName := c.PostForm("functionName") // 合约函数名
		args := c.PostForm("args")                 // 函数参数
		valid := validation.Validation{}
		valid.Required(contractName, "contractName").Message("合约名不能为空")
		valid.Required(functionName, "functionName").Message("函数名不能为空")
		valid.Required(args, "args").Message("参数不能为空")

		if valid.HasErrors() {
			app.MarkErrors("CreateNewTransaction", valid.Errors)
			appG.Response(http.StatusOK, errorMsg.INVALID_PARAMS, nil)
			return
		}

		tx := "tx" + "-" + contractName + "-" + functionName + "-" + args

		// TODO:为了提高整体的TPS,可以在获取这一笔交易后,将其添加到等待池,凑够足够数量的交易后再统一处理

		// TODO:如何为了提高读取的TPS，可以在此处检测所有"读"的合约函数,在本地完成读取操作,然后返回

		future := raft.Apply([]byte(tx), 5*time.Second) // 使用fsm.Apply()完成交易的执行(后面的时间是等待raft执行的时间，若超时相当于本次共识失败,返回超时作为错误原因)
		if err := future.Error(); err != nil {
			appG.Response(http.StatusOK, errorMsg.ERROR, err.Error())
			return
		}

		appG.Response(http.StatusOK, errorMsg.SUCCESS, future.Response())
	}
}
