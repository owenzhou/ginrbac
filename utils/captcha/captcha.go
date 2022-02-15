package captcha

import (
	"bytes"
	"ginrbac/bootstrap/support/facades"
	"ginrbac/bootstrap/utils/php"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"strings"
	"time"

	"github.com/golang/freetype"
)

//Rand生成字符串
func random(i int) string {
	strs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	str := ""
	rand.Seed(time.Now().UnixNano())
	for j := 0; j < i; j++ {
		rd := rand.Int63n(int64(len(strs) - 1))
		str += string(strs[rd])
	}
	return str
}

//获取验证码图片
func NewCaptcha() map[string]string {
	var result = map[string]string{"code": "", "data": ""}

	//获取验证码字符串
	randStr := random(facades.Config.Captcha.Num)

	width := facades.Config.Captcha.Width
	height := facades.Config.Captcha.Height
	secretKey := facades.Config.Captcha.SecretKey
	// 创建画布
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			//给画布填充背景颜色
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	/*
		//获取embed字体文件
		fontFS, err := fs.Sub(facades.Views, "views/layouts/assets/font")
		if err != nil {
			return result
		}
		//读取字体文件内容
		fontBytes, err := fs.ReadFile(fontFS, "Duality.ttf")
		fmt.Println(len(fontBytes), len(fontBase64))
		if err != nil {
			return result
		}
	*/
	fontBytes := []byte(php.Base64_decode(fontBase64))
	//freetype 转为字体
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return result
	}

	//调整字体并写入画布
	f := freetype.NewContext()
	f.SetDPI(92)
	f.SetFont(font)
	f.SetFontSize(26)
	f.SetClip(img.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}))

	pt := freetype.Pt(img.Bounds().Dx()-width+5, img.Bounds().Dy()-(height/2)+10)
	f.DrawString(randStr, pt)

	fpng := bytes.NewBuffer([]byte{})

	if err := png.Encode(fpng, img); err != nil {
		return result
	}
	pngStr := fpng.String()
	fpng.Reset()

	t := php.Date("Y-m-d H:i:s")
	result["code"] = php.Base64_encode(php.Md5(secretKey+php.Strtolower(randStr)+t) + "&" + t)
	result["data"] = "data:image/png;base64," + php.Base64_encode(pngStr)
	return result
}

func Validate(captcha, code string) bool {
	decode := php.Base64_decode(code)
	s := strings.Split(decode, "&")
	if len(s) < 2 {
		return false
	}
	validateCode := s[0]
	generateTime := s[1]

	expired := facades.Config.Captcha.Expired
	generateTimestamp := php.Strtotime(generateTime)
	now := php.Time()
	if now-generateTimestamp > expired {
		return false
	}
	if md5Cap := php.Md5(facades.Config.Captcha.SecretKey + php.Strtolower(captcha) + generateTime); md5Cap == validateCode {
		return true
	}
	return false
}
