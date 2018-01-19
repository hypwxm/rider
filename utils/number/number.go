package number

import "regexp"

const (
	// 匹配整数
	interger string = "^-?[0-9]+$"
	// 匹配正整数
	postiveInterger string = "^[0-9]*[1-9]{1,}[0-9]*$"

	// 匹配非负整数
	nonNegativeInteger string = "^([0-9]+|(-0+))$"

	// 匹配非负数
	nonNegative string = "^[0-9]*([0-9]+|\\.[0-9]+|[0-9]+\\.)[0-9]*$"

	// 匹配正数
	postiveNum string = "^[0-9]*([1-9]+|\\.[0-9]*[1-9]+|[1-9]+\\.)[0-9]*$"

	// 匹配11位手机号
	phoneNumber11 string = "^1[0-9]{10}$"

	// 匹配18位身份证号
	idCardNumber string = "^[0-9]{17}[a-zA-Z0-9]$"
)

var (
	regInteger *regexp.Regexp

	regPosInt *regexp.Regexp

	regNonNegInt *regexp.Regexp

	regNonNeg *regexp.Regexp

	regPosNum *regexp.Regexp

	regPhoneNum11 *regexp.Regexp

	regIdCardNum *regexp.Regexp
)

func init() {
	regInteger, _ = regexp.Compile(interger)
	regPosInt, _ = regexp.Compile(postiveInterger)
	regNonNegInt, _ = regexp.Compile(nonNegativeInteger)
	regNonNeg, _ = regexp.Compile(nonNegative)
	regPosNum, _ = regexp.Compile(postiveNum)
	regPhoneNum11, _ = regexp.Compile(phoneNumber11)
	regIdCardNum, _ = regexp.Compile(idCardNumber)
}

// 判断字符串是不是整数
func IsInteger(numStr string) bool {
	return regInteger.MatchString(numStr)
}

// 判断字符串是否为正整数
func IsPosInt(numStr string) bool {
	return regPosInt.MatchString(numStr)
}

// 判断非负整数
func IsNonNegInt(numStr string) bool {
	return regNonNegInt.MatchString(numStr)
}

// 判断是否为非负数
func IsNonNeg(numStr string) bool {
	return regNonNeg.MatchString(numStr)
}

// 判断是否为正数
func IsPosNum(numStr string) bool {
	return regPosNum.MatchString(numStr)
}

// 判断是否为11为手机号码
func Is11PhoneNum(numStr string) bool {
	return regPhoneNum11.MatchString(numStr)
}

// 判断是否为18位身份证号
func IsIdCardNum(numStr string) bool {
	return regIdCardNum.MatchString(numStr)
}
