/* 
   OctoWS2811 powered by fastled on Teensy 4.1 and using NativeEthernet
   Code modified from http://doityourselfchristmas.com/forums/showthread.php?55073-Teensy-4-1-E1-31-receiver-controller
   
   Note that you need to download the library from https://github.com/PaulStoffregen/OctoWS2811 as the version in the Arduino 
   library manager does not have all the code we want. 
*/

#include <OctoWS2811.h>
#include <FastLED.h>
#include <NativeEthernet.h>
#include <Arduino.h>

//change #defines where necessary
#define ETHERNET_SOCKET_SIZE 8192
#define ETHERNET_BUFFER 6490
#define NUM_LEDS_PER_STRIP 250
#define NUM_STRIPS 8
#define NUM_LEDS_OVERALL 1017

//Number of LEDS in each chain
#define two 150
#define seven 106
#define eight 199
#define fourteen 187
#define five 0
#define six 209
#define twenty 107
#define twentyone 85

//ethernet and timing variables
uint8_t packetBuffer[ETHERNET_BUFFER];
int Counter = 0;
int c = 0;
float fps = 0;
unsigned long currentMillis = 0;
unsigned long previousMillis = 0;

#define BOARD_ID 1
#if BOARD_ID < 1
#error BOARD_ID must be greater than or equal to 1. Board ID 0 is reserved for ledsim
#endif

EthernetUDP Udp;
uint8_t mac[6]; //this is determined from the Teensy 4.1 board.  
//IPAddress ip(192, 168, 0, 100); //static ip address of Teensy board (J-og)
IPAddress ip(169, 254, 2, 1); //static ip address of Teensy board (J-20211025)
//IPAddress ip(192, 168, 6, 1); //static ip address of Teensy board (R)
//#define UDP_PORT 8888 //RPi UDP Port Number (J-og)
#define UDP_PORT 5151 //RPi UDP Port Number (R)

//set up OctoWS2811 and FastLED for Teensy 4.1
int unsigned NUM_LEDS = NUM_STRIPS * NUM_LEDS_PER_STRIP;
//byte pinList[NUM_STRIPS] = {6, 20, 21, 5, 2, 14, 7, 8};
byte pinList[NUM_STRIPS] = {2,14,7,8,6,20,21,5};
/*generic same number of leds in each strip*/
CRGB rgbarray[NUM_STRIPS * NUM_LEDS_PER_STRIP];
/*Using actual number of leds*/
//CRGB rgbarray[NUM_LEDS_OVERALL];

// These buffers need to be large enough for all the pixels.
// The total number of pixels is "ledsPerStrip * numPins".
// Each pixel needs 3 bytes, so multiply by 3.  An "int" is
// 4 bytes, so divide by 4.  The array is created using "int"
// so the compiler will align it to 32 bit memory.
/* Generic same number of LEDS per pin*/
DMAMEM int displayMemory[NUM_STRIPS * NUM_LEDS_PER_STRIP * 3 / 4];
int drawingMemory[NUM_STRIPS * NUM_LEDS_PER_STRIP * 3 / 4];
/*Using actual number of leds*/
//DMAMEM int displayMemory[NUM_LEDS_OVERALL * 3 / 4];
//int drawingMemory[NUM_LEDS_OVERALL * 3 / 4];
/* Keep the intialisation below to the same maximum number of leds for each strip.
 *  All it does is give a lot of unused leds on each pin
 */
OctoWS2811 octo(NUM_LEDS_PER_STRIP, displayMemory, drawingMemory, WS2811_GRB | WS2811_800kHz, NUM_STRIPS, pinList);

/*---------------------------Custom Class---------------------------------*/
//custom c++ classes and templates for OctoWS2811
template <EOrder RGB_ORDER = RGB,
          uint8_t CHIP = WS2811_800kHz>
class CTeensy4Controller : public CPixelLEDController<RGB_ORDER, 8, 0xFF>
{
    OctoWS2811 *pocto;

public:
    CTeensy4Controller(OctoWS2811 *_pocto)
        : pocto(_pocto){};

