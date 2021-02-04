package controllers

import (
	"database/sql"
	"dingding/json"
	"dingding/method"
	"dingding/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

// 同步新增的部门
func RsyncDepNew() (bool, error) {
	queryStr := "SELECT * FROM [DD_DepNew] WHERE toDD=0 ORDER BY depCode ASC" // 先添加顶级部门，再添加子部门
	//queryStr := "SELECT * FROM [DD_DepNew] WHERE toDD=1 AND UpDepCode=''"
	data, err := DbQuery(queryStr, (*models.DepNewStruct)(nil))
	if err != nil {
		fmt.Println("select depnew error:" + err.Error())
		Logger.Error(fmt.Sprintf("select depnew error:%s", err.Error()))
		return false, err
	}
	Logger.Info(fmt.Sprintf("depnew:%v", data))
	// 有新增数据
	if len(data) > 0 {
		// 写入钉钉
		for _, v := range data {
			t := v.(*models.DepNewStruct)
			name := t.DepName
			parentid := t.UpDepCode
			depCode := t.DepCode
			if parentid == "" {
				parentid = "1"
			} else {
				// 获取钉钉的depid
				parentid, err = getDepidByHrid(parentid)
				if err != nil {
					Logger.Error(fmt.Sprintf("dep create error: name:%s, depcode:%s, parentid:%s, error:%s", name, depCode, parentid, err.Error()))
				}
			}
			// 获取dingding ip
			depId, err := DepartmentCreate(name, parentid, depCode)
			fmt.Println(depId)
			if err != nil {
				fmt.Println(err.Error())
				Logger.Error(fmt.Sprintf("dep create error: name:%s, depcode:%s, parentid:%s, error:%s", name, depCode, parentid, err.Error()))
			} else {
				// 添加部门的映射关系
				/*
					depIdStr := strconv.FormatInt(depId,10)
					convertStr := fmt.Sprintf("INSERT INTO [DD_DepConvert](DD_DepCode, HR_DepCode) VALUES('%s', '%s')", depIdStr, depCode)
					fmt.Println(convertStr)
					Logger.Info("converstr:"+convertStr)
					_, err := DbExec(convertStr)
					if err != nil {
						fmt.Println(err.Error())
						Logger.Error(fmt.Sprintf("update error:name:%s, depcode:%s, error:%s", name, depCode, err.Error()))
					}
				*/
				// 更新9段数据库
				nowStr := method.NowStr()
				updateStr := fmt.Sprintf("UPDATE [DD_DepNew] SET toDD='1', runtime='%s' WHERE DepCode='%s'", nowStr, depCode)
				// UPDATE [DD_DepNew] SET toDD='1', runtime='2019-06-30 10:19:36' WHERE DepCode='05'
				//fmt.Println(updateStr)
				_, err = DbExec(updateStr)
				if err != nil {
					fmt.Println(err.Error())
					Logger.Error(fmt.Sprintf("update error:name:%s, depcode:%s, error:%s", name, depCode, err.Error()))
				}
			}
		}
	}
	return true, nil
}

// 同步部门的修改
func RsyncDepEdit() (bool, error) {
	queryStr := "SELECT * FROM [DD_DepEdit] WHERE toDD=0"
	data, err := DbQuery(queryStr, (*models.DepEditStruct)(nil))
	if err != nil {
		Logger.Error(fmt.Sprintf("select depedit error:%s", err.Error()))
		return false, err
	}
	if len(data) > 0 {
		// 钉钉侧修改数据
		for _, v := range data {
			t := v.(*models.DepEditStruct)
			depName := t.DepName
			depCode := t.DepCode
			// 获取钉钉id
			//fmt.Println("depCode:", depCode)
			ddid, err := getDepidByHrid(depCode)
			//fmt.Println("ddid:", ddid)
			if err != nil {
				fmt.Println("get error:", err.Error())
				Logger.Error("getDepidByHrid error:%s" + err.Error())
			} else {
				// 修改数据
				ddidInt, _ := strconv.ParseInt(ddid, 10, 64)
				_, err := DepartmentUpdate(ddidInt, depName)
				if err != nil {
					Logger.Error(fmt.Sprintf("eidt error:name:%s, depcode:%s, error:%s", depName, depCode, err.Error()))
				} else {
					// 修改更新
					nowStr := method.NowStr()
					updateStr := fmt.Sprintf("UPDATE [DD_DepEdit] SET toDD='1', runtime='%s' WHERE DepCode='%s'", nowStr, depCode)
					//fmt.Println(updateStr)
					_, err = DbExec(updateStr)
					if err != nil {
						Logger.Error(fmt.Sprintf("update error:name:%s, depcode:%s, error:%s", depName, depCode, err.Error()))
					}
				}
			}
		}
	}
	return true, nil
}

// 同步删除部门
func RsyncDepDel() (bool, error) {
	queryStr := "SELECT * FROM [DD_DepDel] WHERE toDD=0 ORDER BY depCode DESC" // 先删除子部门，在删除顶级部门
	//queryStr := "SELECT * FROM [DD_DepDel]"
	data, err := DbQuery(queryStr, (*models.DepDelStruct)(nil))
	if err != nil {
		return false, err
	}
	if len(data) > 0 {
		// 删除钉钉侧
		for _, v := range data {
			t := v.(*models.DepDelStruct)
			depName := t.DepName
			depCode := t.DepCode
			//fmt.Println("depCode:", depCode)
			ddid, err := getDepidByHrid(depCode)
			//fmt.Println("ddid:", ddid)
			if err != nil {
				fmt.Println("get error:", err.Error())
				Logger.Error("getDepidByHrid error:" + err.Error())
			} else {
				// 修改数据
				ddidInt, _ := strconv.ParseInt(ddid, 10, 64) // 转为int64
				_, err := DepartmentDelete(ddidInt)
				if err != nil {
					fmt.Println("del error:", err.Error())
					Logger.Error(fmt.Sprintf("del error:name:%s, depcode:%s, error:%s", depName, depCode, err.Error()))
				} else {
					// 修改更新
					nowStr := method.NowStr()
					updateStr := fmt.Sprintf("UPDATE [DD_DepDel] SET toDD='1', runtime='%s' WHERE DepCode='%s'", nowStr, depCode)
					//fmt.Println(updateStr)
					_, err = DbExec(updateStr)
					if err != nil {
						fmt.Println("exec error:", err.Error())
						Logger.Error(fmt.Sprintf("del error:name:%s, depcode:%s, error:%s", depName, depCode, err.Error()))
					}
					// 删除部门相应的映射
					delStr := fmt.Sprintf("DELETE FROM [DD_DepConvert] WHERE HR_DepCode='%s'", depCode)
					//fmt.Println(delStr)
					Logger.Info(delStr)
					_, err = DbExec(delStr)
					if err != nil {
						fmt.Println("exec error:", err.Error())
						Logger.Error(fmt.Sprintf("del convert error:%s, depcode:%s", err.Error(), depCode))
					}
				}
			}
		}
	}
	return true, nil
}

// 同步新增的用户
func RsyncEmpNew() (bool, error) {
	queryStr := "SELECT * FROM [DD_EmpNew] WHERE toDD=0"
	//queryStr := "SELECT * FROM [DD_EmpNew]"
	data, err := DbQuery(queryStr, (*models.EmpNewStruct)(nil))
	if err != nil {
		Logger.Error(fmt.Sprintf("select empnew error:%s", err.Error()))
		return false, err
	}
	//Logger.Info(fmt.Sprintf("empnew:%v", data))
	if len(data) > 0 {
		for _, v := range data {
			t := v.(*models.EmpNewStruct)
			userid := t.Id
			// 获取员工的基本信息
			empInfo, err := getUserInfoByHrid(userid)
			if err != nil {
				fmt.Println("get empinfo error:", err.Error())
				Logger.Error(fmt.Sprintf("get empinfo error:%s", err.Error()))
			} else {
				empId := empInfo.Id           // 工号
				name := empInfo.Empname       // 名字
				department := empInfo.Depcode // 部门编号
				duty := empInfo.Duty          // 职务
				mobile := empInfo.Movephone   // 手机号
				hiredDate := empInfo.InDate   // 入职日期

				// 查找出钉钉部门id
				ddDepid, _ := getDepidByHrid(department)
				ddDepidInt, _ := strconv.ParseInt(ddDepid, 10, 64)
				var departmentIds []int64
				departmentIds = append(departmentIds, ddDepidInt)

				// 手机号不能为空
				if len(mobile) == 0 {
					mobile = ""
				}

				// 要为时间戳
				var hiredTime int64
				if len(hiredDate) > 0 {
					timeLayout := "2006-01-02 15:04:05"
					loc, _ := time.LoadLocation("Local")
					tmp, _ := time.ParseInLocation(timeLayout, hiredDate, loc)
					hiredTime = tmp.Unix() //转化为时间戳 类型是int64
				} else {
					hiredTime = time.Now().Unix()
				}
				_, err := UserOp("create", empId, name, departmentIds, duty, mobile, "", "", "", "", "", empId, false, false, "", hiredTime)
				if err != nil {
					fmt.Println("user op error:%s", err.Error())
					Logger.Error(fmt.Sprintf("user op error:%s, userid:%s", err.Error(), empId))
				} else {
					// 修改更新
					nowStr := method.NowStr()
					updateStr := fmt.Sprintf("UPDATE [DD_EmpNew] SET toDD='1', runtime='%s' WHERE ID='%s'", nowStr, empId)
					//fmt.Println(updateStr)
					_, err = DbExec(updateStr)
					if err != nil {
						fmt.Println("exec error:", err.Error())
						Logger.Error(fmt.Sprintf("user op error:%s, userid:%s", err.Error(), empId))
					}
				}
			}
		}
	}
	return true, nil
}

// 同步员工的修改
func RsyncEmpEdit() (bool, error) {
	queryStr := "SELECT * FROM [DD_EmpEdit] WHERE toDD=0"
	//queryStr := "SELECT * FROM [DD_EmpEdit]"
	data, err := DbQuery(queryStr, (*models.EmpEditStruct)(nil))
	if err != nil {
		Logger.Error(fmt.Sprintf("select empedit error:%s", err.Error()))
		return false, err
	}
	//Logger.Info(fmt.Sprintf("empedit:%v", data))
	if len(data) > 0 {
		for _, v := range data {
			t := v.(*models.EmpEditStruct)
			userid := t.Id
			// 获取员工的基本信息
			empInfo, err := getUserInfoByHrid(userid)
			if err != nil {
				fmt.Println("get empinfo error:", err.Error())
				Logger.Error(fmt.Sprintf("get empinfo error:%s", err.Error()))
			} else {
				empId := empInfo.Id           // 工号
				name := empInfo.Empname       // 名字
				department := empInfo.Depcode // 部门编号
				duty := empInfo.Duty          // 职务
				mobile := empInfo.Movephone   // 手机号
				hiredDate := empInfo.InDate   // 入职日期

				// 查找出钉钉部门id
				ddDepid, _ := getDepidByHrid(department)
				ddDepidInt, _ := strconv.ParseInt(ddDepid, 10, 64)
				var departmentIds []int64
				departmentIds = append(departmentIds, ddDepidInt)

				// 手机号不能为空
				if len(mobile) == 0 {
					mobile = ""
				}

				// test
				// name = "修改为02"

				// 要为时间戳
				var hiredTime int64
				if len(hiredDate) > 0 {
					timeLayout := "2006-01-02 15:04:05"
					loc, _ := time.LoadLocation("Local")
					tmp, _ := time.ParseInLocation(timeLayout, hiredDate, loc)
					hiredTime = tmp.Unix() //转化为时间戳 类型是int64
				}

				_, err := UserOp("update", empId, name, departmentIds, duty, mobile, "", "", "", "", "", empId, false, false, "", hiredTime)
				if err != nil {
					fmt.Println("user op error:%s", err.Error())
					Logger.Error(fmt.Sprintf("user op error:%s, userid:%s", err.Error(), empId))
				} else {
					// 修改更新
					nowStr := method.NowStr()
					updateStr := fmt.Sprintf("UPDATE [DD_EmpEdit] SET toDD='1', runtime='%s' WHERE ID='%s'", nowStr, empId)
					//fmt.Println(updateStr)
					_, err = DbExec(updateStr)
					if err != nil {
						fmt.Println("exec error:", err.Error())
						Logger.Error(fmt.Sprintf("user op error:%s, userid:%s", err.Error(), empId))
					}
				}
			}
		}
	}
	return true, nil
}

// 用户删除
func RsyncEmpDel() (bool, error) {
	queryStr := "SELECT * FROM [DD_EmpDel] WHERE toDD=0"
	data, err := DbQuery(queryStr, (*models.EmpDelStruct)(nil))
	if err != nil {
		Logger.Error(fmt.Sprintf("select empedit error:%s", err.Error()))
		return false, err
	}
	//Logger.Info(fmt.Sprintf("empedit:%v", data))
	if len(data) > 0 {
		// 删除钉钉侧
		for _, v := range data {
			t := v.(*models.EmpDelStruct)
			userid := t.Id
			// 获取员工的基本信息
			empInfo, err := getUserInfoByHrid(userid)
			if err != nil {
				fmt.Println("get empinfo error:", err.Error())
				Logger.Error(fmt.Sprintf("get empinfo error:%s", err.Error()))
			} else {
				empId := empInfo.Id // 工号
				_, err := UserDel(empId)
				if err != nil {
					fmt.Println("user del error:%s", err.Error())
					Logger.Error(fmt.Sprintf("user op error:%s, userid:%s", err.Error(), empId))
				} else {
					// 修改更新
					nowStr := method.NowStr()
					updateStr := fmt.Sprintf("UPDATE [DD_EmpDel] SET toDD='1', runtime='%s' WHERE ID='%s'", nowStr, empId)
					//fmt.Println(updateStr)
					_, err = DbExec(updateStr)
					if err != nil {
						fmt.Println("exec error:", err.Error())
						Logger.Error(fmt.Sprintf("user del error:%s, userid:%s", err.Error(), empId))
					}
				}
			}
		}
	}
	return true, nil
}

// 通过钉钉来同步打卡记录
func RsyncRecord2() error {
	// 通过钉钉来获取员工列表
	token, err := getAccessToken()
	if err != nil {
		return errors.New("get accesstoken error:" + err.Error())
	}

	url := fmt.Sprintf("%s/queryonjob?access_token=%s", SmartworkUrl, token)
	//fmt.Println(url)
	data := make(map[string]interface{})
	data["status_list"] = "2,3,5,-1"
	data["offset"] = 0
	data["size"] = 10
	// 发送请求
	dataStr, _ := json.Marshal(data)
	//fmt.Println(string(dataStr))
	body, err := method.HttpPostJson(url, dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(body))
	return nil
}

// 推送工资信息
func PushCash() {
	type UserIdStruct struct {
		UserId  string  `sql:"userID"`
	}
	type UserNoteStruct struct {
		Note   string  `sql:"note"`
	}
	// 获取推送的用户列表
	queryStr := "SELECT userID FROM [DD_SendInfoUserID]"
	idDatas, err := DbQuery(queryStr, (*UserIdStruct)(nil))
	if err != nil {
		Logger.Error(fmt.Sprintf("select SendInfoUserID error:%s", err.Error()))
		fmt.Println(fmt.Sprintf("select SendInfoUserID error:%s", err.Error()))
		return
	}

	token, err := getAccessToken()
	if err != nil {
		Logger.Error(fmt.Sprintf("getAccessToken error:%s", err.Error()))
		fmt.Println(fmt.Sprintf("getAccessToken error:%s", err.Error()))
		return
	}

	for _, v := range idDatas {
		t := v.(*UserIdStruct)
		if len(t.UserId) > 0 {
			// 查询用户具体工资信息
			queryStr := fmt.Sprintf("SELECT note FROM [DD_SendInfo] WHERE userid='%s' order by SN", t.UserId)
			noteDatas, err := DbQuery(queryStr, (*UserNoteStruct)(nil))
			if err != nil {
				Logger.Error(fmt.Sprintf("select note error:%s, userId:%s", err.Error(), t.UserId))
				fmt.Println(fmt.Sprintf("select note error:%s, userId:%s", err.Error(), t.UserId))
				continue
			}
			// 发送推送通知
			var pushMsgArr []string
			for _, n := range noteDatas {
				t2 := n.(*UserNoteStruct)
				pushMsgArr = append(pushMsgArr, t2.Note)
			}
			//fmt.Println(t.UserId)
			//fmt.Println(pushMsgArr)
			_, err = PushTextMsg(token, t.UserId, pushMsgArr)
			if err != nil {
				fmt.Println(err.Error)
				Logger.Error(fmt.Sprintf("pushMsg error:%s, userId:%s", err.Error(), t.UserId))
			} else {
				fmt.Println(fmt.Sprintf("push success, userId:%s", t.UserId))
			}
		}
	}
}

// 同步打卡记录
func RsyncRecord(startDate, endDate string) (bool, error) {
	type recordStruct struct {
		ID string `sql:"ID"`
	}
	sqlStr := fmt.Sprintf("exec DD_GetAllEmpInfoFromHR")
	data, err := DbQuery(sqlStr, (*recordStruct)(nil))

	wrongSourceType := []string{"SYSTEM", "AUTO_CHECK"}

	if err != nil {
		//fmt.Println("exec DD_GetAllEmpInfoFromHR error:%s"+err.Error())
		Logger.Error("exec DD_GetAllEmpInfoFromHR error:%s" + err.Error())
	} else {
		dataLen := len(data)

		Logger.Info(fmt.Sprintf("datalen:%d", data))

		if dataLen > 0 {

			maxNum := 50
			if dataLen > maxNum {
				count := math.Ceil(float64(dataLen) / float64(maxNum))

				Logger.Info(fmt.Sprintf("count:%d", int(count)))

				for i := 0; i < int(count); i++ {
					start := i * maxNum
					end := start + maxNum
					if end > dataLen {
						end = dataLen
					}
					arr := data[start:end]
					var ids []string
					for _, v := range arr {
						t := v.(*recordStruct)
						ids = append(ids, t.ID)
					}

					Logger.Info(fmt.Sprintf("ids:%s", ids))

					_, list, err := AttendanceRecordList(ids, startDate, endDate)
					if err != nil {
						Logger.Error("AttendanceList error:%s" + err.Error())
					} else {
						timeLayout3 := "2006-01-02"
						timeLayout4 := "15:04"

						for _, v2 := range list {
							if v2.UserCheckTime > 0 && !method.InArray(v2.SourceType, wrongSourceType) {
								sqlStr2 := fmt.Sprintf("exec DD_InsertCardToHR '%s','%s','%s','%s'", v2.UserId, time.Unix(v2.UserCheckTime/1000, 0).Format(timeLayout3), time.Unix(v2.UserCheckTime/1000, 0).Format(timeLayout4), v2.UserAddress)
								//fmt.Println(sqlStr2)
								Logger.Info(sqlStr2)
								_, err := DbExec(sqlStr2)
								if err != nil {
									Logger.Error("DD_InsertCardToHR error:" + err.Error())
								}
							}
						}
					}
				}
			} else {
				var ids []string
				for _, v := range data {
					t := v.(*recordStruct)
					ids = append(ids, t.ID)
				}

				_, list, err := AttendanceRecordList(ids, startDate, endDate)
				if err != nil {
					Logger.Error("AttendanceList error:%s" + err.Error())
				} else {
					timeLayout3 := "2006-01-02"
					timeLayout4 := "15:04"

					for _, v2 := range list {
						if v2.UserCheckTime > 0 && !method.InArray(v2.SourceType, wrongSourceType) {
							sqlStr2 := fmt.Sprintf("exec DD_InsertCardToHR '%s','%s','%s','%s'", v2.UserId, time.Unix(v2.UserCheckTime/1000, 0).Format(timeLayout3), time.Unix(v2.UserCheckTime/1000, 0).Format(timeLayout4), v2.UserAddress)
							//fmt.Println(sqlStr2)
							Logger.Info(sqlStr2)
							_, err := DbExec(sqlStr2)
							if err != nil {
								Logger.Error("DD_InsertCardToHR error:" + err.Error())
							}
						}
					}
				}
			}
		}
	}
	return true, nil
}

// 根据hr部门id获取钉钉部门id
func getDepidByHrid(hrdepid string) (string, error) {
	sqlStr := fmt.Sprintf("SELECT * FROM [DD_DepConvert] WHERE HR_DepCode='%s'", hrdepid)
	data, err := DbQuery(sqlStr, (*models.DepConvertStruct)(nil))
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("数据不存在")
	}
	t := data[0].(*models.DepConvertStruct)
	ddDepCode := t.DDdepCode
	return ddDepCode, nil
}

