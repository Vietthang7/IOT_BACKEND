package consts

import (
	"math"

	"github.com/google/uuid"
)

// query and pagination

type Pagination struct {
	CurrentPage  int     `json:"current_page"`
	TotalPages   float64 `json:"total_pages"`
	TotalResults int64   `json:"total"`
}

func (p *Pagination) GetTotalPages(len int) float64 {
	return math.Ceil(float64(p.TotalResults) / float64(len))
}

type Query struct {
	ID        uuid.UUID `json:"id"`
	Search    string    `json:"search"`
	Active    string    `query:"active"`
	Page      int       `json:"page"`
	Length    int       `json:"length" form:"length"`
	Order     string    `json:"order"`
	Sort      string    `json:"sort"`
	StartTime string    `query:"start_time"`
	EndTime   string    `query:"end_time"`
	Code      string    `query:"code"`
	Type      uint      `query:"type"`
	Status    int       `query:"status"`
	Result    int       `query:"result"`
}

// func (q *Query) GetActive() *bool {
// 	if q.Active != "" {
// 		if status, err := strconv.ParseBool(q.Active); err == nil {
// 			return &status
// 		}
// 	}
// 	return nil
// }
// func (q *Query) GetOffset() int {
// 	return (q.GetPage() - 1) * q.GetPageSize()
// }

// func (q *Query) GetPageSize() int {
// 	if q.Length > 200 {
// 		q.Length = 200
// 	}
// 	if q.Length < 1 {
// 		q.Length = 12
// 	}
// 	return q.Length
// }

// func (q *Query) GetPage() int {
// 	if q.Page < 1 {
// 		q.Page = 1
// 	}
// 	return q.Page
// }

// func (q *Query) GetField(Orders []string, d string) string {
// 	if !utils.Contains(Orders, q.Order) {
// 		q.Order = d
// 	}
// 	return q.Order
// }

// func (q *Query) GetSort() string {
// 	if q.Sort != "asc" {
// 		q.Sort = "desc"
// 	}
// 	return q.Sort
// }

// func (q *Query) GetStartDate(fmt string) string {
// 	t, err := utils.ConvertStringToTime(q.StartAt)
// 	log.Println(t)
// 	if err != nil {
// 		return time.Now().Format(fmt)
// 	}
// 	return t.Format(fmt)
// }

// func (q *Query) GetEndDate(fmt string) string {
// 	t, err := utils.ConvertStringToTime(q.EndAt)
// 	if err != nil {
// 		return time.Now().Format(fmt)
// 	}
// 	return t.Format(fmt)
// }
