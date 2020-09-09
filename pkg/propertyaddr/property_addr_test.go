package propertyaddr

import (
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"testing"
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
		assert.Equal(t, c.expect, NewPropertyAddrFromString(c.input), idx)
	}
}

func TestNewPropertyAddr(t *testing.T) {
	assert.Equal(t, &PropertyAddr{
		owner: "res1",
		addrs: []string{"p1", "p2"},
	}, NewPropertyAddr("res1", "p1", "p2"))
}

func TestNewPropertyAddrFromStringWithOwner(t *testing.T) {
	assert.Equal(t, &PropertyAddr{
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
		assert.Equal(t, addr, propAddr.String())
	}
}

func TestPropertyAddr_MarshalJSON(t *testing.T) {
	cases := []struct {
		addr   PropertyAddr
		expect string
	}{
		{
			addr:   PropertyAddr{},
			expect: "{}",
		},
		{
			addr: PropertyAddr{addrs: []string{"p1"}},
			expect: `{
  "addr": "p1"
}`,
		},
		{
			addr: PropertyAddr{addrs: []string{"p1", "p2"}},
			expect: `{
  "addr": "p1.p2"
}`,
		},
		{
			addr: PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
			expect: `{
  "owner": "res1",
  "addr": "p1.p2"
}`,
		},
	}

	for _, c := range cases {
		actual, err := c.addr.MarshalJSON()
		assert.NoError(t, err)
		assert.JSONEq(t, c.expect, string(actual))
	}
}

func TestPropertyAddr_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		expect PropertyAddr
	}{
		{
			input:  "{}",
			expect: PropertyAddr{},
		},
		{
			input: `{
  "addr": "p1"
}`,
			expect: PropertyAddr{addrs: []string{"p1"}},
		},
		{
			input: `{
  "addr": "p1.p2"
}`,
			expect: PropertyAddr{addrs: []string{"p1", "p2"}},
		},
		{
			input: `{
  "owner": "res1",
  "addr": "p1.p2"
}`,
			expect: PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
		},
	}

	for _, c := range cases {
		var addr PropertyAddr
		err := addr.UnmarshalJSON([]byte(c.input))
		assert.NoError(t, err)
		assert.Equal(t, c.expect, addr)
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
		assert.Equal(t, c.expect, c.addr.Append(c.oaddr), idx)
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
			contains: false,
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
			oaddr:    PropertyAddr{owner: "res1", addrs: []string{"p1", "p2"}},
			contains: true,
		},
	}

	for idx, c := range cases {
		assert.Equal(t, c.contains, c.addr.Contains(c.oaddr), idx)
	}
}

func TestToSwaggerDefinitionRef(t *testing.T) {
	cases := []struct{
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

	for idx, c :=range cases {
		ref, err := ToSwaggerDefinitionRef(c.addr)
		if c.isError {
			assert.Error(t, err, idx)
			continue
		}
		expect, _ := spec.NewRef(c.uri)
		assert.Equal(t, expect.String(), ref.String(), idx)
	}
}
