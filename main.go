package main

import (
	"./generators"
	"./parser"
	"flag"
	"errors"
	"fmt"
	"strings"
)

var (
	lang  = "csharp"
	v = parser.Package{
		Classes: []parser.Class{
			{
				Name: "Apple",
				Parent: parser.Parent{
					Name:   "Fruit",
					Access: "public",
				},
				Fields: []parser.Field{
					{
						Access:  "public",
						Type:    "string",
						Name:    "colour",
						Default: `"red"`,
					},
					{
						Access:  "public",
						Type:    "string",
						Static:  true,
						Name:    "sort",
						Default: `"Golden"`,
					},
					{
						Access:  "private",
						Type:    "int",
						Name:    "size",
						Default: "1",
					},
				},
				Methods: []parser.Method{
					{
						Access: "private",
						Name:   "print",
						Parameters: []parser.Parameter{
							{
								Pass:  "&",
								Name:  "colour",
								Type:  "string",
								Const: true,
							},
						},
					},
					{
						Access: "protected",
						Return: "int",
						Static: true,
						Name:   "getSize",
					},
					{
						Access: "public",
						Return: "string",
						Name:   "getColor",
						Const:  true,
					},
				},
			},
			{
				Access: "private",
				Name:   "Seed",
				Fields: []parser.Field{
					{
						Access: "public",
						Type:   "int",
						Name:   "size",
					},
				},
				Methods: []parser.Method{
					{
						Static: true,
						Access: "public",
						Return: "int",
						Name:   "transform",
						Const:  true,
					},
				},
			},
		},
	}
)


func getArgs() (string, string, string, bool) {
	langPtr := flag.String("l", "", "language")
	xmlPtr := flag.String("f", "", "file")
	xmlUrlPtr := flag.String("u", "", "file url")
	spacesPtr := flag.Bool("s", false, "use spaces")
	flag.Parse()
	return *langPtr, *xmlPtr, *xmlUrlPtr, *spacesPtr
}

func validateArgs(lang, file, url string) error {
	if lang == "" {
		return errors.New("specify language (-l) flag")
	}
	if url == "" && file == "" {
		return errors.New("specify file path (-f) or url path (-u) flag")
	}
	if file != "" && url != "" {
		return errors.New("do not use both -f and -u flags at the same time")
	}
	return nil
}

func getExtension(language string) string {
	switch language {
	case "java":
		return ".java"
	case "go":
		return ".go"
	case "ruby":
		return ".rb"
	case "cpp":
		return ".h"
	case "python":
		return ".py"
	case "js_es6":
		return ".js"
	case "cs":
		return ".cs"
	}
	return ""
}

func getFileFormat(name string) (string, error) {
	arr := strings.Split(name, ".")
	if len(arr) > 0 {
		return arr[len(arr) - 1], nil
	}
	return "", errors.New("invalid input file")
}

func execute() error {
	language, fileName, url, useSpaces := getArgs()
	err := validateArgs(language, fileName, url)
	if err != nil {
		return err
	}
	var byteContext []byte
	if url != "" {
		byteContext, err = parser.Download(url)
		if err != nil {
			return err
		}
		fileName = url
	} else {
		byteContext, err = parser.Read(fileName)
		if err != nil {
			return err
		}
	}
	generator, err := generators.GetGenerator(language)
	if err != nil {
		return nil
	}
	fileFormat, err := getFileFormat(fileName)
	if err != nil {
		return err
	}
	var object parser.Package
	switch fileFormat {
	case "xml":
		object = parser.ParseXml(byteContext)
	case "json":
		object = parser.ParseJson(byteContext)
	case "yml":
		object = parser.ParseYaml(byteContext)
	default:
		return errors.New(fmt.Sprintf("invalid format of '%s' file", fileName))
	}
	object.UseSpaces = useSpaces
	fileContextMap := generator.Generate(object)
	ext := getExtension(language)
	for key, val := range fileContextMap {
		err = parser.Write(key + ext, val)
		if err != nil {
			return err
		}
	}
	fmt.Println("Generated successfully.")
	return nil
}

func main() {
//	generator, err := generators.GetGenerator(lang)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(generator.Generate(object))

	err := execute()
	if err != nil {
		panic(err)
	}

}
