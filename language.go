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

func (s Suffixes) Len() int {
	return len(s)
}

func (s Suffixes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Suffixes) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func GetLanguage(label string) (language *Language, err error) {
	var path string = os.Getenv("GOBSTRACT_LANGS")
	if path == "" {
		return nil, errors.New("'GOBSTRACT_LANGS' enviroment variable not defined.")
	}

	language = &Language{Label: label, basePath: path}
	if s, e := language.isSupported(); e != nil {
		return nil, e
	} else if !s {
		return nil, errors.New("Language not supported")
	}

	if s, r := language.getDataset("stopwords"); r != nil {
		return nil, r
	} else {
		language.Stopwords = s
	}
	
	if p, e := language.getDataset("prefixes"); e != nil {
		return nil, e
	} else {
		language.Prefixes = p
	}

	if s, e := language.getDataset("suffixes"); e != nil {
		return nil, e
	} else {
		sort.Sort(Suffixes(s))
		language.Suffixes = s
	}

	return language, nil
}

func (l *Language) isSupported() (bool, error) {
	var err error
	var langs []os.FileInfo
	if langs, err = ioutil.ReadDir(l.basePath); err != nil {
		return false, err
	}

	for _, lang := range langs {
		if lang.IsDir() && lang.Name() == l.Label {
			return true, nil
		}
	}
	return false, nil
}

func (l *Language) getDataset(dataset string) (data []string, e error) {
	var loc string = fmt.Sprintf("%s/%s/%s", l.basePath, l.Label, dataset)
	if f, e := os.Open(loc); e != nil {
		return nil, e
	} else {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}

		if e = scanner.Err(); e != nil {
			return nil, e
		}
	}

	return data, nil
}
