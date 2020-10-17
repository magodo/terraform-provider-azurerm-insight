package propertyaddr

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func TestNewPropertyAddrFromString(t *testing.T) {
	cases := []struct {
		input  string
		expect *PropertyAddr
	}{
		{
			"",
			&PropertyAddr{
				owner: "",
				addrs: []string{},
			},
		},
		{
			"p1",
			&PropertyAddr{
				owner: "",
				addrs: []string{"p1"},
			},
		},
		{
			"p1.p2",
			&PropertyAddr{
				owner: "",
				addrs: []string{"p1", "p2"},
			},
		},
		{
			"res1:",
			&PropertyAddr{
				owner: "res1",
				addrs: []string{},
			},
		},
		{
			"res1:p1.p2",
			&PropertyAddr{
				owner: "res1",
				addrs: []string{"p1", "p2"},
			},
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.expect, NewPropertyAddrFromString(c.input), idx)
	}
}

func TestNewPropertyAddr(t *testing.T) {
	require.Equal(t, &PropertyAddr{
		owner: "res1",
		addrs: []string{"p1", "p2"},
	}, NewPropertyAddrFromStringWithOwner("res1", "p1.p2"))
}

func TestNewPropertyAddrFromStringWithOwner(t *testing.T) {
	require.Equal(t, &PropertyAddr{
		owner: "res1",
		addrs: []string{"p1", "p2"},
	}, NewPropertyAddrFromStringWithOwner("res1", "p1.p2"))
}

func TestPropertyAddr_String(t *testing.T) {
	addrs := []string{
		"",
		"p1",
		"p1.p2",
		"res1:",
		"res1:p1",
		"res1:p1.p2",
	}

	for _, addr := range addrs {
		propAddr := NewPropertyAddrFromString(addr)
		require.Equal(t, addr, propAddr.String())
	}
}

