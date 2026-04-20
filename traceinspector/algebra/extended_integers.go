package algebra

import "strconv"

// This file defines an "extended" integer type, the integer set Z augmented with positive and negative infinity

type ExtInt struct {
	Value              int
	is_inf, is_neg_inf bool
}

func (eint ExtInt) String() string {
	if eint.is_inf {
		return "∞"
	} else if eint.is_neg_inf {
		return "-∞"
	} else {
		return strconv.Itoa(eint.Value)
	}
}

func ExtInt_Finite(val int) ExtInt {
	return ExtInt{Value: val}
}

func ExtInt_Zero() ExtInt {
	return ExtInt{Value: 0}
}

func ExtInt_Infty() ExtInt {
	return ExtInt{is_inf: true}
}

func ExtInt_NegInfty() ExtInt {
	return ExtInt{is_neg_inf: true}
}

func (eint ExtInt) IsFinite() bool {
	return !(eint.is_inf || eint.is_neg_inf)
}

func (eint ExtInt) IsPositive() bool {
	return (eint.IsFinite() && eint.Value > 0) || eint.is_inf
}

func (eint ExtInt) IsNegative() bool {
	return (eint.IsFinite() && eint.Value < 0) || eint.is_neg_inf
}

func (lhs ExtInt) Eq(rhs ExtInt) bool {
	return (lhs.is_inf && rhs.is_inf) || (lhs.is_neg_inf && rhs.is_neg_inf) || (lhs.IsFinite() && rhs.IsFinite() && lhs.Value == rhs.Value)
}

func (lhs ExtInt) Leq(rhs ExtInt) bool {
	// trivial case
	if lhs.is_neg_inf || rhs.is_inf {
		return true
	}
	if lhs.is_inf {
		return rhs.is_inf
	}
	if rhs.is_neg_inf {
		return lhs.is_neg_inf
	}

	// remaining case is lhs, rhs = Z
	return lhs.Value <= rhs.Value
}

func (eint ExtInt) Neg() ExtInt {
	if eint.is_inf {
		return ExtInt_NegInfty()
	} else if eint.is_neg_inf {
		return ExtInt_Infty()
	} else {
		return ExtInt{Value: -eint.Value}
	}
}

func (lhs ExtInt) Add(rhs ExtInt) ExtInt {
	if lhs.is_inf || rhs.is_inf {
		return ExtInt_Infty()
	}
	if lhs.is_neg_inf || rhs.is_neg_inf {
		return ExtInt_NegInfty()
	}
	return ExtInt{Value: lhs.Value + rhs.Value}
}

func (lhs ExtInt) Sub(rhs ExtInt) ExtInt {
	// note that I model infty ± infty as infty, and -infty ± infty as -infty. This is not mathematically correct
	if lhs.is_inf || rhs.is_inf {
		return ExtInt_Infty()
	}
	if lhs.is_neg_inf || rhs.is_neg_inf {
		return ExtInt_NegInfty()
	}
	return ExtInt{Value: lhs.Value - rhs.Value}
}

func (lhs ExtInt) Mul(rhs ExtInt) ExtInt {
	if lhs.Eq(ExtInt_Zero()) || rhs.Eq(ExtInt_Zero()) {
		// Note that 0 * infty is undefined, but we return 0. Again this is not mathematically correct
		return ExtInt_Zero()
	}
	if lhs.IsFinite() && rhs.IsFinite() {
		return ExtInt{Value: lhs.Value * rhs.Value}
	}
	// at least one value is +- inf, so the value is inf; just have to define the sign
	switch lhs.IsPositive() {
	case true:
		switch rhs.IsPositive() {
		case true:
			return ExtInt_Infty()
		case false:
			return ExtInt_NegInfty()
		}
	case false:
		// lhs = -
		switch rhs.IsNegative() {
		case true:
			return ExtInt_Infty()
		case false:
			return ExtInt_NegInfty()
		}
	}
	panic("This should never ever happen")
}
