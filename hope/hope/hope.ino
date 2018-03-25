#include <FastLED.h>
#include <StandardCplusplus.h>
#include <vector>

// Crystals/LEDs config
//--------------------------------------------------
#define NUM_CHAINS 2
#define LED_TYPE WS2811
#define FastLED_ORDER RGB
#define MAX_NUM_LEDS 50

#define NUM_LEDS_1 50
#define DATA_PIN_1 3
//--------------------------------------------------

CRGB physicalLEDs[NUM_CHAINS][MAX_NUM_LEDS];

class Effect;
class LED;
class Chain;


class Effect {
public:
  virtual ~Effect() {}
  virtual void draw(long unsigned int currentTime, Chain* chain) = 0;
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

	// map this virtual led to another physical led
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

	// copy virtual led value to all physical led mapped values
	// TODO: maybe use physical LED array to store value?
	void writeLED()
	{
		for (int i = 0; i < num; i++)
			physicalLEDs[ledPosition[i].chainIndex][ledPosition[i].ledIndex] = val;
		//physicalLEDs[ledPosition[0].chainIndex][ledPosition[0].ledIndex] += CHSV(1, 1, 1);
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
		for (auto effect : effects) {
			effect->draw(millis(), this);
		}
		//Serial.println("finished exec");
		//Serial.flush();
	};

	// Read a config array
	void setConfiguration(char** config) {
		for (int i = 0; i < leds.size(); i++) {
			int size = config[i][0];
			for (int j = 0; j < size; j++) {
				leds.at(i).addPhysicalLED(config[i][1 + 2 * j], config[i][1 + 2 * j + 1]);
			}
		}
	}

	// write color values from virtual LEDs to physical
	void displayLEDs()
	{
		//Serial.println(String((*(colors[0])).hue));
		//FastLED.clear();

		for (auto led : leds)
		{
			led.writeLED();
		}
		FastLED.show();
	};
};
	

std::vector<Chain> chains;

class Fade : public Effect {
public:
	CHSV toColor;
	int speed;
	fract8 fractspeed;

	virtual void draw(long unsigned int currentTime, Chain* chain) {
		// Serial.println("begin draw");
		//Serial.flush();
		for (auto led : chain->leds) {
			toColor.hue > led.val.hue ? led.val.hue += speed : 1;
			toColor.sat > led.val.sat ? led.val.sat += speed : 1;
			toColor.val > led.val.val ? led.val.val += speed : 1;     
		}
		//Serial.println("COLORS:");
		//Serial.println(String(toColor.val));
		//Serial.println(String(((*colors)[0])->val));
		//Serial.println("finish draw");
		Serial.flush();
	};
};

void configLEDs()
{
	// Add physical LEDs to FastLED
	FastLED.addLeds<LED_TYPE, DATA_PIN_1, FastLED_ORDER>(physicalLEDs[0], NUM_LEDS_1);

	// initialize virtual chains
	Chain chain = Chain(2);
	// For each virtual led, give a list of physical leds (chain, led) that it maps to) 
	// (first element indicates how many physical leds that led maps to)
	const int numVirtualLEDs = 2;
	const int maxLEDs = 2;
	char config1[numVirtualLEDs][maxLEDs*2 + 1] = 
													{ {2, 0,0, 1,0},
													{2, 0,1, 1,1} };
	chain.setConfiguration((char**)config1);
	chains.push_back(chain);
}

void setup() {
	Serial.begin(9600);
	Serial.println("boot");
	delay(1000);

	configLEDs();

	delay(1000);
	Serial.println("finished setup");
	Serial.flush();

	// Just for testing:

	// Add fade effect to chain
	Fade* fade = new Fade();
	fade->toColor = CHSV(255, 255, 255);
	fade->speed = 10;
	fade->fractspeed = 0.01f;

	chains[0].effects.push_back(fade);
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

