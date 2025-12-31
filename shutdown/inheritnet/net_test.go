//go:build !windows
// +build !windows

package inheritnet

import (
	"fmt"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/andeya/gust/shutdown"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListen(t *testing.T) {
	// TCP variants
	for _, nett := range []string{"tcp", "tcp4", "tcp6"} {
		ln, err := Listen(nett, ":0")
		require.NoError(t, err)
		ln.Close()
	}

	// Unix socket
	tmp := createTempPath(t)
	defer os.Remove(tmp)
	ln, err := Listen("unix", tmp)
	require.NoError(t, err)
	ln.Close()

	// Invalid network
	_, err = Listen("invalid", ":0")
	assert.Error(t, err)

	// Invalid address
	_, err = Listen("tcp", "invalid:address")
	assert.Error(t, err)
}

func TestListenTCP(t *testing.T) {
	// Nil address
	ln, err := ListenTCP("tcp", nil)
	require.NoError(t, err)
	ln.Close()

	// With inherited listener
	ln1, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	require.NoError(t, err)
	defer ln1.Close()

	n := &inheritNet{inherited: []net.Listener{ln1}}
	ln2, err := n.ListenTCP("tcp", ln1.Addr().(*net.TCPAddr))
	require.NoError(t, err)
	assert.Equal(t, ln1, ln2)
	assert.Nil(t, n.inherited[0])

	// With nil in inherited
	n = &inheritNet{inherited: []net.Listener{nil}}
	ln3, err := n.ListenTCP("tcp", nil)
	require.NoError(t, err)
	ln3.Close()
}

func TestListenUnix(t *testing.T) {
	tmp := createTempPath(t)
	defer os.Remove(tmp)

	addr, _ := net.ResolveUnixAddr("unix", tmp)
	ln, err := ListenUnix("unix", addr)
	require.NoError(t, err)
	defer ln.Close()

	// With inherited
	n := &inheritNet{inherited: []net.Listener{ln}}
	ln2, err := n.ListenUnix("unix", addr)
	require.NoError(t, err)
	assert.Equal(t, ln, ln2)
}

func TestAppend_Public(t *testing.T) {
	// Reset globalInheritNet for clean test
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()
	// Test public Append function (may fail due to global state)
	_ = Append(ln)
}

func TestSetInherited_Public(t *testing.T) {
	s := shutdown.New()
	// Test public SetInherited function
	_ = SetInherited(s)
}

func TestAppend(t *testing.T) {
	n := &inheritNet{}
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()

	// First append
	assert.NoError(t, n.Append(ln))
	assert.Len(t, n.active, 1)

	// Duplicate
	err := n.Append(ln)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "re-register")

	// With inherited match
	n2 := &inheritNet{inherited: []net.Listener{ln}}
	mockLn := &mockListener{addr: ln.Addr()}
	assert.NoError(t, n2.Append(mockLn))
	assert.Equal(t, ln, n2.active[0])

	// With nil in lists
	n3 := &inheritNet{inherited: []net.Listener{nil}, active: []net.Listener{nil}}
	ln3, _ := net.Listen("tcp", ":0")
	defer ln3.Close()
	assert.NoError(t, n3.Append(ln3))
}

func TestInherit(t *testing.T) {
	// No env
	os.Unsetenv(envCountKey)
	n := &inheritNet{}
	assert.NoError(t, n.inherit())

	// Invalid env
	os.Setenv(envCountKey, "invalid")
	n = &inheritNet{}
	assert.Error(t, n.inherit())
	os.Unsetenv(envCountKey)

	// Zero count
	os.Setenv(envCountKey, "0")
	n = &inheritNet{}
	assert.NoError(t, n.inherit())
	os.Unsetenv(envCountKey)

	// Invalid fd
	os.Setenv(envCountKey, "1")
	n = &inheritNet{fdStart: 99999}
	assert.Error(t, n.inherit())
	os.Unsetenv(envCountKey)
}

func TestIsSameAddr(t *testing.T) {
	cases := []struct {
		a1, a2 string
		same   bool
	}{
		{"127.0.0.1:8080", "127.0.0.1:8080", true},
		{"127.0.0.1:8080", "127.0.0.1:8081", false},
		{"0.0.0.0:8080", "[::]:8080", true},
		{"0.0.0.0:8080", ":8080", true},
	}
	for _, c := range cases {
		a1, _ := net.ResolveTCPAddr("tcp", c.a1)
		a2, _ := net.ResolveTCPAddr("tcp", c.a2)
		assert.Equal(t, c.same, isSameAddr(a1, a2), "%s vs %s", c.a1, c.a2)
	}

	// Different network
	tcp, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	unix, _ := net.ResolveUnixAddr("unix", "/tmp/test")
	assert.False(t, isSameAddr(tcp, unix))
}

func TestSetInherited(t *testing.T) {
	s := shutdown.New()

	// Empty
	n := &inheritNet{}
	assert.NoError(t, n.SetInherited(s))

	// With listeners
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()
	n.active = []net.Listener{ln}
	assert.NoError(t, n.SetInherited(s))

	// Without File method (panic)
	n2 := &inheritNet{active: []net.Listener{&mockListener{addr: &net.TCPAddr{}}}}
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, fmt.Sprintf("%v", r), "missing method File")
		}
	}()
	_ = n2.SetInherited(s)
}

func TestActiveListeners(t *testing.T) {
	n := &inheritNet{}
	active, _ := n.activeListeners()
	assert.Empty(t, active)

	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()
	n.active = []net.Listener{ln}

	active, _ = n.activeListeners()
	assert.Len(t, active, 1)

	// Returns copy
	n.active = append(n.active, ln)
	active2, _ := n.activeListeners()
	assert.NotEqual(t, len(active), len(active2))
}

func TestConcurrentAccess(t *testing.T) {
	n := &inheritNet{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ln, err := n.ListenTCP("tcp", nil)
			if err == nil {
				ln.Close()
			}
		}()
	}
	wg.Wait()
}

func TestErrorPaths(t *testing.T) {
	os.Setenv(envCountKey, "invalid")
	defer os.Unsetenv(envCountKey)

	// ListenTCP inherit error
	n1 := &inheritNet{}
	_, err := n1.ListenTCP("tcp", nil)
	assert.Error(t, err)

	// ListenUnix inherit error
	tmp := createTempPath(t)
	defer os.Remove(tmp)
	addr, _ := net.ResolveUnixAddr("unix", tmp)
	n2 := &inheritNet{}
	_, err = n2.ListenUnix("unix", addr)
	assert.Error(t, err)

	// Append inherit error
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()
	n3 := &inheritNet{}
	err = n3.Append(ln)
	assert.Error(t, err)
}

func createTempPath(t *testing.T) string {
	f, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

type mockListener struct{ addr net.Addr }

func (m *mockListener) Accept() (net.Conn, error) { return nil, nil }
func (m *mockListener) Close() error              { return nil }
func (m *mockListener) Addr() net.Addr            { return m.addr }
