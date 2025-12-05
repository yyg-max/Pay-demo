/*
 * MIT License
 *
 * Copyright (c) 2025 linux.do
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package worker

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/linux-do/pay/internal/apps/dispute"
	"github.com/linux-do/pay/internal/apps/payment"
	"github.com/linux-do/pay/internal/apps/user"
	"github.com/linux-do/pay/internal/config"
	"github.com/linux-do/pay/internal/task"
)

// StartWorker 启动任务处理服务器
func StartWorker() error {
	asynqServer := asynq.NewServer(
		task.RedisOpt,
		asynq.Config{
			Concurrency:     config.Config.Worker.Concurrency,
			ShutdownTimeout: 3 * time.Minute,
			Queues: map[string]int{
				"critical": 10,
				"default":  5,
				"low":      1,
			},
			StrictPriority: true,
		},
	)

	// 注册任务处理器
	mux := asynq.NewServeMux()
	mux.Use(taskLoggingMiddleware)
	mux.HandleFunc(task.UpdateUserGamificationScoresTask, user.HandleUpdateUserGamificationScores)
	mux.HandleFunc(task.UpdateSingleUserGamificationScoreTask, user.HandleUpdateSingleUserGamificationScore)
	mux.HandleFunc(task.AutoRefundExpiredDisputesTask, dispute.HandleAutoRefundExpiredDisputes)
	mux.HandleFunc(task.AutoRefundSingleDisputeTask, dispute.HandleAutoRefundSingleDispute)
	mux.HandleFunc(task.MerchantPaymentNotifyTask, payment.HandleMerchantPaymentNotify)
	// 启动服务器
	return asynqServer.Run(mux)
}
