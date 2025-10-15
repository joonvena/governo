package rules

type Action struct {
	Type   string         `yaml:"type"`
	Params map[string]any `yaml:"params"`
}

type Rule struct {
	Name      string   `yaml:"name"`
	Service   string   `yaml:"service,omitempty"`
	Methods   []string `yaml:"methods,omitempty"`
	AssetType []string `yaml:"asset_type,omitempty"`
	Condition string   `yaml:"condition"`
	Actions   []Action `yaml:"actions"`
}
