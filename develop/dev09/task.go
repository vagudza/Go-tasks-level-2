package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.

Запуск:
go run task.go https://www.iana.org/						(тут много относительных путей) TODO: bug
go run task.go https://golang-blog.blogspot.com/			(тут - абсолютные пути)
*/

// Config - конфигурация проекта
type Config struct {
	url          string
	downloadPath string
	maxDepth     int
	delay        int
}

func NewConfig() *Config {
	var conf Config
	conf.downloadPath = "downloads" //"E:/go/wbschool_exam_L2/develop/dev09/downloads"
	conf.maxDepth = 1
	conf.delay = 0

	if len(os.Args) == 2 {
		conf.url = os.Args[1]
	} else {
		log.Fatalf("wget: missing URL (must be only one)")
	}

	return &conf
}

// Wget - структура, содержащая методы рекурсивного парсинга сайта
type Wget struct {
	visitedLinks               map[string]bool
	conf                       *Config
	regAbsUrl, regHref, regSrc *regexp.Regexp
}

func NewWget(conf *Config) *Wget {
	return &Wget{
		visitedLinks: make(map[string]bool),
		conf:         conf,
		regAbsUrl:    regexp.MustCompile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`),
		regHref:      regexp.MustCompile(`href="\S+\"`),
		regSrc:       regexp.MustCompile(`src="\S+\"`),
	}
}

func main() {
	conf := NewConfig()
	wget := NewWget(conf)
	wget.wget(conf.url, conf.downloadPath, conf.maxDepth)
}

func (w *Wget) wget(sourceUrl string, downloadPath string, maxDepth int) {
	if maxDepth < 0 {
		return
	}

	// валидируем URL
	parsedURI, err := url.ParseRequestURI(sourceUrl)
	if err != nil {
		log.Fatal("validating url: ", err.Error())
	}

	res, err := http.Get(sourceUrl)
	if err != nil {
		log.Fatal("http Get: ", err.Error())
	}
	defer res.Body.Close()

	unescapedUrl, err := url.PathUnescape(sourceUrl)
	if err != nil {
		log.Fatal("url.PathUnescape: ", err.Error())
	}
	log.Printf("Downloading: %s\n", unescapedUrl)
	log.Printf("HTTP request was sended. Waiting for answer... %d %s\n", res.StatusCode, http.StatusText(res.StatusCode))

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll: ", err.Error())
	}

	filename := w.getFilenameWithPath(sourceUrl, parsedURI)
	filepath := downloadPath + filename
	path := w.getPathFromFilepath(filepath)

	isExist, err := w.dirExists(path)
	if err != nil {
		log.Fatal("dirExists: ", err.Error())
	}

	if !isExist {
		err = w.createDir(path)
		if err != nil {
			log.Fatal("createDir: ", err.Error())
		}
	}

	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal("os.Create: ", downloadPath+filename, " --> ", err.Error())
	}
	defer out.Close()

	out.Write(data)
	links := w.getParsedLinks(parsedURI, data)
	for _, innerLink := range links {
		//fmt.Println(">>", innerLink, "<")
		w.wget(innerLink, w.conf.downloadPath, maxDepth-1)
		time.Sleep(time.Millisecond * time.Duration(w.conf.delay))
	}
}

func (w *Wget) getPathFromFilepath(filepath string) string {
	splittedFilepath := strings.Split(filepath, "/")
	return strings.Join(splittedFilepath[:len(splittedFilepath)-1], "/")
}

func (w *Wget) createDir(path string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return fmt.Errorf("can not create folder: '%s'", err.Error())
	}
	return nil
}

