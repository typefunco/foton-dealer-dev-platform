package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// ExcelDealerRepository репозиторий для работы с данными дилеров из Excel таблиц.
type ExcelDealerRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	sq     squirrel.StatementBuilderType
}

// NewExcelDealerRepository создает новый экземпляр репозитория.
func NewExcelDealerRepository(pool *pgxpool.Pool, logger *slog.Logger) *ExcelDealerRepository {
	return &ExcelDealerRepository{
		pool:   pool,
		logger: logger,
		sq:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// GetDealerNetTableName возвращает название таблицы dealer_net для указанного года и квартала.
func (r *ExcelDealerRepository) GetDealerNetTableName(year int, quarter string) string {
	return fmt.Sprintf("dealer_net_%d_%s", year, strings.ToLower(quarter))
}

// TableExists проверяет существование таблицы dealer_net.
func (r *ExcelDealerRepository) TableExists(ctx context.Context, year int, quarter string) (bool, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check table existence",
			slog.String("table_name", tableName),
			slog.String("error", err.Error()),
		)
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}

	return exists, nil
}

// GetDealersWithFilters получает дилеров из таблицы dealer_net с фильтрами.
func (r *ExcelDealerRepository) GetDealersWithFilters(ctx context.Context, year int, quarter string, filters *model.FilterParams) ([]*model.Dealer, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		r.logger.Warn("Dealer net table does not exist",
			slog.String("table_name", tableName),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []*model.Dealer{}, nil
	}

	// Строим запрос
	query := r.sq.Select("dealer", "region", "city", "manager").
		From(tableName).
		Where(squirrel.NotEq{"dealer": nil}).
		Where(squirrel.NotEq{"dealer": ""})

	// Применяем фильтры
	if filters.HasRegionFilter() {
		query = query.Where(squirrel.Eq{"region": filters.Region})
	}

	if filters.HasDealerFilter() {
		// Для Excel таблиц у нас нет ID дилеров, поэтому фильтруем по названию
		dealerNames := make([]string, len(filters.DealerIDs))
		for i, id := range filters.DealerIDs {
			// Здесь нужно получить название дилера по ID из основной таблицы dealers
			// Пока что используем ID как строку
			dealerNames[i] = fmt.Sprintf("Dealer %d", id)
		}
		query = query.Where(squirrel.Eq{"dealer": dealerNames})
	}

	// Сортировка
	if filters.SortBy != "" {
		order := "ASC"
		if filters.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.OrderBy(fmt.Sprintf("%s %s", filters.SortBy, order))
	} else {
		query = query.OrderBy("dealer")
	}

	// Пагинация
	if filters.Limit > 0 {
		query = query.Limit(uint64(filters.Limit))
	}
	if filters.Offset > 0 {
		query = query.Offset(uint64(filters.Offset))
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealersWithFilters: error building query: %w", err)
	}

	r.logger.Info("Executing dealer query",
		slog.String("table_name", tableName),
		slog.String("sql", sql),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealersWithFilters: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		var dealerName, region, city, manager *string
		err = rows.Scan(&dealerName, &region, &city, &manager)
		if err != nil {
			return nil, fmt.Errorf("ExcelDealerRepository.GetDealersWithFilters: error scanning row: %w", err)
		}

		// Создаем объект дилера
		dealer := &model.Dealer{
			DealerNameRu: getStringValue(dealerName),
			Region:       getStringValue(region),
			City:         getStringValue(city),
			Manager:      getStringValue(manager),
		}

		dealers = append(dealers, dealer)
	}

	r.logger.Info("Dealers retrieved from Excel table",
		slog.String("table_name", tableName),
		slog.Int("count", len(dealers)),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	return dealers, nil
}

// GetDealerCardData получает данные карточки дилера из таблицы dealer_net.
func (r *ExcelDealerRepository) GetDealerCardData(ctx context.Context, year int, quarter string, dealerName string) (*model.DealerCardData, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		r.logger.Warn("Dealer net table does not exist",
			slog.String("table_name", tableName),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return &model.DealerCardData{}, nil
	}

	// Получаем данные дилера
	query := r.sq.Select("*").
		From(tableName).
		Where(squirrel.Eq{"dealer": dealerName}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealerCardData: error building query: %w", err)
	}

	r.logger.Info("Executing dealer card query",
		slog.String("table_name", tableName),
		slog.String("dealer_name", dealerName),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealerCardData: error querying: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return &model.DealerCardData{}, nil
	}

	// Получаем названия колонок
	fieldDescriptions := rows.FieldDescriptions()
	values := make([]interface{}, len(fieldDescriptions))
	valuePointers := make([]interface{}, len(fieldDescriptions))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	err = rows.Scan(valuePointers...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealerCardData: error scanning row: %w", err)
	}

	// Создаем карточку дилера
	cardData := &model.DealerCardData{
		DealerNameRu: dealerName,
		Region:       "",
		City:         "",
		Manager:      "",
	}

	// Заполняем данные из колонок
	for i, field := range fieldDescriptions {
		if i < len(values) && values[i] != nil {
			switch field.Name {
			case "region":
				if region, ok := values[i].(string); ok {
					cardData.Region = region
				}
			case "city":
				if city, ok := values[i].(string); ok {
					cardData.City = city
				}
			case "manager":
				if manager, ok := values[i].(string); ok {
					cardData.Manager = manager
				}
			}
		}
	}

	r.logger.Info("Dealer card data retrieved",
		slog.String("table_name", tableName),
		slog.String("dealer_name", dealerName),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	return cardData, nil
}

// GetAvailableRegions получает список доступных регионов из таблицы dealer_net.
func (r *ExcelDealerRepository) GetAvailableRegions(ctx context.Context, year int, quarter string) ([]string, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		r.logger.Warn("Dealer net table does not exist",
			slog.String("table_name", tableName),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []string{}, nil
	}

	query := r.sq.Select("DISTINCT region").
		From(tableName).
		Where(squirrel.NotEq{"region": nil}).
		Where(squirrel.NotEq{"region": ""}).
		OrderBy("region")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetAvailableRegions: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetAvailableRegions: error querying: %w", err)
	}
	defer rows.Close()

	var regions []string
	for rows.Next() {
		var region *string
		err = rows.Scan(&region)
		if err != nil {
			return nil, fmt.Errorf("ExcelDealerRepository.GetAvailableRegions: error scanning row: %w", err)
		}
		regions = append(regions, getStringValue(region))
	}

	r.logger.Info("Available regions retrieved",
		slog.String("table_name", tableName),
		slog.Int("count", len(regions)),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	return regions, nil
}

// GetSalesDataFromExcel получает данные продаж из таблицы dealer_net.
func (r *ExcelDealerRepository) GetSalesDataFromExcel(ctx context.Context, year int, quarter string, region string) ([]*model.SalesWithDetails, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		r.logger.Warn("Dealer net table does not exist",
			slog.String("table_name", tableName),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []*model.SalesWithDetails{}, nil
	}

	// Строим запрос для получения данных продаж
	query := r.sq.Select("dealer", "region", "city", "manager", "hdt", "mdt", "ldt", "hdt_2", "mdt_2", "ldt_2", "sales").
		From(tableName).
		Where(squirrel.NotEq{"dealer": nil}).
		Where(squirrel.NotEq{"dealer": ""})

	// Применяем фильтр по региону если указан
	if region != "" && region != "all-russia" {
		query = query.Where(squirrel.Eq{"region": region})
	}

	query = query.OrderBy("dealer")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetSalesDataFromExcel: error building query: %w", err)
	}

	r.logger.Info("Executing sales data query",
		slog.String("table_name", tableName),
		slog.String("region", region),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetSalesDataFromExcel: error querying: %w", err)
	}
	defer rows.Close()

	var salesData []*model.SalesWithDetails
	for rows.Next() {
		var dealerName, region, city, manager, sales *string
		var hdt, mdt, ldt, hdt2, mdt2, ldt2 *string

		err = rows.Scan(&dealerName, &region, &city, &manager, &hdt, &mdt, &ldt, &hdt2, &mdt2, &ldt2, &sales)
		if err != nil {
			return nil, fmt.Errorf("ExcelDealerRepository.GetSalesDataFromExcel: error scanning row: %w", err)
		}

		// Создаем объект данных продаж
		salesDetail := &model.SalesWithDetails{
			Sales: model.Sales{
				StockHDT:    getIntFromString(hdt),
				StockMDT:    getIntFromString(mdt),
				StockLDT:    getIntFromString(ldt),
				BuyoutHDT:   getIntFromString(hdt2),
				BuyoutMDT:   getIntFromString(mdt2),
				BuyoutLDT:   getIntFromString(ldt2),
				SalesTarget: getStringValue(sales),
			},
			DealerNameRu: getStringValue(dealerName),
			Region:       getStringValue(region),
			City:         getStringValue(city),
			Manager:      getStringValue(manager),
		}

		salesData = append(salesData, salesDetail)
	}

	r.logger.Info("Sales data retrieved from Excel table",
		slog.String("table_name", tableName),
		slog.String("region", region),
		slog.Int("count", len(salesData)),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	return salesData, nil
}

// GetDealerDevDataFromExcel получает данные дилер-девелопмента из таблицы dealer_net.
func (r *ExcelDealerRepository) GetDealerDevDataFromExcel(ctx context.Context, year int, quarter string, region string) ([]*model.DealerDevWithDetails, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		r.logger.Warn("Dealer net table does not exist",
			slog.String("table_name", tableName),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []*model.DealerDevWithDetails{}, nil
	}

	// Строим запрос для получения данных дилер-девелопмента
	query := r.sq.Select("dealer", "region", "city", "manager", "class", "check_list_percent", "marketing_investments", "branding", "dealer_development").
		From(tableName).
		Where(squirrel.NotEq{"dealer": nil}).
		Where(squirrel.NotEq{"dealer": ""})

	// Применяем фильтр по региону если указан
	if region != "" && region != "all-russia" {
		query = query.Where(squirrel.Eq{"region": region})
	}

	query = query.OrderBy("dealer")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealerDevDataFromExcel: error building query: %w", err)
	}

	r.logger.Info("Executing dealer dev data query",
		slog.String("table_name", tableName),
		slog.String("region", region),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetDealerDevDataFromExcel: error querying: %w", err)
	}
	defer rows.Close()

	var dealerDevData []*model.DealerDevWithDetails
	for rows.Next() {
		var dealerName, region, city, manager, class, checkListPercent, marketingInvestments, branding, dealerDevelopment *string

		err = rows.Scan(&dealerName, &region, &city, &manager, &class, &checkListPercent, &marketingInvestments, &branding, &dealerDevelopment)
		if err != nil {
			return nil, fmt.Errorf("ExcelDealerRepository.GetDealerDevDataFromExcel: error scanning row: %w", err)
		}

		// Создаем объект данных дилер-девелопмента
		dealerDevDetail := &model.DealerDevWithDetails{
			DealerDevelopment: model.DealerDevelopment{
				DealershipClass:      getStringValue(class),
				CheckListScore:       getIntFromString(checkListPercent),
				MarketingInvestments: int64(getFloatFromString(marketingInvestments)),
				Branding:             getStringValue(branding) == "Yes" || getStringValue(branding) == "Y",
				DDRecommendation:     getStringValue(dealerDevelopment),
			},
			DealerNameRu: getStringValue(dealerName),
			Region:       getStringValue(region),
			City:         getStringValue(city),
			Manager:      getStringValue(manager),
		}

		dealerDevData = append(dealerDevData, dealerDevDetail)
	}

	r.logger.Info("Dealer dev data retrieved from Excel table",
		slog.String("table_name", tableName),
		slog.String("region", region),
		slog.Int("count", len(dealerDevData)),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	return dealerDevData, nil
}

// GetAfterSalesDataFromExcel получает данные автозапчастей из таблицы dealer_net.
func (r *ExcelDealerRepository) GetAfterSalesDataFromExcel(ctx context.Context, year int, quarter string, region string) ([]*model.AfterSalesWithDetails, error) {
	tableName := r.GetDealerNetTableName(year, quarter)

	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		r.logger.Warn("Dealer net table does not exist",
			slog.String("table_name", tableName),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []*model.AfterSalesWithDetails{}, nil
	}

	// Строим запрос для получения данных автозапчастей
	query := r.sq.Select("dealer", "region", "city", "manager", "service_contracts_sales", "spare_parts_sales_q3", "spare_parts_sales_ytd_percent", "warranty_stock_percent", "recommended_stock_percent", "foton_labour_hours", "foton_labour_hours_share", "warranty_hours", "service_contracts_hours", "as_trainings", "aftersales").
		From(tableName).
		Where(squirrel.NotEq{"dealer": nil}).
		Where(squirrel.NotEq{"dealer": ""})

	// Применяем фильтр по региону если указан
	if region != "" && region != "all-russia" {
		query = query.Where(squirrel.Eq{"region": region})
	}

	query = query.OrderBy("dealer")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetAfterSalesDataFromExcel: error building query: %w", err)
	}

	r.logger.Info("Executing after sales data query",
		slog.String("table_name", tableName),
		slog.String("region", region),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ExcelDealerRepository.GetAfterSalesDataFromExcel: error querying: %w", err)
	}
	defer rows.Close()

	var afterSalesData []*model.AfterSalesWithDetails
	for rows.Next() {
		var dealerName, region, city, manager, serviceContractsSales, sparePartsSalesQ3, sparePartsSalesYtdPercent, warrantyStockPercent, recommendedStockPercent, fotonLabourHours, fotonLabourHoursShare, warrantyHours, serviceContractsHours, asTrainings, aftersales *string

		err = rows.Scan(&dealerName, &region, &city, &manager, &serviceContractsSales, &sparePartsSalesQ3, &sparePartsSalesYtdPercent, &warrantyStockPercent, &recommendedStockPercent, &fotonLabourHours, &fotonLabourHoursShare, &warrantyHours, &serviceContractsHours, &asTrainings, &aftersales)
		if err != nil {
			return nil, fmt.Errorf("ExcelDealerRepository.GetAfterSalesDataFromExcel: error scanning row: %w", err)
		}

		// Создаем объект данных автозапчастей
		afterSalesDetail := &model.AfterSalesWithDetails{
			AfterSales: model.AfterSales{
				RecommendedStock:   getIntFromString(recommendedStockPercent),
				WarrantyStock:      getIntFromString(warrantyStockPercent),
				FotonLaborHours:    getIntFromString(fotonLabourHours),
				FotonWarrantyHours: getIntFromString(warrantyHours),
				ServiceContracts:   getIntFromString(serviceContractsHours),
				ASTrainings:        getStringValue(asTrainings) == "Yes" || getStringValue(asTrainings) == "Y",
				ASDecision:         getStringValue(aftersales),
			},
			DealerNameRu:          getStringValue(dealerName),
			Region:                getStringValue(region),
			City:                  getStringValue(city),
			Manager:               getStringValue(manager),
			SparePartsSalesQ3:     getStringValue(sparePartsSalesQ3),
			SparePartsSalesYtd:    getStringValue(sparePartsSalesYtdPercent),
			FotonLabourHoursShare: getStringValue(fotonLabourHoursShare),
		}

		afterSalesData = append(afterSalesData, afterSalesDetail)
	}

	r.logger.Info("After sales data retrieved from Excel table",
		slog.String("table_name", tableName),
		slog.String("region", region),
		slog.Int("count", len(afterSalesData)),
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	return afterSalesData, nil
}

// getStringValue возвращает строку из указателя или пустую строку если указатель nil.
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getIntValue возвращает int из указателя или 0 если указатель nil.
func getIntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// getIntFromString возвращает int из строки или 0 если строка пустая или невалидная.
func getIntFromString(s *string) int {
	if s == nil || *s == "" {
		return 0
	}

	// Пытаемся преобразовать строку в int
	if val, err := strconv.Atoi(*s); err == nil {
		return val
	}

	return 0
}

// getFloatFromString возвращает float64 из строки или 0 если строка пустая или невалидная.
func getFloatFromString(s *string) float64 {
	if s == nil || *s == "" {
		return 0
	}

	// Пытаемся преобразовать строку в float64
	if val, err := strconv.ParseFloat(*s, 64); err == nil {
		return val
	}

	return 0
}
