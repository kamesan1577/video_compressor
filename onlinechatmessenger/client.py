import socket

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

server_address = "0.0.0.0"
server_port = 9001

address = ''
port = 9050

username = bytes(input("Type your name: "),"utf-8")
message = b'Message to send to the client.'

# 空の文字列も0.0.0.0として使用できます。
sock.bind((address,port))

try:
    print('sending {!r}'.format(message))
    # サーバへのデータ送信
    sent = sock.sendto(len(username).to_bytes(1,"big")+username+message, (server_address, server_port))
    print('Send {} bytes'.format(sent))
    print(f"usernamelen: {len(username)}")

    # 応答を受信
    print('waiting to receive')
    data, server = sock.recvfrom(4096)
    print('received {!r}'.format(data.decode()))
finally:
    print('closing socket')
    sock.close()
