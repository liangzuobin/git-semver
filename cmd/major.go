// Copyright © 2019 liangzuobin <liangzuobin123@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// majorCmd represents the major command
var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "generate a major version tag",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		sv, err := currentversiontag(ctx)
		if err != nil {
			fmt.Printf("get current semver failed: %v \n", err)
			return
		}
		sv.major++
		sv.minor = 0
		sv.patch = 0
		if err := addgittag(ctx, sv); err != nil {
			fmt.Printf("add new tag %s failed: %v", sv.tag(), err)
			return
		}
		fmt.Printf("current version: %s", sv.tag())
	},
}

func init() {
	rootCmd.AddCommand(majorCmd)
}
