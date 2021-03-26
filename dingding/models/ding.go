package models

type AccessTokenRes struct {
	ExpiresIn   int    `json:"expires_in"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ErrCode     int    `json:"errcode"`
}

type DepartmentRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Id      int64  `json:"id"`
}

// 获取用户的具体信息
type GetUser struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`       // 对返回码的文本描述内容
	UserId       string `json:"userid"`       // 用户id
	Unionid      string `json:"unionid"`      // 员工在当前开发者企业账号范围内的唯一标识，系统生成，固定值，不会改变
	Name         string `json:"name"`         // 员工名字
	Tel          string `json:"tel"`          // 分机号（仅限企业内部开发调用）
	WorkPlace    string `json:"workplace"`    // 办公地点
	remark       string `json:"remark"`       // 备注
	Mobile       int32  `json:"mobile"`       // 手机号
	Email        string `json:"email"`        // 员工的电子邮箱
	OrgEmail     string `json:"orgEmail"`     // 员工的企业邮箱，如果员工已经开通了企业邮箱，接口会返回，否则不会返回
	Active       bool   `json:"active"`       // 是否已经激活，true表示已激活，false表示未激活
	OrderInDepts string `json:"orderInDepts"` // 在对应的部门中的排序，Map结构的json字符串，key是部门的Id，value是人员在这个部门的排序值
	IsAdmin      bool   `json:"isAdmin"`      // 是否管理员
}

// 用户结果
type UserRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Userid  string `json:"userid"`
}

// 用户的基本信息
type UserInfoStruct struct {
	Errcode         int    `json:"errcode"`
	Unionid         string `json:"unionid"`
	Remark          string `json:"remark"`
	Userid          string `json:"userid"`
	IsLeaderInDepts string `json:"isLeaderInDepts"`
	IsBoss          bool   `json:"isBoss"`
	HiredDate       int64  `json:"hiredDate"`
	IsSenior        bool   `json:"isSenior"`
	Tel             string `json:"tel"`
	Department      []int  `json:"department"`
	WorkPlace       string `json:"workPlace"`
	Email           string `json:"email"`
	OrderInDepts    string `json:"orderInDepts"`
	Mobile          string `json:"mobile"`
	Errmsg          string `json:"errmsg"`
	Active          bool   `json:"active"`
	Avatar          string `json:"avatar"`
	IsAdmin         bool   `json:"isAdmin"`
	IsHide          bool   `json:"isHide"`
	Jobnumber       string `json:"jobnumber"`
	Name            string `json:"name"`
	Extattr         struct {
	} `json:"extattr"`
	StateCode string `json:"stateCode"`
	Position  string `json:"position"`
	Roles     []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		GroupName string `json:"groupName"`
	} `json:"roles"`
}

// 增加和修改用户信息
type UserOp struct {
	UserId       string `json:"userid"`       // 用户id
	Name         string `json:"name"`         // 员工名字
	Department   string `json:"department"`   // 成员所属部门id列表
	OrderInDepts string `json:"orderInDepts"` // 实际是Map的序列化字符串
	Position     string `json:"position"`     // 职位信息
	Mobile       int32  `json:"mobile"`       // 手机号
	Tel          string `json:"tel"`          // 分机号（仅限企业内部开发调用
	WorkPlace    string `json:"workplace"`    // 办公地点
	remark       string `json:"remark"`       // 备注
	Email        string `json:"email"`        // 员工的电子邮箱
	OrgEmail     string `json:"orgEmail"`     // 员工的企业邮箱，如果员工已经开通了企业邮箱，接口会返回，否则不会返回
	Jobnumber    string `json:"jobnumber"`    // 员工工号，对应显示到OA后台和客户端个人资料的工号栏目，长度为0~64个字符
	IsHide       bool   `json:"isHide"`       // 隐藏手机号后，手机号在个人资料页隐藏
	IsSenior     bool   `json:"isSenior"`     // 是否高管模式
	Extattr      string `json:"extattr"`      // 扩展属性，可以设置多种属性
	HiredDate    int32  `json:"hiredDate"`    // 入职时间，Unix时间戳
}

// 打卡结果列表(简版)
type Record struct {
	Id             int64  `json:"id"`        // 唯一标识
	GroupId        int64  `json:"groupId"`   // 考勤组ID
	PlanId         int64  `json:"planId"`    // 排班ID
	RecordId       int64  `json:"recordId"`  // 打卡记录ID
	WorkDate       int64  `json:"workDate"`  // 工作日
	UserId         string `json:"userId"`    // 用户ID
	CheckType      string `json:"checkType"` // OnDuty：上班   OffDuty：下班
	TimeResult     string `json:"timeResult"`
	LocationResult string `json:"locationResult"`
	ApproveId      string `json:"approveId"` // 关联的审批id，当该字段非空时，表示打卡记录与请假、加班等审批有关
	ProcInstId     string `json:"procInstId"`
	BaseCheckTime  int64  `json:"baseCheckTime"` // 计算迟到和早退，基准时间
	UserCheckTime  int64  `json:"userCheckTime"` // 实际打卡时间,  用户打卡时间的毫秒数
	SourceType     string `json:"sourceType"`
	UserAddress    string `json:"userAddress"` // 打开地址
}

// 打卡结果(简版)
type AttendanceListRes struct {
	ErrCode      int      `json:"errcode"`      // 返回码
	ErrMsg       string   `json:"errmsg"`       // 对返回码的文本描述内容
	HasMore      string   `json:"hasMore"`      // 分页返回参数，表示是否还有更多数据
	Recordresult []Record `json:"recordresult"` // 结果列表
}

