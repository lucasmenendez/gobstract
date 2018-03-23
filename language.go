package gobstract

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var supported = map[string][]string{
	"es": {"a", "al", "algo", "algunas", "algunos", "ante", "antes", "como",
		"con", "contra", "cual", "cuando", "de", "del", "desde", "donde",
		"durante", "e", "el", "ella", "ellas", "ellos", "en", "entre", "era",
		"erais", "eran", "eras", "eres", "es", "esa", "esas", "ese", "eso",
		"esos", "esta", "estaba", "estabais", "estaban", "estabas", "estad",
		"estada", "estadas", "estado", "estados", "estamos", "estando",
		"estar", "estaremos", "estará", "estarán", "estarás", "estaré",
		"estaréis", "estaría", "estaríais", "estaríamos", "estarían",
		"estarías", "estas", "este", "estemos", "esto", "estos", "estoy",
		"estuve", "estuviera", "estuvierais", "estuvieran", "estuvieras",
		"estuvieron", "estuviese", "estuvieseis", "estuviesen", "estuvieses",
		"estuvimos", "estuviste", "estuvisteis", "estuviéramos",
		"estuviésemos", "estuvo", "está", "estábamos", "estáis", "están",
		"estás", "esté", "estéis", "estén", "estés", "fue", "fuera", "fuerais",
		"fueran", "fueras", "fueron", "fuese", "fueseis", "fuesen", "fueses",
		"fui", "fuimos", "fuiste", "fuisteis", "fuéramos", "fuésemos", "ha",
		"habida", "habidas", "habido", "habidos", "habiendo", "habremos",
		"habrá", "habrán", "habrás", "habré", "habréis", "habría", "habríais",
		"habríamos", "habrían", "habrías", "habéis", "había", "habíais",
		"habíamos", "habían", "habías", "han", "has", "hasta", "hay", "haya",
		"hayamos", "hayan", "hayas", "hayáis", "he", "hemos", "hube", "hubiera",
		"hubierais", "hubieran", "hubieras", "hubieron", "hubiese", "hubieseis",
		"hubiesen", "hubieses", "hubimos", "hubiste", "hubisteis", "hubiéramos",
		"hubiésemos", "hubo", "la", "las", "le", "les", "lo", "los", "me", "mi",
		"mis", "mucho", "muchos", "muy", "más", "mí", "mía", "mías", "mío",
		"míos", "nada", "ni", "no", "nos", "nosotras", "nosotros", "nuestra",
		"nuestras", "nuestro", "nuestros", "order", "os", "otra", "otras",
		"otro", "otros", "para", "pero", "poco", "por", "porque", "que",
		"quien", "quienes", "qué", "se", "sea", "seamos", "sean", "seas",
		"seremos", "será", "serán", "serás", "seré", "seréis", "sería",
		"seríais", "seríamos", "serían", "serías", "seáis", "sido", "siendo",
		"sin", "sobre", "sois", "somos", "son", "soy", "su", "sus", "suya",
		"suyas", "suyo", "suyos", "sí", "también", "tanto", "te", "tendremos",
		"tendrá", "tendrán", "tendrás", "tendré", "tendréis", "tendría",
		"tendríais", "tendríamos", "tendrían", "tendrías", "tened", "tenemos",
		"tenga", "tengamos", "tengan", "tengas", "tengo", "tengáis", "tenida",
		"tenidas", "tenido", "tenidos", "teniendo", "tenéis", "tenía",
		"teníais", "teníamos", "tenían", "tenías", "ti", "tiene", "tienen",
		"tienes", "todo", "todos", "tu", "tus", "tuve", "tuviera", "tuvierais",
		"tuvieran", "tuvieras", "tuvieron", "tuviese", "tuvieseis",
		"tuviesen", "tuvieses", "tuvimos", "tuviste", "tuvisteis", "tuviéramos",
		"tuviésemos", "tuvo", "tuya", "tuyas", "tuyo", "tuyos", "tú", "un",
		"una", "uno", "unos", "vosotras", "vosotros", "vuestra", "vuestras",
		"vuestro", "vuestros", "y", "ya", "yo", "él", "éramos",
	},
	"en": {"a", "about", "above", "after", "again", "against", "all", "am",
		"an", "and", "any", "are", "aren'tokens", "as", "at", "be", "because",
		"been", "before", "being", "below", "between", "both", "but", "by",
		"can'tokens", "cannot", "could", "couldn'tokens", "did", "didn'tokens",
		"do", "does",
		"doesn'tokens", "doing", "don'tokens", "down", "during", "each", "few",
		"for", "from", "further", "had", "hadn'tokens", "has", "hasn'tokens",
		"have", "haven'tokens", "having", "he", "he'd", "he'll", "he'sentences",
		"her", "here", "here'sentences", "hers", "herself", "him", "himself",
		"his", "how", "how'sentences", "i", "i'd", "i'll", "i'm", "i've", "if",
		"in", "into", "is", "isn'tokens", "it", "it'sentences", "its", "itself",
		"let'sentences", "me", "more", "most", "mustn'tokens", "my", "myself",
		"no", "nor", "not", "of", "off", "on", "once", "only", "or", "other",
		"ought", "our", "ours", "ourselves", "out", "over", "own", "same",
		"shan'tokens", "she", "she'd", "she'll", "she'sentences", "should",
		"shouldn'tokens", "so", "some", "such", "than", "that",
		"that'sentences", "the", "their", "theirs", "them", "themselves",
		"then", "there", "there'sentences", "these", "they", "they'd",
		"they'll", "they're", "they've", "this", "those", "through", "to",
		"too", "under", "until", "up", "very", "was", "wasn'tokens", "we",
		"we'd", "we'll", "we're", "we've", "were", "weren'tokens", "what",
		"what'sentences", "when", "when'sentences", "where", "where'sentences",
		"which", "while", "who", "who'sentences", "whom", "why",
		"why'sentences", "with", "won'tokens", "would", "wouldn'tokens", "you",
		"you'd", "you'll", "you're", "you've", "your", "yours", "yourself",
		"yourselves",
	},
}

