package controllers

import (
	"dingding/method"
	"dingding/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var APPKEY, APPSECRET string

var Ncrypto *Crypto
var cacheSyncMap sync.Map

const (
	DepartmentURL      = "https://oapi.dingtalk.com/department"                       // 部门url前缀
	UserUrl            = "https://oapi.dingtalk.com/user"                             // 用户url前缀
	AttendanceUrl      = "https://oapi.dingtalk.com/attendance"                       // 考勤url前缀
	SmartworkUrl       = "https://oapi.dingtalk.com/topapi/smartwork/hrm/employee"    // 智能人事
	CallbackUrl        = "https://oapi.dingtalk.com/call_back"                        // 事件回调
	ProcessUrl         = "https://oapi.dingtalk.com/topapi/processinstance"           // 审批实例详情
	PushUrl            = "https://oapi.dingtalk.com/topapi/message/corpconversation/" // 信息推送
)

/*
func DepartmentCreate(c *gin.Context) {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{"code":1001, "msg":"get token error"})
		return
	}
	name := "lua开发组"
	parentid := 1
	sourceIdentifier := "21"
	id, err := departmentCreate(token, name, parentid, sourceIdentifier)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "id":id, "msg":"success"})
}

func DepartmentUpdate(c *gin.Context) {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{"code":1001, "msg":"get token error"})
		return
	}
	name := "python开发组"
	id := int64(114181015)
	_, err = departmentUpdate(token, id, name)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func DepartmentDelete(g *gin.Context) {
	token, err := getAccessToken()
	id := int64(114181015)
	_, err = departmentDelete(token, id)
	if err != nil {
		fmt.Println(err.Error())
	}
}
*/

// 钉钉回调函数
func Callback(c *gin.Context) {
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Logger.Warning(err.Error())
	} else {
		sign := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		type cbodyStruct struct {
			Encrypt string `json:"encrypt"`
		}
		var cbody cbodyStruct
		_ = json.Unmarshal(body, &cbody)
		data, err := Ncrypto.DecryptMsg(sign, timestamp, nonce, cbody.Encrypt)
		if err != nil {
			Logger.Error(err.Error())
			return
		}
		var eventType models.EventType
		_ = json.Unmarshal([]byte(data), &eventType)
		Logger.Info("cb_info:" + string(data))
		/*
			if eventType.EventType == "check_url" {

			}
		*/
		switch eventType.EventType {
		case "check_url":
			callbackCheck(c)
		case "check_in": // 签到
			callbackCheckIn(data)
		case "user_add_org": // 员工增加
			callbackUserAdd(data)
		case "user_modify_org": // 员工修改
			callbackUserModify(data)
		case "user_leave_org": // 员工离职
			callbackUserLeave(data)
		case "org_dept_modify": // 部门修改
			callbackDeptModify(data)
		case "org_dept_create": // 部门添加
			callbackDeptCreate(data)
		case "bpms_instance_change": // 补卡，请假
			callbacInstanceChange(data)
		default:
			Logger.Error("no eventType")
		}
	}
}

// 回调确认 (验证url注册)
func callbackCheck(c *gin.Context) {
	replymsg := "success"
	timestamp := time.Now().Unix()
	nonce := "123456"
	timeStr := strconv.FormatInt(timestamp, 10)
	encrypt, sign, err := Ncrypto.EncryptMsg(replymsg, timeStr, nonce)
	if err != nil {
		Logger.Error(err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"msg_signature": sign, "timeStamp": timeStr, "nonce": nonce, "encrypt": encrypt})
	}
}

// 签到时间回调
func callbackCheckIn(data string) {
	timelayout := "2006-01-02"
	timelayout2 := "15:04"

	var checkIn models.CheckInStruct
	_ = json.Unmarshal([]byte(data), &checkIn)

	t := time.Unix(int64(checkIn.TimeStamp/1000), 0)
	dateStr := t.Format(timelayout)
	timeStr := t.Format(timelayout2)

	staffId := checkIn.StaffId // 工号

	sqlStr := fmt.Sprintf("exec DD_insertcardToHR '%s','%s','%s'", staffId, dateStr, timeStr)
	Logger.Info(sqlStr)

	_, err := DbExec(sqlStr)
	if err != nil {
		//fmt.Println(err.Error())
		Logger.Error(fmt.Sprintf("check_in, error:%s", err.Error()))
	}
	Logger.Info("check_in success")
}

// 审批实例
func callbacInstanceChange(data string) {
	var cbMsg models.CallbackMsgStruct
	_ = json.Unmarshal([]byte(data), &cbMsg)
	title := cbMsg.Title
	pid := cbMsg.ProcessInstanceId // 实例id
	//fmt.Println(pid)
	ptype := cbMsg.Type // 审批正常结束（同意或拒绝）的type为finish，审批终止的type为terminate，
	// 这里本来可以按照 模版id(processCode)来准确匹配，但配置是个麻烦事，这里就使用title匹配关键字来确实是什么类型
	if ptype == "finish" { // type为finish时为结束，这里需要查询审批的具体信息，返给九段最后审批的结果
		// ProcessinstanceInfo(pid)
		if strings.Index(title, "加班") != -1 { // 加班
			//insertOvertimneToHr(pid)
			insertOvertimneToHrV2(pid)
		} else if strings.Index(title, "补卡") != -1 { // 补卡
			// 写入补卡的信息
			insertRecordToHr(pid)
		} else if strings.Index(title, "请假") != -1 { // 请假
			// 写入请假的具体信息
			insertLeaveToHr(pid)
		} else if strings.Index(title, "外出") != -1 { // 外出
			// 写入外出的具体信息
			insertOutsideToHr(pid)
		} else if strings.Index(title, "出差") != -1 { // 出差
			insertTripToHr(pid)
		}
	}
}

