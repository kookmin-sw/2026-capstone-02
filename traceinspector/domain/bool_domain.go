package domain

import (
	"fmt"
	"traceinspector/imp"
)

type BoolDomain struct {
	val            bool
	is_bot, is_top bool
}

func (domain BoolDomain) String() string {
	if domain.is_bot {
		return "⊥_bool"
	} else if domain.is_top {
		return "⊤_bool"
	} else {
		return fmt.Sprintf("%t", domain.val)
	}
}

func (domain BoolDomain) IsTrue() bool {
	return !domain.IsBot() && !domain.IsTop() && domain.val
}

func (domain BoolDomain) IsFalse() bool {
	return !domain.IsBot() && !domain.IsTop() && !domain.val
}

func (domain BoolDomain) Clone() BoolDomain {
	return BoolDomain{val: domain.val, is_bot: domain.is_bot, is_top: domain.is_top}
}

func (domain BoolDomain) CreateTop() BoolDomain {
	return BoolDomain{is_top: true}
}

func (domain BoolDomain) CreateBot() BoolDomain {
	return BoolDomain{is_bot: true}
}

func (domain BoolDomain) From_BoolLitExpr(expr imp.BoolLitExpr) BoolDomain {
	return BoolDomain{val: expr.Value}
}

func (domain BoolDomain) IsBot() bool {
	return domain.is_bot
}

func (domain BoolDomain) IsTop() bool {
	return domain.is_top
}

func (lhs BoolDomain) Join(rhs BoolDomain) (BoolDomain, bool) {
	if lhs.is_top || rhs.is_top {
		return BoolDomain{is_top: true}, false
	} else if lhs.is_bot {
		return BoolDomain{val: rhs.val, is_bot: rhs.is_bot}, rhs.is_bot
	} else if rhs.is_bot {
		return BoolDomain{val: lhs.val, is_bot: lhs.is_bot}, lhs.is_bot
	} else if lhs.val == rhs.val {
		return BoolDomain{val: lhs.val}, false
	} else {
		return BoolDomain{is_top: true}, lhs.IsTop()
	}
}

func (lhs BoolDomain) Incl(rhs BoolDomain) bool {
	if rhs.is_top {
		return true
	} else if lhs.is_bot {
		return true
	} else if lhs.is_top {
		return rhs.is_top
	} else if rhs.is_bot { // lhs concrete
		return false
	} else { // lhs = concrete, rhs = concrete
		return lhs.val == rhs.val
	}
}

func (lhs BoolDomain) Widen(rhs BoolDomain) BoolDomain {
	if lhs.is_bot {
		return rhs
	}
	if rhs.is_bot {
		return lhs
	}
	if lhs.is_top || rhs.is_top {
		return BoolDomain{is_top: true}
	}
	if lhs.val == rhs.val {
		return lhs
	} else {
		return BoolDomain{is_top: true}
	}
}

// expression evaluation
