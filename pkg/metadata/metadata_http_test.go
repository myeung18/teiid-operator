/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metadata

import (
	"testing"

	"github.com/teiid/teiid-operator/pkg/util/test"

	"github.com/stretchr/testify/assert"
	"github.com/teiid/teiid-operator/pkg/apis/teiid/v1alpha1"
)

func TestHttpJavaSource(t *testing.T) {
	code := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "Request.java",
			Content: `
			from("telegram:bots/cippa").to("log:stash");
			from("undertow:uri").to("log:stash");
			from("ine:xistent").to("log:stash");
		`,
		},
		Language: v1alpha1.LanguageJavaSource,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.False(t, meta.PassiveEndpoints)
}

func TestHttpOnlyJavaSource(t *testing.T) {
	code := v1alpha1.SourceSpec{

		DataSpec: v1alpha1.DataSpec{
			Name: "Request.java",
			Content: `
			from("direct:bots/cippa").to("log:stash");
			from("undertow:uri").to("log:stash");
			from("seda:path").to("log:stash");
		`,
		},
		Language: v1alpha1.LanguageJavaSource,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.True(t, meta.PassiveEndpoints)
}

func TestHttpOnlyJavaSourceRest(t *testing.T) {
	code := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "Request.java",
			Content: `
			from("direct:bots/cippa").to("log:stash");
			rest().get("").to("log:stash");
		`,
		},
		Language: v1alpha1.LanguageJavaSource,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.True(t, meta.PassiveEndpoints)
}

func TestHttpOnlyJavaSourceRest2(t *testing.T) {
	code := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "Request.java",
			Content: `
			from("vm:bots/cippa").to("log:stash");
			rest( ).get("").to("log:stash");
		`,
		},
		Language: v1alpha1.LanguageJavaSource,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.True(t, meta.PassiveEndpoints)
}

func TestNoHttpGroovySource(t *testing.T) {
	code := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "Request.groovy",
			Content: `
			from('direct:bots/cippa').to("log:stash");
			from('teelgram:uri').to("log:stash");
			from('seda:path').to("log:stash");
		`,
		},
		Language: v1alpha1.LanguageGroovy,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.False(t, meta.RequiresHTTPService)
	assert.False(t, meta.PassiveEndpoints)
}

func TestHttpOnlyGroovySource(t *testing.T) {
	code := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "Request.groovy",
			Content: `
			from('direct:bots/cippa').to("log:stash");
			from('undertow:uri').to("log:stash");
			from('seda:path').to("log:stash");
		`,
		},
		Language: v1alpha1.LanguageGroovy,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.True(t, meta.PassiveEndpoints)
}

func TestHttpXMLSource(t *testing.T) {
	code := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "routes.xml",
			Content: `
			<from uri="telegram:ciao" />
			<rest path="/">
			</rest>
		`,
		},
		Language: v1alpha1.LanguageXML,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.False(t, meta.PassiveEndpoints)
}

func TestHttpOnlyXMLSource(t *testing.T) {
	code := v1alpha1.SourceSpec{

		DataSpec: v1alpha1.DataSpec{
			Name: "routes.xml",
			Content: `
			<from uri="direct:ciao" />
			<rest path="/">
			</rest>
		`,
		},
		Language: v1alpha1.LanguageXML,
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := Extract(catalog, code)

	assert.True(t, meta.RequiresHTTPService)
	assert.True(t, meta.PassiveEndpoints)
}

func TestMultilangHTTPOnlySource(t *testing.T) {
	codes := []v1alpha1.SourceSpec{
		{
			DataSpec: v1alpha1.DataSpec{
				Name: "routes.xml",
				Content: `
				<from uri="direct:ciao" />
				<rest path="/">
				</rest>
			`,
			},
			Language: v1alpha1.LanguageXML,
		},
		{

			DataSpec: v1alpha1.DataSpec{
				Name: "routes2.groovy",
				Content: `
				from('seda:in').to('seda:out')
			`,
			},
			Language: v1alpha1.LanguageGroovy,
		},
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := ExtractAll(catalog, codes)

	assert.True(t, meta.RequiresHTTPService)
	assert.True(t, meta.PassiveEndpoints)
}

func TestMultilangHTTPSource(t *testing.T) {
	codes := []v1alpha1.SourceSpec{
		{

			DataSpec: v1alpha1.DataSpec{
				Name: "routes.xml",
				Content: `
				<from uri="direct:ciao" />
				<rest path="/">
				</rest>
			`,
			},
			Language: v1alpha1.LanguageXML,
		},
		{

			DataSpec: v1alpha1.DataSpec{
				Name: "routes2.groovy",
				Content: `
				from('seda:in').to('seda:out')
				from('timer:tick').to('log:info')
			`,
			},
			Language: v1alpha1.LanguageGroovy,
		},
	}

	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	meta := ExtractAll(catalog, codes)

	assert.True(t, meta.RequiresHTTPService)
	assert.False(t, meta.PassiveEndpoints)
}
