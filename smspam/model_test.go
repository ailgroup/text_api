package smspam

import (
	"testing"
)

var mov int
var env = ParseEnv("production", "config.toml")
var m = env.LoadExistingModel("")
var classDef = env.GetClassifierDefinition("")

func TestLoadModel(t *testing.T) {
	//t.Logf("model words count: class0= %d, class1= %d, class2= %d, class3= %d", m.Count[0], m.Count[1], m.Count[2], m.Count[3])
	t.Logf("model words count: class0= %d, class1= %d", m.Count[0], m.Count[1])
	//t.Log("model dictionary count =", m.DictCount)

	if m.DictCount <= 10 {
		t.Error("model dictionary count should have more that 10 words, count is:", m.DictCount)
	}
}
