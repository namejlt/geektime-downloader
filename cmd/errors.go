package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/namejlt/geektime-downloader/pconst"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/namejlt/geektime-downloader/internal/geektime"
)

func checkError(err error) {
	if err != nil {
		// special newline case
		if errors.Is(err, pconst.ErrGeekTimeRateLimit) ||
			os.IsTimeout(err) {
			fmt.Println()
		}

		var eg *geektime.ErrGeekTimeAPIBadCode
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, promptui.ErrInterrupt) {
			os.Exit(1)
		} else if errors.As(err, &eg) ||
			errors.Is(err, pconst.ErrWrongPassword) ||
			errors.Is(err, pconst.ErrTooManyLoginAttemptTimes) {
			exitWithMsg(err.Error())
		} else if errors.Is(err, pconst.ErrGeekTimeRateLimit) ||
			errors.Is(err, pconst.ErrAuthFailed) {
			exitWithMsg("请重新登陆")
		} else if os.IsTimeout(err) {
			exitWithMsg("请求超时")
		} else {
			fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err.Error())
			os.Exit(1)
		}
	}
}

func exitWithMsg(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
