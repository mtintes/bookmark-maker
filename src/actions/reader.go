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

func BuildStructure(inputDirectory []string, outputFile string) {

	links, folders, _ := readInput(inputDirectory)

	globalLinks = links

	//log.Println(links)
	//log.Println(folders)

	bookmarkFile := addHeader()

	bookmarkFile = fmt.Sprintf("%s%s", bookmarkFile, filesToPaths(folders))

	bookmarkFile += addFooter()
	//log.Println("paths")
	//log.Println(gohtml.Format(bookmarkFile))

	writefile := []byte(bookmarkFile)
	err := os.WriteFile(outputFile, writefile, 0644)
	if err != nil {
		log.Println("We had trouble writing the file", err)
	}
}

func readInput(inputDirectory []string) ([]Link, map[string]interface{}, error) {

	/////// MERGE 1 //////////

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
	// log.Println(inputFile.Folders)
	//return &inputFile.Links, inputFile.Folders, nil

	// junk := InputFile{Links: []Link{}, Folders: nil}
	return removeDuplicateLinks(inputFile.Links), inputFile.Folders, nil
}

func merge(input []InputFile) InputFile {
	var response InputFile

	//verify only one of each file

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
	// currentString = fmt.Sprintf("%s/%s", currentString, key)
	for key, dir := range folders {
		//log.Println(reflect.TypeOf(dir))

		// log.Printf("current: %s", currentString)
		// log.Printf("key: %s", key)
		//log.Println(reflect.TypeOf(dir))

		if key == "bookmarks" {
			//log.Println(dir)
			response += generateBookmarks(dir.([]interface{}))
			//log.Println("It is a link")
		} else if key == "name" {
			//log.Println("help")
			//response += fmt.Sprintf("%s/%s", currentString, key)
		} else if _, ok := dir.([]interface{}); ok {
			response += fmt.Sprintf("<DT><H3 ADD_DATE=\"1677903092\" LAST_MODIFIED=\"1677903150\">%s</H3>\n", key)

			response += "<DL><p>\n"
			// log.Println("file")
			// log.Println(dir)
			//currentString = fmt.Sprintf("%s/%s", currentString, dir)

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

	// log.Println("test", filesToBookmark)

	for _, bookmark := range filesToBookmark {
		// fmt.Println("typeof", reflect.TypeOf(bookmark))

		if reflect.TypeOf(bookmark).Kind() == reflect.String {
			// log.Println("Is String")
			// log.Println(bookmark)
			name, url := lookupLinkName(bookmark.(string))
			response += fmt.Sprintf("<DT><A HREF=\"%s\" ADD_DATE=\"1677903092\" LAST_MODIFIED=\"1677903150\">%s</A>\n", url, name)
		} else if reflect.TypeOf(bookmark).Kind() == reflect.Map {
			// log.Println("is map string")
			var name string
			var url string
			for key, value := range bookmark.(map[string]interface{}) {
				// log.Println("here", key, value)
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
	// log.Println("globalLinks", globalLinks)

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