// 写入hr加班
func insertOvertimneToHrV2(pid string) {
	// 获取审批的详情
	data, err := ProcessinstanceInfo(pid)
	if err != nil {
		Logger.Error("insert overtime error:" + err.Error())
	} else {
		processInstance := data.ProcessInstance
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			componentValue := processInstance.FormComponentValues[0].Value // 加班详情
			// 表单详情2
			overtimeDesc := processInstance.FormComponentValues[1].Value // 加班原因

			var cvRes []models.NewComponentValues
			_ = json.Unmarshal([]byte(componentValue), &cvRes)
			if len(cvRes) > 0 {

				startTIme := ""    // 加班开始时间
				endTIme := ""      // 加班结束时间
				duration := ""     // 加班时长
				overtimeType := "" // 加班类型
				overtimeType = "1" // 常用加班

				// 查找数据
				for _, v := range cvRes {
					switch v.Props.BizAlias {
					case "startTime":
						{
							startTIme = v.Value
						}
					case "finishTime":
						{
							endTIme = v.Value
						}
					case "duration":
						{
							duration = v.Value
						}

					}
				}

				originatorUserid := processInstance.OriginatorUserid // 发起人

				timelayout1 := "2006-01-02"
				timelayout2 := "15:04"

				// 加班申请日期
				//createTIme := processInstance.CreateTime   // 发起时间
				//createTImeInt, _ := time.Parse("2006-01-02 15:04:05", createTIme)
				//createTimeformat := createTImeInt.Format(timelayout1)

				// 加班开始时间
				startTImeInt, _ := time.Parse("2006-01-02 15:04", startTIme)
				startTImeformat := startTImeInt.Format(timelayout2)
				startTImeformat2 := startTImeInt.Format(timelayout1)

				// 加班结束时间
				endTImeInt, _ := time.Parse("2006-01-02 15:04", endTIme)
				endTImeformat := endTImeInt.Format(timelayout2)

				userinfo, err := GetUserInfo(originatorUserid)
				if err != nil {
					Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
				} else {
					mobile := userinfo.Mobile // userid改为手机号
					// 发送加班到九段
					sqlStr := fmt.Sprintf("exec DD_InsertOvertimeToHR '%s','%s','%s','%s','%s','%s','%s'", mobile, startTImeformat2, startTImeformat, endTImeformat, duration, overtimeType, overtimeDesc)
					Logger.Info("sqlstr:" + sqlStr)
					fmt.Println(sqlStr)
					_, err = DbExec(sqlStr)
					if err != nil {
						Logger.Error("exec DD_InsertOvertimeToHR error:" + err.Error())
					}
				}
			} else {
				Logger.Error("exec DD_InsertOvertimeToHR cvRes empty")
			}
		}
	}
}

// 写入hr加班
func insertOvertimneToHr(pid string) {
	// 获取审批的详情
	data, err := ProcessinstanceInfo(pid)
	if err != nil {
		Logger.Error("insert overtime error:" + err.Error())
	} else {
		processInstance := data.ProcessInstance
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			componentValue := processInstance.FormComponentValues[0].Value // 加班详情
			// 表单详情2
			overtimeDesc := processInstance.FormComponentValues[1].Value // 加班原因

			var cvRes []interface{}
			_ = json.Unmarshal([]byte(componentValue), &cvRes)
			if len(cvRes) > 0 {
				startTIme := cvRes[0]    // 加班开始时间
				endTIme := cvRes[1]      // 加班结束时间
				duration := cvRes[2]     // 加班时长
				overtimeType := cvRes[4] // 加班类型
				if overtimeType == nil {
					overtimeType = "1" // 常用加班
				}
				originatorUserid := processInstance.OriginatorUserid // 发起人
				createTIme := processInstance.CreateTime             // 发起时间

				timelayout1 := "2006-01-02"
				timelayout2 := "15:04"

				// 加班申请日期
				//createTImeInt, _ := time.Parse("2006-01-02 15:04:05", createTIme)
				createTImeInt, _ := time.Parse("2006-01-02 15:04:05", createTIme)
				createTimeformat := createTImeInt.Format(timelayout1)

				// 加班开始时间
				startTImeInt, _ := time.Parse("2006-01-02 15:04", startTIme.(string))
				startTImeformat := startTImeInt.Format(timelayout2)

				// 加班结束时间
				endTImeInt, _ := time.Parse("2006-01-02 15:04", endTIme.(string))
				endTImeformat := endTImeInt.Format(timelayout2)

				userinfo, err := GetUserInfo(originatorUserid)
				if err != nil {
					Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
				} else {
					mobile := userinfo.Mobile // userid改为手机号
					// 发送加班到九段
					sqlStr := fmt.Sprintf("exec DD_InsertOvertimeToHR '%s','%s','%s','%s','%.1f','%s','%s'", mobile, createTimeformat, startTImeformat, endTImeformat, duration, overtimeType, overtimeDesc)
					Logger.Info("sqlstr:" + sqlStr)
					fmt.Println(sqlStr)
					_, err = DbExec(sqlStr)
					if err != nil {
						Logger.Error("exec DD_InsertOvertimeToHR error:" + err.Error())
					}
				}
			} else {
				Logger.Error("exec DD_InsertOvertimeToHR cvRes empty")
			}
		}
	}
}

// 写入hr补卡
func insertRecordToHr(pid string) {
	// 获取审批的险情
	data, err := ProcessinstanceInfo(pid)
	if err != nil {
		Logger.Error("insert overtime error:" + err.Error())
	} else {
		processInstance := data.ProcessInstance
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			recordDesc := ""                                     // processInstance.FormComponentValues[0].Value  // 补卡原因
			originatorUserid := processInstance.OriginatorUserid // 发起人
			recordTime := ""                                     // processInstance.FormComponentValues[2].Value  // 补卡的时间点

			for _, v := range processInstance.FormComponentValues {
				if v.Name == "补卡理由" {
					recordDesc = v.Value
				} else if v.Name == "repairCheckTime" {
					recordTime = v.Value
				}
			}

			timelayout1 := "2006-01-02"
			timelayout2 := "15:04"
			recordTimeInt, _ := strconv.ParseInt(recordTime, 10, 64)
			timeint := int64(recordTimeInt / 1000)
			//fmt.Println(timeint)
			tu := time.Unix(timeint, 0)

			time1 := tu.Format(timelayout1)
			time2 := tu.Format(timelayout2)

			// 补卡类型，再钉钉上没看到类型设置，这里暂时写死  因公
			recordType := 1

			// 获取用户基本信息
			userinfo, err := GetUserInfo(originatorUserid)
			if err != nil {
				Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
			} else {
				mobile := userinfo.Mobile // userid改为手机号
				// 发送补卡类型到hr
				sqlStr := fmt.Sprintf("exec DD_InsertBukaToHR '%s','%s','%s','%d','%s'", mobile, time1, time2, recordType, recordDesc)
				Logger.Info(sqlStr)
				_, err = DbExec(sqlStr)
				if err != nil {
					Logger.Error("exec DD_InsertBukaToHR error:" + err.Error())
				}
			}
		}
	}
}

