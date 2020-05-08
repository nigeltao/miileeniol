// Copyright 2020 Nigel Tao.
//
// Licensed under the MIT license.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

const version0Dot1 = true

var texts = [...]string{
	"" +
		"  Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.\n" +
		"  Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battlefield of that war. We have come to dedicate a portion of that field, as a final resting place for those who here gave their lives that that nation might live%i. It is altogether fitting and proper that we should do this.\n" +
		"  But, in a larger sense, we can not dedicate - we can not consecrate - we can not hallow - this ground. The brave men, living and dead, who struggled here, have consecrated it, far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us - that from these honored dead we take increased devotion to that cause for which they gave the last full measure of devotion - that we here highly resolve that these dead shall not have died in vain - that this nation, under God, shall have a new birth of freedom - and that government of the people, by the people, for the people, shall not perish from the earth.\n" +
		"",

	"" +
		"  I have a friend who's an artist and has sometimes taken a view which I don't " +
		"agree with very well. He'll hold up a flower and say \"look how beautiful it " +
		"is,\" and I'll agree. Then he says \"I as an artist can see how beautiful this is " +
		"but you as a scientist take this all apart and it becomes a dull thing,\" and I " +
		"think that he's kind of nutty. First of all, the beauty that he sees is " +
		"available to other people and to me too, I believe. Although I may not be quite " +
		"as refined aesthetically as he is, I can appreciate the beauty of a flower. " +
		"At the same time, I see much more about the flower than he sees. I could " +
		"imagine the cells in there, the complicated actions inside, which also have a " +
		"beauty. I mean it's not just beauty at this dimension, at one centimeter; " +
		"there's also beauty at smaller dimensions, the inner structure, also the " +
		"processes. The fact that the colors in the flower evolved in order to attract " +
		"insects to pollinate it is interesting; it means that insects can see the " +
		"color. It adds a question: does this aesthetic sense also exist in the lower " +
		"forms? Why is it aesthetic? All kinds of interesting questions which the " +
		"science knowledge only adds to the excitement, the mystery and the awe of a " +
		"flower. It only adds. I don't understand how it subtracts.\n" +
		"",

	"" +
		"And take me disappearing\n" +
		"Through the smoke rings of my mind\n" +
		"Down the foggy ruins of time\n" +
		"Far past the frozen leaves\n" +
		"The haunted frightened trees\n" +
		"Out to the windy beach\n" +
		"Far from the twisted reach of crazy sorrow.\n" +
		"Yes, to dance beneath the diamond sky\n" +
		"With one hand waving free\n" +
		"Silhouetted by the sea\n" +
		"Circled by the circus sands\n" +
		"With all memory and fate\n" +
		"Driven deep beneath the waves\n" +
		"Let me forget about today until tomorrow.\n" +
		"\n" +
		"I must not fear. Fear is the mind killer. Fear is the little death that brings total obliteration. I will face my fear. I will permit it to pass over me and through me. And when it has gone past I will turn the inner eye to see its path. Where the fear has gone there will be nothing. Only I will remain.\n" +
		"\n" +
		"I am so clever that sometimes I don't understand a single word of what I am saying.\n" +
		"",

	"" +
		"I returned, and saw under the sun, that the race is not to the swift, nor the battle to the strong, neither yet bread to the wise, nor yet riches to men of understanding, nor yet favor to men of skill; but time and chance happeneth to them all.\n" +
		"\n" +
		"Turning and turning in the widening gyre\n" +
		"The falcon cannot hear the falconer;\n" +
		"Things fall apart; the center cannot hold;\n" +
		"Mere anarchy is loosed upon the world,\n" +
		"The blood dimmed tide is loosed, and everywhere\n" +
		"The ceremony of innocence is drowned;\n" +
		"The best lack all conviction, while the worst\n" +
		"Are full of passionate intensity.\n" +
		"\n" +
		"According to a researcher at Cambridge University, it doesn't matter in what order the letters in a word are, the only important thing is that the first and last letter be at the right place. The rest can be a total mess and you can still read it without problem. This is because the human mind does not read every letter by itself, but the word as a whole.\n" +
		"\n" +
		"It is not from the benevolence of the butcher, the brewer, or the baker that we expect our dinner, but from their regard to their own interest.\n" +
		"",

	"" +
		"Dearest creature in creation\n" +
		"Studying English pronunciation,\n" +
		"  I will teach you in my verse\n" +
		"  Sounds like corpse, corps, horse and worse.\n" +
		"I will keep you, Susy, busy,\n" +
		"Make your head with heat grow dizzy;\n" +
		"  Tear%i in eye, your dress you'll tear%e;\n" +
		"  Queer, fair seer, hear my prayer.\n" +
		"Pray, console your loving poet,\n" +
		"Make my coat look new, dear, sew it!\n" +
		"  Just compare heart, hear and heard,\n" +
		"  Dies and diet, lord and word.\n" +
		"Sword and sward, retain and Britain\n" +
		"(Mind the latter how it's written).\n" +
		"  Made has not the sound of bade,\n" +
		"  Say - said, pay - paid, laid but plaid.\n" +
		"Now I surely will not plague you\n" +
		"With such words as vague and ague,\n" +
		"  But be careful how you speak,\n" +
		"  Say: gush, bush, steak, streak, break, bleak,\n" +
		"Previous, precious, fuchsia, via\n" +
		"Recipe, pipe, studding-sail, choir;\n" +
		"  Woven, oven, how and low,\n" +
		"  Script, receipt, shoe, poem, toe.\n" +
		"…\n" +
		"",

	"" +
		"  He sighed and opened the black box and took out his rings and slipped them on. Another box held a set of knives and Klatchian steel, their blades darkened with lamp black. Various cunning and intricate devices were taken from velvet bags and dropped into pockets. A couple of long bladed throwing tlingas were slipped into their sheaths inside his boots. A thin silk line and folding grapnel were wound%a around his waist, over the chainmail shirt. A blowpipe was attached to its leather thong and dropped down the back of his cloak; Teppic picked a slim tin container with an assortment of darts, their tips corked and their stems braille coded for ease of selection in the dark.\n" +
		"  He winced, checked the blade of his rapier and slung the baldric over his right shoulder, to balance the bag of lead%e slingshot ammunition. As an afterthought he opened his sock drawer and took a pistol crossbow, a flask of oil, a roll of lockpicks and, after some consideration, a punch dagger, a bag of assorted caltrops and a set of brass knuckles.\n" +
		"  Teppic picked up his hat and checked its lining for the coil of cheesewire. He placed it on his head at a jaunty angle, took a last satisfied look at himself in the mirror, turned on his heel and, very slowly, fell over.\n" +
		"",
}

