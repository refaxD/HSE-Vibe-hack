package usecase

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParseService struct {
	categoryRepo   domain.CategoryRepository
	placeRepo      domain.PlaceRepository
	eventPlaceRepo domain.EventPlaceRepository
	cfg            *config.Config
	httpClient     *http.Client
}

func NewParseService(
	categoryRepo domain.CategoryRepository,
	placeRepo domain.PlaceRepository,
	eventPlaceRepo domain.EventPlaceRepository,
	cfg *config.Config,
) *ParseService {
	return &ParseService{
		categoryRepo:   categoryRepo,
		placeRepo:      placeRepo,
		eventPlaceRepo: eventPlaceRepo,
		cfg:            cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (s *ParseService) Run(ctx context.Context) {
	log.Println("[parseService] starting initCategories")
	if err := s.initCategories(ctx); err != nil {
		log.Printf("[parseService] initCategories error: %v", err)
		return
	}
	log.Println("[parseService] finished initCategories")

	log.Println("[parseService] starting initPlaces")
	if err := s.initPlaces(ctx); err != nil {
		log.Printf("[parseService] initPlaces error: %v", err)
		return
	}
	log.Println("[parseService] finished initPlaces")

	log.Println("[parseService] starting initEventPlaces")
	if err := s.initEventPlaces(ctx); err != nil {
		log.Printf("[parseService] initEventPlaces error: %v", err)
	}
	log.Println("[parseService] finished initEventPlaces")
}

// ---- Categories ----

func (s *ParseService) initCategories(ctx context.Context) error {
	if err := s.categoryRepo.Clear(ctx); err != nil {
		return err
	}

	categories := getDefaultCategories()
	for i, cat := range categories {
		if err := s.categoryRepo.Save(ctx, &cat); err != nil {
			return fmt.Errorf("save category %d: %w", i, err)
		}
	}
	log.Printf("[parseService] saved %d categories", len(categories))
	return nil
}

// ---- Places ----

type placesDataFile struct {
	Data []placeDataDef `json:"data"`
}

type placeDataDef struct {
	Name      string            `json:"name"`
	ID        int               `json:"id"`
	AutoParse bool              `json:"autoParse"`
	Category  string            `json:"category"`
	Type      int               `json:"type"`
	Geometry  string            `json:"geometry"`
	Base      []string          `json:"base"`
	Fields    map[string]any    `json:"fields"`
	Data      map[string]string `json:"data"`
}

func (s *ParseService) initPlaces(ctx context.Context) error {
	if err := s.placeRepo.Clear(ctx); err != nil {
		return err
	}
	if err := s.categoryRepo.ClearPlaces(ctx); err != nil {
		return err
	}

	data, err := os.ReadFile("internal/usecase/places-data.json")
	if err != nil {
		return fmt.Errorf("read places-data.json: %w", err)
	}

	var pf placesDataFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return fmt.Errorf("parse places-data.json: %w", err)
	}

	for i, datum := range pf.Data {
		if err := s.getPlaceFromAPI(ctx, datum); err != nil {
			log.Printf("[parseService] error processing dataset %d (%s): %v", i, datum.Name, err)
		}
	}
	return nil
}

func (s *ParseService) getPlaceFromAPI(ctx context.Context, datum placeDataDef) error {
	apiURL := fmt.Sprintf("https://apidata.mos.ru/v1/datasets/%d/features?api_key=%s", datum.ID, s.cfg.MosDataKey)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("create request dataset %d: %w", datum.ID, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch dataset %d: %w", datum.ID, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body dataset %d: %w", datum.ID, err)
	}

	if resp.StatusCode != http.StatusOK {
		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200]
		}
		return fmt.Errorf("dataset %d: HTTP %d: %s", datum.ID, resp.StatusCode, preview)
	}

	var geoJSON struct {
		Features []map[string]any `json:"features"`
	}
	if err := json.Unmarshal(body, &geoJSON); err != nil {
		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200]
		}
		return fmt.Errorf("unmarshal dataset %d (HTTP %d, body: %s): %w", datum.ID, resp.StatusCode, preview, err)
	}

	log.Printf("[parseService] dataset %d (%s): %d features", datum.ID, datum.Name, len(geoJSON.Features))

	for _, feature := range geoJSON.Features {
		places, categoryNames := s.configurePlace(datum, feature)
		for i, place := range places {
			place.ID = primitive.NewObjectID()
			if err := s.placeRepo.Save(ctx, &place); err != nil {
				continue
			}

			var catIDs []primitive.ObjectID
			for _, catName := range categoryNames[i] {
				cat, err := s.categoryRepo.FindByName(ctx, catName)
				if err != nil {
					continue
				}
				catIDs = append(catIDs, cat.ID)
				_ = s.categoryRepo.AddPlace(ctx, cat.ID, place.ID)
			}
			if len(catIDs) > 0 {
				_ = s.placeRepo.SetCategories(ctx, place.ID, catIDs)
			}
		}
	}
	return nil
}

