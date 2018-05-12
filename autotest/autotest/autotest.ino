#include <EtherCard.h>
#include <IPAddress.h>
#include <enc28j60.h>

// These details MUST be unique per Arduino device. If they are not, network issues will occur.
const static uint8_t cfgIP[] = {192,168,2,11};
const static uint8_t cfgMAC[] = {0x74,0x69,0x69,0x2D,0x30,0x11};

const static uint8_t cfgGateway[] = {192,168,2,1};

uint8_t Ethernet::buffer[600];

int totalReceived = 0;
char hasErrors = 0;
char sentAck = 0;
char receivedAck = 0;

void receive(uint16_t destPort, uint8_t srcIP[IP_LEN], uint16_t srcPort, const char *data, uint16_t len) {
  totalReceived++;
  
  if (len != 500) {
    hasErrors |= 1;
    return;
  }

  receivedAck |= (data[0] == 1);
  
  if (data[0] == 0 || data[0] == 1) {
    for (int i = 1; i < 500; i++) {
      hasErrors |= (data[i] != char(i % 256));
    }
  } else {
    hasErrors |= 1;
  }

  if (!sentAck) {
    char ok[3] = "OK";
    ether.sendUdp(ok, 2, 5151, srcIP, srcPort);
    ether.sendUdp(ok, 2, 5151, srcIP, srcPort);
    ether.sendUdp(ok, 2, 5151, srcIP, srcPort);
    ether.sendUdp(ok, 2, 5151, srcIP, srcPort);
    ether.sendUdp(ok, 2, 5151, srcIP, srcPort);
    sentAck = 1;
  }
}

void wait() {
  while (!Serial.available()) { delay(100); };
  while (Serial.available()) { Serial.read(); };
}

void testNetwork() {
  Serial.println("ETHERNET");

  uint8_t nFirmwareVersion = ether.begin(sizeof Ethernet::buffer, cfgMAC, 3);
  if (0 == nFirmwareVersion) {
    Serial.println("FAIL BEGIN");
    return;
  }

  if (!ether.staticSetup(cfgIP, cfgGateway)) {
    Serial.println("FAIL SETUP");
    return;
  }

  if (!ENC28J60::isLinkUp()) {
    Serial.println("FAIL LINK");
    return;
  }

  Serial.println("OK ETHERNET");

  ether.udpServerListenOnPort(&receive, 5151);

  Serial.println("NETWORK_PREPARE");
  while (!Serial.available()) { ether.packetLoop(ether.packetReceive()); };
  while (Serial.available()) { Serial.read(); };
  Serial.println("NETWORK_START");
  
  totalReceived = 0;
  hasErrors = 0;
  sentAck = 0;
  receivedAck = 0;
  
  unsigned long startTime = millis();
  while ((millis() - startTime) < 3000) {
    ether.packetLoop(ether.packetReceive());
  }

  Serial.println(totalReceived);

  if (!hasErrors) {
    Serial.println("FAIL CORRUPT");
  } else {
    Serial.println("OK CORRUPT");
  }

  if (!receivedAck) {
    Serial.println("FAIL ACK");
  } else {
    Serial.println("OK ACK");
  }
}

void testPins() {
  Serial.println("PINS");

  pinMode(9, INPUT);  
  pinMode(10, OUTPUT);
  analogWrite(10, 0);
  delay(100);

  char failure = 0;
  failure |= digitalRead(9) == HIGH;

  analogWrite(10, 255);
  delay(100);
  
  failure |= digitalRead(9) == LOW;

  if (!failure) {
    Serial.println("OK 10");
  } else {
    Serial.println("FAIL 10");
  }

  pinMode(10, INPUT);
  pinMode(9, OUTPUT);
  analogWrite(9, 0);
  delay(100);

  failure = 0;
  failure |= digitalRead(10) == HIGH;

  analogWrite(9, 255);
  delay(100);
  
  failure |= digitalRead(10) == LOW;

  if (!failure) {
    Serial.println("OK 9");
  } else {
    Serial.println("FAIL 9");
  }

  failure = 0;
  
  pinMode(5, INPUT);
  pinMode(6, OUTPUT);
  analogWrite(6, 0);
  delay(100);

  failure |= digitalRead(5) == HIGH;

  analogWrite(6, 255);
  delay(100);
  
  failure |= digitalRead(5) == LOW;

  if (!failure) {
    Serial.println("OK 6");
  } else {
    Serial.println("FAIL 6");
  }

  pinMode(6, INPUT);
  pinMode(5, OUTPUT);
  analogWrite(5, 0);
  delay(100);

  failure = 0;
  failure |= digitalRead(6) == HIGH;

  analogWrite(5, 255);
  delay(100);
  
  failure |= digitalRead(6) == LOW;

  if (!failure) {
    Serial.println("OK 5");
  } else {
    Serial.println("FAIL 5");
  }
}

void setup() {
  Serial.begin(115200);
  wait();
  Serial.println("HELLO");
  wait();

  testNetwork();
  testPins();
  
  Serial.println("END");
}

void loop() {
  delay(10000000);
}