// 写入hr请假
func insertLeaveToHr(pid string) {
	// 获取审批的险情
	data, err := ProcessinstanceInfo(pid)
	if err != nil {
		Logger.Error("insert overtime error:" + err.Error())
	} else {
		processInstance := data.ProcessInstance
		Logger.Info(processInstance.Status + ":" + processInstance.Result)
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			componentValue := processInstance.FormComponentValues[0].Value // 请假详情
			// 表单详情2
			leaveDesc := processInstance.FormComponentValues[1].Value // 请假原因

			var cvRes []interface{}
			_ = json.Unmarshal([]byte(componentValue), &cvRes)
			Logger.Info(cvRes[0].(string))
			if len(cvRes) > 0 {
				startTIme := cvRes[0]                                // 请假开始时间
				endTIme := cvRes[1]                                  // 请假结束时间
				duration := cvRes[2]                                 // 请假时长
				leaveType := cvRes[4]                                // 请假类型
				originatorUserid := processInstance.OriginatorUserid // 发起人
				//createTIme := processInstance.CreateTime               // 发起时间

				// 获取用户基本信息
				userinfo, err := GetUserInfo(originatorUserid)
				if err != nil {
					Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
				} else {
					mobile := userinfo.Mobile // userid改为手机号
					// 发送加班到九段
					//sqlStr := fmt.Sprintf("exec DD_InsertHolidayToHR '%s','%s','%s','%s','%.1f','%s','%s','%s'", originatorUserid, startTIme, endTIme, leaveType, duration, leaveDesc, "0", "0")
					sqlStr := fmt.Sprintf("exec DD_InsertHolidayToHR '%s','%s','%s','%s','%v','%s','%s','%s'", mobile, startTIme, endTIme, leaveType, duration, leaveDesc, "0", "0")
					Logger.Info(sqlStr)
					_, err = DbExec(sqlStr)
					if err != nil {
						Logger.Error("exec DD_InsertHolidayToHR error:" + err.Error())
					}
				}
			} else {
				Logger.Error("exec DD_InsertHolidayToHR cvRes empty")
			}
		}
	}
}

// 外出写入请假列表 (和请假的类型相同)
func insertOutsideToHr(pid string) {
	data, err := ProcessinstanceInfo(pid)
	if err != nil {
		Logger.Error("insert overtime error:" + err.Error())
	} else {
		processInstance := data.ProcessInstance
		Logger.Info(processInstance.Status + ":" + processInstance.Result)
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			componentValue := processInstance.FormComponentValues[0].Value //  外出详情
			// 表单详情2
			leaveDesc := processInstance.FormComponentValues[1].Value // 外出原因

			var cvRes []interface{}
			_ = json.Unmarshal([]byte(componentValue), &cvRes)
			Logger.Info(cvRes[0].(string))
			if len(cvRes) > 0 {
				startTIme := cvRes[0] // 外出开始时间
				endTIme := cvRes[1]   // 外出结束时间
				duration := cvRes[2]  // 外出时长
				//leaveType := cvRes[4]     // 外出类型
				leaveType := "外出"
				originatorUserid := processInstance.OriginatorUserid // 发起人
				//createTIme := processInstance.CreateTime               // 发起时间

				// 获取用户基本信息
				userinfo, err := GetUserInfo(originatorUserid)
				if err != nil {
					Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
				} else {
					mobile := userinfo.Mobile // userid改为手机号
					// 发送加班到九段
					sqlStr := fmt.Sprintf("exec DD_InsertHolidayToHR '%s','%s','%s','%s','%.1f','%s','%s','%s'", mobile, startTIme, endTIme, leaveType, duration, leaveDesc, "0", "0")
					Logger.Info(sqlStr)
					_, err = DbExec(sqlStr)
					if err != nil {
						Logger.Error("exec DD_InsertHolidayToHR error:" + err.Error())
					}
				}
			} else {
				Logger.Error("exec DD_InsertHolidayToHR cvRes empty")
			}
		}
	}
}

// 出差写入请假
func insertTripToHr(pid string) {
	data, err := ProcessinstanceInfo(pid)
	if err != nil {
		Logger.Error("insert overtime error:" + err.Error())
	} else {
		fmt.Println(data)
	}
}

// 用户增加回调
func callbackUserAdd(data string) {
	type useradd struct {
		CorpId    string   `json:"CorpId"`
		EventType string   `json:"EventType"`
		UserId    []string `json:"UserId"`
		TimeStamp string   `json:"TimeStamp"`
	}
	var cbMsg useradd
	_ = json.Unmarshal([]byte(data), &cbMsg)
	for _, v := range cbMsg.UserId {
		userId := v
		// 获取用户的具体信息
		userInfo, err := GetUserInfo(userId)
		//fmt.Println(userInfo)
		if err != nil {
			//fmt.Println("getuserinfo error:%s", err.Error())
			Logger.Error("getuserinfo error:" + err.Error())
		} else {
			userId := userInfo.Userid         // 员工id
			name := userInfo.Name             // 姓名
			department := userInfo.Department // 部门id
			position := userInfo.Position     // 职位信息
			hiredDate := userInfo.HiredDate   // 入职时间
			mobile := userInfo.Mobile         // 手机号

			// 写入下工号
			jobnumber := userInfo.Jobnumber // 工号

			// 9段需要强制用户jobnumber 和 userid相同
			if userId != jobnumber {
				// 不同 更新jobnumber
				updateCode, err := updateJobNumber(userId, userId)
				if err != nil {
					Logger.Error(fmt.Sprintf("updatecode:%s, error:%s", updateCode, err.Error()))
				}
			}

			// 根据钉钉部门id找出9段部门
			departmentId := ""
			if len(department) > 0 { // 只取第一个部门
				ddid := department[0]
				ddidStr := strconv.FormatInt(int64(ddid), 10)
				depcode, err := getDepidByDdid(ddidStr)
				if err != nil {
					//fmt.Println("getDepidByDdid error:%s", err.Error())
					Logger.Error("getDepidByDdid error:%s" + err.Error())
				} else {
					departmentId = depcode
				}
			}

			if hiredDate == 0 {
				hiredDate = int64(time.Now().Unix() * 1000)
			}

			// 日期格式化
			timelayout := "2006-01-02"
			t := time.Unix(int64(hiredDate/1000), 0)
			hiredDateStr := t.Format(timelayout)

			var resCode int
			sqlStr := fmt.Sprintf("exec DD_InsertEmpToHR '%s','%s','%s','%s','%s','%s', %d", userId, name, departmentId, position, hiredDateStr, mobile, resCode)
			//fmt.Sprintf(sqlStr)
			Logger.Info("sqlstr:" + sqlStr)
			_, err := DbExec(sqlStr)
			if err != nil {
				//fmt.Println(err.Error())
				Logger.Error(fmt.Sprintf("insert emptohr, error:%s", err.Error()))
			} else {
				//fmt.Println("insert emptohr success")
				Logger.Info("insert emptohr success")
			}
		}
	}
}

