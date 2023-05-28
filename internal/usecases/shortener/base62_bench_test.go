package shortener

import (
	"testing"

	"github.com/google/uuid"
)

func Benchmark_encodeUUID(b *testing.B) {
	uuidStr := "624708fa-d258-4b99-b09a-49d95f294626"
	uuidObj, err := uuid.Parse(uuidStr)
	if err != nil {
		panic(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodeUUID(uuidObj)
	}
}

func Benchmark_decodeUUID(b *testing.B) {
	idStr := "2ZrI5IHFnvPscPYKlxFtRQ"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := decodeUUID(idStr)
		if err != nil {
			panic(err)
		}
	}
}
