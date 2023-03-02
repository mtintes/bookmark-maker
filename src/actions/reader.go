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

var globalLinks *[]Link

func BuildStructure(inputDirectory string, directoryRoot string) {

	links, folders, _ := readInput(inputDirectory)

	globalLinks = links

	log.Println(links)
	log.Println(folders)

	bookmarkFile := addHeader()

	bookmarkFile = fmt.Sprintf("%s%s", bookmarkFile, filesToPaths(folders))

	bookmarkFile += addFooter()
	log.Println("paths")
	log.Println(bookmarkFile)

	// for _, data := range paths {
	// 	log.Println(data)
	// }
}

func readInput(inputDirectory string) (*[]Link, map[string]interface{}, error) {

	var inputFile InputFile

	if strings.HasSuffix(inputDirectory, ".yaml") {
		data, err := os.ReadFile(inputDirectory)
		if err != nil {
			log.Println("Bad File Read")
			panic(err)
		}
		err = yaml.Unmarshal(data, &inputFile)
		if err != nil {
			log.Println("Bad Marshal")
			panic(err)
		}
	}

	// log.Println(inputFile.Folders)
	return &inputFile.Links, inputFile.Folders, nil
}

func filesToPaths(folders map[string]interface{}) string {
	var response string
	// currentString = fmt.Sprintf("%s/%s", currentString, key)
	for key, dir := range folders {
		log.Println(reflect.TypeOf(dir))

		// log.Printf("current: %s", currentString)
		// log.Printf("key: %s", key)
		log.Println(reflect.TypeOf(dir))

		if key == "bookmarks" {
			log.Println(dir)
			response += generateBookmarks(dir.([]interface{}))
			log.Println("It is a link")
		} else if key == "name" {
			log.Println("help")
			//response += fmt.Sprintf("%s/%s", currentString, key)
		} else if _, ok := dir.([]interface{}); ok {
			response += fmt.Sprintf("<DT><H3>%s</H3>", key)

			response += `
			<DL><p>
				`
			// log.Println("file")
			// log.Println(dir)
			//currentString = fmt.Sprintf("%s/%s", currentString, dir)

			response += generateBookmarks(dir.([]interface{}))
			response += `
			</DL><p>`
		} else if _, ok := dir.(map[string]interface{}); ok {
			response += fmt.Sprintf("<DT><H3>%s</H3>", key)

			response += `
			<DL><p>
				`
			response += filesToPaths(dir.(map[string]interface{}))
			response += `
			</DL><p>`
		}
	}

	return response
}

func individualPaths(dir []interface{}) string {
	var response string

	for _, data := range dir {
		log.Println(reflect.TypeOf(data))
		name, url := lookupLinkName(data.(string))
		response += fmt.Sprintf("<DT><A HREF=\"%s\">%s</A>\n", url, name)
		// response += fmt.Sprintf("%s/%s", key, data)
	}

	return response
}

func addHeader() string {
	return `<!DOCTYPE NETSCAPE-Bookmark-file-1>
	<!-- This is an automatically generated file.
		 It will be read and overwritten.
		 DO NOT EDIT! -->
	<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
	<meta http-equiv="Content-Security-Policy"
		  content="default-src 'self'; script-src 'none'; img-src data: *; object-src 'none'"></meta>
	<TITLE>Bookmarks</TITLE>
	<H1>Bookmarks Menu</H1>
	
	<DL><p>
    <DT><H3 ADD_DATE="1589210900" LAST_MODIFIED="1589210900">Bookmarks bar</H3>
	`
	//TODO: get rid of the extra bookmarks bar folder
}

func addFooter() string {
	return `
	</DL>
	`
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
			response += fmt.Sprintf("<DT><A HREF=\"%s\">%s</A>\n", url, name)
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

			response += fmt.Sprintf(`<DT><A HREF="%s">%s</A>\n`, url, name)

		}
	}

	return response
}

func lookupLinkName(linkName string) (string, string) {
	// log.Println("globalLinks", globalLinks)

	for _, data := range *globalLinks {
		if linkName == data.Id {
			return data.Name, data.Url
		}
	}
	return "", ""
}
