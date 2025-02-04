// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interop

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"

	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

func (s *Server) verifySenderCertificate(ctx context.Context, senderID string, state *tls.ConnectionState) (addrs []string, err error) {
	// TODO: Support reading TLS client certificate from proxy headers (https://github.com/TheThingsNetwork/lorawan-stack/issues/717).
	senderClientCAs, err := s.SenderClientCAs(ctx, senderID)
	if err != nil {
		return nil, err
	}
	for _, chain := range state.VerifiedChains {
		peerCert, clientCA := chain[0], chain[len(chain)-1]
		for _, senderClientCA := range senderClientCAs {
			if clientCA.Equal(senderClientCA) {
				// If the TLS client certificate contains DNS addresses, use those. Otherwise, fallback to using CommonName as address.
				if len(peerCert.DNSNames) > 0 {
					addrs = append([]string(nil), peerCert.DNSNames...)
				} else {
					addrs = []string{peerCert.Subject.CommonName}
				}
				return
			}
		}
	}
	// TODO: Verify state.PeerCertificates[0] with senderClientCAs as Roots and state.PeerCertificates[1:] as Intermediates.
	// (https://github.com/TheThingsNetwork/lorawan-stack/issues/718).
	return nil, errUnauthenticated.New()
}

func verifySenderNSID(patterns []string, nsID string) error {
	if len(patterns) == 0 {
		return errCallerNotAuthorized.WithAttributes("target", nsID)
	}

	host := nsID
	if url, err := url.Parse(nsID); err == nil && url.Host != "" {
		host = url.Host
	}
	if h, _, err := net.SplitHostPort(nsID); err == nil {
		host = h
	}
	if len(host) == 0 {
		return errCallerNotAuthorized.WithAttributes("target", nsID)
	}
	hostParts := strings.Split(host, ".")

nextPattern:
	for _, pattern := range patterns {
		patternParts := strings.Split(pattern, ".")
		if len(patternParts) != len(hostParts) {
			return errCallerNotAuthorized.WithAttributes("target", nsID)
		}
		for i, patternPart := range patternParts {
			if i == 0 && patternPart == "*" {
				continue
			}
			if patternPart != hostParts[i] {
				continue nextPattern
			}
		}
		return nil
	}
	return errCallerNotAuthorized.WithAttributes("target", nsID)
}

type authInfo interface {
	GetAddresses() []string
}

// NetworkServerAuthInfo contains authentication information of the Network Server.
type NetworkServerAuthInfo struct {
	NetID     types.NetID
	Addresses []string
}

func (n NetworkServerAuthInfo) GetAddresses() []string { return n.Addresses }

type nsAuthInfoKeyType struct{}

var nsAuthInfoKey nsAuthInfoKeyType

// NewContextWithNetworkServerAuthInfo returns a derived context with the given authentication information of the
// Network Server.
func NewContextWithNetworkServerAuthInfo(parent context.Context, authInfo NetworkServerAuthInfo) context.Context {
	return context.WithValue(parent, nsAuthInfoKey, authInfo)
}

// NetworkServerAuthInfoFromContext returns the authentication information of the Network Server from context.
func NetworkServerAuthInfoFromContext(ctx context.Context) (NetworkServerAuthInfo, bool) {
	authInfo, ok := ctx.Value(nsAuthInfoKey).(NetworkServerAuthInfo)
	return authInfo, ok
}

func (s *Server) authenticateNS(ctx context.Context, r *http.Request, data []byte) (context.Context, error) {
	var header NsMessageHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, err
	}
	if !header.ProtocolVersion.SupportsNSID() && header.SenderNSID != nil {
		return nil, ErrMalformedMessage.New()
	}

	// If the client presents a TLS client certificate, use that for authentication.
	if state := r.TLS; state != nil && len(state.PeerCertificates) > 0 {
		addrs, err := s.verifySenderCertificate(ctx, types.NetID(header.SenderID).String(), state)
		if err != nil {
			return nil, err
		}
		if header.SenderNSID != nil {
			if err := verifySenderNSID(addrs, *header.SenderNSID); err != nil {
				return nil, err
			}
		}
		return NewContextWithNetworkServerAuthInfo(ctx, NetworkServerAuthInfo{
			NetID:     types.NetID(header.SenderID),
			Addresses: addrs,
		}), nil
	}

	// TODO: Verify Packet Broker token (https://github.com/TheThingsNetwork/lorawan-stack/issues/4703).

	return ctx, nil
}

type ApplicationServerAuthInfo struct {
	ASID      string
	Addresses []string
}

func (a ApplicationServerAuthInfo) GetAddresses() []string { return a.Addresses }

type asAuthInfoKeyType struct{}

var asAuthInfoKey asAuthInfoKeyType

// NewContextWithApplicationServerAuthInfo returns a derived context with the given authentication information of the
// Application Server.
func NewContextWithApplicationServerAuthInfo(parent context.Context, authInfo ApplicationServerAuthInfo) context.Context {
	return context.WithValue(parent, asAuthInfoKey, authInfo)
}

// ApplicationServerAuthInfoFromContext returns the authentication information of the Application Server from context.
func ApplicationServerAuthInfoFromContext(ctx context.Context) (ApplicationServerAuthInfo, bool) {
	authInfo, ok := ctx.Value(asAuthInfoKey).(ApplicationServerAuthInfo)
	return authInfo, ok
}

func (s *Server) authenticateAS(ctx context.Context, r *http.Request, data []byte) (context.Context, error) {
	var header AsMessageHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, err
	}

	state := r.TLS
	if state == nil {
		return nil, errUnauthenticated.New()
	}
	addrs, err := s.verifySenderCertificate(ctx, header.SenderID, state)
	if err != nil {
		return nil, err
	}
	return NewContextWithApplicationServerAuthInfo(ctx, ApplicationServerAuthInfo{
		ASID:      header.SenderID,
		Addresses: addrs,
	}), nil
}

type senderAuthenticatorFunc func(ctx context.Context, r *http.Request, data []byte) (context.Context, error)

func (f senderAuthenticatorFunc) Authenticate(ctx context.Context, r *http.Request, data []byte) (context.Context, error) {
	return f(ctx, r, data)
}

type senderAuthenticator interface {
	Authenticate(ctx context.Context, r *http.Request, data []byte) (context.Context, error)
}
