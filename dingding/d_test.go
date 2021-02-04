package main

import (
	"database/sql"
	"dingding/controllers"
	"dingding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/mattn/go-adodb"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
//appKey   = "dingebpo4jkck8mkkbeb"   // 9duan
//appSecrt = "BgoAD74YRu-lxS-US4NVNXMRcuQ5f9nH5FsaJEwDe7CD7PG5oZMwm_OVTw9c30T-"  // 9duan

//appKey       = "dingt19d59lo94bksnx1"
//appSecrt     = "00cTPglw0NRNM5knpiA0rhtoFgd0ZOuWMS8mv6COLdCBAhTkM4XRhODrcttF6M_u"
)

// 考勤
func TestAttendanceList(t *testing.T) {
	userIds := []string{
		//"2640262323937802",  // 9段账号
		// "manager2529",       // 9段账号
		//"171802165236043309",
		//"manager1609",
		"654719055736498989",
	}
	// 时间戳模版
	timeLayout := "2006-01-02 15:04:05"
	endUnix := time.Now().Unix()
	startUnix := endUnix - 3600*24
	endDate := time.Unix(endUnix, 0).Format(timeLayout)
	startDate := time.Unix(startUnix, 0).Format(timeLayout)

	startDate = "2020-04-02 00:00:00"
	endDate = "2020-04-04 23:59:59"

	code, msg, list := controllers.AttendanceRecordList(userIds, startDate, endDate)
	if code != 0 {
		fmt.Println("attendlist error:", msg)
		return
	}
	fmt.Println(list)
}

// 请假
func TestAttendLeave(t *testing.T) {
	userIds := []string{
		//"2640262323937802",
		//"manager2529",
		"manager7998", // baicheng
	}
	// 时间戳模版
	// timeLayout := "2006-01-02 15:04:05"
	endUnix := time.Now().Unix()*1000 + 24*3600*1000
	startUnix := endUnix - 180*3600*24*1000

	//startUnix := int64(1556467200*1000)
	//endUnix := startUnix + 24*3600*1000

	//endDate := time.Unix(endUnix, 0).Format(timeLayout)

	//startDate := time.Unix(startUnix, 0).Format(timeLayout)

	code, list, hasmore, err := controllers.AttendanceLeaveStatus(userIds, startUnix, endUnix, 0, 20)
	if code != 0 {
		fmt.Println("attend leave error:", err.Error())
	}
	fmt.Println(hasmore)
	fmt.Println(list)
}

// 注册回调事件
func TestRegisterCallback(t *testing.T) {
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
	_, err := controllers.RegisterCallback(tags)
	if err != nil { // err != nil
		fmt.Println("reg:", err.Error())
	} else {
		// 返回回调成功
		token := "123456"
		aesKey := "1234567890123456789012345678901234567890123"
		//corpid := "ding4cb600184df08bd035c2f4657eb6378f"
		corpid := "ding02634a50287cf0f835c2f4657eb6378f"
		cropty := controllers.NewCrypto(token, aesKey, corpid)

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

		data, _ := json.Marshal(resp)

		fmt.Println(string(data))
	}
}

// 更新注册回调
func TestUpdateCallback(t *testing.T) {
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
	_, err := controllers.UpdateCallback(tags)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("update success")
	}
}

// 注册回调列表
func TestCallbackList(t *testing.T) {
	code, data, err := controllers.GetCallbackList()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(code)
		fmt.Println(data)
	}
}

// oa信息推送
func TestPushOaMsg(t *testing.T) {
	//userId := "171802165236043309"
	/*
	userId := "033435111821927863"
	_, err := controllers.PushMsg(userId, )
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("success")
	}
	*/
}

// text信息推送
func TestPushTextMsg(t *testing.T) {
	/*
	note := []string{
		"工资1:1000",
		"工资1:2000",
		"工资1:3000",
		"工资1:4000",
		"工资1:5000",
		"工资1:6000",
		"工资1:7000",
		"工资1:8000",
		"工资1:9000",
		"工资1:10000",
		"工资1:11000",
	}
	userId := "171802165236043309"

	token, err := controllers.getAccessToken()
	if err != nil {
		fmt.Println(fmt.Sprintf("getAccessToken error:%s", err.Error()))
		return
	}

	_, err = controllers.PushTextMsg(token, userId, note, )
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("success")
	}
	 */
}

// 推送信息列表
func TestPushMsgList(t *testing.T) {
	controllers.PushCash()
}

// 解密
func TestDecryptMsg(t *testing.T) {
	token := "123456"
	aesKey := "1234567890123456789012345678901234567890123"
	corpid := "ding4cb600184df08bd035c2f4657eb6378f"
	cropty := controllers.NewCrypto(token, aesKey, corpid)

	// 钉钉返回数据
	sign := "17455f6cec375012fcf5974e98d548fdaf480193"
	timestamp := "1561177082241"
	nonce := "ul6dU1CF"
	secret := "Q78g5BQPh4vGX1xY1YRdlZxpM8YJn3n0At+8XyGR/LJ4FOwr8DUrFykxRrkqSa+DpHaEWnXAlsY1OQ3bxaoqDTQNGi2J06cots2lfyCLhSk3H52E9y8HuXCCbI6VRNA+"
	data, err := cropty.DecryptMsg(sign, timestamp, nonce, secret)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(data)
	}
}

