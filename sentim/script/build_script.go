package main

import (
	"path/filepath"

	"github.com/jbowles/text_api/sentim"
)

// build the model
func build(trigger bool) {
	if trigger {
		sentim.UMICHSI650CorpusParse()
	}

	absPath, _ := filepath.Abs("config.toml")
	env := sentim.ParseEnv("production", absPath)
	env.BuildModel("Sentiment")
}

func main() {
	// true re-parses the UC Irvine data set, only need to do that once
	build(true)
	//build(false)
}
