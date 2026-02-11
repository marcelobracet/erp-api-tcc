package validation

import "testing"

func TestOnlyDigits(t *testing.T) {
	got := OnlyDigits("12.345.678/0001-99")
	if got != "12345678000199" {
		t.Fatalf("OnlyDigits() = %q; want %q", got, "12345678000199")
	}
}

func TestIsValidCNPJ_InvalidLength(t *testing.T) {
	if IsValidCNPJ("123") {
		t.Fatal("expected invalid for short input")
	}
	if IsValidCNPJ("123456789012345") {
		t.Fatal("expected invalid for long input")
	}
}

func TestIsValidCNPJ_RejectsAllSameDigits(t *testing.T) {
	if IsValidCNPJ("00.000.000/0000-00") {
		t.Fatal("expected invalid for repeated digits")
	}
	if IsValidCNPJ("11111111111111") {
		t.Fatal("expected invalid for repeated digits")
	}
}

func TestIsValidCNPJ_ValidGenerated(t *testing.T) {
	cnpj := generateValidCNPJ("123456780001")
	if !IsValidCNPJ(cnpj) {
		t.Fatalf("expected valid CNPJ, got %q", cnpj)
	}
	masked := cnpj[0:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "/" + cnpj[8:12] + "-" + cnpj[12:14]
	if !IsValidCNPJ(masked) {
		t.Fatalf("expected valid masked CNPJ, got %q", masked)
	}
}

func TestIsValidCNPJ_InvalidCheckDigits(t *testing.T) {
	valid := generateValidCNPJ("987654320001")
	invalid := valid[:13] + string(flipDigit(valid[13]))
	if IsValidCNPJ(invalid) {
		t.Fatalf("expected invalid check digits, got %q", invalid)
	}
}

func generateValidCNPJ(base12 string) string {
	if len(base12) != 12 {
		panic("base12 must have 12 digits")
	}
	d := make([]int, 14)
	for i := 0; i < 12; i++ {
		d[i] = int(base12[i] - '0')
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
	d[12] = dv1

	sum = 0
	for i := 0; i < 13; i++ {
		sum += d[i] * w2[i]
	}
	mod = sum % 11
	dv2 := 0
	if mod >= 2 {
		dv2 = 11 - mod
	}
	d[13] = dv2

	out := make([]byte, 14)
	for i := 0; i < 14; i++ {
		out[i] = byte('0' + d[i])
	}
	return string(out)
}

func flipDigit(b byte) byte {
	if b == '9' {
		return '0'
	}
	return b + 1
}
