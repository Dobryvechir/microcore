/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

const MapKeySeparator = "~^~!~"

func (item *DvFieldInfo) CreateQuickInfoByKeys(ids []string) {
	if item != nil {
		if item.QuickSearch == nil {
			item.QuickSearch = &QuickSearchInfo{Looker: make(map[string]*DvFieldInfo)}
		}
		if item.QuickSearch.Looker == nil {
			item.QuickSearch.Looker = make(map[string]*DvFieldInfo)
		}
		n := len(item.Fields)
		m := len(ids)
		for i := 0; i < n; i++ {
			f := item.Fields[i]
			if f != nil {
				f.CreateQuickInfoForObjectType()
				key := ""
				for j := 0; j < m; j++ {
					id := ids[j]
					fc := f.QuickSearch.Looker[id]
					v := ""
					if fc != nil {
						v = string(fc.Value)
					}
					if j == 0 {
						key = v
					} else {
						key = key + MapKeySeparator + v
					}
				}
				f.QuickSearch.Key = key
				item.QuickSearch.Looker[key] = f
			}
		}
	}
}

func (item *DvFieldInfo) CreateQuickInfoForObjectType() {
	if item != nil {
		if item.QuickSearch == nil {
			item.QuickSearch = &QuickSearchInfo{Looker: make(map[string]*DvFieldInfo)}
		}
		if item.QuickSearch.Looker == nil {
			item.QuickSearch.Looker = make(map[string]*DvFieldInfo)
		}
		n := len(item.Fields)
		for i := 0; i < n; i++ {
			f := item.Fields[i]
			if f != nil {
				name := string(f.Name)
				if f.QuickSearch == nil {
					f.QuickSearch = &QuickSearchInfo{Key: name}
				} else {
					f.QuickSearch.Key = name
				}
				item.QuickSearch.Looker[name] = f
			}
		}
	}
}

func CreateQuickInfoByKeysForAny(data interface{}, ids []string) {
	n := len(ids)
	if n > 0 {
		switch data.(type) {
		case *DvFieldInfo:
			data.(*DvFieldInfo).CreateQuickInfoByKeys(ids)
		}
	} else {
		switch data.(type) {
		case *DvFieldInfo:
			data.(*DvFieldInfo).CreateQuickInfoForObjectType()
		}
	}
}

