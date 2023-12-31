package repositories

type UserSpecificationByTelegramID struct {
	telegramID int64
}

func NewUserSpecificationByID(telegramID int64) *UserSpecificationByTelegramID {
	return &UserSpecificationByTelegramID{telegramID}
}

func (a *UserSpecificationByTelegramID) ToSQL() (string, map[string]any) {
	query := `SELECT id, telegram_id, language, created_at,
	updated_at, deleted_at FROM users WHERE telegram_id = @telegram_id;`
	return query, map[string]any{"telegram_id": a.telegramID}
}

func UserByIDIsEqual(telegramID int64) func(s *UserSpecificationByTelegramID) bool {
	return func(s *UserSpecificationByTelegramID) bool {
		return telegramID == s.telegramID
	}
}
