package data

// type holidays struct {
// 	festival string
// 	startDay string
// 	endDay   string
// 	lasts    int
// }

type Holidays struct {
	Festival string
	StartDay int
	EndDay   int
}

var WeekDays = map[int]string{
	1: "星期一",
	2: "星期二",
	3: "星期三",
	4: "星期四",
	5: "星期五",
	6: "星期六",
	7: "星期天",
}

func Holiday() *map[string][]Holidays {
	// 定义假期数据
	correntHoliday := make(map[string][]Holidays)
	correntHoliday["01"] = []Holidays{
		{
			Festival: "春节",
			StartDay: 21,
			EndDay:   27,
		},
	}
	correntHoliday["04"] = []Holidays{
		{
			Festival: "清明节",
			StartDay: 05,
			EndDay:   05,
		},
		{
			Festival: "劳动节",
			StartDay: 29,
			EndDay:   30,
		},
	}
	correntHoliday["05"] = []Holidays{
		{
			Festival: "劳动节",
			StartDay: 01,
			EndDay:   03,
		},
	}
	correntHoliday["06"] = []Holidays{
		{
			Festival: "端午节",
			StartDay: 22,
			EndDay:   24,
		},
	}
	correntHoliday["09"] = []Holidays{
		{
			Festival: "中秋节",
			StartDay: 29,
			EndDay:   30,
		},
	}
	correntHoliday["10"] = []Holidays{
		{
			Festival: "国庆节",
			StartDay: 01,
			EndDay:   06,
		},
	}

	return &correntHoliday
}
