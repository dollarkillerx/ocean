package enum

type SchemaType string

const (
	SchemaInt64     SchemaType = "int64"
	SchemaFloat64   SchemaType = "float64"
	SchemaString    SchemaType = "string"
	SchemaTimestamp SchemaType = "timestamp"
	SchemaBool      SchemaType = "bool"
)
