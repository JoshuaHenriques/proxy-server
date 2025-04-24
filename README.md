# Reverse Proxy

## Core Functionality  
A high-performance reverse proxy operating at transport (Layer 4) and application (Layer 7) layers, designed to securely manage inbound client-to-server traffic across multiple protocols with optional forward proxy capabilities.

---

## Core Features  
✅ **Layer 4 (TCP/UDP):**  
- **Inbound Traffic**: Front and protect non-HTTP services (game servers, databases, VoIP)  
- **Raw Packet Handling**: Stream unmodified TCP/UDP traffic to backend services  
- **NAT Traversal**: Enable external access to internal services without public IPs through connection rewriting and port mapping  

✅ **Layer 7 (HTTP/HTTPS):**  
- **Web Traffic Management**: Load balancing, path-based routing, and SSL termination  
- **Security Features**: IP filtering, rate limiting
