package core

import (
	"image/color"
	"slices"

	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

var (
	builtins = []string{
		"sin", "cos", "tan", "asin", "acos", "atan", "atan2",
		"pow", "exp", "log", "exp2", "log2", "sqrt", "inversesqrt",
		"abs", "sign", "floor", "ceil", "fract", "mod", "min", "max",
		"clamp", "mix", "step", "smoothstep", "length", "distance",
		"dot", "cross", "normalize", "faceforward", "reflect", "refract",
		"transpose", "dfdx", "dfdy", "fwidth",
		"imageSrc0At", "imageSrc1At", "imageSrc2At", "imageSrc3At",
		"imageSrc0Origin", "imageSrc1Origin", "imageSrc2Origin", "imageSrc3Origin",
		"imageSrc0UnsafeAt", "imageSrc1UnsafeAt", "imageSrc2UnsafeAt", "imageSrc3UnsafeAt",
		"imageSrc0Size", "imageSrc1Size", "imageSrc2Size", "imageSrc3Size",
		"imageDstOrigin", "imageDstSize", "discard",
	}
	builtinsTypes = []string{
		"float", "int", "ivec2", "ivec3", "ivec4", "bool",
		"vec2", "vec3", "vec4", "mat2", "mat3", "mat4",
	}
	builtinOperators = []string{
		":=", "==", "=", "+", "-", "/", "*", "%", "^", "&", "<", ">", "|", "[", "]",
		"+=", "-=", "/=", "*=", ".", "<=", ">=", "&&", "||",
	}
)

type Lexer struct {
	lexer chroma.Lexer
}

func newLexer() *Lexer {
	return &Lexer{
		lexer: lexers.Get("go"),
	}
}

var (
	clrComment     = color.RGBA{64, 156, 64, 255}
	clrKeyword     = color.RGBA{192, 0, 156, 255}
	clrType        = color.RGBA{156, 128, 255, 255}
	clrLit         = color.RGBA{96, 128, 255, 255}
	clrBuiltin     = color.RGBA{96, 96, 255, 255}
	clrName        = color.RGBA{156, 156, 192, 255}
	clrOperator    = color.RGBA{96, 192, 255, 255}
	clrPunctuation = color.RGBA{192, 128, 0, 255}
)

func (p *Lexer) Format(rt *uiex.RichText) error {
	txt := rt.Text()
	iterator, err := p.lexer.Tokenise(nil, txt)
	if err != nil {
		return err
	}
	rt.SetText("")
	rt.Reset()
	for _, t := range iterator.Tokens() {
		var noEffect bool
		switch t.Type.Category() {
		case chroma.Comment.Category():
			rt.PushColorFg(clrComment)
		case chroma.Keyword.Category():
			switch t.Type {
			case chroma.KeywordType:
				rt.PushColorFg(clrType)
			default:
				rt.PushColorFg(clrKeyword)
			}
		case chroma.LiteralNumber.Category():
			rt.PushColorFg(clrLit)
		case chroma.Name.Category():
			switch t.Type {
			case chroma.NameBuiltin, chroma.NameFunction, chroma.NameOther:
				switch {
				case slices.Contains(builtinsTypes, t.Value):
					rt.PushColorFg(clrType)
				case slices.Contains(builtins, t.Value):
					rt.PushColorFg(clrBuiltin)
				default:
					rt.PushColorFg(clrName)
				}
			}
		case chroma.Punctuation.Category(), chroma.Operator.Category():
			switch {
			case slices.Contains(builtinOperators, t.Value):
				rt.PushColorFg(clrOperator)
			default:
				rt.PushColorFg(clrPunctuation)
			}
		default:
			noEffect = true
		}
		rt.Append(t.Value)
		if !noEffect {
			rt.Pop()
		}
	}

	return nil
}
