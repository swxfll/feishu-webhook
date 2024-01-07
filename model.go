package main

import "time"

// AlterManager version 0.26.0
type AlterManager struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname              string `json:"alertname"`
			Instance               string `json:"instance"`
			Job                    string `json:"job"`
			Severity               string `json:"severity"`
			App                    string `json:"app,omitempty"`
			ControllerRevisionHash string `json:"controller_revision_hash,omitempty"`
			KubernetesNamespace    string `json:"kubernetes_namespace,omitempty"`
			KubernetesPodName      string `json:"kubernetes_pod_name,omitempty"`
			PodTemplateGeneration  string `json:"pod_template_generation,omitempty"`
		} `json:"labels"`
		Annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		Severity  string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}

type Payload struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// Grafana version 5.3.4
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