type QueryonjobResult struct {
	NextCursor int      `json:"next_cursor"` // 下一次分页调用的offset值，当返回结果里没有nextCursor时，表示分页结束
	DataList   []string `json:"data_list"`   // 员工userid列表
}

// 企业员工结果
type QueryonjobRes struct {
	ErrCode int              `json:"errcode"` // 返回码
	Errmsg  string           `json:"errmsg"`  // 对返回码的文本描述内容
	Success bool             `json:"success"` // 调用是否成功
	Result  QueryonjobResult `json:"result"`  // 分页结果
}

type LeaveStatus struct {
	DurationUnit    string `json:"duration_unit"`
	DurationPercent int    `json:"duration_percent"`
	EndTime         int64  `json:"end_time"`
	StartTime       int64  `json:"start_time"`
	Userid          string `json:"userid"`
}

// 员工请假状态
type AttendanceLeave struct {
	Errmsg  string `json:"errmsg"`
	Errcode int    `json:"errcode"`
	Result  struct {
		HasMore     bool          `json:"has_more"`
		LeaveStatus []LeaveStatus `json:"leave_status"`
	} `json:"result"`
	Success bool `json:"success"`
}

// 注册回调
type RegCallback struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 回调列表
type CallbackList struct {
	Errcode     int      `json:"errcode"`
	Errmsg      string   `json:"errmsg"`
	CallbackTag []string `json:"call_back_tag"`
	Token       string   `json:"token"`
	AesKey      string   `json:"aes_key"`
	Url         string   `json:"url"`
}

// 事件名称
type EventType struct {
	EventType string `json:"EventType"`
}

// 签到回调
type CheckInStruct struct {
	CorpId    string `json:"CorpId"`
	EventType string `json:"EventType"`
	StaffId   string `json:"StaffId"`
	TimeStamp int64  `json:"TimeStamp"`
}

// 回调时间
type CallbackMsgStruct struct {
	ProcessInstanceId string `json:"processInstanceId"`
	CorpId            string `json:"corpId"`
	EventType         string `json:"EventType"`
	BusinessId        string `json:"businessId"`
	Title             string `json:"title"`
	Type              string `json:"type"`
	Url               string `json:"url"`
	CreateTime        int64  `json:"createTime"`
	FinishTime        int64  `json:"finishTime"`
	Result            string `json:"result"`
	ProcessCode       string `json:"processCode"`
	BizCategoryId     string `json:"bizCategoryId"`
	StaffId           string `json:"staffId"`
	Remark            string `json:"remark"`
}

// 审批列表
type ProcessInstanceListStruct struct {
	Result struct {
		List       []string `json:"list"`
		NextCursor int64    `json:"next_cursor"`
	} `json:"result"`
	Errcode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	RequestId string `json:"request_id"`
}

// 审批实例
type ProcessInstanceStruct struct {
	Errcode         int    `json:"errcode"`
	Errmsg          string `json:"errmsg"`
	ProcessInstance struct {
		Title               string `json:"title"`
		CreateTime          string `json:"create_time"`
		FinishTime          string `json:"finish_time"`
		OriginatorUserid    string `json:"originator_userid"`
		OriginatorDeptID    string `json:"originator_dept_id"`
		Status              string `json:"status"`
		CcUserids           string `json:"cc_userids"`
		FormComponentValues []struct {
			Name     string `json:"name"`
			Value    string `json:"value"`
			ExtValue string `json:"ext_value"`
		} `json:"form_component_values"`
		Result           string `json:"result"`
		BusinessID       string `json:"business_id"`
		OperationRecords []struct {
			Userid          string `json:"userid"`
			Date            string `json:"date"`
			OperationType   string `json:"operation_type"`
			OperationResult string `json:"operation_result"`
			Remark          string `json:"remark"`
		} `json:"operation_records"`
		Tasks []struct {
			Userid     string `json:"userid"`
			TaskStatus string `json:"task_status"`
			TaskResult string `json:"task_result"`
			CreateTime string `json:"create_time"`
			FinishTime string `json:"finish_time"`
			Taskid     string `json:"taskid"`
		} `json:"tasks"`
		OriginatorDeptName         string   `json:"originator_dept_name"`
		BizAction                  string   `json:"biz_action"`
		AttachedProcessInstanceIds []string `json:"attached_process_instance_ids"`
	} `json:"process_instance"`
}

// 部门详情
type DepartmentInfoStruct struct {
	Errcode               int    `json:"errcode"`
	Errmsg                string `json:"errmsg"`
	ID                    int64  `json:"id"`
	Name                  string `json:"name"`
	Order                 int    `json:"order"`
	Parentid              int64  `json:"parentid"`
	CreateDeptGroup       bool   `json:"createDeptGroup"`
	AutoAddUser           bool   `json:"autoAddUser"`
	DeptHiding            bool   `json:"deptHiding"`
	DeptPermits           string `json:"deptPermits"`
	UserPermits           string `json:"userPermits"`
	OuterDept             bool   `json:"outerDept"`
	OuterPermitDepts      string `json:"outerPermitDepts"`
	OuterPermitUsers      string `json:"outerPermitUsers"`
	OrgDeptOwner          string `json:"orgDeptOwner"`
	DeptManagerUseridList string `json:"deptManagerUseridList"`
	SourceIdentifier      string `json:"sourceIdentifier"`
}

// 手机号
type PhoneStruct struct {
	Phone uint64 `sql:"phone"`
}

// 用户userid
type UseridStruct struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Userid  string `json:"userid"`
}

type NewComponentValues struct {
	ComponentName string `json:"component_name"`
	ComponentType string `json:"component_type"`
	Props         struct {
		BizAlias string `json:"bizAlias"`
	} `json:"props"`
	Value string `json:"value"`
}