    virtual void init() {}
    virtual void showPixels(PixelController<RGB_ORDER, 8, 0xFF> &pixels)
    {

        uint32_t i = 0;
        while (pixels.has(1))
        {
            uint8_t r = pixels.loadAndScale0();
            uint8_t g = pixels.loadAndScale1();
            uint8_t b = pixels.loadAndScale2();
            pocto->setPixel(i++, r, g, b);
            pixels.stepDithering();
            pixels.advanceData();
        }

        pocto->show();
    }
};
CTeensy4Controller<RGB, WS2811_800kHz> *pcontroller;

/*----------------------------Program---------------------------------*/
//main setup function
void setup() {

  Serial.begin(115200);
  delay(10);
  
  //teensy ethernet MAC setup
  teensyMAC(mac);
 
  static char teensyMac[23];
  sprintf(teensyMac, "MAC: %02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
  Serial.print("MAC: ");
  Serial.println(teensyMac);

  Ethernet.setSocketSize(ETHERNET_SOCKET_SIZE);
  Ethernet.begin(mac, ip);
  
  Serial.print("IP Address: ");
  Serial.println(Ethernet.localIP());
  Serial.print("MAC: ");
  Serial.println(teensyMac);

  Udp.begin(UDP_PORT);
  octo.begin();
  pcontroller = new CTeensy4Controller<RGB, WS2811_800kHz>(&octo);
  
  FastLED.addLeds(pcontroller, rgbarray, NUM_LEDS);
  FastLED.setBrightness(50);  
  
  //light test (no ethernet)  
  initTest();
  Serial.println("Test Sequence Complete");
}

//main loop for gathering ethernet packets
void loop() {
  //Process packets
   
  int packetSize = Udp.parsePacket(); //Read UDP packet count
    
  if(packetSize){
    Serial.println(packetSize);
    Udp.read(packetBuffer,ETHERNET_BUFFER); //read UDP packet
    pixelDisplay(packetBuffer, 1); //process data function
  }
  /*
  Serial.print(Counter++);
  Serial.print("Packet Size: ");
  Serial.println(packetSize);
  Serial.print("packetBuffer first 5 CRGB contents: ");
  for(int index = 0; index < 5; index++) {
    Serial.print(packetBuffer[index],HEX);
    Serial.print("\t");    // prints a tab
  }
  Serial.println();
  */;
  pixelrefresh(0);
}
/*-------------------------Functions LED----------------------------*/

// (1) timing function for LEDs (DO NOT TOUCH UNLESS NECESSARY)
static inline void pixelrefresh(const int syncrefresh){
  // Create static variables so that the code and variables can
  // all be declared inside a function 
  static unsigned long frametimestart;
  static unsigned long frametimeend;
  static unsigned long frametimechk;
  static unsigned long frameonce;
  unsigned long now = micros(); 

  //start frame time
  frametimestart = now;
  
  //Serial.println(frametimechk);
  //If we have framed no need to frame again update time to most recent
  if  (syncrefresh == 1){
    frametimeend = frametimestart; 
    frameonce = 1;
  }
   
  //If we havent framed this will increment via time and at some point will be true, 
  //if so we need to frame to clear out any buffer and the hold off untill 
  //we receive our next valid dmx packet. We use the pixel protocol to get a general rule of timing to compare to.

  frametimechk = frametimestart - frametimeend;
  // num leds time 30us + 300us reset to simulate the time it would take to write out pixels. 
  //this should help us not loop to fast and risk premature framing and jeopordize ethernet buffer
  if (frametimechk >= (NUM_LEDS * 30) + 300){
    frametimeend = frametimestart;
    if (frameonce == 1){
    octo.show();
    Serial.println ("Partial framing detected");
    frameonce = 0;  
    }
  }
}

// (2) flash red infinitely if Ethernet not working
void errorSequence() {
 while (1){
   LEDS.showColor(CRGB(125, 0, 0)); //turn all pixels on red
   delay(2000);
   LEDS.showColor(CRGB(0, 0, 0)); //turn off all pixels
   delay(2000);
 }
}

// (3) displays LED
//void pixelDisplay(uint8_t* pbuff, int count) {
//  for (int i = 0; i < NUM_LEDS; i++) {
//    byte charValueR = pbuff[i*3];
//    byte charValueG = pbuff[i*3+1];
//    byte charValueB = pbuff[i*3+2];
//    Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
//    octo.setPixel(i, charValueG,charValueR,charValueB/10); //RBG GRB
//  }
//  octo.show();
//  pixelrefresh(1);
//  fps2(10);
//}
void pixelDisplay(uint8_t* pbuff, int count) {
  int j = 0;
  for (int i = 0; i < 2000; i++) {
    if((i<two)){/* pin 2 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
      //octo.setPixel(j, 0, 125, 0);
    }
    else if(249<i && i<250+fourteen){/* pin 14 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    }
    else if(499<i && i<500+seven){/* pin 7 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    }
    else if(749<i && i<750+eight){/* pin 8 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    }
    else if(999<i && i<1000+six){/* pin 6 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    }
    else if(1249<i && i<1250+twenty){/* pin 20 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    }
    else if(1499<i && i<1500+twentyone){/* pin 21 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    }
    else if(1749<i && i<1750+five){/* pin 5 */
      byte charValueR = pbuff[j*3];
      byte charValueG = pbuff[j*3+1];
      byte charValueB = pbuff[j*3+2];
      //Serial.print(charValueR);Serial.print(",");Serial.print(charValueG);Serial.print(",");Serial.println(charValueB);
      octo.setPixel(i, charValueR,charValueG,charValueB); //RBG GRB
      j++;
    } 
    else {
      byte charValueR = 0;
      byte charValueG = 0;
      byte charValueB = 0;
      //Serial.println("I'm Here");
    }
  }
  octo.show();
  pixelrefresh(1);
  fps2(10);
}

// (4) timing function for LEDs (DO NOT TOUCH UNLESS NECESSARY)
static inline void fps2(const int seconds){
  // Create static variables so that the code and variables can
  // all be declared inside a function
  static unsigned long lastMillis;
  static unsigned long frameCount;
  static unsigned int framesPerSecond;
  
  // It is best if we declare millis() only once
  unsigned long now = millis();
  frameCount ++;
  if (now - lastMillis >= seconds * 1000) {
    framesPerSecond = frameCount / seconds;    
    Serial.print("FPS @ ");
    Serial.println(framesPerSecond);
    frameCount = 0;
    lastMillis = now;
  }

}

// (5) start up LED sequence. runs at board boot to make sure pixels are working
void initTest() {
  LEDS.clear(); //clear led assignments
  
  LEDS.showColor(CRGB(120, 250, 5)); //turn all pixels on green
  delay(3000);

  LEDS.showColor(CRGB(110, 250, 5)); //turn all pixels on red
  delay(3000);

  LEDS.showColor(CRGB(100, 250, 5)); //turn all pixels on blue
  delay(3000);

  LEDS.showColor(CRGB(0,0,0)); //turn off all pixels to start

  //process dot trace
//  for (int n=0; n < NUM_LEDS; n++){
//    octo.setPixel(n,100,0,0);
//    octo.setPixel(n-2,0,0,0);
//    octo.show();
//    delay(10);
//  }
//
//  for (int n=0; n < NUM_LEDS; n++){
//    octo.setPixel(n,0,100,0);
//    octo.setPixel(n-2,0,0,0);
//    octo.show();
//    delay(10);
//  }
//
//  for (int n=0; n < NUM_LEDS; n++){
//    octo.setPixel(n,0,0,100);
//    octo.setPixel(n-2,0,0,0);
//    octo.show();
//    delay(10); 
//  }
//  LEDS.showColor(CRGB(0,0,0)); //turn off all pixels to start
}

/*-----------------------Functions Ethernet----------------------*/
//define mac address of teensy board
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
    errorSequence();
  #endif
}
