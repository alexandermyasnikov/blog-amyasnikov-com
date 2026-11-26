func (p *Repository) BuyTicket(userId uint) (*model.Ticket, error) {
	...

	ticket := model.NewTicket(userId)
	if err := tx.Create(&ticket).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	...

	return ticket, nil
}

func NewTicket(userId uint) *Ticket {
	return &Ticket{
		Number: helpers.GenerateTicketNumber(),
		Status: StatusActive,
		UserID: userId,
	}
}

func GenerateTicketNumber() string {
	source := rand.NewSource(time.Now().UnixNano() / int64(time.Millisecond))
	random := rand.New(source)

	return fmt.Sprintf("%06d", random.Intn(1000000))
}
