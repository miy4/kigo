package kigo

type Key int16

const (
	KeyCtrlA = Key(iota + 1)
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
)

const (
	KeyUp = Key(iota + 128)
	KeyDown
	KeyRight
	KeyLeft
	KeyPgUp
	KeyPgDn
	KeyHome
	KeyEnd
	KeyDel
)

var keySeq *node

type node struct {
	children map[byte]*node
	values   []Key
}

func (n *node) insert(seq string, result Key) {
	cur := n
	keys := make([]Key, 0)

	for _, b := range []byte(seq) {
		keys = append(keys, Key(b))

		child, ok := cur.children[b]
		if !ok {
			child = &node{
				children: make(map[byte]*node),
				values:   keys,
			}
			cur.children[b] = child
		}
		cur = child
	}

	cur.values = []Key{result}
}

func init() {
	keySeq = &node{
		children: make(map[byte]*node),
		values:   []Key{},
	}

	keySeq.insert("\x1b[A", KeyUp)
	keySeq.insert("\x1b[B", KeyDown)
	keySeq.insert("\x1b[C", KeyRight)
	keySeq.insert("\x1b[D", KeyLeft)
	keySeq.insert("\x1b[5~", KeyPgUp)
	keySeq.insert("\x1b[6~", KeyPgDn)
	keySeq.insert("\x1b[1~", KeyHome)
	keySeq.insert("\x1b[7~", KeyHome)
	keySeq.insert("\x1b[H", KeyHome)
	keySeq.insert("\x1bOH", KeyHome)
	keySeq.insert("\x1b[4~", KeyEnd)
	keySeq.insert("\x1b[8~", KeyEnd)
	keySeq.insert("\x1b[F", KeyEnd)
	keySeq.insert("\x1bOF", KeyEnd)
	keySeq.insert("\x1b[3~", KeyDel)
}
