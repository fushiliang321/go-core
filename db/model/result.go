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

func (m *Model[t]) First() (res *t) {
	tx := m.Db.First(&res)
	if tx.RowsAffected == 0 {
		return nil
	}
	return
}

func (m *Model[t]) Find() (res *[]t) {
	tx := m.Db.Find(&res)
	if tx.RowsAffected == 0 {
		return nil
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
	total := m.Count()
	if total > 0 {
		lastPage = ((int(total) / limit) + 1)
		offset := (page - 1) * limit
		if res := m.Offset(offset).Limit(limit).Find(); res != nil {
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

func (m *Model[t]) Create(data t) int64 {
	tx := m.Db.Create(data)
	return tx.RowsAffected
}
