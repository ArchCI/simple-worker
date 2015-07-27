package config

// The struct of worker.yml file.
type WorkerConfig struct {
	Debug          bool   `yaml:"debug"`
	ArchciServer   string `yaml:"archci_server"`
	Interval       int    `yaml:"interval"`
	ConcurrentTask int    `yaml:"concurrent_task"`
	CpuLimit       string `yaml:"cpu_limit"`
	MemoryLimit    string `yaml:"memory_limit"`
}
