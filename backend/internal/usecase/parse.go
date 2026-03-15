package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bounding box for Moscow (south, west, north, east).
const (
	moscowSouth = 55.49
	moscowWest  = 37.32
	moscowNorth = 55.92
	moscowEast  = 37.97
)

// bboxGridRows/Cols — разбиение Москвы на тайлы для батч-запросов к Overpass.
const bboxGridRows = 4
const bboxGridCols = 4

// placeBatchSize — размер батча для записи мест в БД.
const placeBatchSize = 200

// overpassURL is the public Overpass API endpoint.
const overpassURL = "https://overpass-api.de/api/interpreter"

// overpassDelay is a polite delay between consecutive Overpass requests.
const overpassDelay = 2 * time.Second

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
			Timeout: 120 * time.Second,
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
	categories := getDefaultCategories()
	created := 0
	for _, cat := range categories {
		if _, err := s.categoryRepo.FindByName(ctx, cat.Name); err == nil {
			// already exists
			continue
		}
		c := cat
		if err := s.categoryRepo.Save(ctx, &c); err != nil {
			log.Printf("[parseService] save category %q: %v", cat.Name, err)
			continue
		}
		created++
	}
	log.Printf("[parseService] categories: %d created, %d already existed", created, len(categories)-created)
	return nil
}

// ---- Places (OpenStreetMap / Overpass API) ----

type osmPlacesFile struct {
	Data []osmPlaceDef `json:"data"`
}

type osmPlaceDef struct {
	Name     string            `json:"name"`
	Category string            `json:"category"`
	Type     int               `json:"type"`
	Filters  map[string]string `json:"filters"`
}

// osmResponse is the top-level Overpass API JSON response.
type osmResponse struct {
	Elements []osmElement `json:"elements"`
}

// osmElement represents a single OSM node, way or relation.
type osmElement struct {
	Type   string            `json:"type"`
	Lat    float64           `json:"lat"`
	Lon    float64           `json:"lon"`
	Center *osmCenter        `json:"center"`
	Tags   map[string]string `json:"tags"`
}

type osmCenter struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (s *ParseService) initPlaces(ctx context.Context) error {
	if err := s.placeRepo.EnsurePlaceUniqueIndex(ctx); err != nil {
		log.Printf("[parseService] EnsurePlaceUniqueIndex: %v (continuing)", err)
	}
	data, err := os.ReadFile("internal/usecase/osm-places.json")
	if err != nil {
		return fmt.Errorf("read osm-places.json: %w", err)
	}

	var pf osmPlacesFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return fmt.Errorf("parse osm-places.json: %w", err)
	}

	for i, def := range pf.Data {
		if err := s.getPlacesFromOSM(ctx, def); err != nil {
			log.Printf("[parseService] error processing %d (%s): %v", i, def.Name, err)
		}

		// Polite delay between Overpass requests.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(overpassDelay):
		}
	}
	return nil
}

// retryDelays defines wait times between retries on HTTP 429.
var retryDelays = []time.Duration{30 * time.Second, 60 * time.Second, 120 * time.Second}

// getPlacesFromOSM запрашивает Overpass по тайлам Москвы и сохраняет результаты в БД батчами.
// При HTTP 429 повторяет запрос с увеличивающейся задержкой.
func (s *ParseService) getPlacesFromOSM(ctx context.Context, def osmPlaceDef) error {
	cat, err := s.categoryRepo.FindByName(ctx, def.Category)
	if err != nil {
		return fmt.Errorf("category %q: %w", def.Category, err)
	}

	tiles := moscowBBoxTiles()
	var batch []*domain.Place
	totalInserted := 0

	flushBatch := func() error {
		if len(batch) == 0 {
			return nil
		}
		insertedIDs, err := s.placeRepo.InsertMany(ctx, batch)
		if err != nil {
			return fmt.Errorf("InsertMany %d places: %w", len(batch), err)
		}
		if len(insertedIDs) > 0 {
			if err := s.categoryRepo.AddPlaces(ctx, cat.ID, insertedIDs); err != nil {
				return fmt.Errorf("AddPlaces: %w", err)
			}
		}
		totalInserted += len(insertedIDs)
		batch = batch[:0]
		return nil
	}

	for ti, tile := range tiles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		osmResp, err := s.fetchOverpassTile(ctx, def, tile)
		if err != nil {
			log.Printf("[parseService] overpass tile %d/%d %s: %v", ti+1, len(tiles), def.Name, err)
			continue
		}

		for _, el := range osmResp.Elements {
			lat, lon, ok := osmCoords(el)
			if !ok || (lat == 0 && lon == 0) {
				continue
			}
			name := osmName(el.Tags)
			if name == "" {
				continue
			}
			place := &domain.Place{
				ID:        primitive.NewObjectID(),
				Type:      def.Type,
				Latitude:  lat,
				Longitude: lon,
				Name:      name,
				Address:   osmAddress(el.Tags),
				Phone:     osmPhone(el.Tags),
				Website:   osmWebsite(el.Tags),
				Schedule:  osmOpeningHours(el.Tags),
				Categories: []primitive.ObjectID{cat.ID},
			}
			batch = append(batch, place)
			if len(batch) >= placeBatchSize {
				if err := flushBatch(); err != nil {
					return err
				}
			}
		}

		if ti < len(tiles)-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(overpassDelay):
			}
		}
	}

	if err := flushBatch(); err != nil {
		return err
	}
	log.Printf("[parseService] OSM %s (%s): inserted %d places in batches", def.Name, def.Category, totalInserted)
	return nil
}

