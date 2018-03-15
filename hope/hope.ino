#include <FastLED.h>
#include <StandardCplusplus.h>
#include <vector>

// Crystals/LEDs config
//--------------------------------------------------
#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811
//--------------------------------------------------

CRGB leds[NUM_LEDS];

using namespace std;

class Effect {
public:
  virtual ~Effect() {}
  virtual void draw(int currentTime, vector<CHSV*> colors) = 0;
  virtual void hello() = 0;
};

class Chain {
public:
  vector<CHSV*> colors;
  vector<Effect*> effects;

  // apply effect to color values
  void executeEffects() {
    int size = colors.size();
    Serial.println(String(size) + " colors");
    Serial.flush();
    colors.clear();
    for (int i = 0; i < size; i++) {
      colors.push_back(&(CHSV(0, 0, 0)));
    }
    for (auto effect : effects) {
      Serial.println("HELLO_0");
      Serial.flush();
      effect->hello();
      effect->draw(millis(), colors);
      Serial.println("HERE_1");
      Serial.flush();
    }
  };

  // write color values to LEDs
  void displayLEDs()
  {
    
  }
};

vector<Chain> chains;

class Fade : public Effect {
public:
  CHSV toColor;
  int speed;

  virtual void hello() {
    Serial.println("Fade effect");
    Serial.flush();
  };

  virtual void draw(int currentTime, vector<CHSV*> colors) {
    Serial.println("Drawing..");
    Serial.flush();
    for (auto color : colors) {
      color->hue = toColor.hue;
      color->sat = toColor.sat;
      color->val = toColor.val;
    }
   Serial.println("Finished drawing...");
   Serial.flush();
  };
};

void setup() {
  Serial.begin(9600);
  Serial.println("boot");
  delay(1000);

  FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(leds, NUM_LEDS); 
  
  Chain chain;
  chain.colors.push_back(&(CHSV(0, 0, 0)));
  chain.colors.push_back(&(CHSV(1, 0, 0)));
  chain.colors.push_back(&(CHSV(2, 0, 0)));

  Fade *fade = new Fade();
  fade->toColor = CHSV(160, 255, 255);
  fade->speed = 1;

  chain.effects.push_back(fade);

  chains.push_back(chain); 
  delay(1000);
  Serial.println("finished setup");
  Serial.flush();
}

void loop() {
  Serial.println("begin loop");
  for (auto chain : chains) {
    chain.executeEffects();
    for (auto col : chain.colors) {
      Serial.println("for loop");
      Serial.println(String(col->val));
      Serial.flush();
    }
  }
  Serial.println("boop");
  Serial.flush();

  delay(1000);
}

