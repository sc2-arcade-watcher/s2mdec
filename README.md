# s2mdec

[![GoDoc](https://godoc.org/github.com/sc2-arcade-watcher/s2mdec?status.svg)](https://godoc.org/github.com/sc2-arcade-watcher/s2mdec)

Package s2mdec is a decoder of s2mi and s2mh files. Go port of [s2m-decoder](https://github.com/sc2-arcade-watcher/s2m-decoder).

- - -

## Use as a Go library

### Install
```cmd
go get -v github.com/sc2-arcade-watcher/s2mdec
```

### Import
```Go
import "github.com/sc2-arcade-watcher/s2mdec"
```

- - -

## Use as a C library

### Build lib
```bash
$ cd ./s2mdeclibc/
$ make
```

### Include
```C
#include "s2mdeclibc.h"
```

### Decode
```C
// Decode s2mi
char *json = s2mdec_read_s2mi(buf, size);

// Decode s2mh
char *json = s2mdec_read_s2mh(buf, size);

// Decode s2ml
char *json = s2mdec_read_s2ml(buf, size);

// Inline translation
char *json = s2mdec_read_s2mh_s2ml(bufH, sizeH, bufL, sizeL);
```

### Example
```C
#include <stdio.h>
#include <stdlib.h>
#include "s2mdeclibc.h"

void main()
{
    FILE *f = fopen("a1f3dfd0ceea4562e5045b6ebb0ad97f416308254c137c28455f55e6535c9f8b.s2mh", "rb");
    fseek(f, 0, SEEK_END);
    long fsize = ftell(f);
    fseek(f, 0, SEEK_SET); // rewind(f);

    char *string = malloc(fsize + 1);
    fread(string, 1, fsize, f);
    fclose(f);

    char *json = s2mdec_read_s2mh(string, fsize);
    printf("%s", json);

    free(string);
    free(json);
}
```
```bash
$ gcc main.c -l:s2mdeclibc -L.
$ ./a
```
```
{"addDefaultPermissions":true,"addMultiMod":false,"arcadeInfo":{"gameInfoScreenshots":[{"caption":{"color":null,"index":0,"table":0},"picture":{"height":600,"index":1,"left":0,"top":0,"width":800}},{"caption":{"color":null,"index":0,"table":0},"picture":{"height":600,"index":2,"left":0,"top":0,"width":800}},{"caption":{"color":null,"index":0,"table":0},"picture":{"height":600,"index":3,"left":0,"top":0,"width":800}},{"caption":{"color":null,"index":0,"table":0},"picture":{"height":600,"index":4,"left":0,"top":0,"width":800}},{"caption":{"color":null,"index":0,"table":0},"picture":{"height":600,"index":5,"left":0,"top":0,"width":800}}],"howToPlayScreenshots":[],"howToPlaySections":[{"items":[{"color":null,"index":9,"table":0},{"color":null,"index":10,"table":0},{"color":null,"index":11,"table":0}],"listType":"bulleted","subtitle":{"color":null,"index":0,"table":0},"title":{"color":null,"index":8,"table":0}},{"items":[{"color":null,"index":13,"table":0}],"listType":"none","subtitle":{"color":null,"index":0,"table":0},"title":{"color":null,"index":12,"table":0}}],"mapIcon":{"height":150,"index":6,"left":0,"top":0,"width":225},"matchmakerTags":[],"patchNoteSections":[{"items":[{"color":null,"index":17,"table":0},{"color":null,"index":18,"table":0},{"color":null,"index":19,"table":0},{"color":null,"index":20,"table":0}],"listType":"none","subtitle":{"color":null,"index":16,"table":0},"title":{"color":null,"index":15,"table":0}},{"items":[{"color":null,"index":23,"table":0},{"color":null,"index":24,"table":0}],"listType":"none","subtitle":{"color":null,"index":22,"table":0},"title":{"color":null,"index":21,"table":0}},{"items":[{"color":null,"index":27,"table":0},{"color":null,"index":28,"table":0},{"color":null,"index":29,"table":0}],"listType":"none","subtitle":{"color":null,"index":26,"table":0},"title":{"color":null,"index":25,"table":0}},{"items":[{"color":null,"index":32,"table":0},{"color":null,"index":33,"table":0}],"listType":"none","subtitle":{"color":null,"index":31,"table":0},"title":{"color":null,"index":30,"table":0}},{"items":[{"color":null,"index":36,"table":0},{"color":null,"index":37,"table":0},{"color":null,"index":38,"table":0},{"color":null,"index":39,"table":0},{"color":null,"index":40,"table":0},{"color":null,"index":41,"table":0},{"color":null,"index":42,"table":0},{"color":null,"index":43,"table":0},{"color":null,"index":44,"table":0}],"listType":"none","subtitle":{"color":null,"index":35,"table":0},"title":{"color":null,"index":34,"table":0}}],"tutorialLink":null,"website":{"color":null,"index":45,"table":0}},"archiveHandle":{"hash":"7067e8b25868263f1c2006c0722983114d026293779d60c738512308a7b4480c","region":"us","type":"s2ma"},"attributes":[],"defaultVariantIndex":0,"extraDependencies":[{"id":288191,"version":0},{"id":12,"version":0}],"filename":"ColdVoyage.SC2Map","header":{"id":289177,"version":65568},"localeTable":[{"locale":"enUS","stringTable":[{"hash":"2177ce1bf9453f634457242b2b8758df5390e2c531732fb7e6fa027c8ba4a5c6","region":"us","type":"s2ml"}]}],"mapNamespace":362949,"mapSize":{"horizontal":256,"vertical":256},"relevantPermissions":[],"specialTags":[],"tileset":{"color":null,"index":3,"table":0},"variants":[{"achievementTags":[],"attributeDefaults":[{"attribute":{"id":1001,"namespace":999},"value":{"_unk_attr_val_1":0,"index":0}},{"attribute":{"id":2000,"namespace":999},"value":{"_unk_attr_val_1":0,"index":10}},{"attribute":{"id":3006,"namespace":999},"value":{"_unk_attr_val_1":0,"index":2}},{"attribute":{"id":3015,"namespace":999},"value":{"_unk_attr_val_1":0,"index":0}},{"attribute":{"id":2018,"namespace":999},"value":[{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":10},{"_unk_attr_val_1":0,"index":20},{"_unk_attr_val_1":0,"index":30},{"_unk_attr_val_1":0,"index":40},{"_unk_attr_val_1":0,"index":50},{"_unk_attr_val_1":0,"index":60},{"_unk_attr_val_1":0,"index":70},{"_unk_attr_val_1":0,"index":80},{"_unk_attr_val_1":0,"index":90},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0}]},{"attribute":{"id":500,"namespace":999},"value":[{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1}]}],"attributeVisibility":[{"attribute":{"id":3006,"namespace":999},"hidden":1},{"attribute":{"id":3000,"namespace":999},"hidden":1},{"attribute":{"id":3010,"namespace":999},"hidden":1},{"attribute":{"id":3015,"namespace":999},"hidden":1},{"attribute":{"id":3004,"namespace":999},"hidden":1},{"attribute":{"id":3003,"namespace":999},"hidden":1},{"attribute":{"id":3001,"namespace":999},"hidden":1}],"categoryDescription":{"color":null,"index":5,"table":0},"categoryId":1,"categoryName":{"color":null,"index":4,"table":0},"lockedAttributes":[{"attribute":{"id":1001,"namespace":999},"value":{"Count":16,"Data":"0x0000"}},{"attribute":{"id":2000,"namespace":999},"value":{"Count":16,"Data":"0x0000"}},{"attribute":{"id":2018,"namespace":999},"value":{"Count":16,"Data":"0xff03"}}],"maxHumanPlayers":10,"maxOpenSlots":16,"maxTeamSize":10,"modeDescription":{"color":null,"index":7,"table":0},"modeId":1,"modeName":{"color":null,"index":6,"table":0}}],"workingSet":{"bigMap":{"height":512,"index":0,"left":0,"top":0,"width":512},"description":{"color":null,"index":2,"table":0},"instances":[{"attribute":{"id":500,"namespace":999},"value":[{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1}]},{"attribute":{"id":3007,"namespace":999},"value":[{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":1}]},{"attribute":{"id":3001,"namespace":999},"value":[{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":1,"index":2},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0}]},{"attribute":{"id":3002,"namespace":999},"value":[{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":1},{"_unk_attr_val_1":0,"index":2},{"_unk_attr_val_1":0,"index":3},{"_unk_attr_val_1":0,"index":4},{"_unk_attr_val_1":0,"index":5},{"_unk_attr_val_1":0,"index":6},{"_unk_attr_val_1":0,"index":7},{"_unk_attr_val_1":0,"index":8},{"_unk_attr_val_1":0,"index":9},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0},{"_unk_attr_val_1":0,"index":0}]}],"localeTable":[{"locale":"enUS","stringTable":[{"hash":"2177ce1bf9453f634457242b2b8758df5390e2c531732fb7e6fa027c8ba4a5c6","region":"us","type":"s2ml"}]}],"maxPlayers":10,"name":{"color":null,"index":1,"table":0},"thumbnail":null,"visualFiles":[{"hash":"ab3f92a8e025b25fca4a8cc8c8736b0b3614c3dcc7875a5b6b3c94f391124a63","region":"us","type":"s2mv"},{"hash":"c50888599249b83e3ec4463cd8467279394839e11668f3cb28312a5545c706fd","region":"us","type":"s2mv"},{"hash":"b952eef72d5009bcaa153fde74511b3b4be5aff7671d03491e1f16b090b07a8f","region":"us","type":"s2mv"},{"hash":"47834a60f73bfaf9d857e10e9aee0596a4b90ba49a0b8b7f63e05c04b72b6d2a","region":"us","type":"s2mv"},{"hash":"61faf4c17b59e368b43693a218cb24af46bd6c7950d36fd1edab869a640f3283","region":"us","type":"s2mv"},{"hash":"d9e0f7b9756900581119faa3337563ffebc455e877c20b10ab0c8df5af14f870","region":"us","type":"s2mv"},{"hash":"d2e06878a567b7a4b1f6957bb75d40c6d356890fe0d498d41444da8c9dc19f4b","region":"us","type":"s2mv"}]}}
```

- - -

## Use as a command-line application

### Build executable
```bash
$ cd ./cmd/s2mdec/
$ make
```

### Decode s2mi
```bash
$ ./s2mdec f1cee150f270b7e4f2c1dabbca0838832f06a9f1336affd1e2f5528a841b5b4f.s2mi
```

### Decode s2mh
```bash
$ ./s2mdec 396811b3e2b6a4abe6396bccc3dca610915cce8cbfe90c32883ff8d8616af85c.s2mh
```

### Decode s2ml
```bash
$ ./s2mdec 68656032d079231b33c6d4c4e1e0710c4d6eee1793a6a640fc05bc6c107e1518.s2ml
```

### Inline translation
```bash
$ ./s2mdec 396811b3e2b6a4abe6396bccc3dca610915cce8cbfe90c32883ff8d8616af85c.s2mh 68656032d079231b33c6d4c4e1e0710c4d6eee1793a6a640fc05bc6c107e1518.s2ml
```

### Produce compact outcome
```bash
$ ./s2mdec -c 396811b3e2b6a4abe6396bccc3dca610915cce8cbfe90c32883ff8d8616af85c.s2mh
```
- - -
