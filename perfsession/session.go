// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package perfsession

import "github.com/aclements/goperf/perffile"

type Session struct {
	kernel  *PIDInfo
	pidInfo map[int]*PIDInfo
}

func New() *Session {
	kernel := &PIDInfo{Comm: "[kernel]"}
	return &Session{
		kernel: kernel,
		pidInfo: map[int]*PIDInfo{
			// The kernel is implicitly PID -1
			-1: kernel,
		},
	}
}

func (s *Session) Update(r perffile.Record) {
	ensurePID := func(pid int) *PIDInfo {
		pidInfo, ok := s.pidInfo[pid]
		if !ok {
			pidInfo = &PIDInfo{kernel: s.kernel}
			s.pidInfo[pid] = pidInfo
		}
		return pidInfo
	}

	switch r := r.(type) {
	case *perffile.RecordComm:
		ensurePID(r.PID).Comm = r.Comm

	case *perffile.RecordExit:
		if r.PID == r.TID {
			delete(s.pidInfo, r.PID)
		}
		// Otherwise this is thread exit

	case *perffile.RecordFork:
		if r.PID == r.TID {
			s.pidInfo[r.PID] = ensurePID(r.PPID).fork()
		}
		// Otherwise this is thread creation

	case *perffile.RecordMmap:
		info := ensurePID(r.PID)
		info.munmap(r.Addr, r.Len)
		info.maps = append(info.maps, &Mmap{nil, *r})

	case *perffile.RecordSample:
		// Sometimes (particularly early in sample files), we
		// see kernel samples before the RecordComm.
		ensurePID(r.PID)
	}
}

func (s *Session) LookupPID(pid int) *PIDInfo {
	return s.pidInfo[pid]
}

type PIDInfo struct {
	Extra ForkableExtra

	Comm   string
	kernel *PIDInfo
	maps   []*Mmap
}

func (p *PIDInfo) fork() *PIDInfo {
	maps := make([]*Mmap, len(p.maps))
	for i, mmap := range p.maps {
		maps[i] = mmap.fork()
	}
	return &PIDInfo{forkExtra(p.Extra), p.Comm, p.kernel, maps}
}

func (p *PIDInfo) munmap(addr, mlen uint64) {
	end := addr + mlen
	removed := false
	nmaps := p.maps
	for i, mmap := range p.maps {
		if addr <= mmap.Addr {
			if end >= mmap.Addr+mmap.Len {
				p.maps[i] = nil
				removed = true
			} else if end > mmap.Addr {
				// Remove beginning of mmap
				mmap.Len -= (end - mmap.Addr)
				mmap.Addr = end
			}
		} else if addr < mmap.Addr+mmap.Len {
			if end >= mmap.Addr+mmap.Len {
				// Remove end of mmap
				mmap.Len = addr - mmap.Addr
			} else {
				// Split mmap in two
				nmmap := *mmap
				nmmap.Len = end - (mmap.Addr + mmap.Len)
				nmaps = append(nmaps, &nmmap)
				mmap.Len = addr - mmap.Addr
			}
		}
	}
	// Fill holes
	if removed {
		d := 0
		for s := 0; s < len(nmaps); s++ {
			if nmaps[d] == nil {
				nmaps[d] = nmaps[s]
			}
			if nmaps[d] != nil {
				d++
			}
		}
		nmaps = nmaps[:d]
	}
	p.maps = nmaps
}

func (p *PIDInfo) mapFind(addr uint64) *Mmap {
	for _, mmap := range p.maps {
		if mmap.Addr <= addr && addr < mmap.Addr+mmap.Len {
			return mmap
		}
	}
	return nil
}

func (p *PIDInfo) LookupMmap(addr uint64) *Mmap {
	m := p.mapFind(addr)
	if m == nil && p.kernel != nil {
		m = p.kernel.mapFind(addr)
	}
	return m
}

type Mmap struct {
	Extra ForkableExtra

	perffile.RecordMmap
}

func (m *Mmap) fork() *Mmap {
	return &Mmap{forkExtra(m.Extra), m.RecordMmap}
}

type ForkableExtra interface {
	Fork() ForkableExtra
}

func forkExtra(extra ForkableExtra) ForkableExtra {
	if extra == nil {
		return nil
	}
	return extra.Fork()
}