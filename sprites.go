package main

import (
	"bytes"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
	"path/filepath"
)

// INPUT: Font file containing glyphs to be used in the sprite sheet
const FontFilename = "NotoSansSC-Regular.otf"
const FontDirectory = "Noto_Sans_SC"

// INTPUT: List of hex-format character codepoints from kUnihanCore2020 source G
const GSourceFilename = "core2020_g.txt"

// OUTPUT: PNG file for writing the glyph grid sprite sheet
const OutputFilename = "hanzi.svg"

// Make a PNG glyph grid sprite sheet
func main() {
	fontPpem := 4
	fontSize := 32
	gridColumns := 20
	fnt := loadFont(filepath.Join(FontDirectory, FontFilename), fontPpem)
	charset := loadChars(GSourceFilename, gridColumns)[:800]
	svg := spriteSheet(fnt, charset, fontPpem, fontSize, gridColumns)
	writeSVG(OutputFilename, svg)
}

// Return a font after loading it from a file, and its ascent
func loadFont(filename string, ppem int) *sfnt.Font {
	fmt.Printf("loading font: %s\n", filename)
	fontData, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	f, err := sfnt.Parse(fontData)
	if err != nil {
		panic(err)
	}
	// Print font metadata
	b := sfnt.Buffer{}
	name, _ := f.Name(&b, sfnt.NameID(6))
	fmt.Printf("  Name(6): %s\n", name)
	mtr, _ := f.Metrics(&b, fixed.I(ppem), font.HintingFull)
	fmt.Printf("  Metrics(ppem:%d): {Height: %v, Ascent: %v, Descent: %v, XHeight: %v, CapHeight: %v}\n",
		ppem, mtr.Height, mtr.Ascent, mtr.Descent, mtr.XHeight, mtr.CapHeight)
	fmt.Printf("  NumGlyphs(): %v\n", f.NumGlyphs())
	return f
}

// Load and parse hex-format character codepoints into runes
func loadChars(filename string, columns int) []CharSpec {
	fmt.Printf("loading charset: %s... ", filename)
	cs := HanziMap(filename, columns)
	fmt.Printf("[len=%d]\n", len(cs))
	return cs
}

// Lay out a grid of character glyphs as a PNG sprite sheet
func spriteSheet(fnt *sfnt.Font, charset []CharSpec, ppem int, size int, columns int) string {
	glyphPaths := []sfnt.Segments{}
	bigBounds := image.Rectangle{}
	for _, c := range charset {
		s := c.GraphemeCluster()
		r := []rune(s)[0]
		// Find index for glyph
		index, err := fnt.GlyphIndex(nil, r)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  GlyphIndex(%s):%d\n", s, index)
		if index == 0 {
			panic("GlyphIndex = 0  (glyph not found)")
		}
		// Load glyph vector data
		glyphSegs, err := fnt.LoadGlyph(nil, index, fixed.I(ppem), nil)
		if err != nil {
			panic(err)
		}
		glyphPaths = append(glyphPaths, glyphSegs)
		// Include this glyph's bounds in the big bounding box
		fixedBounds, _, err := fnt.GlyphBounds(nil, index, fixed.I(ppem), font.HintingNone)
		if err != nil {
			panic(err)
		}
		bigBounds = bigBounds.Union(unfixRect(fixedBounds))
	}
	b := bigBounds
	fmt.Printf("bounds: x:%v y:%v w:%v h:%v\n", b.Min.X, b.Min.Y, b.Max.X-b.Min.X, b.Max.Y-b.Min.Y)
	// Render the glyphs into an SVG sprite sheet
	return renderSvgSpriteSheet(glyphPaths, bigBounds, size, columns)
}

// Convert point from 26_6 fixed-point to regular int (do not scale or truncate)
func unfixPt(p fixed.Point26_6) image.Point {
	return image.Pt(int(p.X), int(p.Y))
}

