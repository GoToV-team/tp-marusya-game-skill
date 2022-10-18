package game

import "github.com/evrone/go-clean-template/pkg/marusia"

type Director interface {
	PlayScene(command SceneRequest) Result
}

type ClosedDirector interface {
	Close()
}

type PlayedSceneResult struct {
	Result
	WorkedDirector ClosedDirector
}

type ScriptRunner interface {
	AttachDirector(sessionId string, op Director)
	RunScene(req marusia.Request) chan PlayedSceneResult
}
