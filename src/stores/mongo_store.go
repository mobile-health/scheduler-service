package stores

type MongoStore struct {
}

func NewStore() Store {
	return &MongoStore{}
}
