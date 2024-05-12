package states

import "errors"

// LoanState represents the state of the loan.
type LoanState string

const (
	proposed  LoanState = "proposed"
	approved  LoanState = "approved"
	invested  LoanState = "invested"
	disbursed LoanState = "disbursed"
)

func InitialState() LoanState {
	return proposed
}

func (ls *LoanState) Approve() error {
	if *ls == proposed {
		*ls = approved
		return nil
	}
	return errors.New("invalid state")
}

func (ls *LoanState) Invest() error {
	if *ls == approved {
		*ls = invested
		return nil
	}
	return errors.New("invalid state")
}

func (ls *LoanState) Disburse() error {
	if *ls == invested {
		*ls = disbursed
		return nil
	}
	return errors.New("invalid state")
}

//
//
//func NextState(current LoanState) (LoanState, error) {
//	switch current {
//	case proposed:
//		return approved, nil
//	case approved:
//		return invested, nil
//	case invested:
//		return disbursed, nil
//	case disbursed:
//		return "", errors.New("reaching end state")
//	default:
//		return "", errors.New("unknown input state")
//	}
//}
