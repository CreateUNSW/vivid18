#include "FastLED.h"


#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811
CRGB samp[NUM_LEDS];

void setup() { 
  FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(samp, NUM_LEDS); 
}
void loop() { 
  static uint8_t hue = 0;
  FastLED.showColor(CHSV(hue++, 255, 255)); 
  delay(10);
}
