package world

import (
	"log"
	"sync/atomic"
)

type counter atomic.Uint64

func (c *counter) inc() {
	atom := (*atomic.Uint64)(c)
	done := false
	for !done {
		val := atom.Load()
		done = atom.CompareAndSwap(val, val+1)
	}
}

func (c *counter) val() uint64 {
	return (*atomic.Uint64)(c).Load()
}

type Stats struct {
	eyeRayCount        counter
	shadowRayCount     counter
	reflectionRayCount counter
	refractionRayCount counter
}

func (s *Stats) EyeRayCount() uint64 {
	return s.eyeRayCount.val()
}
func (s *Stats) ShadowRayCount() uint64 {
	return s.shadowRayCount.val()
}
func (s *Stats) ReflectionRayCount() uint64 {
	return s.reflectionRayCount.val()
}
func (s *Stats) RefractionRayCount() uint64 {
	return s.refractionRayCount.val()
}

func (s *Stats) TotalRayCount() uint64 {
	return s.eyeRayCount.val() + s.shadowRayCount.val() + s.reflectionRayCount.val() + s.refractionRayCount.val()
}

func (s *Stats) Log() {
	log.Printf("eye rays: %v, shadow rays: %v, reflection rays: %v, refraction rays: %v, total: %v",
		s.eyeRayCount.val(), s.shadowRayCount.val(), s.reflectionRayCount.val(), s.refractionRayCount.val(), s.TotalRayCount())
}