func (s *ParseService) configurePlace(datum placeDataDef, feature map[string]any) ([]domain.Place, [][]string) {
	var places []domain.Place
	var allCatNames [][]string

	geom := getMap(feature, "geometry")
	coords := geom["coordinates"]

	switch datum.Geometry {
	case "Point":
		lat, lon := extractPoint(coords)
		places = append(places, domain.Place{Type: datum.Type, Latitude: lat, Longitude: lon})
	case "MultiPoint":
		if arr, ok := coords.([]any); ok {
			for _, c := range arr {
				lat, lon := extractPoint(c)
				places = append(places, domain.Place{Type: datum.Type, Latitude: lat, Longitude: lon})
			}
		}
	case "Polygon":
		if rings, ok := coords.([]any); ok && len(rings) > 0 {
			if ring, ok := rings[0].([]any); ok && len(ring) > 0 {
				var sumLat, sumLon float64
				for _, c := range ring {
					lat, lon := extractPoint(c)
					sumLat += lat
					sumLon += lon
				}
				n := float64(len(ring))
				places = append(places, domain.Place{Type: datum.Type, Latitude: sumLat / n, Longitude: sumLon / n})
			}
		}
	}

	// Navigate to base
	base := navigateToBase(feature, datum.Base)

	for placeIdx := range places {
		// Fields
		fillFields(&places[placeIdx], datum, base, placeIdx)

		// Data
		fillData(&places[placeIdx], datum, base)

		// Categories
		catNames := s.resolveCategories(datum, feature)
		allCatNames = append(allCatNames, catNames)
	}

	// If categories resolved to empty for non-autoParse, skip
	if !datum.AutoParse && len(places) > 0 && len(allCatNames) > 0 && len(allCatNames[0]) == 0 {
		return nil, nil
	}

	return places, allCatNames
}

func fillFields(place *domain.Place, datum placeDataDef, base map[string]any, placeIdx int) {
	for fieldKey, fieldValue := range datum.Fields {
		if fieldValue == nil {
			continue
		}

		// Special handling for MultiPoint addresses
		if datum.Geometry == "MultiPoint" && fieldKey == "address" {
			if objAddr, ok := base["ObjectAddress"]; ok {
				if arr, ok := objAddr.([]any); ok && placeIdx < len(arr) {
					if addrMap, ok := arr[placeIdx].(map[string]any); ok {
						place.Address = toString(addrMap["Address"])
					}
				}
			}
			continue
		}

		switch v := fieldValue.(type) {
		case string:
			if v == "" {
				continue
			}
			val := base[v]
			setPlaceField(place, fieldKey, val)
		case []any:
			// Navigate through array path like ["Email", "0", "Email"]
			var current any = base
			for _, step := range v {
				key := fmt.Sprintf("%v", step)
				switch m := current.(type) {
				case map[string]any:
					current = m[key]
				case []any:
					idx := 0
					fmt.Sscanf(key, "%d", &idx)
					if idx < len(m) {
						current = m[idx]
					} else {
						current = nil
					}
				default:
					current = nil
				}
				if current == nil {
					break
				}
			}
			if current != nil {
				setPlaceField(place, fieldKey, current)
			}
		}
	}
}

