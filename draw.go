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
	"unicode"
	"unicode/utf8"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

const printRoman = false

var texts = [...]string{
	"" +
		"Twinkle, twinkle, little star\n" +
		"How I wonder what you are\n" +
		"  Up above the world so high\n" +
		"  Like a diamond in the sky\n" +
		"When the blazing sun is gone\n" +
		"When he nothing shines upon\n" +
		"  Then you show your little light\n" +
		"  Twinkle, twinkle, all the night\n" +
		"Then the traveler in the dark\n" +
		"Thanks you for your tiny spark\n" +
		"  He could not see which way to go\n" +
		"  If you did not twinkle so\n" +
		"In the dark blue sky you keep\n" +
		"And often through my curtains peep\n" +
		"  For you never shut your eye\n" +
		"  Till the sun is in the sky\n" +
		"As your bright and tiny spark\n" +
		"Lights the traveler in the dark\n" +
		"  Though I know not what you are\n" +
		"  Twinkle, twinkle, little star\n" +
		"Twinkle, twinkle, little star\n" +
		"How I wonder what you are\n" +
		"  Up above the world so high\n" +
		"  Like a diamond in the sky\n" +
		"",

	"" +
		"  Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.\n" +
		"  Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battlefield of that war. We have come to dedicate a portion of that field, as a final resting place for those who here gave their lives%a that that nation might live%i. It is altogether fitting and proper that we should do this.\n" +
		"  But, in a larger sense, we can not dedicate - we can not consecrate - we can not hallow - this ground. The brave men, living and dead, who struggled here, have consecrated it, far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us - that from these honored dead we take increased devotion to that cause for which they gave the last full measure of devotion - that we here highly resolve that these dead shall not have died in vain - that this nation, under God, shall have a new birth of freedom - and that government of the people, by the people, for the people, shall not perish from the earth.\n" +
		"",

	"" +
		"And take me disappearing\n" +
		"Through the smoke rings of my mind\n" +
		"Down the foggy ruins of time\n" +
		"Far past the frozen leaves\n" +
		"The haunted frightened trees\n" +
		"Out to the windy%i beach\n" +
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
		"Time flies like an arrow. Fruit flies like a banana.\n" +
		"\n" +
		"Computers are useless. They can only give you answers.\n" +
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
		"\n" +
		"It is not from the benevolence of the butcher, the brewer, or the baker that we expect our dinner, but from their regard to their own interest.\n" +
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
		"(Internet myth) According to a researcher at Cambridge University, it doesn't matter in what order the letters in a word are, the only important thing is that the first and last letter be at the right place. The rest can be a total mess and you can still read%i it without problem. This is because the human mind does not read%i every letter by itself, but the word as a whole.\n" +
		"\n" +
		"Who would know aught of art must learn, act, and then take his ease.\n" +
		"\n" +
		"Are those shy Eurasian footwear, cowboy chaps, or jolly earthmoving headgear?\n" +
		"",

	"" +
		"Now is the winter of our discontent\n" +
		"Made glorious summer by this sun of York;\n" +
		"And all the clouds that lour'd upon our house\n" +
		"In the deep bosom of the ocean buried.\n" +
		"Now are our brows bound with victorious wreaths;\n" +
		"Our bruised arms hung up for monuments;\n" +
		"Our stern alarums changed to merry meetings,\n" +
		"Our dreadful marches to delightful measures.\n" +
		"Grim visaged war hath smooth'd his wrinkled front;\n" +
		"And now, instead of mounting barbed steeds\n" +
		"To fright the souls of fearful adversaries,\n" +
		"He capers nimbly in a lady's chamber\n" +
		"To the lascivious pleasing of a lute.\n" +
		"But I, that am not shaped for sportive tricks,\n" +
		"Nor made to court an amorous looking glass;\n" +
		"I, that am rudely stamp'd, and want love's majesty\n" +
		"To strut before a wanton ambling nymph;\n" +
		"I, that am curtail'd of this fair proportion,\n" +
		"Cheated of feature by dissembling nature,\n" +
		"Deformed, unfinish'd, sent before my time\n" +
		"Into this breathing world, scarce half made up,\n" +
		"And that so lamely and unfashionable\n" +
		"That dogs bark at me as I halt by them;\n" +
		"Why, I, in this weak piping time of peace,\n" +
		"…\n" +
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
		"\n" +
		"I am so clever that sometimes I don't understand a single word of what I am saying.\n" +
		"",

	"" +
		"Once upon a midnight dreary,\n  while I pondered, weak and weary,\n" +
		"Over many a quaint and curious volume of forgotten lore —\n" +
		"While I nodded, nearly napping,\n  suddenly there came a tapping,\n" +
		"As of some one gently rapping, rapping at my chamber door.\n" +
		"\"'Tis some visitor,\" I muttered, \"tapping at my chamber door —\n" +
		"                            Only this and nothing more.\"\n" +
		"Ah, distinctly I remember\n  it was in the bleak December;\n" +
		"And each separate dying ember wrought its ghost upon the floor.\n" +
		"Eagerly I wished the morrow; —\n  vainly I had sought to borrow\n" +
		"From my books surcease of sorrow — sorrow for the lost Lenore —\n" +
		"For the rare and radiant maiden whom the angels name Lenore —\n" +
		"                            Nameless here for evermore.\n" +
		"And the silken, sad, uncertain\n  rustling of each purple curtain\n" +
		"Thrilled me — filled me with fantastic terrors never felt before;\n" +
		"So that now, to still the beating\n  of my heart, I stood repeating\n" +
		"\"'Tis some visitor entreating entrance at my chamber door —\n" +
		"Some late visitor entreating entrance at my chamber door; —\n" +
		"                            This it is and nothing more.\"\n" +
		"…\n" +
		"",
}

