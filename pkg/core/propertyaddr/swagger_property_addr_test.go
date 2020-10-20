package propertyaddr

import (
	"testing"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/utils"

	"github.com/stretchr/testify/require"
)

func TestParseSwaggerRelPropertyAddr(t *testing.T) {
	cases := []struct {
		input  string
		expect SwaggerRelPropertyAddr
		error  bool
	}{
		{
			input:  "",
			expect: nil,
		},
		{
			input:  "p1",
			expect: SwaggerRelPropertyAddr{{name: "p1"}},
		},
		{
			input:  "p1.p2",
			expect: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}},
		},
		{
			input:  "p1.p2{v1}",
			expect: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}},
		},
		{
			input: "p1.p2}",
			error: true,
		},
	}

	for idx, c := range cases {
		actual, err := ParseSwaggerRelPropertyAddr(c.input)
		if c.error {
			require.Error(t, err, idx)
			continue
		}
		require.NoError(t, err, idx)
		require.Equal(t, c.expect, actual, idx)
	}
}

func TestSwaggerRelPropertyAddr_String(t *testing.T) {
	cases := []struct {
		input  SwaggerRelPropertyAddr
		expect string
	}{
		{
			input:  SwaggerRelPropertyAddr{},
			expect: "",
		},
		{
			input:  SwaggerRelPropertyAddr{{name: "p1"}},
			expect: "p1",
		},
		{
			input:  SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}},
			expect: "p1.p2",
		},
		{
			input:  SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}},
			expect: "p1.p2{v1}",
		},
	}
	for idx, c := range cases {
		require.Equal(t, c.expect, c.input.String(), idx)
	}
}

func TestSwaggerPropertyAddr_Append(t *testing.T) {
	cases := []struct {
		input  SwaggerPropertyAddr
		oaddr  string
		expect SwaggerPropertyAddr
		error  bool
	}{
		{
			input:  SwaggerPropertyAddr{},
			oaddr:  "",
			expect: SwaggerPropertyAddr{},
		},
		{
			input:  SwaggerPropertyAddr{},
			oaddr:  "p1",
			expect: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}}},
		},
		{
			input:  SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}}},
			oaddr:  "p2",
			expect: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
		},
		{
			input:  SwaggerPropertyAddr{},
			oaddr:  "p1.p2",
			expect: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
		},
		{
			input:  SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}}},
			oaddr:  "p2{v1}",
			expect: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}}},
		},
		{
			input: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}}},
			oaddr: "p2}",
			error: true,
		},
	}

	for idx, c := range cases {
		actual, err := c.input.Append(c.oaddr)
		if c.error {
			require.Error(t, err, idx)
			continue
		}
		require.NoError(t, err, idx)
		require.Equal(t, c.expect, actual, idx)
	}
}

func TestParseSwaggerPropertyAddr(t *testing.T) {
	cases := []struct {
		input  string
		expect SwaggerPropertyAddr
		error  bool
	}{
		{
			input:  "",
			expect: SwaggerPropertyAddr{},
		},
		{
			input:  "schema1:",
			expect: SwaggerPropertyAddr{Schema: "schema1"},
		},
		{
			input:  "p1",
			expect: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}}},
		},
		{
			input:  "p1.p2",
			expect: SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
		},
		{
			input:  "schema1:p1.p2",
			expect: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
		},
		{
			input:  "schema1:p1.p2{v1}",
			expect: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}}},
		},
		{
			input:  "schema1:p1.p2{v1}",
			expect: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}}},
		},
		{
			input: "schema1:p1.p2}",
			error: true,
		},
		{
			input: "schema1:schema2:p1",
			error: true,
		},
	}
	for idx, c := range cases {
		actual, err := ParseSwaggerPropertyAddr(c.input)
		if c.error {
			require.Error(t, err, idx)
			continue
		}
		require.NoError(t, err, idx)
		require.Equal(t, c.expect, actual, idx)
	}
}

func TestNewSwaggerPropertyAddr(t *testing.T) {
	cases := []struct {
		schemaName string
		propAddr   string
		expect     SwaggerPropertyAddr
		error      bool
	}{
		{
			schemaName: "",
			propAddr:   "",
			expect:     SwaggerPropertyAddr{},
		},
		{
			schemaName: "schema1",
			propAddr:   "p1.p2{v1}",
			expect: SwaggerPropertyAddr{
				Schema:       "schema1",
				PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}},
			},
		},
		{
			schemaName: "schema1",
			propAddr:   "p1.p2}",
			error:      true,
		},
	}
	for idx, c := range cases {
		actual, err := NewSwaggerPropertyAddr(c.schemaName, c.propAddr)
		if c.error {
			require.Error(t, err, idx)
			continue
		}
		require.NoError(t, err, idx)
		require.Equal(t, c.expect, actual, idx)
	}
}