// 加密
func TestEncryptMsg(t *testing.T) {
	token := "123456"
	aesKey := "1234567890123456789012345678901234567890123"
	corpid := "ding4cb600184df08bd035c2f4657eb6378f"
	cropty := controllers.NewCrypto(token, aesKey, corpid)

	replymsg := "success"
	timestamp := "1783610513"
	nonce := "123456"
	encrypt, sign, err := cropty.EncryptMsg(replymsg, timestamp, nonce)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(encrypt)
		fmt.Println(sign)
	}
}

// 员工数据写入
func TestEmpInsert(t *testing.T) {
	name := "钉钉人员1"
	userid := "DD_01"
	department := []int64{112353960}
	position := "test"
	mobile := "13488875678"
	tel := "12345678"
	workPlace := "深圳市"
	remark := "今天下雨"
	email := "123456@qq.com"
	orgEmail := ""
	jobnumber := "123"
	isHide := false
	isSenior := false
	extattrMap := map[string]string{
		"爱好": "睡觉",
		"电影": "复仇者",
	}
	extattr, _ := json.Marshal(extattrMap)
	hireDate := time.Now().Unix()
	op := "create"
	uid, err := controllers.UserOp(op, userid, name, department, position, mobile, tel, workPlace, remark, email, orgEmail, jobnumber, isHide, isSenior, string(extattr), hireDate)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(uid)
}

// 获取员工的具体数据
func TestEmpInfo(t *testing.T) {
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
	eid := "DD_01"
	sqlStr := fmt.Sprintf("exec DD_GetEmpInfoFromHR '%s'", eid)
	fmt.Println(sqlStr)
	rows, err := conn.Query(sqlStr)
	if err != nil {
		fmt.Println("sql exec error:", err.Error())
	}
	/*
		fmt.Println(rows)
		rowCnt, err := rows.RowsAffected()
		if err != nil {
			fmt.Println("sql exec error:", err.Error())
		}
		fmt.Println(rowCnt)
	*/
	//建立一个列数组
	// 输出列表
	cols, err := rows.Columns()
	var colsdata = make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		colsdata[i] = new(interface{})
		fmt.Print(cols[i])
		fmt.Print("\t")
	}
	// 输出数据
	for rows.Next() {
		rows.Scan(colsdata...) //将查到的数据写入到这行中
		PrintRow(colsdata)     //打印此行
	}
}

// 获取新增用户
func TestNewEmp(t *testing.T) {
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

	stmt, err := conn.Prepare(`select * from [DD_empNew]`)
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
	// 输出列表
	cols, err := rows.Columns()
	var colsdata = make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		colsdata[i] = new(interface{})
		fmt.Print(cols[i])
		fmt.Print("\t")
	}
	// 输出数据
	for rows.Next() {
		rows.Scan(colsdata...) //将查到的数据写入到这行中
		PrintRow(colsdata)     //打印此行
	}

}

