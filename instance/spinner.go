package instance

import (
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) spin(_ discord.Message) {
	if !strings.Contains(in.sdlr.AwaitResumeTrigger(), "use spin") {
		return
	}
