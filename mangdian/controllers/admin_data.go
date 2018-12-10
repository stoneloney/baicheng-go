package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"common/models"
	"common/method"

	"mangdian/models"

	"github.com/tealeg/xlsx"
	"github.com/gin-gonic/gin"
)

var DataMap = map[string]map[string]string {
	"type": {
		"0": "--",
		"1": "企业宣传片",
		"2": "产品宣传片",
		"3": "广告片",
		"4": "微电影",
	},
	"duration": {
		"0": "--",
		"1": "30秒",
		"2": "0-1分钟",
		"3": "1-3分钟",
		"4": "3-5分钟",
		"5": "5-10分钟",
		"6": "10-15分钟",
		"7": "30分钟",
	},
	"director": {
		"0": "否",
		"1": "否",
		"2": "有",
	},
	"model": {
		"0": "--",
		"1": "国内",
		"2": "国外",
	},
	"effect": {
		"0": "--",
		"1": "三维",
		"2": "MG",
		"3": "特效",
	},
	"dubbed": {
		"0": "--",
		"1": "中文",
		"2": "英文",
		"3": "中英文",
	},
}

// 询价
func AdminData(c *gin.Context) {
	H := DefaultH(c)
	H["Title"] = "Admin data"

	page, _ := strconv.Atoi(c.Query("page"))
	pageNum := 20
	count := 0

	db := models.GetDB()

	var dataes []projmodels.Data
	// 总数
	db.Find(&dataes).Count(&count)
	H["count"] = count
	H["pagenum"] = pageNum

	// 获取数据列表
	db.Limit(pageNum).Offset(page*pageNum).Order("createtime desc").Find(&dataes)

	for index, val := range dataes {
		dataes[index].Type = DataMap["type"][val.Type]
		dataes[index].Duration = DataMap["duration"][val.Duration]
		dataes[index].Director = DataMap["director"][val.Director]
		dataes[index].Model = DataMap["model"][val.Model]
		dataes[index].Effect = DataMap["effect"][val.Effect]
		dataes[index].Dubbed = DataMap["dubbed"][val.Dubbed]
	}

	H["dataes"] = dataes

	c.HTML(http.StatusOK, "admin/data", H)
}

// 订制
func AdminMade(c *gin.Context) {
	H := DefaultH(c)
	H["title"] = "Admin made"

	page, _ := strconv.Atoi(c.Query("page"))
	pageNum := 20
	count := 0

	db := models.GetDB()

	var dataes []projmodels.Made
	// 总数
	db.Find(&dataes).Count(&count)
	H["count"] = count
	H["pagenum"] = pageNum

	// 获取数据列表
	db.Limit(pageNum).Offset(page*pageNum).Order("createtime desc").Find(&dataes)

	for index, val := range dataes {
		dataes[index].Type = DataMap["type"][val.Type]
		dataes[index].Duration = DataMap["duration"][val.Duration]
	}
	H["dataes"] = dataes

	c.HTML(http.StatusOK, "admin/made", H)

}

// 删除询价
func AdminDataDelete(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db := models.GetDB()
	data := projmodels.Data{}
	db.First(&data, id)
	if data.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"数据不存在"})
		return
	}
	if err := db.Delete(&data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})

}

// 删除订制
func AdminMadeDelete(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db := models.GetDB()
	data := projmodels.Made{}
	db.First(&data, id)
	if data.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"数据不存在"})
		return
	}
	if err := db.Delete(&data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

// 数据导出
func ExportData(c *gin.Context) {
	mtype := c.Query("type")
	var types = []string{"data", "made"}
	if (!method.InArray(mtype, types)) {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"类型错误"})
		return
	}

	db := models.GetDB()
	var dataes []projmodels.Data
	var madeDataes []projmodels.Made
	var dataTitles = [9]string{"ID","电话","类型","时长","导演","模特","特效","配音","日期"}
	var madeTitles = [7]string{"ID","电话","类型","时长","城市","公司","时间"}

	if mtype == "data" {
		db.Order("createtime desc").Find(&dataes)
	} else {
		db.Order("createtime desc").Find(&madeDataes)
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"code":2, "msg":err.Error()})
        return
    }
    row = sheet.AddRow()
    if mtype == "data" {
    	for _, title := range dataTitles {
	    	cell = row.AddCell()
	    	cell.Value = title
	    }
	    for _, val := range dataes {
	    	values := []string{
	    		strconv.FormatInt(val.ID, 10),
	    		val.Phone,
	    		val.Type,
	    		val.Duration,
	    		val.Director,
	    		val.Model,
	    		val.Effect,
	    		val.Dubbed,
	    		val.Createtime,
	    	}
	    	row = sheet.AddRow()
	    	for _, v := range values {
	    		cell = row.AddCell()
	    		cell.Value = v
	    	}
	    }
    } else {
    	for _, title := range madeTitles {
    		cell = row.AddCell()
    		cell.Value = title
    	}
    	for _, val := range madeDataes {
	    	values := []string{
	    		strconv.FormatInt(val.ID, 10),
	    		val.Phone,
	    		val.Type,
	    		val.Duration,
	    		val.City,
	    		val.Company,
	    		val.Createtime,
	    	}
	    	row = sheet.AddRow()
	    	for _, v := range values {
	    		cell = row.AddCell()
	    		cell.Value = v
	    	}
	    }
    }
    // 输出内容
    /*
    this.Ctx.Output.Header("Accept-Ranges", "bytes")
    this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "orders.xls"))//文件名
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	file.Write(c.ResponseWriter)
	*/
	writer := c.Writer
	writer.Header().Set("Accept-Ranges", "bytes")
	writer.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "data.xls"))
	writer.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	file.Write(writer)
}





