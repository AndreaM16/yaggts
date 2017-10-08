package util

import (
	"os"
	"errors"
)

const csvName = "multiTimeline.csv"

func DeleteCsv () error {
	err := os.Remove(csvName); if err != nil {
		return errors.New("unable to remove csv")
	}
	return nil
}
