package logger

import "fmt"


//其中0x1B是标记，[开始定义颜色，1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色。显示效果为：

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

func BlackText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;30m%s%c[0m", 0x1B, v, 0x1B)
}

func RedText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;31m%s%c[0m", 0x1B, v, 0x1B)
}

func GreenText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;32m%s%c[0m", 0x1B, v, 0x1B)
}

func YellowText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;33m%s%c[0m", 0x1B, v, 0x1B)
}

func BlueText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;34m%s%c[0m", 0x1B, v, 0x1B)
}

func PurpleText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;35m%s%c[0m", 0x1B, v, 0x1B)
}

func CyanText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;36m%s%c[0m", 0x1B, v, 0x1B)
}

func WhiteText(v interface{}) string {
	return fmt.Sprintf("%c[0;0;37m%s%c[0m", 0x1B, v, 0x1B)
}


func BlackBg(v interface{}) string {
	return fmt.Sprintf("%c[0;40;0m%s%c[0m", 0x1B, v, 0x1B)
}

func RedBg(v interface{}) string {
	return fmt.Sprintf("%c[0;41;0m%s%c[0m", 0x1B, v, 0x1B)
}

func GreenBg(v interface{}) string {
	return fmt.Sprintf("%c[0;42;0m%s%c[0m", 0x1B, v, 0x1B)
}

func YellowBg(v interface{}) string {
	return fmt.Sprintf("%c[0;43;0m%s%c[0m", 0x1B, v, 0x1B)
}

func BlueBg(v interface{}) string {
	return fmt.Sprintf("%c[0;44;0m%s%c[0m", 0x1B, v, 0x1B)
}

func PurpleBg(v interface{}) string {
	return fmt.Sprintf("%c[0;45;0m%s%c[0m", 0x1B, v, 0x1B)
}

func CyanBg(v interface{}) string {
	return fmt.Sprintf("%c[0;46;0m%s%c[0m", 0x1B, v, 0x1B)
}

func WhiteBg(v interface{}) string {
	return fmt.Sprintf("%c[0;47;0m%s%c[0m", 0x1B, v, 0x1B)
}


func BlackBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;30m%s%c[0m", 0x1B, v, 0x1B)
}

func RedBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;31m%s%c[0m", 0x1B, v, 0x1B)
}

func GreenBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;32m%s%c[0m", 0x1B, v, 0x1B)
}

func YellowBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;33m%s%c[0m", 0x1B, v, 0x1B)
}

func BlueBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;34m%s%c[0m", 0x1B, v, 0x1B)
}

func PurpleBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;35m%s%c[0m", 0x1B, v, 0x1B)
}

func CyanBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;36m%s%c[0m", 0x1B, v, 0x1B)
}

func WhiteBoldText(v interface{}) string {
	return fmt.Sprintf("%c[1;0;37m%s%c[0m", 0x1B, v, 0x1B)
}




func BlackUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;30m%s%c[0m", 0x1B, v, 0x1B)
}

func RedUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;31m%s%c[0m", 0x1B, v, 0x1B)
}

func GreenUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;32m%s%c[0m", 0x1B, v, 0x1B)
}

func YellowUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;33m%s%c[0m", 0x1B, v, 0x1B)
}

func BlueUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;34m%s%c[0m", 0x1B, v, 0x1B)
}

func PurpleUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;35m%s%c[0m", 0x1B, v, 0x1B)
}

func CyanUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;36m%s%c[0m", 0x1B, v, 0x1B)
}

func WhiteUnderlineText(v interface{}) string {
	return fmt.Sprintf("%c[4;0;37m%s%c[0m", 0x1B, v, 0x1B)
}



func BlackFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;30m%s%c[0m", 0x1B, v, 0x1B)
}

func RedFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;31m%s%c[0m", 0x1B, v, 0x1B)
}

func GreenFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;32m%s%c[0m", 0x1B, v, 0x1B)
}

func YellowFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;33m%s%c[0m", 0x1B, v, 0x1B)
}

func BlueFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;34m%s%c[0m", 0x1B, v, 0x1B)
}

func PurpleFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;35m%s%c[0m", 0x1B, v, 0x1B)
}

func CyanFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;36m%s%c[0m", 0x1B, v, 0x1B)
}

func WhiteFlashText(v interface{}) string {
	return fmt.Sprintf("%c[5;0;37m%s%c[0m", 0x1B, v, 0x1B)
}



func BlackAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;30m%s%c[0m", 0x1B, v, 0x1B)
}

func RedAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;31m%s%c[0m", 0x1B, v, 0x1B)
}

func GreenAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;32m%s%c[0m", 0x1B, v, 0x1B)
}

func YellowAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;33m%s%c[0m", 0x1B, v, 0x1B)
}

func BlueAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;34m%s%c[0m", 0x1B, v, 0x1B)
}

func PurpleAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;35m%s%c[0m", 0x1B, v, 0x1B)
}

func CyanAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;36m%s%c[0m", 0x1B, v, 0x1B)
}

func WhiteAntiWhiteText(v interface{}) string {
	return fmt.Sprintf("%c[7;0;37m%s%c[0m", 0x1B, v, 0x1B)
}