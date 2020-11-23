#!/bin/bash

#GRAY="-colorspace gray +level 0%x85% -background white -alpha remove"
GRAY="-colorspace gray -level 0%,75% -background white -alpha remove"
DITHER2="-ordered-dither o2x2,2"
SRC="hanzi.svg"

# Scale and dither SVG sprite sheet to make 1-bit PNG sprite sheet
echo converting o2x2
convert $GRAY $DITHER2 $SRC hanzi_o2x2.png
