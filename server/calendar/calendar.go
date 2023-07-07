package calendar

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type OwnTime struct {
	Year  string
	Month string
	Day   string
	Week  int
}

func Init(log *zap.Logger) *OwnTime {

	resp, err := http.Get("http://worldtimeapi.org/api/timezone/Asia/Shanghai")
	if err != nil {
		log.Sugar().Errorf("获取网络时间失败：%v", err)
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Sugar().Errorf("获取网络时间失败：%v", err)
	}

	// 从结果中获取时间字段
	timeString := result["datetime"].(string)
	dayOfWeek := int(result["day_of_week"].(float64))

	// 解析时间字符串
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		log.Sugar().Errorf("获取网络时间失败：%v", err)
	}

	yearMonthDay := strings.Split(t.Format("2006-01-02"), "-")

	ownTime := OwnTime{
		Year:  yearMonthDay[0],
		Month: yearMonthDay[1],
		Day:   yearMonthDay[2],
		Week:  dayOfWeek,
	}

	return &ownTime
}
