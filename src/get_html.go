package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const requestURL = ""
const imageHostPattern = ""

func main() {
	// _. DOM を取得
	doc := getDoc(requestURL)

	// 1. 保存用のディレクトリ用意
	dirName := makeDirectory(doc)

	// 2. 画像 URL の書き出しファイルを生成
	outputFilePath := dirName + "/output.txt"
	file := touchOutputFile(outputFilePath)
	defer file.Close()
	defer os.Remove(outputFilePath)

	// 3. 画像 URL 一覧を保存
	writeURLs(doc, file)

	// _. 読み込みのため，ファイルの先頭に戻す
	file.Seek(0, 0)

	// 4. txt 中の URL から画像ダウンロード
	downloadURLs(file, dirName)
}

func downloadURLs(file *os.File, dirName string) {
	scanner := bufio.NewScanner(file)

	// １行毎に読み込み
	i := 1
	for scanner.Scan() {
		url := scanner.Text()
		filename := dirName + "/" + strconv.Itoa(i)

		fmt.Println(scanner.Text())
		err := downloadFile(url, filename + ".jpg")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Re request by changing extension")

			rep := regexp.MustCompile(`\.(jpg)$`)
			url := rep.ReplaceAllString(url, ".png")

			downloadFile(url, filename + ".png")
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func writeURLs(doc *goquery.Document, file *os.File) {
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("data-src")
		if isMatchHost(url) != true {
			return
		}

		_, _ = file.WriteString(url + "\n")
	})
}

func isMatchHost(str string) (bool) {
    result, _ := regexp.MatchString(imageHostPattern, str)
	return result
}

func touchOutputFile(filePath string) *os.File {
	// 後々読み込むため RDWR
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}

	return file
}

func makeDirectory(doc *goquery.Document) string {
	header := doc.Find("h1").Text()
	dirName := "../results/" + header
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Title: 「%s」\n", header)

	return dirName
}

func getDoc(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		panic("Fail get HTML")
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	return doc
}