// 根据钉钉部门id获取hr部门id
func getDepidByDdid(ddid string) (string, error) {
	sqlStr := fmt.Sprintf("SELECT * FROM [DD_DepConvert] WHERE DD_DepCode='%s'", ddid)
	data, err := DbQuery(sqlStr, (*models.DepConvertStruct)(nil))
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("数据不存在")
	}
	t := data[0].(*models.DepConvertStruct)
	hrDepCode := t.HRdepCode
	return hrDepCode, nil
}

// 获取hr系统员工基本信息
func getUserInfoByHrid(hrid string) (*models.EmpInfoStruct, error) {
	var empInfo *models.EmpInfoStruct
	sqlStr := fmt.Sprintf("exec DD_GetEmpInfoFromHR '%s'", hrid)
	data, err := DbQuery(sqlStr, (*models.EmpInfoStruct)(nil))
	if err != nil {
		return empInfo, err
	}
	if len(data) > 0 {
		empInfo = data[0].(*models.EmpInfoStruct)
		return empInfo, nil
	}
	return empInfo, errors.New("empinfo empty")
}

// 生成hr部门id
func MakeHrDepId(ddid int64) (string, error) {
	type makeDepidStruct struct {
		MaxDepCode sql.NullString `sql:"maxDepCode"`
	}
	sqlStr := ""
	hrdepid := ""
	var err error
	if ddid == 1 { // 为1说明是顶级部门
		//sqlStr = "select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE len(UpDepCode)=0"
		sqlStr = "select MAX(DEPARTMENTCODE) AS maxDepCode from [department] WHERE len(DEPARTMENTCODE)=2"
		fmt.Println(sqlStr)
		Logger.Info(sqlStr)
	} else {
		// 查询钉钉部门id对应的hr部门id
		ddidStr := strconv.FormatInt(ddid, 10)
		hrdepid, _ = getDepidByDdid(ddidStr)
		if err != nil {
			Logger.Error("getDepidByDdid error:" + err.Error())
			return "", errors.New("getDepidByDdid error:" + err.Error())
		} else {
			//sqlStr = fmt.Sprintf("select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE UpDepCode='%s'", hrdepid)
			sqlStr = "select MAX(DEPARTMENTCODE) AS maxDepCode from [department] WHERE DEPARTMENTCODE LIKE '" + hrdepid + "%' AND DEPARTMENTCODE != '" + hrdepid + "'"
			fmt.Println(sqlStr)
			Logger.Info(sqlStr)
		}
	}
	if sqlStr == "" {
		return "", errors.New("sql str empty")
	}
	data, err := DbQuery(sqlStr, (*makeDepidStruct)(nil))
	if err != nil {
		return "", err
	}
	/*
		if len(data) == 0 {  // 没查到就是没数据，返回部门为01
			return "01", nil
		}
	*/
	t := data[0].(*makeDepidStruct)
	depid := ""
	fmt.Println(t)
	if t.MaxDepCode.Valid { // maxdepcode不为空时，按正常计算hr应该出现的数目
		Logger.Info("maxstring:" + t.MaxDepCode.String)
		// 查看对照表中是否含有，如果没有的话，
		depid = calculateHrid(t.MaxDepCode.String)
	} else {
		if ddid == 1 {
			return "01", nil
		} else {
			Logger.Info("hrdepid:" + hrdepid)
			return hrdepid + "01", nil
		}
	}

	return depid, nil
}

