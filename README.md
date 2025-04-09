# **Proxy Server**

A multi-protocol proxy server for seamless TCP/UDP and HTTP/HTTPS traffic routing, offering **bidirectional proxy functionality** for both Layer 4 (TCP/UDP) and Layer 7 (HTTP/HTTPS) traffic in a lightweight solution.

_(Work in Progress)_

## **Core Features**

- **Dual-Protocol Support**:
  - **Layer 4 (TCP/UDP)**: Stream raw packets for **outbound** (client → external server) or **inbound** (client → internal service) traffic flows, ideal for non-HTTP protocols like game servers, databases, or VoIP.
  - **Layer 7 (HTTP/HTTPS)**: Route web traffic for **outbound** (masking client IPs) or **inbound** (load balancing, SSL termination) use cases, supporting APIs, websites, and applications.

---
