package services

type Language string

var (
	Finnish = Language("fi")
	English = Language("en")
	Russian = Language("ru")
)

var codePerLanguage = map[string]Language{
	"fi": Finnish,
	"en": English,
	"ru": Russian,
}

func GetLanguagePerCode(code string) (Language, bool) {
	language, ok := codePerLanguage[code]
	return language, ok
}

type BuildingDTO struct {
	ID                int64
	NameFi            *string   `valueLanguage:"fi" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	NameEn            *string   `valueLanguage:"en" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	NameRu            *string   `valueLanguage:"ru" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	Address           string    `valueLanguage:"all" nameFi:"Katuosoite" nameEn:"Address" nameRu:"Адрес"`
	DescriptionFi     *string   `valueLanguage:"fi" nameFi:"Kerrosluku" nameEn:"Description" nameRu:"Описание"`
	DescriptionEn     *string   `valueLanguage:"en" nameFi:"Kerrosluku" nameEn:"Description" nameRu:"Описание"`
	DescriptionRu     *string   `valueLanguage:"ru" nameFi:"Kerrosluku" nameEn:"Description" nameRu:"Описание"`
	CompletionYear    *int      `valueLanguage:"all" nameFi:"Käyttöönottovuosi" nameEn:"Completion_year" nameRu:"Год_постройки"`
	Authors           *[]string `valueLanguage:"all" nameFi:"Suunnittelijat" nameEn:"Authors" nameRu:"Авторы"`
	FacadesFi         *string   `valueLanguage:"fi" nameFi:"Julkisivut" nameEn:"Facades" nameRu:"Фасады"`
	FacadesEn         *string   `valueLanguage:"en" nameFi:"Julkisivut" nameEn:"Facades" nameRu:"Фасады"`
	FacadesRu         *string   `valueLanguage:"ru" nameFi:"Julkisivut" nameEn:"Facades" nameRu:"Фасады"`
	DetailsFi         *string   `valueLanguage:"fi" nameFi:"Erityispiirteet" nameEn:"Interesting_details" nameRu:"Интересные_детали"`
	DetailsEn         *string   `valueLanguage:"en" nameFi:"Erityispiirteet" nameEn:"Interesting_details" nameRu:"Интересные_детали"`
	DetailsRu         *string   `valueLanguage:"ru" nameFi:"Erityispiirteet" nameEn:"Interesting_details" nameRu:"Интересные_детали"`
	NotableFeaturesFi *string   `valueLanguage:"fi" nameFi:"Huomattavia_ominaisuuksia" nameEn:"Notable_features" nameRu:"Примечательные_особенности"`
	NotableFeaturesEn *string   `valueLanguage:"en" nameFi:"Huomattavia_ominaisuuksia" nameEn:"Notable_features" nameRu:"Примечательные_особенности"`
	NotableFeaturesRu *string   `valueLanguage:"ru" nameFi:"Huomattavia_ominaisuuksia" nameEn:"Notable_features" nameRu:"Примечательные_особенности"`
	SurroundingsFi    *string   `valueLanguage:"fi" nameFi:"Ympäristönkuvaus" nameEn:"Surroundings" nameRu:"Окрестности"`
	SurroundingsEn    *string   `valueLanguage:"en" nameFi:"Ympäristönkuvaus" nameEn:"Surroundings" nameRu:"Окрестности"`
	SurroundingsRu    *string   `valueLanguage:"ru" nameFi:"Ympäristönkuvaus" nameEn:"Surroundings" nameRu:"Окрестности"`
	HistoryFi         *string   `valueLanguage:"fi" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	HistoryEn         *string   `valueLanguage:"en" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	HistoryRu         *string   `valueLanguage:"ru" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
}
