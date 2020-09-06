package s2mdec

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/icza/s2prot"
)

// Representation of a decoded struct types.
// - Null:     nil                     (untyped nil)
// - Struct:   map[string]interface{}
// - Array:    []interface{}           (Not []string)
// - Integer:  int64                   (Not int)
// - String:   string                  (Not []byte only castable)
// - Blob:     string                  (Not []byte only castable)
// - Bytes:    string                  (Not []byte only castable)
// - BitArray: s2prot.BitArr           (Value stored in s2prot.BitArr.Data []byte)
// - Optional: (any of above)

// ----------------------------------------------------------

var errStructInvalid = errors.New("invalid struct")

func makeErrVer(ver int) error {
	return fmt.Errorf("unexpected ver: %v", ver)
}

func makeErrStructLen(argStruct s2prot.Struct) error {
	return fmt.Errorf("unexpected struct len: %d", len(argStruct))
}

func makeErrArrayLen(argArray []interface{}) error {
	return fmt.Errorf("unexpected array len: %d", len(argArray))
}

func isIntIn(a int, s []int) bool {
	for _, v := range s {
		if a == v {
			return true
		}
	}
	return false
}

func isStrIn(a string, s []string) bool {
	for _, v := range s {
		if a == v {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------

// toBool from int.
func toBool(d int) (bool, error) {
	switch d {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		// Do nothing. (fallthrough)
	}
	return d != 0, fmt.Errorf("unexpected value casted to bool: 0x%X", d)
}

// verOf arg as in max(dict.keys()).
func verOf(unlabeled s2prot.Struct) (int, error) {
	ver := -1.0
	for k := range unlabeled {
		n, err := strconv.Atoi(k)
		if err != nil { // key was not a number
			return -1, fmt.Errorf("unexpected struct key: %s: %v", k, err)
		}
		ver = math.Max(float64(ver), float64(n))
	}
	if ver >= 0.0 { // not empty struct
		return int(ver), nil
	}
	return -1, errors.New("empty struct")
}

// readArrayOfStructs is newArrayOfStructsFromArrayOfStructs.
func readArrayOfStructs(fnMap func(s2prot.Struct) s2prot.Struct, arr []interface{}) []interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArrayOfStructs: %v", r))
		}
	}()
	//
	ret := make([]interface{}, len(arr))
	for i, v := range arr {
		ret[i] = fnMap(v.(s2prot.Struct))
	}
	return ret
}

// readArrayOfStringBytes is newArrayOfStructsFromArrayOfStringBytes.
func readArrayOfStringBytes(fnMap func([]byte) s2prot.Struct, arr []interface{}) []interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArrayOfStringBytes: %v", r))
		}
	}()
	//
	ret := make([]interface{}, len(arr))
	for i, v := range arr {
		ret[i] = fnMap([]byte(v.(string)))
	}
	return ret
}

// readArrayOfStrings is newArrayOfStringsFromArrayOfStrings.
func readArrayOfStrings(fnMap func(string) string, arr []interface{}) []interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArrayOfStrings: %v", r))
		}
	}()
	//
	ret := make([]interface{}, len(arr))
	for i, v := range arr {
		ret[i] = fnMap(v.(string))
	}
	return ret
}

// ----------------------------------------------------------

// throws error
func readDepotLink(unlabeled []byte) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readDepotLink: %v", r))
		}
	}()
	//
	return s2prot.Struct{
		"type":   string(unlabeled[:4]),
		"region": strings.ToLower(strings.Trim(string(unlabeled[4:8]), "\x00")),
		"hash":   hex.EncodeToString(unlabeled[8:]),
	}
}

// throws error
func readLocalizationLink(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readLocalizationLink: %v", r))
		}
	}()
	//
	return s2prot.Struct{
		"locale":      unlabeled.Stringv("0"),
		"stringTable": readArrayOfStringBytes(readDepotLink, unlabeled.Array("1")),
	}
}

// throws error
func readInstanceHeader(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readInstanceHeader: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"id":      unlabeled.Int("0"),
		"version": unlabeled.Int("1"),
		// "majorVersion": unlabeled.Int("1") >> 16,    //
		// "minorVersion": unlabeled.Int("1") & 0xFFFF, //
	}
}