// 员工修改
func callbackUserModify(data string) {
	type userModify struct {
		CorpId    string   `json:"CorpId"`
		EventType string   `json:"EventType"`
		UserId    []string `json:"UserId"`
		TimeStamp string   `json:"TimeStamp"`
	}
	var cbMsg userModify
	_ = json.Unmarshal([]byte(data), &cbMsg)
	for _, v := range cbMsg.UserId {
		userId := v
		// 获取用户的具体信息
		userInfo, err := GetUserInfo(userId)
		//fmt.Println(userInfo)
		if err != nil {
			//fmt.Println("getuserinfo error:%s", err.Error())
			Logger.Error("getuserinfo error:" + err.Error())
		} else {
			userId := userInfo.Userid         // 员工id
			name := userInfo.Name             // 姓名
			department := userInfo.Department // 部门id
			position := userInfo.Position     // 职位信息
			hiredDate := userInfo.HiredDate   // 入职时间
			mobile := userInfo.Mobile         // 手机号

			// 根据钉钉部门id找出9段部门
			departmentId := ""
			if len(department) > 0 { // 只取第一个部门
				ddid := department[0]
				ddidStr := strconv.FormatInt(int64(ddid), 10)
				depcode, err := getDepidByDdid(ddidStr)
				if err != nil {
					//fmt.Println("getDepidByDdid error:%s", err.Error())
					Logger.Error("getDepidByDdid error:%s" + err.Error())
				} else {
					departmentId = depcode
				}
			}

			if hiredDate == 0 {
				hiredDate = int64(time.Now().Unix() * 1000)
			}

			// 日期格式化
			timelayout := "2006-01-02"
			t := time.Unix(int64(hiredDate/1000), 0)
			hiredDateStr := t.Format(timelayout)

			var resCode int
			sqlStr := fmt.Sprintf("exec DD_UpdateEmpToHR '%s','%s','%s','%s','%s','%s', %d", userId, name, departmentId, position, hiredDateStr, mobile, resCode)
			fmt.Sprintf(sqlStr)
			Logger.Info("sqlstr:" + sqlStr)
			_, err := DbExec(sqlStr)
			if err != nil {
				//fmt.Println(err.Error())
				Logger.Error(fmt.Sprintf("modify emptohr, error:%s", err.Error()))
			} else {
				//fmt.Println("modify emptohr success")
				Logger.Info("modify emptohr success")
			}
		}
	}
}

// 员工离职 (同用户删除)
func callbackUserLeave(data string) {

}

// 部门修改
func callbackDeptModify(data string) {
	type deptModify struct {
		CorpId    string  `json:"CorpId"`
		EventType string  `json:"EventType"`
		DeptId    []int64 `json:"DeptId"`
		TimeStamp string  `json:"TimeStamp"`
	}
	var cbMsg deptModify
	_ = json.Unmarshal([]byte(data), &cbMsg)
	for _, v := range cbMsg.DeptId {
		depId := v
		Logger.Info(fmt.Sprintf("callbackDeptModify:%d", v))
		// 获取用户的具体信息
		depInfo, err := DepartmentInfo(depId)
		//fmt.Println(depInfo)
		if err != nil {
			//fmt.Println("depinfo error:%s", err.Error())
			Logger.Error("depinfo error:" + err.Error())
		} else {
			name := depInfo.Name
			// 找出在9段的部门id
			depIdstr := strconv.FormatInt(depId, 10)
			depcode, err := getDepidByDdid(depIdstr)
			if err != nil {
				//fmt.Println("getDepidByDdid:error:" + err.Error())
				Logger.Error("getDepidByDdid:error:" + err.Error())
			} else {
				sqlStr := fmt.Sprintf("exec DD_UpdateDepToHR '%s','%s'", depcode, name)
				Logger.Info(sqlStr)
				_, err = DbExec(sqlStr)
				if err != nil {
					//fmt.Println("exec DD_UpdateDepToHR error:" + err.Error())
					Logger.Error("exec DD_UpdateDepToHR error:" + err.Error())
				}
			}
		}
	}
}

// 部门添加
func callbackDeptCreate(data string) {
	type deptCreate struct {
		CorpId    string  `json:"CorpId"`
		EventType string  `json:"EventType"`
		DeptId    []int64 `json:"DeptId"`
		TimeStamp string  `json:"TimeStamp"`
	}
	var cbMsg deptCreate
	_ = json.Unmarshal([]byte(data), &cbMsg)
	for _, v := range cbMsg.DeptId {
		// 获取部门详情
		Logger.Info(fmt.Sprintf("depid:%d", v))
		data, err := DepartmentInfo(v)
		if err != nil {
			Logger.Error("get departinfo error:" + err.Error())
		} else {
			// 查看是否顶级部门
			parentid := data.Parentid
			ddid := data.ID
			depName := data.Name
			depid := ""
			if len(data.SourceIdentifier) > 0 {
				depid = data.SourceIdentifier
			} else {
				depid, err = MakeHrDepId(parentid)
				if err != nil {
					Logger.Error("make hrdepid error:" + err.Error())
					return
				}
			}
			Logger.Info("makehrdepid:" + depid)
			/*
				upDepCode := ""
				if parentid != 1 {
					// 写入hr新增的部门
					parentidstr := strconv.FormatInt(parentid,10)
					upDepCode, err = getDepidByDdid(parentidstr)
					if err != nil {
						Logger.Error("getDepidByDdid error:"+err.Error())
					}
				}
				// 写入部门新增
				runtimne := method.NowStr()
				sqlStr := fmt.Sprintf("INSERT INTO [DD_DepNew](UpDepCode, DepCode, DepName, toDD, Runtime)VALUES('%s','%s','%s','%s','%s')", upDepCode, depid, depName, "1", runtimne)
				_, err := DbExec(sqlStr)
				if err != nil {
					Logger.Error("INSERT INTO [DD_DepNew] error:"+err.Error())
				}
			*/
			// 写入部门映射
			sqlStr := fmt.Sprintf("INSERT INTO [DD_DepConvert](DD_DepCode, HR_DepCode) VALUES('%d', '%s')", ddid, depid)
			//fmt.Println("insert into depcovert sql:" + sqlStr)
			Logger.Info("insert into depcovert sql:" + sqlStr)
			_, err = DbExec(sqlStr)
			if err == nil {
				// 写入hr
				sqlStr := fmt.Sprintf("exec DD_InsertDepToHR '%s','%s'", depid, depName)
				Logger.Info(sqlStr)
				_, err := DbExec(sqlStr)
				if err != nil {
					Logger.Error("exec DD_InsertDepToHR error:" + err.Error())
				}
			} else {
				Logger.Error(fmt.Sprintf("insert into depconver error:%s, ddid:%d, depid:%s", err.Error(), ddid, depid))
			}
		}
	}
}

/**
设置appkey
@params appkey
*/
func SetAppKey(appkey string) {
	APPKEY = appkey
}

/**
设置secret
@params secret
*/
func SetSecret(secret string) {
	APPSECRET = secret
}

/**
  钉钉增加部门
  @params name  部门名称
  @params parentid  父部门id，根部门id为1
  @params sourceIdentifier  部门标识字段，开发者可用该字段来唯一标识一个部门，并与钉钉外部通讯录里的部门做映射
*/
func DepartmentCreate(name string, parentid string, sourceIdentifier string) (int64, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1001, err
	}
	if len(name) == 0 {
		return 0, errors.New("name is empty")
	}
	if !method.IsNumeric(parentid) {
		return 0, errors.New("parentid is not number")
	}
	if !method.IsNumeric(sourceIdentifier) {
		return 0, errors.New("sourceIdentifier is not number")
	}
	url := fmt.Sprintf("%s/create?access_token=%s", DepartmentURL, token)
	//fmt.Println(url)
	// 构造数据
	data := make(map[string]interface{})
	data["name"] = name
	data["parentid"] = parentid
	data["sourceIdentifier"] = sourceIdentifier
	dataStr, _ := json.Marshal(data)
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		return 0, err
	}
	//fmt.Println(string(body))
	var department models.DepartmentRes
	_ = json.Unmarshal(body, &department)
	if department.ErrCode != 0 {
		return 0, errors.New(department.ErrMsg)
	}
	return department.Id, nil
}

