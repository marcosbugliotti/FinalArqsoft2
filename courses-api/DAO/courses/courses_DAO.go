package courses

type Course struct {
	ID           int64   `bson:"id"`
	Name         string  `bson:"name"`
	Description  string  `bson:"description"`
	Category     string  `bson:"category"`
	Duration     string  `bson:"duration"`
	InstructorID int64   `bson:"instructor_id"`
	ImageID      string  `bson:"image_id"`
	Capacity     int     `bson:"capacity"`
	Rating       float64 `bson:"rating"`
}
