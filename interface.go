package sqlx

import "gorm.io/gorm"

type GORMAndErrHandler interface {
	FindAll(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (total int64)
	FindOne(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB)
	FindAllWithoutCount(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB)
	Count(model interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (count int64)
	Create(model interface{})
	Exists(wrapDBs ...func(*gorm.DB) *gorm.DB) (exists bool)
	TxProcess(txFuncGroup ...func(h GORMAndErrHandler))
	Err(errWrappers ...func(err error) error) error
}
