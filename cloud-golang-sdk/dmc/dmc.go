package dmc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"reflect"
	"time"
)

const (
	lrpc2Version             uint32 = 4          // LRPC version for handshake
	handShakeAck             uint32 = 1          // handshare acknowledged
	msgIDAdminClientGetToken uint16 = 1500       // Command
	msgDataTypeString        byte   = 4          // String data type
	msgEnd                   byte   = 0          // End of message
	magicNumber              uint32 = 0xABCD8012 // magic number
	headerLength             uint16 = 34         // Header length
	msgDataTypeInt32         byte   = 2          // Int32 data type
	msgDataTypeSet           byte   = 7          // Set data type
)

// LRPC2 represents local RPC data structure
type LRPC2 struct {
	unixSocketFile   string
	maxPayloadLength uint32
	pid              int
	client           net.Conn
}

// NewLRPC2 initiates a new local RPC client
func NewLRPC2() *LRPC2 {
	lrpc := LRPC2{}
	lrpc.unixSocketFile = "/var/centrify/cloud/daemon2"

	return &lrpc
}

func reader(r io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := r.Read(buf[:])
	if err != nil {
		return nil, fmt.Errorf("Error reading from server: %v", err)
	}

	return buf, nil
}

func (lrpc *LRPC2) connect() error {
	c, err := net.Dial("unix", lrpc.unixSocketFile)
	if err != nil {
		fmt.Printf("Error connecting to local rpc server: %v", err)
		return fmt.Errorf("Error connecting to local rpc server: %v", err)
	}
	lrpc.client = c
	//defer c.Close()

	lrpc.pid = os.Getegid()
	// Do handshake to check version number
	_, err = lrpc.client.Write(uint32ToByteArray(lrpc2Version))
	if err != nil {
		return fmt.Errorf("Error in handshake: %v", err)
	}

	// Read response for version
	data, err := reader(lrpc.client, 4)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	answer := byteArrayToUInt32(data)
	if answer != handShakeAck {
		return fmt.Errorf("Server doesn't support LRPC2 version 4")
	}
	data, err = reader(lrpc.client, 4)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	lrpc.maxPayloadLength = byteArrayToUInt32(data)

	//fmt.Printf("Successfully connected to Centrify Client\n")

	return nil
}

