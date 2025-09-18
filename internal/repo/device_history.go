package repo

import (
	"backend/app"
	"backend/internal/consts"
	"backend/internal/model"
	"context"
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
