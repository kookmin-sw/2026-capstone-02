package traceinspector

type IntervalDomainValue struct {
	value              int64
	is_inf, is_neg_inf bool // disregard value field if any of these values is true
}

// returns whether v is a finite value
func (v *IntervalDomainValue) is_finite() bool {
	return !(v.is_inf || v.is_neg_inf)
}

// Compute the minimum of two IntervalDomainValues
func min_IntervalDomainValue(l1 IntervalDomainValue, l2 IntervalDomainValue) IntervalDomainValue {
	if l1.is_neg_inf || l2.is_neg_inf {
		return IntervalDomainValue{is_neg_inf: true} // zero values are 0 and false
	} else if l1.is_inf && l2.is_inf {
		return IntervalDomainValue{is_inf: true}
	} else if l1.is_inf {
		return IntervalDomainValue{value: l2.value}
	} else if l2.is_inf {
		return IntervalDomainValue{value: l1.value}
	} else {
		return IntervalDomainValue{value: min(l1.value, l2.value)}
	}
}

// Compute the maximum of two IntervalDomainValues
func max_IntervalDomainValue(l1 IntervalDomainValue, l2 IntervalDomainValue) IntervalDomainValue {
	if l1.is_inf || l2.is_inf {
		return IntervalDomainValue{is_inf: true}
	} else if l1.is_neg_inf && l2.is_neg_inf {
		return IntervalDomainValue{is_neg_inf: true}
	} else if l1.is_neg_inf {
		return IntervalDomainValue{value: l2.value}
	} else if l2.is_neg_inf {
		return IntervalDomainValue{value: l1.value}
	} else {
		return IntervalDomainValue{value: max(l1.value, l2.value)}
	}
}

//////////////////////////////////

type IntervalDomain struct {
	lower, upper IntervalDomainValue
}

func (domain *IntervalDomain) is_bounded() bool {
	return domain.lower.is_finite() && domain.upper.is_finite()
}

func (l1 *IntervalDomain) Join(l2 IntervalDomain) IntervalDomain {
	return IntervalDomain{lower: min_IntervalDomainValue(l1.lower, l2.lower), upper: max_IntervalDomainValue(l1.upper, l2.upper)}
}
