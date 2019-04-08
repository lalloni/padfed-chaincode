// MIT License
//
// Copyright (c) 2018 Pablo Ignacio Lalloni
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package periodo

import "strconv"

var (
	// MinYear establece el año mínimo válido
	MinYear uint = 1000
	// MaxYear establece el año máximo válido
	MaxYear uint = 9999
)

var (
	days     = []uint{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	minMonth = uint(0)
	maxMonth = uint(len(days) - 1)
)

const periododesconocido = "tipo período desconocido"

type tipoPeriodo uint8

const (
	// Diario es el tipo de los períodos con forma YYYYMMDD
	Diario tipoPeriodo = iota
	// Mensual es el tipo de los períodos con forma YYYYMM
	Mensual
	// Anual es el tipo de los períodos con forma YYYY
	Anual
)

// Parse intenta extraer un período fiscal de v del tipo especificado en t.
// En caso de éxito retorna true en match y las partes del período extraído
// correspondientes en y, m y d.
func Parse(t tipoPeriodo, v string) (match bool, y, m, d uint) {
	i, err := strconv.ParseUint(v, 10, strconv.IntSize)
	if err != nil {
		return false, 0, 0, 0
	}
	switch t {
	case Diario:
		y, m, d := DecomposePeriodoDiario(uint(i))
		return CheckPeriodoDiario(y, m, d), y, m, d
	case Mensual:
		y, m := DecomposePeriodoMensual(uint(i))
		return CheckPeriodoMensual(y, m), y, m, 0
	case Anual:
		y := uint(i)
		return CheckPeriodoAnual(y), y, 0, 0
	default:
		// Nunca debería ocurrir porque los usuarios de la librería no pueden
		// crear otras instancias de tipoPeriodo porque no se exporta.
		panic(periododesconocido)
	}
}

// ComposePeriodoDiario permite armar un período diario desde sus componentes.
// No realiza ningún tipo de validación, para eso usar CheckPeriodoDiario.
func ComposePeriodoDiario(y, m, d uint) uint {
	return y*10000 + m*100 + d
}

// DecomposePeriodoDiario permite desarmar un período diario a sus componentes.
// No realiza ningún tipo de validación, para eso usar CheckPeriodoDiarioCompound.
func DecomposePeriodoDiario(p uint) (y, m, d uint) {
	return p / 10000, p / 100 % 100, p % 100
}

// ComposePeriodoMensual permite armar un período mensual desde sus componentes.
// No realiza ningún tipo de validación, para eso usar CheckPeriodoMensual.
func ComposePeriodoMensual(y, m uint) uint {
	return y*100 + m
}

// DecomposePeriodoMensual permite desarmar un período mensual a sus componentes.
// No realiza ningún tipo de validación, para eso usar CheckPeriodoMensualCompound.
func DecomposePeriodoMensual(p uint) (y, m uint) {
	return p / 100, p % 100
}

// CheckPeriodoDiario valida que los componentes conformen un período diario correcto.
//
// Valida que:
//
// - El año (y) esté dentro del rango [MinYear, MaxYear] (que por defecto es [1000,9999])
// - El mes (m) esté dentro del rango [0,12]
// - Que el día (d):
//   - Si m = 0: sea igual a 0
//   - Si m > 0: esté dentro del rango [0, ds] siendo ds el correcto según el mes y año (considerando años bisiestos)
func CheckPeriodoDiario(y, m, d uint) bool {
	if y < MinYear || y > MaxYear || m < minMonth || m > maxMonth {
		return false
	}
	if m == 0 {
		return d == 0
	}
	ds := days[m]
	if m == 2 && y%4 == 0 {
		ds = 29
	}
	return d <= ds
}

// CheckPeriodoDiarioCompound valida que el período compuesto sea correcto.
// Descompone usando DecomposePeriodoDiario y delega las validaciones a CheckPeriodoDiario.
func CheckPeriodoDiarioCompound(v uint) bool {
	return CheckPeriodoDiario(DecomposePeriodoDiario(v))
}

// CheckPeriodoMensual valida que los componentes conformen un período mensual correcto.
//
// Valida que:
//
// - El año (y) esté dentro del rango [MinYear, MaxYear] (que por defecto es [1000,9999])
// - El mes (m) esté dentro del rango [0,12]
func CheckPeriodoMensual(y, m uint) bool {
	return y >= MinYear && y <= MaxYear && m >= minMonth && m <= maxMonth
}

// CheckPeriodoMensualCompound valida que el período compuesto sea correcto.
// Descompone usando DecomposePeriodoMensual y delega las validaciones a CheckPeriodoMensual.
func CheckPeriodoMensualCompound(v uint) bool {
	return CheckPeriodoMensual(DecomposePeriodoMensual(v))
}

// CheckPeriodoAnual valida que el período anual suministrado sea correcto.
//
// Valida que:
//
// - El año (y) esté dentro del rango [MinYear, MaxYear] (que por defecto es [1000,9999])
func CheckPeriodoAnual(y uint) bool {
	return y >= MinYear && y <= MaxYear
}
