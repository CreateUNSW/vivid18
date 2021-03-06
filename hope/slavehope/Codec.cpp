#include "Arduino.h"
#include "Codec.h"

/*
3 bits: address
1 bit: direction

2 bits: priority
2 bits: speed

3 bits: type of effect
1 bit padding

8 bits
Color A:
3 bits: lumosity (0 - 7)
5 bits: hue (0 - 31)

8 bits
Color B:
3 bits: lumosity (0 - 7)
5 bits: hue (0 - 31)

4 bit padding at the end as a terminator

total bits: 32 bits

example:
decimal: 3920017104
1110 1001 1010 011 01011 011 01101 0000

breakdown:
- address: 7
- is direction up: false
- priority 2, speed at 1
- effect type 5 follow with padding

colorA
- lumosity 3
- hue is 11

colorB
- lumosity 3
- hue is 13

and padding 4 bits 0000


*/

Codec::Codec () {}; 

DataStream Codec::decode (int input []) {
  //need to malloc
  // little endian or big boy?
  //unpack the stream into a datastruct
  DataStream result;
  int cutTemp;
  //one byte goes here
  result.addr = input[0] >> 5;

  //prob add an if condition here
  if(input[0] & (1<<4)){
    //Serial.println("this is true");
    result.updir = true;
  }else{
    result.updir = false;
  }

  //might wanna make a temp bit so it will have snipets
  cutTemp = result.addr<<5 | result.updir<<3;
  result.prior = (input[0] & ~cutTemp) >>2 ;
  cutTemp |= result.prior <<2;
  //Serial.println(cutTemp, BIN);
  result.pace = (input[0] & ~cutTemp);

  /*
  1. 111 0 10 01
  2. 1010 011 0
  3. 1011 011 0
  4. 1101 0000
  */

  //3bits and a pad + another 3 bits
  result.eff = input[1] >> 5 & 0xff;
  result.lumA = (input[1] & ~(result.eff << 5)) >> 1;

  //Serial.print("the lum a color: ");
  //Serial.println(result.lumA, BIN);
  cutTemp = input[2] >> 4;
  //need the last bit of the previous array + 5 bits
  result.colorA = (input[1] >> 7) | (cutTemp);
  //remain bits after 4th bit
  result.lumB = (input[2] & ~(cutTemp << 4)) >> 1;


  //need the last bit of the previous array and clear out the padding
  result.colorB = (input[2] >> 7 ) | (input[3] >>4);
  int bitNum = 0;

  //debugging code here
  Serial.println("debug\n===========================");
  Serial.print("Byte 1: ");
  Serial.println(input[0],BIN);
  Serial.print("Byte 2: ");
  Serial.println(input[1], BIN);
  Serial.print("Byte 3: ");
  Serial.println(input[2], BIN);
  Serial.print("Byte 4: ");
  Serial.println(input[3], BIN);

  //check the struct
  Serial.println("===========================\nConverted values\n===========================");
  Serial.print("address: ");
  Serial.println(result.addr);
  Serial.print("up boolean: ");
  Serial.println(result.updir);
  Serial.print("priority: ");
  Serial.println(result.prior);
  Serial.print("speed: ");
  Serial.println(result.pace);
  Serial.print("effect: ");
  Serial.println(result.eff);

  Serial.print("lum A: ");
  Serial.println(result.lumA);
  Serial.print("color A: ");
  Serial.println(result.colorA);
  Serial.print("lum B: ");
  Serial.println(result.lumB);
  Serial.print("color B: ");
  Serial.println(result.colorB);

  return result;
}


int Codec::encode (DataStream memes) {
  //need to malloc and free
  int tempBits;

  //temp clearing
  int results [4];

  results[0] = memes.addr & 0xff;
  results[0] = memes.updir | results[0] << 1;
  results[0] = memes.prior & 0xff | results[0] << 2;
  results[0] = memes.pace & 0xff | results[0] << 2;
  /*
  1. 11101001
  2. 1010 011 0
  3. 1011 011 0
  4. 1101 0000
  */
  results[1] =  memes.eff & 0xff | results[1] << 3;
  //add padding heres
  results[1] = 0 & 0xff | results[1] << 1;
  results[1] = memes.lumA & 0xff | results[1] << 3;

  tempBits = memes.colorA;

  if (tempBits >> 4 == 0b0){
    //Serial.println("this is trueee!");
    results[1] = 0 & 0xff | results[1] << 1;
    results[2] = (memes.colorA & 0xff);
  }else{
    results[1] = 1 & 0xff | results[1] <<1;
    tempBits = memes.colorA;
    results[2] = (tempBits & ~(1<<4)) & 0xff;
  }
  /*
  //if the leading bit is 0 then just pad or else shift it
  */
  results[2] = memes.lumB & 0xff | results[2] << 3;

  tempBits = memes.colorB;
  if (tempBits >> 4 == 0b0){
    //Serial.println("this is true number 2");
    results[2] = 0 & 0xff | results[2] << 1;
    results[3] = (memes.colorB & 0xff);
  }else{
    //Serial.println("this is false number 2");
    results[2] = 1 & 0xff | results[2] << 1;

    //Serial.print("checking number ");
    //Serial.println(memes.colorB, BIN);
    results[3] = (memes.colorB & ~(1<<4)) & 0xff;
  }

  results[3] = results[3] << 4;

  //devugging stuff
  Serial.println(results[0], BIN);
  Serial.println(results[1], BIN);
  Serial.println(results[2], BIN);
  Serial.println(results[3], BIN);
  Serial.println("========================");

  //return a pointer addresss
  return *results;
}

//helper functions
boolean Codec::validate (int stream[]){
  //check data integrity
  //check if the size is correct

  return true;
}