// 生成hr部门id
func calculateHrid(depid string) string {
	len1 := len(depid)
	strint, _ := strconv.Atoi(depid)
	strint += 1
	intstr := strconv.Itoa(strint)
	len2 := len(intstr)
	if len1 > len2 {
		intstr = "0" + intstr
	}
	if len(intstr) == 1 {
		intstr = "0" + intstr
	}
	return intstr
}

// 获取考勤异常
func GetWorkError(c *gin.Context) {
	if db == nil {
		db = GetDb()
	}
	// 获取参数
	uid := c.Query("uid")
	date := c.Query("date")

	if len(uid) == 0 || len(date) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": -1001, "msg": "参数错误"})
		return
	}

	// 获取用户基本信息
	userinfo, err := GetUserInfo(uid)
	if err != nil {
		Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
		c.JSON(http.StatusOK, gin.H{"code": -3001, "msg": "get userinfo error"})
	} else {
		mobile := userinfo.Mobile // userid改为手机号
		sqlStr := fmt.Sprintf("exec DD_GetWorkError '%s','%s'", mobile, date)
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			//log.Fatal("Prepare failed:", err.Error())
			fmt.Println("Prepare failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2001, "msg": err.Error()})
			return
		}
		defer stmt.Close()

		//通过Statement执行查询
		rows, err := stmt.Query()
		if err != nil {
			//log.Fatal("Query failed:", err.Error())
			fmt.Println("Query failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2002, "msg": err.Error()})
			return
		}

		//建立一个列数组
		var respData [][]string
		cols, err := rows.Columns()
		var colsdata = make([]interface{}, len(cols))
		var title []string
		for i := 0; i < len(cols); i++ {
			colsdata[i] = new(interface{})
			title = append(title, cols[i])
		}
		respData = append(respData, title)

		for rows.Next() {
			rows.Scan(colsdata...)
			var tempData []string
			for _, val := range colsdata {
				switch v := (*(val.(*interface{}))).(type) {
				case nil:
					tempData = append(tempData, "")
				case bool:
					if v {
						tempData = append(tempData, "true")
					} else {
						tempData = append(tempData, "false")
					}
				case []byte:
					//fmt.Print(string(v))
					tempData = append(tempData, string(v))
				case time.Time:
					// fmt.Print(v.Format("2016-01-02 15:05:05.999"))
					tempData = append(tempData, v.Format("2016-01-02 15:04:05"))
				default:
					tempData = append(tempData, fmt.Sprintf("%v", v))
				}
			}
			respData = append(respData, tempData)
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": respData})
	}
}

// 获取考勤
func GetWorkDetail(c *gin.Context) {
	if db == nil {
		db = GetDb()
	}

	// 获取参数
	uid := c.Query("uid")
	date := c.Query("date")

	if len(uid) == 0 || len(date) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": -1001, "msg": "参数错误"})
		return
	}

	userinfo, err := GetUserInfo(uid)
	if err != nil {
		Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
		c.JSON(http.StatusOK, gin.H{"code": -3001, "msg": "get userinfo error"})
	} else {
		mobile := userinfo.Mobile // userid改为手机号
		sqlStr := fmt.Sprintf("exec DD_GetWorkDetail '%s','%s'", mobile, date)
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			//log.Fatal("Prepare failed:", err.Error())
			fmt.Println("Prepare failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2001, "msg": err.Error()})
			return
		}
		defer stmt.Close()

		//通过Statement执行查询
		rows, err := stmt.Query()
		if err != nil {
			//log.Fatal("Query failed:", err.Error())
			fmt.Println("Query failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2002, "msg": err.Error()})
			return
		}

		//建立一个列数组
		var respData [][]string
		cols, err := rows.Columns()
		var colsdata = make([]interface{}, len(cols))
		var title []string
		for i := 0; i < len(cols); i++ {
			colsdata[i] = new(interface{})
			title = append(title, cols[i])
		}
		respData = append(respData, title)

		for rows.Next() {
			rows.Scan(colsdata...)
			var tempData []string
			for _, val := range colsdata {
				switch v := (*(val.(*interface{}))).(type) {
				case nil:
					tempData = append(tempData, "")
				case bool:
					if v {
						tempData = append(tempData, "true")
					} else {
						tempData = append(tempData, "false")
					}
				case []byte:
					//fmt.Print(string(v))
					tempData = append(tempData, string(v))
				case time.Time:
					// fmt.Print(v.Format("2016-01-02 15:05:05.999"))
					tempData = append(tempData, v.Format("2016-01-02 15:04:05"))
				default:
					tempData = append(tempData, fmt.Sprintf("%v", v))
				}
			}
			respData = append(respData, tempData)
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": respData})
	}
}

/*
func GetWorkDetail(c *gin.Context) {
	if db == nil {
		db = GetDb()
	}
	sqlStr := fmt.Sprintf("exec DD_GetWorkDetail '%s','%s'", "001", "2019-02")
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		//log.Fatal("Prepare failed:", err.Error())
		fmt.Println("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	//通过Statement执行查询
	rows, err := stmt.Query()
	if err != nil {
		//log.Fatal("Query failed:", err.Error())
		fmt.Println("Query failed:", err.Error())
	}

	//建立一个列数组
	cols, err := rows.Columns()
	var colsdata = make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		colsdata[i] = new(interface{})
	}

	type respData struct {
		Date       string   `json:"date"`
		Week       string   `json:"week"`
		AmOnduty   string   `json:"amonduty"`
		AmOffduty  string   `json:"amoffduty"`
		PmOnduty   string   `json:"pmonduty"`
		PmOffduty  string   `json:"pmoffduty"`
		Ovstart    string   `json:"ovstart"`
		Ovend      string   `json:"ovend"`
	}

	var respDatas []respData

	for rows.Next() {
		rows.Scan(colsdata...)
		var t respData
		t.Date = printInterface(colsdata[0])
		t.Week = printInterface(colsdata[1])
		t.AmOnduty = printInterface(colsdata[2])
		t.AmOffduty = printInterface(colsdata[3])
		t.PmOnduty = printInterface(colsdata[4])
		t.PmOffduty = printInterface(colsdata[5])
		t.Ovstart = printInterface(colsdata[6])
		t.Ovend = printInterface(colsdata[7])

		fmt.Println(t.Date)

		respDatas = append(respDatas, t)
	}

	c.JSON(http.StatusOK, gin.H{"code":0, "data":respDatas})
}
*/

// 获取资金
func GetSalaryDetail(c *gin.Context) {
	if db == nil {
		db = GetDb()
	}

	// 获取参数
	uid := c.Query("uid")
	date := c.Query("date")

	if len(uid) == 0 || len(date) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": -1001, "msg": "参数错误"})
		return
	}

	userinfo, err := GetUserInfo(uid)
	if err != nil {
		Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
		c.JSON(http.StatusOK, gin.H{"code": -3001, "msg": "get userinfo error"})
	} else {
		mobile := userinfo.Mobile // userid改为手机号
		sqlStr := fmt.Sprintf("exec DD_GetSalaryDetail '%s','%s'", mobile, date)
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			//log.Fatal("Prepare failed:", err.Error())
			fmt.Println("Prepare failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2001, "msg": err.Error()})
			return
		}
		defer stmt.Close()

		//通过Statement执行查询
		rows, err := stmt.Query()
		if err != nil {
			//log.Fatal("Query failed:", err.Error())
			fmt.Println("Query failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2002, "msg": err.Error()})
			return
		}

		//建立一个列数组
		var respData [][]string
		cols, err := rows.Columns()
		var colsdata = make([]interface{}, len(cols))
		var title []string
		for i := 0; i < len(cols); i++ {
			colsdata[i] = new(interface{})
			title = append(title, cols[i])
		}
		respData = append(respData, title)

		for rows.Next() {
			rows.Scan(colsdata...)
			var tempData []string
			for _, val := range colsdata {
				switch v := (*(val.(*interface{}))).(type) {
				case nil:
					tempData = append(tempData, "")
				case bool:
					if v {
						tempData = append(tempData, "true")
					} else {
						tempData = append(tempData, "false")
					}
				case []byte:
					//fmt.Print(string(v))
					tempData = append(tempData, string(v))
				case time.Time:
					// fmt.Print(v.Format("2016-01-02 15:05:05.999"))
					tempData = append(tempData, v.Format("2016-01-02 15:04:05"))
				default:
					tempData = append(tempData, fmt.Sprintf("%v", v))
				}
			}
			respData = append(respData, tempData)
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": respData})
	}
}

// 获取月报
func GetMonthDetail(c *gin.Context) {
	if db == nil {
		db = GetDb()
	}

	// 获取参数
	uid := c.Query("uid")
	date := c.Query("date")

	if len(uid) == 0 || len(date) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": -1001, "msg": "参数错误"})
		return
	}

	userinfo, err := GetUserInfo(uid)
	if err != nil {
		Logger.Error(fmt.Sprintf("get userinfo err:%s", err.Error()))
		c.JSON(http.StatusOK, gin.H{"code": -3001, "msg": "get userinfo error"})
	} else {
		mobile := userinfo.Mobile // userid改为手机号
		sqlStr := fmt.Sprintf("exec DD_GetWorkMonth '%s','%s'", mobile, date)
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			//log.Fatal("Prepare failed:", err.Error())
			fmt.Println("Prepare failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2001, "msg": err.Error()})
			return
		}
		defer stmt.Close()

		//通过Statement执行查询
		rows, err := stmt.Query()
		if err != nil {
			//log.Fatal("Query failed:", err.Error())
			fmt.Println("Query failed:", err.Error())
			c.JSON(http.StatusOK, gin.H{"code": -2002, "msg": err.Error()})
			return
		}

		//建立一个列数组
		var respData [][]string
		cols, err := rows.Columns()
		var colsdata = make([]interface{}, len(cols))
		var title []string
		for i := 0; i < len(cols); i++ {
			colsdata[i] = new(interface{})
			title = append(title, cols[i])
		}
		respData = append(respData, title)

		for rows.Next() {
			rows.Scan(colsdata...)
			var tempData []string
			for _, val := range colsdata {
				switch v := (*(val.(*interface{}))).(type) {
				case nil:
					tempData = append(tempData, "")
				case bool:
					if v {
						tempData = append(tempData, "true")
					} else {
						tempData = append(tempData, "false")
					}
				case []byte:
					//fmt.Print(string(v))
					tempData = append(tempData, string(v))
				case time.Time:
					// fmt.Print(v.Format("2016-01-02 15:05:05.999"))
					tempData = append(tempData, v.Format("2016-01-02 15:04:05"))
				default:
					tempData = append(tempData, fmt.Sprintf("%v", v))
				}
			}
			respData = append(respData, tempData)
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": respData})
	}

}

// 根据用户手机号获取userid
func GetUserIdByPhone() {

	sqlStr := fmt.Sprintf("exec DD_GetAllEmpPhoneFromHR")
	data, err := DbQuery(sqlStr, (*models.PhoneStruct)(nil))

	if err != nil {
		Logger.Error("GetUserIdByPhone error:" + err.Error())
		return
	}

	if len(data) == 0 {
		Logger.Info("GetUserIdByPhone data empty")
		return
	}

	for _, v := range data {
		t := v.(*models.PhoneStruct)

		userId, err := getDDUserIdByPhone(t.Phone)
		if err != nil {
			Logger.Error(fmt.Sprintf("getDDUserIdByPhone error:%s, phone:%d", err.Error(), t.Phone))
			continue
		}

		// 更新userid
		sqlStr := fmt.Sprintf("exec DD_UpdateUserIDToHR '%d','%s'", t.Phone, userId)
		_, err = DbExec(sqlStr)
		if err != nil {
			Logger.Error("exec DD_GetAllEmpInfoFromHR error:" + err.Error())
		}
	}

}

func printInterface(val interface{}) string {
	switch v := (*(val.(*interface{}))).(type) {
	case nil:
		//fmt.Print("NULL")
		return ""
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case []byte:
		//fmt.Print(string(v))
		return string(v)
	case time.Time:
		t := v.Format("2016-01-02 15:05:05.999")
		return t
	default:
		return v.(string)
	}
}
