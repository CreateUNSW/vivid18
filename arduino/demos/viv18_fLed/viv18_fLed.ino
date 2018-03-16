// a rough adaption of code from viv17

#include "FastLED.h"

// Crystals/LEDs config
//--------------------------------------------------
#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811
#define NUM_PATTERNS 2
//--------------------------------------------------


// ====================
// Time variables
// t updates every frame, currently runs at x fps ie, seconds * x = t
uint8_t t = 0;
// t1 updates every x seconds, where x is fps
uint8_t t1 = 0;

unsigned long wait = 0; 

// ====================
// Fade variables
double fadeSpeed = 1;
// the initial rate of fade
#define FADE_AMOUNT 1
// the rate at which fade changes
#define FADE_DELTA 0.01
#define PATTERN_DELTA 3
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

//  if (t > 255) { 
//    t = 0; 
//  }

//  Serial.print("t: "); 
//  Serial.print(t); 
//    Serial.print(" wait: "); 
//  Serial.print(wait); 
//      Serial.print(" millis: "); 
//  Serial.print(millis()); 
//  Serial.print(F(" Change to "));
//      Serial.print(chooseFern); 
//      Serial.print("\n");  

  // update fade speed with each cycle
  if(fadeSpeed > 1 + FADE_DELTA) fadeSpeed -= FADE_DELTA;

  // change pattern at set cycle rate
//  if (t == 0) {
//    t1++;
//    wait = millis(); 
    if (wait < 0) wait = 0; 
//    if (t1 == PATTERN_DELTA) {
//      t1 = 0;
//      chooseFern = rand()% NUM_PATTERNS;
//      chooseTransition = rand() ; 
//    }
//  }

    if (millis() - abs(wait) > PATTERN_DELTA*1000) {
        Serial.print("CHANGE"); 
        wait = millis(); 
        if (wait < 0) wait = 0; 
        chooseFern = rand() % NUM_PATTERNS; 
        chooseTransition = rand(); 
    }



//  // choose a pattern
//  switch(chooseFern % NUM_PATTERNS) {
//    case 0: rainbow(); break;
//    case 1: mondrianColors(); break;
////    case 2: trinkleStar(); break;
////    case 3: flood(); break; 
////    case 4: fadeIn(128); break; 
//    // case 5: fadeOut(0); break; 
//    default: rainbow(); break; 
//
//    // TODO
//      //shimmerCenter(centre);
//      //crystalGradient();
//      //muzzLight();
//      //solidHue();
//      //colorToWhiteHue();
//      //complementaryHue();
//
//  }
  
    //Serial.print("Execute"); Serial.println( millis()); 
    //  randomFern(); 
    //  mondrianColors(); 
//    rainbow(); 
    flood(); 
//    trinkleStar(); 

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

unsigned long start = 0; 
void randomFern() {
//  Serial.println(F("Random"));

  unsigned long curr = millis(); 
    if (curr - start > 500) {
        Serial.print(start); Serial.print(" "); Serial.print( millis()); Serial.print(" "); Serial.println( curr); 
      
        start = curr; 
//        if (start < 0) start = 0; 
    
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
  unsigned long curr = millis(); 
  if (curr - start > 500) {
    start = curr; 
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
  unsigned long curr = millis(); 
  if( curr - start > 500) {
    start = curr; 
    for(int i = 0; i < NUM_LEDS; i++) {
      crystalHSV(i, hue++,  255, 255); 
    }
    changeColor();
  }
}


unsigned long start2 = millis(); 
//Basic, will light up 3 at a time
void flood () {
  Serial.println(F("flood"));
  //0-5
  unsigned long curr = millis(); 
  
}

//void flood () {
//  //0-5
//  int randHue = random() % 360;
//  for (int bVal = 255; bVal != 0; bVal -=5 ) { //Slowly turn off after lighting up
//    for (int i = 0; i < NUM_LEDS; i++) { 
//      leds[i].setHSV(randHue, 255, bVal);
//    }
//    FastLED.show(); //Flash
//    delay(20); //Wait 500ms
//  }
//}

//Random light blinks 
void trinkleStar (){
  Serial.println(F("trinkle"));
  //random a led position 
  //binkle for 1 sec then fades in saturation 
  //sleep 1 or 2 second
  unsigned long curr = millis(); 
  if(curr - start > 500) {
    start = curr; 
    int randNum = random(NUM_LEDS+1);
    for (int bVal = 255; bVal != 0; bVal -=5 ){
      //add the colour here
      //fill in here angle, sat, brightness
      fill_solid(target, NUM_LEDS, CRGB::White);
      //leds[randNum].setHSV(42, 255, bVal);
      crystalHSV(randNum, 42, 255, bVal); 
    //    delay(20); // sleep 
    
    //    FastLED.show(); //Blink
    }
    changeColor(); 
  }
//   delay(500); //Wait 500ms
}

//Function for fade
void fadeIn(int hueValue) {
  unsigned long curr = millis(); 
    if(curr - start > 500) {
    for (uint8_t bVal = 0; bVal != 255; bVal +=5 ){
      for (int j = 0; j < NUM_LEDS; ++j) {
        leds[j].setHSV(hueValue, 255, bVal);
  //      delay(1); // sleep 
       }
  //     FastLED.show(); //Blink
    } 
  //  changeColor(); 
  }
}

void fadeOut(int hueValue) {
  unsigned long curr = millis(); 
  if(curr - start > 500) {
    for (uint8_t bVal = 255; bVal >= 0; bVal -=5 ){
      for (int i = 0; i < NUM_LEDS; ++i) {
        leds[i].setHSV(hueValue, 255, bVal);
  //      delay(1); // sleep 
       }
  //     FastLED.show(); //Blink
    }
//    changeColor(); 
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

