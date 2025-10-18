package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const dealerTableName = "dealers"

// DealerRepository репозиторий для работы с дилерами.
type DealerRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewDealerRepository конструктор.
func NewDealerRepository(pool *pgxpool.Pool) *DealerRepository {
	return &DealerRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает нового дилера.
func (r *DealerRepository) Create(ctx context.Context, dealer *model.Dealer) (int, error) {
	now := time.Now()
	dealer.CreatedAt = now
	dealer.UpdatedAt = now

	query := r.sq.Insert(dealerTableName).
		Columns("name", "city", "region", "manager", "created_at", "updated_at").
		Values(dealer.DealerNameRu, dealer.City, dealer.Region, dealer.Manager, dealer.CreatedAt, dealer.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("DealerRepository.Create: error building query: %w", err)
	}

	var id int64
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("DealerRepository.Create: error inserting: %w", err)
	}

	dealer.DealerID = int(id)
	return int(id), nil
}

// GetByID получает дилера по ID.
func (r *DealerRepository) GetByID(ctx context.Context, id int) (*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByID: error building query: %w", err)
	}

	dealer := &model.Dealer{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
		&dealer.CreatedAt, &dealer.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByID: error scanning: %w", err)
	}

	return dealer, nil
}

// GetAll получает всех дилеров.
func (r *DealerRepository) GetAll(ctx context.Context) ([]*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		OrderBy("name")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetAll: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetAll: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		dealer := &model.Dealer{}
		err = rows.Scan(
			&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetAll: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// GetByRegion получает дилеров по региону.
func (r *DealerRepository) GetByRegion(ctx context.Context, region string) ([]*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		Where(squirrel.Eq{"region": region}).
		OrderBy("name")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByRegion: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByRegion: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		dealer := &model.Dealer{}
		err = rows.Scan(
			&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetByRegion: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// GetWithFilters получает дилеров с применением фильтров.
func (r *DealerRepository) GetWithFilters(ctx context.Context, filters *model.FilterParams) ([]*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName)

	// Применяем фильтры
	if filters.HasRegionFilter() {
		query = query.Where(squirrel.Eq{"region": filters.Region})
	}

	if filters.HasDealerFilter() {
		query = query.Where(squirrel.Eq{"id": filters.DealerIDs})
	}

	// Сортировка
	if filters.SortBy != "" {
		order := "ASC"
		if filters.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.OrderBy(fmt.Sprintf("%s %s", filters.SortBy, order))
	} else {
		query = query.OrderBy("name")
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
		return nil, fmt.Errorf("DealerRepository.GetWithFilters: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetWithFilters: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		dealer := &model.Dealer{}
		err = rows.Scan(
			&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetWithFilters: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// Update обновляет данные дилера.
func (r *DealerRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()
	query := r.sq.Update(dealerTableName).SetMap(updates).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.Update: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.Update: error updating: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DealerRepository.Update: no rows affected, record with id %d not found", id)
	}

	return nil
}

// Delete удаляет дилера.
func (r *DealerRepository) Delete(ctx context.Context, id int) error {
	query := r.sq.Delete(dealerTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.Delete: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.Delete: error deleting: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DealerRepository.Delete: no rows affected, record with id %d not found", id)
	}

	return nil
}

// GetDealerCardData получает данные карточки дилера за период.
func (r *DealerRepository) GetDealerCardData(ctx context.Context, dealerID int, period time.Time) (*model.DealerCardData, error) {
	// Определяем квартал и год из периода
	quarter := getQuarterFromTime(period)
	year := period.Year()

	// Получаем основную информацию о дилере
	dealer, err := r.GetByID(ctx, dealerID)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetDealerCardData: error getting dealer: %w", err)
	}

	// Создаем структуру карточки дилера
	cardData := &model.DealerCardData{
		DealerID:      dealer.DealerID,
		Ruft:          dealer.Ruft,
		DealerNameRu:  dealer.DealerNameRu,
		DealerNameEn:  dealer.DealerNameEn,
		City:          dealer.City,
		Region:        dealer.Region,
		Manager:       dealer.Manager,
		JointDecision: dealer.JointDecision,
		Period:        period,
	}

	// Получаем данные Dealer Development
	ddData, err := r.getDealerDevData(ctx, dealerID, quarter, year)
	if err == nil && ddData != nil {
		checklistScore := float64(ddData.CheckListScore)
		cardData.CheckListScore = &checklistScore
		class := model.DealershipClass(ddData.DealershipClass)
		cardData.DealershipClass = &class
		branding := model.BrandingStatus("No")
		if ddData.Branding {
			branding = model.BrandingStatus("Yes")
		}
		cardData.Branding = &branding
		marketingInvestments := float64(ddData.MarketingInvestments)
		cardData.MarketingInvestments = &marketingInvestments
		cardData.DDRecommendation = &ddData.DDRecommendation
	}

	// Получаем бренды дилера
	brands, err := r.GetBrands(ctx, dealerID)
	if err == nil {
		cardData.Brands = brands
	}

	// Получаем побочный бизнес дилера
	businesses, err := r.GetBusinesses(ctx, dealerID)
	if err == nil && len(businesses) > 0 {
		businessStr := ""
		for i, business := range businesses {
			if i > 0 {
				businessStr += ", "
			}
			businessStr += business
		}
		cardData.BySideBusinesses = &businessStr
	}

	// Получаем данные Sales
	salesData, err := r.getSalesData(ctx, dealerID, quarter, year)
	if err == nil && salesData != nil {
		cardData.StockHDT = &salesData.StockHDT
		cardData.StockMDT = &salesData.StockMDT
		cardData.StockLDT = &salesData.StockLDT
		cardData.BuyoutHDT = &salesData.BuyoutHDT
		cardData.BuyoutMDT = &salesData.BuyoutMDT
		cardData.BuyoutLDT = &salesData.BuyoutLDT
		cardData.FotonSalesPersonnel = &salesData.FotonSalesmen
		salesTrainings := model.SalesTrainingsStatus("No")
		if salesData.SalesTrainings {
			salesTrainings = model.SalesTrainingsStatus("Yes")
		}
		cardData.SalesTrainings = &salesTrainings
		serviceContractsSales := float64(salesData.ServiceContractsSales)
		cardData.ServiceContractsSales = &serviceContractsSales
		cardData.SalesRecommendation = &salesData.SalesDecision
	}

	// Получаем данные Performance (пока используем mock данные)
	// TODO: Реализовать получение данных из таблиц performance_sales и performance_aftersales

	// Получаем данные After Sales
	asData, err := r.getAfterSalesData(ctx, dealerID, quarter, year)
	if err == nil && asData != nil {
		asRevenue := float64(asData.ServiceContracts)
		cardData.ASRevenue = &asRevenue
		asMargin := float64(asData.FotonLaborHours)
		cardData.ASMargin = &asMargin
		asMarginPct := float64(asData.FotonWarrantyHours)
		cardData.ASMarginPct = &asMarginPct
		asProfitPct := float64(asData.RecommendedStock)
		cardData.ASProfitPct = &asProfitPct
	}

	return cardData, nil
}

// getQuarterFromTime определяет квартал по времени
func getQuarterFromTime(t time.Time) string {
	month := t.Month()
	switch {
	case month >= 1 && month <= 3:
		return "Q1"
	case month >= 4 && month <= 6:
		return "Q2"
	case month >= 7 && month <= 9:
		return "Q3"
	case month >= 10 && month <= 12:
		return "Q4"
	default:
		return "Q1"
	}
}

// getDealerDevData получает данные развития дилера
func (r *DealerRepository) getDealerDevData(ctx context.Context, dealerID int, quarter string, year int) (*model.DealerDevelopment, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"check_list_score", "dealer_ship_class", "branding",
		"marketing_investments", "dealer_dev_recommendation",
		"created_at", "updated_at",
	).From("dealer_dev").Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("getDealerDevData: error building query: %w", err)
	}

	dd := &model.DealerDevelopment{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dd.ID, &dd.DealerID, &dd.Quarter, &dd.Year,
		&dd.CheckListScore, &dd.DealershipClass, &dd.Branding,
		&dd.MarketingInvestments, &dd.DDRecommendation,
		&dd.CreatedAt, &dd.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getDealerDevData: error scanning: %w", err)
	}

	return dd, nil
}

// getSalesData получает данные продаж
func (r *DealerRepository) getSalesData(ctx context.Context, dealerID int, quarter string, year int) (*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_target", "stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_salesmen", "service_contracts_sales", "sales_trainings", "sales_decision",
		"created_at", "updated_at",
	).From("sales").Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("getSalesData: error building query: %w", err)
	}

	sales := &model.Sales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&sales.ID, &sales.DealerID, &sales.Quarter, &sales.Year,
		&sales.SalesTarget, &sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
		&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
		&sales.FotonSalesmen, &sales.ServiceContractsSales, &sales.SalesTrainings, &sales.SalesDecision,
		&sales.CreatedAt, &sales.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getSalesData: error scanning: %w", err)
	}

	return sales, nil
}