// 获取新增部门 (查询各种数据)
func TestNewDep(t *testing.T) {
	var isdebug = true

	/*
		var server = "61.142.247.115"
		var port = 23284
		var user = "sa"
		var password = "gt2019"
		var database = "GTKQ"
	*/

	var server = "61.142.247.115"
	var port = 23284
	var user = "sa"
	var password = "123"
	var database = "tophr"

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", server, user, password, port, database)
	if isdebug {
		fmt.Println(connString)
	}
	//建立连接
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		//log.Fatal("Open Connection failed:", err.Error())
		fmt.Println("Open Connection failed:", err.Error())
	}
	err = conn.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer conn.Close()

	//stmt, err := conn.Prepare(`select * from [DD_DepNew]`)   // 部门新增
	//stmt, err := conn.Prepare(`select * from [employee]`) // 用户列表
	//stmt, err := conn.Prepare(`select * from [DD_DepEdit]`)   // 部门修改
	//stmt, err := conn.Prepare(`select * from [DD_DepDel]`) // 部门删除
	//stmt, err := conn.Prepare(`select * from [DD_DepConvert]`) // 部门映射
	//stmt, err := conn.Prepare(`select * from [DD_EmpNew]`)  // 新用户
	//stmt, err := conn.Prepare(`select * from [DD_EmpEdit]`)  // 修改用户
	//stmt, err := conn.Prepare(`select * from [DD_EmpDel]`)  // 删除用户
	//stmt, err := conn.Prepare(`select * from [department]`)  // 部门表
	stmt, err := conn.Prepare(`exec DD_GetAllEmpInfoFromHR`) // 获取所有用户
	//stmt, err := conn.Prepare(`exec DD_GetWorkError '001','2019-02'`)  // 考勤异常
	//stmt, err := conn.Prepare(`exec DD_GetWorkDetail '001','2019-02'`)  // 考勤日报
	//stmt, err := conn.Prepare(`exec DD_GetSalaryDetail '001','2019-07'`)  // 工资明细
	//stmt, err := conn.Prepare(`exec DD_GetWorkMonth '033435111821927863','2019-07'`)  // 考勤月报
	//stmt, err := conn.Prepare(`select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE UpDepCode='07'`)

	// 找出没有上级需要生成的最大值
	//stmt, err := conn.Prepare(`select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE len(UpDepCode)=0`)
	//stmt, err := conn.Prepare(`select MAX(DEPARTMENTCODE) AS maxDepCode from [department] WHERE len(DEPARTMENTCODE)=2`)
	// 如果有上级找出需要生成的最大值
	//stmt, err := conn.Prepare(`select MAX(DepCode) AS maxDepCode from [DD_DepNew] WHERE UpDepCode='06'`)
	//hrid := "06"
	//sqlStr := "select MAX(DEPARTMENTCODE) AS maxDepCode from [department] WHERE DEPARTMENTCODE LIKE '"+ hrid +"%'"
	//fmt.Println(sqlStr)
	//stmt, err := conn.Prepare(sqlStr)

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
		fmt.Print(cols[i])
		fmt.Print("\t")
	}

	for rows.Next() {
		rows.Scan(colsdata...) //将查到的数据写入到这行中
		PrintRow(colsdata)     //打印此行
		//fmt.Println(string(colsdata[0].(string)))
		//PrintInterface(colsdata[0])
		//fmt.Println(colsdata[0].(*string))
	}
}

func TestNewDep2(t *testing.T) {
	/*
		data, err := controllers.GetDepNew()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			a := data[0].(*models.DepNewStruct)
			fmt.Println(a.DepName)
		}
	*/
}

type SA struct {
	user   string
	passwd string
	port   int
}

type Mssql struct {
	*sql.DB
	dataSource string
	database   string
	windows    bool
	sa         SA
}

func (m *Mssql) Open() (err error) {
	var conf []string
	conf = append(conf, "Provider=SQLOLEDB")
	conf = append(conf, "Data Source="+m.dataSource)
	if m.windows {
		// Integrated Security=SSPI 这个表示以当前WINDOWS系统用户身去登录SQL SERVER服务器(需要在安装sqlserver时候设置)，
		// 如果SQL SERVER服务器不支持这种方式登录时，就会出错。
		conf = append(conf, "integrated security=SSPI")
	}
	conf = append(conf, "Initial Catalog="+m.database)
	conf = append(conf, "user id="+m.sa.user)
	conf = append(conf, "password="+m.sa.passwd)
	conf = append(conf, "port="+fmt.Sprint(m.sa.port))

	m.DB, err = sql.Open("adodb", strings.Join(conf, ";"))
	if err != nil {
		return err
	}
	return nil
}

func TestConnect2(t *testing.T) {
	db, err := gorm.Open("mssql", "sqlserver://sa:gt2019@61.142.247.115:23284?database=KQ")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
}

func TestConnect(t *testing.T) {

	db := Mssql{
		dataSource: "61.142.247.115\\SQLEXPRESS",
		database:   "KQ",
		// windwos: true 为windows身份验证，false 必须设置sa账号和密码
		windows: false,
		sa: SA{
			user:   "sa",
			passwd: "gt2019",
			port:   23284,
		},
	}
	// 连接数据库
	err := db.Open()
	if err != nil {
		fmt.Println("sql open:", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from [employee]")
	if err != nil {
		fmt.Println("query error : ", err.Error())
		return
	}
	for rows.Next() {
		var name string
		var number int
		rows.Scan(&name, &number)
		fmt.Printf("Name: %s \t Number: %d\n", name, number)
	}
}

func TestConfig(t *testing.T) {
	fmt.Println("appkey:" + controllers.Cfg.Dingding.Appkey)
}

// 查询钉钉在职人员
func TestJob(t *testing.T) {
	err := controllers.RsyncRecord2()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func PrintInterface(val interface{}) string {
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

func PrintRow(colsdata []interface{}) {
	for _, val := range colsdata {
		switch v := (*(val.(*interface{}))).(type) {
		case nil:
			fmt.Print("NULL")
		case bool:
			if v {
				fmt.Print("True")
			} else {
				fmt.Print("False")
			}
		case []byte:
			fmt.Print(string(v))
		case time.Time:
			fmt.Print(v.Format("2016-01-02 15:05:05.999"))
		default:
			fmt.Print(v)
		}
		fmt.Print("\t")
	}
}

func init() {
	// 加载配置
	controllers.LoadConfig()
	controllers.SetAppKey(controllers.Cfg.Dingding.Appkey)
	controllers.SetSecret(controllers.Cfg.Dingding.Appsecret)
}
