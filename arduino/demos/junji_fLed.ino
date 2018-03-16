#include "FastLED.h"
#include "QList.h"

// How many leds in your strip?
#define NUM_LEDS 50
#define NUM_CHAINS 6

// For led chips like Neopixels, which have a data line, ground, and power, you just
// need to define DATA_PIN.  For led chipsets that are SPI based (four wires - data, clock,
// ground, and power), like the LPD8806 define both DATA_PIN and CLOCK_PIN
#define DATA_PIN 3
#define CLOCK_PIN 13

#define FORWARD 0
#define BACKWARD 1

// Define the array of leds
CRGB leds[NUM_CHAINS][NUM_LEDS];
CRGB samp[NUM_LEDS]; 
uint8_t gHue = 0;


void setup() { 

      Serial.begin(9600); 
      // Uncomment/edit one of the following lines for your leds arrangement.
      // FastLED.addLeds<TM1803, DATA_PIN, RGB>(leds, NUM_LEDS);
      // FastLED.addLeds<TM1804, DATA_PIN, RGB>(leds, NUM_LEDS);
      // FastLED.addLeds<TM1809, DATA_PIN, RGB>(leds, NUM_LEDS);
      FastLED.addLeds<WS2811, DATA_PIN, RGB>(samp, NUM_LEDS);
      FastLED.addLeds<WS2811, 4, RGB>(leds[1], NUM_LEDS);
      FastLED.addLeds<WS2811, 5, RGB>(leds[2], NUM_LEDS);
      FastLED.addLeds<WS2811, 6, RGB>(leds[3], NUM_LEDS);
      FastLED.addLeds<WS2811, 7, RGB>(leds[4], NUM_LEDS);
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

// list of patterns to implement
enum pattern_l {
  FLOOD, FADE, PULSE, TWINKLE, WIPE
};

// pattern type stipulating LED chain address, target colour, execution time, direction etc. 
//typedef struct _pattern_t pattern_t; 

class pattern_t {
    public: 
    pattern_l pattern; 
    int direction; 
    int speed; 
    
    CRGB*  ledAddress; 
    int    numLeds; 
    CRGB   colorA; 
    CRGB   colorB; 
    int    start; 
//    int   execTimeEnd; 
    int    period; 

    pattern_t(
      pattern_l l_pattern = WIPE, 
      int l_direction = FORWARD, 
      int l_speed = 250, 
      CRGB* l_ledAddress = samp, 
      int l_numLeds = NUM_LEDS, 
      CRGB l_colorA = CRGB::Red, 
      CRGB l_colorB = CRGB::Green, 
      unsigned long l_start = millis(), 
      unsigned long l_period = 1000
    ) {
      pattern = l_pattern; 
      direction = l_direction; 
      speed = l_speed; 
      ledAddress = l_ledAddress; 
      numLeds = l_numLeds; 
      colorA = l_colorA; 
      colorB = l_colorB; 
      start = l_start;  
      period = l_period; 
    };
    
}; 

QList <pattern_t> jobQ; 

//
//concurrentJobQueue = [
//    flood,
//    fade,
//    flood,
//    flood,
//]

  
void executeEffect(struct pattern_t pattern, unsigned long since) {
  switch(pattern.pattern) { 
    case FLOOD: 
    break; 
    case FADE: 
    break; 
    case PULSE: 
    break; 
    case TWINKLE: 
    break; 
    case WIPE: 
    {
//      unsigned long currentTime = millis(); 
        Serial.print("start: "); Serial.println( pattern.start); 
        while (since < pattern.start + pattern.period) {
          
          for(int i=0; i<pattern.numLeds; i++){
          if(pattern.direction == FORWARD){
            pattern.ledAddress[i] = pattern.colorA;
          }
          else{
            pattern.ledAddress[NUM_LEDS-1-i] = pattern.colorA;
          }
          FastLED.show();
          delay(pattern.speed); 
          
          }

          for(int i=0; i<pattern.numLeds; i++){
          if(pattern.direction == FORWARD){
            pattern.ledAddress[i] = CRGB(0,0,0);
          }
          else{
            pattern.ledAddress[NUM_LEDS-1-i] = CRGB(0,0,0);
          }
          FastLED.show();
          delay(pattern.speed); 
          
          }
          since = millis(); 
          Serial.print("end: "); Serial.println( pattern.start + pattern.period); 
          Serial.print("clock: "); Serial.println( since) ;
    }
    break; 
  }

  }
    
}

//render() {
//    since = delta()
//    // something
//
//    lighting = [] // colors/lights for this render
//
//    for each job in concurrentJobQueue {
//        job.executeEffect(lighting, since)
//    }
//}

void render() {
    unsigned long since = millis(); 
    int i = 0; 
    while (jobQ.size() > 0) {
        since = millis(); 
        pattern_t cpat = jobQ.front();
        jobQ.pop_front(); 
        Serial.println(i); 
        executeEffect(cpat, since); 
        i++; 
    }
}

void loop() {
//  static uint8_t hue = 0;
//  FastLED.showColor(CHSV(hue++, 255, 255));
//  delay(1);

    pattern_t wipe1(
      WIPE, 
      FORWARD, 
      10, 
      samp, 
      NUM_LEDS, 
      CRGB::Blue, 
      CRGB::Red, 
      millis(), 
      1000
    );

    pattern_t wipe2(
      WIPE, 
      FORWARD, 
      10, 
      samp, 
      NUM_LEDS, 
      CRGB::Green, 
      CRGB::Green, 
      millis(),  
      1000 
    );

     pattern_t wipe3(
      WIPE, 
      FORWARD, 
      10, 
      samp, 
      NUM_LEDS, 
      CRGB::Blue, 
      CRGB::Blue, 
      millis(), 
      1000 
    );

    jobQ.push_back(wipe1); 
    jobQ.push_back(wipe2); 
    jobQ.push_back(wipe3); 

    render(); 

    delay(5000); 
  
}
