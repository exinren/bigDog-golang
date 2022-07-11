package app

import (
	"bigDog-golang/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}
// 构造一个返回的结构体
func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}
// 返回数据
func (r *Response) ToResponse(data interface{})  {
	if nil == data{
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}
// 返回列表
//func (r Response) ToResponseList(list interface{}, totalRows int)  {
//	r.Ctx.JSON(http.StatusOK, gin.H{
//		"list": list,
//		"pager": Pager{
//			Page: GetPage(r.Ctx),
//			PageSize: GetPageSize(r.Ctx),
//			TotalRows: totalRows,
//		},
//	})
//}

// 错误返回
func (r Response) ToErrorResponse(err *errcode.Error)  {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}