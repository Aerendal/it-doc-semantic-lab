package domain

// Industry identifies the vertical sector for which a document is scoped.
type Industry string

const (
	IndustryFinance       Industry = "finance"
	IndustryHealthcare    Industry = "healthcare"
	IndustryManufacturing Industry = "manufacturing"
	IndustryRetail        Industry = "retail"
	IndustryPublicSector  Industry = "public_sector"
	IndustryGeneric       Industry = "generic"
)

// IndustryMatrix maps an industry to its required document classes and coverage expectations.
type IndustryMatrix struct {
	Industry        Industry        `json:"industry"`
	RequiredClasses []DocumentClass `json:"required_classes"`
	OptionalClasses []DocumentClass `json:"optional_classes"`
}