func TestPropertyAddr_MarshalJSON(t *testing.T) {
	cases := []struct {
		addr   PropertyAddr
		expect string
	}{
		{
			addr:   PropertyAddr{},
			expect: `""`,
		},
		{
			addr:   PropertyAddr{addrs: []string{"p1"}},
			expect: `"p1"`,
		},
		{
			addr:   PropertyAddr{addrs: []string{"p1", "p2"}},
			expect: `"p1.p2"`,
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
			expect: `"res1:p1.p2"`,
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

func TestPropertyAddr_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		expect PropertyAddr
	}{
		{
			input:  `""`,
			expect: PropertyAddr{addrs: []string{}},
		},
		{
			input:  `"p1"`,
			expect: PropertyAddr{addrs: []string{"p1"}},
		},
		{
			input:  `"p1.p2"`,
			expect: PropertyAddr{addrs: []string{"p1", "p2"}},
		},
		{
			input:  `"res1:p1.p2"`,
			expect: PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
		},
	}

	for _, c := range cases {
		var addr PropertyAddr
		err := addr.UnmarshalJSON([]byte(c.input))
		require.NoError(t, err)
		require.Equal(t, c.expect, addr)
	}
}

func TestPropertyAddr_Append(t *testing.T) {
	cases := []struct {
		addr   PropertyAddr
		oaddr  string
		expect PropertyAddr
	}{
		{
			addr:   PropertyAddr{},
			oaddr:  "",
			expect: PropertyAddr{},
		},
		{
			addr:   PropertyAddr{},
			oaddr:  "p1",
			expect: PropertyAddr{addrs: []string{"p1"}},
		},
		{
			addr:   PropertyAddr{owner: "res1"},
			oaddr:  "",
			expect: PropertyAddr{owner: "res1"},
		},
		{
			addr:   PropertyAddr{owner: "res1"},
			oaddr:  "p1",
			expect: PropertyAddr{owner: "res1", addrs: []string{"p1"}},
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p2"}},
			oaddr:  "p1",
			expect: PropertyAddr{owner: "res1", addrs: []string{"p2", "p1"}},
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.expect, c.addr.Append(c.oaddr), idx)
	}
}

func TestPropertyAddr_Contains(t *testing.T) {
	cases := []struct {
		addr     PropertyAddr
		oaddr    PropertyAddr
		contains bool
	}{
		{
			addr:     PropertyAddr{},
			oaddr:    PropertyAddr{},
			contains: false,
		},
		{
			addr:     PropertyAddr{owner: "res1"},
			oaddr:    PropertyAddr{},
			contains: false,
		},
		{
			addr:     PropertyAddr{addrs: []string{"p1"}},
			oaddr:    PropertyAddr{},
			contains: false,
		},
		{
			addr:     PropertyAddr{},
			oaddr:    PropertyAddr{owner: "res1"},
			contains: false,
		},
		{
			addr:     PropertyAddr{},
			oaddr:    PropertyAddr{addrs: []string{"p1"}},
			contains: true,
		},
		{
			addr:     PropertyAddr{owner: "res1", addrs: []string{}},
			oaddr:    PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			contains: true,
		},
		{
			addr:     PropertyAddr{addrs: []string{"p1"}},
			oaddr:    PropertyAddr{addrs: []string{"p2"}},
			contains: false,
		},
		{
			addr:     PropertyAddr{addrs: []string{"p1"}},
			oaddr:    PropertyAddr{addrs: []string{"p1"}},
			contains: false,
		},
		{
			addr:     PropertyAddr{addrs: []string{"p1"}},
			oaddr:    PropertyAddr{addrs: []string{"p1", "p2"}},
			contains: true,
		},
		{
			addr:     PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:    PropertyAddr{owner: "res1", addrs: []string{"p2"}},
			contains: false,
		},
		{
			addr:     PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:    PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			contains: false,
		},
		{
			addr:     PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:    PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
			contains: true,
		},
		{
			addr:     PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:    PropertyAddr{owner: "res2", addrs: []string{"p1"}},
			contains: false,
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.contains, c.addr.Contains(c.oaddr), idx)
	}
}

func TestPropertyAddr_Equals(t *testing.T) {
	cases := []struct {
		addr   PropertyAddr
		oaddr  PropertyAddr
		equals bool
	}{
		{
			addr:   PropertyAddr{},
			oaddr:  PropertyAddr{},
			equals: true,
		},
		{
			addr:   PropertyAddr{owner: "res1"},
			oaddr:  PropertyAddr{},
			equals: false,
		},
		{
			addr:   PropertyAddr{addrs: []string{"p1"}},
			oaddr:  PropertyAddr{},
			equals: false,
		},
		{
			addr:   PropertyAddr{},
			oaddr:  PropertyAddr{owner: "res1"},
			equals: false,
		},
		{
			addr:   PropertyAddr{owner: "res1"},
			oaddr:  PropertyAddr{owner: "res1"},
			equals: true,
		},
		{
			addr:   PropertyAddr{},
			oaddr:  PropertyAddr{addrs: []string{"p1"}},
			equals: false,
		},
		{
			addr:   PropertyAddr{addrs: []string{"p1"}},
			oaddr:  PropertyAddr{addrs: []string{"p2"}},
			equals: false,
		},
		{
			addr:   PropertyAddr{addrs: []string{"p1"}},
			oaddr:  PropertyAddr{addrs: []string{"p1"}},
			equals: true,
		},
		{
			addr:   PropertyAddr{addrs: []string{"p1"}},
			oaddr:  PropertyAddr{addrs: []string{"p1", "p2"}},
			equals: false,
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:  PropertyAddr{owner: "res1", addrs: []string{"p2"}},
			equals: false,
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:  PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			equals: true,
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:  PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
			equals: false,
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:  PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
			equals: false,
		},
		{
			addr:   PropertyAddr{owner: "res1", addrs: []string{"p1"}},
			oaddr:  PropertyAddr{owner: "res2", addrs: []string{"p1"}},
			equals: false,
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.equals, c.addr.Equals(c.oaddr), idx)
	}
}

func TestToSwaggerDefinitionRef(t *testing.T) {
	cases := []struct {
		addr    PropertyAddr
		uri     string
		isError bool
	}{
		{
			addr:    PropertyAddr{},
			uri:     "#definitions",
			isError: false,
		},
		{
			addr:    PropertyAddr{addrs: []string{"p1", "p2"}},
			uri:     "#definitions/p1/p2",
			isError: false,
		},
		{
			addr:    PropertyAddr{owner: "schema1"},
			uri:     "#definitions/schema1",
			isError: false,
		},
		{
			addr:    PropertyAddr{owner: "schema1", addrs: []string{"p1", "p2"}},
			uri:     "#definitions/schema1/p1/p2",
			isError: false,
		},
	}

	for idx, c := range cases {
		ref, err := ToSwaggerDefinitionRef(c.addr)
		if c.isError {
			require.Error(t, err, idx)
			continue
		}
		expect, _ := spec.NewRef(c.uri)
		require.Equal(t, expect.String(), ref.String(), idx)
	}
}
