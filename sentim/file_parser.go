package sentim

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

// UMICHSI650CorpusParse parse and split the michingan university sentiment data set (https://www.kaggle.com/c/si650winter11/data)
// NOTE: i stripped all double quote marks
func UMICHSI650CorpusParse() {
	baseDir := "build_data/training/"
	baseFile := filepath.Join(baseDir, "base/michigan_train.csv")
	posFileOut := filepath.Join(baseDir, "pos/pos_umichsi650.csv")
	negFileOut := filepath.Join(baseDir, "neg/neg_umichsi650.csv")

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

	posf, err := os.Create(posFileOut)
	if err != nil {
		panic(err)
	}
	defer func() {
		ferr := posf.Close()
		if ferr != nil {
			fmt.Printf("ERROR on defer close file %v %v", posf.Name(), ferr)
		}
	}()

	negf, err := os.Create(negFileOut)
	if err != nil {
		panic(err)
	}
	defer func() {
		ferr := negf.Close()
		if ferr != nil {
			fmt.Printf("ERROR on defer close file %v %v", negf.Name(), ferr)
		}
	}()

	// construct readers/writers
	reader := csv.NewReader(bufio.NewReader(bf))
	reader.Comma = '	'
	reader.LazyQuotes = true
	posWriter := csv.NewWriter(posf)
	negWriter := csv.NewWriter(negf)

	baseData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for i, row := range baseData {
		/*
			if i == 0 {
				continue
			}
		*/
		switch row[0] {
		case "1":
			posWriter.Write([]string{row[1]})
			fmt.Printf("ROW POS %d %v\n", i, row)
		case "0":
			negWriter.Write([]string{row[1]})
			fmt.Printf("ROW NEG %d %v\n", i, row)
		}
	}

	posWriter.Flush()
	negWriter.Flush()

	if err = posWriter.Error(); err != nil {
		fmt.Printf("EOF ERROR: %v", err)
	}

	if err = negWriter.Error(); err != nil {
		fmt.Printf("EOF ERROR: %v", err)
	}
}
