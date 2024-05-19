package app

import (
	"math/big"
	"slices"
	"strings"
)

type Calculator struct {
	numString string
	lastOpe   rune

	topBuffer string

	Operator rune
	Current  *big.Float
	Result   *big.Float
}

func New() *Calculator {
	return &Calculator{
		Current: big.NewFloat(0),
		Result:  big.NewFloat(0),
	}
}

func (c *Calculator) TopString() string {
	return strings.Replace(c.topBuffer, ".", ",", 1)
}

func (c *Calculator) InputString() string {
	s := c.Current.Text('g', 16)
	s = strings.Replace(s, ".", ",", 1)
	if strings.Contains(s, "e") {
		return s
	}
	i, f, ok := strings.Cut(s, ",")
	if len(i) <= 3 {
		return s
	}
	s = ""
	for j := 0; len(i)-j > 0; j++ {
		s += string(i[len(i)-j-1])
		if (j+1)%3 == 0 {
			s += " "
		}
	}
	runes := []rune(s[:])
	slices.Reverse(runes)
	s = strings.TrimSpace(string(runes))
	if ok {
		s += "," + f
	}

	return s
}

func (c *Calculator) ProcessToken(t rune) {
	switch t {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		if c.lastOpe != 0 {
			c.numString = ""
		}
		if c.lastOpe == '=' {
			c.Result.SetFloat64(0)
		}
		if t == '.' && strings.Contains(c.numString, ".") {
			return
		}
		c.numString += string(t)
		c.Current.Parse(c.numString, 10)
		c.lastOpe = 0
	case '<':
		c.numString = c.numString[:max(len(c.numString)-1, 0)]
		c.Current.Parse(c.numString, 10)
		c.lastOpe = 0
	case '±':
		c.Current.Mul(c.Current, big.NewFloat(-1))
		c.lastOpe = 0
	case '²', '³':
		if c.lastOpe != 0 && c.Operator != '=' {
			return
		}

		if c.Operator == 0 || c.Operator == '=' {
			c.topBuffer = c.Current.Text('g', 16) + string(t)
		} else {
			c.topBuffer = c.Result.Text('g', 16) + " " + string(c.Operator) + " " + c.Current.Text('g', 16) + string(t)
		}
		switch t {
		case '²':
			c.Current.Mul(c.Current, c.Current)
		case '³':
			c.Current.Mul(c.Current, c.Current)
			c.Current.Mul(c.Current, c.Current)
		}
		//c.Result.Set(c.Current)
		c.lastOpe = 0
	case '+', '–', '÷', '×', '=':
		if c.lastOpe == 0 {
			c.numString = ""
			switch c.Operator {
			case '+':
				c.Result.Add(c.Result, c.Current)
			case '–':
				c.Result.Sub(c.Result, c.Current)
			case '÷':
				c.Result.Quo(c.Result, c.Current)
			case '×':
				c.Result.Mul(c.Result, c.Current)
			case '=':
				c.Result.Set(c.Current)
			case 0:
				c.Result.Set(c.Current)
			}
		}

		c.Operator = t
		c.Current.Set(c.Result)
		c.topBuffer = c.Result.Text('g', 16) + " " + string(t)
		c.lastOpe = t
	default:
	}
}
