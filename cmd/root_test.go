package cmd

import (
	"context"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegExAndSort(t *testing.T) {
	tagexps := [][]string{
		[]string{"v0.0.1", "v", "0", "0", "1", ""},
		[]string{"v2.1.2", "v", "2", "1", "2", ""},
		[]string{"2.1.1", "", "2", "1", "1", ""},
		[]string{"2.0.1", "", "2", "0", "1", ""},
		[]string{"v3.0.0-ga", "v", "3", "0", "0", "-ga"},
		[]string{"4.0.4rc", "", "4", "0", "4", "rc"},
	}
	svs := make([]semver, 0, len(tagexps))
	for _, tagexp := range tagexps {
		assert.True(t, reg.MatchString(tagexp[0]), "%s not match", tagexp[0])
		sv, err := parsesemver([]byte(tagexp[0]))
		assert.NoError(t, err)
		assert.Equal(t, tagexp[1], sv.prefix)
		assert.Equal(t, tagexp[5], sv.suffix)
		assert.Equal(t, tagexp[2], strconv.Itoa(sv.major))
		assert.Equal(t, tagexp[3], strconv.Itoa(sv.minor))
		assert.Equal(t, tagexp[4], strconv.Itoa(sv.patch))
		svs = append(svs, sv)
	}

	sort.Sort(semvers(svs))
	assert.Equal(t, "4.0.4rc", svs[0].tag())
	assert.Equal(t, "v3.0.0-ga", svs[1].tag())
	assert.Equal(t, "v2.1.2", svs[2].tag())
	assert.Equal(t, "2.1.1", svs[3].tag())
	assert.Equal(t, "2.0.1", svs[4].tag())
	assert.Equal(t, "v0.0.1", svs[5].tag())
}

func TestAddGitTag(t *testing.T) {
	ctx := context.Background()
	sv, err := currentversiontag(ctx)
	assert.NoError(t, err)

	sv.patch++
	err = addgittag(context.Background(), sv, "test")
	assert.NoError(t, err)
}

func TestNewTag(t *testing.T) {
	sv := semver{major: 1, minor: 1, patch: 1}

	sv.newtag(patch)
	assert.Equal(t, 1, sv.major)
	assert.Equal(t, 1, sv.minor)
	assert.Equal(t, 2, sv.patch)

	sv.newtag(minor)
	assert.Equal(t, 1, sv.major)
	assert.Equal(t, 2, sv.minor)
	assert.Equal(t, 0, sv.patch)

	sv.newtag(major)
	assert.Equal(t, 2, sv.major)
	assert.Equal(t, 0, sv.minor)
	assert.Equal(t, 0, sv.patch)
}
