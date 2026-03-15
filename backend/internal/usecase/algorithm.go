package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
	"github.com/hse-vibe-hack/backend/pkg/config"
)

const basePrompt = `### ИНСТРУКЦИЯ ДЛЯ НЕЙРОСЕТИ

Анализируй пользовательский запрос и преобразуй его в JSON-структуру по правилам:

1. **Распознавание тегов**:
   - ` + "`<#:fixed:[lon]:[lat]:[name]>`" + ` → тип ` + "`fixed`" + ` (статичная точка)
     - Поля: [` + "`lon`" + ` , ` + "`lat`" + `] (список чисел), ` + "`name`" + ` (строка), ` + "`isPivotPoint`" + ` (true/false)
   - ` + "`<#:embeding:[name]>`" + ` → тип ` + "`embeding`" + ` (группа точек)
     - Поле: ` + "`name`" + ` (строка)
   - ` + "`<@:[name]>`" + ` → тип ` + "`event`" + ` (событие)
     - Поле: ` + "`name`" + ` (строка)

2. **Обработка текста без тегов**:
   → Если есть прямое упоминание категории → тип ` + "`category`" + ` (см. список ниже)
   → Иначе → выполни операцию ImagineRoute

3. Операция ImagineRoute:
   • Сгенерируй уточняющий подпромпт на основе изначального raw_prompt этого подпромпта, в котором придумай возможный вариант маршрута по местам из списка категорий и их сочетаний, подходящий под запрос (допустимо разбиение на несколько последовательных мест, последовательный маршрут). Сгенерированный промпт может быть написан как сценарий прогулки, где места идут по-порядку, или на выбр одно из нескольких, в зависимости от условий запроса пользователя
   • Рекурсивно примени эти же правила к сгенерированному тексту и добавь все места в изначальный промпт. Разные части подмаршрута помечай как отдельную локацию

4. **Категории (ТОЛЬКО из списка)**:
Бассейн
Спорт площадка
Спорт зал
Катание на лошадях
Теннис
Каток
Тир
Пейнтбол
Футбол
Регби
Скалодром
Ресторан
Бар
Фаст-фуд
Кафе
Кафетерий
Закусочная
Магазин еды
Ночной клуб
Парк
Ботанический сад
Пикник
Кинотеатр
Клуб
Аквапарк
Игровая площадка
Детский технопарк
Досуг
Аттракционы
Библиотека
Выставка
Музей
Театр
Дом культуры
Концертный зал
Мечеть
Католический храм
Монастырь
Синагога
Православная церковь
Коворкинг
Гипермаркет
Дискаунтер
Промтовары
Магазин
Алкомаркет
Магазин цветов
Магазин одежды
Супермаркет
Автосалон

5. **Формат ответа**:
   - Сохраняй исходный порядок элементов относительно текста запроса в итовых данных.
   - Ответ должен быть предоставлен в формате json по указанной спецификации в виде обычного текста (без форматирования ` + "```json...```" + `)

---

### ВХОДНЫЕ ДАННЫЕ (пример)
"Вечерний маршрут: ужин в <#:embeding:TopRestaurants>, затем культурная программа"

---

### ОЖИДАЕМЫЙ ОТВЕТ (шаблон)
[
    {
      "type": "embeding",
      "raw_prompt": "<#:embeding:TopRestaurants>",
      "name": "TopRestaurants"
    },
    {
      "type": "route",
      "raw_prompt": "культурная программа",
      "generated_prompt": "Посещение театра или музея вечером",
      "parsed_elements": [
        {
          "type": "category",
          "categories": {"Театр": 60, "Музей": 40}
        }
      ]
    }
  ]

---

### ВХОДНЫЕ ДАННЫЕ (пример)
"Прогулка может включать посещение парка с живописными видами, ужин в романтическом ресторане и завершиться посещением аквапарка для совместного отдыха"

---

### ОЖИДАЕМЫЙ ОТВЕТ (шаблон)

[{"type":"route","raw_prompt":"Романтическая прогулка по москве","generated_prompt":"Прогулка может включать посещение парка с живописными видами и завершится ужином в романтическом ресторане","parsed_elements":[{"type":"category","raw_prompt":"парк с живописными видами","categories":{"Парк":100}},{"type":"category","raw_prompt":"ужином в романтическом ресторане","categories":{"Ресторан":100}}]},{"type":"category","raw_prompt":"с посещением аквапарка","categories":{"Аквапарк":100}}]

---

### ВХОДНЫЕ ДАННЫЕ (пример)
"Сходить в кино с другом, а потом посетить кафе с коворкингом"

---

### ОЖИДАЕМЫЙ ОТВЕТ (шаблон)

[{"type":"category","raw_prompt":"Сходить в кино с другом","categories":{"Кинотеатр":100}},{"type":"category","raw_prompt":"посетить кафе с коворкингом","categories":{"Кафе":60, "Коворкинг": 70}}]

### ВХОДНЫЕ ДАННЫЕ

{USER_INPUT}

### ОТВЕТ`

