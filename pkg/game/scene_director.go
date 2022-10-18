package game

import (
	"context"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/evrone/go-clean-template/pkg/marusia"
	"github.com/evrone/go-clean-template/pkg/stack"
	"strings"
)

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
	ctx := req.toSceneContext(so.ctx)

	sceneInfo := scene.Info{}
	if so.currentScene != nil {
		sceneInfo, _ = so.currentScene.GetSceneInfo(ctx)
	}

	errCmd := scene.NoCommand
	var sceneCmd scene.Command

	switch strings.ToLower(req.Command) {
	case marusia.OnStart, "debug":
		so.currentScene = so.cf.StartScene

	case marusia.OnInterrupt, strings.ToLower(so.cf.EndCommand):
		so.stashedScene.Push(so.currentScene)
		sceneCmd = so.cf.GoodbyeScene.React(ctx)
		so.currentScene = so.cf.GoodbyeScene.Next()
		so.reactSceneCommand(sceneCmd)

	default:
		cmd := so.matchCommands(req.Command, sceneInfo.ExpectedMessages)
		if cmd != "" {
			req.Command = cmd

			sceneCmd = so.currentScene.React(ctx)
			so.currentScene = so.currentScene.Next()
			so.reactSceneCommand(sceneCmd)
		} else {
			so.stashedScene.Push(so.currentScene)
			errCmd = so.cf.ErrorScene.React(ctx)
			so.currentScene = so.cf.ErrorScene.Next()
		}
	}

	info := scene.Info{}
	withReact := true
	if so.isEndOfScript {
		info = scene.Info{Text: so.cf.GoodbyeMessage}
	} else {
		info, withReact = so.currentScene.GetSceneInfo(ctx)
		for !withReact {
			so.currentScene = so.currentScene.Next()
			oldInfo := info
			info, withReact = so.currentScene.GetSceneInfo(ctx)
			info = scene.Info{
				Text: scene.Text{
					BaseText:     oldInfo.Text.BaseText + "\n" + info.Text.BaseText,
					TextToSpeech: oldInfo.Text.TextToSpeech + "\n" + info.Text.TextToSpeech,
				},
				Buttons:          info.Buttons,
				ExpectedMessages: info.ExpectedMessages,
			}
		}

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

func (so *ScriptDirector) matchCommands(command string, commands []scene.MessageMatcher) string {
	for _, cmd := range commands {
		if matched, msg := cmd.Match(command); matched {
			return msg
		}
	}
	return ""
}
