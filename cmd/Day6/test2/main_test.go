package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestSafeMap(t *testing.T) {
	m := SafeMap{
		data: make(map[string]string),
	}
	m.Set("name", "xiaoming")
	got, ok := m.Get("name")
	if !ok {
		t.Fatal("key name 应该存在")
	}
	if got != "xiaoming" {
		t.Errorf("Get(name) = %v,期望: %v", got, "xiaoming")
	}
	m.Set("name", "lihua")

	got, ok = m.Get("name")
	if !ok {
		t.Fatal("覆盖后 key name 应该仍然存在")
	}
	if got != "lihua" {
		t.Errorf("Get(name) = %q,期望 %q", got, "lihua")
	}

	// 3. 删除后不应再读取到数据
	m.Delete("name")

	_, ok = m.Get("name")
	if ok {
		t.Error("删除后 key name 不应该存在")
	}

	// 4. 从未写入的 key 也应返回不存在
	_, ok = m.Get("unknown")
	if ok {
		t.Error("不存在的 key 不应该返回 ok=true")
	}

}

func TestSafeMapConcurrent(t *testing.T) {
	m := SafeMap{
		data: make(map[string]string),
	}
	const goroutineCount = 1000
	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("user-%d", id)
			value := fmt.Sprintf("name-%d", id)
			m.Set(key, value)
		}(i)
	}
	wg.Wait()

	// 第二阶段：检查并发写入结果
	for i := 0; i < goroutineCount; i++ {
		key := fmt.Sprintf("user-%d", i)
		want := fmt.Sprintf("name-%d", i)

		got, ok := m.Get(key)
		if !ok {
			t.Errorf("key %q 不存在", key)
			continue
		}

		if got != want {
			t.Errorf("Get(%q) = %q，期望 %q", key, got, want)
		}
	}

	// 第三阶段：并发删除
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()

			key := fmt.Sprintf("user-%d", id)
			m.Delete(key)
		}(i)
	}

	wg.Wait()

	// 第四阶段：确认所有数据都已删除
	for i := 0; i < goroutineCount; i++ {
		key := fmt.Sprintf("user-%d", i)

		if _, ok := m.Get(key); ok {
			t.Errorf("删除后 key %q 仍然存在", key)
		}
	}
}

func TestSafeMapMixedConcurrent(t *testing.T) {
	m := SafeMap{
		data: make(map[string]string),
	}

	const (
		keyCount       = 100
		workerCount    = 30
		operationCount = 1000
	)

	// 预先放入数据，让读和删除操作一开始就有数据可处理。
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key-%d", i)
		m.Set(key, "initial")
	}

	var wg sync.WaitGroup

	// 充当起跑信号，让所有 goroutine 尽量同时开始。
	start := make(chan struct{})

	// 30 个写 goroutine。
	for workerID := 0; workerID < workerCount; workerID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			<-start

			for i := 0; i < operationCount; i++ {
				key := fmt.Sprintf("key-%d", i%keyCount)
				m.Set(key, "updated")
			}
		}(workerID)
	}

	// 30 个读 goroutine。
	for workerID := 0; workerID < workerCount; workerID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			<-start

			for i := 0; i < operationCount; i++ {
				key := fmt.Sprintf("key-%d", i%keyCount)
				m.Get(key)
			}
		}(workerID)
	}

	// 30 个删除 goroutine。
	for workerID := 0; workerID < workerCount; workerID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			<-start

			for i := 0; i < operationCount; i++ {
				key := fmt.Sprintf("key-%d", i%keyCount)
				m.Delete(key)
			}
		}(workerID)
	}

	// 关闭 channel 后，所有 `<-start` 同时解除阻塞。
	close(start)

	wg.Wait()

	// 最终某个 key 可能存在，也可能被删除。
	// 如果存在，它只能是初始值或者写 goroutine 设置的新值。
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key-%d", i)

		value, ok := m.Get(key)
		if !ok {
			continue
		}

		if value != "initial" && value != "updated" {
			t.Errorf("key %q 得到非法值 %q", key, value)
		}
	}
}
func BenchmarkSafeMapGet(b *testing.B) {
	m := SafeMap{
		data: make(map[string]string),
	}
	m.Set("name", "xiaoming")

	// 将准备测试数据的时间排除在 Benchmark 之外。
	b.ResetTimer()

	// 输出每次操作的内存分配信息。
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m.Get("name")
	}
}

func BenchmarkSafeMapSet(b *testing.B) {
	m := SafeMap{
		data: make(map[string]string),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 始终覆盖同一个 key，避免把 map 扩容成本算进测试。
		m.Set("name", "xiaoming")
	}
}

func BenchmarkSafeMapGetParallel(b *testing.B) {
	m := SafeMap{
		data: make(map[string]string),
	}
	m.Set("name", "xiaoming")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get("name")
		}
	})
}

func BenchmarkSafeMapSetParallel(b *testing.B) {
	m := SafeMap{
		data: make(map[string]string),
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("name", "xiaoming")
		}
	})
}