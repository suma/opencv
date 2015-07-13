package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo LDFLAGS: -ljsonconfig
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
#include <stdlib.h>
#include "tracker_bridge.h"
*/
import "C"

func convertCandidatezToPointer(candidatez [][]Candidate) ([]*C.Candidate, []int32) {
	candidatesPointers := []*C.Candidate{}
	lengths := []int32{}
	for _, cz := range candidatez {
		candidatePointers := convertCandidatesToPointer(cz)
		candidatesPointers = append(candidatesPointers, (*C.Candidate)(&candidatePointers[0]))
		lengths = append(lengths, int32(len(candidatePointers)))
	}
	return candidatesPointers, lengths
}

func GetMatching(kThreashold float32, frames ...[]Candidate) [][]Candidate {
	candidatesPointers, ls := convertCandidatezToPointer(frames)
	views := C.MVOM_GetMatching((**C.Candidate)(&candidatesPointers[0]),
		(*C.int)(&ls[0]), C.int(len(frames)), C.float(kThreashold))
	defer C.Candidatez_Delete(views)

	l := int(views.length)
	viewsP := make([]C.Candidates, l)
	C.ResolveCandidatez(views, (*C.Candidates)(&viewsP[0]))

	ret := make([][]Candidate, l)
	for i := 0; i < l; i++ {
		regionLength := int(viewsP[i].length)
		candidates := make([]C.Candidate, regionLength)
		C.ResolveCandidates(viewsP[i], (*C.Candidate)(&candidates[0]))

		mvRegions := make([]Candidate, regionLength)
		for j := 0; j < regionLength; j++ {
			mvRegions[j] = Candidate{p: candidates[j]}
		}
		ret[i] = mvRegions
	}
	return ret
}