/**
  钉钉更新部门
  @params token  访问token
  @params id     部门id
  @params name   名称
*/
func DepartmentUpdate(id int64, name string) (int64, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1001, err
	}
	if len(token) == 0 {
		return 0, errors.New("token is empty")
	}
	if !method.IsNumeric(id) {
		return 0, errors.New("id is not number")
	}
	if len(name) == 0 {
		return 0, errors.New("name is empty")
	}
	url := fmt.Sprintf("%s/update?access_token=%s", DepartmentURL, token)
	//fmt.Println(url)
	// 构造数据
	data := make(map[string]interface{})
	data["id"] = id
	data["name"] = name
	dataStr, _ := json.Marshal(data)
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		//fmt.Println(err.Error())
		return 0, err
	}
	//fmt.Println(string(body))
	var department models.DepartmentRes
	_ = json.Unmarshal(body, &department)
	if department.ErrCode != 0 {
		return 0, errors.New(department.ErrMsg)
	}
	return department.Id, nil
}

/**
  钉钉删除部门
  @params token  访问token
  @params id     部门id
*/
func DepartmentDelete(id int64) (bool, error) {
	token, err := getAccessToken()
	if err != nil {
		return false, err
	}
	if len(token) == 0 {
		return false, errors.New("token is empty")
	}
	if !method.IsNumeric(id) {
		return false, errors.New("id is not number")
	}
	url := fmt.Sprintf("%s/delete?access_token=%s&id=%d", DepartmentURL, token, id)
	//fmt.Println(url)
	body, err := method.HttpGet(url)
	if err != nil {
		return false, err
	}
	var department models.DepartmentRes
	_ = json.Unmarshal(body, &department)
	if department.ErrCode != 0 {
		return false, errors.New(department.ErrMsg)
	}
	return true, nil
}

/**
  钉钉部门详情
  @params token  访问token
  @params id     部门id
*/
func DepartmentInfo(id int64) (models.DepartmentInfoStruct, error) {
	var department models.DepartmentInfoStruct
	token, err := getAccessToken()
	if err != nil {
		return department, err
	}
	if !method.IsNumeric(id) {
		return department, errors.New("id is not number")
	}
	url := fmt.Sprintf("%s/get?access_token=%s&id=%d", DepartmentURL, token, id)
	//fmt.Println(url)
	body, err := method.HttpGet(url)
	if err != nil {
		return department, err
	}
	_ = json.Unmarshal(body, &department)
	if department.Errcode != 0 {
		return department, errors.New(department.Errmsg)
	}
	return department, nil
}

/**
  钉钉用户添加，修改
*/
func UserOp(utype, userid, name string, department []int64, position, mobile, tel, workPlace, remark, email, orgEmail, jobnumber string, isHide, isSenior bool, extattr string, hiredDate int64) (int, error) {
	if utype != "create" && utype != "update" {
		return -1001, errors.New("type is error")
	}
	if len(name) == 0 {
		return -1002, errors.New("name is empty")
	}
	token, err := getAccessToken()
	if err != nil {
		return -1003, errors.New("token is empty")
	}
	url := ""
	data := make(map[string]interface{})
	if utype == "create" {
		data["hiredDate"] = hiredDate // 入职时间
		data["mobile"] = mobile       // 激活的用户，不能修改手机号
		url = fmt.Sprintf("%s/create?access_token=%s", UserUrl, token)
	} else {
		url = fmt.Sprintf("%s/update?access_token=%s", UserUrl, token)
	}
	//fmt.Println(url)
	//fmt.Println(department)
	// 构造数据
	data["userid"] = userid
	data["name"] = name
	data["department"] = department
	data["position"] = position
	data["tel"] = tel
	data["workPlace"] = workPlace
	data["remark"] = remark
	data["email"] = email
	data["orgEmail"] = orgEmail
	data["jobnumber"] = jobnumber
	data["isHide"] = isHide
	data["isSenior"] = isSenior
	data["extattr"] = extattr
	// 发送请求
	dataStr, _ := json.Marshal(data)
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		//fmt.Println(err.Error())
		return -1004, err
	}
	//fmt.Println(string(body))
	var res models.UserRes
	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return -1005, errors.New(res.ErrMsg)
	}
	return res.ErrCode, nil
}

/**
 * 更新jobnumber
 */
func updateJobNumber(userid, jobnumber string) (int, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1003, errors.New("token is empty")
	}
	data := make(map[string]string)
	url := fmt.Sprintf("%s/update?access_token=%s", UserUrl, token)
	data["jobnumber"] = jobnumber
	data["userid"] = userid
	// 发送请求
	dataStr, _ := json.Marshal(data)
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1004, err
	}
	//fmt.Println(string(body))
	var res models.UserRes
	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return -1005, errors.New(res.ErrMsg)
	}
	return res.ErrCode, nil
}

/*
  钉钉用户删除
  @params access_token  接口凭证
  @params userid  用户id
*/
func UserDel(userid string) (int, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1001, err
	}
	url := fmt.Sprintf("%s/delete?access_token=%s&userid=%s", UserUrl, token, userid)
	//fmt.Println(url)
	body, err := method.HttpGet(url)
	if err != nil {
		return -1002, err
	}
	var res models.UserRes
	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return -1003, err
	}
	return res.ErrCode, nil
}

/**
  钉钉根据用户手机号获取userId
*/
func getDDUserIdByPhone(phone uint64) (string, error) {
	token, err := getAccessToken()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/get_by_mobile?access_token=%s&mobile=%d", UserUrl, token, phone)
	//fmt.Println(url)

	body, err := method.HttpGet(url)
	if err != nil {
		return "", err
	}

	var res models.UseridStruct

	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return "", errors.New(res.ErrMsg)
	}

	return res.Userid, nil

}

/*
  获取用户基本资料
*/
func GetUserInfo(userid string) (models.UserInfoStruct, error) {
	var userInfo models.UserInfoStruct
	token, err := getAccessToken()
	if err != nil {
		return userInfo, err
	}
	url := fmt.Sprintf("%s/get?access_token=%s&userid=%s", UserUrl, token, userid)
	//fmt.Println(url)
	body, err := method.HttpGet(url)
	if err != nil {
		return userInfo, err
	}
	_ = json.Unmarshal(body, &userInfo)
	if userInfo.Errcode != 0 {
		return userInfo, errors.New(userInfo.Errmsg)
	}
	return userInfo, nil
}

