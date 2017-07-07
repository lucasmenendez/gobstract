package language

import (
	"os"
	"fmt"
	"errors"
	"regexp"
	"io/ioutil"
	"path/filepath"
)

var langs_dir = "./language/data"

type Language struct {
	Label string
	Stopwords []string
}

func (language *Language) isSupported() (bool, error) {
	var err error
	var location string

	if location, err = filepath.Abs(langs_dir); err != nil {
		return false, err
	}

	var langs []os.FileInfo
	if langs, err = ioutil.ReadDir(location); err != nil {
		return false, err
	}

	for _, lang := range langs {
		if !lang.IsDir() && lang.Name() == language.Label {
			return true, nil
		}
	}

	return false, nil
}

func (language *Language) getStopwords() error {
	var err error
	var location string = fmt.Sprintf("%s/%s", langs_dir, language.Label)

	var raw_stopwords []byte
	if raw_stopwords, err = ioutil.ReadFile(location); err != nil {
		return err
	}

	var rgx_stopwords *regexp.Regexp = regexp.MustCompile(`\n`)
	language.Stopwords = rgx_stopwords.Split(string(raw_stopwords), -1)
	return nil
}

func GetLanguage(label string) (*Language, error) {
	var stopwords []string
	var language *Language = &Language{label, stopwords}

	var err error
	var supported bool
	if supported, err = language.isSupported(); err != nil {
		return nil, err
	}

	if !supported {
		return nil, errors.New("Language not supported")
	}

	if err = language.getStopwords(); err != nil {
		return nil, err
	}

	return language, nil
}
