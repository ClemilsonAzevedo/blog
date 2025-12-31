package types

type Entity interface {
	TableName() string
	GetID() any
}