// 注册回调
func RegCallback(c *gin.Context) {
	tags := []string{
		"bpms_task_change",     //  审批任务开始，结束，转交
		"bpms_instance_change", // 审批实例开始，结束
		"user_add_org",         // 通讯录用户增加
		"user_modify_org",      // 通讯录用户更改
		"user_leave_org",       // 用户离职
		"org_dept_create",      // 通讯录企业部门创建
		"org_dept_modify",      // 通讯录企业部门修改
		"org_dept_remove",      // 通讯录企业部门删除
		"check_in",             // 用户签到
	}
	//callbackUrl := controllers.Cfg.Dingding.Callback
	_, err := RegisterCallback(tags)
	if err != nil { // err != nil
		//fmt.Println("reg:", err.Error())
		c.JSON(http.StatusOK, gin.H{"code": -4001, "msg": err.Error()})
	} else {
		// 返回回调成功
		token := "123456"
		/*
			aesKey := "1234567890123456789012345678901234567890123"
			corpid := "ding02634a50287cf0f835c2f4657eb6378f"
		*/
		aesKey := Cfg.Dingding.Aeskey
		corpid := Cfg.Dingding.Corpid

		cropty := NewCrypto(token, aesKey, corpid)

		replymsg := "success"

		timestamp := time.Now().Unix()
		timeStr := strconv.FormatInt(timestamp, 10)

		nonce := "123456"
		encrypt, sign, err := cropty.EncryptMsg(replymsg, timeStr, nonce)
		if err != nil {
			//fmt.Println(err.Error())
			panic(err.Error())
		}

		type respStruct struct {
			Sign      string `json:"msg_signature"`
			TimeStamp string `json:"timeStamp"`
			Nonce     string `json:"nonce"`
			Encrypt   string `json:"encrypt"`
		}

		var resp respStruct
		resp.Sign = sign
		resp.TimeStamp = timeStr
		resp.Nonce = nonce
		resp.Encrypt = encrypt

		//data, _ := json.Marshal(resp)

		//fmt.Println(string(data))
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	}
}

/**
  钉钉获取用户打卡结果 (只返回上午和下午的打卡结果)
  @params userid    用户列表
  @params stardate  考勤开始日期
  @params endate    考勤结束日期
*/
func AttendanceList(userids []string, startdate string, enddate string) (int, []models.Record, error) {
	var records []models.Record
	token, err := getAccessToken()
	if err != nil {
		return -1001, records, err
	}
	url := fmt.Sprintf("%s/list?access_token=%s", AttendanceUrl, token)
	fmt.Println(url)

	// 构造数据
	data := make(map[string]interface{})
	data["workDateFrom"] = startdate
	data["workDateTo"] = enddate
	data["userIdList"] = userids
	data["offset"] = 0
	data["limit"] = 50

	dataStr, _ := json.Marshal(data)
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1004, records, err
	}
	fmt.Println(string(body))
	var res models.AttendanceListRes
	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return -1005, records, err
	}
	return res.ErrCode, res.Recordresult, nil
}

/**
  钉钉获取用户打卡结果 (包含多次打卡记录)
  @params userid    用户列表
  @params stardate  考勤开始日期
  @params endate    考勤结束日期
*/
func AttendanceRecordList(userids []string, startdate string, enddate string) (int, []models.Record, error) {
	var records []models.Record
	token, err := getAccessToken()
	if err != nil {
		return -1001, records, err
	}
	url := fmt.Sprintf("%s/listRecord?access_token=%s", AttendanceUrl, token)
	//fmt.Println(url)

	// 构造数据
	data := make(map[string]interface{})
	data["checkDateFrom"] = startdate
	data["checkDateTo"] = enddate
	data["userIds"] = userids
	data["isI18n"] = false

	dataStr, _ := json.Marshal(data)

	//fmt.Println(string(dataStr))

	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1004, records, err
	}
	//fmt.Println(string(body))
	var res models.AttendanceListRes
	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return -1005, records, err
	}
	return res.ErrCode, res.Recordresult, nil
}

/**
  @查询请假状态
  @params  userid_list  用户id列表，最多100
  @params  start_time  Unix时间戳，最多180天
  @params  end_time  Unix时间戳
  @params  offset  分页偏移
  @params  size    分页大小，最大20
*/
func AttendanceLeaveStatus(userIds []string, startTime, endTime int64, offset, size int) (int, []models.LeaveStatus, bool, error) {
	var list []models.LeaveStatus
	token, err := getAccessToken()
	if err != nil {
		return -1001, list, false, err
	}
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/getleavestatus?access_token=%s", token)
	fmt.Println(url)

	// userid以逗号(,)分割
	var users string
	users = strings.Join(userIds, ",")
	fmt.Println(users)
	// 构造数据
	data := make(map[string]interface{})
	data["userid_list"] = users
	data["start_time"] = startTime
	data["end_time"] = endTime
	data["offset"] = offset
	data["size"] = size
	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1002, list, false, err
	}
	fmt.Println(string(body))
	var res models.AttendanceLeave
	_ = json.Unmarshal(body, &res)
	if res.Errcode != 0 {
		return -1003, list, false, err
	} else {
		return 0, res.Result.LeaveStatus, res.Result.HasMore, nil
	}
}

/**
获取在职员工列表
@params status_list (2,3,5,-1) 在职员工子状态筛选，其他状态无效。2，试用期；3，正式；5，待离职；-1，无状态
@params offset 分页游标，从0开始
@parms  size   分页大小，最大20
*/
func Queryonjob(offset, size int) (int, bool, int, []string, error) {
	var list []string
	token, err := getAccessToken()
	if err != nil {
		return -1001, false, 0, list, err
	}
	url := fmt.Sprintf("%s/queryonjob?access_token=%s", SmartworkUrl, token)
	fmt.Println(url)

	// 构造数据
	data := make(map[string]interface{})
	data["status_list"] = "2,3,5,-1"
	data["offset"] = offset
	data["size"] = size

	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1002, false, 0, list, err
	}
	fmt.Println(string(body))
	var res models.QueryonjobRes
	_ = json.Unmarshal(body, &res)
	if res.ErrCode != 0 {
		return -1003, false, 0, list, err
	}
	fmt.Println(res)
	return res.ErrCode, res.Success, res.Result.NextCursor, res.Result.DataList, errors.New(res.Errmsg)
}

/**
  推送text类型信息
*/
func PushTextMsg(accessToken, userId string, note []string) (int, error) {
	if len(userId) == 0 {
		return -1001, errors.New("userId empty")
	}
	if len(note) == 0 {
		return -1002, errors.New("note empty")
	}
	if len(accessToken) == 0 {
		return -1003, errors.New("accessToken empty")
	}
	/*
		token, err := getAccessToken()
		if err != nil {
			return -1001, err
		}
	*/
	var apiUrl = PushUrl + "/asyncsend_v2?access_token=" + accessToken

	// 构造数据
	data := make(map[string]interface{})
	data["agent_id"] = Cfg.Dingding.Agentid
	data["userid_list"] = userId

	type Msg struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}

	var msg Msg
	msg.MsgType = "text"
	msg.Text.Content = strings.Join(note, "\n")

	data["msg"] = msg
	dataStr, _ := json.Marshal(data)

	res, err := method.HttpPostJson(apiUrl, dataStr)
	if err != nil {
		return -3001, err
	}
	fmt.Println(string(res))
	return 0, nil
}