// getAfterSalesData получает данные послепродажного обслуживания
func (r *DealerRepository) getAfterSalesData(ctx context.Context, dealerID int, quarter string, year int) (*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"recommended_stock", "warranty_stock", "foton_labor_hours",
		"service_contracts", "as_trainings", "csi", "foton_warranty_hours", "as_decision",
		"created_at", "updated_at",
	).From("after_sales").Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("getAfterSalesData: error building query: %w", err)
	}

	as := &model.AfterSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&as.ID, &as.DealerID, &as.Quarter, &as.Year,
		&as.RecommendedStock, &as.WarrantyStock, &as.FotonLaborHours,
		&as.ServiceContracts, &as.ASTrainings, &as.CSI, &as.FotonWarrantyHours, &as.ASDecision,
		&as.CreatedAt, &as.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getAfterSalesData: error scanning: %w", err)
	}

	return as, nil
}

// AddBrand добавляет бренд дилеру.
func (r *DealerRepository) AddBrand(ctx context.Context, dealerID int, brandName string) error {
	query := r.sq.Insert("dealer_brands").
		Columns("dealer_id", "brand_name", "created_at").
		Values(dealerID, brandName, time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBrand: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBrand: error inserting: %w", err)
	}

	return nil
}

// RemoveBrand удаляет бренд у дилера.
func (r *DealerRepository) RemoveBrand(ctx context.Context, dealerID int, brandName string) error {
	query := r.sq.Delete("dealer_brands").
		Where(squirrel.Eq{"dealer_id": dealerID, "brand_name": brandName})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBrand: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBrand: error deleting: %w", err)
	}

	return nil
}

