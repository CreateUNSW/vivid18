#include <FastLED.h>
#include <StandardCplusplus.h>
#include <vector>

using namespace std;

class Effect {
public:
	virtual ~Effect() {}
	virtual void draw(int currentTime, vector<CHSV*> colors) = 0;
};

class Chain {
public:
	vector<CHSV*> colors;
	vector<Effect*> effects;
	void executeEffects() {
		int size = colors.size();
		colors.clear();
		for (int i = 0; i < size; i++) {
			colors.push_back(&(CHSV(0, 0, 0)));
		}
		for (auto effect : effects) {
			effect->draw(millis(), colors);
		}
	};
};

vector<Chain> chains;

class Fade : public Effect {
public:
	CHSV toColor;
	int speed;

	virtual void draw(int currentTime, vector<CHSV*> colors) {
		for (auto color : colors) {
			color->hue = toColor.hue;
			color->sat = toColor.sat;
			color->val = toColor.val;
		}
	};
};

void setup() {
	Chain chain;
	chain.colors.push_back(&(CHSV(0, 0, 0)));

	Fade fade;
	fade.toColor = CHSV(160, 255, 255);
	fade.speed = 1;

	chain.effects.push_back(&fade);

	chains.push_back(chain);


	Serial.begin(9600);
}

void loop() {
	for (auto chain : chains) {
		chain.executeEffects();
		for (auto col : chain.colors) {
			Serial.println(col->val);
		}
	}
	delay(1000);
}
