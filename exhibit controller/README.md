# Exhibit controller

## Description
All code that runs on the RPI which controls the entire exhibit.

## Data Structure/contraint for communication:
##### Total bits: 32 bits
* 3 bits: address
* 1 bit: direction

* 2 bits: priority
* 2 bits: speed

* 3 bits: type of effect 
* 1 bit padding

8 bits
Color A:
* 3 bits: lumosity (0 - 7)
* 5 bits: hue (0 - 31)

8 bits
Color B:
* 3 bits: lumosity (0 - 7)
* 5 bits: hue (0 - 31) 

* 4 bit padding at the end as a terminator
