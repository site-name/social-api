// +build windows

package graceful

import "net"

// GetListener obtains a listener for the local network address.
// On windows this is basically just a shim around net.Listen.
func GetListener(network, address string) (net.Listener, error) {
	// Add a deferral to say that we've tried to grab a listener
	defer GetManager().InformCleanup()

	return net.Listen(network, address)
}
