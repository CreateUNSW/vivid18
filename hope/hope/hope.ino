#include <FastLED.h>
#include <StandardCplusplus.h>
#include <vector>

// Crystals/LEDs config
//--------------------------------------------------
#define NUM_LEDS 23
#define DATA_PIN 3
#define LED_TYPE WS2811

#define NUM_CHAINS 2
//--------------------------------------------------

CRGB physicalLEDs[NUM_CHAINS][NUM_LEDS];

class Effect;
class LED;
class Chain;


class Effect {
public:
  virtual ~Effect() {}
  virtual void draw(long unsigned int currentTime, Chain& chain) = 0;
};

struct _LedPosition {
	char chainIndex;
	char ledIndex;
};
typedef struct _LedPosition LedPosition;
// Represents a LED in a "virtual" chain
struct LED {
public:
	char num;
	LedPosition* ledPosition;	 // Physical mapping of this virtual LED to real LEDs
	CHSV val;

	LED(CHSV val, char numLEDs)
		:
		val(val),
		num(0)
	{
		ledPosition = (LedPosition*)malloc(numLEDs * sizeof(LedPosition));
	}

	void addPhysicalLED(char chainNum, char ledNum)
	{
		ledPosition[num].chainIndex = chainNum;
		ledPosition[num].ledIndex = ledNum;
		num++;
	}

	void setValue(CHSV inValue)
	{
		this->val = inValue;
	}

	void displayLED()
	{
		for (int i = 0; i < num; i++)
			physicalLEDs[ledPosition[i].chainIndex][ledPosition[i].ledIndex] = val;
	}
};

// Represents a visual chain (virtual chain)
class Chain {
public:
	std::vector<Effect*> effects;
	std::vector<LED> leds;
	//LED* leds;

	// Initialize with number of virtual LEDs
	Chain::Chain(char numLEDs)
	{
		leds.reserve(numLEDs);
		for (int i = 0; i < numLEDs; i++) {
			leds.push_back(LED(CHSV(0, 0, 0), numLEDs));
		}
	};

	// apply effect to color values
	void executeEffects() {
		//Serial.println("executing");
		//Serial.flush();
		int size = leds.size();
		for (int i = 0; i < size; i++) {
			//colors.push_back(&(CHSV(0, 0, 0)));
		}
		for (auto effect : effects) {
			effect->draw(millis(), *this);
		}
		//Serial.println("finished exec");
		//Serial.flush();
	};

	// write color values to LEDs
	void displayLEDs()
	{
		//Serial.println(String((*(colors[0])).hue));
		//FastLED.clear();
		for (int i = 0; i < 23; i++)
		{
			leds[i].displayLED();
		}

		/*for (auto led : leds)
		{
			led.displayLED();
		}*/
		FastLED.show();
	};
};
	

std::vector<Chain> chains;

class Fade : public Effect {
public:
	CHSV toColor;
	int speed;
	fract8 fractspeed;

	virtual void draw(long unsigned int currentTime, Chain& chain) {
		// Serial.println("begin draw");
		//Serial.flush();
		for (auto color : chain.leds) {
			toColor.hue > color.val.hue ? color.val.hue += speed : 1;
			toColor.sat > color.val.sat ? color.val.sat += speed : 1;
			toColor.val > color.val.val ? color.val.val += speed : 1;     
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

	FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(physicalLEDs[0], NUM_LEDS); 

	// initialize one chain with 23 LEDs
	Chain chain = Chain(NUM_LEDS);

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

