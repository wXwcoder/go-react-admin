package v1

import (
	"net/http"
	"strconv"

	"go-react-admin/model"
	"go-react-admin/service"

	"github.com/gin-gonic/gin"
)

type DynamicTableApi struct{}

var dynamicTableService = service.DynamicTableService{}

// CreateTable 创建动态表
// @Tags DynamicTable
// @Summary 创建动态表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DynamicTable true "动态表信息"
// @Success 200 {object} response.Response{data=model.DynamicTable,msg=string} "创建成功"
// @Router /dynamicTable/createTable [post]
func (dta *DynamicTableApi) CreateTable(c *gin.Context) {
	var table model.DynamicTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err := dynamicTableService.CreateTable(&table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    table,
	})
}

// GetTableList 获取动态表列表
// @Tags DynamicTable
// @Summary 获取动态表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页大小"
// @Param search query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /dynamicTable/getTableList [get]
func (dta *DynamicTableApi) GetTableList(c *gin.Context) {
	// 设置缓存控制头，防止浏览器缓存
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	tables, total, err := dynamicTableService.GetTableList(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"list":     tables,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetTableByID 根据ID获取动态表
// @Tags DynamicTable
// @Summary 根据ID获取动态表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "表ID"
// @Success 200 {object} response.Response{data=model.DynamicTable,msg=string} "获取成功"
// @Router /dynamicTable/getTable/{id} [get]
func (dta *DynamicTableApi) GetTableByID(c *gin.Context) {
	// 设置缓存控制头，防止浏览器缓存
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	table, err := dynamicTableService.GetTableByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    table,
	})
}

// UpdateTable 更新动态表
// @Tags DynamicTable
// @Summary 更新动态表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DynamicTable true "动态表信息"
// @Success 200 {object} response.Response{data=model.DynamicTable,msg=string} "更新成功"
// @Router /dynamicTable/updateTable [put]
func (dta *DynamicTableApi) UpdateTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	var table model.DynamicTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	table.ID = uint(id)
	if err := dynamicTableService.UpdateTable(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    table,
	})
}

// DeleteTable 删除动态表
// @Tags DynamicTable
// @Summary 删除动态表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "表ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /dynamicTable/deleteTable/{id} [delete]
func (dta *DynamicTableApi) DeleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	if err := dynamicTableService.DeleteTable(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// ToggleTableStatus 切换表状态
// @Tags DynamicTable
// @Summary 切换表状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "表ID"
// @Success 200 {object} response.Response{msg=string} "操作成功"
// @Router /dynamicTable/toggleStatus/{id} [patch]
func (dta *DynamicTableApi) ToggleTableStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	if err := dynamicTableService.ToggleTableStatus(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "状态切换成功",
	})
}

// GetTableSchema 获取表结构信息
// @Tags DynamicTable
// @Summary 获取表结构信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param tableName path string true "表名"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /dynamicTable/getSchema/{tableName} [get]
func (dta *DynamicTableApi) GetTableSchema(c *gin.Context) {
	tableName := c.Param("tableName")

	schema, err := dynamicTableService.GetTableSchema(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": schema,
	})
}

// ValidateTableName 验证表名
// @Tags DynamicTable
// @Summary 验证表名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param tableName query string true "表名"
// @Success 200 {object} response.Response{data=bool,msg=string} "验证结果"
// @Router /dynamicTable/validateTableName [get]
func (dta *DynamicTableApi) ValidateTableName(c *gin.Context) {
	tableName := c.Query("tableName")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "表名不能为空",
		})
		return
	}

	isValid, err := dynamicTableService.ValidateTableName(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"isValid": isValid,
		},
	})
}

// CreateDynamicTable 创建动态表
func (dta *DynamicTableApi) CreateDynamicTable(c *gin.Context) {
	var req model.DynamicTable
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 验证表名
	if valid, err := dynamicTableService.ValidateTableName(req.TableName); !valid {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	table, err := dynamicTableService.CreateDynamicTable(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": table})
}
