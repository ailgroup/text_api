package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/cdipaolo/goml/text"
	"github.com/jbowles/text_api/smspam"
	"github.com/klauspost/cld2"
)

//read the instructions on installing
//https://github.com/klauspost/cld2

var (
	// this will be the template instance handlers uses
	tpl *template.Template
	//this will the sms text spam model loaded in memory so multiple handlers can use it
	smsSpamModel *text.NaiveBayes
)

// setup a templates path globber
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	smsSpamModel = smspam.SMSSpamModel()
}

// Detector holds fields from cld2 project
type Detector struct {
	LanguageCode string
	Accuracy     int
	NormalScore  float64
}

// DetectorResults holds fields from cld2 project and a set of results
type DetectorResults struct {
	DetectResults []*Detector
	Reliable      bool
	NoRelevant    int
}

// DetectBrowse is the handler for the browser friendly endpoint
func DetectBrowse(w http.ResponseWriter, req *http.Request) {
	text := req.FormValue("text")
	predictLangs := cld2.DetectThree(text)
	/*
			DETECT: {Estimates:[{Language:English Percent:98 NormScore:1427}] TextBytes:95 Reliable:true}
		log.Printf("DETECT: %+v", predictLangs)
	*/
	tpl.ExecuteTemplate(w, "index.html", predictLangs)
}

// Detect is the handler for the API endpoint
func Detect(w http.ResponseWriter, req *http.Request) {
	text := req.FormValue("text")
	predict := cld2.DetectThree(text)

	results := DetectorResults{Reliable: predict.Reliable, NoRelevant: predict.TextBytes}
	for _, p := range predict.Estimates {
		d := Detector{
			LanguageCode: p.Language.String(),
			Accuracy:     p.Percent,
			NormalScore:  p.NormScore,
		}
		results.DetectResults = append(results.DetectResults, &d)
	}
	resDetect, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resDetect)
}

// SmsSpamPrediction holds results for text spam classification
type SmsSpamPrediction struct {
	Class            string
	ClassMessage     string
	ClassIndex       uint8
	ProbabilityClass uint8
	Probability      float64
}

// SmsSpamClassify is the handler for classifying SMS text messages as spam or not.
func SmsSpamClassify(w http.ResponseWriter, req *http.Request) {
	msg := req.FormValue("msg")
	predict := smsSpamModel.Predict(msg)
	pclass, probability := smsSpamModel.Probability(msg)
	label, err := smspam.GetClassLabelFor(predict)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
	}

	result := SmsSpamPrediction{
		Class:            label.Name,
		ClassMessage:     label.Msg,
		ClassIndex:       predict,
		ProbabilityClass: pclass,
		Probability:      probability,
	}

	resSpam, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resSpam)
}
