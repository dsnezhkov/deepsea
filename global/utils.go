package global

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/lucasjones/reggen"
	"io"
	"math/rand"
	"os"
	"regexp"
	"time"
)

var EmailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var CSVFileRe = regexp.MustCompile(`.*\.csv$`)
var DBFileRe = regexp.MustCompile(`.*\.db$`)

func CSV2Json(mfile string) ([]byte, error) {

	var marks []Mark

	csvFile, err := os.Open(mfile)
	if err != nil {
		fmt.Printf("ERROR: CSV File: %s %v\n", mfile, err)
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	_, err = reader.Read() // ignore first line (header)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		marks = append(marks, Mark{
			Identifier: line[0],
			Email:      line[1],
			Firstname:  line[2],
			Lastname:   line[3],
			//Metadata: &global.Metadata{ State: line[4], },
		})
	}
	marksJson, _ := json.Marshal(marks)
	return marksJson, nil
}

// Exists reports whether the named file exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RegToString(regPattern string) (string, error) {

	ident, err := reggen.Generate(regPattern, 1)
	if err != nil {
		return "", err
	}
	return ident, nil
}

func RandString(len int) string {

	var s string
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		s += fmt.Sprintf("%d", rand.Intn(9))
	}

	return s
}
