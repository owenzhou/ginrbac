// 常用工具类
package php

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/*****************加密相关函数 开始******************/

//md5加密
func Md5(str string) string {
	data := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", data)
}

//md5文件加密
func Md5_file(filename string, raw ...bool) (string, error) {
	b, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer b.Close()
	m := md5.New()
	if _, err := io.Copy(m, b); err != nil {
		return "", err
	}
	if len(raw) > 0 && raw[0] {
		return fmt.Sprintf("%b", m.Sum(nil)), nil
	}
	return fmt.Sprintf("%x", m.Sum(nil)), nil
}

//sha1加密
func Sha1(str string) string {
	data := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", data)
}

func Password_hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Password_verify(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*****************加密相关函数 结束******************/

/*****************文件相关函数 开始******************/

//判断是否是文件夹
func Is_dir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

//判断是否是文件
func Is_file(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

//创建文件夹
func Mkdir(dir string, mode os.FileMode, recursive ...bool) (bool, error) {
	var err error
	if len(recursive) > 0 {
		err = os.MkdirAll(dir, mode)
	} else {
		err = os.Mkdir(dir, mode)
	}
	if err != nil {
		return false, err
	}
	return true, err
}

//创建文件或打开文件
func Fopen(name string) (*os.File, error) {
	flag := os.O_APPEND | os.O_CREATE
	f, err := os.OpenFile(name, flag, 0777)
	return f, err
}

//写入文件，与fopen一起使用
func Fwrite(f *os.File, content string) (int, error) {
	return f.Write([]byte(content))
}

//关闭文件
func Fclose(f *os.File) (bool, error) {
	err := f.Close()
	if err != nil {
		return false, err
	}
	return true, err
}

//删除文件
func Unlink(name string) (bool, error) {
	err := os.RemoveAll(name)
	if err != nil {
		return false, err
	}
	return true, err
}

//文件或文件夹重命名
func Rename(oldname, newname string) (bool, error) {
	err := os.Rename(oldname, newname)
	if err != nil {
		return false, err
	}
	return true, err
}

//复制文件或文件夹
func Copy(src, dst string) (int64, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !stat.Mode().IsRegular() {
		err := errors.New("not a regular file or directory")
		return 0, err
	}
	fd, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	fs, err := os.Open(src)
	defer fd.Close()
	defer fs.Close()
	if err != nil {
		return 0, err
	}
	return io.Copy(fd, fs)

}

//删除目录
func Rmdir(path string) error {
	return os.RemoveAll(path)
}

//读取文件
func ReadFile(name string) (string, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(content), err
}

//写文件
func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

//追加写入
func WriteAppendFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	n, _ := f.Seek(0, 2)
	_, err = f.WriteAt([]byte(content), n)
	return err
}

//file_get_contents
func File_get_contents(fileName string, skipVerify ...bool) (string, error) {
	//使用正则判断是网络还是本地
	matched, err := regexp.MatchString(`^http.*`, fileName)
	if err != nil {
		return "", err
	}

	//如果是网络
	if matched {
		var(
			res *http.Response
			err error
		)
		if len(skipVerify) > 0 && skipVerify[0] {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := http.Client{Transport: tr}
			res, err = client.Get(fileName)
		}else{
			res, err = http.Get(fileName)
		}

		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		content, err := io.ReadAll(res.Body)
		return string(content), err
	}

	return ReadFile(fileName)
}

//file_put_contents
func File_put_contents(fileName string, content string) error {
	return WriteFile(fileName, content)
}

//判断文件或目录是否存在
func File_exists(path string) bool {
	if _, err := os.Stat(path);err != nil {
		return false
	}
	return true
}

//更改权限
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

/*****************文件相关函数 结束******************/

/*****************时间相关函数 开始******************/
//返回格式化的时间字符串
func Date(format string, durations ...interface{}) string {
	flag := []string{
		"Y", "2006",
		"y", "2006",
		"m", "01",
		"d", "02",
		"H", "15",
		"i", "04",
		"s", "05",
	}
	replacer := strings.NewReplacer(flag...)
	t := time.Now()
	if len(durations) > 0 {
		switch duration := durations[0].(type) {
		case int:
			t = time.Unix(int64(duration), 0)
		case int64:
			t = time.Unix(duration, 0)
		case string:
			add, _ := time.ParseDuration(duration)
			t = t.Add(add)
		case time.Time:
			t = duration
		}
	}

	return t.Format(replacer.Replace(format))
}

//当前时间戳
func Time() int64 {
	return time.Now().Unix()
}

//时间字符转时间戳
func Strtotime(str string, loc ...string) int64 {
	re := regexp.MustCompile(`([\d-/:]+)`)
	if !re.MatchString(str) {
		return 0
	}

	format := re.ReplaceAllStringFunc(str, func(s string) string {
		if strings.Contains(s, ":") {
			s = "15:04:05"
		}
		if strings.Contains(s, "/") {
			s = "2006/01/02"
		}
		if strings.Contains(s, "-") {
			s = "2006-01-02"
		}
		return s
	})
	//替换后的format空格太多
	if strings.Count(format, " ") > 1 {
		format = "2006-01-02 15:04:05"
	}
	l := "Asia/Shanghai"
	if len(loc) > 0 {
		l = loc[0]
	}
	location, _ := time.LoadLocation(l)
	t, err := time.ParseInLocation(format, str, location)
	if err != nil {
		return 0
	}
	return t.Unix()
}

/*****************时间相关函数 结束******************/

/*****************字符串相关函数 开始******************/

//ascii 转 字符
func Chr(i int) string {
	return fmt.Sprintf("%c", i)
}

//字符转ascii
func Ord(str string) int {
	r := []rune(str)
	return int(r[0])
}

//返回分隔后的字符串数组
func Split(str string, sep string) []string {
	return strings.Split(str, sep)
}

//将字符串数组转为字符串
func Join(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

//字符串转小写
func Strtolower(str string) string {
	return strings.ToLower(str)
}

//字符串转大写
func Strtoupper(str string) string {
	return strings.ToUpper(str)
}

//字符串每个单词首字母转大写
func Ucwords(str string) string {
	return cases.Title(language.Und).String(str)
}

//字符串首字母转大写
func Ucfirst(str string) string {
	b := []byte(str)
	for k, v := range b {
		b[k] = byte(unicode.ToUpper(rune(v)))
		break
	}
	return string(b)
}

//截取字符串
func Substr(str string, start int, length ...int) string {
	strLen := len(str)
	//start
	//正数 - 在字符串的指定位置开始
	//负数 - 在从字符串结尾的指定位置开始
	if start < 0 {
		start = strLen + start
	}
	//length
	//可选。规定要返回的字符串长度。默认是直到字符串的结尾。
	//正数 - 从 start 参数所在的位置返回
	//负数 - 从字符串末端返回
	if len(length) > 0 {
		length0 := length[0]
		end := start + length0
		if length0 < 0 {
			end = strLen + length0
		}
		if end > strLen {
			return str[start:]
		}
		return str[start:end]
	}
	return str[start:]
}

//截取中文字符串
func Mb_substr(str string, start int, length ...int) string {
	ru := []rune(str)
	strLen := len(ru)
	//start
	//正数 - 在字符串的指定位置开始
	//负数 - 在从字符串结尾的指定位置开始
	if start < 0 {
		start = strLen + start
	}
	//length
	//可选。规定要返回的字符串长度。默认是直到字符串的结尾。
	//正数 - 从 start 参数所在的位置返回
	//负数 - 从字符串末端返回
	if len(length) > 0 {
		length0 := length[0]
		end := start + length0
		if length0 < 0 {
			end = strLen + length0
		}
		if end > strLen {
			return string(ru[start:])
		}
		return string(ru[start:end])
	}
	return string(ru[start:])
}

//查找子字符串在字符串中第一次出现的位置
func Strpos(str string, substr string) int {
	return strings.Index(str, substr)
}

//查找子字符串在字符串中最后一次出现的位置
func Strrpos(str string, substr string) int {
	return strings.LastIndex(str, substr)
}

//去掉字符串左右的子字符串
func Trim(str string, sep string) string {
	return strings.Trim(str, sep)
}

//去掉字符串中左边的子字符串
func Ltrim(str string, sep string) string {
	return strings.TrimLeft(str, sep)
}

//去掉字符串中右边的子字符串
func Rtrim(str string, sep string) string {
	return strings.TrimRight(str, sep)
}

//返回指定子字符串及以后的字符串
func Strstr(str string, sep string) string {
	pos := strings.Index(str, sep)
	if pos == -1 {
		return ""
	}
	return str[pos:]
}

//返回字符串的长度
func Strlen(str string) int {
	return len(str)
}

//返回中文字符串的长度
func Mb_strlen(str string) int {
	r := []rune(str)
	return len(r)
}

//替换字符串
func Str_replace(old string, new string, str string) string {
	return strings.ReplaceAll(str, old, new)
}

//字符串反转
func Strrev(str string) string {
	ru := []rune(str)
	if len(ru) <= 0 {
		return ""
	}
	return Strrev(string(ru[1:])) + string(ru[0:1])
}

//返回数组，数组的每个元素的长度为n
func Str_split(str string, n ...int) []string {
	var result []string
	var strs string
	runes := []rune(str)
	if len(n) > 0 {
		for k, v := range runes {
			if (k+1)%n[0] < n[0] {
				strs += string(v)
			}

			if len([]rune(strs)) == n[0] || k == len([]rune(str))-1 {
				result = append(result, strs)
				strs = ""
			}
		}
	} else {
		result = strings.Split(str, "")
	}
	return result
}

//区分大小写的字符串比较
func Strcasecmp(str, str2 string) bool {
	result := strings.Compare(str, str2)
	if result == 0 {
		return true
	} else {
		return false
	}
}

//不区分大小写的字符串比较
func Strcmp(str, str2 string) bool {
	return strings.EqualFold(str, str2)
}

//随机数字
func Mt_rand(min, max int64) int64 {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rd.Int63n(max-min) + min
}

//base64编码
func Base64_encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

//base64解码
func Base64_decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

//struct转json字符串
func Json_encode(in interface{}) (string, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

//json字符串转struct
func Json_decode(str string, out interface{}) error {
	return json.Unmarshal([]byte(str), out)
}

/*****************字符串相关函数 结束******************/

/*****************数组相关函数 开始******************/

//获取接口的数据类型
func Gettype(i interface{}) string {
	return reflect.TypeOf(i).String()
}

//获取接口的具体分类
func Getkind(i interface{}) string {
	return reflect.TypeOf(i).Kind().String()
}

func In_array(need interface{}, arr interface{}) bool {
	kind := Getkind(arr)
	if kind != "array" && kind != "slice" && kind != "map" {
		return false
	}
	val := reflect.ValueOf(arr)
	needVal := reflect.ValueOf(need)
	if kind == "map" {
		mapkeys := val.MapKeys()
		for _, v := range mapkeys {
			if reflect.DeepEqual(needVal.Interface(), val.MapIndex(v).Interface()) {
				return true
			}
		}
		return false
	}

	for i := 0; i < val.Len(); i++ {
		if reflect.DeepEqual(needVal.Interface(), val.Index(i).Interface()) {
			return true
		}
	}

	return false
}

//数组入栈
func Array_push(arr interface{}, elem ...interface{}) {
	kind := Getkind(arr)
	if kind != "ptr" {
		return
	}
	kind = reflect.TypeOf(arr).Elem().Kind().String()
	if kind != "slice" && kind != "array" {
		return
	}

	val := reflect.ValueOf(arr).Elem()
	var elems []reflect.Value
	for _, v := range elem {
		elems = append(elems, reflect.ValueOf(v))
	}
	result := reflect.Append(val, elems...)
	val.Set(result)
}

//数组出栈
func Array_pop(arr interface{}) interface{} {
	kind := Getkind(arr)
	if kind != "ptr" {
		return ""
	}
	kind = reflect.TypeOf(arr).Elem().Kind().String()
	if kind != "slice" && kind != "array" {
		return ""
	}

	val := reflect.ValueOf(arr).Elem()
	result := val.Index(val.Len() - 1)
	newArr := val.Slice(0, val.Len()-1)
	val.Set(newArr)
	return result.Interface()
}

//将数组开头的单元移出数组
func Array_shift(arr interface{}) interface{} {
	kind := Getkind(arr)
	if kind != "ptr" {
		return ""
	}
	kind = reflect.TypeOf(arr).Elem().Kind().String()
	if kind != "slice" && kind != "array" {
		return ""
	}

	val := reflect.ValueOf(arr).Elem()
	result := val.Index(0)
	newArr := val.Slice(1, val.Len())
	val.Set(newArr)
	return result.Interface()
}

//在数组开头插入一个或多个单元
func Array_unshift(arr interface{}, elem ...interface{}) {
	kind := Getkind(arr)
	if kind != "ptr" {
		return
	}
	kind = reflect.TypeOf(arr).Elem().Kind().String()
	if kind != "slice" && kind != "array" {
		return
	}

	val := reflect.ValueOf(arr).Elem()
	if len(elem) <= 0 {
		return
	}
	typ := reflect.TypeOf(arr).Elem()
	elems := reflect.MakeSlice(typ, 0, 0)
	for _, v := range elem {
		elems = reflect.Append(elems, reflect.ValueOf(v))
	}
	result := reflect.AppendSlice(elems, val)
	val.Set(result)
}

//删除数组指定元素
func Unset(arr interface{}, i int) {
	kind := Getkind(arr)
	if kind != "ptr" {
		return
	}
	kind = reflect.TypeOf(arr).Elem().Kind().String()
	if kind != "slice" && kind != "array" {
		return
	}

	val := reflect.ValueOf(arr).Elem()
	result := reflect.AppendSlice(val.Slice(0, i), val.Slice(i+1, val.Len()))
	val.Set(result)
}

//从数组中取出一段
func Array_slice(arr interface{}, n ...int) interface{} {
	kind := Getkind(arr)
	if kind != "slice" && kind != "array" {
		return ""
	}
	val := reflect.ValueOf(arr)
	if len(n) == 1 {
		return val.Slice(n[0], val.Len()).Interface()
	}
	if len(n) >= 2 {
		return val.Slice(n[0], n[1]).Interface()
	}

	return ""
}

//去掉数组中的某一部分并用其它值取代
func Array_splice(arr interface{}, offset int, length int, replacement ...interface{}) interface{} {
	kind := Getkind(arr)
	if kind != "ptr" {
		return ""
	}
	kind = reflect.TypeOf(arr).Elem().Kind().String()
	if kind != "slice" && kind != "array" {
		return ""
	}

	val := reflect.ValueOf(arr).Elem()
	result := val.Slice(offset, offset+length)
	data := reflect.AppendSlice(val.Slice(0, offset), val.Slice(offset+length, val.Len()))

	if len(replacement) > 0 {
		rep := replacement[0]
		kind = Getkind(rep)
		if kind == "slice" || kind == "array" {
			data = reflect.AppendSlice(reflect.AppendSlice(val.Slice(0, offset), reflect.ValueOf(rep)), val.Slice(offset+length, val.Len()))
		} else {
			right := val.Slice(offset+length, val.Len())
			data = reflect.AppendSlice(reflect.Append(val.Slice(0, offset), reflect.ValueOf(rep)), right)
		}
	}

	val.Set(data)
	return result
}

//数组差集
//只支持 []int, []float32, []float64, []string, []byte
func Array_diff(arr1 interface{}, arrs ...interface{}) (interface{}, error) {
	arr1Value := reflect.ValueOf(arr1)
	if arr1Value.Kind().String() != "slice" && arr1Value.Kind().String() != "array" {
		return "", errors.New("first argument is not slice or array")
	}
	for _, v := range arrs {
		result := make([]interface{}, 0)
		itemValue := reflect.ValueOf(v)
		if itemValue.Kind().String() != arr1Value.Kind().String() {
			return "", fmt.Errorf("first argument is %s, current argument is %s", arr1Value.Kind().String(), itemValue.Kind().String())
		}
		for i := 0; i < arr1Value.Len(); i++ {
			finded := false
			for j := 0; j < itemValue.Len(); j++ {
				if reflect.DeepEqual(arr1Value.Index(i).Interface(), itemValue.Index(j).Interface()) {
					finded = true
				}
			}
			if !finded {
				result = append(result, arr1Value.Index(i).Interface())
			}
		}
		arr1Value = reflect.ValueOf(result)
	}
	if Gettype(arr1) == "[]int" {
		result := make([]int, arr1Value.Len())
		for i := 0; i < arr1Value.Len(); i++ {
			result[i] = arr1Value.Index(i).Interface().(int)
		}
		return result, nil
	}
	if Gettype(arr1) == "[]float32" {
		result := make([]float32, arr1Value.Len())
		for i := 0; i < arr1Value.Len(); i++ {
			result[i] = arr1Value.Index(i).Interface().(float32)
		}
		return result, nil
	}
	if Gettype(arr1) == "[]float64" {
		result := make([]float64, arr1Value.Len())
		for i := 0; i < arr1Value.Len(); i++ {
			result[i] = arr1Value.Index(i).Interface().(float64)
		}
		return result, nil
	}
	if Gettype(arr1) == "[]string" {
		result := make([]string, arr1Value.Len())
		for i := 0; i < arr1Value.Len(); i++ {
			result[i] = arr1Value.Index(i).Interface().(string)
		}
		return result, nil
	}
	if Gettype(arr1) == "[]byte" {
		result := make([]byte, arr1Value.Len())
		for i := 0; i < arr1Value.Len(); i++ {
			result[i] = arr1Value.Index(i).Interface().(byte)
		}
		return result, nil
	}
	return arr1Value.Interface(), nil
}

/*****************数组相关函数 结束******************/

/*********************其它函数*************************/

//去除html标签
func Strip_tags(str string) string {
	reg := regexp.MustCompile(`<[\s\S]+?>`)
	if reg.MatchString(str){
		str = reg.ReplaceAllString(str, "")
	}
	return str
}