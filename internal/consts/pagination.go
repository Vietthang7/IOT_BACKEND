package consts

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RequestTable struct {
	Page      int    `json:"page,omitempty" form:"page"`
	Length    int    `json:"length,omitempty" form:"length"`
	Search    string `json:"search,omitempty" form:"search"`
	Order     string `json:"order_by,omitempty" form:"order_by"`
	Dir       string `json:"order_dir,omitempty" form:"order_dir"`
	DirNumber int    `json:"-"`
	Total     int64  `json:"total,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
}

func BindRequestTable(c *fiber.Ctx, order string) RequestTable {
	var request RequestTable
	_ = c.QueryParser(&request)

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Length <= 0 {
		request.Length = 10
	}
	if request.Search != "" {
		request.Search = strings.TrimSpace(request.Search)
	}
	if request.Length < 5 {
		request.Length = 5
	}
	if request.Order == "" {
		request.Order = order
	}
	if strings.ToLower(request.Dir) != "asc" {
		request.DirNumber = -1
		request.Dir = "desc"
	} else {
		request.DirNumber = 1
		request.Dir = "asc"
	}

	return request
}

func (u *RequestTable) CustomOptions(DB *gorm.DB) *gorm.DB {
	if u.Order == "" {
		return DB.Limit(u.Length).Offset((u.Page - 1) * u.Length)
	}
	return DB.Limit(u.Length).Offset((u.Page - 1) * u.Length).Order(fmt.Sprintf("%s %s", u.Order, u.Dir))
}
