package ref

import (
  "strings"
  "reflect"
  "time"
  "math"
  "strconv"
  "github.com/google/uuid"
  "github.com/golang/glog"
)

var timeKind = reflect.TypeOf(time.Time{}).Kind()

func ConvertToMap(a interface{}) map[string]interface{} {
  res := make(map[string]interface{})
  v := reflect.ValueOf(a)
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }
  if v.Kind() == reflect.Struct {
    for i := 0; i < v.NumField(); i++ {
      field := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0]
      if field != "" && field != "-" {
        if v.Field(i).IsValid() {
          if !valueIsZero(v.Field(i)) {
            switch v.Field(i).Kind() {
            case reflect.Struct:
                    AppendChildMap(&res, field, ConvertToMap(v.Field(i).Interface()))
                    break;
            case reflect.Slice:
                    s := v.Field(i)
                    for j := 0; j < s.Len(); j++ {
                      ei := field + ARRAY_SEPARATOR + strconv.Itoa(j)
                      if s.Index(j).Kind() == reflect.Struct {
                        AppendChildMap(&res, ei, ConvertToMap(s.Index(j).Interface()))
                      } else {
                        res[strings.ToLower(ei)] = s.Index(j).Interface()
                      }
                    }
                    break;
            case reflect.Map:
                    for _, e := range v.Field(i).MapKeys() {
                      ei := field + MAP_SEPARATOR + e.String()
                      mi := v.Field(i).MapIndex(e)
                      switch t := mi.Interface().(type) {
                      case int:
                          res[strings.ToLower(ei)] = t
                      case string:
                          res[strings.ToLower(ei)] = t
                      case bool:
                          res[strings.ToLower(ei)] = t
                      default:
                          AppendChildMap(&res, ei, ConvertToMap(t))
                      }
                    }
                    break;
            case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                    res[strings.ToLower(field)] = strconv.FormatInt(v.Field(i).Int(), 10)
                    break
            //case reflect.Float32, reflect.Float64:
            //case reflect.String:
            //        break;
            default:
                    if glog.V(9) {
                      glog.Infof("DBG: ConvertToMap Model(%s:%s) default", v.Field(i).Kind(), field)
                    }
                    res[strings.ToLower(field)] = v.Field(i).Interface()
                    break;
            }
          }
        }
      }
    }
  }
  return res
}

