package ai

import (
	"app/internal/config"
	"encoding/json"
	"errors"
	"golang.org/x/exp/slog"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type AddressInfo struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Detail   string `json:"detail"`
}

func FormatAddress(reqAddress string) (*AddressInfo, error) {
	type Request struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Stream bool `json:"stream"`
	}

	type ResponseOllama struct {
		Model     string    `json:"model"`
		CreatedAt time.Time `json:"created_at"`
		Message   struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		DoneReason         string `json:"done_reason"`
		Done               bool   `json:"done"`
		TotalDuration      int64  `json:"total_duration"`
		LoadDuration       int    `json:"load_duration"`
		PromptEvalCount    int    `json:"prompt_eval_count"`
		PromptEvalDuration int64  `json:"prompt_eval_duration"`
		EvalCount          int    `json:"eval_count"`
		EvalDuration       int64  `json:"eval_duration"`
	}

	type ResponseOllamaOpenAI struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			Logprobs     interface{} `json:"logprobs"`
			FinishReason string      `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
			PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
			PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
		} `json:"usage"`
		SystemFingerprint string `json:"system_fingerprint"`
	}

	// 获取AI配置
	appConfig, err := config.GetAppConfig()
	if err != nil {
		return nil, err
	}
	if !appConfig.AIConfig.Enable {
		return nil, errors.New("AI未启用")
	}
	url := appConfig.AIConfig.Gateway

	request := Request{
		Model: appConfig.AIConfig.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role: "system",
				Content: `你的工作是将用户提供的地址按照中国地理行政区划进行解析并拆分为json格式的结果返回，取出地址中的省、市、区县和详细风吹返回格式为：
{"province":"四川省","city":"成都市","district":"高新区","detail":"详细地址"}

需注意几点：
1.必须严格按照三级区划解析，province=省级（如果是直辖市，如重庆市，province和city都为市的名称），city=市级，district=区县级（注意区县也可能有地级市），detail=详细地址（仅返回区县后面的内容）。
2.如果输入地址未从省开始省的，需要分析出该市的省名称，city一定是市结尾的。
3.特别注意直辖市的解析，如重庆市巫溪县城厢镇解放街171号，应该解析为 province=重庆市,city=重庆市,district=巫溪县,detail=解放街171号。
4.返回内容是纯json文本，不要携带` + "```" + `json`,
			},
			{
				Role:    "user",
				Content: reqAddress,
			},
		},
	}

	marshal, _ := json.Marshal(request)
	payload := strings.NewReader(string(marshal))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+appConfig.AIConfig.Key)

	res, err := client.Do(req)
	if err != nil {
		slog.Warn("ai_address", slog.Any("err", err))
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	addressStr := ""

	if appConfig.AIConfig.GatewayType == "ollama" {
		var response ResponseOllama
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		addressStr = response.Message.Content
	} else {
		var response ResponseOllamaOpenAI
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		addressStr = response.Choices[0].Message.Content
	}

	// 移除两侧```
	addressStr = strings.Trim(addressStr, "`")
	addressStr = strings.Trim(addressStr, "json")

	var address AddressInfo
	err = json.Unmarshal([]byte(addressStr), &address)
	if err != nil {
		return nil, err
	}
	return &address, nil
}
