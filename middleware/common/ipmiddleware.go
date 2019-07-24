package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"snack/db"
	"time"
)

const (
	APPCODE = "0798ec0cf1bf420ab8adca0cdfbb61b7"
)

type AreaInfo struct {
	Area string `json:"area"`
	City string `json:"city"`
	CityID string `json:"city_id"`
	Country string `json:"country"`
	CountryID string `json:"country_id"`
	IP string `json:"ip"`
	ISP string `json:"isp"`
	LongIP string `json:"long_ip"`
	Region string `json:"region"`
	RegionID string `json:"region_id"`
	CreateTime int64 `json:"create_time"`
}

type Data struct {
	Data AreaInfo `json:"data"`
	LogID string `json:"log_id"`
	Msg string `json:"msg"`
	Ret int `json:"ret"`
}

func CheckIp(c *gin.Context) {

	go func() {
		ip := c.Request.Header.Get("X-Forwarded-For")

		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://api01.aliyun.venuscn.com/ip", nil)
		if err != nil {
			log.Println("ERROR: request aliyun ip parse failed, error: ", err)
			return
		}

		query := req.URL.Query()
		query.Add("ip", ip)
		req.URL.RawQuery = query.Encode()
		req.Header.Set("Authorization", "APPCODE "+APPCODE)

		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Println("ERROR: get aliyun ip parse response failed, error: ", err)
			return
		}

		if resp.StatusCode == http.StatusUnauthorized {
			log.Println("ERROR: aliyun ip parse Authorization failed")
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ERROR: get aliyun ip parse data failed, error: ", err)
			return
		}
		var dataEntry Data
		if err := json.Unmarshal(data, &dataEntry); err != nil {
			log.Println("ERROR: unmarshal aliyun ip data failed, error: ", err)
			return
		}

		if dataEntry.Ret != 200 {
			log.Println("ERROR: aliyun ip parse response err: ", dataEntry.Msg)
			return
		}

		session := db.GetMgoSession()
		defer session.Close()

		dataEntry.Data.CreateTime = time.Now().Unix()
		if err := session.DB(db.DB_NAME).C(db.COLL_IP).Insert(dataEntry.Data); err != nil {
			log.Println("ERROR: save aliyun ip parse failed: ", err)
			return
		}
	}()

	c.Next()
}
