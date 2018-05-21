package gojay

import (
	"math"
)

var digits []int8

const maxInt64toMultiply = math.MaxInt64 / 10
const maxInt32toMultiply = math.MaxInt32 / 10
const maxInt16toMultiply = math.MaxInt16 / 10
const maxInt8toMultiply = math.MaxInt8 / 10
const maxUint8toMultiply = math.MaxUint8 / 10
const maxUint16toMultiply = math.MaxUint16 / 10
const maxUint32toMultiply = math.MaxUint32 / 10
const maxUint64toMultiply = math.MaxUint64 / 10
const maxUint32Length = 10
const maxUint64Length = 20
const maxUint16Length = 5
const maxUint8Length = 3
const maxInt32Length = 10
const maxInt64Length = 19
const maxInt16Length = 5
const maxInt8Length = 3
const invalidNumber = int8(-1)

var pow10uint64 = [20]uint64{
	0,
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

func init() {
	digits = make([]int8, 256)
	for i := 0; i < len(digits); i++ {
		digits[i] = invalidNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		digits[i] = i - int8('0')
	}
}

func (dec *Decoder) skipNumber() (int, error) {
	end := dec.cursor + 1
	// look for following numbers
	for j := dec.cursor + 1; j < dec.length || dec.read(); j++ {
		switch dec.data[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = j + 1
			continue
		case '.':
			end = j + 1
			continue
		case ',', '}', ']':
			return end, nil
		case ' ', '\n', '\t', '\r':
			continue
		}
		// invalid json we expect numbers, dot (single one), comma, or spaces
		return end, dec.raiseInvalidJSONErr(dec.cursor)
	}
	return end, nil
}

func (dec *Decoder) getExponent() int64 {
	start := dec.cursor
	end := dec.cursor
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch dec.data[dec.cursor] { // is positive
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = dec.cursor + 1
		case '-':
			dec.cursor++
			return -dec.getExponent()
		case '+':
			dec.cursor++
			return dec.getExponent()
		default:
			// if nothing return 0
			// could raise error
			if start == end {
				dec.raiseInvalidJSONErr(dec.cursor)
				return 0
			}
			return dec.atoi64(start, end-1)
		}
	}
	if start == end {
		dec.raiseInvalidJSONErr(dec.cursor)
		return 0
	}
	return dec.atoi64(start, end-1)
}
