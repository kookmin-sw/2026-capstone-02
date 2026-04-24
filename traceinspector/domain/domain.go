package domain

import (
	"traceinspector/imp"
)

type AbstractDomain[DomainImpl any] interface {
	IsBot() bool
	IsTop() bool
	Incl(DomainImpl) bool               // inclusion operator `lhs ⊑ rhs`
	Join(DomainImpl) (DomainImpl, bool) // abstract join operator `lhs ⊔ rhs`. returns result and whether result is changed
	Widen(DomainImpl) DomainImpl        // widening operator `lhs ▽ rhs`
	String() string                     // return string representation of the domain value
	Clone() DomainImpl                  // Return a copy of the domain
}

type IntegerDomain[DomainImpl any] interface {
	AbstractDomain[DomainImpl]
	From_IntLitExpr(imp.IntLitExpr) DomainImpl
	Add(DomainImpl) DomainImpl
	Sub(DomainImpl) DomainImpl
	Mul(DomainImpl) DomainImpl
	Div(DomainImpl) DomainImpl
	Mod(DomainImpl) DomainImpl
	Eq(DomainImpl) BoolDomain
	Neq(DomainImpl) BoolDomain
	Lessthan(DomainImpl) BoolDomain
	Greaterthan(DomainImpl) BoolDomain
	Leq(DomainImpl) BoolDomain
	Geq(DomainImpl) BoolDomain
	Neg() DomainImpl
	Filter(FilterQueryType, DomainImpl) DomainImpl // compute the result of filtering the current domain
	// x.Filter(<=, y) = x', where x' ⊑ x and x' <= y
}
