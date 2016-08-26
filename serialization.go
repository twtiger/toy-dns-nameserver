package nameserver

import (
	"encoding/binary"
	"errors"
)

func (m *message) serialize() ([]byte, error) {
	q, err := m.query.serialize()
	if err != nil {
		return nil, err
	}

	return append(serializeHeaders(), append(q, serializeAnswer(m.answers)...)...), nil
}

func (q *query) serialize() ([]byte, error) {
	l, err := serializeLabels(q.qname)
	if err != nil {
		return nil, err
	}

	return append(l, append(serializeUint16(uint16(q.qtype)), serializeUint16(uint16(q.qclass))...)...), nil
}

func (r *record) serialize() (b []byte) {
	l, _ := serializeLabels(r.Name)
	b = append(b, l...)

	b = append(b, serializeUint16(uint16(r.Type))...)

	b = append(b, serializeUint16(uint16(r.Class))...)

	b = append(b, serializeUint32(uint32(r.TTL))...)

	b = append(b, serializeUint16(r.RDLength)...)

	b = append(b, []byte(r.RData)...)
	return
}

func (l label) serialize() (b []byte) {
	b = append(b, byte(len(l)))
	b = append(b, []byte(l)...)
	return
}

func serializeLabels(l []label) ([]byte, error) {
	var b []byte
	if len(l) == 0 {
		return nil, errors.New("no labels to serialize")
	}

	for _, e := range l {
		b = append(b, e.serialize()...)
	}
	b = append(b, 0)
	return b, nil
}

func serializeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}

func serializeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func serializeAnswer(r []*record) (b []byte) {
	for _, e := range r {
		b = append(b, e.serialize()...)
	}
	return
}

func serializeHeaders() []byte {
	return make([]byte, 12)
}