func loadDict() {
	f, err := os.Open("third-party/Britfone/britfone.main.3.0.1.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Bytes()

		if i := bytes.IndexByte(line, ','); i >= 0 {
			k, v := string(line[:i]), strings.TrimSpace(string(line[i+1:]))
			if (k == "") || (v == "") {
				log.Fatalf("bad Britfone line: %q\n", line)
			}
			if _, ok := dict[k]; ok {
				log.Fatalf("duplicate Britfone key: %q\n", k)
			}
			dict[k] = v

		} else if _, ok := dict[string(line)]; ok {
			log.Fatalf("duplicate Britfone key: %q\n", line)
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

	for r, cluster := range letters {
		if cluster == "" {
			continue
		}
		s, diacritic := cluster[:len(cluster)-1], cluster[len(cluster)-1]
		if (s == "") || ((diacritic != '\'') && (diacritic != '~')) {
			s, diacritic = cluster, 0
		}

		width := 16 * utf8.RuneCountInString(s)
		m := image.NewGray(image.Rect(0, 0, width, 28))
		c.SetClip(m.Bounds())
		c.SetDst(m)
		c.DrawString(s, freetype.Pt(0, 26))

		if false {
			// Change s/false/true/ in the line above to draw plain vowels,
			// without diacritics.

		} else if diacritic == '\'' {
			const y = 0
			const x = 0
			m.SetGray(x+7, y+5, color.Gray{0xFF})
			m.SetGray(x+8, y+5, color.Gray{0xFF})
			m.SetGray(x+6, y+6, color.Gray{0xFF})
			m.SetGray(x+7, y+6, color.Gray{0xFF})
			m.SetGray(x+8, y+6, color.Gray{0xFF})
			m.SetGray(x+9, y+6, color.Gray{0xFF})
			m.SetGray(x+6, y+7, color.Gray{0xFF})
			m.SetGray(x+7, y+7, color.Gray{0xFF})
			m.SetGray(x+8, y+7, color.Gray{0xFF})
			m.SetGray(x+9, y+7, color.Gray{0xFF})
			m.SetGray(x+6, y+8, color.Gray{0xFF})
			m.SetGray(x+7, y+8, color.Gray{0xFF})
			m.SetGray(x+8, y+8, color.Gray{0xFF})
			m.SetGray(x+9, y+8, color.Gray{0xFF})
			m.SetGray(x+7, y+9, color.Gray{0xFF})
			m.SetGray(x+8, y+9, color.Gray{0xFF})

		} else if diacritic == '~' {
			const y = 7
			x0 := 3
			x1 := width - 2
			for ; x0 < x1; x0++ {
				m.SetGray(x0, y+0, color.Gray{0xFF})
				m.SetGray(x0, y+1, color.Gray{0xFF})
			}
			for x0 = x1 - 2; x0 < x1; x0++ {
				m.SetGray(x0, y+2, color.Gray{0xFF})
			}
		}

		glyphs[int64(r)] = &image.Alpha{
			Pix:    m.Pix,
			Stride: m.Stride,
			Rect:   m.Rect,
		}
	}
}

func parse(s string) (word string, remaining string) {
	for i := 0; i < len(s); i++ {
		if r, n := utf8.DecodeRuneInString(s[i:]); r <= ' ' {
			return strings.ToUpper(s[:i]), s[i:]
		} else if i != 0 {
			// No-op.
		} else if !isAlpha(r) && (r != '%') {
			return strings.ToUpper(s[:i+n]), s[i+n:]
		}
	}
	return strings.ToUpper(s), ""
}

func isConsonant(r rune) bool {
	switch r {
	case 'K', 'S', 'T', 'N', 'H', 'L', 'B', 'V', 'F', 'X',
		'G', 'Z', 'D', 'M', 'J', 'R', 'P', 'W', 'C', 'Y':
		return true
	}
	return false
}

func isVowel(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U':
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
		g := glyphs[int64(r)]
		if dst != nil {
			draw.DrawMask(dst, dst.Bounds().Add(image.Point{x, y}),
				fg, image.Point{}, g, image.Point{}, draw.Over)
		}
		x += g.Bounds().Dx() * 15 / 16
		return x
	}

	dictKey, suffix := "", ""
	for i := len(englishWord) - 1; i >= 0; i-- {
		if r := rune(englishWord[i]); ('A' <= r) && (r <= 'Z') {
			dictKey, suffix = englishWord[:i+1], englishWord[i+1:]
			break
		}
	}

	spelling := dict[dictKey]
	if spelling == "" {
		if dst != nil {
			log.Printf("%q not in Britfone", dictKey)
		}
		incompleteDict = true
		return x
	}

	runes := []rune(nil)
	if suffix == "" {
		runes = []rune(spelling)
	} else {
		runes = []rune(spelling + " " + suffix)
	}

	underDot, seenUnderDot := false, false
	for i := 0; i < len(runes); {
		r := runes[i]
		if (r == ' ') || (r == 'ˌ') || (r == 'ː') {
			i++
			continue
		} else if r == 'ˈ' {
			underDot, seenUnderDot = true, true
			i++
			continue
		}

		glyphsKey := int64(r)
		i++
		if i >= len(runes) {
			// No-op.
		} else if r = runes[i]; unicode.IsLetter(r) && (r != 'ː') {
			glyphsKey = (glyphsKey << 32) | int64(r)
			i++
		}

		g := glyphs[glyphsKey]
		if g == nil {
			log.Fatalf("couldn't draw %q (%q)", englishWord, spelling)
		}
		if dst != nil {
			if printRoman {
				print(roman[glyphsKey])
			}

			draw.DrawMask(dst, dst.Bounds().Add(image.Point{x, y}),
				fg, image.Point{}, g, image.Point{}, draw.Over)

			if underDot {
				underDot = false
				drawLowDot(dst, x, y, fg)
			}
		}

		x += (g.Bounds().Dx() * 15 / 16)
	}

	if !seenUnderDot {
		println("no underdot:", englishWord)
	}

	if (dst != nil) && printRoman {
		print(" ")
	}
	return x
}

func do(outName string, text string) {
	const pageInset = 25

	fg := &image.Uniform{C: color.RGBA{0x00, 0x00, 0x7F, 0xFF}}
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
		if true {
			draw.Draw(
				rgba, image.Rect((imageWidth/2), 0, (imageWidth/2)+1, imageHeight),
				guide, image.Point{}, draw.Src)
		}
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
				if printRoman {
					println()
				}

				line := originalText[:len(originalText)-len(s)]
				originalText = originalText[len(line):]
				drawEnglish(rgba, (imageWidth/2)+x, y-50, red, line)
				continue
			}

			word, remaining := parse(s)
			x1 := drawWord(nil, x, y, fg, word)
			if (x > pageInset) && (x1 > ((imageWidth / 2) - pageInset)) {
				x, y = pageInset, y+50
				if printRoman {
					println()
				}

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
		do(fmt.Sprintf("miileeniol-example-%d.png", i), text)
	}
}

