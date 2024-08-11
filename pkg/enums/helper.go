package enums

import (
	pb "budgeting-service/generated/budgeting"
	"fmt"
)

// TransactionType represents the type of a transaction
func TypeTrnsactionParse(transactionTypeString string) (*pb.TypeTrnsaction, error) {
	var transactionType pb.TypeTrnsaction

	switch transactionTypeString {
	case "INCOME":
		transactionType = pb.TypeTrnsaction_INCOME
	case "EXPENSE":
		transactionType = pb.TypeTrnsaction_EXPENSE
	default:
		return nil, fmt.Errorf("invalid transaction type: %s", transactionTypeString)
	}

	return &transactionType, nil
}

func CategoryTypeParse(categoryTypeString string) (*pb.CategoryType, error) {
	var categoryType pb.CategoryType

	switch categoryTypeString {
	case "INCOME":
		categoryType = pb.CategoryType_INCOME
	case "EXPENSE":
		categoryType = pb.CategoryType_EXPENSE
	default:
		return nil, fmt.Errorf("invalid category type: %s", categoryTypeString)
	}

	return &categoryType, nil
}

func PeriodParse(periodString string) (*pb.Period, error) {
	var period pb.Period

	switch periodString {
	case "DAILY":
		period = pb.Period_DAILY
	case "WEEKLY":
		period = pb.Period_WEEKLY
	case "MONTHLY":
		period = pb.Period_MONTHLY
	case "YEARLY":
		period = pb.Period_YEARLY
		return nil, fmt.Errorf("invalid period: %s", periodString)
	}

	return &period, nil
}

func GoalStatusParse(statusString string) (*pb.Status, error) {
	var status pb.Status

    switch statusString {
    case "INPROGRESS":
        status = pb.Status_INPORGRESS
    case "ACHIEVED":
        status = pb.Status_ACHIEVED
    case "FAILED":
        status = pb.Status_FAILED
    default:
        return nil, fmt.Errorf("invalid goal status: %s", statusString)
    }

    return &status, nil
}