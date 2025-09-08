package repo

import (
	"backend/app"
	"backend/internal/consts"
	"backend/internal/model"
	"context"
)

type (
	DataSensor      model.DataSensor
	List_DataSensor []DataSensor
)

func (u *DataSensor) Find(p *consts.RequestTable, query interface{}, args []interface{}) (entries List_DataSensor, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
		DB          = p.CustomOptions(app.Database.DB).WithContext(ctx).Where(query, args...)
	)
	defer cancel()
	err = DB.Debug().Find(&entries).Error
	return
}

func (u *DataSensor) Count(query interface{}, args []interface{}) int64 {
	var (
		count       int64
		ctx, cancel = context.WithTimeout(context.Background(), app.CTimeOut)
	)
	defer cancel()
	app.Database.DB.Where(query, args...).Model(&model.DataSensor{}).WithContext(ctx).Count(&count)
	return count
}
