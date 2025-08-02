package router

import (
	v1 "go-react-admin/api/v1"

	"github.com/gin-gonic/gin"
)

// InitDynamicRoutes 初始化动态数据管理路由
func InitDynamicRoutes(Router *gin.RouterGroup) {
	dynamicTableApi := v1.DynamicTableApi{}
	dynamicFieldApi := v1.DynamicFieldApi{}
	dynamicDataApi := v1.DynamicDataApi{}

	// 动态表管理路由
	dynamicTableRouter := Router.Group("dynamicTable")
	{
		dynamicTableRouter.POST("createTable", dynamicTableApi.CreateTable)             // 创建动态表
		dynamicTableRouter.GET("getTableList", dynamicTableApi.GetTableList)            // 获取动态表列表
		dynamicTableRouter.GET("getTable/:id", dynamicTableApi.GetTableByID)            // 根据ID获取动态表
		dynamicTableRouter.PUT("updateTable/:id", dynamicTableApi.UpdateTable)              // 更新动态表
		dynamicTableRouter.DELETE("deleteTable/:id", dynamicTableApi.DeleteTable)       // 删除动态表
		dynamicTableRouter.PATCH("toggleStatus/:id", dynamicTableApi.ToggleTableStatus) // 切换表状态
		dynamicTableRouter.GET("getSchema/:tableName", dynamicTableApi.GetTableSchema)  // 获取表结构信息
		dynamicTableRouter.GET("validateTableName", dynamicTableApi.ValidateTableName)  // 验证表名
	}

	// 动态字段管理路由
	dynamicFieldRouter := Router.Group("dynamicField")
	{
		dynamicFieldRouter.POST("createField", dynamicFieldApi.CreateField)                // 创建动态字段
		dynamicFieldRouter.GET("getFields/:tableId", dynamicFieldApi.GetFieldsByTableID)   // 根据表ID获取字段列表
		dynamicFieldRouter.GET("getField/:id", dynamicFieldApi.GetFieldByID)               // 根据ID获取字段
		dynamicFieldRouter.PUT("updateField/:id", dynamicFieldApi.UpdateField)                 // 更新字段
		dynamicFieldRouter.DELETE("deleteField/:id", dynamicFieldApi.DeleteField)          // 删除字段
		dynamicFieldRouter.PATCH("updateOrder/:tableId", dynamicFieldApi.UpdateFieldOrder) // 更新字段排序
		dynamicFieldRouter.PATCH("toggleStatus/:id", dynamicFieldApi.ToggleFieldStatus)    // 切换字段状态
		dynamicFieldRouter.POST("batchCreate", dynamicFieldApi.BatchCreateFields)          // 批量创建字段
		dynamicFieldRouter.GET("getFieldTypes", dynamicFieldApi.GetFieldTypes)             // 获取支持的字段类型
	}

	// 动态数据管理路由
	dynamicDataRouter := Router.Group("dynamicData")
	{
		dynamicDataRouter.POST(":tableName/create", dynamicDataApi.CreateData)             // 创建动态数据
		dynamicDataRouter.GET(":tableName/list", dynamicDataApi.GetDynamicDataList)        // 获取动态数据列表
		dynamicDataRouter.GET(":tableName/get/:id", dynamicDataApi.GetDataByID)            // 根据ID获取动态数据
		dynamicDataRouter.PUT(":tableName/update/:id", dynamicDataApi.UpdateData)          // 更新动态数据
		dynamicDataRouter.DELETE(":tableName/delete/:id", dynamicDataApi.DeleteData)       // 删除动态数据
		dynamicDataRouter.DELETE(":tableName/batchDelete", dynamicDataApi.BatchDeleteData) // 批量删除动态数据
		dynamicDataRouter.GET(":tableName/statistics", dynamicDataApi.GetDataStatistics)  // 获取数据统计
	}

	// 动态视图管理路由
	dynamicViewRouter := Router.Group("dynamicView")
	{
		dynamicViewRouter.POST("create", dynamicDataApi.CreateView)              // 创建数据视图
		dynamicViewRouter.GET("list/:tableId", dynamicDataApi.GetViewList)       // 获取视图列表
		dynamicViewRouter.GET("get/:id", dynamicDataApi.GetViewByID)             // 获取视图详情
		dynamicViewRouter.PUT("update/:id", dynamicDataApi.UpdateView)             // 更新视图
		dynamicViewRouter.DELETE("delete/:id", dynamicDataApi.DeleteView)        // 删除视图
		dynamicViewRouter.POST("apply/:viewId", dynamicDataApi.ApplyView)        // 应用视图
	}
}
