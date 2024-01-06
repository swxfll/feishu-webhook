package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type Payload struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type Grafana struct {
	EvalMatches []struct {
		Value  float64 `json:"value"`
		Metric string  `json:"metric"`
		Tags   struct {
			Name                   string `json:"__name__"`
			App                    string `json:"app"`
			ControllerRevisionHash string `json:"controller_revision_hash"`
			Instance               string `json:"instance"`
			Job                    string `json:"job"`
			KubernetesNamespace    string `json:"kubernetes_namespace"`
			KubernetesPodName      string `json:"kubernetes_pod_name"`
			PodTemplateGeneration  string `json:"pod_template_generation"`
		} `json:"tags"`
	} `json:"evalMatches"`
	Message  string `json:"message"`
	RuleId   int    `json:"ruleId"`
	RuleName string `json:"ruleName"`
	RuleUrl  string `json:"ruleUrl"`
	State    string `json:"state"`
	Title    string `json:"title"`
}

func main() {
	r := gin.Default()

	r.POST("/feishu-webhook", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// 错误处理
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		fmt.Println(string(body)) // 打印原始 JSON 数据到控制台

		var grafana Grafana
		if err := json.Unmarshal(body, &grafana); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
			return
		}

		// 发送 POST 请求
		url := "https://open.feishu.cn/open-apis/bot/v2/hook/9e44a9bf-8952-48ac-beb5-e11f77c25692"

		payload := Payload{
			MsgType: "text",
			Content: struct {
				Text string `json:"text"`
			}(struct{ Text string }{fmt.Sprintf(
				"%s\n实例:%s\n告警项: %s\n告警数值: %f\nPod名称: %s\n所属名称空间: %s\n自定义信息:%s\n",
				grafana.Title,
				grafana.EvalMatches[0].Tags.Instance,
				grafana.EvalMatches[0].Tags.Name,
				grafana.EvalMatches[0].Value,
				grafana.EvalMatches[0].Tags.KubernetesPodName,
				grafana.EvalMatches[0].Tags.KubernetesNamespace,
				grafana.Message)}),
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("JSON 编码失败:", err)
			return
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("POST 请求失败:", err)
			return
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应失败:", err)
			return
		}

		fmt.Println("POST 响应:", string(body))

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
