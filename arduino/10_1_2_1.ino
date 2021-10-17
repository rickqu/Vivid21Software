#include <EtherCard.h>
#include <IPAddress.h>
#include <enc28j60.h>
#include <FastLED.h>
#include "test.cpp"

#define BOARD_ID 1

#if BOARD_ID < 1 
#error BOARD_ID must be greater than or equal to 1. Board ID 0 is reserved for ledsim
#endif

#define ORDER_M20 RBG
#define ORDER_M30 RBG

#define ORDER_1 ORDER_M30
#define ORDER_2 ORDER_M30
#define ORDER_3 ORDER_M30
#define ORDER_4 ORDER_M30

// missing pins 5-8.
#define NUM_1 36 +  102
#define NUM_2 51
#define NUM_3 104 + 94
#define NUM_4 53

#define NUM_PINS 4
#define NUM_TOTAL NUM_1 + NUM_2 + NUM_3 + NUM_4

// TODO!
#define LED_TYPE UCS1903

#define LISTEN_PORT 5151
#define SERVER_PORT 5050

const static uint8_t cfgIP[] = {10,1,2,BOARD_ID};
const static uint8_t cfgMAC[] = {0x74,0x69,0x69,0x2D,0x30,BOARD_ID};

// ledsim is the gateway
const static uint8_t cfgGateway[] = {10,1,2,0};

// not indicative of final board
#define PIN_1 2
#define PIN_2 6
#define PIN_3 14
#define PIN_4 20

unsigned long lastReceive = 0;
unsigned long nextPing = 1000;
uint8_t countReply[] = {'C', 0};
uint8_t ackReply[] = {'A', 1};

//unsigned long curMillis = 0;
unsigned long curMillis = SOME_CONSTANT;

uint8_t Ethernet::buffer[1200];
CRGB *leds = (CRGB*)(Ethernet::buffer + UDP_DATA_P);

// People say this is a way to reset an Arduino in software.
// But really, this is causing a segfault lol.
void(* resetFunc) (void);

void blinkInit() {
  RXLED0;
  if ((curMillis / 1000) % 2 == 0) {
    TXLED1;
  } else {
    TXLED0;
  }

  flashRed();
}

void flashRed() {
  int frame = (curMillis / 100) % 50;
  CRGB red = CRGB(0xFF, 0x00, 0x00);
  CRGB off = CRGB(0x00, 0x00, 0x00);

  if (frame < 2) {
    setAllColor(red);
  } else if (frame < 4) {
    setAllColor(off);
  } else if (frame < 6) {
    setAllColor(red);
  } else if (frame < 8) {
    setAllColor(off);
  } else if (frame < 10) {
    setAllColor(red);
  } else if (frame < 12) {
    setAllColor(off);

  } else if (frame < 17) {
    setAllColor(red);
  } else if (frame < 19) {
    setAllColor(off);
  } else if (frame < 24) {
    setAllColor(red);
  } else if (frame < 26) {
    setAllColor(off);
  } else if (frame < 31) {
    setAllColor(red);
  } else if (frame < 33) {
    setAllColor(off);

  } else if (frame < 35) {
    setAllColor(red);
  } else if (frame < 37) {
    setAllColor(off);
  } else if (frame < 39) {
    setAllColor(red);
  } else if (frame < 41) {
    setAllColor(off);
  } else if (frame < 43) {
    setAllColor(red);
  } else {
    setAllColor(off);
  }
}

void setAllColor(CRGB color) {
  for (int led = 0; led < NUM_TOTAL; led++) {
    leds[led] = color;
  }

  FastLED.show();
}

void startNetwork() {
  while (1) {
    uint8_t nFirmwareVersion = ether.begin(sizeof Ethernet::buffer, cfgMAC, 3);
    if (0 == nFirmwareVersion) {
      curMillis = millis();
      blinkInit();
      continue;
    }

    if (!ether.staticSetup(cfgIP, cfgGateway)) {
      RXLED0;
      TXLED0;
      delay(10000000);
    }

    while (ether.clientWaitingGw()) {
      for (int i = 0; i < 10; i++) {
        ether.packetLoop(ether.packetReceive());
      }
      curMillis = millis();
      blinkInit();
    }

    return;
  }
}

void setupLEDs() {
  #if NUM_1 > 0
    pinMode(PIN_1, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_1, ORDER_1>(leds, NUM_1);
  #endif
  #if NUM_2 > 0
    pinMode(PIN_2, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_2, ORDER_2>(leds + NUM_1, NUM_2);
  #endif
  #if NUM_3 > 0
    pinMode(PIN_3, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_3, ORDER_3>(leds + NUM_1 + NUM_2, NUM_3);
  #endif
  #if NUM_4 > 0
    pinMode(PIN_4, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_4, ORDER_4>(leds + NUM_1 + NUM_2 + NUM_3, NUM_4);
  #endif
}

void receive(uint16_t destPort, uint8_t srcIP[IP_LEN], uint16_t srcPort, const char *data, uint16_t len) {
  if (srcPort != SERVER_PORT) {
    return;
  }

  if (len == 2) {
    if (data[0] == 'R') {
      resetFunc();
      return;
    } else if (data[0] == 'S') {
      ackReply[1] = data[1];
      ether.sendUdp((char*)(&ackReply), 2, LISTEN_PORT, cfgGateway, SERVER_PORT);
      return;
    }
  }

  FastLED.show();
  lastReceive = curMillis;
  countReply[1]++;
}

void setup() {
  setupLEDs();
  startNetwork();

  ether.udpServerListenOnPort(&receive, LISTEN_PORT);
}

void loop() {
  for (int i = 0; i < 10; i++) {
    ether.packetLoop(ether.packetReceive());
  }

  curMillis = millis();

  if (curMillis >= nextPing) {
    ether.sendUdp((char*)(&countReply), 2, LISTEN_PORT, cfgGateway, SERVER_PORT);
    countReply[1] = 0;
    nextPing += 2000;
  }

  if (!ENC28J60::isLinkUp()) {
    blinkInit();
  } else if (curMillis - lastReceive < 1000) {
    RXLED1;
    TXLED1;
  } else {
    RXLED0;
    TXLED1;
    flashRed();
  }
}
