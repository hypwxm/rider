package number

import "testing"

func TestIsPosInt(t *testing.T) {
	var testStr []string = []string{
		"asdada",
		"ad9090ada",
		"123jkj123",
		"0",
		"00",
		"000",
		"哈哈",
		"-123",
		"-0",
		"-000",
		"-00123",
		"-01",
		"0.123",
		"1.1",
		"1.0",
		"",
	}
	for k, _ := range testStr {
		if IsPosInt(testStr[k]) {
			t.Error("正整数匹配错误")
		}
	}

	testStr = []string{
		"123",
		"010",
		"01",
		"10",
		"00001",
		"1000",
		"123010",
		"00102130",
	}
	for k, _ := range testStr {
		if !IsPosInt(testStr[k]) {
			t.Error("正整数匹配错误")
		}
	}
}

func TestIsNonNegInt(t *testing.T) {
	var testStr []string = []string{
		"asdada",
		"ad9090ada",
		"123jkj123",
		"哈哈",
		"-123",
		"-00123",
		"-01",
		"1.23",
		"-0.1",
		"0.1",
		"1.01",
		"",
	}
	for k, _ := range testStr {
		if IsNonNegInt(testStr[k]) {
			t.Error("正整数匹配错误")
		}
	}

	testStr = []string{
		"123",
		"010",
		"01",
		"10",
		"00001",
		"1000",
		"123010",
		"00102130",
		"-0",
		"-00",
	}
	for k, _ := range testStr {
		if !IsNonNegInt(testStr[k]) {
			t.Error("正整数匹配错误")
		}
	}
}

func TestIsPosNum(t *testing.T) {
	var testStr []string = []string{
		"asdada",
		"ad9090ada",
		"123jkj123",
		"哈哈",
		"-123",
		"-00123",
		"-01",
		"-0.1",
		"0",
		".00",
		"0.00",
		"0.",
		"",
	}
	for k, _ := range testStr {
		if IsPosNum(testStr[k]) {
			t.Error("正数匹配错误", testStr[k])
		}
	}

	testStr = []string{
		"123",
		"010",
		"01",
		"10",
		"00001",
		"1000",
		"123010",
		"00102130",
		"0.1",
		"0.01",
		"1.01",
		"01.01",
		".01",
		"0.1",
		"1.",
		"010.",
		"10.00",
	}
	for k, _ := range testStr {
		if !IsPosNum(testStr[k]) {
			t.Error("正数匹配错误", testStr[k])
		}
	}
}

func TestIsNonNeg(t *testing.T) {
	var testStr []string = []string{
		"asdada",
		"ad9090ada",
		"123jkj123",
		"哈哈",
		"-123",
		"-00123",
		"-01",
		"-0.1",
		"",
	}
	for k, _ := range testStr {
		if IsNonNeg(testStr[k]) {
			t.Error("非负数匹配错误", testStr[k])
		}
	}

	testStr = []string{
		"123",
		"010",
		"01",
		"10",
		"00001",
		"1000",
		"123010",
		"00102130",
		"0.1",
		"0.01",
		"1.01",
		"01.01",
		".01",
		"0",
		"0.0",
		"0.",
		".0",
	}
	for k, _ := range testStr {
		if !IsNonNeg(testStr[k]) {
			t.Error("非负数匹配错误", testStr[k])
		}
	}
}

func BenchmarkIsNonNeg(b *testing.B) {
	var testStr []string = []string{
		"asdada",
		"ad9090ada",
		"123jkj123",
		"哈哈",
		"-123",
		"-00123",
		"-01",
		"-0.1",
		"",
	}
	for k, _ := range testStr {
		if IsNonNeg(testStr[k]) {
			b.Error("非负数匹配错误", testStr[k])
		}
	}

	testStr = []string{
		"123",
		"010",
		"01",
		"10",
		"00001",
		"1000",
		"123010",
		"00102130",
		"0.1",
		"0.01",
		"1.01",
		"01.01",
		".01",
		"0",
		"0.0",
		"0.",
		".0",
	}

	for k, _ := range testStr {
		if !IsNonNeg(testStr[k]) {
			b.Error("非负数匹配错误", testStr[k])
		}
	}

}
