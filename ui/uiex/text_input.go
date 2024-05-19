package uiex

import (
	"image"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextInput struct {
	*block

	multiline bool
	readonly  bool

	caretAlphaSign float64
	caret          *Bar
	charIndex      int
	text           Text
}

func NewTextInput() *TextInput {
	tb := &TextInput{
		caret:          NewBar(),
		caretAlphaSign: -1,
		text:           NewBasicText(""),
	}
	tb.block = newBlock(tb)
	tb.SetContent(tb.text)
	tb.text.SetAlignX(ui.AlignOffset)
	tb.text.SetAlignY(ui.AlignOffset)
	tb.caret = NewBar()
	tb.caret.SetVisible(true)
	tb.caret.SetEventHandling(false)
	tb.caret.SetAlign(ui.AlignOffset, ui.AlignOffset)
	tb.caret.SetDirection(DirectionVertical)
	metrics := tb.text.Face().Metrics()
	tb.caret.SetLength(int(metrics.HAscent + metrics.HDescent))
	tb.AddDeco(tb.caret)
	tb.SetItemOptions(
		ui.WithCustomStyleFunc(func(itb *TextInput, is ui.InputState) {
			// Caret's alpha animation
			alpha := itb.caret.Alpha()
			alpha += itb.caretAlphaSign * 10. / 255.
			if alpha < 0 || alpha > 1 {
				alpha = min(max(alpha, 0), 1)
				itb.caretAlphaSign *= -1
			}
			itb.caret.SetColorAlpha(alpha)
		}),
		ui.WithCustomUpdateFunc(func(itb *TextInput, is ui.InputState) {
			txt := itb.text.Text()
			prevCharIndex := itb.charIndex
			paddX, _, paddY, _ := itb.text.Padding()
			textPadding := image.Pt(paddX, paddY)
			// If text changed need to reset the charindex
			itb.charIndex = min(itb.charIndex, len(txt))
			prevTextLength := len(txt)
			metrics := tb.text.Face().Metrics()
			tb.caret.SetLength(int(metrics.HAscent + metrics.HDescent))
			// Register focus and new caret position
			if itb.Hovered() {
				ebiten.SetCursorShape(ebiten.CursorShapeText)
			} else {
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			}
			// Accept inputs if focused
			if itb.Focused() {
				// Register new caret position on press
				if itb.JustPressed() {
					cs := itb.LastCursor()
					cs = cs.Sub(itb._block.LastRegion().Min)
					cs = cs.Sub(itb.caret.LastRegion().Size().Div(2))
					cs = cs.Sub(textPadding)
					itb.updateCaretIndex(cs)
				}
				// Text content
				runes := ebiten.AppendInputChars(nil)
				if len(runes) > 0 {
					txt = txt[:itb.charIndex] + string(runes) + txt[itb.charIndex:]
					for _, r := range runes {
						itb.charIndex += utf8.RuneLen(r)
					}
				}
				// Key presses
				tps := ebiten.TPS()
				_, charIncr := utf8.DecodeRuneInString(txt[itb.charIndex:])
				_, charDecr := utf8.DecodeLastRuneInString(txt[:itb.charIndex])
				charDecr = -charDecr
				wordSeps := " \n\t.,!;+-*/()[]{}"
				wordStops := "\n.,!;+-*/()[]{}"
				// Words skip in case ctrl is held
				if is.KeyPressDuration(ebiten.KeyControl) > 0 {
					charIncr = itb.charIndex
					if itb.charIndex < len(txt) {
						sepIndex := strings.IndexRune(wordSeps, rune(txt[charIncr]))
						for ; charIncr < len(txt); charIncr++ {
							if sepIndex != strings.IndexRune(wordSeps, rune(txt[charIncr])) {
								sepIndex = strings.IndexRune(wordSeps, rune(txt[charIncr]))
								break
							}
						}
						if sepIndex > -1 && !strings.ContainsRune(wordStops, rune(wordSeps[sepIndex])) {
							for ; charIncr < len(txt); charIncr++ {
								if sepIndex != strings.IndexRune(wordSeps, rune(txt[charIncr])) {
									break
								}
							}
						}
					}
					charIncr -= itb.charIndex
					charDecr = itb.charIndex - 1
					if itb.charIndex > 0 {
						sepIndex := strings.IndexRune(wordSeps, rune(txt[charDecr]))
						if sepIndex > -1 && !strings.ContainsRune(wordStops, rune(txt[charDecr])) {
							for ; charDecr > 0; charDecr-- {
								if sepIndex != strings.IndexRune(wordSeps, rune(txt[charDecr])) {
									sepIndex = strings.IndexRune(wordSeps, rune(txt[charDecr]))
									break
								}
							}
						}
						for ; charDecr > 0; charDecr-- {
							if sepIndex != strings.IndexRune(wordSeps, rune(txt[charDecr])) {
								charDecr++
								break
							}
						}
					}
					charDecr -= itb.charIndex
				}
				// Command keys
				switch {
				case is.KeyPressDuration(ebiten.KeyLeft) == 1, is.KeyPressDuration(ebiten.KeyLeft) >= tps/2:
					itb.charIndex = max(itb.charIndex+charDecr, 0)
				case is.KeyPressDuration(ebiten.KeyRight) == 1, is.KeyPressDuration(ebiten.KeyRight) >= tps/2:
					itb.charIndex = min(itb.charIndex+charIncr, len(txt))
				case is.KeyPressDuration(ebiten.KeyUp) == 1, is.KeyPressDuration(ebiten.KeyUp) >= tps/2:
					cs := itb.caret.SourceOffset()
					cs.Y -= int(itb.text.Layout().LineSpacing)
					cs = cs.Sub(textPadding)
					itb.updateCaretIndex(cs)
				case is.KeyPressDuration(ebiten.KeyDown) == 1, is.KeyPressDuration(ebiten.KeyDown) >= tps/2:
					cs := itb.caret.SourceOffset()
					cs.Y += int(itb.text.Layout().LineSpacing)
					cs = cs.Sub(textPadding)
					itb.updateCaretIndex(cs)
				case is.KeyPressDuration(ebiten.KeyBackspace) == 1, is.KeyPressDuration(ebiten.KeyBackspace) >= tps/2:
					if itb.charIndex > 0 {
						txt = txt[:itb.charIndex+charDecr] + txt[itb.charIndex:]
						itb.charIndex += charDecr
					}
				case is.KeyPressDuration(ebiten.KeyDelete) == 1, is.KeyPressDuration(ebiten.KeyDelete) >= tps/2:
					if itb.charIndex < len(txt) {
						txt = txt[:itb.charIndex] + txt[itb.charIndex+charIncr:]
					}
				case is.KeyPressDuration(ebiten.KeyTab) == 1, is.KeyPressDuration(ebiten.KeyTab) >= tps/2:
					txt = txt[:itb.charIndex] + "\t" + txt[itb.charIndex:]
					itb.charIndex++
				case is.KeyPressDuration(ebiten.KeyEnter) == 1, is.KeyPressDuration(ebiten.KeyEnter) >= tps/2,
					is.KeyPressDuration(ebiten.KeyNumpadEnter) == 1, is.KeyPressDuration(ebiten.KeyNumpadEnter) >= tps/2:
					txt = txt[:itb.charIndex] + "\n" + txt[itb.charIndex:]
					itb.charIndex++
				}
			} else {
				// Hide caret if not focused
				itb.caret.SetColorAlpha(0)
			}

			// Update caret source offset
			var y, bi int
			for i, r := range txt {
				if i == itb.charIndex {
					break
				}
				if r == '\n' {
					bi = i + 1
					y++
				}
			}

			adv := advance(txt[bi:itb.charIndex], itb.text.Face(), itb.text.TabSpaces())
			off := itb.SourceOffset()
			off = off.Add(image.Pt(
				int(adv),
				int(float64(y)*itb.text.Layout().LineSpacing),
			))
			// Scroll to content
			if prevCharIndex != itb.charIndex || prevTextLength != len(txt) {
				fullRegion := tb.block.LastFullRegion()
				caretRegion := image.Rect(
					off.X, off.Y, off.X+itb.caret.thickness, off.Y+itb.caret.length,
				).Add(fullRegion.Min)
				if !caretRegion.In(fullRegion) {
					diff := fullRegion.Max.Sub(caretRegion.Max)
					diff.X, diff.Y = min(diff.X, 0), min(diff.Y, 0)
					off = off.Add(diff)
					src := itb.SourceOffset()
					itb.SetSourceOffset(src.Add(diff))
				}
				if off.X < 0 || off.Y < 0 {
					diff := off
					diff.X, diff.Y = min(diff.X, 0), min(diff.Y, 0)
					off = off.Sub(diff)
					src := itb.SourceOffset()
					itb.SetSourceOffset(src.Sub(diff))
				}
			}
			// If caret has been sollicitated, refresh alpha for visibility
			if prevCharIndex != itb.charIndex || prevTextLength != len(txt) {
				itb.caret.SetColorAlpha(1)
			}
			itb.text.SetText(txt)
			// Add HDescent to center the caret vertically
			itb.caret.SetSourceOffset(off.
				// TODO: caret is not centered vertically anymore
				//int(metrics.HDescent))).
				Add(textPadding),
			)
		}),
	)
	// TODO: update caret length based on font size dynamically
	defaultTheme.apply(tb)

	return tb
}

func advance(str string, face text.Face, tabSpaces int) float64 {
	spaces := ""
	for i := 0; i < tabSpaces; i++ {
		spaces += " "
	}
	str = strings.ReplaceAll(str, "\t", spaces)
	return text.Advance(str, face)
}

func nextRuneLen(s string) int {
	_, count := utf8.DecodeRuneInString(s)
	return count
}

func (tb *TextInput) updateCaretIndex(cursor image.Point) {
	cursor = cursor.Sub(tb.SourceOffset())
	glyphs := tb.text.Glyphs()
	txt := tb.text.Text()
	y := 0
	lh := tb.text.Layout().LineSpacing
	lastRow := max(strings.Count(txt, "\n"), 0)
	targetRow := min(int((float64(cursor.Y)+lh/2)/lh), lastRow)
	var ir, ig int
	for ig <= len(glyphs) {
		for ir < len(txt) && txt[ir] == '\n' && y < targetRow {
			y++
			ir++
		}
		if y == targetRow {
			break
		}
		if txt[ir] == '\t' {
			ir++
			ig += tb.text.TabSpaces()
		} else {
			ir += (glyphs[ig].EndIndexInBytes - glyphs[ig].StartIndexInBytes) //nextRuneLen(txt[ir:])
			ig++
		}
	}
	if ir == len(txt) || txt[ir] == '\n' {
		tb.charIndex = ir
		return
	}

	// Look for the closest X spot in the row
	var p image.Point
	var closestPoint image.Point
	var closestDx = math.MaxInt
	var closestIndex = 0
	var closestWidth = 0.
	p.Y = int(float64(targetRow) * lh)
	spaces := ""
	// TODO: this is ugly
	for i := 0; i < tb.text.TabSpaces(); i++ {
		spaces += " "
	}

	for ig < len(glyphs) { //, g := range glyphs[ig:] {
		g := glyphs[ig]
		if ir >= len(txt) || txt[ir] == '\n' {
			break
		}

		s := txt[ir : ir+(g.EndIndexInBytes-g.StartIndexInBytes)]
		if txt[ir] == '\t' {
			s = spaces
		}
		width := advance(s, tb.text.Face(), tb.text.TabSpaces())

		p.X = int(g.X)

		dx := cursor.X - p.X
		if dx < 0 {
			dx = -dx
		}
		if dx <= closestDx {
			closestPoint = p
			closestDx = dx
			closestIndex = ir
			closestWidth = width
		}
		if txt[ir] == '\t' {
			ir++
			ig += tb.text.TabSpaces()
		} else {
			ir += len(s) //(g.EndIndexInBytes - g.StartIndexInBytes) //nextRuneLen(txt[ir:])
			ig++
		}
	}
	if closestDx != math.MaxInt && cursor.X > closestPoint.X+int(closestWidth/2) {
		end := strings.IndexRune(txt[closestIndex:], '\n')
		if end == -1 {
			end = len(txt) - 1
		} else {
			end += closestIndex
		}
		if end-closestIndex > 0 {
			closestIndex += nextRuneLen(txt[closestIndex:])
		}
	}
	tb.charIndex = closestIndex
}

func (tb *TextInput) InsertAtCaret(str string) {
	txt := tb.text.Text()
	txt = txt[:tb.charIndex] + str + txt[tb.charIndex:]
	tb.text.SetText(txt)
	tb.charIndex += len(str)
}

func (tb *TextInput) CaretPosition() int {
	return tb.charIndex
}

func (tb *TextInput) SetCaretPosition(pos int) {
	tb.charIndex = min(pos, max(len(tb.text.Text())-1, 0))
}

func (tb *TextInput) Text() Text {
	return tb.text
}

func (tb *TextInput) SetText(txt Text) {
	tb._block.SetContent(txt)
	tb.text = txt
	tb.text.SetAlignX(ui.AlignOffset)
	tb.text.SetAlignY(ui.AlignOffset)
}

func (tb *TextInput) LastFullRegion() image.Rectangle {
	return tb.text.LastFullRegion().Add(tb._block.LastFullRegion().Min)
}

func (tb *TextInput) LastRegion() image.Rectangle {
	return tb.text.LastRegion().Add(tb._block.LastRegion().Min)
}

func (tb *TextInput) SourceOffset() image.Point {
	return tb.text.SourceOffset()
}

func (tb *TextInput) SetSourceOffset(offset image.Point) {
	tb.text.SetSourceOffset(offset)
}

func (tb *TextInput) SetContent(txt Text) {
	tb.SetText(txt)
}

// Options

type TextInputOption func(*TextInput)

func (tb *TextInput) WithOptions(opts ...TextInputOption) *TextInput {
	for _, o := range opts {
		o(tb)
	}
	return tb
}

func (tb *TextInput) SetMultiline(multiline bool) {
	tb.multiline = multiline
}

func (tb *TextInput) SetCaretThickness(thickness int) {
	tb.caret.SetThickness(thickness)
}

func (tb *TextInput) SetCaretOptions(opts ...ui.ItemOption) {
	tb.caret.SetItemOptions(opts...)
}
