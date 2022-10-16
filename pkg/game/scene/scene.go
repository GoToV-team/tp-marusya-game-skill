package scene

import "encoding/json"

type Command uint64

const (
	NoCommand         = Command(0)
	ApplyStashedScene = Command(1)
	FinishScene       = Command(2)
)

type Button struct {
	Title   string
	URL     string
	Payload interface{}
}

type Text struct {
	BaseText     string
	TextToSpeech string
}

type Info struct {
	Text             Text
	Buttons          []Button
	ExpectedMessages []string
}

type Scene interface {
	GetSceneInfo() Info
	React(message string, fullUserMessage string, payload json.RawMessage) (Scene, Command)
}
