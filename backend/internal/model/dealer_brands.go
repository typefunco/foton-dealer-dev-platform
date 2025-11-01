package model

// BrandsUpdate представляет данные для обновления брендов дилера.
type BrandsUpdate struct {
	DealerName       string `json:"dealer_name"`       // Название дилера
	City             string `json:"city"`              // Город дилера
	Brands           string `json:"brands"`            // Список брендов через запятую
	BysideBusinesses string `json:"byside_businesses"` // Список побочных бизнесов через запятую
}

// BrandsUploadResponse представляет ответ на загрузку файла с брендами.
type BrandsUploadResponse struct {
	Status          string   `json:"status"`            // success или error
	Message         string   `json:"message"`           // Сообщение о результате
	UpdatedCount    int      `json:"updated_count"`     // Количество обновленных дилеров
	NotFoundDealers []string `json:"not_found_dealers"` // Список дилеров, которые не найдены
	ProcessingTime  string   `json:"processing_time"`   // Время обработки
}

// BrandsFileInfo содержит метаданные о файле с брендами.
type BrandsFileInfo struct {
	FileName  string `json:"file_name"`
	Year      int    `json:"year"`       // Год из названия файла
	Quarter   string `json:"quarter"`    // Квартал (Q1, Q2, Q3, Q4)
	TableName string `json:"table_name"` // Название таблицы dealer_net
}
