/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package procfile_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/initializ/procfile/v5/procfile"
)

func testProcfile(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect         = NewWithT(t).Expect
		path, bindPath string
		bindings       libcnb.Bindings
	)

	it.Before(func() {
		var err error
		path, err = ioutil.TempDir("", "procfile")
		bindPath, err = ioutil.TempDir("", "bindProcfile")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	it("returns an empty Procfile when file does not exist", func() {
		Expect(procfile.NewProcfileFromPath(path)).To(HaveLen(0))
	})

	it("returns a parsed Profile", func() {
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte("test-type: test-command"), 0644)).To(Succeed())

		Expect(procfile.NewProcfileFromPath(path)).To(Equal(procfile.Procfile{"test-type": "test-command"}))
	})

	it("returns a Procfile from given binding", func() {
		Expect(ioutil.WriteFile(filepath.Join(bindPath, "Procfile"), []byte("test-type-bind: test-command"), 0644)).To(Succeed())
		bindings = libcnb.Bindings{libcnb.Binding{
			Name:   "name1",
			Type:   "Procfile",
			Path:   bindPath,
			Secret: map[string]string{"Procfile": filepath.Join(bindPath, "Procfile")},
		}}

		Expect(procfile.NewProcfileFromBinding(bindings)).To(Equal(procfile.Procfile{"test-type-bind": "test-command"}))
	})

	it("returns a Procfile with only file contents, if no binding", func() {
		bindings = libcnb.Bindings{}
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte("test-type-path: test-command"), 0644)).To(Succeed())

		Expect(procfile.NewProcfileFromPathOrBinding(path, bindings)).To(Equal(procfile.Procfile{"test-type-path": "test-command"}))

	})

	it("returns a Procfile with only binding contents, if no file", func() {
		Expect(ioutil.WriteFile(filepath.Join(bindPath, "Procfile"), []byte("test-type-bind: test-command"), 0644)).To(Succeed())
		bindings = libcnb.Bindings{libcnb.Binding{
			Name:   "name1",
			Type:   "Procfile",
			Path:   bindPath,
			Secret: map[string]string{"Procfile": filepath.Join(bindPath, "Procfile")},
		}}

		Expect(procfile.NewProcfileFromPathOrBinding(path, bindings)).To(Equal(procfile.Procfile{"test-type-bind": "test-command"}))

	})

	it("returns a merged Procfile from given binding + file", func() {
		Expect(ioutil.WriteFile(filepath.Join(bindPath, "Procfile"), []byte("test-type-bind: test-command"), 0644)).To(Succeed())
		bindings = libcnb.Bindings{libcnb.Binding{
			Name:   "name1",
			Type:   "Procfile",
			Path:   bindPath,
			Secret: map[string]string{"Procfile": filepath.Join(bindPath, "Procfile")},
		}}
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte("test-type-path: test-command"), 0644)).To(Succeed())

		Expect(procfile.NewProcfileFromPathOrBinding(path, bindings)).To(Equal(procfile.Procfile{"test-type-path": "test-command", "test-type-bind": "test-command"}))

	})

	it("returns a merged Procfile from given binding + file, binding takes precedence on duplicates", func() {
		Expect(ioutil.WriteFile(filepath.Join(bindPath, "Procfile"), []byte("test-type: bind-test-command\ntest-type-2: another-test-command"), 0644)).To(Succeed())
		bindings = libcnb.Bindings{libcnb.Binding{
			Name:   "name1",
			Type:   "Procfile",
			Path:   bindPath,
			Secret: map[string]string{"Procfile": filepath.Join(bindPath, "Procfile")},
		}}
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte("test-type: path-test-command"), 0644)).To(Succeed())

		Expect(procfile.NewProcfileFromPathOrBinding(path, bindings)).To(Equal(procfile.Procfile{"test-type": "bind-test-command", "test-type-2": "another-test-command"}))

	})

}
