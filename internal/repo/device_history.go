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
func (d *DeviceHistory) GetDistinctDeviceNames() ([]string, error) {
	var deviceNames []string

	err := app.Database.DB.Model(&model.DeviceHistory{}).
		Distinct("device_name").
		Pluck("device_name", &deviceNames).Error

	return deviceNames, err
}
