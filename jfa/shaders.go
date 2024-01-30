package jfa

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"sort"
	"text/template"

	"github.com/hajimehoshi/ebiten/v2"
)

type hash = string

var (
	//go:embed shaders/jfa.tmpl.kage
	jfaTmplSrc string
	jfaTmpl    *template.Template

	jfaShaders = map[hash]*ebiten.Shader{}

	//go:embed shaders/encoding.kage
	encodingSrc    []byte
	encodingShader *ebiten.Shader
)

func jfaShaderHash(opts *GenerateOptions) hash {
	var masks []ColorMask
	if len(opts.PlainValueThresholds) > 0 {
		masks = make([]ColorMask, len(opts.PlainValueThresholds))
		for cm := range opts.PlainValueThresholds {
			masks = append(masks, cm)
		}
		sort.SliceStable(masks, func(i, j int) bool {
			return masks[i] < masks[j]
		})
	}

	return fmt.Sprintf("%d%v%v",
		opts.DistanceType,
		opts.EdgesPlain,
		masks,
	)
}

func JFAShader(opts *GenerateOptions) *ebiten.Shader {
	h := jfaShaderHash(opts)
	if s, ok := jfaShaders[h]; ok && s != nil {
		return s
	}

	tresholds := make([]string, 0, len(opts.PlainValueThresholds))
	for cm := range opts.PlainValueThresholds {
		var stmt string
		switch cm {
		case ColorMaskAlpha:
			stmt = "t *= (1-step(clr.a, ColorMaskAlpha))"
		case ColorMaskGreyscale:
			stmt = "t *= (1-step(clr.r*0.299+clr.g*0.587+clr.b*0.114, ColorMaskGreyscale))"
		case ColorMaskR:
			stmt = "t *= (1-step(clr.r, ColorMaskR))"
		case ColorMaskG:
			stmt = "t *= (1-step(clr.g, ColorMaskG))"
		case ColorMaskB:
			stmt = "t *= (1-step(clr.b, ColorMaskB))"
		}
		tresholds = append(tresholds, stmt)
	}
	thresholdCmp := ">"
	if opts.DistanceType == DistanceInterior {
		thresholdCmp = "<"
	}
	emptyValue := "s := step(uv.xy-0.5, vec2(0))"
	oobValue := "return minseed"
	if opts.EdgesPlain {
		emptyValue = "s := vec2(0)"
		oobValue = "return vec3(current, minseed.z)"
	}

	out := &bytes.Buffer{}
	err := jfaTmpl.Execute(out, map[string]any{
		"ThresholdChecks": tresholds,
		"ThresholdCmp":    thresholdCmp,
		"EmptyValue":      emptyValue,
		"OOBValue":        oobValue,
	})
	if err != nil {
		log.Fatal("can't execute jfa shader template: ", err)
	}

	s, err := ebiten.NewShader(out.Bytes())
	if err != nil {
		log.Fatal("can't compile jfa shader: ", err)
	}
	jfaShaders[h] = s

	return s
}

func encodingShaderHash(opts *GenerateOptions) hash {
	return fmt.Sprintf("%d",
		opts.Encoding,
	)
}

func EncodingShader() *ebiten.Shader {
	return encodingShader
}

func init() {
	var err error

	jfaTmpl = template.New("jfa")
	_, err = jfaTmpl.Parse(jfaTmplSrc)
	if err != nil {
		log.Fatalf("couldn't parse jfa shader template: %v\n", err)
	}

	encodingShader, err = ebiten.NewShader(encodingSrc)
	if err != nil {
		log.Fatalf("couldn't compile encoding shader: %v\n", err)
	}
}
