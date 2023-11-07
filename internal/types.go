package internal

import "time"


type Address struct {
	ID              int64
	StreetAddress   string
	NeighbourhoodID *int64
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

type Building struct {
	ID                    int64
	Code                  *string
	NameFi                *string
	NameEn                *string
	NameRu                *string
	Address               string
	ConstructionStartYear *int
	CompletionYear        *int
	ComplexFi             *string
	ComplexEn             *string
	ComplexRu             *string
	HistoryFi             *string
	HistoryEn             *string
	HistoryRu             *string
	ReasoningFi           *string
	ReasoningEn           *string
	ReasoningRu           *string
	ProtectionStatusFi    *string
	ProtectionStatusEn    *string
	ProtectionStatusRu    *string
	InfoSourceFi          *string
	InfoSourceEn          *string
	InfoSourceRu          *string
	SurroundingsFi        *string
	SurroundingsEn        *string
	SurroundingsRu        *string
	FoundationFi          *string
	FoundationEn          *string
	FoundationRu          *string
	FrameFi               *string
	FrameEn               *string
	FrameRu               *string
	FloorDescriptionFi    *string
	FloorDescriptionEn    *string
	FloorDescriptionRu    *string
	FacadesFi             *string
	FacadesEn             *string
	FacadesRu             *string
	SpeciaFeaturesFi      *string
	SpeciaFeaturesEn      *string
	SpeciaFeaturesRu      *string
	latitude_ETRSGK25     *float32
	longitude_ERRSGK25    *float32
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	DeletedAt             *time.Time
}