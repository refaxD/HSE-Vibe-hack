package mongo

import (
	"context"
	"errors"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type placeRepo struct {
	col *mongo.Collection
}

func NewPlaceRepository(db *mongo.Database) domain.PlaceRepository {
	return &placeRepo{col: db.Collection("place")}
}

func (r *placeRepo) FindAll(ctx context.Context) ([]domain.Place, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	places := make([]domain.Place, 0)
	if err := cursor.All(ctx, &places); err != nil {
		return nil, err
	}
	return places, nil
}

func (r *placeRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Place, error) {
	var place domain.Place
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&place)
	if err != nil {
		return nil, err
	}
	return &place, nil
}

func (r *placeRepo) FindByName(ctx context.Context, name string) ([]domain.Place, error) {
	cursor, err := r.col.Find(ctx, bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	places := make([]domain.Place, 0)
	if err := cursor.All(ctx, &places); err != nil {
		return nil, err
	}
	return places, nil
}

func (r *placeRepo) Save(ctx context.Context, place *domain.Place) error {
	if place.ID.IsZero() {
		place.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, place)
	return err
}

// InsertMany вставляет места с ordered: false. При дубликате по (name, lat, lon) запись пропускается.
// Возвращает ID только реально вставленных документов (для продолжения после паузы).
func (r *placeRepo) InsertMany(ctx context.Context, places []*domain.Place) ([]primitive.ObjectID, error) {
	if len(places) == 0 {
		return nil, nil
	}
	docs := make([]interface{}, len(places))
	for i, p := range places {
		if p.ID.IsZero() {
			p.ID = primitive.NewObjectID()
		}
		docs[i] = p
	}
	opts := options.InsertMany().SetOrdered(false)
	result, err := r.col.InsertMany(ctx, docs, opts)
	if err == nil {
		ids := make([]primitive.ObjectID, 0, len(result.InsertedIDs))
		for _, v := range result.InsertedIDs {
			if oid, ok := v.(primitive.ObjectID); ok {
				ids = append(ids, oid)
			}
		}
		return ids, nil
	}
	// При частичном успехе (дубликаты) драйвер не отдаёт InsertedIDs — вставляем по одному и собираем ID.
	var bwe mongo.BulkWriteException
	if !errors.As(err, &bwe) {
		return nil, err
	}
	ids := make([]primitive.ObjectID, 0, len(places))
	for _, p := range places {
		_, err := r.col.InsertOne(ctx, p)
		if err == nil {
			ids = append(ids, p.ID)
			continue
		}
		var we mongo.WriteException
		if errors.As(err, &we) && len(we.WriteErrors) > 0 && we.WriteErrors[0].Code == 11000 {
			continue // дубликат — пропускаем
		}
		return ids, err // иная ошибка
	}
	return ids, nil
}

// EnsurePlaceUniqueIndex создаёт уникальный индекс по name+latitude+longitude для защиты от дубликатов при повторном запуске.
func (r *placeRepo) EnsurePlaceUniqueIndex(ctx context.Context) error {
	idx := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "latitude", Value: 1},
			{Key: "longitude", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("place_name_lat_lon_unique"),
	}
	_, err := r.col.Indexes().CreateOne(ctx, idx)
	// При повторном запуске индекс может уже существовать (85/86) — не прерываем парсинг
	if err != nil {
		return nil
	}
	return nil
}

func (r *placeRepo) SetCategories(ctx context.Context, id primitive.ObjectID, categoryIDs []primitive.ObjectID) error {
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"categories": categoryIDs}})
	return err
}

func (r *placeRepo) Clear(ctx context.Context) error {
	return r.col.Drop(ctx)
}
