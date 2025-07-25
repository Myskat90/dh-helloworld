/*
Copyright 2025 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/deckhouse/module-sdk/pkg"
)

func ReadinessFunc(ctx context.Context, input *pkg.HookInput) error {
	input.Logger.Info("start user logic for readiness probe")

	c := input.DC.GetHTTPClient()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1/readyz", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	input.Logger.Debug("readiness probe done successfully", slog.Any("body", string(respBody)))

	return nil
}
