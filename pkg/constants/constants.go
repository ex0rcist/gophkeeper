// Holds useful app wide constants
package constants

import "time"

type CtxKey string

var (
	// Access token key name for gRPC calls
	AccessTokenHeader = "Access-Token"

	// Header for client ID
	ClientIDHeader = "Client-ID"

	// Context key name for user_id storage
	CtxUserIDKey CtxKey = "user_id"

	// How long should client wait for server response
	DefaultClientTimeout = time.Second * 5

	// Default time format
	TimeFormat = "2006-01-02 15:04"

	// Chunk size for file uplaod and download
	ChunkSize = 1024 * 1024 * 2
)
