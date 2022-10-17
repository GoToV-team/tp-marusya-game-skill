package game

import (
	"context"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/evrone/go-clean-template/pkg/marusia"
	"github.com/evrone/go-clean-template/pkg/stack"
	"strings"
)

type Director interface {
	PlayScene(command SceneRequest) Result
}

type ScriptDirector struct {
	stashedScene  stack.Stack[scene.Scene]
	currentScene  scene.Scene
	isEndOfScript bool
	ctx           context.Context
	cf            SceneDirectorConfig
}

func NewScriptDirector(cf SceneDirectorConfig) *ScriptDirector {
	return &ScriptDirector{
		stashedScene:  stack.NewStack[scene.Scene](),
		cf:            cf,
		currentScene:  nil,
		ctx:           context.Background(),
		isEndOfScript: false,
	}
}

func (so *ScriptDirector) PlayScene(req SceneRequest) Result {
	errCmd := scene.NoCommand
	switch strings.ToLower(req.Command) {
	case marusia.OnStart, "debug":
		so.currentScene = so.cf.StartScene
		break
	case marusia.OnInterrupt, strings.ToLower(so.cf.EndCommand):
		so.stashedScene.Push(so.currentScene)
		so.currentScene, _ = so.cf.GoodbyeScene.React(req.toSceneContext(so.ctx))
		break
	default:
		cmd := so.matchCommands(req.Command, so.currentScene.GetSceneInfo().ExpectedMessages)
		if cmd != "" {
			req.Command = cmd
			var sceneCmd scene.Command
			so.currentScene, sceneCmd = so.currentScene.React(req.toSceneContext(so.ctx))
			so.reactSceneCommand(sceneCmd)
		} else {
			so.stashedScene.Push(so.currentScene)
			so.currentScene, errCmd = so.cf.ErrorScene.React(req.toSceneContext(so.ctx))
		}
		break
	}

	info := scene.Info{}
	if so.isEndOfScript {
		info = scene.Info{Text: so.cf.GoodbyeMessage}
	} else {
		info = so.currentScene.GetSceneInfo()
		if errCmd == scene.ApplyStashedScene && !so.stashedScene.Empty() {
			so.currentScene, _ = so.stashedScene.Pop()
		}
		info.Buttons = append(info.Buttons, scene.Button{Title: so.cf.EndCommand})
	}

	return Result{
		Text:          info.Text,
		Buttons:       info.Buttons,
		IsEndOfScript: so.isEndOfScript,
	}
}

func (so *ScriptDirector) reactSceneCommand(command scene.Command) {
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

func (so *ScriptDirector) matchCommands(command string, commands []string) string {
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
