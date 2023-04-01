package ix

import (
	"encoding/json"
	"os"
	"path"
)

type SqlItem struct {
	Category string `json:"category"`
	Sql      string `json:"sql"`
}

func loadSqlItems() []SqlItem {
	file := path.Join("sqls.json")
	bs, err := os.ReadFile(file)
	if err != nil {
		panic("read sql items file err -> " + err.Error())
	}

	var sqlItems []SqlItem
	if err = json.Unmarshal(bs, &sqlItems); err != nil {
		panic("unmarshal sql items err -> " + err.Error())
	}

	return sqlItems
}

type SqlData struct {
	items []SqlItem
	sqls  map[string][]string
}

func (d *SqlData) Items() []SqlItem {
	return d.items
}

func (d *SqlData) Sqls(category string) []string {
	sqls, ok := d.sqls[category]
	if !ok {
		return nil
	}
	return sqls
}

func LoadSqlData() *SqlData {
	sqlItems := loadSqlItems()

	sqlData := &SqlData{
		items: sqlItems,
		sqls:  make(map[string][]string),
	}

	for _, item := range sqlItems {
		category, sql := item.Category, item.Sql
		sqlList, ok := sqlData.sqls[category]
		if !ok {
			sqlData.sqls[category] = []string{sql}
		} else {
			sqlData.sqls[category] = append(sqlList, sql)
		}
	}

	return sqlData
}

type CheckResult struct {
	Passed bool                   `json:"passed"`
	Err    string                 `json:"err"`
	Detail map[string]interface{} `json:"detail"`
}

func (r *CheckResult) SetPassed() {
	r.Passed = true
	r.Err = ""
}

func (r *CheckResult) SetWarning(warn string) {
	r.Passed = true
	r.Err = warn
}

func (r *CheckResult) SetError(err error) {
	r.Passed = false
	r.Err = err.Error()
}

func (r *CheckResult) SetDetail(key string, value any) {
	if r.Detail == nil {
		r.Detail = make(map[string]any)
	}
	r.Detail[key] = value
}

func (r *CheckResult) GetDetail(key string) any {
	if r.Detail != nil {
		v, ok := r.Detail[key]
		if ok {
			return v
		}
	}
	return nil
}

func NewCheckResult() *CheckResult {
	return &CheckResult{
		Passed: true,
		Err:    "no result output",
		Detail: map[string]interface{}{},
	}
}

type Checker func(sql string) *CheckResult
