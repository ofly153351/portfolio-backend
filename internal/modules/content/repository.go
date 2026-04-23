package content

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Ping(ctx context.Context) error
	GetByLocaleStatus(ctx context.Context, locale, status string) (HistoryItem, bool, error)
	UpsertContent(ctx context.Context, item HistoryItem, status string) error
	AppendHistory(ctx context.Context, item HistoryItem) error
	ListHistory(ctx context.Context, locale string) ([]HistoryItem, error)
}

type MongoRepository struct {
	db        *mongo.Database
	projects  *mongo.Collection
	technical *mongo.Collection
	info      *mongo.Collection
	history   *mongo.Collection
}

type sectionProjectsDoc struct {
	Locale    string        `bson:"locale"`
	Status    string        `bson:"status"`
	Version   int           `bson:"version"`
	Projects  []ProjectItem `bson:"projects_json"`
	UpdatedBy string        `bson:"updated_by"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type sectionTechnicalDoc struct {
	Locale    string          `bson:"locale"`
	Status    string          `bson:"status"`
	Version   int             `bson:"version"`
	Technical []TechnicalItem `bson:"technical_json"`
	UpdatedBy string          `bson:"updated_by"`
	UpdatedAt time.Time       `bson:"updated_at"`
}

type sectionInfoDoc struct {
	Locale    string        `bson:"locale"`
	Status    string        `bson:"status"`
	Version   int           `bson:"version"`
	Info      PortfolioInfo `bson:"info_json"`
	UpdatedBy string        `bson:"updated_by"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type historyDoc struct {
	Locale    string      `bson:"locale"`
	Version   int         `bson:"version"`
	Content   ContentBody `bson:"content_json"`
	UpdatedBy string      `bson:"updated_by"`
	UpdatedAt time.Time   `bson:"updated_at"`
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	if db == nil {
		return nil
	}
	return &MongoRepository{
		db:        db,
		projects:  db.Collection("portfolio_projects"),
		technical: db.Collection("portfolio_technical"),
		info:      db.Collection("portfolio_info"),
		history:   db.Collection("portfolio_content_history"),
	}
}

func (r *MongoRepository) Ping(ctx context.Context) error {
	if r == nil || r.db == nil {
		return errors.New("content repository unavailable")
	}
	return r.db.Client().Ping(ctx, nil)
}

func (r *MongoRepository) GetByLocaleStatus(ctx context.Context, locale, status string) (HistoryItem, bool, error) {
	if r == nil || r.db == nil {
		return HistoryItem{}, false, errors.New("content repository unavailable")
	}
	filter := bson.M{"locale": locale, "status": status}

	var (
		pDoc sectionProjectsDoc
		tDoc sectionTechnicalDoc
		iDoc sectionInfoDoc
	)

	pFound, err := findOne(ctx, r.projects, filter, &pDoc)
	if err != nil {
		return HistoryItem{}, false, err
	}
	tFound, err := findOne(ctx, r.technical, filter, &tDoc)
	if err != nil {
		return HistoryItem{}, false, err
	}
	iFound, err := findOne(ctx, r.info, filter, &iDoc)
	if err != nil {
		return HistoryItem{}, false, err
	}
	if !pFound && !tFound && !iFound {
		return HistoryItem{}, false, nil
	}

	out := HistoryItem{
		Locale: locale,
		Content: ContentBody{
			Technical:     []TechnicalItem{},
			Projects:      []ProjectItem{},
			PortfolioInfo: PortfolioInfo{},
		},
	}
	if pFound {
		out.Content.Projects = pDoc.Projects
	}
	if tFound {
		out.Content.Technical = tDoc.Technical
	}
	if iFound {
		out.Content.PortfolioInfo = iDoc.Info
	}

	// Keep metadata consistent even if only one section exists.
	out.Version, out.UpdatedBy, out.UpdatedAt = latestMetaFromSections(pDoc, pFound, tDoc, tFound, iDoc, iFound)
	if out.Version <= 0 {
		out.Version = 1
	}
	if out.UpdatedAt.IsZero() {
		out.UpdatedAt = time.Now().UTC()
	}

	return out, true, nil
}

func (r *MongoRepository) UpsertContent(ctx context.Context, item HistoryItem, status string) error {
	if r == nil || r.db == nil {
		return errors.New("content repository unavailable")
	}
	filter := bson.M{"locale": item.Locale, "status": status}

	projectSet := bson.M{
		"locale":        item.Locale,
		"status":        status,
		"version":       item.Version,
		"projects_json": item.Content.Projects,
		"updated_by":    item.UpdatedBy,
		"updated_at":    item.UpdatedAt,
	}
	technicalSet := bson.M{
		"locale":         item.Locale,
		"status":         status,
		"version":        item.Version,
		"technical_json": item.Content.Technical,
		"updated_by":     item.UpdatedBy,
		"updated_at":     item.UpdatedAt,
	}
	infoSet := bson.M{
		"locale":     item.Locale,
		"status":     status,
		"version":    item.Version,
		"info_json":  item.Content.PortfolioInfo,
		"updated_by": item.UpdatedBy,
		"updated_at": item.UpdatedAt,
	}

	if _, err := r.projects.UpdateOne(ctx, filter, bson.M{"$set": projectSet}, options.Update().SetUpsert(true)); err != nil {
		return err
	}
	if _, err := r.technical.UpdateOne(ctx, filter, bson.M{"$set": technicalSet}, options.Update().SetUpsert(true)); err != nil {
		return err
	}
	if _, err := r.info.UpdateOne(ctx, filter, bson.M{"$set": infoSet}, options.Update().SetUpsert(true)); err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) AppendHistory(ctx context.Context, item HistoryItem) error {
	if r == nil || r.history == nil {
		return errors.New("content repository unavailable")
	}
	_, err := r.history.InsertOne(ctx, bson.M{
		"locale":       item.Locale,
		"version":      item.Version,
		"content_json": item.Content,
		"updated_by":   item.UpdatedBy,
		"updated_at":   item.UpdatedAt,
	})
	return err
}

func (r *MongoRepository) ListHistory(ctx context.Context, locale string) ([]HistoryItem, error) {
	if r == nil || r.history == nil {
		return nil, errors.New("content repository unavailable")
	}
	cur, err := r.history.Find(ctx, bson.M{"locale": locale}, options.Find().SetSort(bson.D{{Key: "version", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]HistoryItem, 0)
	for cur.Next(ctx) {
		var doc historyDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, HistoryItem{
			Locale:    doc.Locale,
			Version:   doc.Version,
			Content:   doc.Content,
			UpdatedBy: doc.UpdatedBy,
			UpdatedAt: doc.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func findOne(ctx context.Context, col *mongo.Collection, filter interface{}, out interface{}) (bool, error) {
	if col == nil {
		return false, errors.New("collection unavailable")
	}
	err := col.FindOne(ctx, filter).Decode(out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func latestMetaFromSections(p sectionProjectsDoc, pFound bool, t sectionTechnicalDoc, tFound bool, i sectionInfoDoc, iFound bool) (int, string, time.Time) {
	version := 0
	updatedBy := ""
	updatedAt := time.Time{}

	consider := func(v int, by string, at time.Time) {
		if v > version {
			version = v
		}
		if at.After(updatedAt) {
			updatedAt = at
			updatedBy = by
		}
	}

	if pFound {
		consider(p.Version, p.UpdatedBy, p.UpdatedAt)
	}
	if tFound {
		consider(t.Version, t.UpdatedBy, t.UpdatedAt)
	}
	if iFound {
		consider(i.Version, i.UpdatedBy, i.UpdatedAt)
	}

	return version, updatedBy, updatedAt
}
