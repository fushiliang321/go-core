package model

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

func (m *Model[t]) Count() (total int64, err error) {
	tx := m.Db.Count(&total)
	return total, tx.Error
}

func (m *Model[t]) Exists() (bool, error) {
	var total int64
	var err error
	tx := m.Db.Limit(1).Select("1").Row()
	err = tx.Err()
	if err != nil {
		return false, err
	}
	err = tx.Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return false, err
	}
	return total != 0, tx.Err()
}

func (m *Model[t]) First() (*t, error) {
	var res t
	tx := m.Db.First(&res)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, tx.Error
		}
	}
	return &res, nil
}

func (m *Model[t]) Find() ([]t, error) {
	var res []t
	tx := m.Db.Find(&res)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, tx.Error
		}
	}
	return res, nil
}

func (m *Model[t]) Paginate(page int, limit int) *PaginateData[t] {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	lastPage := 0
	var lists []t
	total, _ := m.Count()
	if total > 0 {
		totalInt := int(total)
		lastPage = totalInt / limit
		if (totalInt % limit) > 0 {
			lastPage++
		}
		offset := (page - 1) * limit
		if res, _ := m.Offset(offset).Limit(limit).Find(); res != nil {
			lists = res
		} else {
			lists = []t{}
		}
	} else {
		lists = []t{}
	}
	return &PaginateData[t]{
		Page:     page,
		LastPage: lastPage,
		Limit:    limit,
		Total:    int(total),
		Lists:    lists,
	}
}

func (m *Model[t]) Delete() (int64, error) {
	var _t t
	tx := m.Db.Delete(&_t)
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Update(column string, value any) (int64, error) {
	tx := m.Db.Update(column, value)
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Updates(values any) (int64, error) {
	tx := m.Db.Updates(values)
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Inc(column string, values ...int) (int64, error) {
	var v int
	if len(values) > 0 && values[0] > 1 {
		v = values[0]
	} else {
		v = 1
	}
	tx := m.Db.Update(column, gorm.Expr(column+" + ?", v))
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Dec(column string, values ...int) (int64, error) {
	var v int
	if len(values) > 0 && values[0] > 1 {
		v = values[0]
	} else {
		v = 1
	}
	tx := m.Db.Update(column, gorm.Expr(column+" - ?", v))
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Pluck(column string, dest any) error {
	tx := m.Db.Pluck(column, dest)
	return tx.Error
}

func (m *Model[t]) Create(data any) (int64, error) {
	tx := m.Db.Create(data)
	return tx.RowsAffected, tx.Error
}
