package actions

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type InputFile struct {
	Links   []Link
	Folders map[string]interface{}
}

type Link struct {
	Id   string
	Name string
	Url  string
}

var globalLinks []Link

func BuildStructure(inputDirectory []string, outputFile string, readmeOutputFile string) {

	links, folders, _ := readInput(inputDirectory)

	globalLinks = links

	bookmarkFile := addHeader()

	bookmarkFile = fmt.Sprintf("%s%s", bookmarkFile, filesToPaths(folders))

	bookmarkFile += addFooter()

	writefile := []byte(bookmarkFile)
	err := os.WriteFile(outputFile, writefile, 0644)
	if err != nil {
		log.Println("We had trouble writing the file", err)
	}

	readme := BuildReadme(folders, 0)
	writefile = []byte(readme)
	err = os.WriteFile(readmeOutputFile, writefile, 0644)
	if err != nil {
		log.Println("We had trouble writing the file", err)
	}
}

func readInput(inputDirectory []string) ([]Link, map[string]interface{}, error) {

	var inputFiles []InputFile

	for _, inputFilePath := range inputDirectory {
		if strings.HasSuffix(inputFilePath, ".yaml") {
			data, err := os.ReadFile(inputFilePath)
			if err != nil {
				log.Println("Bad File Read")
				panic(err)
			}
			var inputFile InputFile
			err = yaml.Unmarshal(data, &inputFile)
			if err != nil {
				log.Println("Bad Marshal")
				panic(err)
			}

			inputFiles = append(inputFiles, inputFile)
		}
	}

	inputFile := merge(inputFiles)

	return removeDuplicateLinks(inputFile.Links), inputFile.Folders, nil
}

func merge(input []InputFile) InputFile {
	var response InputFile

	oneFolders := false

	for _, inputItem := range input {
		if inputItem.Folders != nil && !oneFolders {
			response.Folders = inputItem.Folders
			oneFolders = true
		}

		if inputItem.Links != nil {
			response.Links = append(response.Links, inputItem.Links...)
		}
	}

	return response
}

func filesToPaths(folders map[string]interface{}) string {
	var response string

	for key, dir := range folders {

		if key == "bookmarks" {
			response += generateBookmarks(dir.([]interface{}))
		} else if key == "name" {
		} else if _, ok := dir.([]interface{}); ok {
			response += fmt.Sprintf("<DT><H3 ADD_DATE=\"1677903092\" LAST_MODIFIED=\"1677903150\">%s</H3>\n", key)

			response += "<DL><p>\n"

			response += generateBookmarks(dir.([]interface{}))
			response += "</DL><p>\n"
		} else if _, ok := dir.(map[string]interface{}); ok {
			response += fmt.Sprintf("<DT><H3 ADD_DATE=\"1677903092\" LAST_MODIFIED=\"1677903150\">%s</H3>\n", key)

			response += "<DL><p>\n"
			response += filesToPaths(dir.(map[string]interface{}))
			response += "</DL><p>\n"
		}
	}

	return response
}

func addHeader() string {
	return `<!DOCTYPE NETSCAPE-Bookmark-file-1>
	<!-- This is an automatically generated file.
		 It will be read and overwritten.
		 DO NOT EDIT! -->
	<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
	<TITLE>Bookmarks</TITLE>
	<H1>Bookmarks</H1>
	<DL><p>
    <DT><H3 ADD_DATE="1589210900" LAST_MODIFIED="1589210900" PERSONAL_TOOLBAR_FOLDER="true">Bookmarks bar</H3>
	<DL><p>
	`
	//TODO: get rid of the extra bookmarks bar folder
}

func addFooter() string {
	return "</DL><p>\n</DL><p>\n"
}

func generateBookmarks(filesToBookmark []interface{}) string {
	var response string

	for _, bookmark := range filesToBookmark {

		if reflect.TypeOf(bookmark).Kind() == reflect.String {

			name, url := lookupLinkName(bookmark.(string))
			response += fmt.Sprintf("<DT><A HREF=\"%s\" ADD_DATE=\"1677903092\" LAST_MODIFIED=\"1677903150\">%s</A>\n", url, name)
		} else if reflect.TypeOf(bookmark).Kind() == reflect.Map {

			var name string
			var url string
			for key, value := range bookmark.(map[string]interface{}) {

				if key == "name" {
					name = value.(string)
				} else if key == "url" {
					_, url = lookupLinkName(value.(string))
				}
			}

			response += fmt.Sprintf("<DT><A HREF=\"%s\" ADD_DATE=\"1677903092\" LAST_MODIFIED=\"1677903150\">%s</A>\n", url, name)

		}
	}

	return response
}

func lookupLinkName(linkName string) (string, string) {

	for _, data := range globalLinks {
		if linkName == data.Id {
			return data.Name, data.Url
		}
	}
	return "", ""
}

func removeDuplicateLinks(items []Link) []Link {
	var outlinks []Link

	for _, link := range items {
		found := false
		for _, outlink := range outlinks {
			if link.Id == outlink.Id {
				found = true
				log.Fatalln("you have a duplicate link: ", link.Id)
			}
		}
		if !found {
			outlinks = append(outlinks, link)
		}
	}

	return outlinks

}

func BuildReadme(folders map[string]interface{}, num int) string {
	var response string

	tabs := ""
	for i := 0; i < num; i++ {
		tabs += "\t"
	}

	for key, dir := range folders {
		if key == "bookmarks" {
			response += generateReadme(dir.([]interface{}), num)
		} else if key == "name" {
		} else if _, ok := dir.([]interface{}); ok {

			response += fmt.Sprintf("%s* %s\n\n", tabs, key)
			response += generateReadme(dir.([]interface{}), num+1)
		} else if _, ok := dir.(map[string]interface{}); ok {
			response += fmt.Sprintf("%s* %s\n\n", tabs, key)
			response += BuildReadme(dir.(map[string]interface{}), num+1)
		}
	}

	return response
}

func generateReadme(filesToBookmark []interface{}, num int) string {
	var response string

	tabs := ""
	for i := 0; i < num; i++ {
		tabs += "\t"
	}

	for _, bookmark := range filesToBookmark {
		if reflect.TypeOf(bookmark).Kind() == reflect.String {
			name, url := lookupLinkName(bookmark.(string))
			response += fmt.Sprintf("%s* [%s](%s)\n\n", tabs, name, url)
		} else if reflect.TypeOf(bookmark).Kind() == reflect.Map {
			var name string
			var url string
			for key, value := range bookmark.(map[string]interface{}) {
				if key == "name" {
					name = value.(string)
				} else if key == "url" {
					_, url = lookupLinkName(value.(string))
				}
			}
			response += fmt.Sprintf("%s* [%s](%s)\n\n", tabs, name, url)
		}
	}

	return response
}
