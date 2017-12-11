package smspam

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

// UciIrvineMLSpamCorpusParse parse and split the uc irvine sms spam data set (http://archive.ics.uci.edu/ml/datasets/SMS+Spam+Collection#)
func UciIrvineMLSpamCorpusParse() {
	baseDir := "build_data/training/"
	baseFile := filepath.Join(baseDir, "base/SMSSpamCollection")
	spamFileOut := filepath.Join(baseDir, "spam/spam.csv")
	hamFileOut := filepath.Join(baseDir, "not_spam/ham.csv")

	// os all our files; overwirte existing spam/ham files
	bf, err := os.Open(baseFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		ferr := bf.Close()
		if ferr != nil {
			fmt.Printf("ERROR on defer close file %v %v", bf.Name(), ferr)
		}
	}()

	sf, err := os.Create(spamFileOut)
	if err != nil {
		panic(err)
	}
	defer func() {
		ferr := sf.Close()
		if ferr != nil {
			fmt.Printf("ERROR on defer close file %v %v", sf.Name(), ferr)
		}
	}()

	hf, err := os.Create(hamFileOut)
	if err != nil {
		panic(err)
	}
	defer func() {
		ferr := hf.Close()
		if ferr != nil {
			fmt.Printf("ERROR on defer close file %v %v", hf.Name(), ferr)
		}
	}()

	// construct readers/writers
	reader := csv.NewReader(bufio.NewReader(bf))
	reader.Comma = '	'
	reader.LazyQuotes = true
	spamWriter := csv.NewWriter(sf)
	hamWriter := csv.NewWriter(hf)

	baseData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for i, row := range baseData {
		fmt.Printf("ROW %d %v\n", i, row)
		if i == 0 {
			continue
		}
		switch row[0] {
		case "ham":
			hamWriter.Write([]string{row[1]})
		case "spam":
			spamWriter.Write([]string{row[1]})
		}
	}

	hamWriter.Flush()
	spamWriter.Flush()

	if err = hamWriter.Error(); err != nil {
		fmt.Printf("EOF ERROR: %v", err)
	}

	if err = spamWriter.Error(); err != nil {
		fmt.Printf("EOF ERROR: %v", err)
	}
}
