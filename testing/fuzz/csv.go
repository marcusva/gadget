package fuzz

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	// MaxLinesCSV is the default maximum amount of CSV lines to geenrate
	MaxLinesCSV = 2500
	// MinLinesCSV is the default minimum amount of CSV lines to geenrate
	MinLinesCSV = 0
	// MaxLenStringCSV is the default maximum length of a single CSV column
	MaxLenStringCSV = 50
)

var (
	maxLinesCSV  = MaxLinesCSV
	minLinesCSV  = MinLinesCSV
	maxLenString = MaxLenStringCSV
	csvCharset   = []byte("")
	mu           = sync.Mutex{}
)

func init() {
	// Use a ASCII charset by default, excluding the non-printable characters
	mu.Lock()
	csvCharset = createASCII()
	mu.Unlock()
}

// createASCII populates a byte array with the ASCII charset in the range
// 0x20 (space) to 0x7E (tilde), excluding non-printable characters.
func createASCII() []byte {
	charset := make([]byte, 0x7E-0x20)
	for i := 0x0; i < (0x7E - 0x20); i++ {
		charset[i] = byte(i + 0x20)
	}
	return charset
}

// FuzzedCSV is an io.Rader that contains randomly generated CSV data.
type FuzzedCSV struct {
	io.ReadSeeker

	// Columns contains the column count of the CSV.
	Columns int

	// Lines holds the line count of the generated CSV, excluding the optional
	// header of the CSV.
	Lines int
}

// SetCharset sets the character set to choose from. if charset is nil,
// the ASCII charset, excluding non-printable characters, will be used.
func SetCharset(charset []byte) {
	if charset == nil {
		charset = createASCII()
	}
	mu.Lock()
	csvCharset = charset
	mu.Unlock()
}

// SetLines sets the minimum and maximum number of CSV lines to generate.
// If minlines is smaller than 0, minlines is set to 0. If maxlines is
// smaller than 1, maxlines is set to 0.
func SetLines(min, max int) error {
	if min > max {
		return errors.New("min must be smaller than or equal to max")
	}
	if max <= 0 {
		max = 1
	}
	if min < 0 {
		min = 0
	}
	mu.Lock()
	maxLinesCSV = max
	minLinesCSV = min
	mu.Unlock()
	return nil
}

// SetMaxLenString sets the maximum length of a single string columns.
// If maxlen is smaller than 1, 1 will be set.
func SetMaxLenString(maxlen int) {
	if maxlen <= 0 {
		maxlen = 1
	}
	mu.Lock()
	maxLenString = maxlen
	mu.Unlock()
}

func createRecord(types []string) []string {
	record := make([]string, len(types))
	for idx, t := range types {
		switch t {
		case "int":
			record[idx] = strconv.FormatInt(rand.Int63(), 10)
		case "float":
			record[idx] = strconv.FormatFloat(rand.Float64(), 'e', rand.Intn(24), 64)
		case "bool":
			record[idx] = strconv.FormatBool(rand.Int63n(2) > 0)
		case "string":
			strlen := rand.Intn(maxLenString) + 1
			buf := make([]byte, strlen)
			chMax := len(csvCharset)
			for i := 0; i < strlen; i++ {
				buf[i] = byte(csvCharset[rand.Intn(chMax)])
			}
			record[idx] = string(buf)
		default:
			record[idx] = ""
		}
	}
	return record
}

func validColumnTypes(types []string) error {
	for _, t := range types {
		switch t {
		case "int":
		case "float":
		case "bool":
		case "string":
			continue
		default:
			return fmt.Errorf("invalid column type '%s'", t)
		}
	}
	return nil
}

// CSV returns an in-memory io.Reader containing random CSV data.
func CSV(types []string, delim rune, headers bool) (*FuzzedCSV, error) {
	if err := validColumnTypes(types); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = delim
	switch runtime.GOOS {
	case "windows":
		writer.UseCRLF = true
	default:
		writer.UseCRLF = false
	}

	rand.Seed(time.Now().UnixNano())

	if headers {
		headline := make([]string, len(types))
		for idx, t := range types {
			headline[idx] = fmt.Sprintf("Header [%s]", t)
		}
		if err := writer.Write(headline); err != nil {
			return nil, err
		}
	}
	mu.Lock()
	defer mu.Unlock()
	maxlines := rand.Intn(maxLinesCSV)
	if maxlines < minLinesCSV {
		maxlines = minLinesCSV
	}
	for i := 0; i < maxlines; i++ {
		if err := writer.Write(createRecord(types)); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return &FuzzedCSV{
		ReadSeeker: bytes.NewReader(buf.Bytes()),
		Lines:      maxlines,
		Columns:    len(types),
	}, nil
}
