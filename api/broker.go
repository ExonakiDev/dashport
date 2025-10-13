// api/broker.go
package api

type Broker interface {
	Authenticate() error
	RefreshToken() error
	// GetPortfolio() {interface{}, error}
}
