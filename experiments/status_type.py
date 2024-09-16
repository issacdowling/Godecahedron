import socketserver
import struct

HOST = "0.0.0.0"
PORT = 1234


class StatusType(socketserver.BaseRequestHandler):
    # This should just check whether a client is trying to get server status or logging in (but not respond)
    def handle(self):
        self.data = self.request.recv(1024).strip()

        print(f"{self.client_address[0]} {(self.data)}")
        # print(self.data[-1])
        if self.data[-1] == 2:
            print("Login")
        elif self.data[-1] == 1:
            print("Status")


if __name__ == "__main__":
    with socketserver.TCPServer((HOST, PORT), StatusType) as server:
        server.serve_forever()
