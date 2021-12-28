package sqlx

import (
	"github.com/ytbiu/errh"
	"gorm.io/gorm"
)

type Base struct {
	GORMHandler
	errh.ErrHandler
}

func New(db *gorm.DB, h errh.ErrHandler) *Base {
	return &Base{
		GORMHandler: GORMHandler{
			db: db,
		},
		ErrHandler: h,
	}
}

// FindAll 查列表
func (b *Base) FindAll(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (total int64) {
	if b.HasErr() {
		return 0
	}
	total, err := b.GORMHandler.FindAll(result, wrapDBs...)
	b.TryToSetErr(err)
	return
}

// FindOne 查详情
func (b *Base) FindOne(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) {
	if b.HasErr() {
		return
	}
	err := b.GORMHandler.FindOne(result, wrapDBs...)
	b.TryToSetErr(err)
}

// FindAllWithoutCount 查列表
func (b *Base) FindAllWithoutCount(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) {
	if b.HasErr() {
		return
	}
	err := b.GORMHandler.FindAllWithoutCount(result, wrapDBs...)
	b.TryToSetErr(err)
}

// Count 数量
func (b *Base) Count(model interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (count int64) {
	if b.HasErr() {
		return
	}
	count, err := b.GORMHandler.Count(model, wrapDBs...)
	b.TryToSetErr(err)
	return count
}

// Create 创建
func (b *Base) Create(model interface{}) {
	if b.HasErr() {
		return
	}
	err := b.GORMHandler.Create(model)
	b.TryToSetErr(err)
}

// Exists 判断是否存在
func (b *Base) Exists(wrapDBs ...func(*gorm.DB) *gorm.DB) (exists bool) {
	if b.HasErr() {
		return
	}
	exists, err := b.GORMHandler.Exists(wrapDBs...)
	b.TryToSetErr(err)
	return
}

// TxProcess 事物
func (b *Base) TxProcess(txFuncGroup ...func(h GORMAndErrHandler)) {
	if b.HasErr() {
		return
	}
	db := b.db
	b.db = db.Begin()
	tx := b.db
	if err := tx.Error; err != nil {
		b.TryToSetErr(err)
		return
	}

	defer func() {
		if b.HasErr() {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, doTx := range txFuncGroup {
		if b.HasErr() {
			return
		}
		doTx(b)
	}
	return
}