var glyphs = map[int64]*image.Alpha{}

var letters = map[int64]string{
	'\'': "'",
	'"':  "\"",
	'+':  "+",
	'-':  "-",
	'?':  "?",
	'!':  "!",
	',':  ",",
	'.':  ".",
	';':  ";",
	':':  ":",
	'(':  "(",
	')':  ")",
	'…':  "…",
	'—':  "—",

	(int64('a') << 32) | int64('ɪ'): "aı~",
	(int64('a') << 32) | int64('ʊ'): "au~",
	(int64('e') << 32) | int64('ɪ'): "eı~",
	(int64('i')):                    "ı'",
	(int64('u')):                    "u'",
	(int64('æ')):                    "a'",
	(int64('ɐ')):                    "ε'",
	(int64('ɑ')):                    "a~",
	(int64('ɒ')):                    "o'",
	(int64('ɔ')):                    "o~",
	(int64('ɔ') << 32) | int64('ɪ'): "oı~",
	(int64('ə')):                    "ε~",
	(int64('ə') << 32) | int64('ʊ'): "εu~",
	(int64('ɛ')):                    "e~",
	(int64('ɛ') << 32) | int64('ə'): "eε~",
	(int64('ɜ')):                    "e'",
	(int64('ɪ')):                    "ı~",
	(int64('ɪ') << 32) | int64('ə'): "ıε~",
	(int64('ʊ')):                    "u~",
	(int64('ʊ') << 32) | int64('ə'): "uε~",

	(int64('b')):                    "B",
	(int64('d')):                    "D",
	(int64('d') << 32) | int64('ʒ'): "J",
	(int64('f')):                    "F",
	(int64('g')):                    "G",
	(int64('h')):                    "H",
	(int64('j')):                    "Y",
	(int64('k')):                    "K",
	(int64('l')):                    "L",
	(int64('m')):                    "M",
	(int64('n')):                    "N",
	(int64('p')):                    "P",
	(int64('s')):                    "S",
	(int64('t')):                    "T",
	(int64('t') << 32) | int64('ʃ'): "Ч", // tx
	(int64('v')):                    "V",
	(int64('w')):                    "W",
	(int64('z')):                    "Z",
	(int64('ð')):                    "Δ", // dh
	(int64('ŋ')):                    "Γ", // ng
	(int64('ɹ')):                    "R",
	(int64('ʃ')):                    "X",
	(int64('ʒ')):                    "Ж", // zh
	(int64('θ')):                    "Θ", // th
}

