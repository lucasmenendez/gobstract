package gobstract

import (
	"os"
	"fmt"
	"sort"
	"bufio"
	"errors"
	"io/ioutil"
)

type Language struct {
	Label string
	Stopwords []string
	Prefixes []string
	Suffixes []string
	basePath string
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
	var basePath string = os.Getenv("GOBSTRACT_LANGS")
	if basePath == "" {
		return nil, errors.New("'GOBSTRACT_LANGS' enviroment variable not defined.")
	}

	var language *Language = &Language{Label: label, basePath: basePath}
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
	var langs []os.FileInfo
	if langs, err = ioutil.ReadDir(language.basePath); err != nil {
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
	var location string = fmt.Sprintf("%s/%s/%s", language.basePath, language.Label, dataset)

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
