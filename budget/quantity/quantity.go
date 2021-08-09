package quantity

// Quantity describes a human-friendly quantity
type Quantity interface {
	ValueOf() float64
}
