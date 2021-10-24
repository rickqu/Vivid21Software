/* 
   OctoWS2811 powered by fastled on Teensy 4.1 and using NativeEthernet
   Code modified from http://doityourselfchristmas.com/forums/showthread.php?55073-Teensy-4-1-E1-31-receiver-controller
   
   Note that you need to download the library from https://github.com/PaulStoffregen/OctoWS2811 as the version in the Arduino 
   library manager does not have all the code we want. 
*/
#include <OctoWS2811.h>
#define USE_OCTOWS2811
#include <FastLED.h>
#include <NativeEthernet.h>
#include <Arduino.h>

#define BOARD_ID 3
#if BOARD_ID < 2 
#error BOARD_ID must be greater than or equal to 1. Board ID 0 and 1 are reserved for ledsim
#endif

EthernetUDP Udp;
uint8_t mac[6]; //this is determined from the Teensy 4.1 board.  
IPAddress ip(10, 1, 2, BOARD_ID); //static ip address of Teensy board
#define UDP_PORT 5050 //RPi UDP Port Number

// missing pins 5-8.
#define NUM_1 115
#define NUM_2 85 + 69 + 60
#define NUM_3 64 + 116
#define NUM_4 117

#define NUM_PINS 4
#define NUM_TOTAL NUM_1 + NUM_2 + NUM_3 + NUM_4
byte pinList[NUM_PINS] = {2, 6, 14, 20};
#define LED_TYPE OCTOWS2811

#define ETHERNET_BUFFER NUM_TOTAL * 3
uint8_t ethernetBuffer[ETHERNET_BUFFER]

CRGB * rgbarray = (CRGB*)ethernetBuffer;

void pixelDisplay() {
  FastLED.Show();
}

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
  CRGB red = CRGB(0xFF, 0x00, 0x00);
  setAllColor(red);
  FastLED.show();
}

void setAllColor(CRGB color) {
  for (int led = 0; led < NUM_TOTAL; led++) {
    rgbarray[led] = color;
  }

  FastLED.show();
}

//define mac address of board
void teensyMAC(uint8_t *mac) {

  static char teensyMac[23];
  
  #if defined(HW_OCOTP_MAC1) && defined(HW_OCOTP_MAC0)
    Serial.println("using HW_OCOTP_MAC* - see https://forum.pjrc.com/threads/57595-Serial-amp-MAC-Address-Teensy-4-0");
    for(uint8_t by=0; by<2; by++) mac[by]=(HW_OCOTP_MAC1 >> ((1-by)*8)) & 0xFF;
    for(uint8_t by=0; by<4; by++) mac[by+2]=(HW_OCOTP_MAC0 >> ((3-by)*8)) & 0xFF;

    #define MAC_OK

  #else
    
    mac[0] = 0x04;
    mac[1] = 0xE9;
    mac[2] = 0xE5;

    uint32_t SN=0;
    __disable_irq();
    
    #if defined(HAS_KINETIS_FLASH_FTFA) || defined(HAS_KINETIS_FLASH_FTFL)
      Serial.println("using FTFL_FSTAT_FTFA - vis teensyID.h - see https://github.com/sstaub/TeensyID/blob/master/TeensyID.h");
      
      FTFL_FSTAT = FTFL_FSTAT_RDCOLERR | FTFL_FSTAT_ACCERR | FTFL_FSTAT_FPVIOL;
      FTFL_FCCOB0 = 0x41;
      FTFL_FCCOB1 = 15;
      FTFL_FSTAT = FTFL_FSTAT_CCIF;
      while (!(FTFL_FSTAT & FTFL_FSTAT_CCIF)) ; // wait
      SN = *(uint32_t *)&FTFL_FCCOB7;

      #define MAC_OK
      
    #elif defined(HAS_KINETIS_FLASH_FTFE)
      Serial.println("using FTFL_FSTAT_FTFE - vis teensyID.h - see https://github.com/sstaub/TeensyID/blob/master/TeensyID.h");
      
      kinetis_hsrun_disable();
      FTFL_FSTAT = FTFL_FSTAT_RDCOLERR | FTFL_FSTAT_ACCERR | FTFL_FSTAT_FPVIOL;
      *(uint32_t *)&FTFL_FCCOB3 = 0x41070000;
      FTFL_FSTAT = FTFL_FSTAT_CCIF;
      while (!(FTFL_FSTAT & FTFL_FSTAT_CCIF)) ; // wait
      SN = *(uint32_t *)&FTFL_FCCOBB;
      kinetis_hsrun_enable();

      #define MAC_OK
      
    #endif
    
    __enable_irq();

    for(uint8_t by=0; by<3; by++) mac[by+3]=(SN >> ((2-by)*8)) & 0xFF;

  #endif

  #ifdef MAC_OK
    sprintf(teensyMac, "MAC: %02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
    Serial.println(teensyMac);
  #else
    Serial.println("ERROR: could not get MAC");
  #endif
}

void setup() {

  Serial.begin(115200);
  delay(10);

  teensyMAC(mac);
 
  static char teensyMac[23];
  sprintf(teensyMac, "MAC: %02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
  Serial.print("MAC: ");
  Serial.println(teensyMac);
  
  Ethernet.begin(mac, ip);
  
  Serial.print("IP Address: ");
  Serial.println(Ethernet.localIP());
  Serial.print("MAC: ");
  Serial.println(teensyMac);
 
  Udp.begin(UDP_PORT);
  
  #if NUM_1 > 0
    pinMode(PIN_1, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_1, ORDER_1>(rgbarray, NUM_1);
  #endif
  #if NUM_2 > 0
    pinMode(PIN_2, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_2, ORDER_2>(rgbarray + NUM_1, NUM_2);
  #endif
  #if NUM_3 > 0
    pinMode(PIN_3, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_3, ORDER_3>(rgbarray + NUM_1 + NUM_2, NUM_3);
  #endif
  #if NUM_4 > 0
    pinMode(PIN_4, OUTPUT);
    FastLED.addLeds<LED_TYPE, PIN_4, ORDER_4>(rgbarray + NUM_1 + NUM_2 + NUM_3, NUM_4);
  #endif 
  FastLED.setBrightness(50);
  
  blinkInit();
  Serial.println("Test Sequence Complete");

} //end setup

void loop() {
  //Process packets
   
  int packetSize = Udp.parsePacket(); //Read UDP packet count
    
  if(packetSize){
    Serial.println(packetSize);
    Udp.read(rgbarray,ETHERNET_BUFFER); //read UDP packet
    pixelDisplay(); //process data function
  }  
}
