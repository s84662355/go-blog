package helper

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// 判断是否为slcie数据
func isSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return val, ok
}

//是否是文件
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

//获取随机数  数字和文字
func GetRandomBoth(n int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//获取本机ip
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func Recover() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

//sha1加密
func Sha1En(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//分页方法，根据传递过来的页数，每页数，总数，返回分页的内容 7个页数 前 1，2，3，4，5 后 的格式返回,小于5页返回具体页数
func Paginator(page int, prepage int, nums int64) map[string]int {

	var firstpage int //前一页地址
	var lastpage int  //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //最后5页
		start := totalpages - 5 + 1
		firstpage = page - 1
		lastpage = int(math.Min(float64(totalpages), float64(page+1)))
		pages = make([]int, 5)
		for i := range pages {
			pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		pages = make([]int, 5)
		firstpage = page - 3
		for i := range pages {
			pages[i] = start + i
		}
		firstpage = page - 1
		lastpage = page + 1
	default:
		/*
		   pages = make([]int, int(math.Min(5, float64(totalpages))))
		   for i, _ := range pages {
		       pages[i] = i + 1
		   }
		*/
		firstpage = int(math.Max(float64(1), float64(page-1)))
		lastpage = page + 1
		//fmt.Println(pages)
	}
	paginatorMap := make(map[string]int)
	//paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["firstpage"] = firstpage
	paginatorMap["lastpage"] = lastpage
	paginatorMap["currpage"] = page
	paginatorMap["total"] = int(nums)
	paginatorMap["size"] = int(prepage)

	return paginatorMap
}

//截取指定字符子串
func SubstrContains(s, substr string) string {
	n := strings.Index(s, substr)
	return s[n:]
}

func readXlsx(filename string) map[string][]string {
	var listOra map[string][]string
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}
	for _, sheet := range xlFile.Sheets {

		//fmt.Printf("Sheet Name: %s\n", sheet.Name)

		// 获取标签页(时间)
		//tmpOra.TIME = sheet.Name
		for _, row := range sheet.Rows {

			var strs []string

			for _, cell := range row.Cells {
				text := cell.String()
				strs = append(strs, text)
			}

			listOra[sheet.Name] = strs
		}
	}
	return listOra
}

/*
func writingXlsx(oraList [][]string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.SetHeightCM(0.5)
	cell = row.AddCell()
	cell.Value = "单位"
	cell = row.AddCell()
	cell.Value = "业务系统"
	cell = row.AddCell()
	cell.Value = "进程名"
	cell = row.AddCell()
	cell.Value = "V1000"
	cell = row.AddCell()
	cell.Value = "V2000"
	cell = row.AddCell()
	cell.Value = "H"
	cell = row.AddCell()
	cell.Value = "L"
	cell = row.AddCell()
	cell.Value = "A"
	cell = row.AddCell()
	cell.Value = "TIME"

	for _, i := range oraList {
		if i.corp == "单位" {
			continue
		}

		// 判断是否为-9999，是的变为0.0
		var row1 *xlsx.Row
		if i.v1000 == "-9999" {
			i.v1000 = "0.0"
		}
		if i.v2000 == "-9999" {
			i.v2000 = "0.0"
		}
		if i.H == "-9999" {
			i.H = "0.0"
		}
		if i.L == "-9999" {
			i.L = "0.0"
		}

		row1 = sheet.AddRow()
		row1.SetHeightCM(0.5)

		cell = row1.AddCell()
		cell.Value = i.corp
		cell = row1.AddCell()
		cell.Value = i.name
		cell = row1.AddCell()
		cell.Value = i.name2

		// 判断值是大于7200，大于变成红色
		v1, _ := strconv.ParseFloat(i.v1000, 64)
		if v1 > 7200 {
			cell = row1.AddCell()
			cell.Value = i.v1000
			cell.GetStyle().Font.Color = "00FF0000"
		} else {
			cell = row1.AddCell()
			cell.Value = i.v1000
		}

		//v2, _ := strconv.Atoi(i.v2000)
		v2, _ := strconv.ParseFloat(i.v2000, 64)
		if v2 > 7200 {
			cell = row1.AddCell()
			cell.Value = i.v2000
			cell.GetStyle().Font.Color = "00FF0000"
		} else {
			cell = row1.AddCell()
			cell.Value = i.v2000
		}

		//vH, _ := strconv.Atoi(i.H)
		vH, _ := strconv.ParseFloat(i.H, 64)
		if vH > 7200 {
			cell = row1.AddCell()
			cell.Value = i.H
			cell.GetStyle().Font.Color = "00FF0000"
		} else {
			cell = row1.AddCell()
			cell.Value = i.H
		}

		//vL, _ := strconv.Atoi(i.L)
		vL, _ := strconv.ParseFloat(i.L, 64)
		if vL > 7200 {
			cell = row1.AddCell()
			cell.Value = i.L
			cell.GetStyle().Font.Color = "00FF0000"
		} else {
			cell = row1.AddCell()
			cell.Value = i.L

		}

		//vA, _ := strconv.Atoi(i.A)
		vA, _ := strconv.ParseFloat(i.A, 64)
		if vA > 7200 {
			cell = row1.AddCell()
			cell.Value = i.A
			cell.GetStyle().Font.Color = "00FF0000"
		} else {
			cell = row1.AddCell()
			cell.Value = i.A
		}

		// 打印时间
		cell = row1.AddCell()
		cell.Value = i.TIME
	}

	err = file.Save("2019-_-_-2019-_-_Lag延时数据.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

*/
