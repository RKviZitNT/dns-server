import socket

def build_dns_query(domain: str, recursive: bool):
    transaction_id = 0x0001  # Transaction ID
    flags = 0x0000  # Flags: Standard query
    if recursive:
        flags |= 0x0100  # Set RD bit for recursive query

    header = bytearray([
        (transaction_id >> 8) & 0xFF,  # High byte of Transaction ID
        transaction_id & 0xFF,         # Low byte of Transaction ID
        (flags >> 8) & 0xFF,           # High byte of Flags
        flags & 0xFF,                  # Low byte of Flags
        0x00, 0x01,                    # Questions
        0x00, 0x00,                    # Answer RRs
        0x00, 0x00,                    # Authority RRs
        0x00, 0x00                     # Additional RRs
    ])

    question = bytearray()
    labels = domain.split('.')
    for label in labels:
        question.append(len(label))
        question.extend(label.encode('utf-8'))
    question.extend([0x00, 0x00, 0x01, 0x00, 0x01])  # QTYPE and QCLASS (A record, Internet)

    return header + question

def send_dns_query(query, server, port=2053):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.sendto(query, (server, port))
        response, _ = s.recvfrom(1024)
    return response

if __name__ == "__main__":
    dns_server = "127.0.0.1"
    domain = input("domain: ")
    recursive = input("Recursive query? (y/n): ").lower() == 'y'
    dns_query = build_dns_query(domain, recursive)
    dns_response = send_dns_query(dns_query, dns_server)
    print(dns_response)