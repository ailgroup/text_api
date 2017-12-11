package main

import (
	"path/filepath"

	"github.com/jbowles/text_api/smspam"
)

// build the model
func build(trigger bool) {
	if trigger {
		smspam.UciIrvineMLSpamCorpusParse()
	}

	absPath, _ := filepath.Abs("config.toml")
	env := smspam.ParseEnv("production", absPath)
	env.BuildModel("SmsTextSpam")
}

func main() {
	// true re-parses the UC Irvine data set, only need to do that once
	//build(true)
	build(false)
}
