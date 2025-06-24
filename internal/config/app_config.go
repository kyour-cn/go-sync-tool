package config

type AiConfig struct {
	Enable      bool   `toml:"enable" json:"enable" comment:"是否启用"`
	Gateway     string `toml:"gateway" json:"gateway" comment:"API网关"`
	GatewayType string `toml:"gateway_type" json:"gateway_type" comment:"API网关类型 openai或ollama"`
	Model       string `toml:"model" json:"model" comment:"API模型"`
	Key         string `toml:"key" json:"key" comment:"API密钥"`
}

// AppConfig 应用配置
type AppConfig struct {
	IsDark      bool     `toml:"is_dark" json:"is_dark" comment:"是否为深色模式"`
	AutoStart   bool     `toml:"auto_start" json:"auto_start" comment:"是否自动启动"`
	ProjectName string   `toml:"project_name" json:"project_name" comment:"项目名"`
	ErpEncoding int      `toml:"erp_encoding" json:"erp_encoding" comment:"erp编码 0=gbk 1=uft8 2=智能识别"`
	AIConfig    AiConfig `toml:"ai" json:"ai" comment:"AI配置"`
}

var appConf *AppConfig

// GetAppConfig 获取App配置
func GetAppConfig() (*AppConfig, error) {
	key := "app"

	// 初始化配置
	if appConf == nil {
		_conf := &AppConfig{
			IsDark: true,
		}

		// 如果配置不存在，则创建默认配置
		if !Exists(key) {
			err := SetAppConfig(_conf)
			if err != nil {
				return nil, err
			}
		}

		err := Unmarshal(key, _conf)
		if err != nil {
			return nil, err
		}
		appConf = _conf
	}

	return appConf, nil
}

// SetAppConfig 设置App配置
func SetAppConfig(conf *AppConfig) error {
	key := "app"
	appConf = conf
	return Marshal(key, conf)
}
