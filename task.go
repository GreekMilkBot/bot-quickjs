package main

type Task struct {
	Title        string                       `yaml:"title"`
	Describe     string                       `yaml:"describe"`
	Parameters   map[string]interface{}       `yaml:"parameters"`
	Script       string                       `yaml:"script,flow"`
	At           bool                         `yaml:"at"`
	GlobalConfig map[string]string            `yaml:"global_config"`
	ScopeConfig  map[string]map[string]string `yaml:"scope_config"`

	Bytecode []byte `yaml:"-"` // js 可用的 bytecode
}

type TaskParam struct {
}
