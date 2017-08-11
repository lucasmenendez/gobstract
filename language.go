package gobstract

import (
	"os"
	"fmt"
	"sort"
	"bufio"
	"errors"
	"io/ioutil"
	"path/filepath"
)

const langs_dir string = "./language"

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
	var language *Language = &Language{Label: label}


	if supported, err := language.isSupported(); err != nil {
		return nil, err
	} else if !supported {
		return nil, errors.New("Language not supported")
	}

	if stopwords, err := language.getDataset("stopwords"); err != nil {
		return nil, err
	} else {
		language.Stopwords = stopwords
	}
	
	if prefixes, err := language.getDataset("prefixes"); err != nil {
		return nil, err
	} else {
		language.Prefixes = prefixes
	}

	if suffixes, err := language.getDataset("suffixes"); err != nil {
		return nil, err
	} else {
		sort.Sort(Suffixes(suffixes))
		language.Suffixes = suffixes
	}

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
	var data []string
	var location string = fmt.Sprintf("%s/%s/%s", langs_dir, language.Label, dataset)

	if file, err := os.Open(location); err != nil {
		return nil, err
	} else {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		}
	}

	return data, nil
}
