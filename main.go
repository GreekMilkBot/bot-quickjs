package main

import (
	"context"
	"flag"
	_ "github.com/GreekMilkBot/GreekMilkBot/adapter/standard"
	gmb_bot "github.com/GreekMilkBot/GreekMilkBot/bot"
	"github.com/GreekMilkBot/GreekMilkBot/gmb"
	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

var (
	config = &Config{
		Tasks: make([]*Task, 0),
	}
	configPath = ""
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "config file")
	flag.Parse()
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		def, _ := yaml.Marshal(config)
		_ = os.WriteFile(configPath, def, 0666)
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, config); err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		log.Printf("flush config")
		def, _ := yaml.Marshal(config)
		_ = os.WriteFile(configPath, def, 0666)
	}()

	bot, err := gmb.NewGreekMilkBot(gmb.WithAdapterURL(ctx, config.BotURL))
	if err != nil {
		log.Panic(err)
	}

	rt := quickjs.NewRuntime()
	defer rt.Close()
	bot.HandleMessageFunc(func(ctx context.Context, id string, message gmb_bot.Message) {
		uid := bot.GetMeta(id)["id"]

		content := message.Content
		if item, ok := content[0].(gmb_bot.ContentAt); ok && item.Uid == uid {
			jsCtx := rt.NewContext()
			polyfill.InjectAll(jsCtx)
			output := ""
			jsCtx.Globals().Set("push", jsCtx.Function(func(context *quickjs.Context, value quickjs.Value, values []quickjs.Value) quickjs.Value {
				output = values[0].String()
				return jsCtx.Undefined()
			}))
			defer jsCtx.Close()
			script := content[1].String()
			ret, err := jsCtx.Eval(script)
			if err != nil {
				log.Printf("err %v", err)
			}
			defer ret.Free()
			jsCtx.Loop()
			if output != "" {
				re := make(gmb_bot.Contents, 0)
				re = append(re, gmb_bot.ContentText{Text: output})
				_, err = gmb.NewClientBus(id, bot.ClientCall).SendMessage(&message, &re)
			}

		}
	})
	err = bot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
