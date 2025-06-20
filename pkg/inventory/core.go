package inventory

type Core struct {
	language string
}

func (core *Core) Init() {
	core.language = "english"
}

func (core *Core) SetLanguage(language string) { core.language = language }
