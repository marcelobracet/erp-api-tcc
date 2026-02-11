package validation

import "unicode"

// OnlyDigits removes all non-digit characters.
func OnlyDigits(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return string(out)
}

// IsValidCNPJ validates a Brazilian CNPJ number.
// It accepts masked inputs (e.g. 12.345.678/0001-99).
func IsValidCNPJ(cnpj string) bool {
	cnpj = OnlyDigits(cnpj)
	if len(cnpj) != 14 {
		return false
	}

	// Reject sequences like 00000000000000
	allSame := true
	for i := 1; i < 14; i++ {
		if cnpj[i] != cnpj[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	d := make([]int, 14)
	for i := 0; i < 14; i++ {
		d[i] = int(cnpj[i] - '0')
	}

	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := 0; i < 12; i++ {
		sum += d[i] * w1[i]
	}
	mod := sum % 11
	dv1 := 0
	if mod >= 2 {
		dv1 = 11 - mod
	}
	if d[12] != dv1 {
		return false
	}

	sum = 0
	for i := 0; i < 13; i++ {
		val := d[i]
		if i == 12 {
			val = dv1
		}
		sum += val * w2[i]
	}
	mod = sum % 11
	dv2 := 0
	if mod >= 2 {
		dv2 = 11 - mod
	}

	return d[13] == dv2
}
