package main

import (
	"syscall/js"
	"github.com/riverzhou/lute"
)

var luteEngine *lute.Lute

func md2html(md string) string {
	return luteEngine.MarkdownStr("", md)
}

func WrapperMd2Html() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            return md2html(args[0].String())
        })
}

func init(){
	luteEngine = lute.New()
}

func main(){
        js.Global().Set("Md2Html", WrapperMd2Html())
	println("Go Assembly Initialed\n")
	<-make(chan bool)
}