func fillData(place *domain.Place, datum placeDataDef, base map[string]any) {
	pd := &domain.PlaceData{}
	hasData := false

	for dataKey, dataField := range datum.Data {
		if dataField == "" {
			continue
		}
		val := base[dataField]
		if val == nil {
			continue
		}
		hasData = true

		switch dataKey {
		case "hasFoodPoint", "hasFood":
			b := toBool(val)
			pd.HasFoodPoint = &b
		case "hasChangeRoom":
			b := toBool(val)
			pd.HasChangeRoom = &b
		case "hasToilet":
			b := toBool(val)
			pd.HasToilet = &b
		case "hasWIFI":
			b := toBool(val)
			pd.HasWIFI = &b
		case "hasWater":
			b := toBool(val)
			pd.HasWater = &b
		case "hasChild":
			b := toBool(val)
			pd.HasChild = &b
		case "hasSport":
			b := toBool(val)
			pd.HasSport = &b
		case "info":
			pd.Info = toString(val)
		case "priceInfo":
			pd.PriceInfo = toString(val)
		case "conditions":
			pd.Conditions = toString(val)
		case "time":
			pd.Time = val
		case "subway":
			pd.Subway = toString(val)
		}
	}

	if hasData {
		place.Data = pd
	}
}

func (s *ParseService) resolveCategories(datum placeDataDef, feature map[string]any) []string {
	if datum.AutoParse {
		return []string{datum.Category}
	}

	attrs := navigateToBase(feature, datum.Base)
	typeObj := toString(attrs["TypeObject"])

	switch datum.Name {
	case "Стационарные торговые объекты":
		return resolveShopCategory(typeObj)
	case "Общественное питание в Москве":
		return resolveFoodCategory(typeObj)
	}
	return nil
}

func resolveShopCategory(typeObj string) []string {
	switch typeObj {
	case "Магазин «Цветы»":
		return []string{"Магазин цветов"}
	case "Прочие специализированные непродовольственные предприятия торговли",
		"Прочие специализированные продовольственные предприятия торговли",
		"Магазин «Промтовары»",
		"Товары для дома",
		"Магазин «Мир садовода»",
		"Магазин товаров повседневного спроса",
		"Торговый Дом",
		"Магазин-склад непродовольственный (опт)",
		"Магазин «Хозяйственные товары»":
		return []string{"Промтовары"}
	case "Магазин «Алкогольные напитки»":
		return []string{"Алкомаркет"}
	case "Универсам", "Минимаркет", "Гастроном", "Универмаг", "Супермаркет":
		return []string{"Супермаркет"}
	case "Магазин «Секонд Хенд»", "Магазин «Дискаунтер»", "Магазин «Дисконт»", "Комиссионный магазин":
		return []string{"Дискаунтер"}
	case "Магазин «Одежда»", "Магазин «Спорт и туризм»", "Магазин «Обувь»":
		return []string{"Магазин одежды"}
	case "Автосалон":
		return []string{"Автосалон"}
	case "Гипермаркет (продовольственный)":
		return []string{"Гипермаркет"}
	default:
		return []string{"Магазин"}
	}
}

func resolveFoodCategory(typeObj string) []string {
	switch typeObj {
	case "ресторан":
		return []string{"Ресторан"}
	case "бар":
		return []string{"Бар"}
	case "предприятие быстрого обслуживания":
		return []string{"Фаст-фуд"}
	case "кафе":
		return []string{"Кафе"}
	case "кафетерий":
		return []string{"Кафетерий"}
	case "закусочная":
		return []string{"Закусочная"}
	case "магазин (отдел кулинарии)":
		return []string{"Магазин еды"}
	case "ночной клуб (дискотека)":
		return []string{"Ночной клуб"}
	default:
		return nil // skip unknown types
	}
}