// throws error
func readLocalizationTableKey(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readLocalizationTableKey: %v", r))
		}
	}()
	//
	if unlabeled == nil { // arg nil being no error nor panic
		return nil
	}
	if len(unlabeled) != 3 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"color": unlabeled.Value("0"), // optional: int or nil
		"table": unlabeled.Int("1"),
		"index": unlabeled.Int("2"),
	}
}

// throws error
func readPicture(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readPicture: %v", r))
		}
	}()
	//
	if unlabeled == nil { // arg nil being no error nor panic
		return nil
	}
	if len(unlabeled) != 5 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"index":  unlabeled.Int("0"),
		"top":    unlabeled.Int("1"),
		"left":   unlabeled.Int("2"),
		"height": unlabeled.Int("3"),
		"width":  unlabeled.Int("4"),
	}
}

// throws error
func readScreenshotEntry(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readScreenshotEntry: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"picture": readPicture(unlabeled.Structv("0")),
		"caption": readLocalizationTableKey(unlabeled.Structv("1")),
	}
}

// throws error
func readWorkingSet(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readWorkingSet: %v", r))
		}
	}()
	//
	// ver
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		panic(errVer) // throw
	}
	if !isIntIn(ver, []int{8, 10, 11}) { // assert
		panic(makeErrVer(ver)) // throw
	}
	if v := unlabeled.Int("5"); v != 22 { // assert data[5] == 22
		panic(fmt.Errorf("unexpected value of workingSet[5]: %v", v)) // throw
	}
	if ver >= 10 {
		if a := unlabeled.Array("9"); len(a) != 0 { // assert data[9] == 0
			panic(makeErrArrayLen(a)) // throw
		}
		if a := unlabeled.Array("10"); len(a) != 0 { // assert data[10] == 0
			panic(makeErrArrayLen(a)) // throw
		}
	}
	if ver >= 11 {
		// 11: [{0: {0: 999, 1: 3004}, 1: 1011, 2: [{0: b'\x00Lic', 1: 162}]}]}
		// Do nothing. (fallthrough)
	}
	// return
	return s2prot.Struct{
		"name":        readLocalizationTableKey(unlabeled.Structv("0")),
		"description": readLocalizationTableKey(unlabeled.Structv("1")),
		"thumbnail":   readPicture(unlabeled.Structv("2")),
		"bigMap":      readPicture(unlabeled.Structv("3")),
		"maxPlayers":  unlabeled.Int("4"),
		"instances":   readArrayOfStructs(readVariantAttributeDefaults, unlabeled.Array("6")),
		"visualFiles": readArrayOfStringBytes(readDepotLink, unlabeled.Array("7")),
		"localeTable": readArrayOfStructs(readLocalizationLink, unlabeled.Array("8")),
	}
}

// throws error
func readPremiumInfo(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readPremiumInfo: %v", r))
		}
	}()
	//
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		panic(errVer) // throw
	}
	if !isIntIn(ver, []int{0}) { // assert
		panic(makeErrVer(ver)) // throw
	}
	return s2prot.Struct{
		"license": unlabeled.Int("0"),
	}
}

// throws error
func readAttributeVisual(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readAttributeVisual: %v", r))
		}
	}()
	//
	if len(unlabeled) != 3 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"text": readLocalizationTableKey(unlabeled.Structv("0")),
		"tip":  readLocalizationTableKey(unlabeled.Structv("1")),
		"art":  readPicture(unlabeled.Structv("2")),
	}
}

// throws error
func readAttributeValueDefinition(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readAttributeValueDefinition: %v", r))
		}
	}()
	//
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		panic(errVer) // throw
	}
	if !isIntIn(ver, []int{1, 2}) { // assert
		panic(makeErrVer(ver)) // throw
	}
	if ver >= 2 {
		if a := unlabeled.Array("2"); len(a) != 0 { // assert // array?
			panic(makeErrArrayLen(a)) // throw // array?
		}
	}
	return s2prot.Struct{
		"value":  strings.Trim(unlabeled.Stringv("0"), "\x00"),
		"visual": readAttributeVisual(unlabeled.Structv("1")),
	}
}

