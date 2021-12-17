package dblayer

import "testing"

func BenchmarkHashPassword(b *testing.B) {
	// 해싱할 문자열
	text := "A string to be Hashed"
	// 소요된 시간 초기화 (앞에서 변수 초기화하는데 걸리는 시간을 제거하기 위함)
	b.ResetTimer()
	// testing.B타입을 사용해 해당 코드를 b.N번 실행합니다.
	// b.N의 값은 코드의 성능을 측정할 수 있는 적절한 값으로 자동 조정됩니다.
	for i := 0; i < b.N; i++ {
		hashPassword(&text)
	}
	// 병렬 테스트 예제 go test -cpu {core} -bench .
	//b.RunParallel(func(pb *testing.PB) {
	//	for pb.Next() {
	//		hashPassword(&text)
	//	}
	//})
}
