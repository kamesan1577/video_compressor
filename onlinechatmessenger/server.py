import socket

def relay(clients: list,data:bytes,sock:socket):
    for client in clients:
        sent = sock.sendto(data,client)
        print('sent {} bytes back to {}'.format(sent, client))
    return
# AF_INETを使用し、UDPソケットを作成
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

server_address = '0.0.0.0'
server_port = 9001
print('starting up on port {}'.format(server_port))

# ソケットを特殊なアドレス0.0.0.0とポート9001に紐付け
sock.bind((server_address, server_port))
clients = []
client_dict = {}
while True:
    print('\nwaiting to receive message')
    data, address = sock.recvfrom(4096)
    print('received {} bytes from {}'.format(len(data), address))
    print(data.decode())
    if not address in clients:
        clients.append(address)
        print(f"new address added: {address}")


    if data:
        relay(clients,data,sock)
