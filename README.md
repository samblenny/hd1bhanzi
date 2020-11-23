# hd1bhanzi
hanzi for high-dpi 1-bit monochrome

## Usage

### Generating List of Character Codepoints

To re-generate the `core2020_g.txt` and `iicore_g.txt` lists of Simplified Chinese
character codepoints taken from the kUnihanCore2020 and kIICore lists:

1. Download a copy of `Unihan.zip` from https://www.unicode.org/Public/13.0.0/ucd/

2. Expand the archive, creating a `Unihan` directory next to `gen_gsource_lists.rb`

3. Download https://blogs.adobe.com/CCJKType/files/2014/06/china-8105-06062014.txt
   and move the file next to `gen_gsource_lists.rb` (related article on TGH-2013 by
   Dr. Ken Lunde is at https://blogs.adobe.com/CCJKType/2014/03/china-8105.html)

3. Run `ruby gen_gsource_lists.rb`


### Generating PNG Glyph Grid Sprite Sheet

1. Make sure you have a working Go compiler, Imagemagick, and a bash shell

2. Download Noto Sans SC font family using the download button at
   https://fonts.google.com/specimen/Noto+Sans+SC?subset=chinese-simplified

3. Unzip the `Noto_Sans_SC.zip` font archive, creating a `Noto_Sans_SC` directory
   of font files next to `sprites.go`

4. Do `go run sprites.go` to create temporary SVG sprite sheet(s)

5. Do `./makeSpriteSheet.sh` to convert sprite sheet(s) from SVG to PNG


## Legal

## Unicode® Data Files

This project uses character lists derived from Unicode® Data Files at
https://www.unicode.org/Public/13.0.0/ucd/

Unicode and the Unicode Logo are registered trademarks of Unicode, Inc. in the
United States and other countries.

The Unicode copyright and permission notice:

```
Copyright © 1991-2020 Unicode, Inc. All rights reserved.
Distributed under the Terms of Use in https://www.unicode.org/copyright.html.

Permission is hereby granted, free of charge, to any person obtaining
a copy of the Unicode data files and any associated documentation
(the "Data Files") or Unicode software and any associated documentation
(the "Software") to deal in the Data Files or Software
without restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, and/or sell copies of
the Data Files or Software, and to permit persons to whom the Data Files
or Software are furnished to do so, provided that either
(a) this copyright and permission notice appear with all copies
of the Data Files or Software, or
(b) this copyright and permission notice appear in associated
Documentation.

THE DATA FILES AND SOFTWARE ARE PROVIDED "AS IS", WITHOUT WARRANTY OF
ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT OF THIRD PARTY RIGHTS.
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR HOLDERS INCLUDED IN THIS
NOTICE BE LIABLE FOR ANY CLAIM, OR ANY SPECIAL INDIRECT OR CONSEQUENTIAL
DAMAGES, OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE,
DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THE DATA FILES OR SOFTWARE.

Except as contained in this notice, the name of a copyright holder
shall not be used in advertising or otherwise to promote the sale,
use or other dealings in these Data Files or Software without prior
written authorization of the copyright holder.
```


## Noto Sans SC

This project uses glyphs from the Noto Sans SC typeface.

Noto Sans SC License Notice:
```
These fonts are licensed under the Open Font License.

You can use them freely in your products & projects - print or digital,
commercial or otherwise. However, you can't sell the fonts on their own.

This isn't legal advice, please consider consulting a lawyer and see the
full license for all details.
```
Open Font License Link: https://scripts.sil.org/OFL
