package language

import (
	"os"
	"fmt"
	"sort"
	"errors"
	"regexp"
	"io/ioutil"
	"path/filepath"
)

const langs_dir string = "./language/data"

type Language struct {
	Label string
	Stopwords []string
	Prefixes []string
	Suffixes []string
}

type Suffixes []string

func (sufs Suffixes) Len() int {
	return len(sufs)
}

func (sufs Suffixes) Swap(i, j int) {
	sufs[i], sufs[j] = sufs[j], sufs[i]
}

func (sufs Suffixes) Less(i, j int) bool {
	return len(sufs[i]) > len(sufs[j])
}

func GetLanguage(label string) (*Language, error) {
	var stopwords []string
	var prefixes []string
	var suffixes []string
	var language *Language = &Language{label, stopwords, prefixes, suffixes}

	var err error
	var supported bool
	if supported, err = language.isSupported(); err != nil {
		return nil, err
	}

	if !supported {
		return nil, errors.New("Language not supported")
	}

	if stopwords, err = language.getDataset("stopwords"); err != nil {
		return nil, err
	}
	
	if prefixes, err = language.getDataset("prefixes"); err != nil {
		return nil, err
	}

	if suffixes, err = language.getDataset("suffixes"); err != nil {
		return nil, err
	}

	language.Stopwords = stopwords
	language.Prefixes = prefixes

	sort.Sort(Suffixes(suffixes))
	language.Suffixes = suffixes
	return language, nil
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
		if lang.IsDir() && lang.Name() == language.Label {
			return true, nil
		}
	}

	return false, nil
}

func (language *Language) getDataset(dataset string) ([]string, error) {
	var err error
	var location string = fmt.Sprintf("%s/%s/%s", langs_dir, language.Label, dataset)

	var raw []byte
	if raw, err = ioutil.ReadFile(location); err != nil {
		return nil, err
	}

	var rgx_linebreak *regexp.Regexp = regexp.MustCompile(`\n`)
	return rgx_linebreak.Split(string(raw), -1), nil
}
