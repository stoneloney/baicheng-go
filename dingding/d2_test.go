package main

import (
	"database/sql"
	"dingding/controllers"
	"dingding/json"
	"dingding/models"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 9段部门同步到钉钉
func Test_depNew(t *testing.T) {
	res, err := controllers.RsyncDepNew()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 9段修改同步到钉钉
func Test_depEdit(t *testing.T) {
	res, err := controllers.RsyncDepEdit()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 9段删除同步到钉钉
func Test_depDel(t *testing.T) {
	res, err := controllers.RsyncDepDel()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 9段员工新增同步到钉钉
func Test_EmpNew(t *testing.T) {
	res, err := controllers.RsyncEmpNew()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 9段员工修改同步到钉钉
func Test_EmpEdit(t *testing.T) {
	res, err := controllers.RsyncEmpEdit()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 9段员工删除同步到钉钉
func Test_EmpDel(t *testing.T) {
	res, err := controllers.RsyncEmpDel()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Test_del(t *testing.T) {
	var isdebug = true
	var server = "2g399409j8.zicp.vip"
	var port = 23284
	var user = "sa"
	var password = "123"
	var database = "topHR"

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)
	if isdebug {
		fmt.Println(connString)
	}
	//建立连接
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		//log.Fatal("Open Connection failed:", err.Error())
		fmt.Println("Open Connection failed:", err.Error())
	}
	defer conn.Close()

	//sqlStr := fmt.Sprintf("DELETE FROM [DD_DepConvert]")
	//sqlStr := fmt.Sprintf("DELETE FROM [DD_DepEdit]")
	sqlStr := fmt.Sprintf("DELETE FROM [DD_DepNew] WHERE depcode=03")

	stmt, err := conn.Prepare(sqlStr)
	if err != nil {
		//log.Fatal("Prepare failed:", err.Error())
		fmt.Println("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	//通过Statement执行查询
	rows, err := stmt.Exec()
	if err != nil {
		//log.Fatal("Query failed:", err.Error())
		fmt.Println("exec failed:", err.Error())
	}
	fmt.Println(rows)
}

func Test_sql(t *testing.T) {
	var isdebug = true
	var server = "61.142.247.115"
	var port = 23284
	var user = "sa"
	var password = "gt2019"
	var database = "topHR"

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)
	if isdebug {
		fmt.Println(connString)
	}
	//建立连接
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		//log.Fatal("Open Connection failed:", err.Error())
		fmt.Println("Open Connection failed:", err.Error())
	}
	defer conn.Close()

	depId := "123456"
	depCode := "05"
	fmt.Println(depId)
	fmt.Println(depCode)

	//nowStr := method.NowStr()
	//sqlStr := fmt.Sprintf("UPDATE [DD_DepNew] SET toDD='1', runtime='%s' WHERE DepCode='%s'", nowStr, depCode)

	//sqlStr := fmt.Sprintf("INSERT INTO [DD_DepConvert](DD_DepCode, HR_DepCode) VALUES('%s', '%s')", depId, depCode)
	//sqlStr := fmt.Sprintf("INSERT INTO [DD_DepEdit](DepCode, DepName, toDD, runtime) VALUES('%s','%s', '%d', '')", "06", "钉钉部门02to03", 0)
	//sqlStr := fmt.Sprintf("INSERT INTO [DD_DepDel](DepCode, DepName, toDD, runtime) VALUES('%s','%s', '%d', '')", "05", "钉钉部门01", 0)
	//sqlStr := fmt.Sprintf("INSERT INTO [DD_DepNew](DepCode, DepName, toDD, runtime) VALUES('%s','%s', '%d', '')", "05", "钉钉部门01", 0)
	//sqlStr := fmt.Sprintf("DELETE FROM [DD_DepConvert] WHERE DD_DepCode='%s'", "120756162")
	//sqlStr := fmt.Sprintf("INSERT INTO [DD_EmpEdit](ID, toDD, runtime) VALUES('%s','%s','')", "DD_01", "0")
	sqlStr := fmt.Sprintf("INSERT INTO [DD_EmpDel](ID, toDD, runtime) VALUES('%s','%s','')", "DD_01", "0")

	stmt, err := conn.Prepare(sqlStr) // 部门新增
	//stmt, err := conn.Prepare(`select * from [DD_DepEdit]`)   // 部门修改

	//stmt, err := conn.Prepare(`select * from [DD_EmpNew]`)  // 新用户

	if err != nil {
		//log.Fatal("Prepare failed:", err.Error())
		fmt.Println("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	//通过Statement执行查询
	rows, err := stmt.Exec()
	if err != nil {
		//log.Fatal("Query failed:", err.Error())
		fmt.Println("exec failed:", err.Error())
	}
	fmt.Println(rows)

	return
}

// 审批实例
func Test_process(t *testing.T) {
	pid := "ad91761c-7adb-4ddb-afdd-f7c70671b295"   // 加班通过
	// pid := "b6ce1145-e6d7-422d-939e-24ee1d792ef9"   // 加班未通过
	//pid := "673c9163-923e-4186-845b-0a12d6bda1fd"
	//pid := "7cb660d6-2a78-4e4d-a1c3-efba408aeafd" // 请假
	//pid := "cd3d9018-097b-45c9-9e3f-e9a970ba9e0d"     //外出
	//pid := "1488e2b6-111e-4936-a5c6-a67585073844"      // 出差
	//pid := "e83c3cfe-66d0-4c02-81d0-500c5afdcc9d"      // 补卡
	res, err := controllers.ProcessinstanceInfo(pid)
	if err != nil {
		fmt.Println("error:%s" + err.Error())
	} else {
		processInstance := res.ProcessInstance
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			recordDesc := ""                                     //processInstance.FormComponentValues[0].Value // 补卡原因
			originatorUserid := processInstance.OriginatorUserid // 发起人
			recordTime := ""                                     // processInstance.FormComponentValues[2].Value // 补卡的时间点

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
			fmt.Println(timeint)
			tu := time.Unix(timeint, 0)

			time1 := tu.Format(timelayout1)
			time2 := tu.Format(timelayout2)

			// 补卡类型，再钉钉上没看到类型设置，这里暂时写死  因公
			recordType := 1

			// 发送补卡类型到hr
			sqlStr := fmt.Sprintf("exec DD_InsertBukaToHR '%s','%s','%s','%d','%s'", originatorUserid, time1, time2, recordType, recordDesc)
			/*
				Logger.Info(sqlStr)
				_, err = DbExec(sqlStr)
				if err != nil {
					Logger.Error("exec DD_InsertBukaToHR error:"+err.Error())
				}
			*/
			fmt.Println(sqlStr)
		}
	}
}

// 加班
func Test_overtime(t *testing.T) {
	//pid := "ad91761c-7adb-4ddb-afdd-f7c70671b295"   // 加班通过
	pid := "5dcef264-3e83-4831-ae71-d1238194a873"
	data, err := controllers.ProcessinstanceInfo(pid)
	if err != nil {
		fmt.Println("insert overtime error:" + err.Error())
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
				createTImeInt, err := time.Parse("2006-01-02 15:04:05", createTIme)
				createTimeformat := createTImeInt.Format(timelayout1)

				// 加班开始时间
				fmt.Println(startTIme.(string))
				startTImeInt, _ := time.Parse("2006-01-02 15:04", startTIme.(string))
				startTImeformat := startTImeInt.Format(timelayout2)

				// 加班结束时间
				fmt.Println(endTIme.(string))
				endTImeInt, _ := time.Parse("2006-01-02 15:04", endTIme.(string))
				endTImeformat := endTImeInt.Format(timelayout2)

				// 发送加班到九段
				sqlStr := fmt.Sprintf("exec DD_InsertOvertimeToHR '%s','%s','%s','%s','%.1f','%s','%s'", originatorUserid, createTimeformat, startTImeformat, endTImeformat, duration, overtimeType, overtimeDesc)
				fmt.Println(sqlStr)
				if err != nil {
					fmt.Println("exec DD_InsertOvertimeToHR error:" + err.Error())
				}
			} else {
				fmt.Println("exec DD_InsertOvertimeToHR cvRes empty")
			}
		}
	}
}

func Test_overtime2(t *testing.T) {
	pid := "651327db-4f9a-4b3b-8f3b-aa38ecd9d757"
	data, err := controllers.ProcessinstanceInfo(pid)
	if err != nil {
		fmt.Println("insert overtime error:" + err.Error())
	} else {
		processInstance := data.ProcessInstance
		// 流程完成并且结果为同意  写入9段
		if processInstance.Status == "COMPLETED" && processInstance.Result == "agree" {
			// 表单详情
			componentValue := processInstance.FormComponentValues[0].Value // 加班详情
			// 表单详情2
			overtimeDesc := processInstance.FormComponentValues[1].Value // 加班原因

			originatorUserid := processInstance.OriginatorUserid // 发起人
			//createTIme := processInstance.CreateTime             // 发起时间

			var cvRes []models.NewComponentValues
			_ = json.Unmarshal([]byte(componentValue), &cvRes)

			startTIme := ""    // 加班开始时间
			endTIme := ""      // 加班结束时间
			duration := ""     // 加班时长
			overtimeType := "" // 加班类型
			overtimeType = "1" // 常用加班

			// 查找数据
			for _, v := range cvRes {
				switch v.Props.BizAlias {
				case "startTime":{
					startTIme = v.Value
				}
				case "finishTime":{
					endTIme = v.Value
				}
				case "duration":{
					duration = v.Value
				}

				}
			}

			timelayout1 := "2006-01-02"
			timelayout2 := "15:04"

			/*
			createTImeInt, err := time.Parse("2006-01-02 15:04:05", createTIme)
			createTimeformat := createTImeInt.Format(timelayout1)
			*/
			// 加班开始时间
			//fmt.Println(startTIme)
			startTImeInt, _ := time.Parse("2006-01-02 15:04", startTIme)
			startTImeformat := startTImeInt.Format(timelayout2)
			startTImeformat2 := startTImeInt.Format(timelayout1)

			// 加班结束时间
			//fmt.Println(endTIme)
			endTImeInt, _ := time.Parse("2006-01-02 15:04", endTIme)
			endTImeformat := endTImeInt.Format(timelayout2)

			// 发送加班到九段
			sqlStr := fmt.Sprintf("exec DD_InsertOvertimeToHR '%s','%s','%s','%s','%s','%s','%s'", originatorUserid, startTImeformat2, startTImeformat, endTImeformat, duration, overtimeType, overtimeDesc)
			fmt.Println(sqlStr)
			if err != nil {
				fmt.Println("exec DD_InsertOvertimeToHR error:" + err.Error())
			}
		}
	}
}

// 包含测试
func Test_index(t *testing.T) {
	str := "xx提交的加班"
	find := "buka"
	i := strings.Index(str, find)
	fmt.Println(i)
}

// 获取员工具体信息
func Test_userinfo(t *testing.T) {
	userid := "033435111821927863 "
	info, err := controllers.GetUserInfo(userid)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(info)
	}
	fmt.Println(info.Mobile)
}

// 获取部门详情
func Test_deptinfo(t *testing.T) {
	//depid := 120863221
	depid := int64(130100781)
	/*
		info, err := controllers.DepartmentInfo(depid)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(info)
			fmt.Println(len(info.SourceIdentifier))
		}
	*/
	data, err := controllers.DepartmentInfo(depid)
	if err != nil {
		//Logger.Error("get departinfo error:"+err.Error())
		fmt.Println("get departinfo error:" + err.Error())
	} else {
		// 查看是否顶级部门
		depid := ""
		if len(data.SourceIdentifier) > 0 {
			depid = data.SourceIdentifier
		} else {
			depid = "sss"
		}
		fmt.Println(depid)
	}
}

// 获取部门id
func Test_hrid(t *testing.T) {
	type makeDepidStruct struct {
		MaxDepCode sql.NullString `sql:"maxDepCode"`
	}
	//sqlStr := "select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE UpDepCode='06'"
	sqlStr := "select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE UpDepCode='09'"
	data, err := controllers.DbQuery(sqlStr, (*makeDepidStruct)(nil))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		depid := ""
		t := data[0].(*makeDepidStruct)
		fmt.Println(t)
		if t.MaxDepCode.Valid {
			depid = t.MaxDepCode.String
		} else {
			depid = "0801"
		}
		//t := data[0].(*makeDepidStruct)
		//depid := t.MaxDepCode
		//depid = controllers.calculateHrid(depid)
		fmt.Println(depid)
	}
}

func Test_exec(t *testing.T) {
	/*
		userId := "ceshi3"
		name := "ceshi3"
		departmentId := "05"
		position := "美工"
		hiredDateStr := "2019-07-02"
		mobile := "13488873456"
		//sqlStr := fmt.Sprintf("exec DD_InsertEmpToHR '%s','%s','%s','%s','%s','%s'", userId, name, departmentId, position, hiredDateStr, mobile)
		var code int
		//sqlStr := "exec DD_InsertEmpToHR 'ceshi3','ceshi3','05','美工','2019-07-02','13488873456'"
		//sqlStr := fmt.Sprintf("exec DD_InsertEmpToHR '%s','%s','%s','%s','%s','%s', %d", userId, name, departmentId, position, hiredDateStr, mobile, code)
	*/
	//sqlStr := "exec DD_insertcardToHR '001','2019-06-01','12:00'"  //成功
	sqlStr := "exec DD_GetAllEmpInfoFromHR"
	fmt.Println(sqlStr)
	res, err := controllers.DbExec(sqlStr)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(res)
	}
}

func Test_time(t *testing.T) {
	timelayout := "2006-01-02"
	timelayout2 := "15:04"

	timeint := int64(1561903018)
	fmt.Println(timeint)
	to := time.Unix(timeint, 0)

	dateStr := to.Format(timelayout)
	timeStr := to.Format(timelayout2)

	fmt.Println(dateStr)
	fmt.Println(timeStr)

	staffId := "171802165236043309"
	sqlStr := fmt.Sprintf("exec DD_insertcardToHR '%s','%s','%s'", staffId, dateStr, timeStr)
	sqlStr = "exec DD_insertcardToHR '171802165236043309','2019-06-30','22:47'"

	_, err := controllers.DbExec(sqlStr)
	if err != nil {
		fmt.Println("error:", err.Error())
	}
}

func Test_str(t *testing.T) {
	var res []string
	val := "[\"2019-06-30 00:00\",\"2019-06-30 01:00\",1,\"hour\",null,\"加班类型\"]"
	_ = json.Unmarshal([]byte(val), &res)
	fmt.Println(res[0])
	fmt.Println(res[1])
}

func Test_add(t *testing.T) {
	str := "05"
	len1 := len(str)
	strint, _ := strconv.Atoi(str)
	strint += 1
	intstr := strconv.Itoa(strint)
	len2 := len(intstr)
	fmt.Println(strint)
	fmt.Println(len1)
	fmt.Println(len2)
	if len1 > len2 {
		intstr = "0" + intstr
	}
	fmt.Println(intstr)

}

// 考勤异常
func Test_workerror(t *testing.T) {
	userId := "171802165236043309"
	timeDesc := "2019-05"
	//sqlStr := fmt.Sprintf("exec DD_GetWorkError '%s','%s'", userId, timeDesc)
	sqlStr := fmt.Sprintf("exec DD_GetWorkDetail '%s','%s'", userId, timeDesc)
	fmt.Println(sqlStr)
	//res, err := controllers.DbQuery(sqlStr)

}

func Test_record(t *testing.T) {
	//_, _ = controllers.RsyncRecord()
}

func init() {
	controllers.LoadConfig()
	controllers.SetAppKey(controllers.Cfg.Dingding.Appkey)
	controllers.SetSecret(controllers.Cfg.Dingding.Appsecret)
}
