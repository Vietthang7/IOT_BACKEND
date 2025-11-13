package repo

import (
	"backend/app"
	"backend/internal/consts"
	"backend/internal/model"
	"context"
	"time"
)

type (
	DeviceHistory      model.DeviceHistory
	List_DeviceHistory []DeviceHistory
)

func (u *DeviceHistory) Find(p *consts.RequestTable, query interface{}, args []interface{}) (entries List_DeviceHistory, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB).WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	err = DB.Debug().Find(&entries).Error
	return
}

func (u *DeviceHistory) Count(query interface{}, args []interface{}) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Where(query, args...).Model(&model.DeviceHistory{}).WithContext(ctx).Count(&count)
	return count
}
func (d *DeviceHistory) GetDevicesWithLatestStatus() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := app.Database.DB.Raw(`
        SELECT 
            dh1.device_name,
            dh1.action,
            dh1.time as last_updated
        FROM device_histories dh1
        INNER JOIN (
            SELECT device_name, MAX(time) as max_time
            FROM device_histories
            WHERE deleted_at IS NULL
            GROUP BY device_name
        ) dh2 ON dh1.device_name = dh2.device_name AND dh1.time = dh2.max_time
        WHERE dh1.deleted_at IS NULL
        ORDER BY dh1.device_name
    `).Scan(&results).Error

	return results, err
}
func (d *DeviceHistory) Create() error {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()

	err := app.Database.DB.WithContext(ctx).Create(d).Error
	return err
}
func (d *DeviceHistory) GetDeviceStatsByDateRange(startDate, endDate string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := app.Database.DB.Raw(`
        SELECT 
            DATE(time) as date,
            device_name,
            SUM(CASE WHEN action = 'ON' THEN 1 ELSE 0 END) as on_count,
            SUM(CASE WHEN action = 'OFF' THEN 1 ELSE 0 END) as off_count,
            COUNT(*) as total_actions
        FROM device_histories
        WHERE DATE(time) BETWEEN ? AND ? AND deleted_at IS NULL
        GROUP BY DATE(time), device_name
        ORDER BY DATE(time) DESC, device_name
    `, startDate, endDate).Scan(&results).Error

	return results, err
}
func (d *DeviceHistory) GetTodayDeviceStats() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	today := time.Now().Format("2006-01-02")

	err := app.Database.DB.Raw(`
        SELECT 
            device_name,
            SUM(CASE WHEN action = 'ON' THEN 1 ELSE 0 END) as on_count,
            SUM(CASE WHEN action = 'OFF' THEN 1 ELSE 0 END) as off_count,
            COUNT(*) as total_actions
        FROM device_histories
        WHERE DATE(time) = ? AND deleted_at IS NULL
        GROUP BY device_name
        ORDER BY device_name
    `, today).Scan(&results).Error

	return results, err
}
