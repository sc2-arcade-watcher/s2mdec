package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/icza/s2prot"
	"github.com/sc2-arcade-watcher/s2mdec"
)

var args []string // non-flag args
var bFlagCompact bool
var bFlagUnlabeled bool

func init() {
	flag.BoolVar(&bFlagCompact, "c", false, "Compact: print out json without indentations")
	flag.BoolVar(&bFlagUnlabeled, "u", false, "Unlabeled: print out json labeled with numbers instead of each field's respective name (applies only to s2mi and s2mh files)")
	flag.Parse()
	args = flag.Args()
}

func run() error {
	// len args
	switch len(args) {
	case 1: // decode a single file
		// fileIn
		fileIn, errFileIn := os.Open(args[0])
		if errFileIn != nil {
			return errFileIn
		}
		defer fileIn.Close()
		// dataIn
		dataIn, errDataIn := ioutil.ReadAll(fileIn)
		if errDataIn != nil {
			return errDataIn
		}
		// switch ext
		switch ext := strings.ToLower(filepath.Ext(fileIn.Name())); ext {
		case ".s2mi":
			// unlabeled
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn).ReadStruct().(s2prot.Struct)
			if !ok {
				return errors.New("invalid s2mi")
			}
			// bFlagUnlabeled
			var output s2prot.Struct
			if bFlagUnlabeled {
				output = unlabeled
			} else {
				labeled, errLabeled := s2mdec.ReadS2MI(unlabeled)
				if errLabeled != nil {
					return fmt.Errorf("s2mi: %v", errLabeled)
				}
				output = labeled
			}
			// bFlagCompact
			if errJSON := writeJSON(os.Stdout, output, !bFlagCompact); errJSON != nil {
				return fmt.Errorf("s2mi: %v", errJSON)
			}
			return nil
		case ".s2mh":
			// unlabeled
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn).ReadStruct().(s2prot.Struct)
			if !ok {
				return errors.New("invalid s2mh")
			}
			// bFlagUnlabeled
			var output s2prot.Struct
			if bFlagUnlabeled {
				output = unlabeled
			} else {
				labeled, errLabeled := s2mdec.ReadS2MH(unlabeled)
				if errLabeled != nil {
					return fmt.Errorf("s2mh: %v", errLabeled)
				}
				output = labeled
			}
			// bFlagCompact
			if errJSON := writeJSON(os.Stdout, output, !bFlagCompact); errJSON != nil {
				return fmt.Errorf("s2mh: %v", errJSON)
			}
			return nil
		case ".s2ml":
			// translation
			translation, errTranslation := s2mdec.ReadS2ML(dataIn)
			if errTranslation != nil {
				return fmt.Errorf("s2ml: %v", errTranslation)
			}
			// bFlagCompact
			if errJSON := writeJSON(os.Stdout, translation, !bFlagCompact); errJSON != nil {
				return fmt.Errorf("s2ml: %v", errJSON)
			}
			return nil
		case ".s2gs":
			// zlib
			rZlib, errZlib := zlib.NewReader(bytes.NewReader(dataIn[16:]))
			if errZlib != nil {
				return fmt.Errorf("s2gs: %v", errZlib)
			}
			defer rZlib.Close()
			// dataIn2
			dataIn2, errDataIn2 := ioutil.ReadAll(rZlib)
			if errDataIn2 != nil {
				return fmt.Errorf("s2gs: %v", errDataIn2)
			}
			// unlabeled
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn2).ReadStruct().(s2prot.Struct)
			if !ok {
				return errors.New("invalid s2gs")
			}
			// bFlagUnlabeled
			var output s2prot.Struct
			if bFlagUnlabeled {
				output = unlabeled
			} else {
				output = unlabeled // not supported yet
			}
			// bFlagCompact
			if errJSON := writeJSON(os.Stdout, output, !bFlagCompact); errJSON != nil {
				return fmt.Errorf("s2gs: %v", errJSON)
			}
			return nil
		default:
			// unlabeled
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn).ReadStruct().(s2prot.Struct)
			if !ok {
				return fmt.Errorf("Unsupported file extension: %v", ext)
			}
			// bFlagCompact
			if errJSON := writeJSON(os.Stdout, unlabeled, !bFlagCompact); errJSON != nil {
				return fmt.Errorf("Unsupported file extension: %v", ext)
			}
			return nil
		}
	case 2: // merging two files s2mh and s2ml
		const nFiles = 2
		fileIn := make([]*os.File, nFiles)
		dataIn := make([][]byte, nFiles)
		for i := 0; i < nFiles; i++ {
			// fileIn
			var errFileIn error
			fileIn[i], errFileIn = os.Open(args[i])
			if errFileIn != nil {
				return errFileIn
			}
			defer fileIn[i].Close()
			// dataIn
			var errDataIn error
			dataIn[i], errDataIn = ioutil.ReadAll(fileIn[i])
			if errDataIn != nil {
				return errDataIn
			}
		}
		// prepare
		s2mh, s2ml := s2prot.Struct(nil), s2mdec.MapLocale(nil)
		// switch ext s2mh
		switch ext := strings.ToLower(filepath.Ext(fileIn[0].Name())); ext {
		case ".s2mh":
			unlabeled, ok := s2mdec.NewVersionedDec(dataIn[0]).ReadStruct().(s2prot.Struct)
			if !ok {
				return errors.New("invalid s2mh")
			}
			var errS2MH error
			s2mh, errS2MH = s2mdec.ReadS2MH(unlabeled)
			if errS2MH != nil {
				return fmt.Errorf("s2mh: %v", errS2MH)
			}
		default:
			return fmt.Errorf("Unsupported file extension: %v", ext)
		}
		// switch ext s2ml
		switch ext := strings.ToLower(filepath.Ext(fileIn[1].Name())); ext {
		case ".s2ml":
			var errS2ML error
			s2ml, errS2ML = s2mdec.ReadS2ML(dataIn[1])
			if errS2ML != nil {
				return fmt.Errorf("s2ml: %v", errS2ML)
			}
		default:
			return fmt.Errorf("Unsupported file extension: %v", ext)
		}
		// merged
		merged, errMerged := s2mdec.S2MHApplyS2ML(s2mh, s2ml, nil)
		if errMerged != nil {
			return fmt.Errorf("s2mh plus s2ml: %v", errMerged)
		}
		// bFlagCompact
		if errJSON := writeJSON(os.Stdout, merged, !bFlagCompact); errJSON != nil {
			return fmt.Errorf("s2mh plus s2ml: %v", errJSON)
		}
		return nil
	default: // unexpected len(args)
		return errors.New("Invalid argument")
	}
}

func writeJSON(w io.Writer, v interface{}, indent bool) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if indent {
		enc.SetIndent("", "  ")
	} else {
		enc.SetIndent("", "")
	}
	return enc.Encode(v)
}

// Exit with the status code.
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln("Unexpected error:", r)
		}
	}()
	if errMain := run(); errMain != nil {
		log.Fatalln(errMain)
	}
}
