package main

//extern char *s2mdec_read_s2mi(char *buf, int size);
//extern char *s2mdec_read_s2mh(char *buf, int size);
//extern char *s2mdec_read_s2ml(char *buf, int size);
//extern char *s2mdec_read_s2mh_s2ml(char *bufS2MH, int sizeS2MH, char *bufS2ML, int sizeS2ML);
import "C"
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"unsafe"

	"github.com/icza/s2prot"
	"github.com/sc2-arcade-watcher/s2mdec"
)

// This function takes bitpacked bytes and returns a raw JSON in C-type string.
// Note: The returned C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be released
// from memory, such as by calling free() of "stdlib.h".
//export s2mdec_read_s2mi
func s2mdec_read_s2mi(buf *C.char, size C.int) (ret *C.char) {
	defer func() {
		if r := recover(); r != nil {
			ret = C.CString(fmt.Sprint("decoding error: ", r))
		}
	}()
	//
	unlabeled, ok := s2mdec.NewVersionedDec(C.GoBytes(unsafe.Pointer(buf), size)).ReadStruct().(s2prot.Struct)
	if !ok {
		return C.CString("invalid s2mi")
	}
	labeled, err := s2mdec.ReadS2MI(unlabeled)
	if err != nil {
		return C.CString(fmt.Sprint("s2mi: ", err))
	}
	sb := strings.Builder{}
	if errJSON := writeJSON(&sb, labeled); errJSON != nil {
		log.Println("s2mi:", errJSON)
		return
	}
	return C.CString(sb.String())
}

// This function takes bitpacked bytes and returns a raw JSON in C-type string.
// Note: The returned C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be released
// from memory, such as by calling free() of "stdlib.h".
//export s2mdec_read_s2mh
func s2mdec_read_s2mh(buf *C.char, size C.int) (ret *C.char) {
	defer func() {
		if r := recover(); r != nil {
			ret = C.CString(fmt.Sprint("decoding error: ", r))
		}
	}()
	//
	unlabeled, ok := s2mdec.NewVersionedDec(C.GoBytes(unsafe.Pointer(buf), size)).ReadStruct().(s2prot.Struct)
	if !ok {
		return C.CString("invalid s2mh")
	}
	labeled, err := s2mdec.ReadS2MH(unlabeled)
	if err != nil {
		return C.CString(fmt.Sprint("s2mh: ", err))
	}
	sb := strings.Builder{}
	if errJSON := writeJSON(&sb, labeled); errJSON != nil {
		log.Println("s2mh:", errJSON)
		return
	}
	return C.CString(sb.String())
}

// This function takes bitpacked bytes and returns a raw JSON in C-type string.
// Note: The returned C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be released
// from memory, such as by calling free() of "stdlib.h".
//export s2mdec_read_s2ml
func s2mdec_read_s2ml(buf *C.char, size C.int) (ret *C.char) {
	defer func() {
		if r := recover(); r != nil {
			ret = C.CString(fmt.Sprint("decoding error: ", r))
		}
	}()
	//
	translation, err := s2mdec.ReadS2ML(C.GoBytes(unsafe.Pointer(buf), size))
	if err != nil {
		return C.CString(fmt.Sprint("s2ml: ", err))
	}
	sb := strings.Builder{}
	if errJSON := writeJSON(&sb, translation); errJSON != nil {
		log.Println("s2ml:", errJSON)
		return
	}
	return C.CString(sb.String())
}

// This function takes bitpacked bytes and returns a raw JSON.
// Note: The returned C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be released
// from memory, such as by calling free() of "stdlib.h".
//export s2mdec_read_s2mh_s2ml
func s2mdec_read_s2mh_s2ml(bufS2MH *C.char, sizeS2MH C.int, bufS2ML *C.char, sizeS2ML C.int) (ret *C.char) {
	defer func() {
		if r := recover(); r != nil {
			ret = C.CString(fmt.Sprint("decoding error: ", r))
		}
	}()
	//
	s2mh, s2ml := s2prot.Struct(nil), s2mdec.MapLocale(nil)
	{ // s2mh
		unlabeled, ok := s2mdec.NewVersionedDec(C.GoBytes(unsafe.Pointer(bufS2MH), sizeS2MH)).ReadStruct().(s2prot.Struct)
		if !ok {
			return C.CString("invalid s2mh")
		}
		var err error
		s2mh, err = s2mdec.ReadS2MH(unlabeled)
		if err != nil {
			return C.CString(fmt.Sprint("s2mh: ", err))
		}
	}
	{ // s2ml
		var err error
		s2ml, err = s2mdec.ReadS2ML(C.GoBytes(unsafe.Pointer(bufS2ML), sizeS2ML))
		if err != nil {
			return C.CString(fmt.Sprint("s2ml: ", err))
		}
	}
	{ // s2mh + s2ml
		merged, err := s2mdec.S2MHApplyS2ML(s2mh, s2ml, nil)
		if err != nil {
			return C.CString(fmt.Sprint("s2mh plus s2ml: ", err))
		}
		sb := strings.Builder{}
		if errJSON := writeJSON(&sb, merged); errJSON != nil {
			log.Println("s2mh plus s2ml:", errJSON)
			return
		}
		return C.CString(sb.String())
	}
}

func main() {
	// Entry point: Do nothing.
}

func writeJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	return enc.Encode(v)
}