var roman = map[int64]string{
	(int64('a') << 32) | int64('ɪ'): "ai",
	(int64('a') << 32) | int64('ʊ'): "au",
	(int64('e') << 32) | int64('ɪ'): "ei",
	(int64('i')):                    "ia",
	(int64('u')):                    "ue",
	(int64('æ')):                    "ae",
	(int64('ɐ')):                    "ua",
	(int64('ɑ')):                    "aa",
	(int64('ɒ')):                    "oe",
	(int64('ɔ')):                    "oa",
	(int64('ɔ') << 32) | int64('ɪ'): "oi",
	(int64('ə')):                    "oo",
	(int64('ə') << 32) | int64('ʊ'): "eu",
	(int64('ɛ')):                    "ee",
	(int64('ɛ') << 32) | int64('ə'): "eo",
	(int64('ɜ')):                    "ea",
	(int64('ɪ')):                    "ii",
	(int64('ɪ') << 32) | int64('ə'): "ie",
	(int64('ʊ')):                    "uu",
	(int64('ʊ') << 32) | int64('ə'): "ue",

	(int64('b')):                    "b",
	(int64('d')):                    "d",
	(int64('d') << 32) | int64('ʒ'): "j",
	(int64('f')):                    "f",
	(int64('g')):                    "g",
	(int64('h')):                    "h",
	(int64('j')):                    "y",
	(int64('k')):                    "k",
	(int64('l')):                    "l",
	(int64('m')):                    "m",
	(int64('n')):                    "n",
	(int64('p')):                    "p",
	(int64('s')):                    "s",
	(int64('t')):                    "t",
	(int64('t') << 32) | int64('ʃ'): "tx",
	(int64('v')):                    "v",
	(int64('w')):                    "w",
	(int64('z')):                    "z",
	(int64('ð')):                    "dh",
	(int64('ŋ')):                    "ng",
	(int64('ɹ')):                    "r",
	(int64('ʃ')):                    "x",
	(int64('ʒ')):                    "zh",
	(int64('θ')):                    "th",
}

