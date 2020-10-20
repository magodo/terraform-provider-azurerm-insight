package propertyaddr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPropertyAddrFromString(t *testing.T) {
	cases := []struct {
		input  string
		expect *TerraformPropertyAddr
	}{
		{
			"",
			&TerraformPropertyAddr{
				ResourceName: "",
				PropertyAddr: []string{},
			},
		},
		{
			"p1",
			&TerraformPropertyAddr{
				ResourceName: "",
				PropertyAddr: []string{"p1"},
			},
		},
		{
			"p1.p2",
			&TerraformPropertyAddr{
				ResourceName: "",
				PropertyAddr: []string{"p1", "p2"},
			},
		},
		{
			"res1:",
			&TerraformPropertyAddr{
				ResourceName: "res1",
				PropertyAddr: []string{},
			},
		},
		{
			"res1:p1.p2",
			&TerraformPropertyAddr{
				ResourceName: "res1",
				PropertyAddr: []string{"p1", "p2"},
			},
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.expect, ParseTerraformPropertyAddr(c.input), idx)
	}
}

func TestNewPropertyAddr(t *testing.T) {
	require.Equal(t, &TerraformPropertyAddr{
		ResourceName: "res1",
		PropertyAddr: []string{"p1", "p2"},
	}, NewTerraformPropertyAddr("res1", "p1.p2"))
}

func TestNewPropertyAddrFromStringWithOwner(t *testing.T) {
	require.Equal(t, &TerraformPropertyAddr{
		ResourceName: "res1",
		PropertyAddr: []string{"p1", "p2"},
	}, NewTerraformPropertyAddr("res1", "p1.p2"))
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
		propAddr := ParseTerraformPropertyAddr(addr)
		require.Equal(t, addr, propAddr.String())
	}
}

func TestPropertyAddr_MarshalJSON(t *testing.T) {
	cases := []struct {
		addr   TerraformPropertyAddr
		expect string
	}{
		{
			addr:   TerraformPropertyAddr{},
			expect: `""`,
		},
		{
			addr:   TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			expect: `"p1"`,
		},
		{
			addr:   TerraformPropertyAddr{PropertyAddr: []string{"p1", "p2"}},
			expect: `"p1.p2"`,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1", "p2"}},
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
		expect TerraformPropertyAddr
	}{
		{
			input:  `""`,
			expect: TerraformPropertyAddr{PropertyAddr: []string{}},
		},
		{
			input:  `"p1"`,
			expect: TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
		},
		{
			input:  `"p1.p2"`,
			expect: TerraformPropertyAddr{PropertyAddr: []string{"p1", "p2"}},
		},
		{
			input:  `"res1:p1.p2"`,
			expect: TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1", "p2"}},
		},
	}

	for _, c := range cases {
		var addr TerraformPropertyAddr
		err := addr.UnmarshalJSON([]byte(c.input))
		require.NoError(t, err)
		require.Equal(t, c.expect, addr)
	}
}

func TestPropertyAddr_Append(t *testing.T) {
	cases := []struct {
		addr   TerraformPropertyAddr
		oaddr  string
		expect TerraformPropertyAddr
	}{
		{
			addr:   TerraformPropertyAddr{},
			oaddr:  "",
			expect: TerraformPropertyAddr{},
		},
		{
			addr:   TerraformPropertyAddr{},
			oaddr:  "p1",
			expect: TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1"},
			oaddr:  "",
			expect: TerraformPropertyAddr{ResourceName: "res1"},
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1"},
			oaddr:  "p1",
			expect: TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p2"}},
			oaddr:  "p1",
			expect: TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p2", "p1"}},
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.expect, c.addr.Append(c.oaddr), idx)
	}
}

func TestPropertyAddr_Contains(t *testing.T) {
	cases := []struct {
		addr     TerraformPropertyAddr
		oaddr    TerraformPropertyAddr
		contains bool
	}{
		{
			addr:     TerraformPropertyAddr{},
			oaddr:    TerraformPropertyAddr{},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{ResourceName: "res1"},
			oaddr:    TerraformPropertyAddr{},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{},
			oaddr:    TerraformPropertyAddr{ResourceName: "res1"},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{},
			oaddr:    TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			contains: true,
		},
		{
			addr:     TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{}},
			oaddr:    TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			contains: true,
		},
		{
			addr:     TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{PropertyAddr: []string{"p2"}},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{PropertyAddr: []string{"p1", "p2"}},
			contains: true,
		},
		{
			addr:     TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p2"}},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			contains: false,
		},
		{
			addr:     TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1", "p2"}},
			contains: true,
		},
		{
			addr:     TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:    TerraformPropertyAddr{ResourceName: "res2", PropertyAddr: []string{"p1"}},
			contains: false,
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.contains, c.addr.Contains(c.oaddr), idx)
	}
}

func TestPropertyAddr_Equals(t *testing.T) {
	cases := []struct {
		addr   TerraformPropertyAddr
		oaddr  TerraformPropertyAddr
		equals bool
	}{
		{
			addr:   TerraformPropertyAddr{},
			oaddr:  TerraformPropertyAddr{},
			equals: true,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1"},
			oaddr:  TerraformPropertyAddr{},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{},
			oaddr:  TerraformPropertyAddr{ResourceName: "res1"},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1"},
			oaddr:  TerraformPropertyAddr{ResourceName: "res1"},
			equals: true,
		},
		{
			addr:   TerraformPropertyAddr{},
			oaddr:  TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{PropertyAddr: []string{"p2"}},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			equals: true,
		},
		{
			addr:   TerraformPropertyAddr{PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{PropertyAddr: []string{"p1", "p2"}},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p2"}},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			equals: true,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1", "p2"}},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1", "p2"}},
			equals: false,
		},
		{
			addr:   TerraformPropertyAddr{ResourceName: "res1", PropertyAddr: []string{"p1"}},
			oaddr:  TerraformPropertyAddr{ResourceName: "res2", PropertyAddr: []string{"p1"}},
			equals: false,
		},
	}

	for idx, c := range cases {
		require.Equal(t, c.equals, c.addr.Equals(c.oaddr), idx)
	}
}
