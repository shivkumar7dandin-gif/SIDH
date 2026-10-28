package repository

import (
	"context"
	"errors"
	"strings"

	subjectModel "github.com/shivkumar7dandin-gif/students-api/internal/subject/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubjectRepository struct {
	collection *mongo.Collection
}

func NewSubjectRepository(db *mongo.Database) *SubjectRepository {
	return &SubjectRepository{
		collection: db.Collection("subjects"),
	}
}

func (r *SubjectRepository) Create(
	ctx context.Context,
	subject *subjectModel.Subject,
) error {
	_, err := r.collection.InsertOne(ctx, subject)
	return err
}

func (r *SubjectRepository) GetAllByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]subjectModel.Subject, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"college_id": collegeID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subjects []subjectModel.Subject

	if err := cursor.All(ctx, &subjects); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *SubjectRepository) GetByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) (*subjectModel.Subject, error) {
	var subject subjectModel.Subject

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id":        id,
			"college_id": collegeID,
		},
	).Decode(&subject)

	if err != nil {
		return nil, err
	}

	return &subject, nil
}

func (r *SubjectRepository) ExistsByName(
	ctx context.Context,
	collegeID bson.ObjectID,
	name string,
) (bool, error) {
	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{
			"college_id": collegeID,
			"name": bson.M{
				"$regex":   "^" + escapeRegex(strings.TrimSpace(name)) + "$",
				"$options": "i",
			},
		},
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *SubjectRepository) Update(
	ctx context.Context,
	subject *subjectModel.Subject,
) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":        subject.ID,
			"college_id": subject.CollegeID,
		},
		bson.M{
			"$set": bson.M{
				"name":       subject.Name,
				"code":       subject.Code,
				"status":     subject.Status,
				"updated_at": subject.UpdatedAt,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *SubjectRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) error {
	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{
			"_id":        id,
			"college_id": collegeID,
		},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func escapeRegex(value string) string {
	replacer := strings.NewReplacer(
		`\\`, `\\\\`,
		`.`, `\.`,
		`+`, `\+`,
		`*`, `\*`,
		`?`, `\?`,
		`(`, `\(`,
		`)`, `\)`,
		`[`, `\[`,
		`]`, `\]`,
		`{`, `\{`,
		`}`, `\}`,
		`^`, `\^`,
		`$`, `\$`,
		`|`, `\|`,
	)

	return replacer.Replace(value)
}

var ErrSubjectNotFound = errors.New("subject not found")
