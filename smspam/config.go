package smspam

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/naoina/toml"
)

type tomlConfig struct {
	Title        string
	Description  string
	Environments []Environment
	Owner        struct {
		Name string
	}
}

// Environment creds and conf and all
// Name is the name of env, BaseDir is for reading and writing data, ClfDef meta-data about classifiers, Client is a scheduler and job runner, Error or any errors.
type Environment struct {
	Name        string
	Version     float32
	BaseDir     string
	PackageDir  string
	ModelDir    string
	TrainingDir string
	ClfDefs     []ClfDef
	Client
	Error error
}

// ClfDef is the config-driven difinition of the classifiers to be created
type ClfDef struct {
	Name            string
	Typ             string
	ModelOut        string
	TextChannelSize int
	Version         float32
	Labels          []Label
	Errors          []error
}

// Label is the config driven definition of each classifier label
type Label struct {
	Val                  uint8
	Name                 string
	TrainDir             string
	Msg                  string
	TrainingFileScanners []*bufio.Scanner
}

// Client controls scheduling and building of classifiers
type Client struct {
	SyncDuration time.Duration
	WaitDuration time.Duration
}

// GetLabelInfo is helper to get the name of the class based on label
func (clf *ClfDef) GetLabelInfo(class uint8) (*Label, error) {
	for _, l := range clf.Labels {
		if l.Val == class {
			return &l, nil
		}
	}
	return &Label{}, fmt.Errorf("class label '%d' does not exist", class)
}

// ParseEnv parses the toml file, loads it, and returns the Environment struct to use
func ParseEnv(env, file string) *Environment {
	f, ferr := os.Open(file)
	if ferr != nil {
		log.Fatal("OPEN CONFIG FILE ERROR:", ferr)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Print("DEFER CLOSE CONFIG FILE ERROR:", err)
		}
	}()

	buf, rerr := ioutil.ReadAll(f)
	if rerr != nil {
		log.Fatal("READ CONFIG FILE ERROR:", rerr)
	}
	var config tomlConfig
	if merr := toml.Unmarshal(buf, &config); merr != nil {
		log.Fatal("UNMARSHAL CONFIG FILE ERROR:", merr)
	}

	//return correct environment
	for _, e := range config.Environments {
		if e.Name == env {
			return &e
		}
	}
	return &Environment{
		//Error: fmt.Errorf(fmt.Sprintf("No configuration in file %s for environment: %s", file, env)),
		Error: fmt.Errorf("No configuration in file %s for environment: %s", file, env),
	}
}
