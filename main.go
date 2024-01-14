package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

var defaultLastSentTime time.Time = time.Now()
var devopsLastSentTime time.Time = time.Now()
var dbLastSentTime time.Time = time.Now()

func main() {
	r := gin.Default()

	// test
	r.POST("/grafana-feishu-webhook", func(c *gin.Context) {
		body := PrintAndParseOriginJSON("grafana-feishu-webhook", c)

		var grafana Grafana
		if err := json.Unmarshal(body, &grafana); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
			return
		}

		msg := fmt.Sprintf(
			"%s\n实例:%s\n告警项: %s\n告警数值: %f\nPod名称: %s\n所属名称空间: %s\n自定义信息:%s\n",
			grafana.Title,
			grafana.EvalMatches[0].Tags.Instance,
			grafana.EvalMatches[0].Tags.Name,
			grafana.EvalMatches[0].Value,
			grafana.EvalMatches[0].Tags.KubernetesPodName,
			grafana.EvalMatches[0].Tags.KubernetesNamespace,
			grafana.Message)

		c.JSON(200, gin.H{
			"SendMessage": SendMessage("", msg),
		})
	})

	r.POST("/alertmanager-feishu-webhook-default", func(c *gin.Context) {
		body := PrintAndParseOriginJSON("alertmanager-feishu-webhook-default", c)

		var alert AlterManager
		if err := json.Unmarshal(body, &alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON" + err.Error()})
			return
		}

		// 将数据转换为格式化的JSON字符串
		jsonStr, err := json.MarshalIndent(alert, "", "  ")
		if err != nil {
			fmt.Println("JSON formatting error:", err)
			return
		}

		now := time.Now()
		duration := now.Sub(defaultLastSentTime)
		c.JSON(200, gin.H{
			"SendMessage": SendMessage("https://open.feishu.cn/open-apis/bot/v2/hook/9f53885e-2225-4e9e-95de-e8616c2ef7bd",
				string(jsonStr)+"\n"+time.Now().String()+"\n距上一次发送间隔:"+duration.String()),
		})
		defaultLastSentTime = now

	})

	r.POST("/alertmanager-feishu-webhook-devops", func(c *gin.Context) {
		body := PrintAndParseOriginJSON("alertmanager-feishu-webhook-devops", c)

		var alert AlterManager
		if err := json.Unmarshal(body, &alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON" + err.Error()})
			return
		}

		// 将数据转换为格式化的JSON字符串
		jsonStr, err := json.MarshalIndent(alert, "", "  ")
		if err != nil {
			fmt.Println("JSON formatting error:", err)
			return
		}

		now := time.Now()
		duration := now.Sub(devopsLastSentTime)
		c.JSON(200, gin.H{
			"SendMessage": SendMessage("https://open.feishu.cn/open-apis/bot/v2/hook/65c903f4-1283-4725-9ecb-2ad9fd2fd48c",
				string(jsonStr)+"\n"+time.Now().String()+"\n距上一次发送间隔:"+duration.String()),
		})
		devopsLastSentTime = now

	})

	r.POST("/alertmanager-feishu-webhook-db'", func(c *gin.Context) {
		body := PrintAndParseOriginJSON("alertmanager-feishu-webhook-db", c)

		var alert AlterManager
		if err := json.Unmarshal(body, &alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON" + err.Error()})
			return
		}

		// 将数据转换为格式化的JSON字符串
		jsonStr, err := json.MarshalIndent(alert, "", "  ")
		if err != nil {
			fmt.Println("JSON formatting error:", err)
			return
		}

		now := time.Now()
		duration := now.Sub(dbLastSentTime)
		c.JSON(200, gin.H{
			"SendMessage": SendMessage("https://open.feishu.cn/open-apis/bot/v2/hook/bece0a8d-a2ea-4228-8cfa-a437a007b87b",
				string(jsonStr)+"\n"+time.Now().String()+"\n距上一次发送间隔:"+duration.String()),
		})
		dbLastSentTime = now

	})

	r.GET("/hello", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"SendMessage": "Hello World2",
		})

	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func PrintAndParseOriginJSON(route string, c *gin.Context) []byte {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// 错误处理
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return nil
	}

	fmt.Printf("===================[%s] origin json start====================================\n", route)
	fmt.Println(string(body)) // 打印原始 JSON 数据到控制台1
	fmt.Printf("===================[%s] origin json end====================================\n", route)

	return body
}

func SendMessage(url string, msg string) string {
	payload := Payload{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}(struct{ Text string }{msg}),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("SendMessage - JSON 编码失败: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Sprintf("SendMessage - POST 请求失败: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("SendMessage - 读取响应失败: %s", err.Error())
	}

	return string(body)
}