// GetBrands получает список брендов дилера.
func (r *DealerRepository) GetBrands(ctx context.Context, dealerID int) ([]string, error) {
	query := r.sq.Select("brand_name").
		From("dealer_brands").
		Where(squirrel.Eq{"dealer_id": dealerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBrands: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBrands: error querying: %w", err)
	}
	defer rows.Close()

	var brands []string
	for rows.Next() {
		var brand string
		err = rows.Scan(&brand)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetBrands: error scanning: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// AddBusiness добавляет тип бизнеса дилеру.
func (r *DealerRepository) AddBusiness(ctx context.Context, dealerID int, businessType string) error {
	query := r.sq.Insert("dealer_businesses").
		Columns("dealer_id", "business_type", "created_at").
		Values(dealerID, businessType, time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBusiness: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBusiness: error inserting: %w", err)
	}

	return nil
}

// RemoveBusiness удаляет тип бизнеса у дилера.
func (r *DealerRepository) RemoveBusiness(ctx context.Context, dealerID int, businessType string) error {
	query := r.sq.Delete("dealer_businesses").
		Where(squirrel.Eq{"dealer_id": dealerID, "business_type": businessType})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBusiness: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBusiness: error deleting: %w", err)
	}

	return nil
}

// GetBusinesses получает список типов бизнеса дилера.
func (r *DealerRepository) GetBusinesses(ctx context.Context, dealerID int) ([]string, error) {
	query := r.sq.Select("business_type").
		From("dealer_businesses").
		Where(squirrel.Eq{"dealer_id": dealerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBusinesses: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBusinesses: error querying: %w", err)
	}
	defer rows.Close()

	var businesses []string
	for rows.Next() {
		var business string
		err = rows.Scan(&business)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetBusinesses: error scanning: %w", err)
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}

// UpdateFull обновляет всю запись дилера целиком.
func (r *DealerRepository) UpdateFull(ctx context.Context, dealer *model.Dealer) error {
	dealer.UpdatedAt = time.Now()

	query := r.sq.Update(dealerTableName).
		Set("name", dealer.DealerNameRu).
		Set("region", dealer.Region).
		Set("city", dealer.City).
		Set("manager", dealer.Manager).
		Set("updated_at", dealer.UpdatedAt).
		Where(squirrel.Eq{"id": dealer.DealerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.UpdateFull: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.UpdateFull: error updating: %w", err)
	}

	return nil
}
