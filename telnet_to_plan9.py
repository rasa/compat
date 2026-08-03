#!/usr/bin/env python3
"""telnet_to_plan9.py."""

import os
import socket
import sys
import time

HOST = "127.0.0.1"
PORT = int(os.environ.get("SSH_PORT", "10023"))
TIMEOUT = 30 * 60

IAC = 255
SE = 240
SB = 250
WILL = 251
WONT = 252
DO = 253
DONT = 254

PASS = b"__PLAN9_PASS__"
FAIL = b"__PLAN9_FAIL__"


def connect():
    """Doc me."""
    last_error = None

    for _ in range(20):
        try:
            return socket.create_connection(
                (HOST, PORT),
                timeout=15,
            )
        except OSError as error:
            last_error = error
            time.sleep(3)

    raise RuntimeError(f"Could not connect to Plan 9: {last_error}")


sock = connect()
sock.settimeout(0.5)

pending = bytearray()
output = bytearray()


def process_telnet(data, _pending):
    """Doc me."""
    data = _pending + data
    _pending = bytearray()

    visible = bytearray()
    i = 0

    while i < len(data):
        if data[i] != IAC:
            visible.append(data[i])
            i += 1
            continue

        if i + 1 >= len(data):
            _pending.extend(data[i:])
            break

        cmd = data[i + 1]

        if cmd == IAC:
            visible.append(IAC)
            i += 2
            continue

        if cmd in (DO, DONT, WILL, WONT):
            if i + 2 >= len(data):
                _pending.extend(data[i:])
                break

            option = data[i + 2]

            if cmd == DO:
                sock.sendall(bytes((IAC, WONT, option)))
            elif cmd == WILL:
                sock.sendall(bytes((IAC, DONT, option)))

            i += 3
            continue

        if cmd == SB:
            end = data.find(bytes((IAC, SE)), i + 2)

            if end == -1:
                _pending.extend(data[i:])
                break

            i = end + 2
            continue

        i += 2

    return bytes(visible), _pending


def read_for(seconds, _pending):
    """Doc me."""
    end = time.monotonic() + seconds

    while time.monotonic() < end:
        try:
            data = sock.recv(4096)
        except socket.timeout:
            continue

        if not data:
            return False, _pending

        visible, _pending = process_telnet(data, _pending)

        if visible:
            sys.stdout.buffer.write(visible)
            sys.stdout.buffer.flush()
            output.extend(visible)

        if PASS in output or FAIL in output:
            return True, _pending

    return False, _pending


# Consume the initial Telnet negotiation/banner.
_, pending = read_for(2, pending)

sock.sendall(b"cd /usr/glenda/work\r\n")
_, pending = read_for(1, pending)

TEST_CMD = (
    "./compat.test "
    "-test.count 1 "
    "-test.timeout 20m "
    "-test.v "
    "-test.coverprofile coverage.out; "
    "teststatus=$status; "
    "if(~ $teststatus '') echo __PLAN9_^PASS__; "
    "if(! ~ $teststatus '') echo __PLAN9_^FAIL__"
    "\r\n"
)

sock.sendall(TEST_CMD.encode())

deadline = time.monotonic() + TIMEOUT

while time.monotonic() < deadline:
    rv, pending = read_for(1, pending)
    if rv:
        break
else:
    print(
        "\nTimed out running Plan 9 tests",
        file=sys.stderr,
    )
    raise SystemExit(124)

passed = PASS in output

try:
    sock.sendall(b"fshalt\r\n")
    _, pending = read_for(3, pending)
except OSError:
    pass
finally:
    sock.close()

raise SystemExit(0 if passed else 1)
