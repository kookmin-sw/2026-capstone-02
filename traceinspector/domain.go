package traceinspector

type AbstractDomain interface {
	IsBot() bool
	IsTop() bool
	Incl(lhs AbstractDomain, rhs AbstractDomain) bool          // inclusion operator `lhs ⊑ rhs`
	Join(a1 AbstractDomain, a2 AbstractDomain) AbstractDomain  // abstract join operator `a1 ⊔ a2`
	Widen(a1 AbstractDomain, a2 AbstractDomain) AbstractDomain // widening operator `a1 ▽ a2`
	ToString() string                                          // return string representation of the domain value
}
