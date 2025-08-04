package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"go-react-admin/model"
	"go-react-admin/service"

	"github.com/gin-gonic/gin"
)

type DynamicFieldApi struct {
	dynamicFieldService service.DynamicFieldService
}

var dynamicFieldService = service.DynamicFieldService{}

// CreateField 创建动态字段
func (api *DynamicFieldApi) CreateField(c *gin.Context) {
	var field model.DynamicField
	if err := c.ShouldBindJSON(&field); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := dynamicFieldService.CreateField(&field); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    field,
	})
}

// GetFieldsByTableID 根据表ID获取字段列表
func (api *DynamicFieldApi) GetFieldsByTableID(c *gin.Context) {
	// 设置缓存控制头，防止浏览器缓存
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	tableID, err := strconv.ParseUint(c.Param("tableId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的表ID",
		})
		return
	}

	fields, err := dynamicFieldService.GetFieldsByTableID(uint(tableID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	fmt.Printf("GetFieldsByTableID tableID: %d, fields: %+v\n", tableID, fields)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    fields,
	})
}

// GetFieldByID 根据ID获取字段
func (api *DynamicFieldApi) GetFieldByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	field, err := dynamicFieldService.GetFieldByID(uint(id))
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
		"data":    field,
	})
}

// UpdateField 更新字段
func (api *DynamicFieldApi) UpdateField(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	var field model.DynamicField
	if err := c.ShouldBindJSON(&field); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	field.ID = uint(id)
	if err := dynamicFieldService.UpdateField(&field); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    field,
	})
}

// DeleteField 删除字段
func (api *DynamicFieldApi) DeleteField(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	if err := dynamicFieldService.DeleteField(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// UpdateFieldOrder 更新字段排序
func (api *DynamicFieldApi) UpdateFieldOrder(c *gin.Context) {
	var req struct {
		FieldIDs []uint `json:"fieldIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := dynamicFieldService.UpdateFieldOrder(req.FieldIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "排序更新成功",
	})
}

// ToggleFieldStatus 切换字段状态
func (api *DynamicFieldApi) ToggleFieldStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	if err := dynamicFieldService.ToggleFieldStatus(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "状态切换成功",
	})
}

// BatchCreateFields 批量创建字段
func (api *DynamicFieldApi) BatchCreateFields(c *gin.Context) {
	var fields []model.DynamicField
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := dynamicFieldService.BatchCreateFields(fields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "批量创建成功",
		"data":    fields,
	})
}

// GetFieldTypes 获取字段类型列表
func (api *DynamicFieldApi) GetFieldTypes(c *gin.Context) {
	fieldTypes := dynamicFieldService.GetFieldTypes()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    fieldTypes,
	})
}
