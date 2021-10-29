package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bogem/id3v2"
	"github.com/jessevdk/go-flags"
	"github.com/seehait-c/bulk-mp3-edit/models"
)

var opts struct {
	MappingFile string `short:"m" long:"mapping-file" description:"Mapping file in JSON" default:"mapping.json"`
	Apply       bool   `short:"a" long:"apply" description:"Apply changes"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error while getting the working directory, %s\n", err)
	}

	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Printf("Error while listing files from %s\n", err)
	}

	mappingFile, err := os.Open(opts.MappingFile)
	if err != nil {
		log.Printf("Error while reading a file %s, %s\n", opts.MappingFile, err)
	}
	mappingByte, err := ioutil.ReadAll(mappingFile)
	if err != nil {
		log.Printf("Error while reading a file %s, %s\n", opts.MappingFile, err)
	}
	var mapping models.Mapping
	json.Unmarshal(mappingByte, &mapping)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp3" {
			mp3File, err := id3v2.Open(file.Name(), id3v2.Options{Parse: true})
			if err != nil {
				log.Printf("Error while reading a file %s, %s\n", file.Name(), err)
			}

			title := mp3File.Title()
			newTitle := title
			for _, nameMapper := range mapping.Name {
				pattern := regexp.MustCompile(nameMapper.Pattern)
				newTitle = pattern.ReplaceAllString(newTitle, nameMapper.Target)
			}
			log.Printf("%24s\t-> %24s", title, newTitle)

			if opts.Apply {
				mp3File.AddTextFrame(mp3File.CommonID("Title"), mp3File.GetTextFrame(mp3File.CommonID("Title")).Encoding, newTitle)
				if err := mp3File.Save(); err != nil {
					log.Printf("Error while applying changes to %s, %s\n", file.Name(), err)
				}
			}
		}
	}
}
