import dns.message, dns.rdatatype, socket

msg = dns.message.make_query('codecrafters.io', dns.rdatatype.A)
msg.question.append(
    dns.message.make_query('codecrafters.io', dns.rdatatype.A).question[0]
)

data = msg.to_wire()
print(f"Sending {len(msg.question)} questions: {[q.name.to_text() for q in msg.question]}")
print(f"Raw request ({len(data)} bytes): {data.hex()}")

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.settimeout(2)
sock.sendto(data, ('127.0.0.1', 2053))

try:
    response, _ = sock.recvfrom(512)
    print(f"Raw response ({len(response)} bytes): {response.hex()}")
    parsed = dns.message.from_wire(response)
    print(f"Response: {parsed}")
except socket.timeout:
    print("Timed out — is the server running?")
finally:
    sock.close()
