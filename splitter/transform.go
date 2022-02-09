package main

import (
	"image"
	"math"

	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

// FourPointTransform does magic
func FourPointTransform(img gocv.Mat, pts *mat.Dense) gocv.Mat {

	rect := orderPoints(pts)
	tl := rect.RawRowView(0)
	tr := rect.RawRowView(1)
	br := rect.RawRowView(2)
	bl := rect.RawRowView(3)

	widthA := math.Sqrt(math.Pow((br[0]-bl[0]), 2) + math.Pow((br[1]-bl[1]), 2))
	widthB := math.Sqrt(math.Pow((tr[0]-tl[0]), 2) + math.Pow((tr[1]-tl[1]), 2))
	maxWidth := int(math.Max(widthA, widthB))

	heightA := math.Sqrt(math.Pow((tr[0]-br[0]), 2) + math.Pow((tr[1]-br[1]), 2))
	heightB := math.Sqrt(math.Pow((tl[0]-bl[0]), 2) + math.Pow((tl[1]-bl[1]), 2))
	maxHeight := int(math.Max(heightA, heightB))

	dst := mat.NewDense(4, 2, []float64{
		0, 0,
		(float64(maxWidth) - 1), 0,
		(float64(maxWidth) - 1), (float64(maxHeight) - 1),
		0, (float64(maxHeight) - 1),
	})

	M := gocv.GetPerspectiveTransform(gocv.NewPointVectorFromPoints(convertDenseToImagePoint(rect)), gocv.NewPointVectorFromPoints(convertDenseToImagePoint(dst)))
	gocv.WarpPerspective(img, &img, M, image.Point{X: maxWidth, Y: maxHeight})

	return img
}

func convertDenseToImagePoint(pts *mat.Dense) []image.Point {
	var sd []image.Point

	r, c := pts.Dims()
	if c != 2 {
		return sd
	}
	for i := 0; i < r; i++ {
		row := pts.RowView(i)
		sd = append(sd, image.Point{
			X: int(row.AtVec(0)),
			Y: int(row.AtVec(1)),
		})
	}
	return sd
}
func orderPoints(pts *mat.Dense) *mat.Dense {

	rect := mat.NewDense(4, 2, nil)

	sumMinIndex, sumMaxIndex := findMinMaxSumIndex(*pts)
	rect.SetRow(0, pts.RawRowView(sumMinIndex))
	rect.SetRow(2, pts.RawRowView(sumMaxIndex))

	diffMinIndex, diffMaxIndex := findMinMaxDiffIndex(*pts)
	rect.SetRow(1, pts.RawRowView(diffMinIndex))
	rect.SetRow(3, pts.RawRowView(diffMaxIndex))

	return rect
}
func findMinMaxSumIndex(pts mat.Dense) (int, int) {
	r, c := pts.Dims()

	maxIndex := 0
	maxValue := 0.0
	minIndex := 0
	minValue := 0.0

	for i := 0; i < r; i++ {
		row := pts.RowView(i)
		sum := 0.0
		for j := 0; j < c; j++ {
			sum += row.AtVec(j)
		}

		if i == 0 {
			maxValue = sum
			minValue = sum
		}

		if sum > maxValue {
			maxValue = sum
			maxIndex = i
		}
		if sum < minValue {
			minValue = sum
			minIndex = i
		}
	}
	return minIndex, maxIndex
}

func findMinMaxDiffIndex(pts mat.Dense) (int, int) {
	r, c := pts.Dims()

	maxIndex := 0
	maxValue := 0.0
	minIndex := 0
	minValue := 0.0

	for i := 0; i < r; i++ {
		row := pts.RowView(i)
		diff := row.AtVec(c - 1)
		for j := c - 2; j >= 0; j-- {
			diff -= row.AtVec(j)
		}

		if i == 0 {
			maxValue = diff
			minValue = diff
		}

		if diff > maxValue {
			maxValue = diff
			maxIndex = i
		}

		if diff < minValue {
			minValue = diff
			minIndex = i
		}
	}
	return minIndex, maxIndex
}
