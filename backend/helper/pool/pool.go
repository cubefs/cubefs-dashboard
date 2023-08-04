// Copyright 2023 The CubeFS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package pool

type TaskPool struct {
	pool chan func()
}

func New(poolSize, poolCacheSize int) TaskPool {
	pool := make(chan func(), poolCacheSize)
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				task, ok := <-pool
				if !ok {
					break
				}
				task()
			}
		}()
	}

	return TaskPool{pool: pool}
}

func (t TaskPool) Run(task func()) {
	t.pool <- task
}

func (t TaskPool) TryRun(task func())  bool {
	select {
	case t.pool<-task:
		return true
	default:
		return false
	}
}

func (t TaskPool) Close()  {
	close(t.pool)
}

func GetPoolSize(poolSize, length int) int  {
	if length < poolSize {
		return length
	}
	return poolSize
}