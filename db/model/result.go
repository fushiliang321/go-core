package model

func (m *Model[t]) Count() (total int64) {
	m.Db.Count(&total)
	return
}

func (m *Model[t]) Exists() bool {
	var total int64
	m.Db.Limit(1).Count(&total)
	return total > 0
}

func (m *Model[t]) First() (res t) {
	m.Db.First(&res)
	return
}

func (m *Model[t]) Find() (res []t) {
	m.Db.Find(&res)
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
	var res []t
	total := m.Count()
	if total > 0 {
		lastPage = ((int(total) / limit) + 1)
		offset := (page - 1) * limit
		res = m.Offset(offset).Limit(limit).Find()
	}
	return &PaginateData[t]{
		Page:     page,
		LastPage: lastPage,
		Limit:    limit,
		Total:    int(total),
		Lists:    res,
	}
}

func (m *Model[t]) Delete() int64 {
	var _t t
	tx := m.Db.Delete(&_t)
	return tx.RowsAffected
}

func (m *Model[t]) Update(column string, value interface{}) int64 {
	tx := m.Db.Update(column, value)
	return tx.RowsAffected
}

func (m *Model[t]) Updates(values interface{}) int64 {
	tx := m.Db.Updates(values)
	return tx.RowsAffected
}

func (m *Model[t]) Pluck(column string, dest any) {
	m.Db.Pluck(column, dest)
	return
}
