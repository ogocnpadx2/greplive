// Package splice fans in multiple line-producing channels into a single
// merged output channel.
//
// It is useful when greplive consumes from more than one source simultaneously
// (e.g. multiple tailed files or a combination of stdin and a named pipe) and
// the downstream pipeline expects a single stream of lines.
//
// Lines from all sources are forwarded on a first-ready basis; no ordering
// guarantee is provided across sources.
//
// The output channel is closed automatically once every source channel has
// been drained, or when the supplied context is cancelled.
package splice
