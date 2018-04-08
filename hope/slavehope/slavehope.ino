#include <Wire.h>
#include <FastLED.h>
#include "Codec.h"
// Crystals/LEDs config
//--------------------------------------------------
#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811

#define NUM_CHAINS 16
//--------------------------------------------------

uint32_t data32 = 0;
int      data[4] = {0};
DataStream pattern;
Codec codec;

CRGB leds[NUM_LEDS];

//byte readByte()
//{
//  while(!Wire.available()) {}
//  return Wire.read();
//}

void setup() {
  // put your setup code here, to run once:
  Wire.begin(0x08);                // join i2c bus with address #8
  Wire.onReceive(receiveEvent); // register event
//  Wire.onRequest(sendEvent);
  Serial.begin(9600);           // start serial for output

  Serial.println("Ready");
  FastLED.addLeds<LED_TYPE, DATA_PIN, RBG>(leds, NUM_LEDS);
}

void loop() {
  // put your main code here, to run repeatedly:
  /*
  byte id = readByte();
  Serial.println(String(id));
  byte numLEDs = readByte();
  Serial.println(String(numLEDs));
  byte hue; byte lum;
  for(int i = 0; i < numLEDs; i++)
  {
    hue = readByte();
    lum = readByte();

    leds[i] = CHSV(hue, 255, lum);
  }

  FastLED.show();
  */
  delay(100);
  if(data[0] != 0) {
      pattern = codec.decode(data);
      Serial.println(pattern.addr);
      Serial.println(pattern.updir);
      Serial.println(pattern.prior);
      Serial.println(pattern.pace);
      Serial.println(pattern.eff);
      Serial.println(pattern.lumA);
      Serial.println(pattern.colorA);
      Serial.println(pattern.lumB);
      Serial.println(pattern.colorB);
  }
  // Serial.println("tick");
}
/*
  1 -> data pin
  2 -> Number LEDs
  n -> every 2 bytes (hue, lum)

*/
void receiveEvent(int howMany) {
  Serial.println("Begin receive");
  Serial.print("received: ");
  Serial.println(howMany);
  int i = 0;
  while (Wire.available()) {
    // data32 <<= 8;
    // data32 |= Wire.read();
//    Serial.println(data32);
    if(i < howMany) {
      data[i] = Wire.read();
      i++;
    }
  }

  for (i = 0; i < howMany; i++) {
    Serial.println(data[i]);
  }


      Serial.println("End receive");

//  int led = 0;
//  byte pin = Wire.read();
//  byte numLEDs = Wire.read();
//  if (numLEDs != (howMany - 2) / 2)
//  {
//    Serial.println("ERROR WOWOWOWO");
//  }
//  byte hue;
//  byte lum;
//  while(led < numLEDs && Wire.available())
//  {
//    hue = Wire.read();
//    lum = Wire.read();
//    leds[led] = CHSV(hue, 255, lum);
//    led++;
//  }
//  if (led != numLEDs)
//  {
//    Serial.println(String("ERROR BOOP BYE ") + led);
//  }
//
//  FastLED.show();
}

void sendEvent() {
  // Wire.write(4);
}
