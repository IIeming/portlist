package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"portlist/calendar"
	"portlist/data"
	"portlist/logger"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	log      *zap.Logger
	current  *calendar.OwnTime
	interval int
	port     string
)

func toolWrite(post string, add int) (string, error) {
	num, _ := strconv.Atoi(post)
	value := fmt.Sprint(num + add)
	if toolGetPosts(value) {
		err := os.WriteFile(".config", []byte(value), 0644)
		return value, err
	}
	return value, fmt.Errorf("the port has listen")
}

func toolRead() ([]byte, error) {
	// 从.config文件读取内容
	data, err := os.ReadFile(".config")
	if err != nil {
		log.Error(fmt.Sprint(err))
	}
	return data, err
}

func toolGetPosts(post string) bool {
	cmd := exec.Command("bash", "-c", "ss -ntl|tail -n +2|awk '{print $4}'|awk -F':' '{print $NF}'")
	// 执行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		log.Error("failed to perform the command")
	}
	posts := strings.Split(string(output), "\n")

	for _, v := range posts {
		if v == post {
			log.Error("the port has listen")
			return false
		}
	}
	return true
}

func toolJsonOut(value, msg string) string {
	data := map[string]interface{}{
		"value": value,
		"msg":   msg,
	}
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func vacationInterval() (interval, weekday, festival string) {
	weekOfDate := data.WeekDays
	holiday := data.Holiday()
	mon, _ := strconv.Atoi(current.Month)
	day, _ := strconv.Atoi(current.Day)
	for i := mon; i <= 12; i++ {
		var index string
		if i < 10 {
			index = fmt.Sprintf("0%v", i)
		} else {
			index = fmt.Sprint(i)
		}
		if _, ok := (*holiday)[index]; ok {
			for _, v := range (*holiday)[index] {
				if current.Month == index {
					if day < v.StartDay {
						startDate := current.Year + "-" + current.Month + "-" + current.Day
						endDate := current.Year + fmt.Sprintf("-%v-%v", index, v.StartDay)
						return dayInterval(startDate, endDate), weekOfDate[current.Week], v.Festival
					} else if day <= v.EndDay {
						return "0", weekOfDate[current.Week], v.Festival
					}
				} else {
					startDate := current.Year + "-" + current.Month + "-" + current.Day
					endDate := current.Year + fmt.Sprintf("-%v-%v", index, v.StartDay)
					return dayInterval(startDate, endDate), weekOfDate[current.Week], v.Festival
				}
			}
		}
	}
	return "365", weekOfDate[current.Week], "下一个"
}

func dayInterval(startDate, endDate string) string {
	layout := "2006-01-02"
	startTime, _ := time.Parse(layout, startDate)
	endTime, _ := time.Parse(layout, endDate)
	duration := endTime.Sub(startTime)
	days := int(duration.Hours() / 24)
	return fmt.Sprint(days)
}

func handleAddPort(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := toolRead()
	if value, err := toolWrite(string(data), interval); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
	} else {
		log.Sugar().Infof("%v,端口更新成功", value)
		fmt.Fprintln(w, toolJsonOut(value, "端口更新成功"))
	}
}

func handleModPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := r.FormValue("curvePost")
	if value, err := toolWrite(string(post), 0); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
	} else {
		log.Sugar().Infof("%v,端口更新成功", value)
		fmt.Fprintln(w, toolJsonOut(value, "端口更新成功"))
	}
}

func handleGetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if value, err := toolRead(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	} else {
		log.Sugar().Infof("%v,端口更新成功", value)
		fmt.Fprintln(w, toolJsonOut(string(value), "端口获取成功"))
	}
}

func handleGetDate(w http.ResponseWriter, r *http.Request) {
	interval, weekday, festival := vacationInterval()
	data := map[string]interface{}{
		"weekday":  weekday,
		"value":    interval,
		"festival": festival,
	}
	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	log.Sugar().Infof("日期获取成功")
	fmt.Fprintln(w, string(jsonData))
}

func main() {
	// 定义命令行参数
	flag.StringVar(&port, "port", "5000", "listen port")
	flag.IntVar(&interval, "interval", 1, "port increase interval")
	// 解析命令行参数
	flag.Parse()

	// 初始化日志服务
	log = logger.Init()
	defer log.Sync()

	// 每天凌晨初始化日期服务
	current = calendar.Init(log)

	// 添加计划任务
	go func() {
		// 获取当前时间
		now := time.Now()
		// 计算下一个0点的时间
		next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		// 计算需要等待的时间
		duration := next.Sub(now)

		log.Sugar().Infof("距离0点还有: ", duration)
		// 等待需要等待的时间
		time.Sleep(duration)

		// 等到0点执行
		current = calendar.Init(log)

		// 创建一个每24小时执行一次的定时器
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop() // 在程序结束时停止定时器

		for range ticker.C {
			log.Info("开始执行计划任务")
			current = calendar.Init(log)
		}
	}()

	fileServer := http.FileServer(http.Dir("./html"))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/add", handleAddPort)
	http.HandleFunc("/api/reset", handleModPost)
	http.HandleFunc("/api/get", handleGetPost)
	http.HandleFunc("/api/date", handleGetDate)

	log.Info("Service listening on port 5000")
	server := &http.Server{
		Addr: ":" + port,
	}

	server.ListenAndServe()
}
