package main

type Config struct {
	BotURL string `yaml:"bot"`

	ToolsDir      string `yaml:"tools"`
	CodegenPrompt string `yaml:"prompt"`

	OpenAIURL   string `yaml:"openai_url"`
	OpenAIToken string `yaml:"openai_token"`

	Tasks []*Task `yaml:"tasks"`
}
