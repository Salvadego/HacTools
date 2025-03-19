package models

type PKAnalyzeResponse struct {
	ComposedTypeCode string `json:"pkComposedTypeCode"`
	ItemPK           string `json:"itemPK,omitempty"`
	ItemType         string `json:"itemType,omitempty"`
	ComposedPK       string `json:"composedPK,omitempty"`
	ComposedType     string `json:"composedType,omitempty"`
	ComposedTypeUID  string `json:"composedTypeUID,omitempty"`
}
