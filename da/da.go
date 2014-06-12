package da

import (
	"sort"
)

const (
	_BUFSIZE     = 51200
	_EXPANDRATIO = 2
	_TERMINATOR  = '\x00'
	_ROOTID      = 0
)

type DoubleArray []struct {
	base, check int
}

func NewDoubleArray() *DoubleArray {
	da := new(DoubleArray)
	*da = make(DoubleArray, _BUFSIZE)

	(*da)[_ROOTID].base = 1
	(*da)[_ROOTID].check = -1

	size := len(*da)
	for i := 1; i < size; i++ {
		(*da)[i].base = -(i - 1)
		(*da)[i].check = -(i + 1)
	}

	(*da)[1].base = -(size - 1)
	(*da)[size-1].check = -1

	return da
}

func (this *DoubleArray) Search(a_keyword string) (id int, ok bool) {
	p, q, _ := this.search([]byte(a_keyword + string(_TERMINATOR)))
	if (*this)[q].check != p || (*this)[q].base > 0 {
		return 0, false
	}
	return -(*this)[q].base, true
}

func (this *DoubleArray) CommonPrefixSearch(a_keyword string) (keywords []string, ids []int) {
	keywords, ids = make([]string, 0), make([]int, 0)
	p, q, i := 0, 0, 0
	str := []byte(a_keyword)
	buf_size := len(*this)
	for size := len(str); i < size; i++ {
		p = q
		ch := int(str[i])
		q = (*this)[p].base + ch
		if q >= buf_size || (*this)[q].check != p {
			break
		}
		ahead := (*this)[q].base + _TERMINATOR
		if ahead < buf_size && (*this)[ahead].check == q && (*this)[ahead].base < 0 {
			keywords = append(keywords, string(str[0:i+1]))
			ids = append(ids, -(*this)[ahead].base)
		}
	}
	return
}

func (this *DoubleArray) PrefixSearch(a_keyword string) (keyword string, id int, ok bool) {
	p, q, i := 0, 0, 0
	str := []byte(a_keyword)
	buf_size := len(*this)
	for size := len(str); i < size; i++ {
		p = q
		ch := int(str[i])
		q = (*this)[p].base + ch
		if q >= buf_size || (*this)[q].check != p {
			break
		}
		ahead := (*this)[q].base + _TERMINATOR
		if ahead < buf_size && (*this)[ahead].check == q && (*this)[ahead].base < 0 {
			keyword = string(str[0 : i+1])
			id = -(*this)[ahead].base
			ok = true
		}
	}
	return
}

func (this *DoubleArray) expand() {
	srcSize := len(*this)
	dst := new(DoubleArray)
	dstSize := srcSize * _EXPANDRATIO
	*dst = make(DoubleArray, dstSize)
	copy(*dst, *this)

	for i := srcSize; i < dstSize; i++ {
		(*dst)[i].base = -(i - 1)
		(*dst)[i].check = -(i + 1)
	}

	start := -(*this)[0].check
	end := -(*dst)[start].base
	(*dst)[srcSize].base = -end
	(*dst)[start].base = -(dstSize - 1)
	(*dst)[end].check = -srcSize
	(*dst)[dstSize-1].check = -start

	*this = *dst
}

func (this *DoubleArray) search(a_str []byte) (p, q, i int) {
	p, q, i = 0, 0, 0
	bufSize := len(*this)
	for size := len(a_str); i < size; i++ {
		p = q
		ch := int(a_str[i])
		q = (*this)[p].base + ch
		if q >= bufSize || (*this)[q].check != p {
			return p, q, i
		}
	}
	return p, q, i
}

func (this *DoubleArray) Build(a_keywords []string) {
	list := a_keywords
	if len(list) == 0 {
		return
	}
	if !sort.StringsAreSorted(list) {
		sort.Strings(list)
	}
	branches := make([]int, len(a_keywords))
	for i, size := 0, len(a_keywords); i < size; i++ {
		branches[i] = i
	}
	this.append(0, 0, branches, a_keywords)
}

func (this *DoubleArray) setBase(a_p, a_base int) {
	if a_p == _ROOTID {
		return
	}
	if (*this)[a_p].check < 0 {
		if (*this)[a_p].base == (*this)[a_p].check {
			this.expand()
		}
		prev := -(*this)[a_p].base
		next := -(*this)[a_p].check
		if -a_p == (*this)[_ROOTID].check {
			(*this)[_ROOTID].check = (*this)[a_p].check
		}
		(*this)[next].base = (*this)[a_p].base
		(*this)[prev].check = (*this)[a_p].check
	}
	(*this)[a_p].base = a_base
}

func (this *DoubleArray) setCheck(a_p, a_check int) {
	if (*this)[a_p].base == (*this)[a_p].check {
		this.expand()
	}
	prev := -(*this)[a_p].base
	next := -(*this)[a_p].check
	if -a_p == (*this)[_ROOTID].check {
		(*this)[_ROOTID].check = (*this)[a_p].check
	}

	(*this)[next].base = (*this)[a_p].base
	(*this)[prev].check = (*this)[a_p].check
	(*this)[a_p].check = a_check

}

func (this *DoubleArray) seekAndMark(a_p int, a_chars []byte) { // chars != nil
	free := _ROOTID
	rep := int(a_chars[0])
	var base int
	for {
	L_start:
		if free != _ROOTID && (*this)[free].check == (*this)[_ROOTID].check {
			this.expand()
		}
		free = -(*this)[free].check
		base = free - rep
		if base <= 0 {
			continue
		}
		for _, ch := range a_chars {
			q := base + int(ch)
			if q < len(*this) && (*this)[q].check >= 0 {
				goto L_start
			}
		}
		break
	}
	this.setBase(a_p, base)
	for _, ch := range a_chars {
		q := (*this)[a_p].base + int(ch)
		if q >= len(*this) {
			this.expand()
		}
		this.setCheck(q, a_p)
	}
}

func (this *DoubleArray) append(a_p, a_i int, a_branches []int, a_keywords []string) {
	chars := make([]byte, 0)
	subtree := make(map[byte][]int)
	for _, keyId := range a_branches {
		str := []byte(a_keywords[keyId])
		var ch byte
		if a_i >= len(str) {
			ch = _TERMINATOR
		} else {
			ch = str[a_i]
		}
		if size := len(chars); size == 0 || chars[len(chars)-1] != ch {
			chars = append(chars, ch)
		}
		if ch != _TERMINATOR {
			subtree[ch] = append(subtree[ch], keyId)
		}
	}
	this.seekAndMark(a_p, chars)
	for _, ch := range chars {
		q := (*this)[a_p].base + int(ch)
		if len(subtree[ch]) == 0 {
			(*this)[q].base = -a_branches[0]
		} else {
			this.append(q, a_i+1, subtree[ch], a_keywords)
		}
	}
}
