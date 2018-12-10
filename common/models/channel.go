package models

type Channel struct {
	Model
	Name    string     `form:"email" binding:"required"`
	Pid     string     `form:"pid"`
	Status  string     `form:"status"`
	Weight  string     `form:"weight"`
}

// 获取列表
func ChannelList() []Channel {
	var channels []Channel
	db := GetDB()
	db.Order("weight desc").Find(&channels)
	return channels
}

