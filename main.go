package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huichen/pinyin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"net/http"
	"strings"
	"gopkg.in/gcfg.v1"
	"time"
)

// 接收去GET到的json
type Devroom struct {
	Devs []struct {
		Mac string `json:"mac"`
		SubMac string `json:"subMac"`
		Key string `json:"key"`
		Mid string `json:"mid"`
		DevName string `json:"devName"`
		RoomName string `json:"roomName"`
	} `json:"devs"`
	R int `json:"r"`
}
// 客户端post本地服务参数
type Queryparams struct {
	Query string `json:"query"`
	Mac string `json:"mac"`
	HomeId string `json:"home_id"`
}
// 本地服务输出返回json
type ResponseDev struct {
	Code     int    `json:"code"`
	OriQuery string `json:"oriQuery"`
	OriPy    string `json:"oriPy"`
	Query    string `json:"query,omitempty"`
	Devs     []TarDev `json:"devs,omitempty"`
		//Mac string `json:"mac"`
		//SubMac string `json:"subMac"`
		//Key string `json:"key"`
		//Mid string `json:"mid"`
		//DevName string `json:"devName"`
		//RoomName string `json:"roomName"`
		//MidDes string `json:"midDes"`
		//MidSe string `json:"midSe"`
		//DevNamePy string `json:"devNamePy"`
	TarDevs []TarDev `json:"tarDevs,omitempty"`
		//Mac string `json:"mac"`
		//SubMac string `json:"subMac"`
		//Mid string `json:"mid"`
		//Key string `json:"key"`
		//DevName string `json:"devName"`
		//RoomName string `json:"roomName"`
		//MidDes string `json:"midDes"`
		//MidSe string `json:"midSe"`
		//DevNamePy string `json:"devNamePy"`
}

type TarDev struct {
	Mac string `json:"mac"`
	SubMac string `json:"subMac"`
	Mid string `json:"mid"`
	Key string `json:"key"`
	DevName string `json:"devName"`
	RoomName string `json:"roomName"`
	MidDes string `json:"midDes"`
	MidSe string `json:"midSe"`
	DevNamePy string `json:"devNamePy"`
}
// 数据库中模型绑定
type mid_user struct {
	Mid string `json:"mid"`
	Definedname string `json:"definedname"`
	Midse string `json:"midSe"`
}

var (
	DB *gorm.DB
)

type Config struct {
	Custom struct{
		Username string
		Password string
		Sqlhost string
		Sqlport string
		Dbname string
		Geturl string
		Port string
	}
}

func InitMySQL(dsn string)(err error){
	// 绑定数据库 (本人数据库没有密码，库名为mid_dev)

	DB, err = gorm.Open("mysql",dsn)
	if err != nil{
		return err
	}
	err = DB.DB().Ping()
	return
}

func Py(s string) string {
	// 拼音声明
	var py pinyin.Pinyin
	var devpy string
	// 初始化，载入汉字拼音映射文件
	py.Init("data/pinyin_table.txt")

	pyrune := []rune(s)
	for j := 0;j < len(pyrune);j ++ {
		devpy += py.GetPinyin(pyrune[j],false)
	}
	return devpy
}