// exists returns whether the given file or directory exists
func (w *Wget) dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	// известно ли, что ошибка сообщает о том, что файл или каталог не существуют
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (w *Wget) getFilenameWithPath(sourceUrl string, parsedURI *url.URL) string {
	var filename string

	trimmedPath := strings.TrimRight(parsedURI.Path, "/")
	splittedTrimmedPath := strings.Split(trimmedPath, "#")
	trimmedPath = splittedTrimmedPath[0]

	if trimmedPath == "" {
		filename = "/index.html"
	} else {
		splittedPath := strings.Split(trimmedPath, "/")
		lastPath := splittedPath[len(splittedPath)-1]

		// https://go.dev/pkg/path/ --> /pkg/path.html
		if path.Ext(lastPath) == "" {
			// parsedURI.RawQuery == "m=old#example_Base"
			// https://go.dev/pkg/path/?m=old#example_Base --> /pkg/path_m=old.html
			filename = trimmedPath + w.getQueryParamsString(sourceUrl) + ".html"
		} else {
			filename = trimmedPath
		}
	}

	if filename[:1] != "/" {
		filename = "/" + filename
	}

	filename = w.validateStringToFilepath(filename, false)
	return filename
}

func (w *Wget) getQueryParamsString(sourceUrl string) string {
	if sourceUrl == "" {
		return sourceUrl
	}

	parsedURL, err := url.Parse(sourceUrl)
	if err != nil {
		log.Print("getQueryParamsString", err.Error())
		return ""
	}
	params := parsedURL.Query()

	queryParamsString := ""
	for key, value := range params {
		values := strings.Join(value, "")
		queryParamsString += key + "=" + w.validateStringToFilepath(values, true)
	}
	return queryParamsString
}

func (w *Wget) validateStringToFilepath(str string, isQueryParam bool) string {
	result := str

	if isQueryParam {
		result = strings.ReplaceAll(result, "/", "")
	}

	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, "\\", "")
	result = strings.ReplaceAll(result, ":", "")
	result = strings.ReplaceAll(result, "*", "")
	result = strings.ReplaceAll(result, "?", "")
	result = strings.ReplaceAll(result, "\"", "")
	result = strings.ReplaceAll(result, "<", "")
	result = strings.ReplaceAll(result, ">", "")
	result = strings.ReplaceAll(result, "!", "")
	result = strings.ReplaceAll(result, "+", "")
	return result
}

// getParsedLinks - возвращет все url из data, причем url - уникальны (ранее не обрабатывались) и принадлежат
// домену из исходного url (при запуске программы)
func (w *Wget) getParsedLinks(parentParsedUrl *url.URL, data []byte) []string {
	// регулярка для абсолютных ссылок
	result := w.regAbsUrl.FindAll(data, -1)

	// для относительных ссылок (href)
	hrefs := w.regHref.FindAll(data, -1)
	for _, href := range hrefs {
		// ссылки в href могут быть как относительными, так и абсолютными
		hrefData := string(href[6 : len(href)-1])

		parsedLink, err := parentParsedUrl.Parse(hrefData)
		if err != nil {
			log.Print("getParsedLinks:", err.Error())
			continue
		}

		result = append(result, []byte(parsedLink.String()))
	}

	// для относительных ссылок (src)
	srcs := w.regSrc.FindAll(data, -1)
	for _, src := range srcs {
		hrefData := string(src[5 : len(src)-1])

		parsedLink, err := parentParsedUrl.Parse(hrefData)
		if err != nil {
			log.Print("getParsedLinks:", err.Error())
			continue
		}

		result = append(result, []byte(parsedLink.String()))
	}
	innerUrls := []string{}

	for i := 0; i < len(result); i++ {
		// сохраняем ссылки только для текущего домена
		innerUrl := string(result[i])
		parsedUrl, err := url.Parse(innerUrl)
		if err != nil {
			log.Print("getParsedLinks:", err.Error())
			continue
		}

		// ссылки вида
		// https://www.iana.org/go/rfc3743#section-5.2
		// https://www.iana.org/go/rfc3743#section-6
		// приводим к одному виду без якорей https://www.iana.org/go/rfc3743
		innerUrl = strings.Split(innerUrl, "#")[0]

		if parsedUrl.Host == parentParsedUrl.Host && !w.isUrlExist(innerUrl) {
			innerUrls = append(innerUrls, innerUrl)
		}
	}

	return innerUrls
}

func (w *Wget) isUrlExist(url string) bool {
	var urlExist bool
	_, urlExist = w.visitedLinks[url]

	if urlExist {
		return true
	}

	w.visitedLinks[url] = true
	return false
}
