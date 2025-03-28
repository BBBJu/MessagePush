package service

import (
	"fmt"
	"messagePush/utils"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// 压力测试参数配置
const (
	TotalMessages  = 1000 // 总消息量
	Workers        = 10   // 并发协程数
	MaxQPS         = 500  // 最大QPS限制
	MessageSubject = "压力测试"
)

// 原子计数器
var (
	sentCount    atomic.Int64
	SuccessCount atomic.Int64
	FailedCount  atomic.Int64
	maxTPS       atomic.Int64
	currentTPS   atomic.Int64
)
var (
	//实际处理成功的数据，需要埋点
	ProcessCount atomic.Int64
)

// 启动压力测试
func StartStressTest() {
	fmt.Printf("开始压力测试，总量：%d，并发数：%d\n", TotalMessages, Workers)

	// 初始化计数器
	resetCounters()

	var wg sync.WaitGroup
	startTime := time.Now()

	// 启动监控协程
	go monitorPerformance(&wg, startTime)
	wg.Add(1)

	// 启动工作协程
	for i := 0; i < Workers; i++ {
		wg.Add(1)
		// 创建限流器
		limiter := time.Tick(time.Second / time.Duration(MaxQPS/Workers))
		go func() {
			defer wg.Done()
			for sentCount.Load() < int64(TotalMessages) {
				<-limiter
				sendTestMessage()
				sentCount.Add(1)
			}
		}()
	}

	wg.Wait()

	// 输出最终报告
	printFinalReport(startTime)
}

// 发送测试消息
func sendTestMessage() {
	params := CreateMessageParams{
		subject:      MessageSubject,
		to:           utils.ReceiveId784,
		channel:      1,
		sourceID:     "stress_test",
		templateId:   4,
		templateData: `{"username":"784"}`,
		priority:     VIPPriority,
	}

	CreateMessage(params)
}

// 监控性能指标
// 修改监控函数退出逻辑
func monitorPerformance(wg *sync.WaitGroup, start time.Time) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastCount int64

	// 添加退出条件检测
	for range ticker.C {
		if SuccessCount.Load()+FailedCount.Load() >= int64(TotalMessages) {
			wg.Done()
			return
		}

		current := SuccessCount.Load() + FailedCount.Load()
		tps := current - lastCount
		lastCount = current

		if tps > maxTPS.Load() {
			maxTPS.Store(tps)
		}
		currentTPS.Store(tps)

		printLiveStatus(start)
	}
}

var (
	resultFile     *os.File
	resultFileOnce sync.Once
)

func printFinalReport(start time.Time) {
	elapsed := time.Since(start).Seconds()
	content := fmt.Sprintf(`
测试完成，耗时 %.2f 秒
发送TPS:%d(%.2f msg/s)
总处理量：%d (%.2f msg/s)
最大瞬时TPS:%d`,
		elapsed,
		sentCount.Load(),
		float64(sentCount.Load())/float64(elapsed),
		SuccessCount.Load(),
		float64(SuccessCount.Load())/elapsed,
		maxTPS.Load(),
	)

	if err := writeToFile(content); err != nil {
		fmt.Printf("写入测试结果失败: %v", err)
	}
}

// 新增文件写入函数
func writeToFile(content string) error {
	resultFileOnce.Do(func() {
		fileName := fmt.Sprintf("TestLog/stress_test.log_%s.log", time.Now().Format("20060102150405"))
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			resultFile = f
		}
	})

	if resultFile == nil {
		return fmt.Errorf("无法创建日志文件")
	}

	_, err := resultFile.WriteString(content + "\n")
	return err
}

// 修改实时状态输出
func printLiveStatus(start time.Time) {
	elapsed := time.Since(start).Seconds()
	status := fmt.Sprintf("[%.1fs] 发送: %d/%d | 成功: %d | 失败: %d | 当前TPS: %d | 最大TPS: %d",
		elapsed,
		sentCount.Load(),
		TotalMessages,
		SuccessCount.Load(),
		FailedCount.Load(),
		currentTPS.Load(),
		maxTPS.Load(),
	)

	if err := writeToFile(status); err != nil {
		fmt.Printf("写入状态失败: %v", err)
	}
}

// 重置计数器
func resetCounters() {
	sentCount.Store(0)
	SuccessCount.Store(0)
	FailedCount.Store(0)
	maxTPS.Store(0)
	currentTPS.Store(0)
}
