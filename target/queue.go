package target

import (
	"math"
	"os"
	"strconv"
)

// Target of keeping the number of items in a queue under a certain length
type QueueLengthTarget struct {
	length     int
	minLength  int
	lastLength int
	vel        []int
	velSamples int
	cumErr     int
	kP         float64
	kI         float64
	kD         float64
	startCount int
}

const queueLengthExceedingPercent float64 = 0.7
const queueAverageSamples int = 1

func NewQueueLengthTarget(length int) Target {
	// TODO!! Better ways to calculate these heuristics and/or pass them in
	kUStr := os.Getenv("F12_KU")
	tUStr := os.Getenv("F12_TU")
	kU, err := strconv.ParseFloat(kUStr, 64)
	if kUStr == "" || err != nil {
		kU = 0.05
	}

	tU, err := strconv.ParseFloat(tUStr, 64)
	if tUStr == "" || err != nil {
		tU = 10.0
	}

	var velSamples int
	velSamplesStr := os.Getenv("F12_VEL_SAMPLES")
	velSamples, err = strconv.Atoi(velSamplesStr)
	if err != nil || velSamples == 0 {
		velSamples = queueAverageSamples
	}

	// Ziegler-Nichols PID
	// kP := 0.6 * kU
	// kI := kP * 2.0 / tU
	// kD := kP * tU / 8.0

	// Ziegler-Nichols PD
	kP := float64(0.8 * kU)
	kD := float64(tU / 8.0)
	kI := float64(0) // Seems like PD works better than PID, but this needs more testing
	log.Debugf("[new ql] kU = %f, tU = %f", kU, tU)
	log.Debugf("[new ql] kP = %f, kI = %f, kD = %f", kP, kI, kD)

	return &QueueLengthTarget{
		length:     length,
		minLength:  int(float64(length) * queueLengthExceedingPercent),
		kP:         kP,
		kI:         kI,
		kD:         kD,
		vel:        make([]int, velSamples+1),
		velSamples: velSamples,
	}
}

func (t *QueueLengthTarget) Meeting(current int) bool {
	meeting := (current <= t.length)
	if !meeting {
		log.Debugf("[ql] not meeting: current %d target %d", current, t.length)
	}
	return meeting
}

func (t *QueueLengthTarget) Exceeding(current int) bool {
	exceeding := (current <= t.minLength)
	if exceeding {
		log.Debugf("[ql] exceeding: current %d target %d", current, t.length)
	}
	return exceeding
}

// Number of additional containers
func (t *QueueLengthTarget) Delta(currentLength int) (delta int) {
	var deltafloat float64
	var currErr int

	currErr = currentLength - t.length
	t.cumErr = t.cumErr + currErr

	var aveVel float64
	// Store the new value at the end of the array (this is one more than we need), but only average over the right number of samples
	t.vel[t.velSamples] = currentLength - t.lastLength
	for i := 0; i <= t.velSamples-1; i++ {
		t.vel[i] = t.vel[i+1]
		aveVel = aveVel + float64(t.vel[i])
	}

	aveVel = aveVel / float64(t.velSamples)

	// There is a point beyond which there is no point letting cumErr grow, because our max containers can't
	// necessarily keep up (and also a question of symmetry, since a queue length can't go below 0?)
	if t.cumErr > 10*t.length {
		t.cumErr = 10 * t.length
	}
	if t.cumErr < -10*t.length {
		t.cumErr = -10 * t.length
	}

	t.lastLength = currentLength

	// To start with, velocity isn't valid
	if t.startCount < t.velSamples {
		log.Debugf("[ql] err %d, cumErr %d", currErr, t.cumErr)
		log.Debugf("[ql] err * kp %f, cumErr * kI %f", t.kP*float64(currErr), t.kI*float64(t.cumErr))
		deltafloat = t.kP*float64(currErr) + t.kI*float64(t.cumErr)
		t.startCount = t.startCount + 1
	} else {
		log.Debugf("[ql] err %d, cumErr %d, vel %f", currErr, t.cumErr, aveVel)
		log.Debugf("[ql] err * kp %f, cumErr * kI %f, vel * kd %f", t.kP*float64(currErr), t.kI*float64(t.cumErr), t.kD*float64(aveVel))
		deltafloat = t.kP*float64(currErr) + t.kI*float64(t.cumErr) + t.kD*float64(aveVel)
	}

	log.Debugf("[ql] => deltaf %f", deltafloat)

	// Round to the nearest integer
	delta = int(math.Floor(deltafloat + 0.5))
	return
}