var dict = map[string]string{
	"A":             "ˈə",
	"ADVERSARIES":   "ˈæ d v ə s ə ɹ i z",
	"AESTHETIC":     "ɛ s θ ˈɛ t ɪ k",
	"AESTHETICALLY": "ɛ s θ ˈɛ t ɪ k ə l i",
	"AFTERTHOUGHT":  "ˈɑː f t ə θ ɔː t",
	"AGUE":          "ˈeɪ g j uː",
	"ALARUMS":       "ə l ˈɑː ɹ ɐ m z",
	"AMBLING":       "ˈæ m b l ɪ ŋ",
	"AMMUNITION":    "æ m j u n ˈɪ ʃ ə n",
	"AMOROUS":       "ˈæ m ə ɹ ə s",
	"AN":            "ˈæ n",
	"ANARCHY":       "ˈæ n ɑː k i",
	"AND":           "ˈæ n d",
	"ARE":           "ˈɑː",
	"AS":            "ˈæ z",
	"ASSORTED":      "ə s ˈɔː t ɪ d",
	"ASSORTMENT":    "ə s ˈɔː t m ə n t",
	"AT":            "ˈæ t",
	"AUGHT":         "ˈɔː t",
	"BALDRIC":       "b ˈɔː l d ɹ ɪ k",
	"BARBED":        "b ˈɑː b d",
	"BE":            "b ˈiː",
	"BECAUSE":       "b ɪ k ˈɒ z",
	"BENEVOLENCE":   "b ɛ n ˈɛ v ə l ə n s",
	"BLADED":        "b l ˈeɪ d ə d",
	"BLAZING":       "b l ˈeɪ z ɪ ŋ",
	"BLOWPIPE":      "b l ˈəʊ p aɪ p",
	"BRAILLE":       "b ɹ ˈeɪ l",
	"BREWER":        "b ɹ ˈuː ə",
	"BROWS":         "b ɹ ˈaʊ z",
	"CALTROPS":      "k ˈæ l t ɹ ə p s",
	"CAN":           "k ˈæ n",
	"CAPERS":        "k ˈeɪ p ə z",
	"CENTIMETER":    "s ˈɛ n t ɪ m iː t ə",
	"CHAINMAIL":     "tʃ ˈeɪ n m eɪ l",
	"CHEESEWIRE":    "tʃ ˈiː z w aɪ ə",
	"CIRCLED":       "s ˈɜː k ə l d",
	"CONCEIVED":     "k ə n s ˈiː v d",
	"CONSECRATE":    "k ə n s ˈɛ k ɹ ˈeɪ t",
	"CONSECRATED":   "k ə n s ˈɛ k ɹ ˈeɪ t ɪ d",
	"CONSOLE":       "k ə n s ˈəʊ l",
	"CORKED":        "k ˈɔː k d",
	"CORPS":         "k ˈɔː",
	"CROSSBOW":      "k ɹ ˈɒ s b əʊ",
	"CUNNING":       "k ˈɐ n ɪ ŋ",
	"CURTAIL'D":     "k ɜː t ˈeɪ l d",
	"DAGGER":        "d ˈæ g ə",
	"DARKENED":      "d ˈɑː k ə n d",
	"DEFORMED":      "d ɪ f ˈɔː m d",
	"DETRACT":       "d ɪ t ɹ ˈæ k t",
	"DEVOTION":      "d ɪ v ˈəʊ ʃ ə n",
	"DID":           "d ˈɪ d",
	"DIMENSION":     "d aɪ m ˈɛ n ʃ ə n",
	"DIMENSIONS":    "d aɪ m ˈɛ n ʃ ə n z",
	"DIMMED":        "d ˈɪ m d",
	"DISAPPEARING":  "d ɪ s ə p ˈɪə ɹ ɪ ŋ",
	"DISCONTENT":    "d ɪ s k ə n t ˈɛ n t",
	"DISSEMBLING":   "d ɪ s ˈɛ m b l ɪ ŋ",
	"DOES":          "d ˈɐ z",
	"DREARY":        "d ɹ ˈɪə ɹ i",
	"EAGERLY":       "ˈiː g ə l i",
	"EARTHMOVING":   "ˈɜː θ m uː v ɪ ŋ",
	"EMBER":         "ˈɛ m b ə",
	"ENGAGED":       "ɪ n g ˈeɪ dʒ d",
	"ENTREATING":    "ɛ n t ɹ ˈiː t ɪ ŋ",
	"EURASIAN":      "j ʊə ɹ ˈeɪ ʒ ə n",
	"EVERMORE":      "ˈɛ v ə m ɔː",
	"EVERY":         "ˈɛ v ɹ i",
	"FALCON":        "f ˈæ l k ə n",
	"FALCONER":      "f ˈæ l k ə n ə",
	"FOR":           "f ˈɔː",
	"GOVERNMENT":    "g ˈɐ v ə n m ə n t",
	"GRAPNEL":       "g ɹ ˈæ p n ə l",
	"GYRE":          "dʒ ˈaɪ ə",
	"HAD":           "h ˈæ d",
	"HAPPENETH":     "h ˈæ p ə n ɛ θ",
	"HAS":           "h ˈæ z",
	"HATH":          "h ˈæ θ",
	"HAUNTED":       "h ˈɔː n t ɪ d",
	"HAVE":          "h ˈæ v",
	"HE'LL":         "h ˈiː l",
	"HE'S":          "h ˈiː z",
	"HEADGEAR":      "h ˈɛ d g ɪə",
	"HIS":           "h ˈɪ z",
	"I'LL":          "ˈaɪ l",
	"IF":            "ˈɪ f",
	"IN":            "ˈɪ n",
	"INSIDE":        "ɪ n s ˈaɪ d",
	"INTEREST":      "ˈɪ n t ɹ ɪ s t",
	"INTERESTING":   "ˈɪ n t ɹ ɪ s t ɪ ŋ",
	"INTO":          "ɪ n t ˈuː",
	"INTRICATE":     "ˈɪ n t ɹ ɪ k ə t",
	"IS":            "ˈɪ z",
	"IT":            "ˈɪ t",
	"IT'S":          "ˈɪ t s",
	"ITS":           "ˈɪ t s",
	"JAUNTY":        "dʒ ˈɔː n t i",
	"JUST":          "dʒ ˈɐ s t",
	"KLATCHIAN":     "k l ˈæ tʃ ɪə n",
	"KNUCKLES":      "n ˈɐ k ə l z",
	"LADY'S":        "l ˈeɪ d i z",
	"LAMELY":        "l ˈeɪ m l i",
	"LASCIVIOUS":    "l ə s ˈɪ v i ə s",
	"LEAD%E":        "l ˈɛ d",
	"LENORE":        "l ɛ n ˈɔː",
	"LIVE%I":        "l ˈɪ v",
	"LIVES%A":       "l ˈaɪ v z",
	"LOCKPICKS":     "l ˈɒ k p ɪ k s",
	"LOOSED":        "l ˈuː s t",
	"LORE":          "l ˈɔː",
	"LOUR'D":        "l ˈɔː d",
	"LOVE'S":        "l ˈɐ v z",
	"LUTE":          "l ˈuː t",
	"MARCHES":       "m ˈɑː tʃ ɪ z",
	"MONUMENTS":     "m ˈɒ n j ʊ m ə n t s",
	"MORROW":        "m ˈɒ ɹ əʊ",
	"NAMELESS":      "n ˈeɪ m l ɪ s",
	"NAPPING":       "n ˈæ p ɪ ŋ",
	"NEITHER":       "n ˈaɪ ð ə",
	"NIMBLY":        "n ˈɪ m b l i",
	"NOBLY":         "n ˈəʊ b l i",
	"NYMPH":         "n ˈɪ m f",
	"OBLITERATION":  "ə b l ɪ t ə ɹ ˈeɪ ʃ ə n",
	"OF":            "ˈɒ v",
	"OFTEN":         "ˈɒ f ə n",
	"OR":            "ˈɔː",
	"PERMIT":        "p ə m ˈɪ t",
	"PIPING":        "p ˈaɪ p ɪ ŋ",
	"POLLINATE":     "p ˈɒ l ə n eɪ t",
	"PONDERED":      "p ˈɒ n d ə d",
	"POOR":          "p ˈɔː",
	"PRAYER":        "p ɹ ˈɛə",
	"QUAINT":        "k w ˈeɪ n t",
	"RADIANT":       "ɹ ˈeɪ d ɪə n t",
	"RAPIER":        "ɹ ˈeɪ p ɪə",
	"RAPPING":       "ɹ ˈæ p ɪ ŋ",
	"READ%I":        "ɹ ˈiː d",
	"REPEATING":     "ɹ ɪ p ˈiː t ɪ ŋ",
	"RESEARCHER":    "ɹ ˈiː s ɜː t ə",
	"RUDELY":        "ɹ ˈuː d l i",
	"RUINS":         "ɹ ˈuː ɪ n z",
	"RUSTLING":      "ɹ ˈɐ s l ɪ ŋ",
	"SANDS":         "s ˈæ n d z",
	"SEER":          "s ˈɪə",
	"SEPARATE":      "s ˈɛ p ɹ ɪ t",
	"SHALL":         "ʃ ˈæ l",
	"SHEATHS":       "ʃ ˈiː θ s",
	"SHINES":        "ʃ ˈaɪ n z",
	"SILHOUETTED":   "s ɪ l ʊ w ˈɛ t ɪ d",
	"SILKEN":        "s ˈɪ l k ə n",
	"SLINGSHOT":     "s l ˈɪ ŋ ʃ ɒ t",
	"SMOOTH'D":      "s m ˈuː ð d",
	"SOMETIMES":     "s ˈɐ m t ˌaɪ m z",
	"SPORTIVE":      "s p ˈɔː t ɪ v",
	"STAMP'D":       "s t ˈæ m p t",
	"STEEDS":        "s t ˈiː d z",
	"STEMS":         "s t ˈɛ m z",
	"STUDDING-SAIL": "s t ˈɐ n s ə l",
	"SUBTRACTS":     "s ə b t ɹ ˈæ k t s",
	"SURCEASE":      "s ˈɜː s iː s",
	"SUSY":          "s ˈuː z i",
	"SWARD":         "s w ˈɔː d",
	"TAPPING":       "t ˈæ p ɪ ŋ",
	"TEAR%E":        "t ˈɛə",
	"TEAR%I":        "t ˈɪə",
	"TEPPIC":        "t ˈɛ p ɪ k",
	"TERRORS":       "t ˈɛ ɹ ə z",
	"THAN":          "ð ˈæ n",
	"THAT":          "ð ˈæ t",
	"THE":           "ð ˈə",
	"THEM":          "ð ˈɛ m",
	"THERE'S":       "ð ˈɛə z",
	"THIS":          "ð ˈɪ s",
	"THRILLED":      "θ ɹ ˈɪ l d",
	"TIS":           "t ˈɪ z",
	"TLINGAS":       "t l ˈɪ ŋ g ə z",
	"TO":            "t ˈuː",
	"TODAY":         "t ə d ˈeɪ",
	"UNFASHIONABLE": "ɐ n f ˈæ ʃ ə n ə b ə l",
	"UNFINISH'D":    "ɐ n f ˈɪ n ɪ ʃ t",
	"UNFINISHED":    "ɐ n f ˈɪ n ɪ ʃ t",
	"VAINLY":        "v ˈeɪ n l i",
	"VIA":           "v ˈiː ə",
	"VICTORIOUS":    "v ɪ k t ˈɔː ɹ ɪə s",
	"VISAGED":       "v ɪ z ˈɑː dʒ d",
	"WANTON":        "w ˈɒ n t ə n",
	"WAS":           "w ˈɒ z",
	"WHO'S":         "h ˈuː z",
	"WIDENING":      "w ˈaɪ d ə n ɪ ŋ",
	"WILL":          "w ˈɪ l",
	"WINCED":        "w ˈɪ n s t",
	"WINDY%I":       "w ˈɪ n d i",
	"WITH":          "w ˈɪ ð",
	"WITHOUT":       "w ɪ ð ˈaʊ t",
	"WOUND%A":       "w ˈaʊ n d",
	"WREATHS":       "ɹ ˈiː θ s",
	"WRINKLED":      "ɹ ˈɪ ŋ k ə l d",
	"YEARS":         "j ˈɪə z",
	"YOU":           "j ˈuː",
	"YOU'LL":        "j ˈuː l",
	"YOUR":          "j ˈɔː",
}
