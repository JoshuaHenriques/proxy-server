# Reverse Proxy

A reverse proxy operating at both Layer 4 (TCP/UDP) and Layer 7 (HTTP/HTTPS), designed to securely manage inbound client-to-server traffic across multiple protocols. 

> Work in Progress 🚧

---

## Core Features

- Inbound Traffic Handling: Front, protect, and forward both web and non-HTTP services (game servers, databases, VoIP, web applications) to backend services without modifying TCP/UDP or HTTP/HTTPS traffic
- NAT Traversal (Layer 4): Enable external access to internal services without public IPs through connection rewriting and port mapping

---

## Planned & WIP Features

- Load Balancing: Distribute traffic across multiple backend servers for high availability and performance
- Session Persistence: Maintain client connections to the same backend server for stateful applications
- SSL Termination: Handle HTTPS SSL/TLS encryption and decryption
- TLS/SSL Offloading: Terminate encrypted TCP connections at the proxy, reducing backend server load
- Source IP Preservation: Relay the original client IP to backend services for auditing and security
- Security Features: IP filtering and rate limiting
- Caching (Layer 7): Store and serve static or dynamic web content to reduce backend load and improve performance