// throws error
func readAttributeDefinition(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readAttributeDefinition: %v", r))
		}
	}()
	//
	// alias definition
	type arbitration = int64
	const (
		arbitrationAlways arbitration = iota // 0x00 always
		arbitrationFCFS                      // 0x01 first come first serve
	)
	//
	type visibility = int64
	const (
		visibilityNone visibility = iota // 0: "none",
		visibilitySelf                   // 1: "self",
		visibilityHost                   // 2: "host",
		visibilityAll                    // 3: "all",
	)
	//
	type flagOption = int64
	const (
		flagOptionUnknown          flagOption = 0x01 // 0x01 unknown
		flagOptionLockedWhenPublic flagOption = 0x02 // 0x02 locked when public
		flagOptionHidden           flagOption = 0x04 // 0x04 hidden
	)
	//
	return s2prot.Struct{
		"instance": readAttributeLink(unlabeled.Structv("0")),
		"values":   readArrayOfStructs(readAttributeValueDefinition, unlabeled.Array("1")),
		"visual":   readAttributeVisual(unlabeled.Structv("2")),
		// "_requirements": unlabeled.Value("3"), // unknown type
		"arbitration": arbitration(unlabeled.Int("4")),
		"visibility":  visibility(unlabeled.Int("5")),
		"access":      visibility(unlabeled.Int("6")),
		"options":     flagOption(unlabeled.Int("7")),
		"default":     readAttributeDefaultValueOrValues(unlabeled.Value("8")), // optional type
		"sortOrder":   unlabeled.Int("9"),
	}
}

// throws error
func readAttributeLink(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readAttributeLink: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"namespace": unlabeled.Int("0"),
		"id":        unlabeled.Int("1"),
	}
}

// throws error
func readAttributeDefaultValue(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readAttributeDefaultValue: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"index": unlabeled.Int("0"),
		// "_unk_attr_val_1": unlabeled.Int("1"),
	}
}

// throws error
func readAttributeDefaultValueOrValues(argDataValue interface{}) interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readAttributeDefaultValueOrValues: %v", r))
		}
	}()
	//
	// either a list of structs or a single struct
	switch attrValue := argDataValue.(type) {
	case []s2prot.Struct: // never but just in case
		retAttrValues := make([]s2prot.Struct, len(attrValue))
		for i, v := range attrValue {
			retAttrValues[i] = readAttributeDefaultValue(v)
		}
		return retAttrValues
	case []interface{}:
		return readArrayOfStructs(readAttributeDefaultValue, attrValue)
	case s2prot.Struct:
		return readAttributeDefaultValue(attrValue)
	default:
		// Do nothing. (fallthrough)
	}
	panic("unexpected type in attribute value") // throw
}

// throws error
func readVariantAttributeDefaults(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readVariantAttributeDefaults: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"attribute": readAttributeLink(unlabeled.Structv("0")),
		"value":     readAttributeDefaultValueOrValues(unlabeled.Value("1")), // optional type
	}
}

// throws error
func readVariantAttributeLocked(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readVariantAttributeLocked: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"attribute":    readAttributeLink(unlabeled.Structv("0")),
		"lockedScopes": int64(binary.BigEndian.Uint16(unlabeled.BitArr("1").Data)), // 16-bit integer whose bit is for each slot in lobby
	}
}

// throws error
func readVariantAttributeVisibility(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readVariantAttributeVisibility: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"attribute": readAttributeLink(unlabeled.Structv("0")),
		"hidden":    unlabeled.Int("1"),
	}
}

