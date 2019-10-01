package types

// Ticket represents a validator ticket
type Ticket struct {
	DecayBy        uint64 `gorm:"column:decayBy" json:"decayBy"`               // Block height when the ticket becomes decayed
	MatureBy       uint64 `gorm:"column:matureBy" json:"matureBy"`             // Block height when the ticket enters maturity.
	Hash           string `gorm:"column:hash" json:"hash"`                     // Hash of the ticket purchase transaction
	ChildOf        string `gorm:"column:childOf" json:"childOf,omitempty"`     // The hash of another ticket which this ticket is derived from
	ProposerPubKey string `gorm:"column:proposerPubKey" json:"proposerPubKey"` // The public key of the validator that owns the ticket.
	Height         uint64 `gorm:"column:height" json:"height"`                 // The block height where this ticket was seen.
	Index          int    `gorm:"column:index" json:"index"`                   // The index of the ticket in the transactions list.
	Value          string `gorm:"column:value" json:"value"`                   // The value paid for the ticket (as a child - then for the parent ticket)
}

// QueryOptions describe how a query should be executed.
type QueryOptions struct {
	Limit   int    `json:"limit" mapstructure:"limit"`
	Offset  int    `json:"offset" mapstructure:"offset"`
	Order   string `json:"order" mapstructure:"order"`
	NoChild bool   `json:"noChild" mapstructure:"noChild"`
}

// TicketManager describes a ticket manager
// Get finds tickets belonging to the given proposer.
type TicketManager interface {
	// Index adds a ticket (and child tickets) to the ticket index.
	Index(tx *Transaction, proposerPubKey string, blockHeight uint64, txIndex int) error
	// Get finds tickets belonging to the given proposer.
	Get(proposerPubKey string, queryOpt QueryOptions) ([]*Ticket, error)
	// Stop stops the ticket manager
	Stop() error
}