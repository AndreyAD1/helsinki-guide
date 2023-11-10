package internal

import "time"

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type Neighbourhood struct {
	ID           int64
	Name         string
	Municipality *string
	Timestamps
}

type Address struct {
	ID              int64
	StreetAddress   string
	NeighbourhoodID *int64
	Timestamps
}

type Actor struct {
	ID      int64
	Name    string
	TitleFi *string
	TitleEn *string
	TitleRu *string
	Timestamps
}

type UseType struct {
	ID     int64
	NameFi string
	NameEn string
	NameRu string
	Timestamps
}

type Building struct {
	ID                    int64
	Code                  *string
	NameFi                *string
	NameEn                *string
	NameRu                *string
	Address               Address
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
	SpecialFeaturesFi     *string
	SpecialFeaturesEn     *string
	SpecialFeaturesRu     *string
	Latitude_ETRSGK25     *float32
	Longitude_ERRSGK25    *float32
	AuthorIds             []int64
	InitialUses           []UseType
	CurrentUses           []UseType
	Timestamps
}