// throws error
func readVariantInfo(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readVariantInfo: %v", r))
		}
	}()
	//
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		panic(errVer) // throw
	}
	if !isIntIn(ver, []int{8, 11, 12, 13, 14, 15}) { // assert
		panic(makeErrVer(ver)) // throw
	}
	if v := unlabeled.Structv("0"); len(v) != 2 { // assert
		panic(makeErrStructLen(v)) // throw
	}
	if v := unlabeled.Structv("5"); len(v) != 3 { // assert
		panic(makeErrStructLen(v)) // throw
	}
	ret := s2prot.Struct{
		"categoryId":          unlabeled.Int("0", "0"),
		"modeId":              unlabeled.Int("0", "1"),
		"categoryName":        readLocalizationTableKey(unlabeled.Structv("1")),
		"modeName":            readLocalizationTableKey(unlabeled.Structv("2")),
		"categoryDescription": readLocalizationTableKey(unlabeled.Structv("3")),
		"modeDescription":     readLocalizationTableKey(unlabeled.Structv("4")),
		"attributeDefaults":   readArrayOfStructs(readVariantAttributeDefaults, unlabeled.Array("6")),
		"lockedAttributes":    readArrayOfStructs(readVariantAttributeLocked, unlabeled.Array("7")),
		"maxTeamSize":         unlabeled.Int("8"),
	}
	if ver >= 11 {
		ret["attributeVisibility"] = readArrayOfStructs(readVariantAttributeVisibility, unlabeled.Array("9"))
		// TODO: data[10] - attribute value restrictions?
		// 10: [{0: {0: 999, 1: 500}, 1: {0: [8, 8, 8, 8]}}], // bitarray?
		ret["achievementTags"] = readArrayOfStrings(func(s string) string { return strings.Trim(s, "\x00") }, unlabeled.Array("11"))
	}
	if ver >= 12 {
		ret["maxHumanPlayers"] = unlabeled.Value("12") // optional: int or nil
	}
	if ver >= 13 {
		ret["maxOpenSlots"] = unlabeled.Value("13") // optional: int or nil
	}
	if ver >= 14 {
		ret["premiumInfo"] = func(argStruct s2prot.Struct) s2prot.Struct {
			if argStruct == nil {
				return nil
			}
			return readPremiumInfo(argStruct)
		}(unlabeled.Structv("14"))
	}
	if ver >= 15 {
		ret["teamNames"] = readArrayOfStructs(readLocalizationTableKey, unlabeled.Array("15"))
	}
	return ret
}

// throws error
func readArcadeSectionHeader(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArcadeSectionHeader: %v", r))
		}
	}()
	//
	type listFormat = int64 // LIST_TYPE
	const (
		listFormatBulleted listFormat = iota // 0: "bulleted",
		listFormatNumbered                   // 1: "numbered",
		listFormatNone                       // 2: "none",
	)
	//
	if len(unlabeled) != 4 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"title":       readLocalizationTableKey(unlabeled.Structv("0")),
		"startOffset": unlabeled.Int("1"),
		"listType":    listFormat(unlabeled.Int("2")),
		"subtitle":    readLocalizationTableKey(unlabeled.Structv("3")),
	}
}

// throws error
func readArcadeSectionRaw(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArcadeSectionRaw: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	return s2prot.Struct{
		"headers": readArrayOfStructs(readArcadeSectionHeader, unlabeled.Array("0")),
		"items":   readArrayOfStructs(readLocalizationTableKey, unlabeled.Array("1")),
	}
}

// throws error
func readArcadeSection(unlabeled s2prot.Struct) []interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArcadeSection: %v", r))
		}
	}()
	//
	if len(unlabeled) != 2 { // assert
		panic(makeErrStructLen(unlabeled)) // throw
	}
	sectHeaders := readArrayOfStructs(readArcadeSectionHeader, unlabeled.Array("0"))
	sectItems := readArrayOfStructs(readLocalizationTableKey, unlabeled.Array("1"))
	for i, prevOffset := len(sectHeaders)-1, int64(len(sectItems)); i >= 0; i-- { // reversed
		sectHeader := sectHeaders[i].(s2prot.Struct)
		startOffset := sectHeader.Int("startOffset")
		sectHeader["items"] = sectItems[startOffset:prevOffset]
		prevOffset = startOffset
		delete(sectHeader, "startOffset")
	}
	return sectHeaders
}

