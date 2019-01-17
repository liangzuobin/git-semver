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
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitsemver",
	Short: "git semver kit",
	Long:  `git sub command to generate semver tags`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("root called")
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type semver struct {
	prefix []byte
	major  int
	minor  int
	patch  int
	suffix []byte
}

func (s semver) tag() string {
	return fmt.Sprintf("%s%d.%d.%d%s", string(s.prefix), s.major, s.minor, s.patch, string(s.suffix))
}

type semvers []semver

func (s semvers) Len() int { return len(s) }

func (s semvers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// desc
func (s semvers) Less(i, j int) (b bool) {
	if s[i].major != s[j].major {
		return s[i].major > s[j].major
	}
	if s[i].minor != s[j].major {
		return s[i].minor > s[j].major
	}
	return s[i].patch > s[j].patch
}

var reg = regexp.MustCompile(`(\D*?)(\d*?)\.(\d*?)\.(\d*)(\D*)`)
var errnotmatch = errors.New("not match")

func parsesemver(tag []byte) (semver, error) {
	s := reg.FindSubmatch(tag)
	m, err := strconv.Atoi(string(s[2]))
	if err != nil {
		return semver{}, fmt.Errorf("parse major %s failed, err: %v", string(s[2]), err)
	}
	n, err := strconv.Atoi(string(s[3]))
	if err != nil {
		return semver{}, fmt.Errorf("parse minor %s failed, err: %v", string(s[3]), err)
	}
	p, err := strconv.Atoi(string(s[4]))
	if err != nil {
		return semver{}, fmt.Errorf("parse patch %s failed, err: %v", string(s[4]), err)
	}
	return semver{prefix: s[1], major: m, minor: n, patch: p, suffix: s[5]}, nil
}

func currentversiontag(ctx context.Context) (semver, error) {
	r, err := gittags(ctx)
	if err != nil {
		return semver{}, err
	}
	return currenttag(ctx, r)
}

func gittags(ctx context.Context) (io.Reader, error) {
	cmd := exec.CommandContext(ctx, "git", "tag", "--sort=v:refname")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &out, nil
}

func currenttag(ctx context.Context, r io.Reader) (semver, error) {
	rd := bufio.NewReader(r)
	svs := make([]semver, 0, 10)
	for {
		b, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return semver{}, err
		}
		if len(b) > 0 && reg.Match(b) {
			sv, err := parsesemver(b)
			if err != nil {
				fmt.Printf("tag %s not valid, will be ignored, err: %v", string(b), err)
				continue
			}
			svs = append(svs, sv)
		}
	}
	if len(svs) == 0 {
		return semver{prefix: []byte("v"), major: 0, minor: 0, patch: 0, suffix: []byte("")}, nil
	}
	sort.Sort(semvers(svs))
	return svs[0], nil
}

func addgittag(ctx context.Context, sv semver) error {
	cmd := exec.CommandContext(ctx, "git", "tag", "-a", sv.tag(), "-m", sv.tag()) // FIXME(liangzuobin) use a real message in args?
	// var out bytes.Buffer
	// cmd.Stdout = &out
	return cmd.Run()
}