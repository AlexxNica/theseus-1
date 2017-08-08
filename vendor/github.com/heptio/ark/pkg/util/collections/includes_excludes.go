/*
Copyright 2017 Heptio Inc.

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

package collections

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

// IncludesExcludes is a type that manages lists of included
// and excluded items. The logic implemented is that everything
// in the included list except those items in the excluded list
// should be included. '*' in the includes list means "include
// everything", but it is not valid in the exclude list.
type IncludesExcludes struct {
	includes sets.String
	excludes sets.String
}

func NewIncludesExcludes() *IncludesExcludes {
	return &IncludesExcludes{
		includes: sets.NewString(),
		excludes: sets.NewString(),
	}
}

// Includes adds items to the includes list. '*' is a wildcard
// value meaning "include everything".
func (ie *IncludesExcludes) Includes(includes ...string) *IncludesExcludes {
	ie.includes.Insert(includes...)
	return ie
}

// GetIncludes returns the items in the includes list
func (ie *IncludesExcludes) GetIncludes() []string {
	return ie.includes.List()
}

// Excludes adds items to the excludes list
func (ie *IncludesExcludes) Excludes(excludes ...string) *IncludesExcludes {
	ie.excludes.Insert(excludes...)
	return ie
}

// GetExcludes returns the items in the excludes list
func (ie *IncludesExcludes) GetExcludes() []string {
	return ie.excludes.List()
}

// ShouldInclude returns whether the specified item should be
// included or not. Everything in the includes list except those
// items in the excludes list should be included.
func (ie *IncludesExcludes) ShouldInclude(s string) bool {
	if ie.excludes.Has(s) {
		return false
	}

	return ie.includes.Has("*") || ie.includes.Has(s)
}

func ValidateIncludesExcludes(includesList, excludesList []string) []error {
	// TODO we should not allow an IncludesExcludes object to be created that
	// does not meet these criteria. Do a more significant refactoring to embed
	// this logic in object creation/modification.

	var errs []error

	includes := sets.NewString(includesList...)
	excludes := sets.NewString(excludesList...)

	if includes.Len() == 0 {
		errs = append(errs, errors.New("includes list cannot be empty"))
	}

	if includes.Len() > 1 && includes.Has("*") {
		errs = append(errs, errors.New("includes list must either contain '*' only, or a non-empty list of items"))
	}

	if excludes.Has("*") {
		errs = append(errs, errors.New("excludes list cannot contain '*'"))
	}

	for _, itm := range excludes.List() {
		if includes.Has(itm) {
			errs = append(errs, errors.New(fmt.Sprintf("excludes list cannot contain an item in the includes list: %v", itm)))
		}
	}

	return errs
}
