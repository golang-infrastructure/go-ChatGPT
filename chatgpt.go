package chatgpt

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	if_expression "github.com/golang-infrastructure/go-if-expression"
	"github.com/golang-infrastructure/go-pointer"
	"github.com/google/uuid"
	"strings"
)

// ------------------------------------------------- 请求相关 ------------------------------------------------------------

type Request struct {
	Action          string           `json:"action"`
	Messages        []RequestMessage `json:"messages"`
	ConversationID  *string          `json:"conversation_id"`
	ParentMessageID *string          `json:"parent_message_id"`
	Model           string           `json:"model"`
}

type RequestMessage struct {
	ID      string         `json:"id"`
	Role    string         `json:"role"`
	Content RequestContent `json:"content"`
}

type RequestContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

func NewRequest(question, conversationID, parentMessageID string) *Request {
	return &Request{
		Action:         "next",
		ConversationID: pointer.ToPointerOrNil(conversationID),
		Messages: []RequestMessage{
			{
				ID:   uuid.New().String(),
				Role: "user",
				Content: RequestContent{
					ContentType: "text",
					Parts:       []string{question},
				},
			},
		},
		ParentMessageID: pointer.ToPointerOrNil(parentMessageID),
		Model:           "text-davinci-002-render",
	}
}

// ------------------------------------------------- 响应相关 ------------------------------------------------------------

type Response struct {
	Message        ResponseMessage `json:"message"`
	ConversationID string          `json:"conversation_id"`
	Error          any             `json:"error"`
}

type ResponseMessage struct {
	ID         string           `json:"id"`
	Role       string           `json:"role"`
	User       any              `json:"user"`
	CreateTime any              `json:"create_time"`
	UpdateTime any              `json:"update_time"`
	Content    ResponseContent  `json:"content"`
	EndTurn    any              `json:"end_turn"`
	Weight     float64          `json:"weight"`
	Metadata   ResponseMetadata `json:"metadata"`
	Recipient  string           `json:"recipient"`
}

type ResponseContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ResponseMetadata struct {
}

// ------------------------------------------------- --------------------------------------------------------------------

const ConversationAPIURL = "https://chat.openai.com/backend-api/conversation"

type ChatGPT struct {
	jwt           string
	authorization string

	userAgent string

	conversationID  string
	parentMessageID string
}

func NewChat(jwt string) *ChatGPT {
	return &ChatGPT{
		jwt:           jwt,
		authorization: "Bearer " + jwt,
		userAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	}
}

func (x *ChatGPT) Talk(question string) (*Response, error) {
	conversationID := if_expression.Return(x.conversationID != "", x.conversationID, "")
	parentMessageID := if_expression.Return(x.parentMessageID != "", x.parentMessageID, uuid.New().String())
	response, err := x.sendConversationMessage(question, conversationID, parentMessageID)
	if err != nil {
		return nil, err
	}
	x.conversationID = response.ConversationID
	x.parentMessageID = response.Message.ID
	return response, nil
}

func (x *ChatGPT) sendConversationMessage(question, conversationID, parentMessageID string) (*Response, error) {
	request := NewRequest(question, conversationID, parentMessageID)
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	response, err := resty.
		New().
		R().
		SetHeader("User-Agent", x.userAgent).
		SetHeader("Authorization", x.authorization).
		SetHeader("Content-Type", "application/json").
		SetBody(string(requestBytes)).
		Post(ConversationAPIURL)
	if err != nil {
		return nil, err
	}
	responseString := response.String()
	arr := strings.Split(responseString, "\n\n")
	index := len(arr) - 3
	if index >= len(arr) || index < 1 {
		return nil, x.buildErrMessage(responseString)
	}
	arr = strings.Split(arr[index], "data: ")
	if len(arr) < 2 {
		return nil, x.buildErrMessage(responseString)
	}
	chatResponse := &Response{}
	if err = json.Unmarshal([]byte(arr[1]), chatResponse); err != nil {
		return nil, err
	}
	return chatResponse, nil
}

// 在调用失败的时候尝试返回更友好的信息
func (x *ChatGPT) buildErrMessage(responseString string) error {
	errContent := x.tryExtractErrorMessage(responseString)
	if errContent != "" {
		return errors.New(errContent)
	} else {
		return errors.New(responseString)
	}
}

// 尝试从响应中抽取出错误信息
func (x *ChatGPT) tryExtractErrorMessage(responseString string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(responseString))
	if err != nil {
		return ""
	}
	return doc.Find("#content").Text()
}

func (x *ChatGPT) GetConversationID() string {
	return x.conversationID
}

func (x *ChatGPT) SetConversationID(conversationID string) {
	x.conversationID = conversationID
}

func (x *ChatGPT) GetParentMessageID() string {
	return x.parentMessageID
}

func (x *ChatGPT) SetParentMessageID(parentMessageID string) {
	x.parentMessageID = parentMessageID
}

func (x *ChatGPT) GetUserAgent() string {
	return x.userAgent
}

func (x *ChatGPT) SetUserAgent(userAgent string) {
	x.userAgent = userAgent
}

func (x *ChatGPT) SetJWT(jwt string) {
	x.jwt = jwt
	x.authorization = "Bearer " + jwt
}
