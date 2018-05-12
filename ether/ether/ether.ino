#include <EtherCard.h>
#include <IPAddress.h>

#include <Wire.h>
#include <FastLED.h>

uint8_t Ethernet::buffer[500]; // configure buffer size to 700 octets

#define NUM_LEDS 40
#define DATA_PIN 5
#define LED_TYPE UCS1903

CRGB leds[NUM_LEDS];
CRGB ledA[NUM_LEDS];
CRGB ledB[NUM_LEDS];

void receive(uint16_t dest_port, uint8_t src_ip[IP_LEN], uint16_t src_port, const char *data, uint16_t len) {
  Serial.println("data receieved!");
  //Serial.println(data);
  memcpy(leds, data, len);
  FastLED.show();
}

void setup() {
  Serial.begin(115200);
  //Serial.setTimeout(20);
  delay(3000);
  Serial.println("setting up...");

  static uint8_t mymac[] = { 0x74,0x69,0x69,0x2D,0x30,0x31 }; // define (unique on LAN) hardware (MAC) address
  uint8_t nFirmwareVersion = ether.begin(sizeof Ethernet::buffer, mymac, 3);
  if(0 == nFirmwareVersion) {
    //Serial.println("failure");
    return;
  }

  Serial.println("ethernet started");

  const static uint8_t ip[] = {192,168,2,11};
  const static uint8_t gw[] = {192,168,2,1};
  if (!ether.dhcpSetup()) {
      Serial.println("failure static ip");
      return;
      // handle failure to configure static IP address (current implementation always returns true!)
  }

  //Serial.println("ethernet setup");

  ether.udpServerListenOnPort(&receive, 6969);

  pinMode(DATA_PIN, OUTPUT);
  FastLED.addLeds<LED_TYPE, DATA_PIN, RGB>(leds, NUM_LEDS); 
  Serial.println("ready");
}

void loop() {
  ether.packetLoop(ether.packetReceive());
}

