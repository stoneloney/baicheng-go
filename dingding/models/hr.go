package models

// 新增部门
type DepNewStruct struct {
	UpDepCode string `sql:"UpDepCode"`
	DepCode   string `sql:"DepCode"`
	DepName   string `sql:"DepName"`
	ToDD      int    `sql:"toDD"`
	Runtime   string `sql:"runtime"`
}

// 修改部门
type DepEditStruct struct {
	DepCode string `sql:"DepCode"`
	DepName string `sql:"DepName"`
	ToDD    int    `sql:"toDD"`
	Runtime string `sql:"runtime"`
}

// 删除部门
type DepDelStruct struct {
	DepCode string `sql:"DepCode"`
	DepName string `sql:"DepName"`
	ToDD    int    `sql:"toDD"`
	Runtime string `sql:"runtime"`
}

// 新用户
type EmpNewStruct struct {
	Id      string `sql:"ID"`
	ToDD    int    `sql:"toDD"`
	Runtime string `sql:"runtime"`
}

// 用户基本信息
type EmpInfoStruct struct {
	Id        string `sql:"ID"`
	Empname   string `sql:"empname"`
	Depcode   string `sql:"departmentcode"`
	Duty      string `sql:"duty"`
	InDate    string `sql:"inDate"`
	Movephone string `sql:"movephone"`
}

// 修改用户
type EmpEditStruct struct {
	Id      string `sql:"ID"`
	ToDD    int    `sql:"toDD"`
	Runtime string `sql:"runtime"`
}

// 用户删除
type EmpDelStruct struct {
	Id      string `sql:"ID"`
	ToDD    int    `sql:"toDD"`
	Runtime string `sql:"runtime"`
}

// 部门映射
type DepConvertStruct struct {
	DDdepCode string `sql:"DD_DepCode"`
	HRdepCode string `sql:"HR_DepCode"`
}
