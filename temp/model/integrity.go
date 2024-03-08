package model

import (
	"bytes"
	"errors"
)

type OrphanedRecord struct {
	ParentId *string `json:"parent_id"`
	ChildId  *string `json:"child_id"`
}

type RelationalIntegrityCheckData struct {
	ParentName   string           `json:"parent_name"`
	ChildName    string           `json:"child_name"`
	ParentIdAttr string           `json:"parent_id_attr"`
	ChildIdAttr  string           `json:"child_id_attr"`
	Records      []OrphanedRecord `json:"records"`
}

type IntegrityCheckResult struct {
	Data any   `json:"data"`
	Err  error `json:"err"`
}

func (r *IntegrityCheckResult) UnmarshalJSON(b []byte) error {
	var data map[string]any
	if err := ModelFromJson(&data, bytes.NewReader(b)); err != nil {
		return err
	}
	if d, ok := data["data"]; ok && d != nil {
		var rdata RelationalIntegrityCheckData
		m := d.(map[string]any)
		rdata.ParentName = m["parent_name"].(string)
		rdata.ChildName = m["child_name"].(string)
		rdata.ParentIdAttr = m["parent_id_attr"].(string)
		rdata.ChildIdAttr = m["child_id_attr"].(string)
		for _, recData := range m["records"].([]any) {
			var record OrphanedRecord
			m := recData.(map[string]any)
			if val := m["parent_id"]; val != nil {
				record.ParentId = GetPointerOfValue(val.(string))
			}
			if val := m["child_id"]; val != nil {
				record.ChildId = GetPointerOfValue(val.(string))
			}
			rdata.Records = append(rdata.Records, record)
		}
		r.Data = rdata
	}
	if err, ok := data["err"]; ok && err != nil {
		r.Err = errors.New(data["err"].(string))
	}
	return nil
}