// Struct to define language object with its code and stopwords
type language struct {
	code      string
	stopwords []string
	model     string
}

// loadLanguage function loads language checking 'STOPWORDS' and 'MODEL'
// environment variables path. Then loads stopwords list from local storage if
// exists or assigns default list. Checks if language model exits.
// Receives language code. Return language struct or error.
func loadLanguage(c string) (l language, e error) {
	var (
		mPath string = os.Getenv("MODELS")
		sPath string = os.Getenv("STOPWORDS")
	)

	l.code = c
	if mPath != "" {
		var fm string = filepath.Join(mPath, c)
		if _, e = os.Stat(fm); e != nil {
			return
		}
		l.model = fm
	}

	if sPath != "" {
		var (
			fs  string = filepath.Join(sPath, c)
			fds *os.File
		)
		if fds, e = os.Open(fs); e != nil {
			return
		}
		defer fds.Close()

		var s *bufio.Scanner = bufio.NewScanner(fds)
		s.Split(bufio.ScanLines)

		for s.Scan() {
			var w string = s.Text()
			if len(w) > 0 {
				l.stopwords = append(l.stopwords, w)
			}
		}
	} else {
		var ok bool
		if l.stopwords, ok = supported[c]; !ok {
			var supported []string
			for _, code := range supported {
				supported = append(supported, code)
			}

			var cs string = strings.Join(supported, ",")
			e = errors.New("language not supported: " + cs)
			return
		}
	}
	return
}

// isStopword function check if provided string is contained into associated
// language stopword list. Previously transform string to lower case to prevent
// false negatives.
func (l language) isStopword(s string) (is bool) {
	var _s string = strings.ToLower(s)
	for _, stw := range l.stopwords {
		is = is || stw == _s
	}
	return
}
