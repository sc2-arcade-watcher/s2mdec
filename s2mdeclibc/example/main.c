#include <stdio.h>
#include <stdlib.h>
#include "s2mdeclibc.h"

// $ gcc main.c -l:s2mdeclibc -L.
void main()
{
    // files
    char *filename1 = "396811b3e2b6a4abe6396bccc3dca610915cce8cbfe90c32883ff8d8616af85c.s2mh";
    char *filename2 = "68656032d079231b33c6d4c4e1e0710c4d6eee1793a6a640fc05bc6c107e1518.s2ml";
    // http://kr.depot.battle.net:1119/396811b3e2b6a4abe6396bccc3dca610915cce8cbfe90c32883ff8d8616af85c.s2mh
    // http://kr.depot.battle.net:1119/68656032d079231b33c6d4c4e1e0710c4d6eee1793a6a640fc05bc6c107e1518.s2mh

    // str1
    FILE *f1 = fopen(filename1, "rb");
    fseek(f1, 0, SEEK_END);
    long fsize1 = ftell(f1);
    fseek(f1, 0, SEEK_SET); // rewind(f1);
    char *str1 = malloc(fsize1 + 1);
    fread(str1, 1, fsize1, f1);
    fclose(f1);

    // str2
    FILE *f2 = fopen(filename2, "rb");
    fseek(f2, 0, SEEK_END);
    long fsize2 = ftell(f2);
    fseek(f2, 0, SEEK_SET); // rewind(f2);
    char *str2 = malloc(fsize2 + 1);
    fread(str2, 1, fsize2, f2);
    fclose(f2);

    // json outcome
    char *json = s2mdec_read_s2mh_s2ml(str1, fsize1, str2, fsize2);
    printf("%s", json);

    // release memory
    free(str1);
    free(str2);
    free(json);
}
