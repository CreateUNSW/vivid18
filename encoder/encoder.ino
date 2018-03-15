#include "FastLED.h"

/*some struct to be added here*/
typedef struct {
//declar some structure here
  int addr;
  bool updir;
  int prior;
  int pace;
  int eff;
  int lumA;
  int colorA;
  int lumB;
  int colorB;
}DataStream;

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
void setup() {
  Serial.begin(9600);
}

int decode (char input []) {
  // takes a 4 byte brinary input and returns a 32 bit int
  // little endian or big boy?
  //ie 0000 0000 0000 0000
  // need to run some test to confirm
  unsigned long result;
  //result = input[3];
  //result = (result<<8)|input[2];
  //result = (result<<8)|input[1];
  result = input[0] & 0xff;
  result = input[1] & 0xff | result<<8;
  result = input[2] & 0xff | result<<8;
  result = input[3] & 0xff | result<<8;


  Serial.println(result);
  return result;
}


int encode (DataStream memes) {
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

void loop(){
    //do some test here
  /*char test1[4];
  test1[0] = 0b11101001;
  test1[1] = 0b10100110;
  test1[2] = 0b10110110;
  test1[3] = 0b11010000;
  char *p = test1;
  decode (p);
  */

  DataStream memes;
  //inputs 1110 1001 1010 011 01011 011 01101 0000
  memes.addr = 7;
  memes.updir = false;
  
  memes.prior = 2;
  memes.pace = 1;
  
  memes.eff = 5;
  
  memes.lumA = 3;
  memes.colorA = 11;
  
  memes.lumB = 3;
  memes.colorB = 29;

  //terminator
  //DataStream *memesAdd = memes;
  
  long test2;
  test2= encode(memes);

  
  delay (1000);

  
  
}