func main() {
	// 加载配置文件
	var config Config
	err := gcfg.ReadFileInto(&config, "config.ini")
	if err != nil {
		fmt.Println("Failed to parse config file: %s", err)
	}
	dsn := config.Custom.Username + ":" + config.Custom.Password + "@tcp(" + config.Custom.Sqlhost + ":" + config.Custom.Sqlport + ")/" + config.Custom.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	port := ":" + config.Custom.Port

	// 连接、绑定数据库（库名mid_dev）

	err = InitMySQL(dsn)
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	DB.AutoMigrate(&mid_user{})

	// 建立默认路由
	r := gin.Default()

	r.POST("/semantic/custom", func(c *gin.Context) {
		// 声明一个Queryparams变量
		var u Queryparams
		err = c.ShouldBind(&u)
		if err != nil {
			panic(err)
		}
		// 声明一个本地的返回（查询mac后有devname的时候）
		var localreturn ResponseDev
		// devrep是通过pickup函数查询mac后返回回来的body

		client := &http.Client{Timeout: time.Second}

		req, err := http.NewRequest("GET",config.Custom.Geturl,nil)
		if err != nil {
			var nomac ResponseDev
			nomac.Code = 1501
			nomac.OriQuery = u.Query
			nomac.OriPy = Py(u.Query)
			fmt.Println(err)
			c.JSON(http.StatusOK,nomac)
			return
		}
		q := req.URL.Query()

		q.Set("mac",u.Mac)
		req.URL.RawQuery = q.Encode()

		resp, er := client.Do(req)
		if er != nil {
			var nomac ResponseDev
			nomac.Code = 1501
			nomac.OriQuery = u.Query
			nomac.OriPy = Py(u.Query)
			fmt.Println(er)
			c.JSON(http.StatusOK,nomac)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)


		//params := url.Values{}
		//Url, err := url.Parse("")
		//if err != nil {
		//	var nomac ResponseDev
		//	nomac.Code = 1501
		//	nomac.OriQuery = u.Query
		//	nomac.OriPy = Py(u.Query)
		//	fmt.Println(err)
		//	c.JSON(http.StatusOK,nomac)
		//	return
		//}
		//params.Set("mac",u.Mac)
		//Url.RawQuery = params.Encode()
		//urlPath := Url.String()
		////fmt.Println(urlPath)
		//
		////client := http.Client{Timeout: time.Second}
		//
		//resp,err := http.Get(urlPath)
		//if err != nil {
		//	var nomac ResponseDev
		//	nomac.Code = 1501
		//	nomac.OriQuery = u.Query
		//	nomac.OriPy = Py(u.Query)
		//	fmt.Println(err)
		//	c.JSON(http.StatusOK,nomac)
		//	return
		//}
		//defer resp.Body.Close()
		//body, _ := ioutil.ReadAll(resp.Body)
		////fmt.Println(string(body))
		var devrep Devroom
		err = json.Unmarshal(body,&devrep)

		// 返回JSON格式的body,绑定结构体

		var n = len(devrep.Devs)
		// 如果n为0，则说明返回的没有devName
		if n == 0 {
			var nodev ResponseDev
			nodev.Code = 1401
			nodev.OriQuery = u.Query
			nodev.OriPy = Py(u.Query)
			//nodev.Devs = devrep.Devs
			c.JSON(http.StatusOK,nodev)
		}else {
			var mark []int
			max := 0
			// 通过比较找出所有query包含的devname绑定的mid
			for i := 0;i < n;i ++ {
				if strings.Contains(u.Query,devrep.Devs[i].DevName) {
					if len(devrep.Devs[i].DevName) > max {
						max = len(devrep.Devs[i].DevName)
						mark = append(mark[0:0],i)
					}else if len(devrep.Devs[i].DevName) == max {
						mark = append(mark,i)
					}
				}
			}
			if len(mark) == 0 {
				localreturn.Code = 1301
				localreturn.OriQuery = u.Query
				localreturn.OriPy = Py(u.Query)
				//localreturn.TarDevs = nil
				for i := 0;i < n;i ++ {
					var sample TarDev
					sample.Mid = devrep.Devs[i].Mid
					sample.Mac = devrep.Devs[i].Mac
					sample.DevName = devrep.Devs[i].DevName
					sample.DevNamePy = Py(devrep.Devs[i].DevName)
					sample.Key = devrep.Devs[i].Key
					sample.SubMac = devrep.Devs[i].SubMac
					sample.RoomName = devrep.Devs[i].RoomName

					var sqlfind mid_user
					DB.Table("mid_user").Select("definedname,midse").Where("mid = ?",devrep.Devs[i].Mid).First(&sqlfind)
					sample.MidSe = sqlfind.Midse
					sample.MidDes = sqlfind.Definedname

					localreturn.Devs = append(localreturn.Devs,sample)

				}
				//localreturn.Devs = devrep.Devs
				c.JSON(http.StatusOK,localreturn)
			}else if len(mark) == 1 {
				// 通过mid去数据库中查询，然后把相应的name返回赋值
				localreturn.Code = 0
				localreturn.OriQuery = u.Query
				localreturn.OriPy = Py(u.Query)


				var sqlfind mid_user
				DB.Table("mid_user").Select("definedname, midse").Where("mid = ?",devrep.Devs[mark[0]].Mid).First(&sqlfind)

				var box TarDev
				box.Mid = devrep.Devs[mark[0]].Mid
				box.SubMac = ""
				box.Mac = devrep.Devs[mark[0]].Mac
				box.Key = devrep.Devs[mark[0]].Key
				box.DevName = devrep.Devs[mark[0]].DevName
				box.RoomName = devrep.Devs[mark[0]].RoomName
				box.MidDes = sqlfind.Definedname
				box.DevNamePy = Py(devrep.Devs[mark[0]].DevName)
				box.MidSe = sqlfind.Midse

				//fmt.Println(sqlfind.Definedname,sqlfind.MidSe)
				var replace string
				replace = strings.ReplaceAll(u.Query, box.DevName, box.MidDes)
				localreturn.Query = replace

				localreturn.TarDevs = append(localreturn.TarDevs,box)
				c.JSON(http.StatusOK,localreturn)
			}else {
				// 通过mid去数据库中查询，然后把相应的name返回赋值
				localreturn.Code = 0
				localreturn.OriQuery = u.Query
				localreturn.OriPy = Py(u.Query)
				localreturn.Query = u.Query
				//localreturn.Devs = nil
				//localreturn.TarDevs = make([]TarDev, 0, len(mark))
				for i := 0;i < len(mark);i ++ {
					var sqlfind mid_user
					DB.Table("mid_user").Select("definedname,midse").Where("mid = ?",devrep.Devs[mark[i]].Mid).First(&sqlfind)
					//fmt.Println(sqlfind)
					var box TarDev
					box.Mid = devrep.Devs[mark[i]].Mid
					box.SubMac = devrep.Devs[mark[i]].SubMac
					box.Mac = devrep.Devs[mark[i]].Mac
					box.Key = devrep.Devs[mark[i]].Key
					box.DevName = devrep.Devs[mark[i]].DevName
					box.RoomName = devrep.Devs[mark[i]].RoomName
					box.MidDes = sqlfind.Definedname
					box.DevNamePy = Py(devrep.Devs[mark[i]].DevName)
					box.MidSe = sqlfind.Midse

					localreturn.TarDevs = append(localreturn.TarDevs,box)
				}
				c.JSON(http.StatusOK,localreturn)
			}
		}
	})
	r.Run(port)
}
