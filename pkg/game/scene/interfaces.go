package scene

type Scene interface {
	GetSceneInfo(ctx *Context) (sceneInfo Info, withReact bool)
	React(ctx *Context) Command
	Next() Scene
}

type MessageMatcher interface {
	Match(message string) (isMatch bool, searchedString string)
}
