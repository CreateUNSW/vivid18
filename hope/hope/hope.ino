#include <FastLED.h>
#include <StandardCplusplus.h>
#include <vector>

// Crystals/LEDs config
//--------------------------------------------------
#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811

#define NUM_CHAINS 16
//--------------------------------------------------

CRGB leds[NUM_LEDS];


using namespace std;

class Effect {
public:
  virtual ~Effect() {}
  virtual void draw(int currentTime, vector<CHSV*>* colors) = 0;
};

class Chain {
public:
  vector<CHSV*> colors;
  vector<Effect*> effects;

  Chain()
  {
    colors = vector<CHSV*>();
    for (int i = 0; i < NUM_LEDS; i++) {
      colors.push_back(new CHSV(0, 0, 0));
      ((leds[i])) &= 255;
    }
    //Serial.println("create chain");
    //Serial.println(String(colors.size()));
    //Serial.flush();
  }

  // apply effect to color values
  void executeEffects() {
    //Serial.println("executing");
    //Serial.flush();
    int size = colors.size();
    for (int i = 0; i < size; i++) {
      //colors.push_back(&(CHSV(0, 0, 0)));
    }
    for (auto effect : effects) {
      effect->draw(millis(), &colors);
    }
    //Serial.println("finished exec");
    //Serial.flush();
  };

  // write color values to LEDs
  void displayLEDs()
  {    
    //Serial.println(String((*(colors[0])).hue));
    //FastLED.clear();
    for(int i = 0; i < NUM_LEDS; i++)
    {
      leds[i] = (*(colors[i]));
    }
    FastLED.show();
  }
};

vector<Chain> chains;

class Fade : public Effect {
public:
  CHSV toColor;
  int speed;
  fract8 fractspeed;

  virtual void draw(int currentTime, vector<CHSV*>* colors) {
   // Serial.println("begin draw");
    //Serial.flush();
    for (auto color : *colors) {
      toColor.hue > color->hue ? color->hue += speed : 1;
      toColor.sat > color->sat ? color->sat += speed : 1;
      toColor.val > color->val ? color->val += speed : 1;     
    }
    //Serial.println("COLORS:");
    //Serial.println(String(toColor.val));
    //Serial.println(String(((*colors)[0])->val));
    //Serial.println("finish draw");
    Serial.flush();
  };
};

void setup() {
  Serial.begin(9600);
  Serial.println("boot");
  delay(1000);

  FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(leds, NUM_LEDS); 

  // initialize one chain with 23 LEDs
  Chain chain;

  // Add fade effect to chain
  Fade* fade = new Fade();
  fade->toColor = CHSV(255, 255, 255);
  fade->speed = 10;
  fade->fractspeed = 0.01f;
  
  chain.effects.push_back(fade);

  chains.push_back(chain); 
  delay(1000);
  Serial.println("finished setup");
  Serial.flush();
}

void loop() {
  unsigned long start = micros();
  for (auto chain : chains) {
    chain.executeEffects();
    for (auto col : chain.colors) {
    }
    chain.displayLEDs();
    //Serial.println("chaining");
    //Serial.flush();
  }
  unsigned long end = micros();
  unsigned long delta = end - start;
  Serial.println(String(delta));
  Serial.flush();

  delay(30);
}

