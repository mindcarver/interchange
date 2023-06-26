package types

const (
	ShareAmtKeyPrefix = "ShareAmt/"
	ReadyKeyPrefix    = "Ready/"
)

func ShareAmtKey(
	addr string,
) []byte {
	var key []byte

	indexBytes := []byte(addr)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

func ReadyKeyKey(
	addr string,
) []byte {
	var key []byte

	indexBytes := []byte(addr)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
