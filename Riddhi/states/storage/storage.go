package storage

import "riddhi/states/protocol/chunk"
 
// storage holds persistent file/content data
// independent of active protocol sessions
var ChunkedFile *chunk.ChunkedFile