// ---- Event Places ----

func (s *ParseService) initEventPlaces(ctx context.Context) error {
	if err := s.eventPlaceRepo.Clear(ctx); err != nil {
		return err
	}

	if s.cfg.TimepadKey == "" {
		log.Println("[parseService] TIMEPAD key not set, skipping events")
		return nil
	}

	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	startsMin := now.Format("2006-01-02")
	startsMax := tomorrow.Format("2006-01-02")

	limit := 59
	maxCount := 200
	iterations := (maxCount + limit - 1) / limit

	for i := 0; i < iterations; i++ {
		count := limit
		if maxCount < (i+1)*limit {
			count = maxCount % limit
		}

		url := fmt.Sprintf(
			"https://api.timepad.ru/v1/events?limit=%d&cities=Москва&starts_at_min=%s&starts_at_max=%s",
			count, startsMin, startsMax,
		)

		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+s.cfg.TimepadKey)

		resp, err := s.httpClient.Do(req)
		if err != nil {
			log.Printf("[parseService] timepad request error: %v", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var result struct {
			Values []struct {
				ID int `json:"id"`
			} `json:"values"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}

		for _, val := range result.Values {
			s.addEvent(ctx, val.ID)
		}
	}
	return nil
}

func (s *ParseService) addEvent(ctx context.Context, eventID int) {
	url := fmt.Sprintf("https://api.timepad.ru/v1/events/%d", eventID)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+s.cfg.TimepadKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Name             string `json:"name"`
		DescriptionShort string `json:"description_short"`
		StartsAt         string `json:"starts_at"`
		EndsAt           string `json:"ends_at"`
		Location         struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"location"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	if len(data.Location.Coordinates) < 2 {
		return
	}

	event := &domain.EventPlace{
		Name:       data.Name,
		Desc:       data.DescriptionShort,
		Latitude:   data.Location.Coordinates[0],
		Longitude:  data.Location.Coordinates[1],
		StartTime:  data.StartsAt,
		FinishTime: data.EndsAt,
	}

	_ = s.eventPlaceRepo.Save(ctx, event)
}

// ---- Helpers ----

func extractPoint(coords any) (lat, lon float64) {
	if arr, ok := coords.([]any); ok && len(arr) >= 2 {
		lat = toFloat64(arr[0])
		lon = toFloat64(arr[1])
	}
	return
}

func navigateToBase(feature map[string]any, basePath []string) map[string]any {
	current := feature
	for _, key := range basePath {
		if m, ok := current[key].(map[string]any); ok {
			current = m
		} else {
			return map[string]any{}
		}
	}
	return current
}

func getMap(m map[string]any, key string) map[string]any {
	if v, ok := m[key].(map[string]any); ok {
		return v
	}
	return map[string]any{}
}

func setPlaceField(place *domain.Place, fieldKey string, val any) {
	if val == nil {
		return
	}
	switch fieldKey {
	case "name":
		place.Name = toString(val)
	case "email":
		place.Email = toString(val)
	case "website":
		place.Website = toString(val)
	case "phone":
		place.Phone = toString(val)
	case "schedule":
		place.Schedule = val
	case "isPaid":
		b := toBool(val)
		place.IsPaid = &b
	case "price":
		place.Price = toString(val)
	case "address":
		place.Address = toString(val)
	}
}

func toString(val any) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat64(val any) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		return 0
	}
}

func toBool(val any) bool {
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "да" || v == "true" || v == "1" || v == "Да" || v == "есть" || v == "Есть"
	case float64:
		return v != 0
	default:
		return false
	}
}
