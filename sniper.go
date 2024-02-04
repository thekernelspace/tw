package main

import "log"

// sniper mode
// press f to toggle, esc when in mode to exit out of the mode
// when in sniper mode, the user can press a key to jump to the dirent that matches the key pressed
// we keep a map of the key to press -> dir pos + cursor to jump to

var mapKeyToDirent map[string]*Dirent
var mapDirentToKey map[string]string

func getNthAlpha(n int) string {
	return string(rune(n + 97))
}

func (m Model) buildMapKeyToDirPos() {
	// build the map
	m.root.Walk(func(d Dirent) {
		if m.root.Equals(d) {
			return // we don't need this
		}

		parentAnnotation := ""
		if d.Parent != nil {
			parentAnnotation = mapDirentToKey[d.Parent.Path()]
		}
		annotation := parentAnnotation + getNthAlpha(d.PosInParent)
		mapKeyToDirent[annotation] = &d
		mapDirentToKey[d.Path()] = annotation
	})
}

func clearMaps() {
	mapKeyToDirent = make(map[string]*Dirent)
	mapDirentToKey = make(map[string]string)
}

func (m *Model) ToggleSniper() {
	m.sniper = !m.sniper
	if m.sniper {
		// clear the maps
		clearMaps()

		m.buildMapKeyToDirPos()
		log.Println("sniper mode on")
		log.Printf("map key -> dirent: %v\n", mapKeyToDirent)
		log.Printf("map dirent -> key: %v\n", mapDirentToKey)
	}
}

func (m *Model) QuitSniper() {
	clearMaps()
	m.sniper = false
	m.sniperKeyBuffer = ""
	log.Println("sniper mode off")
}

func getDirentFromKey(k string) *Dirent {
	return mapKeyToDirent[k]
}

func getKeyFromDirent(d Dirent) string {
	return mapDirentToKey[d.Path()]
}
