package decoder

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/Palaciodiego008/toonify/internal/types"
	"github.com/Palaciodiego008/toonify/internal/utils"
	"github.com/Palaciodiego008/toonify/parser"
)

// Decoder handles TOON decoding
type Decoder struct {
	opts *types.DecodeOptions
}

// New creates a new TOON decoder
func New(opts *types.DecodeOptions) *Decoder {
	if opts == nil {
		opts = types.DefaultDecodeOptions()
	}
	return &Decoder{opts: opts}
}

// Decode decodes TOON data into a Go value
func (d *Decoder) Decode(data []byte, v interface{}) error {
	parsed, err := parser.Parse(string(data), d.opts)
	if err != nil {
		return err
	}

	return d.assignValue(parsed, v)
}

func (d *Decoder) assignValue(src interface{}, dst interface{}) error {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return types.NewToonError("destination must be a pointer", 0, 0)
	}

	dstElem := dstValue.Elem()
	srcValue := reflect.ValueOf(src)

	return d.assignReflectValue(srcValue, dstElem)
}

func (d *Decoder) assignReflectValue(src, dst reflect.Value) error {
	if !src.IsValid() {
		// Handle nil values
		if dst.CanSet() {
			dst.Set(reflect.Zero(dst.Type()))
		}
		return nil
	}

	srcType := src.Type()
	dstType := dst.Type()

	// Handle interface{} destination
	if dstType.Kind() == reflect.Interface && dstType.NumMethod() == 0 {
		if dst.CanSet() {
			dst.Set(src)
		}
		return nil
	}

	// Handle pointer destination
	if dstType.Kind() == reflect.Ptr {
		if src.IsNil() {
			if dst.CanSet() {
				dst.Set(reflect.Zero(dstType))
			}
			return nil
		}

		if dst.IsNil() {
			dst.Set(reflect.New(dstType.Elem()))
		}
		return d.assignReflectValue(src, dst.Elem())
	}

	// Direct assignment if types match
	if srcType.AssignableTo(dstType) {
		if dst.CanSet() {
			dst.Set(src)
		}
		return nil
	}

	// Handle conversions
	switch dstType.Kind() {
	case reflect.String:
		return d.assignString(src, dst)
	case reflect.Bool:
		return d.assignBool(src, dst)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return d.assignInt(src, dst)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return d.assignUint(src, dst)
	case reflect.Float32, reflect.Float64:
		return d.assignFloat(src, dst)
	case reflect.Slice:
		return d.assignSlice(src, dst)
	case reflect.Array:
		return d.assignArray(src, dst)
	case reflect.Map:
		return d.assignMap(src, dst)
	case reflect.Struct:
		return d.assignStruct(src, dst)
	default:
		return types.NewToonError(fmt.Sprintf("unsupported destination type: %v", dstType), 0, 0)
	}
}

func (d *Decoder) assignString(src, dst reflect.Value) error {
	var str string
	switch src.Kind() {
	case reflect.String:
		str = src.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.FormatInt(src.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = strconv.FormatUint(src.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat(src.Float(), 'g', -1, 64)
	case reflect.Bool:
		str = strconv.FormatBool(src.Bool())
	default:
		return types.NewToonError(fmt.Sprintf("cannot convert %v to string", src.Type()), 0, 0)
	}

	if dst.CanSet() {
		dst.SetString(str)
	}
	return nil
}

func (d *Decoder) assignBool(src, dst reflect.Value) error {
	var b bool
	switch src.Kind() {
	case reflect.Bool:
		b = src.Bool()
	case reflect.String:
		var err error
		b, err = strconv.ParseBool(src.String())
		if err != nil {
			return types.NewToonError(fmt.Sprintf("cannot parse bool from string: %s", src.String()), 0, 0)
		}
	default:
		return types.NewToonError(fmt.Sprintf("cannot convert %v to bool", src.Type()), 0, 0)
	}

	if dst.CanSet() {
		dst.SetBool(b)
	}
	return nil
}

func (d *Decoder) assignInt(src, dst reflect.Value) error {
	var i int64
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i = src.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i = int64(src.Uint())
	case reflect.Float32, reflect.Float64:
		i = int64(src.Float())
	case reflect.String:
		var err error
		i, err = strconv.ParseInt(src.String(), 10, 64)
		if err != nil {
			return types.NewToonError(fmt.Sprintf("cannot parse int from string: %s", src.String()), 0, 0)
		}
	default:
		return types.NewToonError(fmt.Sprintf("cannot convert %v to int", src.Type()), 0, 0)
	}

	if dst.CanSet() {
		dst.SetInt(i)
	}
	return nil
}

func (d *Decoder) assignUint(src, dst reflect.Value) error {
	var u uint64
	switch src.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u = src.Uint()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		u = uint64(src.Int())
	case reflect.Float32, reflect.Float64:
		u = uint64(src.Float())
	case reflect.String:
		var err error
		u, err = strconv.ParseUint(src.String(), 10, 64)
		if err != nil {
			return types.NewToonError(fmt.Sprintf("cannot parse uint from string: %s", src.String()), 0, 0)
		}
	default:
		return types.NewToonError(fmt.Sprintf("cannot convert %v to uint", src.Type()), 0, 0)
	}

	if dst.CanSet() {
		dst.SetUint(u)
	}
	return nil
}

