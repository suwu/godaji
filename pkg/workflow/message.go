package workflow

import "fmt"

type MessageHandler func(*Message)

func defaultHandler(message *Message) {
	return
}

var handlerRegistry = make(map[string]MessageHandler)

func RegisterHandler(message string, handler MessageHandler) {
	handlerRegistry[message] = handler
}

func GetHandler(message string) MessageHandler {
	handler, ok := handlerRegistry[message]
	if ok {
		return handler
	} else {
		return defaultHandler
	}
}

const (
	CREATE_FLOW string = "CREATE_FLOW"
)

// CREATE_FLOW
func CreateFlowMessage(flowId string) {
	message := &Message{
		Name:    CREATE_FLOW,
		Payload: map[string]interface{}{"flowId": flowId},
	}
	Dispatch(message)
}

func createFlow(message *Message) {
	fmt.Println(message)

}

func init() {
	RegisterHandler(CREATE_FLOW, createFlow)
}
