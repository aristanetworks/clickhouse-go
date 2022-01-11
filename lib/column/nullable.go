package column

import (
	"reflect"

	"github.com/ClickHouse/clickhouse-go/v2/lib/binary"
)

type Nullable struct {
	base  Interface
	nulls UInt8
}

func (col *Nullable) Rows() int {
	return len(col.nulls)
}

func (c *Nullable) Decode(decoder *binary.Decoder, rows int) (err error) {
	if err := c.nulls.Decode(decoder, rows); err != nil {
		return err
	}
	if err := c.base.Decode(decoder, rows); err != nil {
		return err
	}
	return nil
}

func (c *Nullable) RowValue(row int) interface{} {
	if c.nulls[row] == 1 {
		return nil
	}
	return c.base.RowValue(row)
}

func (c *Nullable) ScanRow(dest interface{}, row int) error {
	if len(c.nulls) < row {

	}
	if c.nulls[row] == 1 {
		return nil
	}
	return c.base.ScanRow(dest, row)
}

func (col *Nullable) Append(v interface{}) error {
	// @todo: rm reflect and add Column.AppendWithNulls ([]null, error)
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i).Interface()
			switch {
			case s.Index(i).IsNil():
				col.nulls = append(col.nulls, 1)
				if err := col.base.AppendRow(null{}); err != nil {
					return err
				}
			default:
				col.nulls = append(col.nulls, 0)
				if err := col.base.AppendRow(v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (col *Nullable) AppendRow(v interface{}) error {
	switch v := v.(type) {
	case nil:
		col.nulls = append(col.nulls, 1)
		return col.base.AppendRow(null{})
	default:
		col.nulls = append(col.nulls, 0)
		return col.base.AppendRow(v)
	}
}

func (col *Nullable) Encode(encoder *binary.Encoder) error {
	if err := col.nulls.Encode(encoder); err != nil {
		return err
	}
	if err := col.base.Encode(encoder); err != nil {
		return err
	}
	return nil
}

var _ Interface = (*Nullable)(nil)