func (lrpc *LRPC2) request(payload []byte) ([]interface{}, error) {
	// Construct request header
	var data []byte
	data = append(data, uint32ToByteArray(magicNumber)...)      // magic number
	data = append(data, uint16ToByteArray(headerLength)...)     // header length
	data = append(data, uint32ToByteArray(lrpc2Version)...)     // LRPC2 version
	data = append(data, uint64ToByteArray(uint64(lrpc.pid))...) // process id
	seq := rand.Int31()
	//fmt.Printf("Set sequence number: %d\n", seq)
	//fmt.Printf("Set sequence number bytes: %v\n", uint32ToByteArray(uint32(seq)))
	data = append(data, uint32ToByteArray(uint32(seq))...) // sequence number
	now := time.Now()
	data = append(data, uint64ToByteArray(uint64(now.Unix()))...) // seconds passed since epoch
	length := uint32(len(payload))
	if length > lrpc.maxPayloadLength {
		return nil, fmt.Errorf("LRPC payload length %d exceeds the max limit %d", length, lrpc.maxPayloadLength)
	}
	data = append(data, uint32ToByteArray(length)...)
	//fmt.Printf("Header: %v\n", data)
	data = append(data, payload...) // payload

	// Send request header
	_, err := lrpc.client.Write(data) // send request
	if err != nil {
		return nil, fmt.Errorf("Write error: %v", err)
	}

	// Read response
	//fmt.Printf("Parsing LRPC response...\n")
	scan := true
	var items []interface{}
	for scan {
		// Read the first 34 bytes
		header, _ := reader(lrpc.client, int(headerLength))
		//fmt.Printf("Respond header: %v\n", header)
		bmagicNumber := header[0:4]
		bheaderLength := header[4:6]
		_ = bheaderLength
		bversion := header[6:10]
		_ = bversion
		bpid := header[10:18]
		_ = bpid
		bseq := header[18:22]
		bts := header[22:30]
		_ = bts
		blength := header[30:34]
		payloadLength := byteArrayToUInt16(blength)
		payload, _ = reader(lrpc.client, int(payloadLength))
		//fmt.Printf("Respond payload: %v\n", payload)
		responseSeq := byteArrayToUInt32(bseq)
		//fmt.Printf("Respond sequence byte: %v\n", bseq)
		//fmt.Printf("Respond sequence: %d\n", responseSeq)

		// Compare magic number
		if bytes.Compare(bmagicNumber, uint32ToByteArray(magicNumber)) != 0 {
			return nil, fmt.Errorf("Unrecognized LRPC2 server")
		}
		// Compare sequence number
		if int32(responseSeq) != seq {
			// Skip unmatched seq number.  It could be the response from prevoius requests
			fmt.Printf("Skip unmatched seq number\n")
			continue
		}
		command := payload[0:2]
		_ = command
		payload = payload[2:]

		// Process payload
		for len(payload) > 0 {
			//fmt.Printf("Current payload in loop: %v\n", payload)
			itemType := payload[0:1]
			//fmt.Printf("Item type: %v\n", itemType)
			if bytes.Compare(itemType, []byte{msgDataTypeInt32}) == 0 {
				//fmt.Printf("Detected Int32 item")
				item := byteArrayToUInt32(payload[1:5])
				items = append(items, item)
				payload = payload[5:]
				//fmt.Printf("Int32 item: %v\n", item)
			} else if bytes.Compare(itemType, []byte{msgDataTypeString}) == 0 {
				//fmt.Printf("Detected string item")
				strlen := int32(byteArrayToUInt32(payload[1:5]))
				if strlen < 0 {
					items = append(items, "")
					payload = payload[5:]
					//fmt.Printf("String item: %v\n", "")
				} else {
					item := string(payload[5 : 5+strlen])
					items = append(items, item)
					payload = payload[5+strlen:]
					//fmt.Printf("String item: %v\n", item)
				}
			} else if bytes.Compare(itemType, []byte{msgDataTypeSet}) == 0 {
				//fmt.Printf("Detected set item")
				count := byteArrayToUInt32(payload[1:5])
				payload = payload[5:]
				var strset []string
				for i := 0; i < int(count); i++ {
					strlen := byteArrayToUInt32(payload[1:5])
					if strlen < 0 {
						strset = append(strset, "")
						payload = payload[5:]
					} else {
						item := string(payload[5 : 5+strlen])
						strset = append(strset, item)
						payload = payload[5+strlen:]
					}
				}
				//fmt.Printf("Set item: %v\n", strset)
				items = append(items, strset)
			} else if bytes.Compare(itemType, []byte{msgEnd}) == 0 {
				scan = false
				//fmt.Printf("End of message\n")
				break
			} else {
				return nil, fmt.Errorf("Unrecognized data type %v", itemType)
			}
		}

		//fmt.Printf("Items: %v\n", items)
	}

	return items, nil
}

// GetToken gets dmc token from Centrify Client service
func (lrpc *LRPC2) GetToken(scope string) (string, error) {
	// Connect to lrpc server
	err := lrpc.connect()
	if err != nil {
		return "", fmt.Errorf("Failed to connect to server: %v", err)
	}
	defer lrpc.client.Close()

	// Construct payload
	var payload []byte
	payload = append(payload, uint16ToByteArray(msgIDAdminClientGetToken)...)
	payload = append(payload, []byte{msgDataTypeString}...)
	payload = append(payload, uint32ToByteArray(uint32(len(scope)))...)
	payload = append(payload, []byte(scope)...)
	payload = append(payload, []byte{msgEnd}...)
	//fmt.Printf("Payload: %v\n", payload)

	// Send payload to lrpc server
	reply, err := lrpc.request(payload)
	if err != nil {
		return "", fmt.Errorf("Request error: %v", err)
	}

	//fmt.Printf("reply %+v\n", reply)
	if len(reply) > 2 {
		statusCode := fmt.Sprintf("%v", reply[0])
		if statusCode != "0" {
			return "", fmt.Errorf("%+v", reply[1])
		}
		if reflect.TypeOf(reply[2]).Kind() != reflect.String {
			return "", fmt.Errorf("Invaid reply from server: %+v", reply)
		}
		return reply[2].(string), nil
	}

	return "", fmt.Errorf("Invaid reply from server: %+v", reply)
}

func uint16ToByteArray(num uint16) []byte {
	arr := make([]byte, 2)
	binary.LittleEndian.PutUint16(arr, num)
	return arr
}

func uint32ToByteArray(num uint32) []byte {
	arr := make([]byte, 4)
	binary.LittleEndian.PutUint32(arr, num)
	return arr
}

func uint64ToByteArray(num uint64) []byte {
	arr := make([]byte, 8)
	binary.LittleEndian.PutUint64(arr, num)
	return arr
}

func byteArrayToUInt32(arr []byte) uint32 {
	val := uint32(0)
	val = binary.LittleEndian.Uint32(arr)
	return val
}

func byteArrayToUInt16(arr []byte) uint16 {
	val := uint16(0)
	val = binary.LittleEndian.Uint16(arr)
	return val
}
