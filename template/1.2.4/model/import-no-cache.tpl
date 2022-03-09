
import (
    "database/sql"
	"github.com/yfanswer/zero-gorm/db"
	{{if .time}}"time"{{end}}

	perrors "github.com/pkg/errors"
	"gorm.io/gorm"
)
