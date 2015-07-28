package config

// ArchciConfig represents the struct of .archci.yml file.
type ArchciConfig struct {
	Image  string   `yaml:"image"`
	Env    []string `yaml:env`
	Script []string `yaml:"script"`
	Email  struct {
		Success []string `yaml:"success"`
		Failure []string `yaml:"failure"`
	} `yaml:"email"`
	Webhook struct {
		Success []string `yaml:success`
		Failure []string `yaml:failure`
	}
}
