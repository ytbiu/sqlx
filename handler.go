package sqlx

import "gorm.io/gorm"

type GORMHandler struct {
	db *gorm.DB
}

func apply(queryDB *gorm.DB, wrapDBs ...func(*gorm.DB) *gorm.DB) {
	for _, wrapDB := range wrapDBs {
		queryDB = wrapDB(queryDB)
	}
}

// FindAll 查列表
func (h *GORMHandler) FindAll(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (int64, error) {
	queryDB := h.db.Model(result)
	apply(queryDB, wrapDBs...)
	if err := queryDB.Find(result).Error; err != nil {
		return 0, err
	}

	var total int64
	return total, queryDB.Count(&total).Error
}

// FindOne 查详情
func (h *GORMHandler) FindOne(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) error {
	queryDB := h.db.Model(result)
	apply(queryDB, wrapDBs...)
	return queryDB.First(result).Error
}

// FindAllWithoutCount 查列表
func (h *GORMHandler) FindAllWithoutCount(result interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) error {
	queryDB := h.db.Model(result)
	apply(queryDB, wrapDBs...)
	return queryDB.Find(result).Error
}

// Count 数量
func (h *GORMHandler) Count(model interface{}, wrapDBs ...func(*gorm.DB) *gorm.DB) (int64, error) {
	queryDB := h.db.Model(model)
	var total int64
	apply(queryDB, wrapDBs...)
	return total, queryDB.Count(&total).Error
}

// Create 创建
func (h *GORMHandler) Create(model interface{}) error {
	return h.db.Create(model).Error
}

// Exists 判断是否存在
func (h *GORMHandler) Exists(wrapDBs ...func(*gorm.DB) *gorm.DB) (bool, error) {
	queryDB := h.db
	var total int64
	apply(queryDB, wrapDBs...)
	return total != 0, queryDB.Count(&total).Error
}

// TxProcess 事务处理
func (h *GORMHandler) TxProcess(txFuncGroup ...func(h *GORMHandler) error) (err error) {
	// 开启事物
	db := h.db
	tx := db.Begin()
	h.db = tx
	if err = tx.Error; err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, doTx := range txFuncGroup {
		if err = doTx(h); err != nil {
			return err
		}
	}
	tx.Commit()
	return
}
