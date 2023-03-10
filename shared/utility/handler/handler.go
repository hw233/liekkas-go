package handler

import (
	"log"

	"google.golang.org/protobuf/proto"
)

// TODO: 抽象出来，支持多种解码引擎，实现Unmarshal和Marshal即可
type Handler struct {
	in  []byte
	out []byte
}

func NewHandler(in []byte) *Handler {
	return &Handler{
		in:  in,
		out: []byte{},
	}
}
func (h *Handler) In(in []byte) {
	h.in = in
}

func (h *Handler) Out() []byte {
	return h.out
}

func (h *Handler) UnmarshalProto(m proto.Message) error {
	err := proto.Unmarshal(h.in, m)
	if err != nil {
		log.Fatalln("Failed to parse address book:", err)
		return err
	}

	return nil
}

func (h *Handler) MarshalProto(m proto.Message) error {
	val, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	h.out = val

	return nil
}
