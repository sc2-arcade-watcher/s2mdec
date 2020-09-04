package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/icza/s2prot"
	"github.com/sc2-arcade-watcher/s2mdec"
)

var args []string // non-flag args
var compact bool

func init() {
	flag.BoolVar(&compact, "c", false, "print out compact json")
	flag.Parse()
	args = flag.Args()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("fatal error:", r)
		}
	}()
	//
	// len args
	switch len(args) {
	case 1:
		// fileIn
		fileIn, err := os.Open(args[0])
		if err != nil {
			log.Println(err)
			return
		}
		defer fileIn.Close()
		// dataIn
		dataIn, err := ioutil.ReadAll(fileIn)
		if err != nil {
			log.Println(err)
			return
		}
		// switch ext
		switch strings.ToLower(filepath.Ext(fileIn.Name())) {
		case ".s2mi":
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn).ReadStruct().(s2prot.Struct)
			if !ok {
				log.Println("invalid s2mi")
				return
			}
			labeled, err := s2mdec.ReadS2MI(unlabeled)
			if err != nil {
				log.Println("s2mi", err)
				return
			}
			if compact {
				rawJSON, errJSON := json.Marshal(labeled)
				if errJSON != nil {
					log.Println("s2mi", err)
					return
				}
				fmt.Print(string(rawJSON))
			} else {
				fmt.Print(labeled)
			}
		case ".s2mh":
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn).ReadStruct().(s2prot.Struct)
			if !ok {
				log.Println("invalid s2mh")
				return
			}
			labeled, err := s2mdec.ReadS2MH(unlabeled)
			if err != nil {
				log.Println("s2mh:", err)
				return
			}
			if compact {
				rawJSON, errJSON := json.Marshal(labeled)
				if errJSON != nil {
					log.Println("s2mh:", err)
					return
				}
				fmt.Print(string(rawJSON))
			} else {
				fmt.Print(labeled)
			}
		case ".s2ml":
			translation, err := s2mdec.ReadS2ML(dataIn)
			if err != nil {
				log.Println("s2ml:", err)
				return
			}
			if compact {
				rawJSON, errJSON := json.Marshal(translation)
				if errJSON != nil {
					log.Println("s2ml:", err)
					return
				}
				fmt.Print(string(rawJSON))
			} else {
				fmt.Print(translation)
			}
		default:
			log.Println("Unsupported file extension.")
			return
		}
	case 2:
		const nFiles = 2
		fileIn := make([]*os.File, nFiles)
		dataIn := make([][]byte, nFiles)
		for i := 0; i < nFiles; i++ {
			var err error
			// fileIn
			fileIn[i], err = os.Open(args[i])
			if err != nil {
				log.Println(err)
				return
			}
			defer fileIn[i].Close()
			// dataIn
			dataIn[i], err = ioutil.ReadAll(fileIn[i])
			if err != nil {
				log.Println(err)
				return
			}
		}
		// prepare
		s2mh, s2ml := s2prot.Struct(nil), s2mdec.MapLocale(nil)
		// switch ext
		switch strings.ToLower(filepath.Ext(fileIn[0].Name())) {
		case ".s2mh":
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn[0]).ReadStruct().(s2prot.Struct)
			if !ok {
				log.Println("invalid s2mh")
				return
			}
			var err error
			s2mh, err = s2mdec.ReadS2MH(unlabeled)
			if err != nil {
				log.Println("s2mh:", err)
				return
			}
		default:
			log.Println("Unsupported file extension.")
			return
		}
		// switch ext
		switch strings.ToLower(filepath.Ext(fileIn[1].Name())) {
		case ".s2ml":
			var err error
			s2ml, err = s2mdec.ReadS2ML(dataIn[1])
			if err != nil {
				log.Println("s2ml:", err)
				return
			}
		default:
			log.Println("Unsupported file extension.")
			return
		}
		merged, err := s2mdec.S2MHApplyS2ML(s2mh, s2ml, nil)
		if err != nil {
			log.Println("s2mh plus s2ml:", err)
			return
		}
		if compact {
			rawJSON, errJSON := json.Marshal(merged)
			if errJSON != nil {
				log.Println("s2mh plus s2ml:", err)
				return
			}
			fmt.Print(string(rawJSON))
		} else {
			fmt.Print(merged)
		}
	default:
		log.Println("Invalid argument.")
		return
	}
}