/**
推送oa类型信息
*/
func PushOaMsg(userId string, note []string) (int, error) {
	if len(userId) == 0 {
		return -1001, errors.New("userId empty")
	}
	if len(note) == 0 {
		return -1002, errors.New("note empty")
	}

	token, err := getAccessToken()
	if err != nil {
		return -1001, err
	}

	var apiUrl = PushUrl + "/asyncsend_v2?access_token=" + token

	// 构造数据
	data := make(map[string]interface{})
	data["agent_id"] = Cfg.Dingding.Agentid
	// test
	data["userid_list"] = userId //"171802165236043309"

	type Form struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	type Msg struct {
		MsgType string `json:"msgtype"`
		Oa      struct {
			MessageUrl string `json:"message_url"`
			Head       struct {
				BgColor string `json:"bgcolor"`
				Text    string `json:"text"`
			} `json:"head"`
			Body struct {
				Title string `json:"title"`
				Form  []Form `json:"form"`
			} `json:"body"`
		} `json:"oa"`
	}

	var msg Msg
	msg.MsgType = "oa"
	msg.Oa.Head.Text = "系统通知"
	msg.Oa.Body.Title = "工资明细"

	/*
		var form1 Form
		form1.Key = "正常工资: "
		form1.Value = "3000"

		var form2 Form
		form2.Key = "奖金: "
		form2.Value = "1000"
	*/
	for _, n := range note {
		var tmp Form
		tmp.Key = n
		tmp.Value = ""
		msg.Oa.Body.Form = append(msg.Oa.Body.Form, tmp)
	}

	data["msg"] = msg

	dataStr, _ := json.Marshal(data)

	res, err := method.HttpPostJson(apiUrl, dataStr)
	if err != nil {
		return -3001, err
	}
	fmt.Println(string(res))
	return 0, nil
}

/**
注册业务事件回调
@params call_back_tag 需要监听的事件类型
@params token  加解密需要用到的token，ISV(服务提供商)推荐使用注册套件时填写的token，普通企业可以随机填写
@params aes_key  数据加密密钥。用于回调数据的加密，长度固定为43个字符，从a-z, A-Z, 0-9共62个字符中选取,您可以随机生成，ISV(服务提供商)推荐使用注册套件时填写的EncodingAESKey
@params url  接收事件回调的url，必须是公网可以访问的url地址
*/
func RegisterCallback(tags []string) (int, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1001, err
	}
	url := fmt.Sprintf("%s/register_call_back?access_token=%s", CallbackUrl, token)
	fmt.Println(url)

	// 构造数据
	data := make(map[string]interface{})
	data["call_back_tag"] = tags
	data["token"] = Cfg.Dingding.Token
	data["aes_key"] = Cfg.Dingding.Aeskey
	data["url"] = Cfg.Dingding.Callback

	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1002, err
	}
	fmt.Println(string(body))
	var res models.RegCallback
	_ = json.Unmarshal(body, &res)
	if res.Errcode != 0 {
		return -1003, errors.New(res.Errmsg)
	}
	return 0, nil
}

func UpdateCallbackBackup(c *gin.Context) {
	tags := []string{
		"user_add_org",         // 通讯录用户增加
		"user_modify_org",      // 通讯录用户更改
		"user_leave_org",       // 用户离职
		"org_dept_create",      // 通讯录企业部门创建
		"org_dept_modify",      // 通讯录企业部门修改
		"org_dept_remove",      // 通讯录企业部门删除
		"check_in",             // 用户签到
		"bpms_task_change",     //  审批任务开始，结束，转交
		"bpms_instance_change", // 审批实例开始，结束
	}
	_, err := UpdateCallback(tags)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("update success")
	}
}

/**
  更新事件回调
  @params call_back_tag 需要监听的事件类型
  @params token  加解密需要用到的token，ISV(服务提供商)推荐使用注册套件时填写的token，普通企业可以随机填写
  @params aes_key  数据加密密钥。用于回调数据的加密，长度固定为43个字符，从a-z, A-Z, 0-9共62个字符中选取,您可以随机生成，ISV(服务提供商)推荐使用注册套件时填写的EncodingAESKey
  @params url  接收事件回调的url，必须是公网可以访问的url地址
*/
func UpdateCallback(tags []string) (int, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1001, err
	}
	url := fmt.Sprintf("%s/update_call_back?access_token=%s", CallbackUrl, token)
	fmt.Println(url)

	data := make(map[string]interface{})
	data["call_back_tag"] = tags
	data["token"] = Cfg.Dingding.Token
	data["aes_key"] = Cfg.Dingding.Aeskey
	data["url"] = Cfg.Dingding.Callback

	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return -1002, err
	}
	fmt.Println(string(body))
	var res models.RegCallback
	_ = json.Unmarshal(body, &res)
	if res.Errcode != 0 {
		return -1003, errors.New(res.Errmsg)
	}
	return 0, nil
}

/**
 * 同步列表
 */
func SyncProcessinstance(processCode string, startTime uint64, cursor int64) {
	endTime := uint64(time.Now().Unix())
	processinstanceStartKey := "processinstance_starttime"
	cacheStartTime, ok := cacheSyncMap.Load(processinstanceStartKey)
	if !ok || cacheStartTime == nil {  // 第一次使用配置的时间
		endTime = startTime + 3600*1000  // 第一次查询时间间隔
	} else {
		startTime = cacheStartTime.(uint64)
	}

	res, err := ProcessinstanceList(processCode, startTime, endTime, 10, cursor)
	if err != nil {
		fmt.Println(err.Error())
		Logger.Error(fmt.Sprintf("ProcessinstanceList starttime:%d, endTime:%d, cursor:%d, error:%s", startTime, endTime, cursor, err.Error()))
	} else {
		cursor = res.Result.NextCursor
		for _, v := range res.Result.List {
			insertOvertimneToHrV2(v)
		}

		if cursor > 0 {
			SyncProcessinstance(processCode, startTime, cursor)
		} else {
			cacheSyncMap.Store(processinstanceStartKey, endTime)
		}
	}
}

/**
 * 获取审批列表
 */
