package model

func (m *Model[t]) Count() (total int64, err error) {
	tx := m.Db.Count(&total)
	return total, tx.Error
}

func (m *Model[t]) Exists() (bool, error) {
	var total int64
	tx := m.Db.Limit(1).Count(&total)
	return total > 0, tx.Error
}

func (m *Model[t]) First() (res *t, err error) {
	tx := m.Db.First(&res)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return
}

func (m *Model[t]) Find() (res *[]t, err error) {
	tx := m.Db.Find(&res)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return
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
		lastPage = ((int(total) / limit) + 1)
		offset := (page - 1) * limit
		if res, _ := m.Offset(offset).Limit(limit).Find(); res != nil {
			lists = *res
		}
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

func (m *Model[t]) Update(column string, value interface{}) (int64, error) {
	tx := m.Db.Update(column, value)
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Updates(values interface{}) (int64, error) {
	tx := m.Db.Updates(values)
	return tx.RowsAffected, tx.Error
}

func (m *Model[t]) Pluck(column string, dest any) error {
	tx := m.Db.Pluck(column, dest)
	return tx.Error
}

func (m *Model[t]) Create(data *t) (int64, error) {
	tx := m.Db.Create(data)
	return tx.RowsAffected, tx.Error
}
