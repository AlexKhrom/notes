package converter

import (
	"strconv"
	"strings"
)

type HexConverter interface {
	Convert(hexString string) (decString string)
}

func NewConverter() HexConverter {
	return &Converter{}
}

type Converter struct {
	a []int
}

func (c *Converter) byteToInt(symb byte) int {
	if symb >= 'a' && symb <= 'z' {
		return int(symb-'a') + 10
	} else if symb >= '0' && symb <= '9' {
		return int(symb - '0')
	}
	return -1
}

func (c *Converter) nextNumber() int {
	tmp := 0
	for i := 0; i < len(c.a); i++ {
		tmp = tmp*16 + c.a[i]
		c.a[i] = tmp / 10
		tmp = tmp % 10
	}
	return tmp
}

func (c *Converter) isZero() bool {
	for i := 0; i < len(c.a); i++ {
		if c.a[i] != 0 {
			return false
		}
	}
	return true
}

func (c *Converter) Convert(hexString string) (decString string) {
	hexString = strings.ToLower(hexString)

	c.a = make([]int, len(hexString))
	for i := 0; i < len(hexString); i++ {
		c.a[i] = c.byteToInt(hexString[i])
	}

	for !c.isZero() {
		decString = strconv.Itoa(c.nextNumber()) + decString
	}
	return decString
}