func ConvertFromMap(a interface{}, data *map[string]interface{}) {
  if a == nil {
    glog.Errorf("ERR: Model() is NULL")
    return
  }
  v := reflect.ValueOf(a)
  if v.IsNil() {
    glog.Errorf("ERR: Model() is NULL")
    return
  }
  if v.Kind() != reflect.Ptr {
    glog.Errorf("ERR: Model(%s) not Pointer\n", v.Type())
    return
  }
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }
  uid0, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
  typeUUID := reflect.TypeOf(uid0)
  if v.Kind() == reflect.Struct {
    for i := 0; i < v.NumField(); i++ {
      field := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0]

      if field != "" && field != "-" {
        if v.Field(i).IsValid() && v.Field(i).CanSet() {
          switch v.Field(i).Kind() {
            case reflect.Struct:
                      childMap := GetChildSubmap(data, field, MAP_SEPARATOR)
                      ConvertFromMap(v.Field(i).Addr().Interface(), &childMap)
                      break
            case reflect.Map:
                      childMap := GetChildSubmap(data, field, MAP_SEPARATOR)
                      sz := GetSizeSubmap(&childMap, MAP_SEPARATOR)
                      if glog.V(9) {
                        glog.Infof("DBG: ConvertFromMap Model(%s:%s) sz=%d", v.Field(i).Kind(), field, sz)
                      }

                      v.Field(i).Set( reflect.MakeMap( reflect.TypeOf(v.Field(i).Interface()) ) )
                      for j := 0; j < sz; j++ {
                        childItem := GetChildSubmap(&childMap, strconv.Itoa(j), MAP_SEPARATOR)
                        ConvertFromMap(v.Field(i).Index(j).Addr().Interface(), &childItem)
                      }
                      break
            case reflect.Slice:
                      childMap := GetChildSubmap(data, field, ARRAY_SEPARATOR)
                      sz := GetSizeSubmap(&childMap, MAP_SEPARATOR)
                      if glog.V(9) {
                        glog.Infof("DBG: ConvertFromMap Model(%s:%s) sz=%d", v.Field(i).Kind(), field, sz)
                      }
 
                      // Create a slice to begin with
                      typeItem := reflect.TypeOf(v.Field(i).Interface()).Elem()
                      v.Field(i).Set( reflect.MakeSlice(reflect.SliceOf( typeItem ), sz, sz) )
                      for j := 0; j < sz; j++ {
                        childItem := GetChildSubmap(&childMap, strconv.Itoa(j), MAP_SEPARATOR)
                        ConvertFromMap(v.Field(i).Index(j).Addr().Interface(), &childItem)
                      }
                      break
            case reflect.Array:
                      d, ok := (*data)[field]
                      if ok {
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) array", v.Type(), field)
                        }
                        if v.Field(i).Type() == typeUUID && reflect.TypeOf(d) == reflect.TypeOf(string("")) {
                          str, _ := d.(string)
                          uid1, _ := uuid.Parse(str)
                          if uid1 != uuid.Nil {
                            v.Field(i).Set(reflect.ValueOf(uid1))
                          }
                        } else {
                          v.Field(i).Set(reflect.ValueOf(d))
                        }
                      }
                      break
            case reflect.Float32, reflect.Float64:
                      d, ok := (*data)[field]
                      if ok {
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) float", v.Type(), field)
                        }
                        v.Field(i).SetFloat(d.(float64)) //reflect.ValueOf(d))
                      }
                      break
            case reflect.Int32, reflect.Int64:
                      d, ok := (*data)[field]
                      if ok {
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) float", v.Type(), field)
                        }
                        v.Field(i).SetInt(d.(int64)) //reflect.ValueOf(d))
                      }
                      break
            case timeKind:
                      d, ok := (*data)[field]
                      if ok {
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) Time (%v)", v.Type(), field, v.Field(i).Kind())
                        }
                        v.Field(i).SetString(d.(time.Time).String())
                      }
                      break
            case reflect.String:
                      d, ok := (*data)[field]
                      if ok {
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) string (%v)", v.Type(), field, v.Field(i).Kind())
                        }
                        switch v := d.(type) {
                        case nil:
                            if glog.V(9) {
                              glog.Infof("DBG: x is nil")            // here v has type interface{}
                            }
                            d = ""
                        case int: 
                            if glog.V(9) {
                              glog.Infof("DBG: x is %v", v)             // here v has type int
                            }
                            d = strconv.Itoa(d.(int))
                        case int64: 
                            if glog.V(9) {
                              glog.Infof("DBG: x is %v", v)             // here v has type int64
                            }
                            d = strconv.FormatInt(d.(int64), 10)
                        case float32: 
                            if glog.V(9) {
                              glog.Infof("DBG: x is %v", v)             // here v has type float32
                            }
                            d = strconv.FormatFloat(float64(d.(float32)), 'E', -1, 64)
                        case float64: 
                            if glog.V(9) {
                              glog.Infof("DBG: x is %v", v)             // here v has type float64
                            }
                            d = strconv.FormatFloat(d.(float64), 'E', -1, 64)
                        case bool:
                            if glog.V(9) {
                              glog.Infof("DBG: x is bool") // here v has type interface{}
                            }
                            d = strconv.FormatBool(d.(bool))
                        case string:
                            if glog.V(9) {
                              glog.Infof("DBG: x is string") // here v has type interface{}
                            }
                        case time.Time:
                            if glog.V(9) {
                              glog.Infof("DBG: x is time.Time") // here v has type interface{}
                            }
                            d = d.(time.Time).String()
                        default:
                            if glog.V(9) {
                              glog.Infof("DBG: type unknown")        // here v has type interface{}
                            }
                        }
                        v.Field(i).SetString(d.(string))
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) string (%v) setted", v.Type(), field, v.Field(i).Kind())
                        }
                      }
                      break
            default:
                      d, ok := (*data)[field]
                      if ok {
                        if glog.V(9) {
                          glog.Infof("DBG: ConvertFromMap Model(%s:%s) default (%v)", v.Type(), field, v.Field(i).Kind())
                        }
                        v.Field(i).Set(reflect.ValueOf(d))
                      }
                      break
          }
        }
      }
    }
  }
}

