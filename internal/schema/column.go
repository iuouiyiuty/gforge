package schema

const (
	cUnderScore = "_"
)

var (
	typeWrappers = []typeWrapper{i64TypeWrapper, byteTypeWrapper, intTypeWrapper, float64TypeWrapper, stringTypeWrapper, timeTypeWrapper}
)

// Column stands for a column of a table
type column struct {
	Name       string `json:"COLUMN_NAME"`
	Type       string `json:"COLUMN_TYPE"`
	Comment    string `json:"COLUMN_COMMENT"`
	IsNullable string `json:"IS_NULLABLE"`
	Key        string `json:"COLUMN_KEY"`
	Extra      string `json:"EXTRA"`
}

//GetType returns which built in type the column should be in generated go code
func (c *column) GetType() (string, error) {
	t := getType(c.Type)
	if "" == t {
		return "", errUnknownType(c.Name, c.Type)
	}
	if c.GetIsNullable() {
		switch t {
		case cTypeString:
			t = cNullString
		case cTypeInt, cTypeInt8, cTypeInt64, cTypeUInt, cTypeUInt64, cTypeByte:
			t = cNullInt
		case cTypeFloat64:
			t = cNullFloat
		case cTypeTime:
			t = cNullTime
		}
	}
	return t, nil
}

//GetName returns the Cammel Name of the struct
func (c *column) GetName() string {
	return convertUnderScoreToCammel(c.Name)
}

func (c *column) GetIsNullable() bool {
	return c.IsNullable == "YES"
}

func (c *column) IsPk() bool {
	return c.Key == "PRI"
}

func (c *column) IsAutoIncr() bool {
	return c.Extra == "auto_increment"
}

func getType(t string) string {
	for _, wrapper := range typeWrappers {
		typer := wrapper(t)
		if typer.Match() {
			return typer.Type()
		}
	}
	return ""
}
