package traceinspector

import (
	"fmt"
	"traceinspector/domain"
)

type ArraySummaryDomain[IntDomainImpl domain.IntegerDomain[IntDomainImpl]] struct {
	val               IntDomainImpl
	is_bottom, is_top bool
}

func (domain ArraySummaryDomain[IntDomainImpl]) String() string {
	if domain.is_bottom {
		return "⊥_bool"
	} else if domain.is_top {
		return "⊤_bool"
	} else {
		return fmt.Sprintf("%s", domain.val)
	}
}

func (domain ArraySummaryDomain[ElemDomain]) Clone() ArraySummaryDomain[ElemDomain] {
	return ArraySummaryDomain[ElemDomain]{val: domain.val, is_bottom: domain.is_bottom, is_top: domain.is_top}
}

func (domain ArraySummaryDomain[ElemDomain]) IsBot() bool {
	return domain.is_bottom
}

func (domain ArraySummaryDomain[ElemDomain]) IsTop() bool {
	return domain.is_top
}

func (lhs ArraySummaryDomain[ElemDomain]) Join(rhs ArraySummaryDomain[ElemDomain]) (ArraySummaryDomain[ElemDomain], bool) {
	if lhs.is_bottom {
		return rhs, rhs.is_bottom
	} else if rhs.is_bottom {
		return lhs, lhs.is_bottom
	} else if lhs.is_top || rhs.is_top {
		return ArraySummaryDomain[ElemDomain]{is_top: true}, lhs.is_top
	} else {
		elem_joined, elem_changed := lhs.val.Join(rhs.val)
		return ArraySummaryDomain[ElemDomain]{val: elem_joined}, elem_changed
	}
}

func (lhs ArraySummaryDomain[ElemDomain]) Incl(rhs ArraySummaryDomain[ElemDomain]) bool {
	return lhs.val.Incl(rhs.val)
}

func (lhs ArraySummaryDomain[ElemDomain]) Widen(rhs ArraySummaryDomain[ElemDomain]) ArraySummaryDomain[ElemDomain] {
	if lhs.is_bottom {
		return rhs
	}
	if rhs.is_bottom {
		return lhs
	}
	if lhs.is_top || rhs.is_top {
		return ArraySummaryDomain[ElemDomain]{is_top: true}
	}
	if lhs.val.Incl(rhs.val) {
		return lhs
	} else {
		return ArraySummaryDomain[ElemDomain]{val: lhs.val.Widen(rhs.val)}
	}
}

// expression evaluation

func (arr ArraySummaryDomain[IntDomainImpl]) Index(val IntDomainImpl) AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]] {
	return AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]]{domain_kind: IntDomainKind, int_domain: arr.val}
}
