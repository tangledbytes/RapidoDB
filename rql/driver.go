package rql

import (
	"io"
	"time"
)

type db interface {
	Set(key string, data interface{}, expireIn time.Duration)
	Get(key string) (interface{}, bool)
}

type driver struct {
	writer io.Writer
	db     db
}

func (d *driver) Operate(src string) {
	// Parse the src
	ast, err := Parse(src)
	if err != nil {
		d.err(err.Error())
		return
	}

	for _, stmt := range ast.Statements {
		switch stmt.Typ {
		case SetType:
			d.db.Set(stmt.SetStatement.key, stmt.SetStatement.val, convertToDuration(stmt.SetStatement.exp))
		case GetType:
			// res := []string{}
			// for _, key := range stmt.GetStatement.keys {
			// 	r, ok := d.db.Get(key)
			// 	if ok {

			// 	}
			// }
		}
	}
}

func (d *driver) err(msg string) {
	d.writer.Write([]byte(msg))
}

func (d *driver) res(msg string) {
	d.writer.Write([]byte(msg))
}

// ============================ HELPER FUNCTIONS ===================================

func convertToDuration(t uint) time.Duration {
	return time.Duration(t * 1000)
}