// throws error
func readArcadeTutorialLink(unlabeled s2prot.Struct) s2prot.Struct {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("readArcadeTutorialLink: %v", r))
		}
	}()
	//
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		panic(errVer) // throw
	}
	if ver != 2 { // assert
		panic(makeErrVer(ver)) // throw
	}
	// looks like link to map, but it's an array..
	arr := unlabeled.Array("2")
	if len(arr) != 1 { // assert
		panic(makeErrArrayLen(arr)) // throw
	}
	// data["2"][0]
	svNut := arr[0].(s2prot.Struct)
	if len(svNut) != 2 { // assert
		panic(makeErrStructLen(svNut)) // throw
	}
	// data["2"][0]["1"]
	if v := svNut.Int("1"); v != 0 { // assert
		panic(fmt.Errorf(`unexpected value of data["2"][0]["1"]: %v`, v)) // throw
	}
	return s2prot.Struct{
		"variantIndex": unlabeled.Int("0"),
		"speed":        unlabeled.Stringv("1"),
		"map":          readInstanceHeader(svNut.Structv("0")), // data["2"][0]["0"]
	}
	// "7": {
	// 	"0": 1,
	// 	"1": "Fasr",
	// 	"2": [               # arr
	// 		{                # svNut
	// 			"0": {
	// 				"0": 210321,
	// 				"1": 65551
	// 			},
	// 			"1": 0
	// 		}
	// 	]
	// },
}

// throws error
func readArcadeInfo(unlabeled s2prot.Struct) s2prot.Struct {
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		panic(errVer)
	}
	if ver != 9 {
		panic(makeErrVer(ver))
	}
	// extra visualizationFiles ??
	if a := unlabeled.Array("0"); len(a) != 0 { // assert // array? struct?
		panic(makeErrArrayLen(a)) // throw // array? struct?
	}
	// extra localizationFiles ??
	if a := unlabeled.Array("1"); len(a) != 0 { // assert // array? struct?
		panic(makeErrArrayLen(a)) // throw // array? struct?
	}
	return s2prot.Struct{
		"gameInfoScreenshots":  readArrayOfStructs(readScreenshotEntry, unlabeled.Array("2")),
		"howToPlayScreenshots": readArrayOfStructs(readScreenshotEntry, unlabeled.Array("3")),
		"howToPlaySections":    readArcadeSection(unlabeled.Structv("4")),
		"patchNoteSections":    readArcadeSection(unlabeled.Structv("5")),
		"mapIcon":              readPicture(unlabeled.Structv("6")),
		"tutorialLink": func(argStruct s2prot.Struct) s2prot.Struct {
			if argStruct == nil {
				return nil
			}
			return readArcadeTutorialLink(argStruct)
		}(unlabeled.Structv("7")),
		"matchmakerTags": readArrayOfStrings(func(s string) string { return strings.Trim(s, "\x00") }, unlabeled.Array("8")),
		"website":        readLocalizationTableKey(unlabeled.Structv("9")),
	}
}