func ProcessinstanceList(processCode string, startTime uint64, endTime uint64, size, cursor int64) (models.ProcessInstanceListStruct, error) {
	var respData models.ProcessInstanceListStruct
	token, err := getAccessToken()
	if err != nil {
		return respData, err
	}
	url := fmt.Sprintf("%s/listids?access_token=%s", ProcessUrl, token)
	fmt.Println(url)

	data := make(map[string]interface{})
	data["process_code"] = processCode
	data["start_time"] = startTime
	data["end_time"] = endTime
	data["size"] = size
	data["cursor"] = cursor

	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return respData, err
	}
	fmt.Println(string(body))
	_ = json.Unmarshal(body, &respData)
	if respData.Errcode != 0 {
		return respData, errors.New(respData.ErrMsg)
	}
	return respData, nil
}

/**
 * 获取实例审批的详情
 * 返回：
 * {"errcode":0,"process_instance":{"attached_process_instance_ids":[],"biz_action":"NONE","business_id":"201906302330000062550","create_time":"2019-06-30 23:30:07","finish_time":"2019-06-30 23:30:23","form_component_values":[{"component_type":"DDOvertimeField","ext_value":"{\"compressedValue\":\"1f8b080000000000000095504b4ec33014bccb5b5b95d38650bc43ad5091a8904aba422c9e12975ab876e4e780aa281760c315b800c7425c033bf413513678f766e679e64d0305eaa2d6e8e5dc961244c2a0941e95be51e441dc3750d60ebdb2e6da4c710b820f9221eb81335bbbb036e0bc8fce95a9bd241019ffc5dcc9c29a12c4288bf84aa2af5d14368055a5b757ce6e72b58949ceb2643c4c79f7d80f9bdb3e779eedb8422351f8383a50175a9af2b27096a8bb284cc7bd8b84a7bb3df2e8fc5eb79fff706f1f18285ac8d8c80a354906cfd251700301e1746819ac912631c641f162ddd334f47a724924a2c9a1a374dc399c163dfa57d147786954480ab3dbe502185435ad737c0cc0e7ebfbd7db07b4dfe521d200f6010000\",\"unit\":\"HOUR\",\"extension\":\"{}\",\"featureMap\":{\"overtimeUrl\":\"https:\/\/attend.dingtalk.com\/attend\/index.html?corpId=ding02634a50287cf0f835c2f4657eb6378f&showmenu=false&dd_share=false&overtimeId=121185344#admin\/overtimeRuleDetail\",\"remark\":\" 加班时长以审批单为准；\",\"overtimeSettingId\":\"121185344\"},\"_from\":\"2019-06-30 00:00\",\"pushTag\":\"加班\",\"detailList\":[{\"classInfo\":{\"hasClass\":false,\"sections\":[{\"endAcross\":1,\"startTime\":1561824000000,\"endTime\":1561910400000,\"startAcross\":0}]},\"workDate\":1561824000000,\"isRest\":false,\"workTimeMinutes\":480,\"approveInfo\":{\"fromAcross\":0,\"toAcross\":0,\"fromTime\":1561824000000,\"durationInDay\":0.12,\"toTime\":1561827600000,\"durationInHour\":1}}],\"durationInDay\":0.13,\"_to\":\"2019-06-30 01:00\",\"isModifiable\":true,\"durationInHour\":1}","id":"DDOvertimeField-J2BX4G42","name":"[\"开始时间\",\"结束时间\"]","value":"[\"2019-06-30 00:00\",\"2019-06-30 01:00\",1,\"hour\",null,\"加班类型\"]"},{"component_type":"TextField","id":"加班原因","name":"加班原因","value":"null"}],"operation_records":[{"date":"2019-06-30 23:30:06","operation_result":"NONE","operation_type":"START_PROCESS_INSTANCE","userid":"171802165236043309"},{"date":"2019-06-30 23:30:22","operation_result":"AGREE","operation_type":"EXECUTE_TASK_NORMAL","remark":"","userid":"171802165236043309"},{"date":"2019-06-30 23:30:22","operation_result":"NONE","operation_type":"NONE","remark":"","userid":"171802165236043309"}],"originator_dept_id":"-1","originator_dept_name":"深圳市九段科技有限公司","originator_userid":"171802165236043309","result":"agree","status":"COMPLETED","tasks":[{"create_time":"2019-06-30 23:30:07","finish_time":"2019-06-30 23:30:23","task_result":"AGREE","task_status":"COMPLETED","taskid":"61625821309","userid":"171802165236043309"}],"title":"贾鹏飞提交的加班"},"request_id":"f6265sed2oev"}
 */
func ProcessinstanceInfo(pid string) (models.ProcessInstanceStruct, error) {
	var respData models.ProcessInstanceStruct
	token, err := getAccessToken()
	if err != nil {
		return respData, err
	}
	url := fmt.Sprintf("%s/get?access_token=%s", ProcessUrl, token)
	fmt.Println(url)

	data := make(map[string]string)
	data["process_instance_id"] = pid

	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return respData, err
	}
	fmt.Println(string(body))
	_ = json.Unmarshal(body, &respData)
	if respData.Errcode != 0 {
		return respData, errors.New(respData.Errmsg)
	}
	return respData, nil
}

/**
  查询事件回调接口
*/
func GetCallbackList() (int, string, error) {
	token, err := getAccessToken()
	if err != nil {
		return -1001, "", err
	}
	url := fmt.Sprintf("%s/get_call_back?access_token=%s", CallbackUrl, token)
	fmt.Println(url)

	body, err := method.HttpGet(url)
	if err != nil {
		fmt.Println(err.Error())
		return -1002, "", err
	}
	fmt.Println(string(body))
	return 0, string(body), nil
}

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "test", "code": 0})
}

// 获取access token
func getAccessToken() (string, error) {
	var accessToken models.AccessTokenRes
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", APPKEY, APPSECRET)
	fmt.Println(url)
	body, err := method.HttpGet(url)
	if err != nil {
		return "", err
	}
	_ = json.Unmarshal(body, &accessToken)
	if accessToken.ErrCode != 0 {
		return "", errors.New(accessToken.Errmsg)
	}
	return accessToken.AccessToken, nil
}

// 对外的获取accesstoken
func GetToken(c *gin.Context) {
	token, err := getAccessToken()
	code := 0
	msg := "success"
	if err != nil {
		code = -3001
		msg = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "token": token, "msg": msg})
}

// 小程序获取用户userinfo
func UserInfo(c *gin.Context) {
	code := c.Query("code")
	accessToken := c.Query("access_token")
	url := fmt.Sprintf("%s/getuserinfo?access_token=%s&code=%s", UserUrl, accessToken, code)
	body, err := method.HttpGet(url)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -3001, "msg": err.Error()})
		return
	}
	type userinfoStruct struct {
		Userid  string `json:"userid"`
		Errmsg  string `json:"errmsg"`
		Errcode int    `json:"errcode"`
	}
	var userinfo userinfoStruct
	_ = json.Unmarshal([]byte(body), &userinfo)
	if userinfo.Errcode != 0 {
		c.JSON(http.StatusOK, gin.H{"code": -3002, "msg": userinfo.Errmsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "userid": userinfo.Userid})
}
