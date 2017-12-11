package smspam

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

var (
	environ = &Environment{}
)

func init() {
	// USE this to build training data (go run script/build_script.go)
	//environ = ParseEnv("production", "config.toml")

	// USE this to run server
	environ = ParseEnv("production", "smspam/config.toml")
}

//GetClassLabelFor return the label information for classification
func GetClassLabelFor(class uint8) (*Label, error) {
	clf := environ.ClfDefs[0]
	return clf.GetLabelInfo(class)
}

// getModel makes the NaiveBayes model used publicly accessible through retrieving it based on environment configuration. It is a method for ClfDef, classifier defintion, which assumes an environment config has been loaded and the correct classifier has been selected from configuration (see getClassifier) . It then uses variously defined configuration to load an the NaiveBayes model.
func (clf *ClfDef) getModel(stream chan base.TextDatapoint) *text.NaiveBayes {
	classCount := uint8(len(clf.Labels))
	if clf.Typ == "only_words_and_numbers" {
		// OnlyWordsAndNumbers is a transform function that will only let 0-1a-zA-Z, and spaces through
		// https://godoc.org/github.com/cdipaolo/goml/base#OnlyWordsAndNumbers
		return text.NewNaiveBayes(stream, classCount, base.OnlyWordsAndNumbers)
	}
	// OnlyWords is a transform function that will only let a-zA-Z, and spaces through
	// https://godoc.org/github.com/cdipaolo/goml/base#OnlyWords
	return text.NewNaiveBayes(stream, classCount, base.OnlyWords)

}

// GetClassifierDefinition retrieves the correct configuration from toml
func (env *Environment) GetClassifierDefinition(m string) *ClfDef {
	for _, c := range env.ClfDefs {
		if c.Name == m {
			return &c
		}
	}
	return &ClfDef{
		Errors: []error{fmt.Errorf("No classifier definition exists for modelName %s in environment: %s", m, env.Name)},
	}
}

// GetFileScanners takes root path for a direcotry and walks it recursively, creating a slice of bufio scanners for each file, returning any errors and the slice of scanners.
func GetFileScanners(rootpath string) (scanners []*bufio.Scanner, funcerr error) {
	funcerr = filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
		}
		if info.IsDir() == true {
			return nil
		}
		f, _ := os.Open(path)
		scanner := bufio.NewScanner(f)
		scanners = append(scanners, scanner)
		return nil
	})
	return
}

// BuildLabelScanners creates bufio file scanners for the training data directory with a WalkFunc to move recursively through the directory. This works fine for directories with small numbers of file, but shoudl be changed if there's lots of files as it will be slow --- recommeded by official docs on WalkFunc: https://golang.org/pkg/path/filepath/#WalkFunc
func (clf *ClfDef) BuildLabelScanners(baseDir, baseTrainDir string) {
	baseTrainPath := path.Join(baseDir, baseTrainDir) // base_dir/training_data

	for i, label := range clf.Labels {
		absPath, _ := filepath.Abs(baseTrainPath + "/" + label.TrainDir) // base_dir/training_data/full_refund
		scns, funcErr := GetFileScanners(absPath)
		if funcErr != nil {
			clf.Errors = append(clf.Errors, fmt.Errorf("error walking label %v directory %v ERROR: %v", label.Name, absPath, funcErr))
		}
		clf.Labels[i].TrainingFileScanners = scns
	}
}

// BuildModel is a method off Environment configuration and uses well-formed config definitions for classifiers and labels to build a NaiveBayes text classifictiaion model. It gets the classifier definition first,ensuring it is well-formed, then builds file scanners for each label by waling the label's training directory for all files. If files exist and things are good we create the cahnnels for the streaming Bayesian model, creating and initializing. Next loop through all scanners, read out text, submit label, stream. wash, rinse, repeat.
func (env *Environment) BuildModel(modelName string) {
	// get the classifier definition to make sure it is well defined
	clf := env.GetClassifierDefinition(modelName)
	clf.panicErrors()

	// build the file scanners on labels to verify files and data is there before defining model
	clf.BuildLabelScanners(env.BaseDir, env.TrainingDir)
	clf.panicErrors()

	streamChan := make(chan base.TextDatapoint, clf.TextChannelSize)
	errorChan := make(chan error)
	model := clf.getModel(streamChan)
	go model.OnlineLearn(errorChan)

	for _, label := range clf.Labels {
		for _, scanner := range label.TrainingFileScanners {
			for scanner.Scan() {
				streamChan <- base.TextDatapoint{
					X: scanner.Text(),
					Y: label.Val,
				}
			}
		}
	}

	close(streamChan)

	for {
		err, _ := <-errorChan
		if err != nil {
			fmt.Printf("Error found: %v", err)
		} else {
			//training is done
			fmt.Printf("Training is DONE! There are '%v' errors", err)
			break
		}
	}

	merr := model.PersistToFile(env.baseModelPath(clf.ModelOut))
	if merr != nil {
		fmt.Println(merr)
	}

}

func (env *Environment) packageModelPath(f string) string {
	//return env.BaseDir + "/" + env.ModelDir + "/" + f
	p, _ := filepath.Abs(path.Join(env.PackageDir, env.BaseDir, env.ModelDir, f))
	return p
}

func (env *Environment) baseModelPath(f string) string {
	//return env.BaseDir + "/" + env.ModelDir + "/" + f
	p, _ := filepath.Abs(path.Join(env.BaseDir, env.ModelDir, f))
	return p
}

func (env *Environment) baseTrainPath(f string) string {
	//return env.BaseDir + "/" + env.TrainingDir + "/" + f
	p, _ := filepath.Abs(path.Join(env.BaseDir, env.TrainingDir, f))
	return p
}

// LoadExistingModel is a simple wrapper around the goml text package's naive bayes RestoreFromFilePath function. It loads a previosly built model from a file, based on configuration definitions. If the config file has changed between model build and trying to load existing model you will need to rebuidl the model before this can be used.
func (env *Environment) LoadExistingModel(m string) *text.NaiveBayes {
	// get the classifier definition to make sure it is well defined
	clf := env.GetClassifierDefinition(m)
	clf.panicErrors()

	streamChan := make(chan base.TextDatapoint, clf.TextChannelSize)
	model := clf.getModel(streamChan)
	err := model.RestoreFromFile(env.baseModelPath(clf.ModelOut))
	if err != nil {
		panic(err)
	}
	return model
}

// SMSSpamModel is a simple wrapper around the goml text package's naive bayes RestoreFromFilePath function. It loads a previosly built model from a file, based on configuration definitions. Similar to LoadExistingModel expect it uses an initialized environment so the model can be loaded by an exernal package ignorant of the config. Specifically, here, the handler package for the text API server.
func SMSSpamModel() *text.NaiveBayes {
	// get the classifier definition to make sure it is well defined
	clf := environ.GetClassifierDefinition("SmsTextSpam")
	clf.panicErrors()

	streamChan := make(chan base.TextDatapoint, clf.TextChannelSize)
	model := clf.getModel(streamChan)
	err := model.RestoreFromFile(environ.packageModelPath(clf.ModelOut))
	if err != nil {
		panic(err)
	}
	return model
}

func (clf *ClfDef) panicErrors() {
	if len(clf.Errors) > 0 {
		for _, e := range clf.Errors {
			log.Print(e)
		}
		panic("END OF ERRORS Panicable")
	}
}