func (d *Decoder) assignFloat(src, dst reflect.Value) error {
	var f float64
	switch src.Kind() {
	case reflect.Float32, reflect.Float64:
		f = src.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f = float64(src.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f = float64(src.Uint())
	case reflect.String:
		var err error
		f, err = strconv.ParseFloat(src.String(), 64)
		if err != nil {
			return types.NewToonError(fmt.Sprintf("cannot parse float from string: %s", src.String()), 0, 0)
		}
	default:
		return types.NewToonError(fmt.Sprintf("cannot convert %v to float", src.Type()), 0, 0)
	}

	if dst.CanSet() {
		dst.SetFloat(f)
	}
	return nil
}

func (d *Decoder) assignSlice(src, dst reflect.Value) error {
	if src.Kind() != reflect.Slice {
		return types.NewToonError(fmt.Sprintf("cannot assign %v to slice", src.Type()), 0, 0)
	}

	srcLen := src.Len()
	dstType := dst.Type()

	slice := reflect.MakeSlice(dstType, srcLen, srcLen)

	for i := 0; i < srcLen; i++ {
		srcElem := src.Index(i)
		dstElem := slice.Index(i)

		if err := d.assignReflectValue(srcElem, dstElem); err != nil {
			return err
		}
	}

	if dst.CanSet() {
		dst.Set(slice)
	}
	return nil
}

func (d *Decoder) assignArray(src, dst reflect.Value) error {
	if src.Kind() != reflect.Slice {
		return types.NewToonError(fmt.Sprintf("cannot assign %v to array", src.Type()), 0, 0)
	}

	srcLen := src.Len()
	dstLen := dst.Len()

	if srcLen != dstLen {
		return types.NewToonError(fmt.Sprintf("array length mismatch: source %d, destination %d", srcLen, dstLen), 0, 0)
	}

	for i := 0; i < srcLen; i++ {
		srcElem := src.Index(i)
		dstElem := dst.Index(i)

		if err := d.assignReflectValue(srcElem, dstElem); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) assignMap(src, dst reflect.Value) error {
	if src.Kind() != reflect.Map {
		return types.NewToonError(fmt.Sprintf("cannot assign %v to map", src.Type()), 0, 0)
	}

	dstType := dst.Type()
	keyType := dstType.Key()
	elemType := dstType.Elem()

	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dstType))
	}

	for _, key := range src.MapKeys() {
		srcValue := src.MapIndex(key)

		// Convert key if necessary
		dstKey := reflect.New(keyType).Elem()
		if err := d.assignReflectValue(key, dstKey); err != nil {
			return err
		}

		// Convert value
		dstValue := reflect.New(elemType).Elem()
		if err := d.assignReflectValue(srcValue, dstValue); err != nil {
			return err
		}

		dst.SetMapIndex(dstKey, dstValue)
	}

	return nil
}

func (d *Decoder) assignStruct(src, dst reflect.Value) error {
	if src.Kind() != reflect.Map {
		return types.NewToonError(fmt.Sprintf("cannot assign %v to struct", src.Type()), 0, 0)
	}

	dstType := dst.Type()

	for _, key := range src.MapKeys() {
		keyStr := key.String()
		srcValue := src.MapIndex(key)

		// Find struct field
		field, found := utils.FindStructField(dstType, keyStr)
		if !found {
			if d.opts.Strict {
				return types.NewToonError(fmt.Sprintf("unknown field: %s", keyStr), 0, 0)
			}
			continue
		}

		dstField := dst.Field(field.Index[0])
		if !dstField.CanSet() {
			continue
		}

		if err := d.assignReflectValue(srcValue, dstField); err != nil {
			return err
		}
	}

	return nil
}
