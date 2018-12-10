package models

type Sms struct {
	Model
	Ip     		string   `form:"ip"`
	Phone  		string   `form:"phone"`
 	Type   		int      `form:"type"`
 	Number      string   `form:"number"`
 	Status      int   	 `form:"status"`
 	Createtime  string   
}
