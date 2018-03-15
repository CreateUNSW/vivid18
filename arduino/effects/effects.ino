#include "FastLED.h"

// How many leds in your strip?
#define NUM_LEDS 50

// For led chips like Neopixels, which have a data line, ground, and power, you just
// need to define DATA_PIN.  For led chipsets that are SPI based (four wires - data, clock,
// ground, and power), like the LPD8806 define both DATA_PIN and CLOCK_PIN
#define DATA_PIN 3
#define CLOCK_PIN 13

// Define the array of leds
CRGB leds[NUM_LEDS];
uint8_t gHue = 0;


void setup() { 
  //initialise random to pin 0
  randomSeed(analogRead(0));
  // Uncomment/edit one of the following lines for your leds arrangement.
  // FastLED.addLeds<TM1803, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<TM1804, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<TM1809, DATA_PIN, RGB>(leds, NUM_LEDS);
  FastLED.addLeds<WS2811, DATA_PIN, RBG>(leds, NUM_LEDS);
  // FastLED.addLeds<WS2812, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<WS2812B, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<NEOPIXEL, DATA_PIN>(leds, NUM_LEDS);
  // FastLED.addLeds<APA104, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<UCS1903, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<UCS1903B, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<GW6205, DATA_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<GW6205_400, DATA_PIN, RGB>(leds, NUM_LEDS);
  
  // FastLED.addLeds<WS2801, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<SM16716, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<LPD8806, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<P9813, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<APA102, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<DOTSTAR, RGB>(leds, NUM_LEDS);
  
  // FastLED.addLeds<WS2801, DATA_PIN, CLOCK_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<SM16716, DATA_PIN, CLOCK_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<LPD8806, DATA_PIN, CLOCK_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<P9813, DATA_PIN, CLOCK_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<APA102, DATA_PIN, CLOCK_PIN, RGB>(leds, NUM_LEDS);
  // FastLED.addLeds<DOTSTAR, DATA_PIN, CLOCK_PIN, RGB>(leds, NUM_LEDS);
}

void trinkleStar (){
  //possibly tinkle in a different colour
  
  //random a led position 
  //binkle for 1 sec then fades in brightness
  //and sleeps
   int randNum = random(NUM_LEDS+1);
   for (int bVal = 255; bVal != 0; bVal -=5 ){
    //add the colour here
    //fill in here angle, sat, brightness
    leds[randNum].setHSV(175, 255, bVal);
    delay(20);
    FastLED.show();
   }
   delay(500);
}

void flood(){
  //possibly add different fill speed by adjusting the delay time
  //and diffrent saturation with brightness
   int colorAngle = random (360);
   int index = 0;
   int lightIndex = 0;
   for (int count = NUM_LEDS; count != 0; count -- ){
    leds[index++].setHSV(colorAngle, 255, 100);
    delay(200);
    FastLED.show();
   }
   delay(500);
}

void rainbowWave(){
  //the led strip is filled with solid colour
  //new colour is swiped through the strip like a wave effect
  int fillColor = random(360);
  int midPoint = NUM_LEDS/2;
  int rightPos = midPoint;
  int leftPos = midPoint;
  
  //while last pointer does not equal to end of strip 
  //update the position of the 3 pointers along the strip 
  for (int start = 0; start < midPoint; start ++){
    //might keep track of previous filled leds and change the led sat and brightness
    leds[rightPos].setHSV(fillColor, 255, 100);
    leds[leftPos].setHSV(fillColor, 255, 100);
    rightPos ++;
    leftPos --;
     delay (100);
     FastLED.show();
  }
  
  
  delay(1000);

}
void loop() {
  //the basic effects
  //trinkleStar();
  //flood();
  /*static uint8_t hue = 0;
  FastLED.showColor(CHSV(hue++, 255, 255));
  delay(100);*/
  
}



