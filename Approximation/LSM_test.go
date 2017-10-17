package Approximation

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/matrix/mat64"
	"github.com/philipp-altmann/ContinuousBenchmarkOptimizer/Bohachevsky"
	p "github.com/philipp-altmann/ContinuousBenchmarkOptimizer/Population"
)

/*func FitnessFunctionTest(value []float64) (result float64) {
	return value[0]
}*/
var testPopulation p.Population
var testApproximator mat64.Dense

const size = 64
const dimensions = 16

func TestGenerateFitnessVector(t *testing.T) {
	testPopulation = p.InitRandomPopulation(size, dimensions)
	testFitnessVector := generateFitnessVector(testPopulation)
	//Test Fitness Vector
	for i := 0; i < size; i++ {
		if testPopulation[i].Fitness != testFitnessVector.At(i, 0) {
			t.Errorf("Mismatching Results Testing Fitness Vector (%f||%f)", testPopulation[i].Fitness, testFitnessVector.At(i, 0))
			t.Fail()
		}
	}
}

func TestGenerateValueMatrix(t *testing.T) {
	testPopulation = p.InitRandomPopulation(size, dimensions)
	testValueMatrix := generateValueMatrix(testPopulation)
	//Test Value matrix
	for i := 0; i < size; i++ {
		for j := 0; j < dimensions; j++ {
			if testPopulation[i].Value[j] != testValueMatrix.At(i, j) {
				t.Errorf("Mismatching Results Testing Value Matrix (%f||%f)", testPopulation[i].Value[j], testValueMatrix.At(i, j))
				t.Fail()
			}
		}
	}
}

func TestGetLSM(t *testing.T) {

	testPopulation = p.InitRandomPopulation(size, dimensions)

	//Set fitness
	for i := 0; i < size; i++ {
		testPopulation[i].Fitness = linearFitness(testPopulation[i].Value[0])
	}

	testApproximator = GetLSMApproximator(testPopulation)
	r, c := testApproximator.Dims()
	fmt.Printf("Rows: %d, Cols: %d", r, c)

	/*
		individual := i.GenerateIndividual(value)
		testI := i.GenerateIndividual(test)

		testI.EvaluateFitness()
		testGeneration := p.InitRandomPopulation(20, 10)
		//fmt.Print(testGeneration)
		Approximator := GetLSMApproximator(testGeneration)
		result := ApproximateFitness(test, Approximator)
		fmt.Printf("Evaluated: %f\nApproximated: %f\n", testI.Fitness, result)

		result = ApproximateFitness(testGeneration[0].Value, Approximator)
		fmt.Printf("Evaluated: %f\nApproximated: %f\n", testGeneration[0].Fitness, result)
	*/
}

func TestApproximation(t *testing.T) {
	//Test for right Approximation Behavior
	//\w linear one-dimensional fitness function
	testPopulation = p.InitRandomPopulation(size, 1)

	//Set fitness
	for i := 0; i < size; i++ {
		testPopulation[i].Fitness = linearFitness(testPopulation[i].Value[0])
	}

	testApproximator = GetLSMApproximator(testPopulation)

	tests := 10
	for i := 0; i < tests; i++ {
		testValue := []float64{rand.Float64()*200 - 100}
		approximationResult := ApproximateFitness(testValue, testApproximator)
		if math.Floor(approximationResult) != math.Floor(linearFitness(testValue[0])) {
			t.Errorf("Mismatching Results Testing Linear Fitness Appromimation (%f||%f)∆%f", approximationResult, linearFitness(testValue[0]), approximationResult-linearFitness(testValue[0]))
			t.Fail()
		}
	}

}

func linearFitness(x float64) float64 {
	return 0.8 * x
}

// from fib_test.go
func benchmarkLSM(b *testing.B, dimensions int) {
	size := 32

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		if n%100 == 0 {
			testPopulation = p.InitRandomPopulation(size, dimensions)
			testPopulation.Evaluate(Bohachevsky.EvaluateFitness)
			testApproximator = GetLSMApproximator(testPopulation)
		}

		var vals []float64
		for v := 0; v < dimensions; v++ {
			vals = append(vals, rand.Float64()*200-100)
		}
		ApproximateFitness(vals, testApproximator)

	}
}

func BenchmarkLSM4(b *testing.B)   { benchmarkLSM(b, 4) }
func BenchmarkLSM8(b *testing.B)   { benchmarkLSM(b, 8) }
func BenchmarkLSM16(b *testing.B)  { benchmarkLSM(b, 16) }
func BenchmarkLSM32(b *testing.B)  { benchmarkLSM(b, 32) }
func BenchmarkLSM64(b *testing.B)  { benchmarkLSM(b, 64) }
func BenchmarkLSM128(b *testing.B) { benchmarkLSM(b, 128) }