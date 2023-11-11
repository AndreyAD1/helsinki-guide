package repositories


var insertBuilding = `INSERT INTO building
(	
	code, 
	name_fi, 
	name_en, 
	name_ru, 
	address_id, 
	construction_start_year,
	completion_year, 
	complex_fi, 
	complex_en, 
	complex_ru, 
	history_fi,
	history_en, 
	history_ru, 
	reasoning_fi, 
	reasoning_en, 
	reasoning_ru,
	protection_status_fi, 
	protection_status_en, 
	protection_status_ru,
	info_source_fi,
	info_source_en,
	info_source_ru,
	surroundings_fi,
	surroundings_en,
	surroundings_ru,
	foundation_fi,
	foundation_en,
	foundation_ru,
	frame_fi,
	frame_en,
	frame_ru,
	floor_description_fi,
	floor_description_en,
	floor_description_ru,
	facades_fi,
	facades_en,
	facades_ru,
	special_features_fi,
	special_features_en,
	special_features_ru,
	latitude_ETRSGK25,
	longitude_ETRSGK25,
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11,
	$12,
	$13,
	$14,
	$15,
	$16,
	$17,
	$18,
	$19,
	$20,
	$21,
	$22,
	$23,
	$24,
	$25,
	$26,
	$27,
	$28,
	$29,
	$30,
	$31,
	$32,
	$33,
	$34,
	$35,
	$36,
	$37,
	$38,
	$39,
	$40,
	$41
) RETURNING id;`

var getAddress = `SELECT * FROM addresses WHERE street_address = $1`
var insertAddress = `INSERT INTO addresses (street_address, neighbourhood_id) VALUES ($1, $2, $3, $4) RETURNING id;`

var insertBuildingAuthor = `INSERT INTO building_authors (building_id, actor_id)
VALUES ($1, $2);`

var getUseType = `SELECT * FROM use_types WHERE name_en = $1;`
var insertUseType = `INSERT INTO use_types (name_fi, name_en, name_ru)
VALUES ($1, $2, $3, $4);`

var insertInitialUses = `INSERT INTO initial_uses (building_id, use_type_id)
VALUES ($1, $2);`
var insertCurrentUses = `INSERT INTO current_uses (building_id, use_type_id)
VALUES ($1, $2);`