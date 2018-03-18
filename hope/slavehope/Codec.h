#ifndef CODEC_H
#define CODEC_H

#include "Arduino.h"

/*Data structure containing LED information*/
typedef struct {
    int addr;
    bool updir;
    int prior;
    int pace;
    int eff;
    int lumA;
    int colorA;
    int lumB;
    int colorB;
} DataStream;

void DecodeTest();
void EncodeTest();

class Codec {
public:
    Codec();
    boolean validate(int stream[]);
    int encode(DataStream memes);
    DataStream decode(int input[]);
}; 

#endif