// Check Value
func valueIsZero(v reflect.Value) bool {
  switch v.Kind() {
    case reflect.Bool:
      return !v.Bool()
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
      return v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
      return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
      return math.Float64bits(v.Float()) == 0
    case reflect.Complex64, reflect.Complex128:
      c := v.Complex()
      return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
    case reflect.Array:
      for i := 0; i < v.Len(); i++ {
        if !valueIsZero(v.Index(i)) {
          return false
        }
      }
      return true
    case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
      return v.IsNil()
    case reflect.String:
      return v.Len() == 0
    case reflect.Struct:
      for i := 0; i < v.NumField(); i++ {
        if !valueIsZero(v.Field(i)) {
          return false
        }
      }
      return true
    default:
      // This should never happens, but will act as a safeguard for
      // later, as a default value doesn't makes sense here.
      panic(&reflect.ValueError{"reflect.Value.IsZero", v.Kind()})
  }
}


func GetFieldUUID(info interface{}, fieldname string) (uuid.UUID, bool) {
  v := reflect.ValueOf(info)
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }
  if !v.FieldByName(fieldname).IsValid() ||
     !v.FieldByName(fieldname).CanInterface() {
    return uuid.Nil, false
  }
  val := v.FieldByName(fieldname).Interface()
  uid, ok := val.(uuid.UUID)
  if  ok {
    return uid, true
  }
  return uuid.Nil, false
}

func SetFieldUUID(info interface{}, id uuid.UUID, fieldname string) bool {
  v := reflect.ValueOf(info)
  if v.Kind() == reflect.Ptr {
     v = v.Elem()
  }
  if !v.FieldByName(fieldname).IsValid() ||
     !v.FieldByName(fieldname).CanInterface() {
    return false
  }
  v.FieldByName(fieldname).Set(reflect.ValueOf(id))
  return true
}

func GetFieldString(info interface{}, fieldname string) (string, bool) {
  v := reflect.ValueOf(info)
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }
  if !v.FieldByName(fieldname).IsValid() ||
     !v.FieldByName(fieldname).CanInterface() {
    return "", false
  }
  val := v.FieldByName(fieldname).Interface()
  code, ok := val.(string)
  if  ok {
    return code, true
  }
  return "", false
}

func FieldExists(info interface{}, fieldname string) bool {
  v := reflect.ValueOf(info)
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }
  if !v.FieldByName(fieldname).IsValid() ||
     !v.FieldByName(fieldname).CanInterface() {
    return false
  }
  return true
}

func InitializeStruct(t reflect.Type, v reflect.Value) {
  for i := 0; i < v.NumField(); i++ {
    f := v.Field(i)
    ft := t.Field(i)
    switch ft.Type.Kind() {
    case reflect.Map:
      f.Set(reflect.MakeMap(ft.Type))
    case reflect.Slice:
      f.Set(reflect.MakeSlice(ft.Type, 0, 0))
    case reflect.Chan:
      f.Set(reflect.MakeChan(ft.Type, 0))
    case reflect.Struct:
      InitializeStruct(ft.Type, f)
    case reflect.Ptr:
      fv := reflect.New(ft.Type.Elem())
      InitializeStruct(ft.Type.Elem(), fv.Elem())
      f.Set(fv)
    default:
    }
  }
}

func ValueToString(info interface{}) (string, bool) {
  res := ""
	v := reflect.ValueOf(info)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
  if !v.IsValid() ||
     !v.CanInterface() {
    return res, false
  }
  
	switch v.Kind() {
	case reflect.Bool:
    if v.Bool() {
      return "true", true
    }
		return "false", true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := v.Int()
    return strconv.FormatInt(i, 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i := v.Uint()
    return strconv.FormatUint(i, 10), true
	case reflect.Float32, reflect.Float64:
    f := v.Float()
    return strconv.FormatFloat(f, 'g', -1, 64), true
		//math.Float64bits(v.Float())
	case reflect.Complex64, reflect.Complex128:
		//c := v.Complex()
		//return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.String:
    return v.String(), true
	case timeKind:
    return v.String(), true    
	case reflect.Array:
		return res, false
  case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return res, false
	case reflect.Struct:
		return res, false
	}
  return res, false
}

func GetType(myvar interface{}) string {
  t := reflect.TypeOf(myvar)
  if t == nil {
    return "<nil>"
  }
  if t.Kind() == reflect.Ptr {
    return "*" + t.Elem().Name()
  } else {
    return t.Name()
  }
}
