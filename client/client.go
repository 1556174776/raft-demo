package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"
)

// 测试客户端,发送一堆交易给raft节点集群,获取TPS

type Client struct {
	NodeID string // 客户节点的名称标识符

	NodeURL []string // 所有节点的http url
}

func NewClient(nodeID string, nodeURL []string) *Client {
	return &Client{
		NodeID:  nodeID,
		NodeURL: nodeURL,
	}
}

func (c *Client) CommonWrite(args []string) string {
	// 创建一个缓冲区，用于存储multipart/form-data请求的内容
	var buf bytes.Buffer
	// 创建一个multipart.Writer实例，用于构建multipart/form-data请求的body
	writer := multipart.NewWriter(&buf)
	// 添加form-data参数
	writer.WriteField("contractName", "Common")
	writer.WriteField("functionName", "Write")

	argStr := ""
	for i := 0; i < len(args); i++ {
		argStr += args[i]
		if i != len(args)-1 {
			argStr += " "
		}
	}
	writer.WriteField("args", argStr)

	writer.Close()

	// 随机找一个节点发送
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(c.NodeURL))

	// TODO:或者只向主节点发送(需要先请求一次)
	return sendHttpPost(c.NodeURL[index], buf, writer)
}

func (c *Client) CommonRead(args []string) string {
	// 创建一个缓冲区，用于存储multipart/form-data请求的内容
	var buf bytes.Buffer
	// 创建一个multipart.Writer实例，用于构建multipart/form-data请求的body
	writer := multipart.NewWriter(&buf)
	// 添加form-data参数
	writer.WriteField("contractName", "Common")
	writer.WriteField("functionName", "Read")

	argStr := ""
	for i := 0; i < len(args); i++ {
		argStr += args[i]
		if i != len(args)-1 {
			argStr += " "
		}
	}
	writer.WriteField("args", argStr)

	writer.Close()

	// 随机找一个节点发送
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(c.NodeURL))

	// TODO:或者只向主节点发送(需要先请求一次)
	return sendHttpPost(c.NodeURL[index], buf, writer)
}

func sendHttpPost(dperurl string, buf bytes.Buffer, writer *multipart.Writer) string {
	// 创建一个http请求
	req, err := http.NewRequest("POST", dperurl, &buf)
	if err != nil {
		panic(err)
	}
	// 设置请求头Content-Type为multipart/form-data
	if writer != nil {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 返回响应内容
	return string(body)
}
