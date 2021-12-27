package gsql

import "gorm.io/gorm"

// Base baseDao
type Base struct {
	GORMHandler
}

// New 构造方法
func New(db *gorm.DB, errWrappers ...func(err error) error) *Base {
	return &Base{
		GORMHandler: GORMHandler{
			db: db,
		},
		ErrHandler: ErrHandler{
			defaultErrWrappers: errWrappers,
		},
	}
}

// FindAll 查列表
func (b *Base) FindAll(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (total int64) {
	if b.HasErr() {
		return 0
	}
	total, b.err = FindAll(result, wrapDBs...)
	return
}

// FindOne 查详情
func (b *Base) FindOne(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) {
	if b.HasErr() {
		return
	}
	b.err = FindOne(result, wrapDBs...)
}

// FindAllWithoutCount 查列表
func (b *Base) FindAllWithoutCount(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) {
	if b.HasErr() {
		return
	}
	b.err = FindAllWithoutCount(result, wrapDBs...)
}

// Count 数量
func (b *Base) Count(model interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (count int64) {
	if b.HasErr() {
		return
	}
	count, b.err = Count(model, wrapDBs...)
	return count
}

// Create 创建
func (b *Base) Create(model interface{}) {
	if b.HasErr() {
		return
	}
	b.err = Create(model)
}

// Exists 判断是否存在
func (b *Base) Exists(wrapDBs ...func(*gorm.DB) *gorm.DB) (exists bool) {
	if b.HasErr() {
		return
	}
	exists, b.err = Exists(wrapDBs...)
	return
}

// TxProcess 事物
func (b *Base) TxProcess(txFuncGroup ...func(h GORMAndErrHandler)) {
	if b.HasErr() {
		return
	}
	db := GetClient()
	b.db = db.Begin()
	tx := b.db
	if b.err = tx.Error; b.err != nil {
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