// fetchOverpassTile выполняет один запрос к Overpass API для заданного тайла.
func (s *ParseService) fetchOverpassTile(ctx context.Context, def osmPlaceDef, b bbox) (*osmResponse, error) {
	query := buildOverpassQuery(def.Filters, b.string())
	var body []byte
	for attempt := 0; ; attempt++ {
		formData := url.Values{}
		formData.Set("data", query)
		req, err := http.NewRequestWithContext(ctx, "POST", overpassURL, strings.NewReader(formData.Encode()))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request: %w", err)
		}
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt >= len(retryDelays) {
				return nil, fmt.Errorf("rate limited after %d retries", attempt)
			}
			wait := retryDelays[attempt]
			log.Printf("[parseService] overpass 429 for %q, retry in %v", def.Name, wait)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			preview := string(body)
			if len(preview) > 200 {
				preview = preview[:200]
			}
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, preview)
		}
		break
	}

	var osmResp osmResponse
	if err := json.Unmarshal(body, &osmResp); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &osmResp, nil
}

// bbox represents a bounding box (south, west, north, east).
type bbox struct{ south, west, north, east float64 }

func (b bbox) string() string {
	return fmt.Sprintf("%.4f,%.4f,%.4f,%.4f", b.south, b.west, b.north, b.east)
}

// moscowBBoxTiles returns a slice of bboxes that partition Moscow.
func moscowBBoxTiles() []bbox {
	latStep := (moscowNorth - moscowSouth) / float64(bboxGridRows)
	lonStep := (moscowEast - moscowWest) / float64(bboxGridCols)
	var tiles []bbox
	for row := 0; row < bboxGridRows; row++ {
		for col := 0; col < bboxGridCols; col++ {
			south := moscowSouth + float64(row)*latStep
			west := moscowWest + float64(col)*lonStep
			tiles = append(tiles, bbox{
				south: south,
				west:  west,
				north: south + latStep,
				east:  west + lonStep,
			})
		}
	}
	return tiles
}

// buildOverpassQuery constructs an Overpass QL query for the given tag filters
// within the specified bounding box (south,west,north,east).
func buildOverpassQuery(filters map[string]string, bboxStr string) string {
	var filterStr strings.Builder
	for k, v := range filters {
		fmt.Fprintf(&filterStr, `["%s"="%s"]`, k, v)
	}
	f := filterStr.String()
	return fmt.Sprintf(
		`[out:json][timeout:120];(node%s(%s);way%s(%s);relation%s(%s););out center;`,
		f, bboxStr, f, bboxStr, f, bboxStr,
	)
}

// osmCoords returns lat/lon from a node (direct) or way/relation (center).
func osmCoords(el osmElement) (lat, lon float64, ok bool) {
	switch el.Type {
	case "node":
		return el.Lat, el.Lon, true
	case "way", "relation":
		if el.Center != nil {
			return el.Center.Lat, el.Center.Lon, true
		}
	}
	return 0, 0, false
}

func osmName(tags map[string]string) string {
	for _, k := range []string{"name:ru", "name"} {
		if v := tags[k]; v != "" {
			return v
		}
	}
	return ""
}

func osmAddress(tags map[string]string) string {
	if full := tags["addr:full"]; full != "" {
		return full
	}
	street := tags["addr:street"]
	num := tags["addr:housenumber"]
	switch {
	case street != "" && num != "":
		return street + ", " + num
	case street != "":
		return street
	}
	return ""
}

func osmPhone(tags map[string]string) string {
	for _, k := range []string{"phone", "contact:phone", "telephone"} {
		if v := tags[k]; v != "" {
			return v
		}
	}
	return ""
}

func osmWebsite(tags map[string]string) string {
	for _, k := range []string{"website", "contact:website", "url"} {
		if v := tags[k]; v != "" {
			return v
		}
	}
	return ""
}

func osmOpeningHours(tags map[string]string) any {
	if v := tags["opening_hours"]; v != "" {
		return v
	}
	return nil
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

		u := fmt.Sprintf(
			"https://api.timepad.ru/v1/events?limit=%d&cities=Москва&starts_at_min=%s&starts_at_max=%s",
			count, startsMin, startsMax,
		)

		req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
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
	u := fmt.Sprintf("https://api.timepad.ru/v1/events/%d", eventID)

	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
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

func toString(val any) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
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
