package types

type (
	Config struct {
		Routing *Routing `yaml:"routing"`
		WebDir  string   `yaml:"web_dir"`
	}

	Routing struct {
		Proxy      bool   `yaml:"proxy"`
		WebURL     string `yaml:"web_url"`
		APIURL     string `yaml:"api_url"`
		GatewayURL string `yaml:"gateway_url"`
	}
)
