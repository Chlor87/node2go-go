package node2go

var nilID = []byte("0")

func formatResponse(res []byte, id []byte) []byte {
	return append([]byte(append(id, ';')), append(res, '\n')...)
}
