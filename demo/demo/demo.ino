#include <Wire.h>
#include <FastLED.h>

#define NUM_LEDS 40
#define DATA_PIN 4
#define LED_TYPE UCS1903

CRGB leds[NUM_LEDS];

void setup() {
  Serial.begin(115200);
  Serial.setTimeout(20);
  pinMode(3, OUTPUT);
  FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(leds, NUM_LEDS); 
}

void loop() {
  Serial.readBytes((char*)leds, NUM_LEDS * 3);
  FastLED.show();
}

