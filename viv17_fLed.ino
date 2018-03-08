// a rough adaption of code from viv17

#include "FastLED.h"

// Crystals/LEDs config
//--------------------------------------------------
#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811
//--------------------------------------------------


// ====================
// Time variables
// t updates every frame, currently runs at x fps ie, seconds * x = t
uint8_t t = 0;
// t1 updates every x seconds, where x is fps
uint8_t t1 = 0;

// ====================
// Fade variables
double fadeSpeed = 1;
// the initial rate of fade
#define FADE_AMOUNT 1
// the rate at which fade changes
#define FADE_DELTA 0.01
boolean transition = true;
uint8_t hue = 0;
int chooseFern = 1;
int chooseTransition = 1;

CRGB leds[NUM_LEDS];
CRGB target[NUM_LEDS];

void setup() { 
  Serial.begin(9600);
  srand(0);
  FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(leds, NUM_LEDS); 
}

void loop() { 
  // configure centre led
  int centre = 12;

  // update fade speed with each cycle
  if(fadeSpeed > 1 + FADE_DELTA) fadeSpeed -= FADE_DELTA;

  // change pattern at set cycle rate
  if (t == 0) {
    t1++;
    if (t1 == 25) {
      t1 = 0;
      Serial.println(F("Change"));
      chooseFern = rand();
      chooseTransition = rand();
    }
  }

  // choose a pattern
  switch(chooseFern % 3) {
    case 0: rainbow(); break;
    case 1: randomFern(); break;
    case 2: mondrianColors(); break;
    default: rainbow();

    // TODO
      //shimmerCenter(centre);
      //crystalGradient();
      //muzzLight();
      //solidHue();
      //colorToWhiteHue();
      //complementaryHue();

  }

  // actually updates wall, updating method
  if (transition) {
    switch (chooseTransition % 1) {
      case 0: fadeTo(); break;
      default: jumpTo();
      // case 1: radialTo(centre); break;
    }
  } else {
    transition = true;
  }

  // t is global timer of range 0-255, don't change at all only use, create your own timer if needed
  t++;
  FastLED.show();
}

// ============ TRANSITION PATTERNS ============

// slow transition
void fadeTo() {
  int red, green, blue;
  for(int index = 0; index < NUM_LEDS; index++) {
    red = leds[index].red + ((target[index].r - leds[index].r) / fadeSpeed);
    green = leds[index].green+ ((target[index].g - leds[index].g) / fadeSpeed);
    blue = leds[index].blue + ((target[index].b - leds[index].b) / fadeSpeed);
    leds[index] = CRGB(red, blue, green);
  }
}

// instant transition
void jumpTo() {
  for(int index = 0; index < NUM_LEDS; index++) {
    leds[index] = CRGB(target[index].r, target[index].b, target[index].g);
  }
}

// ============ PATTERNS ============

// random
void randomFern() {
  Serial.println(F("Random"));
  if( t % 127 == 1) {
    for(int i = 0; i < NUM_LEDS; i++) {
      crystalHSV(i, rand() % 255,  255, 255); 
    }
    changeColor();
  }
}

void mondrianColors() {
    Serial.println(F("Mondrian"));
  CRGB red = CRGB(255, 0, 0);
  CRGB yellow = CRGB(255, 0, 255);
  CRGB blue = CRGB(0, 255, 0);
  CRGB white = CRGB(255, 255, 255);
  int rollDice = 0;
  if(t % 30 == 1) {
    for(int i = 0; i < NUM_LEDS; i++) {
      rollDice = rand()%4;
      if(rollDice == 0) crystalRGB(i, red.r,  red.g, red.b); 
      if(rollDice == 1) crystalRGB(i, yellow.r,  yellow.g, yellow.b); 
      if(rollDice == 2) crystalRGB(i, blue.r,  blue.g, blue.b);
      if(rollDice == 3) crystalRGB(i, white.r,  white.g, white.b);       
    }
  }
} 

// rainbow
void rainbow() {
  Serial.println(F("Rainbow"));
  if( t % 10 == 1) {
    for(int i = 0; i < NUM_LEDS; i++) {
      crystalHSV(i, hue++,  255, 255); 
    }
    changeColor();
  }
}

// ============ HELPER FUNCTIONS ============

void crystalRGB(int index, int r, int g, int b) {
  target[index] = CRGB(r, b, g);
}

void crystalHSV(int index, int h, int s, int v) {
//  target[index] = CHSV(h, s, v);
  target[index] = CRGB(0, 0, 0).setHSV(h, s, v);
}

void changeColor() {
    fadeSpeed = FADE_AMOUNT;
}