// ReadS2MH reads s2mh.
func ReadS2MH(unlabeled s2prot.Struct) (retStruct s2prot.Struct, retError error) {
	defer func() {
		if r := recover(); r != nil {
			retStruct, retError = nil, fmt.Errorf("decoding error: %v", r)
		}
	}()
	//
	// assert arg
	if len(unlabeled) != 2 {
		return nil, makeErrStructLen(unlabeled)
	}
	// set arg
	if unlabeled = unlabeled.Structv("0"); unlabeled == nil {
		return nil, errStructInvalid
	}
	// ver
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		return nil, errVer
	}
	// assert ver
	if !isIntIn(ver, []int{13, 14, 18, 22, 23, 24}) {
		return nil, makeErrVer(ver)
	}
	// assert member
	if vStruct := unlabeled.Structv("0"); len(vStruct) != 2 {
		return nil, makeErrStructLen(vStruct)
	}
	// set ret
	retStruct = s2prot.Struct{
		"header":        readInstanceHeader(unlabeled.Structv("0")),
		"filename":      unlabeled.Stringv("1"), // utf8
		"archiveHandle": readDepotLink([]byte(unlabeled.Stringv("2"))),
		"mapNamespace":  unlabeled.Int("3"),
		"workingSet":    readWorkingSet(unlabeled.Structv("4")),
		"attributes":    readArrayOfStructs(readAttributeDefinition, unlabeled.Array("5")),
		"localeTable":   readArrayOfStructs(readLocalizationLink, unlabeled.Array("8")),
		"mapSize": func(argStruct s2prot.Struct) s2prot.Struct {
			if argStruct == nil {
				return nil
			}
			return s2prot.Struct{
				"horizontal": argStruct.Int("0"),
				"vertical":   argStruct.Int("1"),
			}
		}(unlabeled.Structv("9")),
		"tileset":             readLocalizationTableKey(unlabeled.Structv("10")),
		"defaultVariantIndex": unlabeled.Int("12"),
		"variants":            readArrayOfStructs(readVariantInfo, unlabeled.Array("13")),
	}
	// if v := unlabeled.Value("6"); v != nil {
	// 	retStruct["_unk6"] = fmt.Sprintf("%s", v)
	// }
	/*
		# TODO: 7 - score IDs and such?
		retStruct["resultDefinitions"] = []
	*/
	retStruct["specialTags"] = func(arg interface{}) []interface{} {
		// if 11 in data:
		// o['specialTags'] = [data[11].decode('ascii')] if data[11] is not None else []
		if arg == nil {
			return []interface{}{} // an empty list is returned instead of nil
		}
		switch argDiscerned := arg.(type) {
		case string:
			return []interface{}{argDiscerned}
		case []interface{}: // ?
			return readArrayOfStrings(func(s string) string { return strings.Trim(s, "\x00") }, argDiscerned) // ?
		default:
			panic(fmt.Errorf("unexpected type of special tags: %T", argDiscerned))
		}
	}(unlabeled.Value("11")) // optional type
	if ver >= 14 {
		retStruct["extraDependencies"] = readArrayOfStructs(readInstanceHeader, unlabeled.Array("14"))
	}
	if ver >= 18 {
		{ // toBool()
			var errBool error
			retStruct["addDefaultPermissions"], errBool = toBool(int(unlabeled.Int("15")))
			if errBool != nil {
				return nil, fmt.Errorf("field addDefaultPermissions: %v", errBool)
			}
		}
		retStruct["relevantPermissions"] = readArrayOfStructs(func(argStruct s2prot.Struct) s2prot.Struct {
			return s2prot.Struct{
				"name": strings.Trim(argStruct.Stringv("0"), "\x00"),
				"id":   argStruct.Int("1"),
			}
		}, unlabeled.Array("16"))
		// skip 17?
		retStruct["specialTags"] = readArrayOfStrings(func(s string) string { return strings.Trim(s, "\x00") }, unlabeled.Array("18"))
	}
	if ver >= 22 {
		retStruct["arcadeInfo"] = func(argStruct s2prot.Struct) s2prot.Struct {
			if argStruct == nil {
				return nil
			}
			return readArcadeInfo(argStruct)
		}(unlabeled.Structv("19"))
		{ // toBool()
			var errBool error
			retStruct["addMultiMod"], errBool = toBool(int(unlabeled.Int("22")))
			if errBool != nil {
				return nil, fmt.Errorf("field addMultiMod: %v", errBool)
			}
		}
	}
	if ver >= 23 {
		// 23: [b'SC2ParkVoicePack'
		//
		// Do nothing. (fallthrough)
	}
	if ver >= 24 {
		// array with a lot of numbers - possibly reward IDs?
		// 24: [23498
		//
		// Do nothing. (fallthrough)
	}
	{ // assert special tags
		mapSetKnownTags := map[string]struct{}{
			"BLIZ": {}, "TRIL": {}, "FEAT": {}, "PRGN": {},
			"HotS": {}, "LotV": {}, "WoL": {}, "WoLX": {},
			"HoSX": {}, "LoVX": {}, "HerX": {}, "Desc": {},
			"Glue": {}, "Blnc": {}, "PREM": {},
		}
		for _, v := range retStruct.Array("specialTags") {
			specialTag := v.(string)
			if _, ok := mapSetKnownTags[specialTag]; !ok {
				return nil, fmt.Errorf("unexpected special tag: %s", specialTag)
			}
		}
	}
	// catch and return
	return retStruct, retError
}

