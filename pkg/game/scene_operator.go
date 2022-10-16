package game

import (
	"encoding/json"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/evrone/go-clean-template/pkg/marusia"
	"github.com/evrone/go-clean-template/pkg/stack"
	"strings"
)

type Operator interface {
	RunScene(command string, fullUserMessage string, payload json.RawMessage) Answer
}

type Answer struct {
	Text          scene.Text
	Buttons       []scene.Button
	IsEndOfScript bool
}

type SceneOperatorConfig struct {
	StartScene   scene.Scene
	ErrorScene   scene.Scene
	GoodbyeScene scene.Scene
	EndCommand   string
}

type SceneOperator struct {
	stashedScene  stack.Stack[scene.Scene]
	currentScene  scene.Scene
	isEndOfScript bool

	cf SceneOperatorConfig
}

func NewSceneOperator(cf SceneOperatorConfig) *SceneOperator {
	return &SceneOperator{
		stashedScene:  stack.NewStack[scene.Scene](),
		cf:            cf,
		currentScene:  nil,
		isEndOfScript: false,
	}
}

func (so *SceneOperator) RunScene(command string, fullUserMessage string, payload json.RawMessage) Answer {
	errCmd := scene.NoCommand
	switch strings.ToLower(command) {
	case marusia.OnStart, "debug":
		so.currentScene = so.cf.StartScene
		break
	case marusia.OnInterrupt, strings.ToLower(so.cf.EndCommand):
		so.stashedScene.Push(so.currentScene)
		so.currentScene, _ = so.cf.GoodbyeScene.React(command, fullUserMessage, payload)
		break
	default:
		cmd := so.matchCommands(command, so.currentScene.GetSceneInfo().ExpectedMessages)
		if cmd != "" {
			var sceneCmd scene.Command
			so.currentScene, sceneCmd = so.currentScene.React(cmd, fullUserMessage, payload)
			so.reactSceneCommand(sceneCmd)
		} else {
			so.stashedScene.Push(so.currentScene)
			so.currentScene, errCmd = so.cf.ErrorScene.React(command, fullUserMessage, payload)
		}
		break
	}

	info := scene.Info{}
	if so.isEndOfScript {
		info = scene.Info{Text: scene.Text{BaseText: "Пока!", TextToSpeech: "Пока"}}
	} else {
		info = so.currentScene.GetSceneInfo()
		if errCmd == scene.ApplyStashedScene && !so.stashedScene.Empty() {
			so.currentScene, _ = so.stashedScene.Pop()
		}
		info.Buttons = append(info.Buttons, scene.Button{Title: so.cf.EndCommand})
	}
	return Answer{
		Text:          info.Text,
		Buttons:       info.Buttons,
		IsEndOfScript: so.isEndOfScript,
	}
}

func (so *SceneOperator) reactSceneCommand(command scene.Command) {
	switch command {
	case scene.ApplyStashedScene:
		if !so.stashedScene.Empty() {
			so.currentScene, _ = so.stashedScene.Pop()
		}
		break
	case scene.FinishScene:
		so.isEndOfScript = true
		break
	}
}

func (so *SceneOperator) matchCommands(command string, commands []string) string {
	for _, cmd := range commands {
		if cmd == "*" {
			return command
		}

		if strings.Contains(strings.ToLower(command), strings.ToLower(cmd)) {
			return cmd
		}
	}
	return ""
}