func TestSwaggerPropertyAddr_Contains(t *testing.T) {
	cases := []struct {
		addr     SwaggerPropertyAddr
		oaddr    SwaggerPropertyAddr
		contains bool
	}{
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: nil},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			contains: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			contains: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
				{name: "p2"},
			}},
			contains: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
				{name: "p2"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			contains: false,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema2", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
				{name: "p2"},
			}},
			contains: false,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1", discriminatorValue: utils.String("v1")},
			}},
			contains: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1", discriminatorValue: utils.String("v1")},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1", discriminatorValue: utils.String("v1")},
			}},
			contains: false,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1", discriminatorValue: utils.String("v1")},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1", discriminatorValue: utils.String("v2")},
			}},
			contains: false,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1", discriminatorValue: utils.String("v2")},
				{name: "p2"},
			}},
			contains: true,
		},
	}
	for idx, c := range cases {
		require.Equal(t, c.contains, c.addr.Contains(c.oaddr), idx)
	}
}

func TestSwaggerPropertyAddr_Equals(t *testing.T) {
	cases := []struct {
		addr  SwaggerPropertyAddr
		oaddr SwaggerPropertyAddr
		equal bool
	}{
		{
			addr:  SwaggerPropertyAddr{},
			oaddr: SwaggerPropertyAddr{},
			equal: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			equal: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
				{name: "p2"},
			}},
			oaddr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
			}},
			equal: false,
		},
	}
	for idx, c := range cases {
		require.Equal(t, c.equal, c.addr.Equals(c.oaddr), idx)
	}
}

func TestSwaggerPropertyAddr_ToSwaggerDefinitionRef(t *testing.T) {
	cases := []struct {
		addr  SwaggerPropertyAddr
		ref   string
		error bool
	}{
		{
			addr:  SwaggerPropertyAddr{},
			error: true,
		},
		{
			addr:  SwaggerPropertyAddr{PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}}},
			error: true,
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
				{name: "p2"},
			}},
			ref: "#definitions/schema1/p1/p2",
		},
		{
			addr: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{
				{name: "p1"},
				{name: "p2", discriminatorValue: utils.String("v1")},
			}},
			ref: "#definitions/schema1/p1/v1",
		},
	}
	for idx, c := range cases {
		ref, err := c.addr.ToSwaggerDefinitionRef()
		if c.error {
			require.Error(t, err, idx)
			continue
		}
		require.NoError(t, err, idx)
		require.Equal(t, c.ref, ref.String(), idx)
	}
}

func TestSwaggerPropertyAddr_String(t *testing.T) {
	cases := []struct {
		addr   SwaggerPropertyAddr
		expect string
	}{
		{
			addr:   SwaggerPropertyAddr{},
			expect: "",
		},
		{
			addr:   SwaggerPropertyAddr{Schema: "schema1"},
			expect: "schema1:",
		},
		{
			addr:   SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
			expect: "schema1:p1.p2",
		},
		{
			addr:   SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}}},
			expect: "schema1:p1.p2{v1}",
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.expect, c.addr.String(), idx)
	}
}

func TestSwaggerPropertyAddr_MarshalJSON(t *testing.T) {
	cases := []struct {
		addr   SwaggerPropertyAddr
		expect string
	}{
		{
			addr:   SwaggerPropertyAddr{},
			expect: `""`,
		},
		{
			addr:   SwaggerPropertyAddr{Schema: "schema1"},
			expect: `"schema1:"`,
		},
		{
			addr:   SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
			expect: `"schema1:p1.p2"`,
		},
		{
			addr:   SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}}},
			expect: `"schema1:p1.p2{v1}"`,
		},
	}

	for idx, c := range cases {
		actual, err := c.addr.MarshalJSON()
		require.NoError(t, err)
		if c.expect == "" {
			require.Equal(t, c.expect, string(actual), idx)
		} else {
			require.JSONEq(t, c.expect, string(actual), idx)
		}
	}
}

func TestSwaggerPropertyAddr_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		expect SwaggerPropertyAddr
	}{
		{
			input:  `""`,
			expect: SwaggerPropertyAddr{},
		},
		{
			input:  `"schema1:"`,
			expect: SwaggerPropertyAddr{Schema: "schema1"},
		},
		{
			input:  `"schema1:p1.p2"`,
			expect: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2"}}},
		},
		{
			input:  `"schema1:p1.p2{v1}"`,
			expect: SwaggerPropertyAddr{Schema: "schema1", PropertyAddr: SwaggerRelPropertyAddr{{name: "p1"}, {name: "p2", discriminatorValue: utils.String("v1")}}},
		},
	}
	for idx, c := range cases {
		var addr SwaggerPropertyAddr
		err := addr.UnmarshalJSON([]byte(c.input))
		require.NoError(t, err, idx)
		require.Equal(t, c.expect, addr, idx)
	}
}
