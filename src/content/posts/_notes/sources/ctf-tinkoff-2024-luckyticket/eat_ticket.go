func (p *Repository) EatTicket(userId, ticketId uint) (int, string, error) {
	...

	ticket := new(model.Ticket)
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ticketId).First(&ticket).Error; err != nil {
		tx.Rollback()
		return 0, "", errIdNotFound
	}

	...

	newLuck := helpers.GetLuckFromTicket(ticket.Number)

	user := new(model.User)
	if err := tx.Model(&model.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		tx.Rollback()
		return 0, "", err
	}

	updatedLuck := helpers.UpdateLuck(user.Luck, newLuck)

	...

	return newLuck, ticket.Number, nil
}

func GetLuckFromTicket(number string) int {
	if isTicketLucky(number) {
		return 1
	} else {
		return -10
	}
}

func UpdateLuck(currentLuck uint, inLuck int) uint {
	luck := int(currentLuck) + inLuck

	if luck > 100 {
		return 100
	} else if luck < 0 {
		return 0
	} else {
		return uint(luck)
	}
}