type LocationItem struct {
	Location      any // *domain.Place or *domain.EventPlace
	Index         int
	CategoriesSum float64
	Latitude      float64
	Longitude     float64
}

type AlgorithmService struct {
	prompt        string
	startPosition []float64
	keys          []domain.PromptElement
	items         [][]LocationItem
	categories    []domain.Category

	distancesMatrix [][][]float64
	durationsMatrix [][][]float64

	placeRepo domain.PlaceRepository
	eventRepo domain.EventPlaceRepository
	catRepo   domain.CategoryRepository
	cfg       *config.Config
	rng       *rand.Rand
}

const (
	iterations              = 1_000_000
	temperatureMultiplier   = 0.9
	maxItemsCount           = 29
	maxItemsPreCount        = 29
	changesCount = 1
)

func NewAlgorithmService(
	prompt string,
	startPosition []float64,
	placeRepo domain.PlaceRepository,
	eventRepo domain.EventPlaceRepository,
	catRepo domain.CategoryRepository,
	cfg *config.Config,
) *AlgorithmService {
	return &AlgorithmService{
		prompt:        prompt,
		startPosition: startPosition,
		placeRepo:     placeRepo,
		eventRepo:     eventRepo,
		catRepo:       catRepo,
		cfg:           cfg,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (a *AlgorithmService) Generate(ctx context.Context) ([]any, error) {
	var err error
	a.categories, err = a.catRepo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.InternalServer()
	}

	start := time.Now()
	log.Println("[algorithm] generation has been started")

	if err := a.parsePrompt(ctx); err != nil {
		return nil, err
	}
	log.Printf("[algorithm] prompt parsed: %v", time.Since(start))
	start = time.Now()

	if err := a.findPointsForPrompt(ctx); err != nil {
		return nil, err
	}
	log.Printf("[algorithm] points found: %v", time.Since(start))
	start = time.Now()

	if err := a.calculateDistances(); err != nil {
		return nil, err
	}
	log.Printf("[algorithm] distances calculated: %v", time.Since(start))
	start = time.Now()

	result := a.annealing()
	log.Printf("[algorithm] annealing finished: %v", time.Since(start))

	path := make([]any, len(result))
	for i, item := range result {
		path[i] = item.Location
	}
	return path, nil
}

func (a *AlgorithmService) parsePrompt(ctx context.Context) error {
	promptText := strings.Replace(basePrompt, "{USER_INPUT}", a.prompt, 1)

	reqBody := map[string]any{
		"model": "deepseek/deepseek-r1",
		"messages": []map[string]string{
			{"role": "user", "content": promptText},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return apperrors.InternalServer("failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.cfg.OpenAIAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return apperrors.InternalServer("AI request failed")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[algorithm] AI response: %s", string(respBody))

	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &aiResp); err != nil || len(aiResp.Choices) == 0 {
		return apperrors.InternalServer("failed to parse AI response")
	}

	content := aiResp.Choices[0].Message.Content

	var elements []domain.PromptElement
	if err := json.Unmarshal([]byte(content), &elements); err != nil {
		return apperrors.InternalServer("failed to parse AI JSON")
	}

	var parsed []domain.PromptElement
	for _, el := range elements {
		if el.Type == "route" {
			parsed = append(parsed, el.ParsedElements...)
		} else {
			parsed = append(parsed, el)
		}
	}

	log.Printf("[algorithm] parsed elements: %+v", parsed)
	a.keys = parsed
	return nil
}

func (a *AlgorithmService) generateStartPlacement() []LocationItem {
	result := make([]LocationItem, len(a.items))
	for i := range a.items {
		idx := a.rng.Intn(len(a.items[i]))
		result[i] = a.items[i][idx]
	}
	return result
}

func (a *AlgorithmService) randomChangePlacement(placement []LocationItem) []LocationItem {
	result := make([]LocationItem, len(placement))
	copy(result, placement)

	for i := 0; i < min(len(result), changesCount); i++ {
		idx := a.rng.Intn(len(result))
		for len(a.items[idx]) == 1 {
			idx = a.rng.Intn(len(result))
		}
		itemIdx := a.rng.Intn(len(a.items[idx]))
		result[idx] = a.items[idx][itemIdx]
	}
	return result
}

func (a *AlgorithmService) annealing() []LocationItem {
	placement := a.generateStartPlacement()
	resultError := a.errorFunction(placement)

	temperature := 1.0
	for i := 0; i < iterations; i++ {
		temperature *= temperatureMultiplier

		currentPlacement := a.randomChangePlacement(placement)
		currentError := a.errorFunction(currentPlacement)

		if currentError > resultError || a.rng.Float64() < math.Exp((currentError-resultError)/temperature) {
			resultError = currentError
			placement = currentPlacement
		}
	}
	return placement
}

func (a *AlgorithmService) errorFunction(items []LocationItem) float64 {
	if len(items) < 2 {
		return 1
	}

	// distance indicator
	avgDistance := 0.0
	minDistance := math.MaxFloat64
	maxDistance := 0.0
	for i := 1; i < len(items); i++ {
		if i-1 >= len(a.distancesMatrix) {
			continue
		}
		dm := a.distancesMatrix[i-1]
		prevIdx := items[i-1].Index
		nextIdx := len(a.items[i-1]) + items[i].Index
		if prevIdx >= len(dm) || nextIdx >= len(dm[prevIdx]) {
			continue
		}
		meters := dm[prevIdx][nextIdx]
		km := math.Max(1, math.Ceil(meters/1000))
		avgDistance += km
		minDistance = math.Min(minDistance, km)
		maxDistance = math.Max(maxDistance, km)
	}
	avgDistance /= float64(len(items) - 1)
	distanceIndicator := avgDistance * math.Max(1, maxDistance-minDistance)
	if distanceIndicator == 0 {
		return 1
	}

	// beauty
	beauty := 0.0
	for idx, item := range items {
		if idx < len(a.keys) {
			if a.keys[idx].Type == "category" {
				beauty += item.CategoriesSum
			} else if a.keys[idx].Type == "fixed" || a.keys[idx].Type == "event" {
				beauty += 100
			}
		}
	}
	if math.IsNaN(beauty) {
		return 1
	}

	return beauty * (1e4 / distanceIndicator) * 1e4
}

func (a *AlgorithmService) calculatePreSimilarity(place *domain.Place, _ int, categoryIndex int) (float64, float64) {
	similarity := 0.0
	categoriesSum := 0.0
	promptElement := a.keys[categoryIndex]

	// categories similarity
	itemCategories := make([]string, 0)
	for _, objID := range place.Categories {
		idStr := objID.Hex()
		for _, cat := range a.categories {
			if cat.ID.Hex() == idStr {
				itemCategories = append(itemCategories, strings.ToLower(cat.Name))
				break
			}
		}
	}

	for key, pct := range promptElement.Categories {
		for _, ic := range itemCategories {
			if strings.ToLower(key) == ic {
				similarity += pct
				categoriesSum += pct
				break
			}
		}
	}

	// distance to start position
	getDistance := func(lat1, lon1, lat2, lon2 float64) float64 {
		lat1i := math.Ceil(lat1 * 1e3)
		lon1i := math.Ceil(lon1 * 1e3)
		lat2i := math.Ceil(lat2 * 1e3)
		lon2i := math.Ceil(lon2 * 1e3)
		pyth := math.Sqrt((lat1i-lat2i)*(lat1i-lat2i) + (lon1i-lon2i)*(lon1i-lon2i))
		manh := math.Abs(lat1i-lat2i) + math.Abs(lon1i-lon2i)
		return (pyth + manh) / 2
	}

	distance := getDistance(place.Latitude, place.Longitude, a.startPosition[0], a.startPosition[1])
	similarity = similarity * (1e5 / math.Max(1, distance))

	return similarity, categoriesSum
}

func (a *AlgorithmService) findPointsForPrompt(ctx context.Context) error {
	a.items = make([][]LocationItem, len(a.keys))
	for i := range a.items {
		a.items[i] = make([]LocationItem, 0)
	}

	// find events and fixed points
	for i, key := range a.keys {
		if key.Type == "event" {
			events, err := a.eventRepo.FindByName(ctx, key.Name)
			if err != nil {
				continue
			}
			for _, ev := range events {
				if len(key.Coords) >= 2 && ev.Latitude == key.Coords[0] && ev.Longitude == key.Coords[1] {
					a.items[i] = append(a.items[i], LocationItem{
						Location:  &ev,
						Index:     0,
						Latitude:  ev.Latitude,
						Longitude: ev.Longitude,
					})
				}
			}
		}
		if key.Type == "fixed" {
			places, err := a.placeRepo.FindByName(ctx, key.Name)
			if err != nil {
				continue
			}
			for _, p := range places {
				if len(key.Coords) >= 2 && p.Latitude == key.Coords[0] && p.Longitude == key.Coords[1] {
					place := p
					a.items[i] = append(a.items[i], LocationItem{
						Location:  &place,
						Index:     0,
						Latitude:  place.Latitude,
						Longitude: place.Longitude,
					})
					break
				}
			}
		}
	}

	// find category places
	allPlaces, err := a.placeRepo.FindAll(ctx)
	if err != nil {
		return apperrors.InternalServer()
	}

	for i, key := range a.keys {
		if key.Type != "category" {
			continue
		}

		type sortPlace struct {
			place         *domain.Place
			similarity    float64
			categoriesSum float64
		}

		var sorted []sortPlace
		for j := range allPlaces {
			p := &allPlaces[j]
			if !math.IsNaN(p.Longitude) && !math.IsNaN(p.Latitude) {
				sim, catSum := a.calculatePreSimilarity(p, j, i)
				sorted = append(sorted, sortPlace{place: p, similarity: sim, categoriesSum: catSum})
			}
		}

		// sort by similarity descending
		for x := 0; x < len(sorted)-1; x++ {
			for y := x + 1; y < len(sorted); y++ {
				if sorted[y].similarity > sorted[x].similarity {
					sorted[x], sorted[y] = sorted[y], sorted[x]
				}
			}
		}

		if len(sorted) > maxItemsPreCount {
			sorted = sorted[:maxItemsPreCount]
		}

		if len(sorted) == 0 || sorted[0].categoriesSum == 0 {
			continue
		}

		// take top places
		limit := maxItemsCount
		if limit > len(sorted) {
			limit = len(sorted)
		}
		for j := 0; j < limit; j++ {
			a.items[i] = append(a.items[i], LocationItem{
				Location:      sorted[j].place,
				Index:         0,
				CategoriesSum: sorted[j].categoriesSum,
				Latitude:      sorted[j].place.Latitude,
				Longitude:     sorted[j].place.Longitude,
			})
		}
	}

	// set indexes
	for i := range a.items {
		for j := range a.items[i] {
			a.items[i][j].Index = j
		}
	}

	return nil
}

func (a *AlgorithmService) calculateDistances() error {
	a.distancesMatrix = make([][][]float64, 0)
	a.durationsMatrix = make([][][]float64, 0)

	for i := 1; i < len(a.items); i++ {
		locations := make([][]float64, 0)
		for _, item := range a.items[i-1] {
			locations = append(locations, []float64{item.Latitude, item.Longitude})
		}
		for _, item := range a.items[i] {
			locations = append(locations, []float64{item.Latitude, item.Longitude})
		}

		reqBody := map[string]any{
			"locations": locations,
			"metrics":   []string{"distance", "duration"},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/matrix/foot-walking", a.cfg.RouterURL), bytes.NewReader(bodyBytes))
		if err != nil {
			return apperrors.InternalServer("failed to create ORS request")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", a.cfg.RouterKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return apperrors.InternalServer("ORS request failed")
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var data struct {
			Distances [][]float64 `json:"distances"`
			Durations [][]float64 `json:"durations"`
		}
		if err := json.Unmarshal(respBody, &data); err != nil {
			return apperrors.InternalServer("failed to parse ORS response")
		}

		a.distancesMatrix = append(a.distancesMatrix, data.Distances)
		a.durationsMatrix = append(a.durationsMatrix, data.Durations)
	}

	return nil
}
