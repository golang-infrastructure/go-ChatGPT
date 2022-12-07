package chatgpt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFoo(t *testing.T) {

	// 把JWT Token放到这里
	jwt := "xxxx"

	chat := NewChatGPT(jwt)
	talk, err := chat.Talk("你好，我的名字叫陈二！")
	assert.Nil(t, err)
	t.Log(talk.Message.Content)

	talk, err = chat.Talk("我的名字叫什么呀？")
	assert.Nil(t, err)
	t.Log(talk.Message.Content)

	// Output:
	//     chatgpt_test.go:15: {text [你好陈二，很高兴认识你！]}
	//    chatgpt_test.go:19: {text [你已经告诉我你的名字是陈二了。你可以直接说：“我的名字叫陈二。”]}

}
