package main

import (
	"fmt"
	chatgpt "github.com/golang-infrastructure/go-ChatGPT"
)

func main() {

	// 把JWT放到这里
	jwt := "xxx"

	chat := chatgpt.NewChat(jwt)
	talk, err := chat.Talk("你好，我的名字叫陈二！")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(talk.Message.Content)

	talk, err = chat.Talk("我的名字叫什么呀？")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(talk.Message.Content)

	// Output:
	// {text [你好，陈二！很高兴认识你。我是 Assistant，一个大型语言模型，旨在帮助人们了解更多关于世界的信息。如果你有什么问题，我将尽力回答。]}
	// {text [你告诉我你的名字叫陈二。你的名字是一个很普通的中国姓氏和名字，在中国文化中，许多人都有一个姓氏和一个名字。例如，陈是一个常见的姓氏，而二是一个常见的名字。不过，每个人的名字都是独一无二的，所以你的名字只属于你自己。]}

}
