package handler

import (
	"math"
	"math/rand"
	"time"
)

// Randomly decides which chunks require payment.
func GeneratePaymentSchedule(
	totalChunks int,
) map[uint32]bool {

	paid := make(map[uint32]bool)

	if totalChunks == 0 {
		return paid
	}

	// Create a local random generator
	rng := rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	// Approximately 25% of chunks require payment
	paidCount := int(
		math.Ceil(
			float64(totalChunks) * 0.25,
		),
	)

	// Always require at least 2 payments
	if paidCount < 2 {
		paidCount = 2
	}

	// Cannot have more paid chunks than total chunks
	if paidCount > totalChunks {
		paidCount = totalChunks
	}

	// Randomly choose unique paid chunks
	for len(paid) < paidCount {

		index := uint32(
			rng.Intn(totalChunks),
		)

		paid[index] = true
	}

	return paid
}