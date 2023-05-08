# 標準パッケージ

## Command

command を実行する。  


```Go
func (c *Cmd) Run() error {
	// process を生成
	if err := c.Start(); err != nil {
		return err
	}

	// process の終了を待機
	return c.Wait()
}
```

プロセスを開始する。

```Go
func (c *Cmd) Start() error {
	...

	// contextがキャンセルされていたら終了する。
	if c.ctx != nil {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
		}
	}

	// I/Oを設定する。
	childFiles := make([]*os.File, 0, 3+len(c.ExtraFiles))
	stdin, err := c.childStdin()
	if err != nil {
		return err
	}
	childFiles = append(childFiles, stdin)
	stdout, err := c.childStdout()
	if err != nil {
		return err
	}
	childFiles = append(childFiles, stdout)
	stderr, err := c.childStderr(stdout)
	if err != nil {
		return err
	}
	childFiles = append(childFiles, stderr)
	childFiles = append(childFiles, c.ExtraFiles...)

	// env を設定する。
	env, err := c.environ()
	if err != nil {
		return err
	}

	// process を開始してpidなどの情報を返す。
	c.Process, err = os.StartProcess(c.Path, c.argv(), &os.ProcAttr{
		Dir:   c.Dir,
		Files: childFiles,
		Env:   env,
		Sys:   c.SysProcAttr,
	})
	if err != nil {
		return err
	}
	started = true

	// Don't allocate the goroutineErr channel unless there are goroutines to start.
	if len(c.goroutine) > 0 {
		goroutineErr := make(chan error, 1)
		c.goroutineErr = goroutineErr

		type goroutineStatus struct {
			running  int
			firstErr error
		}
		statusc := make(chan goroutineStatus, 1)
		statusc <- goroutineStatus{running: len(c.goroutine)}
		for _, fn := range c.goroutine {
			go func(fn func() error) {
				// 処理は並列で行っている。
				err := fn()

				// 一番最初のエラーと稼働中のgoroutine数を管理する。
				status := <-statusc
				if status.firstErr == nil {
					status.firstErr = err
				}
				status.running--
				if status.running == 0 {
					goroutineErr <- status.firstErr
				} else {
					statusc <- status
				}
			}(fn)
		}
		c.goroutine = nil // Allow the goroutines' closures to be GC'd when they complete.
	}

	if (c.Cancel != nil || c.WaitDelay != 0) && c.ctx != nil && c.ctx.Done() != nil {
		// バッファなしチャネルを生成。
		resultc := make(chan ctxResult)
		c.ctxResult = resultc
		go c.watchCtx(resultc)
	}

	return nil
}
```

プロセスの終了を待機する。

```Go
func (c *Cmd) Wait() error {
	if c.Process == nil {
		return errors.New("exec: not started")
	}
	if c.ProcessState != nil {
		return errors.New("exec: Wait was already called")
	}

	// process が終了するまで待機する。
	state, err := c.Process.Wait()
	if err == nil && !state.Success() {
		err = &ExitError{ProcessState: state}
	}
	c.ProcessState = state

	var timer *time.Timer
	if c.ctxResult != nil {
		// c.Start()の中で生成したバッファなしチャネルから受信する。
		watch := <-c.ctxResult
		timer = watch.timer

		if err == nil && watch.err != nil {
			err = watch.err
		}
	}

	if goroutineErr := c.awaitGoroutines(timer); err == nil {
		err = goroutineErr
	}
	closeDescriptors(c.parentIOPipes)
	c.parentIOPipes = nil

	return err
}
```

プロセスの完了とキャンセルを監視する。

```Go
func (c *Cmd) watchCtx(resultc chan<- ctxResult) {
	select {
	// c.Wait()の watch := <-c.ctxResult が実行されたら実行される。
	case resultc <- ctxResult{}:
		return
	// context がキャンセルされたら処理を継続。
	case <-c.ctx.Done():
	}

	var err error
	if c.Cancel != nil {
		if interruptErr := c.Cancel(); interruptErr == nil {
			err = c.ctx.Err()
		} else if errors.Is(interruptErr, os.ErrProcessDone) {
		} else {
			err = wrappedError{
				prefix: "exec: canceling Cmd",
				err:    interruptErr,
			}
		}
	}
	if c.WaitDelay == 0 {
		resultc <- ctxResult{err: err}
		return
	}

	timer := time.NewTimer(c.WaitDelay)
	select {
	case resultc <- ctxResult{err: err, timer: timer}:
		return
	case <-timer.C:
	}

	killed := false
	// process を kill する。
	if killErr := c.Process.Kill(); killErr == nil {
		killed = true
	} else if !errors.Is(killErr, os.ErrProcessDone) {
		err = wrappedError{
			prefix: "exec: killing Cmd",
			err:    killErr,
		}
	}

	if c.goroutineErr != nil {
		select {
		case goroutineErr := <-c.goroutineErr:
			if err == nil && !killed {
				err = goroutineErr
			}
		default:
			closeDescriptors(c.parentIOPipes)
			_ = <-c.goroutineErr
			if err == nil {
				err = ErrWaitDelay
			}
		}

		c.goroutineErr = nil
	}

	resultc <- ctxResult{err: err}
}
```
