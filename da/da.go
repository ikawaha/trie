package da

const (
	_BUFSIZE    = 1024
	_TERMINATOR = '\x00'
)

type doubleArray []struct {
	base, check int
}

func NewDoubleArray() *doubleArray {
	var da *doubleArray = new(doubleArray)
	*da = make(doubleArray, _BUFSIZE)
	(*da)[0].base = 1
	(*da)[0].check = -1
	for i, len := 0, len(*da); i < len; i++ {
		(*da)[i].check = -1
	}
	return da
}

func (this *doubleArray) Add(a_keyword string, a_id int) {
	str := []byte(a_keyword + string(_TERMINATOR))
	p, q, i := this.search(str)
	for q >= len(*this) {
		this.expand()
	}
	if (*this)[q].check == p && (*this)[q].base < 0 {
		return
	}
	if (*this)[q].check < 0 {
		(*this)[q].check = p
	} else {
		q = this.rearrange(p, q, i, str)
	}
	this.add(q, i+1, str, a_id)
}

func (this *doubleArray) Search(a_keyword string) (id int, ok bool) {
	p, q, _ := this.search([]byte(a_keyword + string(_TERMINATOR)))
	if (*this)[q].check != p || (*this)[q].base > 0 {
		return 0, false
	}
	return -(*this)[q].base, true
}

func (this *doubleArray) CommonPrefixSearch(a_keyword string) (keywords []string, ids []int) {
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

func (this *doubleArray) PrefixSearch(a_keyword string) (keyword string, id int, ok bool) {
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

func (this *doubleArray) Size() int {
	return len(*this)
}

func (this *doubleArray) Efficiency() (int, int, float64) {
	unspent := 0
	for _, pair := range *this {
		if pair.check < 0 {
			unspent++
		}
	}
	return unspent, len(*this), float64(len(*this)-unspent) / float64(len(*this)) * 100
}

func (this *doubleArray) shrink() {
	src_size := len(*this)
	for i, size := 0, src_size; i < size; i++ {
		if (*this)[size-i-1].check < 0 {
			src_size--
		} else {
			break
		}
	}
	if src_size == len(*this) {
		return
	}
	var dst *doubleArray = new(doubleArray)
	*dst = make(doubleArray, src_size)
	copy(*dst, (*this)[:src_size])
	*this = *dst
}

func (this *doubleArray) expand() {
	src_size := len(*this)
	var dst *doubleArray = new(doubleArray)
	*dst = make(doubleArray, src_size+_BUFSIZE)
	copy(*dst, *this)
	for i, size := src_size, len(*dst); i < size; i++ {
		(*dst)[i].check = -1
	}
	*this = *dst
}

func (this *doubleArray) search(a_str []byte) (p, q, i int) {
	p, q, i = 0, 0, 0
	buf_size := len(*this)
	for size := len(a_str); i < size; i++ {
		p = q
		ch := int(a_str[i])
		q = (*this)[p].base + ch
		if q >= buf_size || (*this)[q].check != p {
			return p, q, i
		}
	}
	return p, q, i
}

func (this *doubleArray) seek(a_p, a_ch int, a_stack map[int]struct{ base, check, ch int }) (q int) {
	q = len(*this)
	for i, size := 1, len(*this); i < size; i++ {
	L_start:
		s := a_ch + i
		if s >= size || (*this)[s].check < 0 {
			base := s - a_ch
			for _, tuple := range a_stack {
				if base+tuple.ch < size && (*this)[base+tuple.ch].check >= 0 {
					i++
					goto L_start
				}
			}
			q = s
			break
		}
	}

	for q >= len(*this) {
		this.expand()
	}
	base := q - a_ch
	table := make(map[int]int)
	for old_p, tuple := range a_stack {
		for base+tuple.ch >= len(*this) {
			this.expand()
		}
		neo_p := base + tuple.ch
		for neo_p >= len(*this) {
			this.expand()
		}
		(*this)[neo_p] = struct{ base, check int }{tuple.base, tuple.check}
		table[old_p] = neo_p
	}
	if len(table) > 0 {
		for index, bc := range *this {
			if neo_p, ok := table[bc.check]; ok {
				(*this)[index].check = neo_p
			}
		}
	}
	return q
}

func (this *doubleArray) add(a_p, a_i int, a_str []byte, a_id int) {
	p := a_p
	for i, size := a_i, len(a_str); i < size; i++ {
		ch := int(a_str[i])
		q := this.seek(p, ch, nil)
		(*this)[p].base = q - ch
		(*this)[q].check = p
		p = q
	}
	(*this)[p].base = -a_id
}

func (this *doubleArray) rearrange(a_p, a_q, a_i int, a_str []byte) int {
	stack := make(map[int]struct{ base, check, ch int })
	for i, size := 0, len(*this); i < size; i++ {
		if (*this)[i].check == a_p {
			ch := i - (*this)[a_p].base
			stack[i] = struct{ base, check, ch int }{(*this)[i].base, (*this)[i].check, ch}
			(*this)[i].base, (*this)[i].check = 0, -1
		}
	}
	ch := int(a_str[a_i])
	q := this.seek(a_p, ch, stack)
	(*this)[a_p].base = q - ch
	(*this)[q].check = a_p
	return q
}
