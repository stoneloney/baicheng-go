package projmodels

// 价格
type Data struct {
	ID	        int64   `form:"id" gorm:"primary_key"`
	Phone   	string  `form:"phone" binding:"required"`
    Type        string     `form:"type"`
    Duration    string     `form:"duration"`
    Director    string     `form:"director"`
    Model       string     `form:"model"`
    Effect      string     `form:"effect"`
    Dubbed      string     `form:"dubbed"`
    Price       int     `form:"price"`
    Createtime  string
}

// 定制
type Made struct {
	ID          int64   `form:"id" gorm:"primary_key"`
	Phone       string  `form:"phone" binding:"required"`
	City        string
	Type   	    string
	Duration    string
	Company     string
	Createtime  string
}

// 设置table name
func (Made) TableName() string {
	return "made"
}




