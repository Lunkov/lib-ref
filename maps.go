package ref

import (
  "strings"
  "strconv"
  "github.com/google/uuid"
  "github.com/golang/glog"
)

const MAP_SEPARATOR   = "."
const ARRAY_SEPARATOR = ":"

func AppendChildMap(parentMap *map[string]interface{}, child string, childMap map[string]interface{}) {
  for k, v := range childMap {
    (*parentMap)[strings.ToLower(child + MAP_SEPARATOR + k)] = v
  }
}

func UnionMaps(dstMap *map[string]interface{}, newMap *map[string]interface{}) {
  if dstMap != nil && newMap != nil {
    for k, v := range (*newMap) {
      (*dstMap)[strings.ToLower(k)] = v
    }
  }
}

func UnionMapsStr(dstMap *map[string]interface{}, newMap *map[string]string) {
  if dstMap != nil && newMap != nil {
    for k, v := range (*newMap) {
      (*dstMap)[k] = v
    }
  }
}


func GetChildSubmap(parentMap *map[string]interface{}, child string, separator string) map[string]interface{} {
  res := make(map[string]interface{})
  zsChild := len(child) + 1
  seach := child + separator
  for k, v := range (*parentMap) {
    if zsChild < len(k) {
      if seach == k[:zsChild] {
        k2 := k[zsChild:]
        res[strings.ToLower(k2)] = v
      }
    }
  }
  return res
}

func GetSizeSubmap(parentMap *map[string]interface{}, separator string) int {
  res := 0
  for k, _ := range (*parentMap) {
    i := strings.Index(k, separator)
    if i >= 0 {
      seach := k[0:i]
      t, err := strconv.Atoi(seach)
      if err == nil {
        if res < t + 1 {
          res = t + 1
        }
      }
    }
  }
  return res
}

func GetMapFieldUUID(data *map[string]interface{}, fieldname string) (uuid.UUID, bool) {
  ids, ok := (*data)[strings.ToLower(fieldname)]
  if !ok {
    return uuid.Nil, false
  }
  id, oks := ids.(string)
  if !oks {
    return uuid.Nil, false
  }
  uid, err := uuid.Parse(id)
  if err != nil {
    glog.Errorf("ERR: getFieldUUID(%s) %v", id, err)
    return uuid.Nil, false
  }
  return uid, true
}