// ReadS2MI reads s2mi.
func ReadS2MI(unlabeled s2prot.Struct) (retStruct s2prot.Struct, retError error) {
	defer func() {
		if r := recover(); r != nil {
			retStruct, retError = nil, fmt.Errorf("decoding error: %v", r)
		}
	}()
	//
	// assert arg
	if len(unlabeled) != 2 {
		return nil, makeErrStructLen(unlabeled)
	}
	// set arg
	unlabeled = unlabeled.Structv("0")
	ver, errVer := verOf(unlabeled)
	if errVer != nil {
		return nil, errVer
	}
	// assert ver
	if !isIntIn(ver, []int{22, 23, 26}) {
		return nil, makeErrVer(ver)
	}
	// assert member
	if v := unlabeled.Structv("11"); len(v) != 4 {
		return nil, makeErrStructLen(v)
	}
	if v := unlabeled.Structv("14"); len(v) != 4 {
		return nil, makeErrStructLen(v)
	}
	// def
	readToon := func(argStruct s2prot.Struct) s2prot.Struct {
		return s2prot.Struct{
			"regionId":  argStruct.Int("0"),
			"app":       strings.Trim(argStruct.Stringv("1"), "\x00"),
			"realmId":   argStruct.Int("2"),
			"battleTag": argStruct.Stringv("3"),
		}
	}
	// set ret
	retStruct = s2prot.Struct{
		"header":            readInstanceHeader(unlabeled.Structv("0")),
		"headerCacheHandle": readDepotLink([]byte(unlabeled.Stringv("1"))),
		"uploadTime":        unlabeled.Int("2"),
		"isLinked":          unlabeled.Int("3") != 0, // bool
		"isLocked":          unlabeled.Int("4") != 0, // bool
		"isPrivate":         unlabeled.Int("5") != 0, // bool
		"mapSize":           unlabeled.Int("6"),
		"name":              unlabeled.Stringv("7"),
		// "profileRecordAddress":  unlabeled.Structv("8"),
		"isMod":                 unlabeled.Int("9") != 0, // bool
		"authorToonName":        readToon(unlabeled.Structv("11")),
		"isLatestVersion":       unlabeled.Int("12") != 0, // bool
		"mainLocale":            unlabeled.Stringv("13"),
		"authorToonHandle":      readToon(unlabeled.Structv("14")),
		"isSkipInitialDownload": unlabeled.Int("15") != 0, // bool
		"createdTime":           unlabeled.Int("16"),
		"labels":                unlabeled.Array("17"),    // readArrayOfStrings(func(s string) string { return s }, unlabeled.Array("17")),
		"isMelee":               unlabeled.Int("18") != 0, // bool
		"isCluster":             unlabeled.Int("19") != 0, // bool
		"clusterParent":         unlabeled.Int("20"),
		"clusterChildren":       unlabeled.Array("21"),
		"isHiddenLobby":         unlabeled.Int("22") != 0, // bool
		"isExtensionMod": func(argStruct s2prot.Struct) bool {
			v := argStruct.Value("23")
			if v == nil {
				return false
			}
			return v.(int64) != 0
		}(unlabeled),
	}
	if ver >= 24 {
		retStruct["transitionId"] = unlabeled.Int("24")
		retStruct["lastPublishTime"] = unlabeled.Int("25")
		retStruct["firstPublicPublishTime"] = unlabeled.Int("26")
	}
	// catch and return
	return retStruct, retError
}

// MapLocale translation.
type MapLocale map[string]string

// String returns the indented JSON string representation of the Struct.
// Defined with value receiver so this gets called even if a non-pointer is printed.
func (s2ml MapLocale) String() string {
	b, _ := json.MarshalIndent(s2ml, "", "  ")
	return string(b)
}

