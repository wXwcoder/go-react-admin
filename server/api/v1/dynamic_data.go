package v1

import (
	"net/http"
	"strconv"

	"go-react-admin/model"
	"go-react-admin/service"

	"github.com/gin-gonic/gin"
)

type DynamicDataApi struct {
	dynamicDataService service.DynamicDataService
}

var dynamicDataService = service.DynamicDataService{}

// CreateData 创建动态数据
func (api *DynamicDataApi) CreateData(c *gin.Context) {
	tableName := c.Param("tableName")

	var data map[string]interface{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err = dynamicDataService.CreateData(tableName, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    data,
	})
}

// GetDynamicDataList 获取动态数据列表
func (api *DynamicDataApi) GetDynamicDataList(c *gin.Context) {
	tableName := c.Param("tableName")
	if tableName == "" {
		c.JSON(400, gin.H{"error": "表名不能为空"})
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	orderBy := c.DefaultQuery("orderBy", "")

	// 获取过滤条件
	filters := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if key != "page" && key != "pageSize" && key != "orderBy" && key != "_t" && len(values) > 0 {
			filters[key] = values[0]
		}
	}

	data, total, err := dynamicDataService.GetDataList(tableName, page, pageSize, filters, orderBy)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"list":     data,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetDataByID 根据ID获取动态数据
func (api *DynamicDataApi) GetDataByID(c *gin.Context) {
	tableName := c.Param("tableName")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	data, err := dynamicDataService.GetDataByID(tableName, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// UpdateData 更新动态数据
func (api *DynamicDataApi) UpdateData(c *gin.Context) {
	tableName := c.Param("tableName")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err = dynamicDataService.UpdateData(tableName, uint(id), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    data,
	})
}

// DeleteData 删除动态数据
func (api *DynamicDataApi) DeleteData(c *gin.Context) {
	tableName := c.Param("tableName")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	if err := dynamicDataService.DeleteData(tableName, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// BatchDeleteData 批量删除动态数据
func (api *DynamicDataApi) BatchDeleteData(c *gin.Context) {
	tableName := c.Param("tableName")

	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := dynamicDataService.BatchDeleteData(tableName, req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量删除成功",
	})
}

// GetDynamicDataStatistics 获取动态数据统计
func (api *DynamicDataApi) GetDynamicDataStatistics(c *gin.Context) {
	tableName := c.Param("tableName")
	if tableName == "" {
		c.JSON(400, gin.H{"error": "表名不能为空"})
		return
	}

	stats, err := dynamicDataService.GetDataStatistics(tableName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": stats})
}

// GetDataStatistics 获取数据统计
func (api *DynamicDataApi) GetDataStatistics(c *gin.Context) {
	tableName := c.Param("tableName")

	statistics, err := dynamicDataService.GetDataStatistics(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": statistics,
	})
}

// CreateView 创建数据视图
func (api *DynamicDataApi) CreateView(c *gin.Context) {
	var view model.DynamicView
	if err := c.ShouldBindJSON(&view); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := dynamicDataService.CreateView(&view); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    view,
	})
}

// GetViewList 获取视图列表
func (api *DynamicDataApi) GetViewList(c *gin.Context) {
	tableID, err := strconv.ParseUint(c.Param("tableId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的表ID",
		})
		return
	}

	views, err := dynamicDataService.GetViewList(uint(tableID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": views,
	})
}

// GetViewByID 获取视图详情
func (api *DynamicDataApi) GetViewByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	view, err := dynamicDataService.GetViewByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": view,
	})
}

// UpdateView 更新视图
func (api *DynamicDataApi) UpdateView(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	var view model.DynamicView
	if err := c.ShouldBindJSON(&view); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	view.ID = uint(id)
	if err := dynamicDataService.UpdateView(&view); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    view,
	})
}

// DeleteView 删除视图
func (api *DynamicDataApi) DeleteView(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	if err := dynamicDataService.DeleteView(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// ApplyView 应用视图
func (api *DynamicDataApi) ApplyView(c *gin.Context) {
	viewID, err := strconv.ParseUint(c.Param("viewId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的视图ID",
		})
		return
	}

	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := dynamicDataService.ApplyView(uint(viewID), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
