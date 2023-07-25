package model

type Seo struct {
	SeoTitle       string `json:"seo_title,omitempty" gorm:"type:varchar(70);column:SeoTitle"`
	SeoDescription string `json:"seo_description,omitempty" gorm:"type:varchar(300);column:SeoDescription"`
}

func (s *Seo) commonPre() {
	s.SeoTitle = SanitizeUnicode(s.SeoTitle)
	s.SeoDescription = SanitizeUnicode(s.SeoDescription)
}

// SeoTranslation represents translation for Seo
type SeoTranslation struct {
	SeoTitle       *string `json:"seo_title" gorm:"type:varchar(70);column:SeoTitle"`
	SeoDescription *string `json:"seo_description" gorm:"type:varchar(300);column:SeoDescription"`
}

func (s *SeoTranslation) commonPre() {
	if s.SeoTitle != nil {
		*s.SeoTitle = SanitizeUnicode(*s.SeoTitle)
	}
	if s.SeoDescription != nil {
		*s.SeoDescription = SanitizeUnicode(*s.SeoDescription)
	}
}