// Convert rectangle from 26_6 fixed-point to regular int (do not scale or truncate)
func unfixRect(r fixed.Rectangle26_6) image.Rectangle {
	return image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

// Write a SVG file
func writeSVG(filename string, svg string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	f.WriteString(svg)
	if err := f.Close(); err != nil {
		panic(err)
	}
}

// Holds mappings from extended grapheme clusters to sprite sheet glyph grid coordinates
type CharSpec struct {
	HexCluster string
	Row        int
	Col        int
}

// Parse and return the first codepoint of a hex grapheme cluster string.
// For example, "1f3c4-200d-2640-fe0f" -> 0x1F3C4
func (cs CharSpec) FirstCodepoint() uint32 {
	utf8 := StringFromHexGC(cs.HexCluster)
	codepoints := []rune(utf8)
	return uint32(codepoints[0])
}

// Convert a hex grapheme cluster string to a regular utf8 string.
// For example, "1f3c4-200d-2640-fe0f" -> "\U0001F3C4\u200d\u2640\ufe0f"
func (cs CharSpec) GraphemeCluster() string {
	return StringFromHexGC(cs.HexCluster)
}

// Parse a hex-codepoint format grapheme cluster into a utf-8 string
// For example, "1f3c4-200d-2640-fe0f" -> "\U0001F3C4\u200d\u2640\ufe0f"
func StringFromHexGC(hexGC string) string {
	base := 16
	bits := 32
	cluster := ""
	hexCodepoints := strings.Split(hexGC, "-")
	if len(hexCodepoints) < 1 {
		panic(fmt.Errorf("unexpected value for hexGC: %q", hexGC))
	}
	for _, hc := range hexCodepoints {
		n, err := strconv.ParseUint(hc, base, bits)
		if err != nil {
			panic(fmt.Errorf("unexpected value for hexGC: %q", hexGC))
		}
		cluster += string(rune(n))
	}
	return cluster
}

// Return mapping of hex-codepoint format grapheme clusters to grid coordinates
// in a glyph sprite sheet for the emoji font
func HanziMap(inputFile string, columns int) []CharSpec {
	text, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	// Start at top left corner of the sprite sheet glyph grid
	row := 0
	col := 0
	// Parse hex format grapheme cluster lines that should look like
	// "1f4aa-1f3fc\n" "1f4e1\n", etc. Comments starting with "#" are
	// possible. Order of grapheme cluster lines in the file should match a
	// row-major order traversal of the glyph grid.
	csList := []CharSpec{}
	for _, line := range strings.Split(string(text), "\n") {
		// Trim comments and leading/trailing whitespace
		txt := strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if len(txt) > 0 {
			// Add a CharSpec for this grapheme cluster
			csList = append(csList, CharSpec{txt, row, col})
			// Advance to next glyph position by row-major order
			col += 1
			if col == columns {
				row += 1
				col = 0
			}
		}
		// Skip blank lines and comments
	}
	return csList
}

// Return a string from rendering the given template and context data
func renderTemplate(templateString string, name string, context interface{}) string {
	fmap := template.FuncMap{"ToLower": strings.ToLower}
	t := template.Must(template.New(name).Funcs(fmap).Parse(templateString))
	var buf bytes.Buffer
	err := t.Execute(&buf, context)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Render an svg file with vector glpyhs arranged into a grid on a sprite sheet
func renderSvgSpriteSheet(glyphPaths []sfnt.Segments, bounds image.Rectangle, size int, columns int) string {
	// Calculate scaling factors between vector coordinates and final output size
	b := bounds
	unit := math.Max(float64(b.Max.X-b.Min.X), float64(b.Max.Y-b.Min.Y)) / float64(size)
	gridCell := unit * float64(size+2)
	// Render svg d= attribute values for <path d="..."/>
	glyphPathDVals := []string{}
	border := 2.0 * unit
	for i, segs := range glyphPaths {
		row := float64(i / columns)
		col := float64(i % columns)
		corner := image.Pt(int(border+col*gridCell), int(border+row*gridCell))
		glyphPathDVals = append(glyphPathDVals, renderPath(segs, bounds, corner))
	}
	n := len(glyphPaths)
	vbWidth := int(border + (gridCell * float64(columns)))
	rows := n / columns
	if n%columns > 0 {
		rows += 1
	}
	vbHeight := int(border + (gridCell * float64(rows)))
	context := struct {
		ViewboxMinX   int
		ViewboxMinY   int
		ViewboxWidth  int
		ViewboxHeight int
		Width         int
		Height        int
		GlyphPaths    []string
	}{0, 0, vbWidth, vbHeight, 2 + (size+2)*columns, 2 + (size+2)*rows, glyphPathDVals}
	// Generate rust source code and write it to a file
	svg := renderTemplate(svgTemplate, "svg", context)
	return svg
}

// Render sfnt glyph segments into the r-value of an svg path d=... attribute
func renderPath(glyphSegs sfnt.Segments, bounds image.Rectangle, corner image.Point) string {
	dOps := []string{}
	x := -bounds.Min.X + corner.X
	y := -bounds.Min.Y + corner.Y
	for _, s := range glyphSegs {
		a, b, c := unfixPt(s.Args[0]), unfixPt(s.Args[1]), unfixPt(s.Args[2])
		// Adjust coordinates to move glyph bounding box top left to (0,0) with +x=right, +y=down
		// (because font vectors have y=0 as baseline with ascent in -y and descent in +y)
		ax, ay, bx, by, cx, cy := a.X+x, a.Y+y, b.X+x, b.Y+y, c.X+x, c.Y+y
		switch s.Op {
		case sfnt.SegmentOpMoveTo:
			dOps = append(dOps, fmt.Sprintf("M%v %v", ax, ay))
		case sfnt.SegmentOpLineTo:
			dOps = append(dOps, fmt.Sprintf("L%v %d", ax, ay))
		case sfnt.SegmentOpQuadTo:
			dOps = append(dOps, fmt.Sprintf("Q%v %v %v %v", ax, ay, bx, by))
		case sfnt.SegmentOpCubeTo:
			dOps = append(dOps, fmt.Sprintf("C%v %v %v %v %v %v", ax, ay, bx, by, cx, cy))
		}
	}
	return strings.Join(dOps, " ")
}

// <svg viewBox="-500 -2000 3000 3000" width="600" height="600" xmlns="http://www.w3.org/2000/svg">
const svgTemplate = `<svg version="1.1" baseProfile="full"
viewBox="{{.ViewboxMinX}} {{.ViewboxMinY}} {{.ViewboxWidth}} {{.ViewboxHeight}}"
width="{{.Width}}" height="{{.Height}}"
xmlns="http://www.w3.org/2000/svg">
<style>path{stroke:none;fill:black;fill-opacity:1;/*shape-rendering:crispEdges;*/}</style>
{{ range $_, $d := .GlyphPaths -}}
<path d="{{$d}}"/>
{{ end -}}
</svg>
`