// ReadS2ML reads s2ml.
func ReadS2ML(rawXML []byte) (retTextByID MapLocale, retError error) {
	defer func() {
		if r := recover(); r != nil {
			retTextByID, retError = nil, fmt.Errorf("decoding error: %v", r)
		}
	}()
	//
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawXML); err != nil {
		return nil, fmt.Errorf("cannot parse xml: %v", err)
	}
	if elemRoot := doc.FindElement("Locale"); elemRoot == nil {
		return nil, errors.New("cannot find root element Locale")
	} else if children := elemRoot.ChildElements(); len(children) > 0 {
		retTextByID = map[string]string{}
		for _, child := range children {
			retTextByID[child.SelectAttrValue("id", "")] = child.Text()
		}
	}
	// catch and return
	return retTextByID, retError
}

// S2MHApplyS2ML adds s2ml to s2mh.
func S2MHApplyS2ML(s2mhLabeled s2prot.Struct, translation MapLocale, targetFields interface{}) (retStruct s2prot.Struct, retError error) {
	defer func() {
		if r := recover(); r != nil {
			retStruct, retError = nil, fmt.Errorf("decoding error: %v", r)
		}
	}()
	//
	// def
	translateProp := func(prop s2prot.Struct) interface{} {
		v := prop.Int("index")
		if v == 0 {
			return nil
		}
		return translation[strconv.Itoa(int(v))]
	}
	// targetFields
	if targetFields == nil {
		targetFields = map[string]interface{}{}
		if err := json.Unmarshal([]byte(`{
			"workingSet": {
				"name": true,
				"description": true
			},
			"tileset": true,
			"variants": {
				"categoryName": true,
				"modeName": true,
				"categoryDescription": true,
				"modeDescription": true
			},
			"arcadeInfo": {
				"gameInfoScreenshots": {
					"caption": true
				},
				"howToPlayScreenshots": {
					"caption": true
				},
				"howToPlaySections": {
					"title": true,
					"subtitle": true,
					"items": true
				},
				"patchNoteSections": {
					"title": true,
					"subtitle": true,
					"items": true
				},
				"website": true
			}
		}`), &targetFields); err != nil {
			return nil, fmt.Errorf("cannot determine target fields: %v", err)
		}
	}
	// recursive
	if mapTargetFields, ok := targetFields.(map[string]interface{}); ok {
		for keyTargetField, valTargetField := range mapTargetFields {
			if arr, ok := s2mhLabeled.Value(keyTargetField).([]interface{}); ok {
				for i, vStruct := range arr {
					if vBool, ok := valTargetField.(bool); vBool && ok {
						arr[i] = translateProp(vStruct.(s2prot.Struct))
					} else if _, ok := valTargetField.(map[string]interface{}); ok {
						S2MHApplyS2ML(vStruct.(s2prot.Struct), translation, valTargetField)
					}
				}
			} else {
				vStruct := s2mhLabeled.Value(keyTargetField).(s2prot.Struct) // presumed
				if vBool, ok2 := valTargetField.(bool); vBool && ok2 && vStruct != nil {
					s2mhLabeled[keyTargetField] = translateProp(vStruct)
				} else if mapValTargetField, ok := valTargetField.(map[string]interface{}); ok {
					S2MHApplyS2ML(vStruct, translation, mapValTargetField)
				}
			}
		}
	}
	// catch and return
	retStruct = s2mhLabeled
	return retStruct, retError
	// x          == keyTargetField
	// fields[x]  == valTargetField
	// data[x]    != valTargetField
	// data       == s2mhLabeled
	// data[x]    == s2mhLabeled.Value(keyTargetField)
	//            == arr  or  vStruct on right
	//            == s2mhLabeled[keyTargetField] on left
	// y          == i
	// data[x][y] == arr[i]
	//            == vStruct
	//
	//     for x in fields:
	//         if isinstance(data[x], list):
	//             for y, item in enumerate(data[x]):
	//                 if fields[x] is True:
	//                     data[x][y] = translate_prop(data[x][y])
	//                 elif isinstance(fields[x], dict):
	//                     s2mh_apply_s2ml(data[x][y], translation, fields[x])
	//         else:
	//             if fields[x] is True and data[x] is not None:
	//                 data[x] = translate_prop(data[x])
	//             elif isinstance(fields[x], dict):
	//                 s2mh_apply_s2ml(data[x], translation, fields[x])
	//
	//     return data
}
