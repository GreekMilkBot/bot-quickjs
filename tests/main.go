package main

import (
	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
	"log"
)

func main() {
	rt := quickjs.NewRuntime()
	defer rt.Close()
	jsCtx := rt.NewContext()
	polyfill.InjectAll(jsCtx)
	output := ""
	defer jsCtx.Close()
	jsCtx.Globals().Set("push", jsCtx.Function(func(context *quickjs.Context, value quickjs.Value, values []quickjs.Value) quickjs.Value {
		output = values[0].String()
		return jsCtx.Undefined()
	}))
	f1, err := jsCtx.Eval(`
fetch('https://api.github.com/repos/go-gitea/gitea/releases/latest').then((e)=>{
    e.json().then(data=>{
        push('gitea 最新版本为: '+data.tag_name)
    })
})
`)
	if err != nil {
		log.Printf("err %v", err)
	}
	f1.Free()
	jsCtx.Loop()
	log.Printf(" %s", output)
}