var dict = map[string]string{}

type glyph struct {
	mask    *image.Alpha
	overBar bool
}

var glyphs = map[rune]glyph{}

const lInset = 3
const rInset = 2

func loadDict() {
	f, err := os.Open("dict.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Bytes()
		if i := bytes.IndexByte(line, '\t'); i >= 0 {
			k, v := string(line[:i]), string(line[i+1:])
			if (k == "") || (v == "") {
				log.Fatalf("bad dict.txt line: %q\n", line)
			}
			if _, ok := dict[k]; ok {
				log.Fatalf("duplicate dict.txt key: %q\n", k)
			}
			dict[k] = v

		} else if _, ok := dict[string(line)]; ok {
			// awk '{print $1}' dict.txt | sort | uniq -d | less
			// log.Fatalf("duplicate dict.txt key: %q\n", line)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}

func makeGlyphs() {
	f, err := freetype.ParseFont(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(26.65)
	c.SetSrc(image.White)
	c.SetHinting(font.HintingFull)

	for r, s := range letters {
		if s == "" {
			continue
		}

		width := 16 * utf8.RuneCountInString(s[1:])
		m := image.NewGray(image.Rect(0, 0, width, 28))
		c.SetClip(m.Bounds())
		c.SetDst(m)
		c.DrawString(s[1:], freetype.Pt(0, 26))

		glyphs[r] = glyph{
			mask: &image.Alpha{
				Pix:    m.Pix,
				Stride: m.Stride,
				Rect:   m.Rect,
			},
			overBar: s[0] == '~',
		}
	}
}

func parse(s string) (word string, remaining string) {
	for i := 0; i < len(s); i++ {
		if r, n := utf8.DecodeRuneInString(s[i:]); r <= ' ' {
			return strings.ToLower(s[:i]), s[i:]
		} else if i != 0 {
			// No-op.
		} else if !isAlpha(r) && (r != '%') {
			return strings.ToLower(s[:i+n]), s[i+n:]
		}
	}
	return strings.ToLower(s), ""
}

func isConsonant(r rune) bool {
	switch r {
	case 'k', 's', 't', 'n', 'h', 'l', 'b', 'v', 'f', 'x',
		'g', 'z', 'd', 'm', 'j', 'r', 'p', 'w', 'c', 'y':
		return true
	}
	return false
}

func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func isAlpha(r rune) bool {
	switch {
	case ('A' <= r) && (r <= 'Z'):
		return true
	case ('a' <= r) && (r <= 'z'):
		return true
	}
	return false
}

func drawEnglish(dst *image.RGBA, x int, y int, fg image.Image, line string) {
	if dst == nil {
		return
	}
	for ; (line != "") && (line[len(line)-1] == '\n'); line = line[:len(line)-1] {
	}

	for {
		i := strings.IndexByte(line, '%')
		if (i < 0) || (i+1 > len(line)) {
			break
		}
		line = line[:i] + line[i+2:]
	}

	goreg.SetClip(dst.Bounds())
	goreg.SetDst(dst)
	goreg.SetSrc(fg)
	goreg.DrawString(line, freetype.Pt(x, y+26))
}

const bigNeg = -(1 << 31)

func drawBar(dst *image.RGBA, x0 int, x1 int, y int, fg image.Image) int {
	if (dst == nil) || (x0 == bigNeg) {
		return bigNeg
	}
	c := fg.At(0, 0)
	for ; x0 < x1; x0++ {
		dst.Set(x0, y+0, c)
		dst.Set(x0, y+1, c)
	}
	for x0 = x1 - 2; x0 < x1; x0++ {
		dst.Set(x0, y+2, c)
	}
	return bigNeg
}

func drawHighDot(dst *image.RGBA, x int, y int, fg image.Image) {
	if dst == nil {
		return
	}
	c := fg.At(0, 0)

	dst.Set(x+7, y+5, c)
	dst.Set(x+8, y+5, c)
	dst.Set(x+6, y+6, c)
	dst.Set(x+7, y+6, c)
	dst.Set(x+8, y+6, c)
	dst.Set(x+9, y+6, c)
	dst.Set(x+6, y+7, c)
	dst.Set(x+7, y+7, c)
	dst.Set(x+8, y+7, c)
	dst.Set(x+9, y+7, c)
	dst.Set(x+6, y+8, c)
	dst.Set(x+7, y+8, c)
	dst.Set(x+8, y+8, c)
	dst.Set(x+9, y+8, c)
	dst.Set(x+7, y+9, c)
	dst.Set(x+8, y+9, c)
}

func drawLowDot(dst *image.RGBA, x int, y int, fg image.Image) {
	if dst == nil {
		return
	}
	c := fg.At(0, 0)

	dst.Set(x+7, y+28, c)
	dst.Set(x+8, y+28, c)
	dst.Set(x+6, y+29, c)
	dst.Set(x+7, y+29, c)
	dst.Set(x+8, y+29, c)
	dst.Set(x+9, y+29, c)
	dst.Set(x+6, y+30, c)
	dst.Set(x+7, y+30, c)
	dst.Set(x+8, y+30, c)
	dst.Set(x+9, y+30, c)
	dst.Set(x+6, y+31, c)
	dst.Set(x+7, y+31, c)
	dst.Set(x+8, y+31, c)
	dst.Set(x+9, y+31, c)
	dst.Set(x+7, y+32, c)
	dst.Set(x+8, y+32, c)
}

var incompleteDict = false

func drawWord(dst *image.RGBA, x int, y int, fg image.Image, englishWord string) (newX int) {
	if englishWord == "" {
		return x
	} else if r, _ := utf8.DecodeRuneInString(englishWord); !isAlpha(r) {
		g := glyphs[r]
		down := 0
		if dst != nil {
			draw.DrawMask(dst, dst.Bounds().Add(image.Point{x, y + down}),
				fg, image.Point{}, g.mask, image.Point{}, draw.Over)
		}
		x += g.mask.Bounds().Dx() * 15 / 16
		return x
	}

	key, suffix := "", ""
	for i := len(englishWord) - 1; i >= 0; i-- {
		if r := rune(englishWord[i]); isConsonant(r) || isVowel(r) {
			key, suffix = englishWord[:i+1], englishWord[i+1:]
			break
		}
	}

	spelling := dict[key]
	if spelling == "" {
		if dst != nil {
			println(key)
		}
		incompleteDict = true
		return x
	}

	dot := true
	vowel := rune(0)
	down := 0
	overBar := bigNeg
	overBarDown := 0
	for _, r := range spelling + suffix {
		if isVowel(r) {
			if vowel == 0 {
				vowel = r << 24
				continue
			}
			r = vowel | r
		} else if !isConsonant(r) {
			down = 0
			overBar = drawBar(dst, overBar, x-rInset, y+overBarDown, fg)
			if r == '-' {
				dot = true
				continue
			}
		}

		g := glyphs[r]
		if g.mask == nil {
			log.Fatalf("couldn't draw %q", spelling)
		}

		if dst != nil {
			if !g.overBar {
				overBar = drawBar(dst, overBar, x-rInset, y+overBarDown, fg)
				if ((r >> 24) != 0) && !version0Dot1 {
					drawHighDot(dst, x, y, fg)
				}

			} else if overBar == bigNeg {
				overBar = x + lInset
				overBarDown = down + 2
				if vowel != 0 {
					overBarDown += 5
				}
			}

			draw.DrawMask(dst, dst.Bounds().Add(image.Point{x, y + down}),
				fg, image.Point{}, g.mask, image.Point{}, draw.Over)

			if dot && !version0Dot1 {
				dot = false
				if false {
					drawLowDot(dst, x, y, fg)
				}
			}
		}
		x += g.mask.Bounds().Dx() * 15 / 16

		if vowel != 0 {
			vowel = 0
			if version0Dot1 {
				down = 5
			}
		}
	}

	if overBar != bigNeg {
		overBar = drawBar(dst, overBar, x-rInset, y+overBarDown, fg)
	}
	return x
}

func do(outName string, text string) {
	const pageInset = 20

	fg := &image.Uniform{C: color.RGBA{0x00, 0x00, 0x7F, 0xFF}}
	if version0Dot1 {
		fg.C = color.RGBA{0x00, 0x7F, 0x00, 0xFF}
	}
	red := &image.Uniform{C: color.RGBA{0x7F, 0x00, 0x00, 0xFF}}

	const imageWidth = 256 * 7
	const imageHeight = 256 * 5
	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.ZP, draw.Src)

	// Draw guidelines.
	{
		guide := &image.Uniform{C: color.RGBA{0xDD, 0xDD, 0xDD, 0xFF}}
		for y := 25 + pageInset; y < imageHeight; y += 50 {
			draw.Draw(
				rgba, image.Rect(0, y, imageWidth, y+1),
				guide, image.Point{}, draw.Src)
		}
		draw.Draw(
			rgba, image.Rect((imageWidth/2), 0, (imageWidth/2)+1, imageHeight),
			guide, image.Point{}, draw.Src)
	}

	// Render glyphs.
	{
		originalText := text

		x, y, s := pageInset, pageInset, text
		for s != "" {
			if c := s[0]; c == ' ' {
				x += 15
				s = s[1:]
				continue
			} else if c == '\n' {
				x = pageInset
				y += 50
				s = s[1:]

				line := originalText[:len(originalText)-len(s)]
				originalText = originalText[len(line):]
				drawEnglish(rgba, (imageWidth/2)+x, y-50, red, line)
				continue
			}

			word, remaining := parse(s)
			x1 := drawWord(nil, x, y, fg, word)
			if (x > pageInset) && (x1 > ((imageWidth / 2) - pageInset)) {
				x, y = pageInset, y+50

				line := originalText[:len(originalText)-len(s)]
				originalText = originalText[len(line):]
				drawEnglish(rgba, (imageWidth/2)+x, y-50, red, line)
			}
			x = drawWord(rgba, x, y, fg, word)
			s = remaining
		}
	}
	if incompleteDict {
		log.Fatal("incomplete dict.txt")
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create(outName)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Flush()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", outName)
}

var goreg *freetype.Context

func main() {
	f, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	goreg = freetype.NewContext()
	goreg.SetDPI(72)
	goreg.SetFont(f)
	goreg.SetFontSize(26.5)
	goreg.SetHinting(font.HintingFull)

	loadDict()
	makeGlyphs()
	for i, text := range texts {
		do(fmt.Sprintf("out%d.png", i), text)
	}
}

var letters = map[rune]string{
	'\'': " '",
	'"':  " \"",
	'+':  " +",
	'-':  " -",
	'?':  " ?",
	'!':  " !",
	',':  " ,",
	'.':  " .",
	';':  " ;",
	':':  " :",
	'(':  " (",
	')':  " )",
	'…':  " …",

	('a' << 24) | 'a': "~a",
	('a' << 24) | 'e': " a",
	('a' << 24) | 'i': "~aı",
	('a' << 24) | 'u': "~au",

	('e' << 24) | 'a': " e",
	('e' << 24) | 'e': "~e",
	('e' << 24) | 'i': "~eı",
	('e' << 24) | 'o': "~eε",

	('i' << 24) | 'a': " ε",
	('i' << 24) | 'i': "~ı",
	('i' << 24) | 'o': "~ıε",

	('o' << 24) | 'a': " o",
	('o' << 24) | 'e': "~o",
	('o' << 24) | 'i': "~oı",
	('o' << 24) | 'o': "~ε",
	('o' << 24) | 'u': "~ou",

	('u' << 24) | 'a': " ı",
	('u' << 24) | 'e': " u",
	('u' << 24) | 'u': "~u",
}

func init() {
	if version0Dot1 {
		letters['p'] = " P"
		letters['b'] = "~B"
		letters['t'] = " T"
		letters['d'] = "~T"
		letters['k'] = " K"
		letters['g'] = "~K"
		letters['m'] = " M"
		letters['n'] = "~M"
		letters['l'] = " L"
		letters['r'] = "~L"
		letters['f'] = " F"
		letters['v'] = "~F"
		letters['c'] = " H"
		letters['h'] = "~H"
		letters['s'] = " S"
		letters['z'] = "~S"
		letters['x'] = " J"
		letters['j'] = "~J"
		letters['w'] = " Y"
		letters['y'] = "~Y"

	} else {
		letters['p'] = " P"
		letters['b'] = " B"
		letters['t'] = " T"
		letters['d'] = " D"
		letters['k'] = " K"
		letters['g'] = " G"
		letters['m'] = " M"
		letters['n'] = " N"
		letters['l'] = " L"
		letters['r'] = " R"
		letters['f'] = " F"
		letters['v'] = " V"
		letters['c'] = " Θ"
		letters['h'] = " H"
		letters['s'] = " S"
		letters['z'] = " Z"
		letters['x'] = " Ш"
		letters['j'] = " J"
		letters['w'] = " W"
		letters['y'] = " Y"
	}

	// Γ
}
