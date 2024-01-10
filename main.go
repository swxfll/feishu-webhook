package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

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
			"SendMessage": SendMessage(msg),
		})
	})

	r.POST("/alertmanager-feishu-webhook", func(c *gin.Context) {
		body := PrintAndParseOriginJSON("alertmanager-feishu-webhook", c)

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

		//index := 0
		//var details string
		//for _, item := range alert.Alerts {
		//	index++
		//
		//	annotations := fmt.Sprintf("⚠告警值: %s\n\nℹ︎︎注解: \ndescription: %s\nsummary: %s\n\n",
		//		item.Annotations.Value, item.Annotations.Description, item.Annotations.Summary)
		//
		//	label := fmt.Sprintf("🏷标签: \nalertname: %s\ninstance: %s\njob: %s\nseverity: %s\n\n",
		//		item.Labels.Alertname, item.Labels.Instance, item.Labels.Job, item.Labels.Severity)
		//
		//	details += fmt.Sprintf("====%d====\n", index) + annotations + label
		//
		//}
		//
		//msg := fmt.Sprintf("receiver: %s\nstatus: %s\ngroupLabels: %s\ncommonLabels: %s\n详情(%d 条告警):\n%s",
		//	alert.Receiver,
		//	alert.Status, alert.GroupLabels, alert.CommonLabels,
		//	index, details)

		c.JSON(200, gin.H{
			"SendMessage": SendMessage(string(jsonStr)),
		})

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

func SendMessage(msg string) string {
	// 发送 POST 请求
	url := "https://open.feishu.cn/open-apis/bot/v2/hook/9e44a9bf-8952-48ac-beb5-e11f77c25692"

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
